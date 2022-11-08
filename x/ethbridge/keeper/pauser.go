package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetPause(ctx sdk.Context, pause *types.Pause) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PausePrefix, k.cdc.MustMarshal(pause))
}

func (k Keeper) getPause(ctx sdk.Context) *types.Pause {
	pause := types.Pause{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PausePrefix)
	k.cdc.MustUnmarshal(bz, &pause)
	return &pause
}

func (k Keeper) IsPaused(ctx sdk.Context) bool {
	return k.getPause(ctx).IsPaused
}
