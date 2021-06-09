package stateindex

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

// encodeOrderPreservingVarUint64 returns a byte-representation for a uint64 number such that
// all zero-bits starting bytes are trimmed in order to reduce the length of the array
// For preserving the order in a default bytes-comparison, first byte contains the number of remaining bytes.
// The presence of first byte also allows to use the returned bytes as part of other larger byte array such as a
// composite-key representation in db
func encodeOrderPreservingVarUint64(number uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, number)
	startingIndex := 0
	size := 0
	for i, b := range bytes {
		if b != 0x00 {
			startingIndex = i
			size = 8 - i
			break
		}
	}
	sizeBytes := proto.EncodeVarint(uint64(size))
	if len(sizeBytes) > 1 {
		panic(fmt.Errorf("[]sizeBytes should not be more than one byte because the max number it needs to hold is 8. size=%d", size))
	}
	encodedBytes := make([]byte, size+1)
	encodedBytes[0] = sizeBytes[0]
	copy(encodedBytes[1:], bytes[startingIndex:])
	return encodedBytes
}

// decodeOrderPreservingVarUint64 decodes the number from the bytes obtained from method 'EncodeOrderPreservingVarUint64'.
// It returns the decoded number, the number of bytes that are consumed in the process, and an error if the input bytes are invalid.
func decodeOrderPreservingVarUint64(bytes []byte) (uint64, int, error) {
	s, numBytes := proto.DecodeVarint(bytes)

	switch {
	case numBytes != 1:
		return 0, 0, errors.Errorf("number of consumed bytes from DecodeVarint is invalid, expected 1, but got %d", numBytes)
	case s > 8:
		return 0, 0, errors.Errorf("decoded size from DecodeVarint is invalid, expected <=8, but got %d", s)
	case int(s) > len(bytes)-1:
		return 0, 0, errors.Errorf("decoded size (%d) from DecodeVarint is more than available bytes (%d)", s, len(bytes)-1)
	default:
		// no error
		size := int(s)
		decodedBytes := make([]byte, 8)
		copy(decodedBytes[8-size:], bytes[1:size+1])
		numBytesConsumed := size + 1
		return binary.BigEndian.Uint64(decodedBytes), numBytesConsumed, nil
	}
}

// encodeReverseOrderVarUint64 returns a byte-representation for a uint64 number such that
// the number is first subtracted from MaxUint64 and then all the leading 0xff bytes
// are trimmed and replaced by the number of such trimmed bytes. This helps in reducing the size.
// In the byte order comparison this encoding ensures that encodeReverseOrderVarUint64(A) > encodeReverseOrderVarUint64(B),
// If B > A
func encodeReverseOrderVarUint64(number uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, math.MaxUint64-number)
	numFFBytes := 0
	for _, b := range bytes {
		if b != 0xff {
			break
		}
		numFFBytes++
	}
	size := 8 - numFFBytes
	encodedBytes := make([]byte, size+1)
	encodedBytes[0] = proto.EncodeVarint(uint64(numFFBytes))[0]
	copy(encodedBytes[1:], bytes[numFFBytes:])
	return encodedBytes
}

// decodeReverseOrderVarUint64 decodes the number from the bytes obtained from function 'encodeReverseOrderVarUint64'.
// Also, returns the number of bytes that are consumed in the process
func decodeReverseOrderVarUint64(bytes []byte) (uint64, int) {
	s, _ := proto.DecodeVarint(bytes)
	numFFBytes := int(s)
	decodedBytes := make([]byte, 8)
	realBytesNum := 8 - numFFBytes
	copy(decodedBytes[numFFBytes:], bytes[1:realBytesNum+1])
	numBytesConsumed := realBytesNum + 1
	for i := 0; i < numFFBytes; i++ {
		decodedBytes[i] = 0xff
	}
	return (math.MaxUint64 - binary.BigEndian.Uint64(decodedBytes)), numBytesConsumed
}
