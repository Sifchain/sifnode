package types

import (
	"bytes"
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
)

// GlobalSequenceKey return storage prefix
// func (key GlobalSequenceKey) GetGlobalSequenceKeyPrefix(cdc codec.BinaryCodec) []byte {
// 	buf := cdc.MustMarshal(&key)
// 	return append(GlobalNonceProphecyIDPrefix, buf[:]...)
// }

// Get the GetLockBurnSequenceKeyFromRawKey from storage key in keeper. storage key = WitnessLockBurnNoncePrefix + LockBurnSequenceKey
func GetGlobalSequenceKeyKeyFromRawKey(cdc codec.BinaryCodec, key []byte) (GlobalSequenceKey, error) {
	// check the key which correct prefix
	if bytes.HasPrefix(key, GlobalNonceProphecyIDPrefix) {
		var globalSequenceKey GlobalSequenceKey
		err := cdc.Unmarshal(key[len(GlobalNonceProphecyIDPrefix):], &globalSequenceKey)

		if err == nil {
			return globalSequenceKey, nil
		}
		return globalSequenceKey, err
	}

	return GlobalSequenceKey{}, errors.New("GlobalSequenceKey prefix is invalid")
}
