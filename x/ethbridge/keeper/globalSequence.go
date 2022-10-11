package keeper

import (
	"bytes"
	"errors"

	"github.com/Sifchain/sifnode/x/instrumentation"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
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
func (k Keeper) GetGlobalSequences(ctx sdk.Context) []*types.GenesisGlobalSequence {
	sequences := make([]*types.GenesisGlobalSequence, 0)
	iterator := k.getGlobalSequenceIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {

		networkIdentity, err := oracletypes.GetFromPrefix(k.cdc, iterator.Key(), types.GlobalNoncePrefix)
		if err != nil {
			panic(err)
		}

		var globalSequence oracletypes.GlobalSequence

		k.cdc.MustUnmarshal(iterator.Value(), &globalSequence)

		sequences = append(sequences, &types.GenesisGlobalSequence{
			NetworkDescriptor: networkIdentity.NetworkDescriptor,
			GlobalSequence:    &globalSequence,
		})
	}
	return sequences
}

func (k Keeper) getGlobalSequenceToBlockNumberIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.GlobalNonceToBlockNumberPrefix)
}

// GetGlobalSequenceToBlockNumbers get all data from keeper
func (k Keeper) GetGlobalSequenceToBlockNumbers(ctx sdk.Context) []*types.GenesisGlobalSequenceBlockNumber {
	globalSequenceBlockNumber := make([]*types.GenesisGlobalSequenceBlockNumber, 0)
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
		globalSequenceKey, err := getGlobalSequenceKeyFromRawKey(k.cdc, key, types.GlobalNonceToBlockNumberPrefix)
		if err != nil {
			panic(err)
		}

		var blockNumber oracletypes.BlockNumber
		k.cdc.MustUnmarshal(value, &blockNumber)

		globalSequenceBlockNumber = append(globalSequenceBlockNumber, &types.GenesisGlobalSequenceBlockNumber{
			GlobalSequenceKey: &globalSequenceKey,
			BlockNumber:       &blockNumber,
		})
	}
	return globalSequenceBlockNumber
}

func getGlobalSequenceKeyFromRawKey(cdc codec.BinaryCodec, key []byte, prefix []byte) (oracletypes.GlobalSequenceKey, error) {
	if bytes.HasPrefix(key, prefix) {
		var globalSequenceKey oracletypes.GlobalSequenceKey
		err := cdc.Unmarshal(key[len(prefix):], &globalSequenceKey)

		if err == nil {
			return globalSequenceKey, nil
		}
		return oracletypes.GlobalSequenceKey{}, err
	}

	return oracletypes.GlobalSequenceKey{}, errors.New("prefix for GlobalSequenceKey is invalid")
}
