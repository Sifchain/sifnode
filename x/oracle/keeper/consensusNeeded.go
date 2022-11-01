package keeper

import (
	"errors"
	"fmt"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetConsensusNeeded for a network.
func (k Keeper) SetConsensusNeeded(ctx sdk.Context,
	networkIdentity types.NetworkIdentity,
	consensusNeeded types.ConsensusNeeded) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix(k.cdc)

	bs := k.cdc.MustMarshal(&consensusNeeded)

	store.Set(key, bs)
}

// GetConsensusNeeded for a network
func (k Keeper) GetConsensusNeeded(ctx sdk.Context, networkIdentity types.NetworkIdentity) (uint32, error) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix(k.cdc)

	if !store.Has(key) {
		return 0.0, fmt.Errorf("%s%s", "ConsensusNeeded not set for ", networkIdentity.NetworkDescriptor.String())
	}

	bz := store.Get(key)
	var consensusNeeded types.ConsensusNeeded
	k.cdc.MustUnmarshal(bz, &consensusNeeded)

	if consensusNeeded.ConsensusNeeded > 100 {
		return 0, errors.New("consensusNeeded stored is too large")
	}
	return consensusNeeded.ConsensusNeeded, nil
}

// GetAllConsensusNeeded get consensus needed for all network descriptors
func (k Keeper) GetAllConsensusNeeded(ctx sdk.Context) []*types.NetworkConfigData {
	consensuses := make([]*types.NetworkConfigData, 0)
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ConsensusNeededPrefix)

	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()

		networkIdentity, err := types.GetFromPrefix(k.cdc, key, types.ConsensusNeededPrefix)
		if err != nil {
			panic(err)
		}

		bz := store.Get(key)
		var consensusNeeded types.ConsensusNeeded
		k.cdc.MustUnmarshal(bz, &consensusNeeded)

		consensuses = append(consensuses, &types.NetworkConfigData{
			NetworkDescriptor: networkIdentity.NetworkDescriptor,
			ConsensusNeeded:   &consensusNeeded,
		})
	}
	return consensuses
}
