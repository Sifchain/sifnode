package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetPmtpEpoch(ctx sdk.Context, params types.PmtpEpoch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PmtpEpochPrefix, k.cdc.MustMarshal(&params))
}

func (k Keeper) GetPmtpEpoch(ctx sdk.Context) types.PmtpEpoch {
	epoch := types.PmtpEpoch{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PmtpEpochPrefix)
	k.cdc.MustUnmarshal(bz, &epoch)
	return epoch
}

func (k Keeper) DecrementPmtpEpochCounter(ctx sdk.Context) {
	epoch := k.GetPmtpEpoch(ctx)
	epoch.EpochCounter--
	k.SetPmtpEpoch(ctx, epoch)
}

func (k Keeper) DecrementPmtpBlockCounter(ctx sdk.Context) {
	epoch := k.GetPmtpEpoch(ctx)
	epoch.BlockCounter--
	k.SetPmtpEpoch(ctx, epoch)
}

func (k Keeper) SetPmtpBlockCounter(ctx sdk.Context, epochLength int64) {
	epoch := k.GetPmtpEpoch(ctx)
	epoch.BlockCounter = epochLength
	k.SetPmtpEpoch(ctx, epoch)
}
