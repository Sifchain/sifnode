package keeper

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) SetOracleWhiteList(ctx sdk.Context, validatorList []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	valList := make([]string, 0, len(validatorList))
	for i, entry := range validatorList {
		valList[i] = entry.String()
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(&stakingtypes.ValAddresses{valList}))
}

func (k Keeper) ExistsOracleWhiteList(ctx sdk.Context) bool {
	key := types.WhiteListValidatorPrefix
	return k.Exists(ctx, key)
}

//
func (k Keeper) GetOracleWhiteList(ctx sdk.Context) []sdk.ValAddress {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := store.Get(key)
	valList := []string{}
	k.cdc.MustUnmarshalBinaryBare(bz, &stakingtypes.ValAddresses{valList})

	vl := make([]sdk.ValAddress, len(valList))
	for i, entry := range valList {
		addr, err := sdk.ValAddressFromBech32(entry)
		if err != nil {
			panic(err)
		}
		vl[i] = addr

	}
	return vl
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
