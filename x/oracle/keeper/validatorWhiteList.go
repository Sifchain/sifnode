package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetOracleWhiteList(ctx sdk.Context, validatorList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(validatorList))
}

func (k Keeper) ExistsOracleWhiteList(ctx sdk.Context) bool {
	key := types.WhiteListValidatorPrefix
	return k.Exists(ctx, key)
}

func (k Keeper) GetOracleWhiteList(ctx sdk.Context) (valList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &valList)
	return
}

func (k Keeper) ValidateAddress(ctx sdk.Context, address sdk.ValAddress) bool {
	if !k.ExistsOracleWhiteList(ctx) {
		return false
	}
	valList := k.GetOracleWhiteList(ctx)

	for _, validator := range valList {
		if validator.Equals(address) {
			return true
		}
	}
	return false
}
