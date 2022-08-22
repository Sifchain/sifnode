//go:build FEATURE_TOGGLE_SDK_045
// +build FEATURE_TOGGLE_SDK_045

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/oracle/types"
)

func (k Keeper) GetProphecyIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.ProphecyPrefix)
}
