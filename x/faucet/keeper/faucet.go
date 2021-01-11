package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetBalance(ctx sdk.Context) error {
	err := ctx.KVStore(k.storeKey)
	if err != nil {
		return nil
	}
	return nil
}
