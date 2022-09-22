package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetEthereumLockBurnSequence set the ethereum lock burn nonce for each relayer
func (k Keeper) SetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress, newSequence uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	bs := k.cdc.MustMarshal(&oracletypes.LockBurnNonce{
		LockBurnNonce: newSequence,
	})

	store.Set(key, bs)
}

// GetEthereumLockBurnSequence return ethereum lock burn nonce
func (k Keeper) GetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := k.GetEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	// nonces start from 0, and the first ethereum transaction
	// should have a nonce of 1
	if !store.Has(key) {
		return 0
	}
	var lockBurnNonce oracletypes.LockBurnNonce
	k.cdc.MustUnmarshal(store.Get(key), &lockBurnNonce)

	return lockBurnNonce.LockBurnNonce
}

// GetEthereumLockBurnSequencePrefix return storage prefix
func (k Keeper) GetEthereumLockBurnSequencePrefix(networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) []byte {

	bs := k.cdc.MustMarshal(&oracletypes.LockBurnNonceKey{
		NetworkDescriptor: networkDescriptor,
		ValidatorAddress:  valAccount,
	})
	return append(types.EthereumLockBurnSequencePrefix, bs[:]...)
}

func (k Keeper) getEthereumLockBurnSequenceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.EthereumLockBurnSequencePrefix)
}

// GetEthereumLockBurnSequences get all sequences from keeper
func (k Keeper) GetEthereumLockBurnSequences(ctx sdk.Context) map[string]uint64 {
	sequences := make(map[string]uint64)
	iterator := k.getEthereumLockBurnSequenceIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()
		var lockBurnNonce oracletypes.LockBurnNonce
		k.cdc.MustUnmarshal(value, &lockBurnNonce)
		sequences[string(key)] = lockBurnNonce.LockBurnNonce
	}
	return sequences
}

// SetSequenceViaRawKey used in import sequence from genesis
func (k Keeper) SetSequenceViaRawKey(ctx sdk.Context, key []byte, newSequence uint64) {
	// network, address := DecodeKey(key)
	var lockBurnNonceKey oracletypes.LockBurnNonceKey
	k.cdc.MustUnmarshal(key[len(types.EthereumLockBurnSequencePrefix):], &lockBurnNonceKey)
	k.SetEthereumLockBurnSequence(ctx, lockBurnNonceKey.NetworkDescriptor, lockBurnNonceKey.ValidatorAddress, newSequence)
}
