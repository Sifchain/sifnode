package keeper

import (
	"github.com/Sifchain/sifnode/x/instrumentation"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetGlobalSequence get current sequence.  Default is 1 if there's no existing value stored.
func (k Keeper) GetGlobalSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) uint64 {
	prefix := k.getGlobalSequencePrefix(ctx, networkDescriptor)
	store := ctx.KVStore(k.storeKey)

	if !k.Exists(ctx, prefix) {
		// global nonce start from 1
		return uint64(1)
	}

	var globalSequence oracletypes.GlobalSequence
	k.cdc.MustUnmarshal(store.Get(prefix), &globalSequence)

	return globalSequence.GlobalSequence
}

// UpdateGlobalSequence get current global nonce and update it
func (k Keeper) UpdateGlobalSequence(ctx sdk.Context,
	networkDescriptor oracletypes.NetworkDescriptor,
	blockNumber uint64) {
	prefix := k.getGlobalSequencePrefix(ctx, networkDescriptor)
	store := ctx.KVStore(k.storeKey)
	globalSequence := k.GetGlobalSequence(ctx, networkDescriptor)

	bs := k.cdc.MustMarshal(&oracletypes.GlobalSequence{
		GlobalSequence: globalSequence + 1,
	})
	store.Set(prefix, bs)
	k.SetGlobalSequenceToBlockNumber(ctx, networkDescriptor, globalSequence, blockNumber)
}

// getGlobalSequencePrefix compute the prefix
func (k Keeper) getGlobalSequencePrefix(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) []byte {
	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	bs := k.cdc.MustMarshal(&networkIdentity)
	return append(types.GlobalNoncePrefix, bs[:]...)
}

// GetGlobalSequenceToBlockNumber
func (k Keeper) GetGlobalSequenceToBlockNumber(
	ctx sdk.Context,
	networkDescriptor oracletypes.NetworkDescriptor,
	globalSequence uint64) uint64 {

	store := ctx.KVStore(k.storeKey)
	prefix := k.getGlobalSequenceToBlockNumberPrefix(ctx, networkDescriptor, globalSequence)

	if !k.Exists(ctx, prefix) {
		return uint64(0)
	}

	var blockNumber oracletypes.BlockNumber

	k.cdc.MustUnmarshal(store.Get(prefix), &blockNumber)
	return blockNumber.BlockNumber
}

// SetGlobalSequenceToBlockNumber
func (k Keeper) SetGlobalSequenceToBlockNumber(
	ctx sdk.Context,
	networkDescriptor oracletypes.NetworkDescriptor,
	globalNonce uint64,
	blockNumber uint64) {

	store := ctx.KVStore(k.storeKey)
	prefix := k.getGlobalSequenceToBlockNumberPrefix(ctx, networkDescriptor, globalNonce)

	bs := k.cdc.MustMarshal(&oracletypes.BlockNumber{
		BlockNumber: blockNumber,
	})

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetGlobalSequenceToBlockNumber, "networkDescriptor", networkDescriptor, "globalNonce", globalNonce, "blockNumber", blockNumber)

	store.Set(prefix, bs)
}

// getGlobalSequenceToBlockNumberPrefix
func (k Keeper) getGlobalSequenceToBlockNumberPrefix(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, globalSequence uint64) []byte {
	bs := k.cdc.MustMarshal(&oracletypes.GlobalSequenceKey{
		NetworkDescriptor: networkDescriptor,
		GlobalSequence:    globalSequence,
	})

	return append(types.GlobalNonceToBlockNumberPrefix, bs[:]...)
}
