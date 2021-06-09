package stateindex

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderPreservingEncodingDecoding(t *testing.T) {
	for i := 0; i < 10000; i++ {
		value := encodeOrderPreservingVarUint64(uint64(i))
		nextValue := encodeOrderPreservingVarUint64(uint64(i + 1))
		if !(bytes.Compare(value, nextValue) < 0) {
			t.Fatalf("A smaller integer should result into smaller bytes. Encoded bytes for [%d] is [%x] and for [%d] is [%x]",
				i, i+1, value, nextValue)
		}
		decodedValue, _, err := decodeOrderPreservingVarUint64(value)
		require.NoError(t, err, "Error via calling DecodeOrderPreservingVarUint64")
		if decodedValue != uint64(i) {
			t.Fatalf("Value not same after decoding. Original value = [%d], decode value = [%d]", i, decodedValue)
		}
	}
}

func TestReverseOrderPreservingEncodingDecoding(t *testing.T) {
	for i := 0; i < 10000; i++ {
		value := encodeReverseOrderVarUint64(uint64(i))
		nextValue := encodeReverseOrderVarUint64(uint64(i + 1))
		if !(bytes.Compare(value, nextValue) > 0) {
			t.Fatalf("A smaller integer should result into greater bytes. Encoded bytes for [%d] is [%x] and for [%d] is [%x]",
				i, i+1, value, nextValue)
		}
		decodedValue, _ := decodeReverseOrderVarUint64(value)
		if decodedValue != uint64(i) {
			t.Fatalf("Value not same after decoding. Original value = [%d], decode value = [%d]", i, decodedValue)
		}
	}
}
