// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package replication

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger-labs/orion-server/config"
	"github.com/hyperledger-labs/orion-server/internal/blockstore"
	"github.com/hyperledger-labs/orion-server/internal/comm"
	ierrors "github.com/hyperledger-labs/orion-server/internal/errors"
	"github.com/hyperledger-labs/orion-server/internal/httputils"
	"github.com/hyperledger-labs/orion-server/internal/queue"
	"github.com/hyperledger-labs/orion-server/pkg/logger"
	"github.com/hyperledger-labs/orion-server/pkg/types"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"
	"go.etcd.io/etcd/wal"
)

const (
	// DefaultSnapshotCatchUpEntries is the default number of entries
	// to preserve in memory when a snapshot is taken. This is for
	// slow followers to catch up.
	DefaultSnapshotCatchUpEntries = uint64(4)
)

type BlockLedgerReader interface {
	Height() (uint64, error)
	Get(blockNumber uint64) (*types.Block, error)
}

//go:generate counterfeiter -o mocks/pending_txs.go --fake-name PendingTxsReleaser . PendingTxsReleaser

type PendingTxsReleaser interface {
	ReleaseWithError(txIDs []string, err error)
}

type BlockReplicator struct {
	localConf       *config.LocalConfiguration
	joinBlock       *types.Block
	joinBlockNumber uint64

	proposeCh       chan *types.Block
	raftID          uint64
	raftStorage     *RaftStorage
	raftNode        raft.Node
	oneQueueBarrier *queue.OneQueueBarrier // Synchronizes the block-replication deliver with the block-processor commit
	transport       *comm.HTTPTransport
	ledgerReader    BlockLedgerReader
	pendingTxs      PendingTxsReleaser

	stopCh        chan struct{}
	stopOnce      sync.Once
	doneProposeCh chan struct{}
	doneEventCh   chan struct{}

	// shared state between the propose-loop go-routine and event-loop go-routine
	mutex                           sync.Mutex
	clusterConfig                   *types.ClusterConfig
	joinExistingCluster             bool
	lastKnownLeader                 uint64
	lastKnownLeaderHost             string // cache the leader's Node host:port for client request redirection
	cancelProposeContext            func() // cancels the propose-context if leadership is lost
	lastProposedBlockNumber         uint64
	lastProposedBlockHeaderBaseHash []byte
	lastCommittedBlock              *types.Block
	numInFlightBlocks               uint32 // number of in-flight blocks
	condTooManyInFlightBlocks       *sync.Cond

	appliedIndex uint64

	// needed by snapshotting
	sizeLimit        uint64 // SnapshotIntervalSize in bytes
	accDataSize      uint64 // accumulative data size since last snapshot
	lastSnapBlockNum uint64
	confState        raftpb.ConfState // Etcdraft requires ConfState to be persisted within snapshot

	lg *logger.SugarLogger
}

// Config holds the configuration information required to initialize the block replicator.
type Config struct {
	LocalConf            *config.LocalConfiguration
	ClusterConfig        *types.ClusterConfig
	JoinBlock            *types.Block
	LedgerReader         BlockLedgerReader
	Transport            *comm.HTTPTransport
	BlockOneQueueBarrier *queue.OneQueueBarrier
	PendingTxs           PendingTxsReleaser
	Logger               *logger.SugarLogger
}

// NewBlockReplicator creates a new BlockReplicator.
func NewBlockReplicator(conf *Config) (*BlockReplicator, error) {
	raftID, err := comm.MemberRaftID(conf.LocalConf.Server.Identity.ID, conf.ClusterConfig)
	if err != nil {
		return nil, err
	}

	lg := conf.Logger.With("nodeID", conf.LocalConf.Server.Identity.ID, "raftID", raftID)

	haveWAL := wal.Exist(conf.LocalConf.Replication.WALDir)
	storage, err := CreateStorage(lg, conf.LocalConf.Replication.WALDir, conf.LocalConf.Replication.SnapDir)
	if err != nil {
		return nil, errors.Errorf("failed to restore persisted raft data: %s", err)
	}
	storage.SnapshotCatchUpEntries = DefaultSnapshotCatchUpEntries

	var snapBlkNum uint64
	var confState raftpb.ConfState
	if s := storage.Snapshot(); !raft.IsEmptySnap(s) {
		snapBlock := &types.Block{}
		if err := proto.Unmarshal(s.Data, snapBlock); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal snapshot block")
		}

		snapBlkNum = snapBlock.GetHeader().GetBaseHeader().GetNumber()
		confState = s.Metadata.ConfState
		lg.Debugf("Starting from last snapshot: block number [%d], Raft ConfState: %+v", snapBlkNum, confState)
	}

	br := &BlockReplicator{
		localConf:            conf.LocalConf,
		joinBlock:            conf.JoinBlock,
		joinBlockNumber:      conf.JoinBlock.GetHeader().GetBaseHeader().GetNumber(), // if joinBlock==nil => 0
		proposeCh:            make(chan *types.Block, 1),
		raftID:               raftID,
		raftStorage:          storage,
		oneQueueBarrier:      conf.BlockOneQueueBarrier,
		stopCh:               make(chan struct{}),
		doneProposeCh:        make(chan struct{}),
		doneEventCh:          make(chan struct{}),
		clusterConfig:        conf.ClusterConfig,
		cancelProposeContext: func() {}, //NOOP
		transport:            conf.Transport,
		ledgerReader:         conf.LedgerReader,
		pendingTxs:           conf.PendingTxs,
		sizeLimit:            conf.ClusterConfig.ConsensusConfig.RaftConfig.SnapshotIntervalSize,
		confState:            confState,
		lastSnapBlockNum:     snapBlkNum,
		lg:                   lg,
	}
	br.condTooManyInFlightBlocks = sync.NewCond(&br.mutex)

	height, err := br.ledgerReader.Height()
	if err != nil {
		br.lg.Panicf("Failed to read block height: %s", err)
	}

	if height > 0 {
		br.lastCommittedBlock, err = br.ledgerReader.Get(height)
		if err != nil {
			br.lg.Panicf("Failed to read last block: %s", err)
		}
		br.lastProposedBlockNumber = br.lastCommittedBlock.GetHeader().GetBaseHeader().GetNumber()
		if baseHash, err := blockstore.ComputeBlockBaseHash(br.lastCommittedBlock); err == nil {
			br.lastProposedBlockHeaderBaseHash = baseHash
		} else {
			br.lg.Panicf("Failed to compute last block base hash: %s", err)
		}
	}

	if height > 1 {
		metadata := br.lastCommittedBlock.GetConsensusMetadata()
		br.appliedIndex = metadata.GetRaftIndex()
		br.lg.Debugf("last block [%d], consensus metadata: %+v", height, metadata)
	}

	//DO NOT use Applied option in config, we guard against replay of written blocks with `appliedIndex` instead.
	raftConfig := &raft.Config{
		ID:              raftID,
		ElectionTick:    int(br.clusterConfig.ConsensusConfig.RaftConfig.ElectionTicks),
		HeartbeatTick:   int(br.clusterConfig.ConsensusConfig.RaftConfig.HeartbeatTicks),
		MaxSizePerMsg:   br.localConf.BlockCreation.MaxBlockSize,
		MaxInflightMsgs: int(br.clusterConfig.ConsensusConfig.RaftConfig.MaxInflightBlocks),
		Logger:          lg,
		Storage:         br.raftStorage.MemoryStorage,
		// PreVote prevents reconnected node from disturbing network.
		// See etcd/raft doc for more details.
		PreVote:                   true,
		CheckQuorum:               true,
		DisableProposalForwarding: true, // This prevents blocks from being accidentally proposed by followers
	}

	lg.Debugf("haveWAL: %v, Storage: %v, Raft config: %+v", haveWAL, storage, raftConfig)

	// TODO support node join to an existing cluster: https://github.com/hyperledger-labs/orion-server/issues/260
	br.joinExistingCluster = (conf.JoinBlock != nil) && (height < br.joinBlockNumber)

	if haveWAL {
		br.raftNode = raft.RestartNode(raftConfig)
	} else {
		if br.joinExistingCluster {
			// TODO support node join to an existing cluster: https://github.com/hyperledger-labs/orion-server/issues/260
			return nil, errors.New("not supported yet: BlockReplicator joinExistingCluster")
		} else {
			startPeers := raftPeers(br.clusterConfig)
			br.raftNode = raft.StartNode(raftConfig, startPeers)
		}
	}

	return br, nil
}

func (br *BlockReplicator) RaftID() uint64 {
	return br.raftID
}

// Submit a block for replication.
//
// This call may block if the replication input queue is full.
// Returns an error if the current node is not a leader.
// Returns an error if the component is already closed.
func (br *BlockReplicator) Submit(block *types.Block) error {
	blockNum := block.GetHeader().GetBaseHeader().GetNumber()
	if err := br.IsLeader(); err != nil {
		br.lg.Debugf("Submit of block [%d] refused, not a leader: %s", blockNum, err)
		return err
	}

	select {
	case <-br.stopCh:
		return &ierrors.ClosedError{ErrMsg: "block replicator closed"}
	case br.proposeCh <- block:
		br.lg.Debugf("Submitted block [%d]", blockNum)
		return nil
	}
}

// Start an internal go-routine to serve the main replication loop.
func (br *BlockReplicator) Start() {
	readyRaftCh := make(chan struct{})
	go br.runRaftEventLoop(readyRaftCh)
	<-readyRaftCh

	readyProposeCh := make(chan struct{})
	go br.runProposeLoop(readyProposeCh)
	<-readyProposeCh
}

func (br *BlockReplicator) runRaftEventLoop(readyCh chan<- struct{}) {
	defer close(br.doneEventCh)

	br.lg.Info("Starting the block replicator event loop")
	close(readyCh)

	//If height is smaller than the block number in the last snapshot, it means the node stopped after a
	// snapshot trigger was received, but before catch-up was completed. In order to cover this case, we do
	// catch-up first.
	if lastSnapshot := br.raftStorage.Snapshot(); !raft.IsEmptySnap(lastSnapshot) {
		if err := br.catchUp(lastSnapshot); err != nil {
			br.lg.Panicf("Failed to catch-up to last snapshot: %+v", lastSnapshot)
		}
	}

	// TODO use 'clock.Clock' so that tests can inject a fake clock
	tickInterval, err := time.ParseDuration(br.clusterConfig.ConsensusConfig.RaftConfig.TickInterval)
	if err != nil {
		br.lg.Panicf("Error parsing raft tick interval duration: %s", err)
	}
	raftTicker := time.NewTicker(tickInterval)
	electionTimeout := tickInterval.Seconds() * float64(br.clusterConfig.ConsensusConfig.RaftConfig.ElectionTicks)
	halfElectionTimeout := electionTimeout / 2

	// TODO proactive campaign to speed up leader election on a new cluster

	var raftStatusStr string
Event_Loop:
	for {
		select {
		case <-raftTicker.C:
			if status := br.raftNode.Status().String(); status != raftStatusStr {
				br.lg.Debugf("Raft node status: %+v", status)
				raftStatusStr = status
			}
			br.raftNode.Tick()

		case rd := <-br.raftNode.Ready():
			startStoring := time.Now()
			if err := br.raftStorage.Store(rd.Entries, rd.HardState, rd.Snapshot); err != nil {
				br.lg.Panicf("Failed to persist etcd/raft data: %s", err)
			}
			duration := time.Since(startStoring).Seconds()
			if duration > halfElectionTimeout {
				br.lg.Warningf("WAL sync took %v seconds and the network is configured to start elections after %v seconds. Your disk is too slow and may cause loss of quorum and trigger leadership election.", duration, electionTimeout)
			}

			if !raft.IsEmptySnap(rd.Snapshot) {
				if err := br.catchUp(rd.Snapshot); err != nil {
					br.lg.Panicf("Failed to catch-up to snapshot: %+v", rd.Snapshot)
				}
			}

			br.transport.SendConsensus(rd.Messages)

			if ok := br.deliverEntries(rd.CommittedEntries); !ok {
				br.lg.Warningf("Stopping to deliver committed entries, breaking out of event loop")
				break Event_Loop
			}

			// update last known leader
			if rd.SoftState != nil {
				leader := atomic.LoadUint64(&rd.SoftState.Lead) // etcdraft requires atomic access to this var
				if leader != raft.None {
					br.lg.Debugf("Leader %d is present", leader)
				} else {
					br.lg.Debug("No leader")
				}

				br.processLeaderChanges(leader)
			}

			br.raftNode.Advance()

		case <-br.stopCh:
			br.lg.Info("Stopping block replicator")
			break Event_Loop
		}
	}

	// Notify the propose-loop go-routine in case it is waiting for blocks to commit or a leadership change.
	br.mutex.Lock()
	br.lastKnownLeader = 0      // stop proposing
	br.lastKnownLeaderHost = "" // stop proposing
	br.numInFlightBlocks = 0    // stop waiting for blocks to commit
	br.condTooManyInFlightBlocks.Broadcast()
	br.mutex.Unlock()

	raftTicker.Stop()
	br.raftNode.Stop()
	if err := br.raftStorage.Close(); err != nil {
		br.lg.Errorf("Error while stopping RaftStorage: %s", err) // TODO move to raft main loop
	}

	br.lg.Info("Exiting block replicator event loop")
}

func (br *BlockReplicator) processLeaderChanges(leader uint64) {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	if leader != br.lastKnownLeader {
		br.lg.Infof("Leader changed: %d to %d", br.lastKnownLeader, leader)

		lostLeadership := br.lastKnownLeader == br.raftID
		assumedLeadership := leader == br.raftID

		if lostLeadership {
			br.lg.Info("Lost leadership")
			// cancel the current proposal to free the propose-loop go-routine, as it might block for a long time.
			br.cancelProposeContext()
			br.cancelProposeContext = func() {} // NOOP
		} else if assumedLeadership {
			br.lg.Info("Assumed leadership")
		}

		if lostLeadership || assumedLeadership {
			var err error
			br.lastProposedBlockNumber = br.lastCommittedBlock.GetHeader().GetBaseHeader().GetNumber()
			br.lastProposedBlockHeaderBaseHash, err = blockstore.ComputeBlockBaseHash(br.lastCommittedBlock)
			if err != nil {
				br.lg.Panicf("Error computing base header hash of last commited block: %+v; error: %s",
					br.lastCommittedBlock.GetHeader(), err)
			}
			br.numInFlightBlocks = 0
			br.condTooManyInFlightBlocks.Broadcast()
		}

		br.lastKnownLeader = leader
		br.lastKnownLeaderHost = br.nodeHostPortFromRaftID(leader)
	}
}

// When a node lags behind the cluster more than the last checkpoint of the leader, the leader will send a snapshot to
// it. A snapshot is a block with some raft information. A received snapshot serves as a trigger for the node to
// perform catch-up, or state transfer. It will contact one of the active members of the cluster (preferably the
// leader), and will request the missing blocks up to the block indicated by the snapshot.
func (br *BlockReplicator) catchUp(snap raftpb.Snapshot) error {
	if snap.Metadata.Index <= br.appliedIndex {
		br.lg.Debugf("Skip snapshot taken at index %d, because it is behind current applied index %d", snap.Metadata.Index, br.appliedIndex)
		return nil
	}

	var snapBlock = &types.Block{}
	if err := proto.Unmarshal(snap.Data, snapBlock); err != nil {
		return errors.Errorf("failed to unmarshal snapshot data to block: %s", err)
	}

	initBlockNumber := br.getLastCommittedBlockNumber()
	br.lg.Debugf("initial last block number: %+v", initBlockNumber)
	br.lg.Debugf("snap block: %+v", snapBlock)

	if initBlockNumber >= snapBlock.Header.BaseHeader.Number {
		br.lg.Errorf("Snapshot is at block [%d], local block number is %d, no catch-up needed", snapBlock.Header.BaseHeader.Number, initBlockNumber)
		return nil
	}

	br.lg.Infof("Starting state transfer; From block: %d, index %d; To block: %d, index: %d",
		initBlockNumber, br.appliedIndex, snapBlock.Header.BaseHeader.Number, snap.Metadata.Index)
	br.confState = snap.Metadata.ConfState
	br.appliedIndex = snap.Metadata.Index

	// Pull the missing blocks, starting with one past the last block we have, and ending with the block number from the snapshot.
	for nextBlockNumber := initBlockNumber + 1; nextBlockNumber <= snapBlock.Header.BaseHeader.Number; {
		var blocks []*types.Block
		var err error
		blocksReadyCh := make(chan struct{})
		ctx, cancel := context.WithCancel(context.Background())

		//Try to pull some blocks in a go-routine so that we may cancel it if the server shuts down.
		//Note that `PullBlocks` will not necessarily return all the blocks we requested, hence the enclosing loop.
		go func() {
			defer close(blocksReadyCh)
			blocks, err = br.transport.PullBlocks(ctx, nextBlockNumber, snapBlock.Header.BaseHeader.Number, br.GetLeaderID())
		}()

		select {
		case <-br.stopCh:
			cancel()
			<-blocksReadyCh
			return &ierrors.ClosedError{ErrMsg: "server stopped during catch-up"}
		case <-blocksReadyCh:
			if err != nil {
				lastBlockNumber := br.getLastCommittedBlockNumber()
				switch err.(type) {
				case *ierrors.ClosedError:

					br.lg.Warnf("closing, stopping to pull blocks from cluster; last block number [%d], snapshot: %+v", lastBlockNumber, snap)
					return nil
				default:
					return errors.Wrapf(err, "failed to pull blocks from cluster; last block number [%d], snapshot: %+v", lastBlockNumber, snap)
				}
			}

			br.lg.Infof("Going to commit [%d] blocks", len(blocks)) //Not necessarily the entire range requested!

			for _, blockToCommit := range blocks {
				br.lg.Infof("enqueue for commit block [%d], ConsensusMetadata: [%+v]",
					blockToCommit.GetHeader().GetBaseHeader().GetNumber(),
					blockToCommit.GetConsensusMetadata())

				if err := br.commitBlock(blockToCommit); err != nil {
					lastBlockNumber := br.getLastCommittedBlockNumber()
					switch err.(type) {
					case *ierrors.ClosedError:
						br.lg.Warnf("closing, stopping to pull blocks from cluster; last block number [%d], snapshot: %+v", lastBlockNumber, snap)
						return nil
					default:
						return err
					}
				}

				nextBlockNumber++
			}
		}
	}

	lastBlockNumber := br.getLastCommittedBlockNumber()
	br.lg.Infof("Finished syncing with cluster up to and including block [%d]", lastBlockNumber)

	return nil

}

func (br *BlockReplicator) deliverEntries(committedEntries []raftpb.Entry) bool {
	br.lg.Debugf("Num. entries: %d", len(committedEntries))
	if len(committedEntries) == 0 {
		return true
	}

	if committedEntries[0].Index > br.appliedIndex+1 {
		br.lg.Panicf("First index of committed entry [%d] should <= appliedIndex [%d]+1", committedEntries[0].Index, br.appliedIndex)
	}

	var position int
	for i := range committedEntries {
		br.lg.Debugf("processing commited entry [%d]: %s", i, raftEntryString(committedEntries[i]))

		switch committedEntries[i].Type {
		case raftpb.EntryNormal:
			if len(committedEntries[i].Data) == 0 {
				br.lg.Debugf("commited entry [%d] has empty data, ignoring", i)
				break
			}

			position = i
			br.accDataSize += uint64(len(committedEntries[i].Data))

			// We need to strictly avoid re-applying normal entries,
			// otherwise we are writing the same block twice.
			if committedEntries[i].Index <= br.appliedIndex {
				br.lg.Debugf("Received block with raft index (%d) <= applied index (%d), skip", committedEntries[i].Index, br.appliedIndex)
				break
			}

			var block = &types.Block{}
			if err := proto.Unmarshal(committedEntries[i].Data, block); err != nil {
				br.lg.Panicf("Error unmarshaling entry [#%d], entry: %+v, error: %s", i, committedEntries[i], err)
			}
			block.ConsensusMetadata = &types.ConsensusMetadata{
				RaftTerm:  committedEntries[i].Term,
				RaftIndex: committedEntries[i].Index,
			}

			err := br.commitBlock(block)
			if err != nil {
				br.lg.Errorf("commit block error: %s, stopping block replicator", err.Error())
				return false
			}

		case raftpb.EntryConfChange:
			// For re-config, we use the V2 API, with `raftpb.EntryConfChangeV2` messages. These message happen only
			// when the raft node bootstraps, i.e. when `raft.StartNode` is called.
			var cc raftpb.ConfChange
			if err := cc.Unmarshal(committedEntries[i].Data); err != nil {
				br.lg.Warnf("Failed to unmarshal ConfChange data: %s", err)
				continue
			}

			br.confState = *br.raftNode.ApplyConfChange(cc)

			switch cc.Type {
			case raftpb.ConfChangeAddNode:
				br.lg.Infof("Applied config change to add node %d, current nodes in cluster: %+v", cc.NodeID, br.confState.Voters)
			case raftpb.ConfChangeRemoveNode:
				br.lg.Infof("Applied config change to remove node %d, current nodes in cluster: %+v", cc.NodeID, br.confState.Voters)
			default:
				br.lg.Panic("Programming error, encountered unsupported raft config change")
			}

		case raftpb.EntryConfChangeV2:
			var ccV2 raftpb.ConfChangeV2
			if err := ccV2.Unmarshal(committedEntries[i].Data); err != nil {
				br.lg.Warnf("Failed to unmarshal ConfChangeV2 data: %s", err)
				continue
			}

			if ccV2.Context != nil {
				var block = &types.Block{}
				if err := proto.Unmarshal(ccV2.Context, block); err != nil {
					br.lg.Panicf("Error unmarshaling entry [#%d], entry: %+v, error: %s", i, committedEntries[i], err)
				}

				block.ConsensusMetadata = &types.ConsensusMetadata{
					RaftTerm:  committedEntries[i].Term,
					RaftIndex: committedEntries[i].Index,
				}

				err := br.commitBlock(block) // transport is reconfigured within after the block commits.
				if err != nil {
					br.lg.Errorf("commit block error: %s, stopping block replicator", err.Error())
					return false
				}
			}

			br.confState = *br.raftNode.ApplyConfChange(ccV2)
			br.lg.Infof("Applied config changes: %+v, current nodes in cluster: %+v", ccV2.Changes, br.confState.Voters)

			// TODO detect removal of leader?

			// detect removal of self
			removalOfSelf := true
			for _, id := range br.confState.Voters {
				if id == br.raftID {
					removalOfSelf = false
					break
				}
			}
			if removalOfSelf {
				br.lg.Warning("This node was removed from the cluster, replication is shutting down")
				return false
			}
		}

		// after commit, update appliedIndex
		if br.appliedIndex < committedEntries[i].Index {
			br.appliedIndex = committedEntries[i].Index
		}
	}

	// Take a snapshot if in-memory storage size exceeds the limit
	if br.accDataSize >= br.sizeLimit {
		var snapBlock = &types.Block{}
		if err := proto.Unmarshal(committedEntries[position].Data, snapBlock); err != nil {
			br.lg.Panicf("Error unmarshaling entry [#%d], entry: %+v, error: %s", position, committedEntries[position], err)
		}

		if err := br.raftStorage.TakeSnapshot(br.appliedIndex, br.confState, committedEntries[position].Data); err != nil {
			br.lg.Fatalf("Failed to create snapshot at index %d: %s", br.appliedIndex, err)
		}

		br.lg.Infof("Accumulated %d bytes since last snapshot, exceeding size limit (%d bytes), "+
			"taking snapshot at block [%d] (index: %d), last snapshotted block number is %d, current voters: %+v",
			br.accDataSize, br.sizeLimit, snapBlock.GetHeader().GetBaseHeader().GetNumber(), br.appliedIndex, br.lastSnapBlockNum, br.confState.Voters)

		br.accDataSize = 0
		br.lastSnapBlockNum = snapBlock.GetHeader().GetBaseHeader().GetNumber()
	}

	return true
}

func (br *BlockReplicator) runProposeLoop(readyCh chan<- struct{}) {
	defer close(br.doneProposeCh)

	br.lg.Info("Starting the block replicator propose loop")
	close(readyCh)

Propose_Loop:
	for {
		select {
		case blockToPropose := <-br.proposeCh:
			var addedPeers, removedPeers []*types.PeerConfig
			var isMembershipConfig bool

			if httputils.IsConfigBlock(blockToPropose) {
				// TODO verify the config-tx by itself and vs. current config
				newClusterConfig := blockToPropose.Payload.(*types.Block_ConfigTxEnvelope).ConfigTxEnvelope.GetPayload().GetNewConfig()
				_, consensus, _, _ := ClassifyClusterReConfig(br.clusterConfig, newClusterConfig)
				if consensus {
					var errDetect error
					addedPeers, removedPeers, _, errDetect = detectPeerConfigChanges(br.clusterConfig.ConsensusConfig, newClusterConfig.ConsensusConfig)
					if errDetect != nil {
						br.releasePendingTXs(blockToPropose, "Declined to propose block", errDetect)
						continue Propose_Loop
					}
					isMembershipConfig = len(addedPeers)+len(removedPeers) > 0
				}
			}

			if !isMembershipConfig {
				// Data-Tx, User-admin-Tx, DB-admin-Tx, as well as Config-Tx that do not membership changes
				if ok := br.proposeRegular(blockToPropose); !ok {
					continue Propose_Loop
				}
			} else {
				//  Config-Tx: with consensus membership change
				if ok := br.proposeMembershipConfigChange(blockToPropose, addedPeers, removedPeers); !ok {
					continue Propose_Loop
				}
			}

			br.updateLastProposal(blockToPropose)

		case <-br.stopCh:
			br.lg.Debug("Stopping block replicator")
			break Propose_Loop
		}
	}

	br.lg.Info("Exiting the block replicator propose loop")
}

// proposeRegular proposes the block to Raft as a regular message.
func (br *BlockReplicator) proposeRegular(blockToPropose *types.Block) bool {
	ctx, blockBytes, doPropose := br.prepareProposal(blockToPropose)
	if !doPropose {
		return false
	}

	// Propose to raft: the call to raft.Node.Propose() may block when a leader loses its leadership and has no quorum.
	// It is cancelled when the node loses leadership, by the event-loop go-routine.
	err := br.raftNode.Propose(ctx, blockBytes)
	if err != nil {
		br.releasePendingTXs(blockToPropose, "Failed to propose block", err)
		return false
	}

	return true
}

// proposeMembershipConfigChange propose membership config changes, that is, adding or removing a peer.
// This is proposed in a 'raftpb.ConfChangeV2' message using the 'ProposeConfChange' API. This call consents on the
// Raft membership change as well as on the config block, which is given as the 'ConfChangeV2.Context'.
// The return value signals whether the proposal completed correctly.
func (br *BlockReplicator) proposeMembershipConfigChange(blockToPropose *types.Block, addedPeers, removedPeers []*types.PeerConfig) bool {
	//  consensus membership config change
	ctx, blockBytes, doPropose := br.prepareProposal(blockToPropose)
	if !doPropose {
		return false
	}

	ccV2 := &raftpb.ConfChangeV2{
		Transition: raftpb.ConfChangeTransitionAuto,
		Context:    blockBytes,
	}
	for _, peer := range addedPeers {
		ccV2.Changes = append(ccV2.Changes, raftpb.ConfChangeSingle{
			Type:   raftpb.ConfChangeAddNode,
			NodeID: peer.RaftId,
		})
	}
	for _, peer := range removedPeers {
		ccV2.Changes = append(ccV2.Changes, raftpb.ConfChangeSingle{
			Type:   raftpb.ConfChangeRemoveNode,
			NodeID: peer.RaftId,
		})
	}

	br.lg.Infof("Going to propose membership config changes: %v", ccV2.Changes)
	// ProposeConfChange to raft: the call to raft.Node.ProposeConfChange() may block when a leader loses
	// its leadership and has no quorum.
	// It is cancelled when the node loses leadership, by the event-loop go-routine.
	err := br.raftNode.ProposeConfChange(ctx, ccV2)
	if err != nil {
		br.releasePendingTXs(blockToPropose, "Failed to propose block", err)
		return false
	}

	return true
}

func (br *BlockReplicator) releasePendingTXs(blockToPropose *types.Block, reasonMsg string, reasonErr error) {
	br.lg.Infof("%s: %+v; because: %s", reasonMsg, blockToPropose.GetHeader(), reasonErr)
	if txIDs, err := httputils.BlockPayloadToTxIDs(blockToPropose.GetPayload()); err == nil {
		br.pendingTxs.ReleaseWithError(txIDs, err)
	} else {
		br.lg.Errorf("Failed to extract TxIDs from block, dropping block: %v; error: %s", blockToPropose.GetHeader(), err)
	}
}

// prepareProposal Prepares the Raft proposal context and bytes, and determine whether to propose (only the leader can
// propose). This also numbers the block and sets the base header hash.
func (br *BlockReplicator) prepareProposal(blockToPropose *types.Block) (ctx context.Context, blockBytes []byte, doPropose bool) {
	br.mutex.Lock()

	if errLeader := br.isLeader(); errLeader != nil {
		br.mutex.Unlock() //do not call the pendingTxs component with a mutex locked

		br.releasePendingTXs(blockToPropose, "Declined to propose block", errLeader)

		return nil, nil, false //skip proposing
	}

	// number the block and set the base header hash
	br.insertBlockBaseHeader(blockToPropose)

	var err error
	blockBytes, err = proto.Marshal(blockToPropose)
	if err != nil {
		br.lg.Panicf("Error marshaling a block: %s", err)
	}

	ctx, br.cancelProposeContext = context.WithCancel(context.Background())

	br.mutex.Unlock()

	return ctx, blockBytes, true
}

// updateLastProposal updates the last block proposed in order to keep track of block numbering.
// TODO also keep track of proposed config, and prevent additional config proposals when one is already in flight.
func (br *BlockReplicator) updateLastProposal(lastBlockProposed *types.Block) {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	if br.isLeader() == nil {
		br.lastProposedBlockNumber = lastBlockProposed.GetHeader().GetBaseHeader().GetNumber()
		if baseHash, err := blockstore.ComputeBlockBaseHash(lastBlockProposed); err == nil {
			br.lastProposedBlockHeaderBaseHash = baseHash
		} else {
			br.lg.Panicf("Failed to compute last block base hash: %s", err)
		}
		br.numInFlightBlocks++

		if br.numInFlightBlocks > br.clusterConfig.ConsensusConfig.RaftConfig.MaxInflightBlocks {
			br.lg.Debugf("Number of in-flight blocks exceeds max, %d > %d, waiting for blocks to commit", //Tested side effect
				br.numInFlightBlocks, br.clusterConfig.ConsensusConfig.RaftConfig.MaxInflightBlocks)

			for br.numInFlightBlocks > br.clusterConfig.ConsensusConfig.RaftConfig.MaxInflightBlocks {
				// the go-routine will be notified by the event-loop go-routine when:
				// - a block commits, or
				// - when leadership is lost or assumed, or
				// - when the event-loop go-routine exits. This is done in order to remain
				//   reactive to server shutdown while waiting for blocks to commit.
				br.condTooManyInFlightBlocks.Wait()
			}
			br.lg.Debugf("Number of in-flight blocks back to normal: %d", br.numInFlightBlocks)
		}
	}
}

// Close signals the internal go-routine to stop and waits for it to exit.
// If the component is already closed, and error is returned.
func (br *BlockReplicator) Close() (err error) {
	err = &ierrors.ClosedError{ErrMsg: "block replicator already closed"}
	br.stopOnce.Do(func() {
		br.lg.Info("closing block replicator")
		close(br.stopCh)
		if errQB := br.oneQueueBarrier.Close(); errQB != nil {
			br.lg.Debugf("OneQueueBarrier error: %s", errQB)
		}
		<-br.doneProposeCh
		<-br.doneEventCh

		//after the node stops, it no longer knows who the leader is
		br.mutex.Lock()
		defer br.mutex.Unlock()
		br.lastKnownLeader = 0
		br.lastKnownLeaderHost = ""

		err = nil
	})

	return err
}

func (br *BlockReplicator) IsLeader() *ierrors.NotLeaderError {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	return br.isLeader()
}

func (br *BlockReplicator) isLeader() *ierrors.NotLeaderError {
	if br.lastKnownLeader == br.raftID {
		return nil
	}

	return &ierrors.NotLeaderError{
		LeaderID: br.lastKnownLeader, LeaderHostPort: br.lastKnownLeaderHost}
}

func (br *BlockReplicator) GetLeaderID() uint64 {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	return br.lastKnownLeader
}

func (br *BlockReplicator) commitBlock(block *types.Block) error {
	br.lg.Infof("Enqueue for commit block [%d], ConsensusMetadata: %+v ",
		block.GetHeader().GetBaseHeader().GetNumber(),
		block.GetConsensusMetadata())

	// we can only get a valid config transaction
	reConfig, err := br.oneQueueBarrier.EnqueueWait(block)
	if err != nil {
		return err
	}

	br.setLastCommittedBlock(block)

	if reConfig == nil {
		return nil
	}

	clusterConfig := reConfig.(*types.ClusterConfig)
	if err := br.updateClusterConfig(clusterConfig); err != nil {
		br.lg.Panicf("Failed to update to ClusterConfig during commitBlock: error: %s", err)
	}

	return nil
}

func (br *BlockReplicator) setLastCommittedBlock(block *types.Block) {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	br.lastCommittedBlock = block
	if br.numInFlightBlocks > 0 { // only reduce on the leader
		br.numInFlightBlocks--
		br.condTooManyInFlightBlocks.Broadcast()
	}
}

func (br *BlockReplicator) getLastCommittedBlockNumber() uint64 {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	return br.lastCommittedBlock.GetHeader().GetBaseHeader().GetNumber()
}

func (br *BlockReplicator) updateClusterConfig(clusterConfig *types.ClusterConfig) error {
	br.lg.Infof("New cluster config committed, going to apply to block replicator: %+v", clusterConfig)

	br.mutex.Lock()
	defer br.mutex.Unlock()

	// We must not receive any errors here, because that would render the config-tx invalid, and invalid config txs
	// are not transferred to the BlockReplicator.
	addedPeers, removedPeers, changedPeers, err := detectPeerConfigChanges(br.clusterConfig.ConsensusConfig, clusterConfig.ConsensusConfig)
	if err != nil {
		return errors.Wrap(err, "failed to detect peer config changes")
	}

	err = br.transport.UpdatePeers(addedPeers, removedPeers, changedPeers, clusterConfig)
	if err != nil {
		return errors.Wrap(err, "failed to update peers on transport")
	}

	br.clusterConfig = clusterConfig

	return nil
}

func (br *BlockReplicator) nodeHostPortFromRaftID(raftID uint64) string {
	if raftID == 0 {
		return ""
	}

	var nodeID string
	for _, p := range br.clusterConfig.ConsensusConfig.Members {
		if p.RaftId == raftID {
			nodeID = p.NodeId
			break
		}
	}

	if nodeID == "" {
		br.lg.Warnf("not found: no member with RaftID: %d", raftID)
		return ""
	}

	for _, n := range br.clusterConfig.Nodes {
		if n.Id == nodeID {
			hostPort := fmt.Sprintf("%s:%d", n.Address, n.Port)
			return hostPort
		}
	}

	br.lg.Warnf("not found: no node with NodeID: %s, RaftID: %d", nodeID, raftID)
	return ""
}

func (br *BlockReplicator) Process(ctx context.Context, m raftpb.Message) error {
	br.lg.Debugf("Incoming raft message: %+v", m)
	//TODO look into the cluster config and check Members and reject messages from removed members

	err := br.raftNode.Step(ctx, m)
	if err != nil {
		br.lg.Errorf("Error during raft node Step: %s", err)
	}
	return err
}

func (br *BlockReplicator) IsIDRemoved(id uint64) bool {
	br.lg.Debugf("> IsIDRemoved: %d", id)
	// see: rafthttp.RAFT

	//TODO look into the cluster config and check whether this RaftID was removed.
	// removed RaftIDs may never return.
	// see issue: https://github.com/ibm-blockchain/bcdb-server/issues/40
	return false
}
func (br *BlockReplicator) ReportUnreachable(id uint64) {
	br.lg.Debugf("ReportUnreachable: %d", id)
	br.raftNode.ReportUnreachable(id)
}

func (br *BlockReplicator) ReportSnapshot(id uint64, status raft.SnapshotStatus) {
	br.lg.Debugf("> ReportSnapshot: %d, %+v", id, status)
	// see: rafthttp.RAFT
	//TODO see issue: https://github.com/ibm-blockchain/bcdb-server/issues/41
}

// called inside a br.mutex.Lock()
func (br *BlockReplicator) insertBlockBaseHeader(proposedBlock *types.Block) {
	blockNum := br.lastProposedBlockNumber + 1
	baseHeader := &types.BlockHeaderBase{
		Number:                 blockNum,
		PreviousBaseHeaderHash: br.lastProposedBlockHeaderBaseHash,
	}

	if blockNum > 1 {
		lastCommittedBlockNum := br.lastCommittedBlock.GetHeader().GetBaseHeader().GetNumber()
		lastCommittedBlockHash, err := blockstore.ComputeBlockHash(br.lastCommittedBlock)
		if err != nil {
			br.lg.Panicf("Error while creating block header for proposed block: %d; possible problems at last commited block header: %+v; error: %s",
				blockNum, br.lastCommittedBlock.GetHeader(), err)
		}
		baseHeader.LastCommittedBlockHash = lastCommittedBlockHash
		baseHeader.LastCommittedBlockNum = lastCommittedBlockNum
	}

	proposedBlock.Header = &types.BlockHeader{BaseHeader: baseHeader}
}
