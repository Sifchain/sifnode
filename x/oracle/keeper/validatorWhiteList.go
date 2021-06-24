package keeper

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllWhiteList set the validator list for a network.
func (k Keeper) GetAllWhiteList(ctx sdk.Context) map[types.NetworkID]types.ValidatorWhiteList {
	result := make(map[types.NetworkID]types.ValidatorWhiteList)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WhiteListValidatorPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		networkDescriptor, err := types.GetFromPrefix(iterator.Key())
		if err != nil {
			panic(err.Error())
		}

		result[networkDescriptor.NetworkID] = k.GetOracleWhiteList(ctx, networkDescriptor)
	}

	return result
}

// SetOracleWhiteList set the validator list for a network.
func (k Keeper) SetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validatorList types.ValidatorWhiteList) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	store.Set(key, k.cdc.MustMarshalBinaryBare(&validatorList))
}

// RemoveOracleWhiteList remove the validator list for a network.
func (k Keeper) RemoveOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	store.Delete(key)
}

// ExistsOracleWhiteList check if the key exist
func (k Keeper) ExistsOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) bool {
	key := networkDescriptor.GetPrefix()
	return k.Exists(ctx, key)
}

// GetOracleWhiteList return validator list
func (k Keeper) GetOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) types.ValidatorWhiteList {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetPrefix()
	bz := store.Get(key)
	validators := &types.ValidatorWhiteList{}
	k.cdc.MustUnmarshalBinaryBare(bz, validators)
	return *validators
}

// GetAllValidators return validator list
func (k Keeper) GetAllValidators(ctx sdk.Context, networkDescriptor types.NetworkDescriptor) []sdk.ValAddress {
	valAddresses := k.GetOracleWhiteList(ctx, networkDescriptor)

	vl := []sdk.ValAddress{}
	for i, power := range valAddresses.GetWhiteList() {
		addr, err := sdk.ValAddressFromBech32(i)
		if err != nil {
			panic(err)
		}
		if power > 0 {
			vl = append(vl, addr)
		}
	}

	return vl
}

// ValidateAddress is a validator in whitelist
func (k Keeper) ValidateAddress(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, address sdk.ValAddress) bool {
	if !k.ExistsOracleWhiteList(ctx, networkDescriptor) {
		return false
	}
	valAddresses := k.GetOracleWhiteList(ctx, networkDescriptor)

	for i, power := range valAddresses.GetWhiteList() {
		addr, err := sdk.ValAddressFromBech32(i)
		if err != nil {
			panic(err)
		}
		if power > 0 && addr.Equals(address) {
			return true
		}
	}

	return false
}

// UpdateOracleWhiteList validator's power
func (k Keeper) UpdateOracleWhiteList(ctx sdk.Context, networkDescriptor types.NetworkDescriptor, validator sdk.ValAddress, power uint32) {
	valList := k.GetOracleWhiteList(ctx, networkDescriptor)
	whiteList := valList.GetWhiteList()
	if whiteList == nil {
		whiteList = make(map[string]uint32)
	}
	whiteList[validator.String()] = power

	valList = types.ValidatorWhiteList{WhiteList: whiteList}
	k.SetOracleWhiteList(ctx, networkDescriptor, valList)
}
