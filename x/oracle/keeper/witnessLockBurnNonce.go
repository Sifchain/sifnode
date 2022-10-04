package keeper

import (
	"github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetWitnessLockBurnNonce set the Witness lock burn nonce for each relayer
func (k Keeper) SetWitnessLockBurnNonce(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress, lockBurnSequence uint64) {

	lockBurnSequenceKey := types.LockBurnSequenceKey{
		NetworkDescriptor: networkDescriptor,
		ValidatorAddress:  valAccount,
	}
	lockBurnSequenceObj := types.LockBurnSequence{
		LockBurnSequence: lockBurnSequence,
	}
	k.SetWitnessLockBurnNonceObj(ctx, lockBurnSequenceKey, lockBurnSequenceObj)
}

// SetWitnessLockBurnNonce set the Witness lock burn nonce for each relayer
func (k Keeper) SetWitnessLockBurnNonceObj(ctx sdk.Context, lockBurnSequenceKey types.LockBurnSequenceKey, lockBurnSequence types.LockBurnSequence) {
	store := ctx.KVStore(k.storeKey)
	key := lockBurnSequenceKey.GetWitnessLockBurnSequencePrefix(k.cdc)

	bs := k.cdc.MustMarshal(&lockBurnSequence)
	store.Set(key, bs)
}

// GetWitnessLockBurnSequence return Witness lock burn nonce
func (k Keeper) GetWitnessLockBurnSequence(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)

	lockBurnSequenceKey := types.LockBurnSequenceKey{
		NetworkDescriptor: networkDescriptor,
		ValidatorAddress:  valAccount,
	}

	key := lockBurnSequenceKey.GetWitnessLockBurnSequencePrefix(k.cdc)

	// nonce start from 1, 0 represent the relayer is a new one
	if !store.Has(key) {
		return 0
	}

	var lockBurnSequence types.LockBurnSequence
	k.cdc.MustUnmarshal(store.Get(key), &lockBurnSequence)

	return lockBurnSequence.LockBurnSequence
}

// GetWitnessLockBurnSequencePrefix return storage prefix
// func (k Keeper) GetWitnessLockBurnSequencePrefix(lockBurnSequenceKey types.LockBurnSequenceKey) []byte {
// 	bs := k.cdc.MustMarshal(&lockBurnSequenceKey)
// 	return append(types.WitnessLockBurnNoncePrefix, bs[:]...)
// }

// GetWitnessLockBurnSequencePrefix return storage prefix
// func (k Keeper) GetWitnessLockBurnSequencePrefix(lockBurnSequenceKey types.LockBurnSequenceKey) []byte {
// 	bs := k.cdc.MustMarshal(&lockBurnSequenceKey)
// 	return append(types.WitnessLockBurnNoncePrefix, bs[:]...)
// }

// GetAllWitnessLockBurnSequence get all witnessLockBurnSequence needed for all validators
func (k Keeper) GetAllWitnessLockBurnSequence(ctx sdk.Context) []*types.GenesisWitnessLockBurnSequence {
	sequences := make([]*types.GenesisWitnessLockBurnSequence, 0)
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.WitnessLockBurnNoncePrefix)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var lockBurnSequenceKey types.LockBurnSequenceKey
		var lockBurnSequence types.LockBurnSequence

		lockBurnSequenceKey, err := types.GetWitnessLockBurnSequenceKeyFromRawKey(k.cdc, iterator.Key())
		if err != nil {
			panic(err)
		}

		k.cdc.MustUnmarshal(iterator.Value(), &lockBurnSequence)
		sequences = append(sequences, &types.GenesisWitnessLockBurnSequence{
			WitnessLockBurnSequenceKey: &lockBurnSequenceKey,
			WitnessLockBurnSequence:    &lockBurnSequence,
		})
	}
	return sequences
}

// func (k Keeper) SetWitnessLockBurnNonceViaRawKey(ctx sdk.Context, key []byte, nonce uint64) {

// 	store := ctx.KVStore(k.storeKey)
// 	bs := k.cdc.MustMarshal(&types.LockBurnNonce{
// 		LockBurnNonce: nonce,
// 	})

// 	store.Set(key, bs)
// }
