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
	consensusNeeded uint32) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix(k.cdc)

	bs := k.cdc.MustMarshal(&types.ConsensusNeeded{
		ConsensusNeeded: consensusNeeded,
	})

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
