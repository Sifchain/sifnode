package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetCrossChainFee set the crosschain fee for a network.
func (k Keeper) SetConsensusNeeded(ctx sdk.Context,
	networkIdentity types.NetworkIdentity,
	consensusNeeded float32) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix()

	store.Set(key, k.cdc.MustMarshalBinaryBare(&types.ConsensusNeeded{
		ConsensusNeeded: consensusNeeded,
	}))
}

// GetCrossChainFeeConfig return crosschain fee config
func (k Keeper) GetConsensusNeeded(ctx sdk.Context, networkIdentity types.NetworkIdentity) (float64, error) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetConsensusNeededPrefix()

	if !store.Has(key) {
		return 0.0, fmt.Errorf("%s%s", "ConsensusNeeded not set for ", networkIdentity.NetworkDescriptor.String())
	}

	bz := store.Get(key)
	consensusNeeded := &types.ConsensusNeeded{}
	k.cdc.MustUnmarshalBinaryBare(bz, consensusNeeded)

	return float64(consensusNeeded.ConsensusNeeded), nil
}
