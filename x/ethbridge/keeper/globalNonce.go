package keeper

import (
	"github.com/Sifchain/sifnode/x/instrumentation"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetGlobalSequence get current sequence.  Default is 1 if there's no existing value stored.
func (k Keeper) GetGlobalSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) uint64 {
	prefix := k.GetGlobalSequencePrefix(ctx, networkDescriptor)
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
	prefix := k.GetGlobalSequencePrefix(ctx, networkDescriptor)
	store := ctx.KVStore(k.storeKey)
	globalSequence := k.GetGlobalSequence(ctx, networkDescriptor)

	bs := k.cdc.MustMarshal(&oracletypes.GlobalSequence{
		GlobalSequence: globalSequence + 1,
	})
	store.Set(prefix, bs)
	k.SetGlobalSequenceToBlockNumber(ctx, networkDescriptor, globalSequence, blockNumber)
}

// GetGlobalSequencePrefix compute the prefix
func (k Keeper) GetGlobalSequencePrefix(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor) []byte {
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
	prefix := k.GetGlobalSequenceToBlockNumberPrefix(ctx, networkDescriptor, globalSequence)

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
	prefix := k.GetGlobalSequenceToBlockNumberPrefix(ctx, networkDescriptor, globalNonce)

	bs := k.cdc.MustMarshal(&oracletypes.BlockNumber{
		BlockNumber: blockNumber,
	})

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.SetGlobalSequenceToBlockNumber, "networkDescriptor", networkDescriptor, "globalNonce", globalNonce, "blockNumber", blockNumber)

	store.Set(prefix, bs)
}

// GetGlobalSequenceToBlockNumberPrefix
func (k Keeper) GetGlobalSequenceToBlockNumberPrefix(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, globalSequence uint64) []byte {
	bs := k.cdc.MustMarshal(&oracletypes.GlobalSequenceKey{
		NetworkDescriptor: networkDescriptor,
		GlobalSequence:    globalSequence,
	})

	return append(types.GlobalNonceToBlockNumberPrefix, bs[:]...)
}

func (k Keeper) getGlobalSequenceIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GlobalNoncePrefix)
}

// GetGlobalSequences get all sequences from keeper
func (k Keeper) GetGlobalSequences(ctx sdk.Context) map[uint32]uint64 {
	sequences := make(map[uint32]uint64)
	iterator := k.getGlobalSequenceIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		var networkIdentity oracletypes.NetworkIdentity
		if len(key) < len(types.GlobalNoncePrefix) {
			panic("the key for global nonce is valid")
		}
		k.cdc.MustUnmarshal(key[len(types.GlobalNoncePrefix):], &networkIdentity)
		if networkIdentity.NetworkDescriptor < 0 {
			panic("network identity value is invalid")
		}

		value := iterator.Value()
		var globalSequence oracletypes.GlobalSequence
		k.cdc.MustUnmarshal(value, &globalSequence)

		sequences[uint32(networkIdentity.NetworkDescriptor)] = globalSequence.GlobalSequence
	}
	return sequences
}

// SetGlobalSequenceViaRawKey used in import sequence from genesis
func (k Keeper) SetGlobalSequenceViaRawKey(ctx sdk.Context, networkDescriptor uint32, globalSequence uint64) {
	store := ctx.KVStore(k.storeKey)

	prefix := k.GetGlobalSequencePrefix(ctx, oracletypes.NetworkDescriptor(networkDescriptor))

	bs := k.cdc.MustMarshal(&oracletypes.GlobalSequence{
		GlobalSequence: globalSequence,
	})

	store.Set(prefix, bs)
}

func (k Keeper) getGlobalSequenceToBlockNumberIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GlobalNonceToBlockNumberPrefix)
}

// GetGlobalSequenceToBlockNumbers get all data from keeper
func (k Keeper) GetGlobalSequenceToBlockNumbers(ctx sdk.Context) map[string]uint64 {
	blockNumbers := make(map[string]uint64)
	iterator := k.getGlobalSequenceToBlockNumberIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()
		var blockNumber oracletypes.BlockNumber
		k.cdc.MustUnmarshal(value, &blockNumber)

		blockNumbers[string(key)] = blockNumber.BlockNumber
	}
	return blockNumbers
}

// SetGlobalSequenceToBlockNumberViaRawKey used in import data from genesis
func (k Keeper) SetGlobalSequenceToBlockNumberViaRawKey(ctx sdk.Context, key string, blockNumber uint64) {
	store := ctx.KVStore(k.storeKey)
	bs := k.cdc.MustMarshal(&oracletypes.BlockNumber{
		BlockNumber: blockNumber,
	})

	store.Set([]byte(key), bs)
}
