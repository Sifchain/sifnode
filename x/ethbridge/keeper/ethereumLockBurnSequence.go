package keeper

import (
	"bytes"
	"errors"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetEthereumLockBurnSequence set the ethereum lock burn nonce for each relayer
func (k Keeper) SetEthereumLockBurnSequence(ctx sdk.Context, networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress, newSequence uint64) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetEthereumLockBurnSequencePrefix(networkDescriptor, valAccount)

	bs := k.cdc.MustMarshal(&types.EthereumLockBurnSequence{
		EthereumLockBurnSequence: newSequence,
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
	var EthereumLockBurnSequence types.EthereumLockBurnSequence
	k.cdc.MustUnmarshal(store.Get(key), &EthereumLockBurnSequence)

	return EthereumLockBurnSequence.EthereumLockBurnSequence
}

// GetEthereumLockBurnSequencePrefix return storage prefix
func (k Keeper) GetEthereumLockBurnSequencePrefix(networkDescriptor oracletypes.NetworkDescriptor, valAccount sdk.ValAddress) []byte {

	bs := k.cdc.MustMarshal(&types.EthereumLockBurnSequenceKey{
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
func (k Keeper) GetEthereumLockBurnSequences(ctx sdk.Context) []*types.GenesisEthereumLockBurnSequence {
	sequences := make([]*types.GenesisEthereumLockBurnSequence, 0)
	iterator := k.getEthereumLockBurnSequenceIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {

		ethereumLockBurnSequenceKey, err := getEthereumLockBurnSequenceKeyFromRawKey(k.cdc, iterator.Key(), types.EthereumLockBurnSequencePrefix)
		if err != nil {
			panic(err)
		}

		var lockBurnSequence types.EthereumLockBurnSequence
		k.cdc.MustUnmarshal(iterator.Value(), &lockBurnSequence)

		sequences = append(sequences, &types.GenesisEthereumLockBurnSequence{
			EthereumLockBurnSequenceKey: &ethereumLockBurnSequenceKey,
			EthereumLockBurnSequence:    &lockBurnSequence,
		})
	}
	return sequences
}

func getEthereumLockBurnSequenceKeyFromRawKey(cdc codec.BinaryCodec, key []byte, prefix []byte) (types.EthereumLockBurnSequenceKey, error) {
	if bytes.HasPrefix(key, prefix) {
		var ethereumLockBurnSequenceKey types.EthereumLockBurnSequenceKey
		err := cdc.Unmarshal(key[len(prefix):], &ethereumLockBurnSequenceKey)

		if err == nil {
			return ethereumLockBurnSequenceKey, nil
		}
		return types.EthereumLockBurnSequenceKey{}, err
	}

	return types.EthereumLockBurnSequenceKey{}, errors.New("prefix for EthereumLockBurnSequenceKey is invalid")
}
