package keeper

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validatorList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	store.Set(key, k.cdc.MustMarshalBinaryBare(validatorList))
}

func (k Keeper) ExistsOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) bool {
	key := networkDescriptor.GetPrefix()
	return k.Exists(ctx, key)
}

//
func (k Keeper) GetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) (valList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &valList)
	return
}

// ValidateAddress is a validator in whitelist
func (k Keeper) ValidateAddress(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, address sdk.ValAddress) bool {
	if !k.ExistsOracleWhiteList(ctx, networkDescriptor) {
		return false
	}
	valList := k.GetOracleWhiteList(ctx, networkDescriptor)

	for _, validator := range valList {
		if validator.Equals(address) {
			return true
		}
	}
	return false
}

// AddOracleWhiteList add new validator to whitelist
func (k Keeper) AddOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validator sdk.ValAddress) {
	valList := k.GetOracleWhiteList(ctx, networkDescriptor)
	k.SetOracleWhiteList(ctx, networkDescriptor, append(valList, validator))
}

// RemoveOracleWhiteList remove a validator from whitelist
func (k Keeper) RemoveOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validator sdk.ValAddress) {
	valList := k.GetOracleWhiteList(ctx, networkDescriptor)

	for index, item := range valList {
		if validator.Equals(item) {
			k.SetOracleWhiteList(ctx, networkDescriptor, append(valList[:index], valList[index+1]))
			return
		}
	}
}
