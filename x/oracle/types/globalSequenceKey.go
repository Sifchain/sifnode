package types

import (
	"bytes"
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
)

// GetGlobalSequenceKeyKeyFromRawKey from storage key in keeper. storage key = GlobalSequenceProphecyIDPrefix + LockBurnSequenceKey
func GetGlobalSequenceKeyKeyFromRawKey(cdc codec.BinaryCodec, key []byte) (GlobalSequenceKey, error) {
	// check the key which correct prefix
	if bytes.HasPrefix(key, GlobalSequenceProphecyIDPrefix) {
		var globalSequenceKey GlobalSequenceKey
		err := cdc.Unmarshal(key[len(GlobalSequenceProphecyIDPrefix):], &globalSequenceKey)

		if err == nil {
			return globalSequenceKey, nil
		}
		return globalSequenceKey, err
	}

	return GlobalSequenceKey{}, errors.New("GlobalSequenceKey prefix is invalid")
}
