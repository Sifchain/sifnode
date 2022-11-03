package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetEthereumLockBurnSequence set the ethereum lock burn nonce for each relayer
func (k Keeper) SetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress, newNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.getEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	bs := k.cdc.MustMarshal(&oracletypes.LockBurnNonce{
		LockBurnNonce: newNonce,
	})

	store.Set(key, bs)
}

// GetEthereumLockBurnSequence return ethereum lock burn nonce
func (k Keeper) GetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := k.getEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	// nonces start from 0, and the first ethereum transaction
	// should have a nonce of 1
	if !store.Has(key) {
		return 0
	}
	var lockBurnNonce oracletypes.LockBurnNonce
	k.cdc.MustUnmarshal(store.Get(key), &lockBurnNonce)

	return lockBurnNonce.LockBurnNonce
}

// getEthereumLockBurnSequencePrefix return storage prefix
func (k Keeper) getEthereumLockBurnSequencePrefix(networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) []byte {

	bs := k.cdc.MustMarshal(&oracletypes.LockBurnNonceKey{
		NetworkDescriptor: networkDescriptor,
		ValidatorAddress:  valAccount,
	})
	return append(types.EthereumLockBurnSequencePrefix, bs[:]...)
}
