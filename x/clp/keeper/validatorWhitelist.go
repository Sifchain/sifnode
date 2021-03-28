package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func (k Keeper) SetClpWhiteList(ctx sdk.Context, validatorList []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	valList := make([]string, 0, len(validatorList))
	for i, entry := range validatorList {
		valList[i] = entry.String()
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(&stakingtypes.ValAddresses{valList}))
}

func (k Keeper) ExistsClpWhiteList(ctx sdk.Context) bool {
	key := types.WhiteListValidatorPrefix
	return k.Exists(ctx, key)
}

func (k Keeper) GetClpWhiteList(ctx sdk.Context) []sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	key := types.WhiteListValidatorPrefix
	bz := store.Get(key)

	valList := []string{}
	k.cdc.MustUnmarshalBinaryBare(bz, &stakingtypes.ValAddresses{valList})

	vl := make([]sdk.AccAddress, len(valList))
	for i, entry := range valList {
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
