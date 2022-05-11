package keeper

import (
	"encoding/binary"

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

	value := store.Get(prefix)
	globalNonce := binary.LittleEndian.Uint64(value)

	return globalNonce
}

// UpdateGlobalSequence get current global nonce and update it
func (k Keeper) UpdateGlobalSequence(ctx sdk.Context,
	networkDescriptor oracletypes.NetworkDescriptor,
	blockNumber uint64) {
	prefix := k.getGlobalSequencePrefix(ctx, networkDescriptor)
	store := ctx.KVStore(k.storeKey)

	globalNonce := k.GetGlobalSequence(ctx, networkDescriptor)

	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, globalNonce+1)
	store.Set(prefix, bs)
	k.SetGlobalSequenceToBlockNumber(ctx, networkDescriptor, globalNonce, blockNumber)
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

	value := store.Get(prefix)
	return binary.LittleEndian.Uint64(value)
}

// SetGlobalSequenceToBlockNumber
func (k Keeper) SetGlobalSequenceToBlockNumber(
	ctx sdk.Context,
	networkDescriptor oracletypes.NetworkDescriptor,
	globalNonce uint64,
	blockNumber uint64) {

	store := ctx.KVStore(k.storeKey)
	prefix := k.getGlobalSequenceToBlockNumberPrefix(ctx, networkDescriptor, globalNonce)

	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, blockNumber)

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetGlobalSequenceToBlockNumber, "networkDescriptor", networkDescriptor, "globalNonce", globalNonce, "blockNumber", blockNumber)

	store.Set(prefix, bs)
}

// getGlobalSequenceToBlockNumberPrefix
func (k Keeper) getGlobalSequenceToBlockNumberPrefix(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, globalNonce uint64) []byte {
	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	bs := k.cdc.MustMarshal(&networkIdentity)
	tmpKey := append(types.GlobalNonceToBlockNumberPrefix, bs[:]...)

	bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, globalNonce)

	return append(tmpKey, bs[:]...)
}
