package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetValidatorWhiteList(ctx sdk.Context, validatorList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(validatorList))
}

func (k Keeper) ExistsValidatorWhiteList(ctx sdk.Context) bool {
	key := types.WhiteListValidatorPrefix
	return k.Exists(ctx, key)
}

func (k Keeper) GetValidatorWhiteList(ctx sdk.Context) (valList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &valList)
	return
}

func (k Keeper) ValidateAddress(ctx sdk.Context, address sdk.AccAddress) bool {
	if !k.ExistsValidatorWhiteList(ctx) {
		return false
	}
	valList := k.GetValidatorWhiteList(ctx)

	for _, validator := range valList {
		if validator.Equals(address) {
			return true
		}
	}
	return false
}
