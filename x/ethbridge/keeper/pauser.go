package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetPauser(ctx sdk.Context, pauser *types.Pauser) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PauserPrefix, k.cdc.MustMarshal(pauser))
}

func (k Keeper) GetPauser(ctx sdk.Context) *types.Pauser {
	pauser := types.Pauser{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PauserPrefix)
	k.cdc.MustUnmarshal(bz, &pauser)
	return &pauser
}

func (k Keeper) Pause(ctx sdk.Context) {
	pauser := k.GetPauser(ctx)
	pauser.IsPaused = true
	k.SetPauser(ctx, pauser)
}

func (k Keeper) IsPaused(ctx sdk.Context) bool {
	return k.GetPauser(ctx).IsPaused
}
