package keeper

import (
	"github.com/Sifchain/sifnode/x/instrumentation"
	"github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetWitnessLockBurnNonce set the Witness lock burn nonce for each relayer
func (k Keeper) SetWitnessLockBurnNonce(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress, newNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetWitnessLockBurnSequencePrefix(networkDescriptor, valAccount)

	bs := k.cdc.MustMarshal(&types.LockBurnNonce{
		LockBurnNonce: newNonce,
	})

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetWitnessLockBurnNonce, "networkDescriptor", networkDescriptor, "valAccount", valAccount, "newNonce", newNonce, "key", key)

	store.Set(key, bs)
}

// GetWitnessLockBurnSequence return Witness lock burn nonce
func (k Keeper) GetWitnessLockBurnSequence(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := k.GetWitnessLockBurnSequencePrefix(networkDescriptor, valAccount)

	// nonce start from 1, 0 represent the relayer is a new one
	if !store.Has(key) {
		return 0
	}

	var lockBurnNonce types.LockBurnNonce
	k.cdc.MustUnmarshal(store.Get(key), &lockBurnNonce)

	return lockBurnNonce.LockBurnNonce
}

// GetWitnessLockBurnSequencePrefix return storage prefix
func (k Keeper) GetWitnessLockBurnSequencePrefix(networkDescriptor types.NetworkDescriptor, valAccount sdk.ValAddress) []byte {
	bs := k.cdc.MustMarshal(&types.LockBurnNonceKey{
		NetworkDescriptor: networkDescriptor,
		ValidatorAddress:  valAccount,
	})

	return append(types.WitnessLockBurnNoncePrefix, bs[:]...)
}
