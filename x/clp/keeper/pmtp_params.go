package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetPmtpParams(ctx sdk.Context, params *types.PmtpParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PmtpParamsPrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetPmtpParams(ctx sdk.Context) *types.PmtpParams {
	params := types.PmtpParams{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PmtpParamsPrefix)
	k.cdc.MustUnmarshal(bz, &params)
	return &params
}
