package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetOracleWhiteList
func (k Keeper) SetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validatorList types.ValidatorWhitelist) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	whitelist, err := json.Marshal(validatorList.Whitelist)
	if err != nil {
		panic("whitelist data format is wrong")
	}
	store.Set(key, k.Cdc.MustMarshalBinaryBare(whitelist))
}

func (k Keeper) ExistsOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) bool {
	key := networkDescriptor.GetPrefix()
	return k.Exists(ctx, key)
}

// GetOracleWhiteList
func (k Keeper) GetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) types.ValidatorWhitelist {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	bz := store.Get(key)
	var whitelistByte []byte
	k.Cdc.MustUnmarshalBinaryBare(bz, &whitelistByte)

	var whitelistMap map[string]uint32
	err := json.Unmarshal(whitelistByte, &whitelistMap)
	if err != nil {
		panic("whitelist data format is wrong")
	}
	fmt.Printf("whittelist is %v\n", whitelistMap)
	return types.NewValidatorWhitelistFromData(whitelistMap)
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
	valList.UpdateValidator(validator, power)
	k.SetOracleWhiteList(ctx, networkDescriptor, valList)
}
