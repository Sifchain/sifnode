package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

// Rename this feature to clp admin list to avoid confusion with Whitelist module

func (keeper Keeper) SetClpWhiteList(ctx sdk.Context, validatorList []sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	key := types.WhiteListValidatorPrefix
	valList := make([]string, len(validatorList))
	for i, entry := range validatorList {
		valList[i] = entry.String()
	}
	store.Set(key, keeper.cdc.MustMarshal(&stakingtypes.ValAddresses{Addresses: valList}))
}

func (keeper Keeper) ExistsClpWhiteList(ctx sdk.Context) bool {
	key := types.WhiteListValidatorPrefix
	return keeper.Exists(ctx, key)
}

func (keeper Keeper) GetClpWhiteList(ctx sdk.Context) []sdk.AccAddress {
	store := ctx.KVStore(keeper.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := store.Get(key)
	valAddresses := &stakingtypes.ValAddresses{}
	keeper.cdc.MustUnmarshal(bz, valAddresses)
	vl := make([]sdk.AccAddress, len(valAddresses.Addresses))
	for i, entry := range valAddresses.Addresses {
		addr, err := sdk.AccAddressFromBech32(entry)
		if err != nil {
			panic(err)
		}
		vl[i] = addr
	}
	return vl
}

func (keeper Keeper) ValidateAddress(ctx sdk.Context, address sdk.AccAddress) bool {
	if !keeper.ExistsClpWhiteList(ctx) {
		return false
	}
	valList := keeper.GetClpWhiteList(ctx)
	for _, validator := range valList {
		if validator.Equals(address) {
			return true
		}
	}
	return false
}
