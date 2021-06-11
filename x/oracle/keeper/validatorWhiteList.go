package keeper

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetOracleWhiteList(ctx sdk.Context, validatorList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := k.Cdc.MustMarshalBinaryBare(validatorList)
	if bz == nil {
		// empty arrays marshal to nil, and set panics on nil values, so an empty array must cause key deletion
		store.Delete(key)
		return
	}

	store.Set(key, bz)
}

func (k Keeper) ExistsOracleWhiteList(ctx sdk.Context) bool {
	key := types.WhiteListValidatorPrefix
	return k.Exists(ctx, key)
}

//
func (k Keeper) GetOracleWhiteList(ctx sdk.Context) (valList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := store.Get(key)
	k.Cdc.MustUnmarshalBinaryBare(bz, &valList)
	return
}

// ValidateAddress is a validator in whitelist
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

// AddOracleWhiteList add new validator to whitelist
func (k Keeper) AddOracleWhiteList(ctx sdk.Context, validator sdk.ValAddress) {
	valList := k.GetOracleWhiteList(ctx)
	k.SetOracleWhiteList(ctx, append(valList, validator))
}

// RemoveOracleWhiteList remove a validator from whitelist
func (k Keeper) RemoveOracleWhiteList(ctx sdk.Context, validator sdk.ValAddress) {
	valList := k.GetOracleWhiteList(ctx)

	for index, item := range valList {
		if validator.Equals(item) {
			k.SetOracleWhiteList(ctx, append(valList[:index], valList[index+1]))
			return
		}
	}
}
