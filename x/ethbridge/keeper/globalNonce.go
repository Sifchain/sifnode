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

	if !k.existsGlobalNonce(ctx, prefix) {
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

// existsGlobalNonce check if the global nonce exists
func (k Keeper) existsGlobalNonce(ctx sdk.Context, prefix []byte) bool {
	if !k.Exists(ctx, prefix) {
		// The store doesnt exist.
		return false
	}
	return true
}

// getGlobalSequencePrefix compute the prefix
// TODO: oracletypes.NetworkDescriptor is int32 (default type for enums), we are converting it to uint32 here.
func (k Keeper) getGlobalSequencePrefix(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(networkDescriptor))

	return append(types.GlobalNoncePrefix, bs[:]...)
}

// GetGlobalSequenceToBlockNumber
func (k Keeper) GetGlobalSequenceToBlockNumber(
	ctx sdk.Context,
	networkDescriptor oracletypes.NetworkDescriptor,
	globalSequence uint64) uint64 {

	store := ctx.KVStore(k.storeKey)
	prefix := k.getGlobalSequenceToBlockNumberPrefix(ctx, networkDescriptor, globalNonce)

	if !k.existsGlobalNonce(ctx, prefix) {
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
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(networkDescriptor))
	tmpKey := append(types.GlobalNonceToBlockNumberPrefix, bs[:]...)

	bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, globalNonce)

	return append(tmpKey, bs[:]...)
}
