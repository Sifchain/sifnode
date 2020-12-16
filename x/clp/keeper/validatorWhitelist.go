package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetClpWhiteList(ctx sdk.Context, validatorList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(validatorList))
}

func (k Keeper) ExistsClpWhiteList(ctx sdk.Context) bool {
	key := types.WhiteListValidatorPrefix
	return k.Exists(ctx, key)
}

func (k Keeper) GetClpWhiteList(ctx sdk.Context) (valList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &valList)
	return
}

func (k Keeper) ValidateAddress(ctx sdk.Context, address sdk.AccAddress) bool {
	if !k.ExistsClpWhiteList(ctx) {
		return false
	}
	valList := k.GetClpWhiteList(ctx)

	for _, validator := range valList {
		if validator.Equals(address) {
			return true
		}
	}
	return false
}
