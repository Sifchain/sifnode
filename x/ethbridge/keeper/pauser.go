package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetPauser(ctx sdk.Context, pauser *types.Pauser) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PauserPrefix, k.cdc.MustMarshal(pauser))
}

func (k Keeper) getPauser(ctx sdk.Context) *types.Pauser {
	pauser := types.Pauser{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PauserPrefix)
	k.cdc.MustUnmarshal(bz, &pauser)
	return &pauser
}

func (k Keeper) IsPaused(ctx sdk.Context) bool {
	return k.getPauser(ctx).IsPaused
}
