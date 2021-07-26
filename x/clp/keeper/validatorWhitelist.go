package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

// Rename this feature to clp admin list to avoid confusion with Whitelist module

func (k Keeper) SetClpWhiteList(ctx sdk.Context, validatorList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	valList := make([]string, len(validatorList))
	for i, entry := range validatorList {
		valList[i] = entry.String()
	}

	store.Set(key, k.cdc.MustMarshalBinaryBare(&stakingtypes.ValAddresses{Addresses: valList}))
}

func (k Keeper) ExistsClpWhiteList(ctx sdk.Context) bool {
	key := types.WhiteListValidatorPrefix
	return k.Exists(ctx, key)
}

func (k Keeper) GetClpWhiteList(ctx sdk.Context) []sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := store.Get(key)

	valAddresses := &stakingtypes.ValAddresses{}
	k.cdc.MustUnmarshalBinaryBare(bz, valAddresses)

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
