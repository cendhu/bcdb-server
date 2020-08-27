package blockstore

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	proto "github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.ibm.com/blockchaindb/server/pkg/fileops"
)

var (
	// blocks are stored in an append-only file. As the
	// the file size could grow significantly in a longer
	// run, we use file chunks so that it would be easy
	// to archive chunks to free some storage space
	chunkPrefix    = "chunk_"
	chunkSizeLimit = int64(64 * 1024 * 1024)

	// block file chunks are stored inside fileChunksDir
	// while the index to the block file's offset to fetch
	// a given block number is stored inside blockIndexDir
	fileChunksDirName = "filechunks"
	blockIndexDirName = "blockindex"

	// underCreationFlag is used to mark that the store
	// is being created. If a failure happens during the
	// creation, the retry logic will use this file to
	// detect the partially created store and do cleanup
	// before creating a new store
	underCreationFlag = "undercreation"
)

// Store maintains a chain of blocks in an append-only
// filesystem
type Store struct {
	fileChunksDirPath     string
	currentFileChunk      *os.File
	currentOffset         int64
	currentChunkNum       uint64
	lastCommittedBlockNum uint64
	blockIndexDirPath     string
	blockIndex            *leveldb.DB
	reusableBuffer        []byte
	mu                    sync.RWMutex
}

// Open opens the store to maintains a chain of blocks
func Open(storeDir string) (*Store, error) {
	exist, err := fileops.Exists(storeDir)
	if err != nil {
		return nil, err
	}
	if !exist {
		return openNewStore(storeDir)
	}

	partialStoreExist, err := isExistingStoreCreatedPartially(storeDir)
	if err != nil {
		return nil, err
	}

	switch {
	case partialStoreExist:
		if err := fileops.RemoveAll(storeDir); err != nil {
			return nil, errors.Wrap(err, "error while removing the existing partially created store")
		}

		return openNewStore(storeDir)
	default:
		return openExistingStore(storeDir)
	}
}

func isExistingStoreCreatedPartially(storeDir string) (bool, error) {
	empty, err := fileops.IsDirEmpty(storeDir)
	if err != nil || empty {
		return true, err
	}

	return fileops.Exists(filepath.Join(storeDir, underCreationFlag))
}

func openNewStore(storeDir string) (*Store, error) {
	if err := fileops.CreateDir(storeDir); err != nil {
		return nil, errors.WithMessagef(err, "error while creating directory [%s]", storeDir)
	}

	underCreationFlagPath := filepath.Join(storeDir, underCreationFlag)
	if err := fileops.CreateFile(underCreationFlagPath); err != nil {
		return nil, err
	}

	fileChunksDirPath := filepath.Join(storeDir, fileChunksDirName)
	blockIndexDirPath := filepath.Join(storeDir, blockIndexDirName)

	for _, d := range []string{fileChunksDirPath, blockIndexDirPath} {
		if err := fileops.CreateDir(d); err != nil {
			return nil, errors.WithMessagef(err, "error while creating directory [%s]", d)
		}
	}

	file, err := openFileChunk(fileChunksDirPath, 0)
	if err != nil {
		return nil, err
	}

	indexDB, err := leveldb.OpenFile(blockIndexDirPath, &opt.Options{})
	if err != nil {
		return nil, errors.WithMessage(err, "error while creating an index database")
	}

	if err := fileops.Remove(underCreationFlagPath); err != nil {
		return nil, errors.WithMessagef(err, "error while removing the under creation flag [%s]", underCreationFlagPath)
	}

	return &Store{
		fileChunksDirPath:     fileChunksDirPath,
		currentFileChunk:      file,
		currentOffset:         0,
		currentChunkNum:       0,
		lastCommittedBlockNum: 0,
		blockIndexDirPath:     blockIndexDirPath,
		blockIndex:            indexDB,
		reusableBuffer:        make([]byte, binary.MaxVarintLen64),
	}, nil
}

func openExistingStore(storeDir string) (*Store, error) {
	fileChunksDirPath := filepath.Join(storeDir, fileChunksDirName)
	blockIndexDirPath := filepath.Join(storeDir, blockIndexDirName)

	currentFileChunk, currentChunkNum, err := findAndOpenLastFileChunk(fileChunksDirPath)
	if err != nil {
		return nil, err
	}

	chunkFileInfo, err := currentFileChunk.Stat()
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting the metadata of file [%s]", currentFileChunk.Name())
	}

	indexDB, err := leveldb.OpenFile(blockIndexDirPath, &opt.Options{})
	if err != nil {
		return nil, errors.WithMessagef(err, "error while opening leveldb file for index")
	}

	s := &Store{
		fileChunksDirPath: fileChunksDirPath,
		currentFileChunk:  currentFileChunk,
		currentOffset:     chunkFileInfo.Size(),
		currentChunkNum:   currentChunkNum,
		blockIndexDirPath: blockIndexDirPath,
		blockIndex:        indexDB,
		reusableBuffer:    make([]byte, binary.MaxVarintLen64),
	}
	return s, s.recoverIfNeeded()
}

func (s *Store) recoverIfNeeded() error {
	// TODO:
	// if there was a failure during the last block commit, the index may not be
	// upto date. Even the last block commit may be written partially to the file.
	// We need to add the logic to recover the store here.

	lastBlockNumberInIndex, lastBlockLocation, err := s.getLastBlockLocationInIndex()
	if err != nil {
		return err
	}

	s.lastCommittedBlockNum = lastBlockNumberInIndex
	_ = lastBlockLocation

	return nil
}

func (s *Store) getLastBlockLocationInIndex() (uint64, *BlockLocation, error) {
	itr := s.blockIndex.NewIterator(&util.Range{}, &opt.ReadOptions{})
	if err := itr.Error(); err != nil {
		return 0, nil, errors.Wrap(err, "error while finding the last committed block number in the index")
	}
	if !itr.Last() {
		return 0, nil, nil
	}

	key := itr.Key()
	val := itr.Value()

	blockNumber, _, err := decodeOrderPreservingVarUint64(key)
	if err != nil {
		return 0, nil, errors.Wrap(err, "error while decoding the last block index key")
	}

	blockLocation := &BlockLocation{}
	if err := proto.Unmarshal(val, blockLocation); err != nil {
		return 0, nil, errors.Wrap(err, "error while unmarshaling block location")
	}

	return blockNumber, blockLocation, nil
}

// Close closes the store
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.currentFileChunk.Close(); err != nil {
		return errors.WithMessage(err, "error while closing the store")
	}

	if err := s.blockIndex.Close(); err != nil {
		return errors.WithMessage(err, "error while closing the block index database")
	}

	return nil
}

func openFileChunk(dir string, chunkNum uint64) (*os.File, error) {
	path := constructBlockFileChunkPath(dir, chunkNum)
	file, err := fileops.OpenFile(path, 0644)
	if err != nil {
		return nil, errors.WithMessagef(err, "error while opening the file chunk")
	}

	return file, nil
}

func constructBlockFileChunkPath(dir string, chunkNum uint64) string {
	chunkName := fmt.Sprintf("%s%d", chunkPrefix, chunkNum)
	return filepath.Join(dir, chunkName)
}

func findAndOpenLastFileChunk(fileChunksDirPath string) (*os.File, uint64, error) {
	files, err := ioutil.ReadDir(fileChunksDirPath)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "error while listing file chunks in [%s]", fileChunksDirPath)
	}

	lastChunkNum := uint64(len(files) - 1)
	lastFileChunk, err := openFileChunk(fileChunksDirPath, lastChunkNum)
	if err != nil {
		return nil, 0, err
	}

	return lastFileChunk, lastChunkNum, nil
}