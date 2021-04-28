package keeper

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validatorList types.ValidatorWhitelist) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	store.Set(key, k.Cdc.MustMarshalBinaryBare(validatorList))
}

func (k Keeper) ExistsOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) bool {
	key := networkDescriptor.GetPrefix()
	return k.Exists(ctx, key)
}

//
func (k Keeper) GetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) (valList types.ValidatorWhitelist) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	bz := store.Get(key)
	k.Cdc.MustUnmarshalBinaryBare(bz, &valList)
	return
}

// ValidateAddress is a validator in whitelist
func (k Keeper) ValidateAddress(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, address sdk.ValAddress) bool {
	if !k.ExistsOracleWhiteList(ctx, networkDescriptor) {
		return false
	}
	valList := k.GetOracleWhiteList(ctx, networkDescriptor)
	return valList.ContainValidator(address)
}

// UpdateOracleWhiteList validator's power
func (k Keeper) UpdateOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validator sdk.ValAddress, power uint32) {
	valList := k.GetOracleWhiteList(ctx, networkDescriptor)
	valList.AddValidator(validator, power)
	k.SetOracleWhiteList(ctx, networkDescriptor, valList)
}
