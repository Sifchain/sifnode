package types

import (
	"bytes"
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
)

// GetWitnessLockBurnSequencePrefix return storage prefix
func (key LockBurnSequenceKey) GetWitnessLockBurnSequencePrefix(cdc codec.BinaryCodec) []byte {
	buf := cdc.MustMarshal(&key)
	return append(WitnessLockBurnNoncePrefix, buf[:]...)
}

// Get the GetLockBurnSequenceKeyFromRawKey from storage key in keeper. storage key = WitnessLockBurnNoncePrefix + LockBurnSequenceKey
func GetWitnessLockBurnSequenceKeyFromRawKey(cdc codec.BinaryCodec, key []byte) (LockBurnSequenceKey, error) {
	// check the key which correct prefix
	if bytes.HasPrefix(key, WitnessLockBurnNoncePrefix) {
		var lockBurnSequenceKey LockBurnSequenceKey
		err := cdc.Unmarshal(key[len(WitnessLockBurnNoncePrefix):], &lockBurnSequenceKey)

		if err == nil {
			return lockBurnSequenceKey, nil
		}
		return lockBurnSequenceKey, err
	}

	return LockBurnSequenceKey{}, errors.New("LockBurnSequenceKey prefix is invalid")
}
