package keeper

import (
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetNativeToken set the validator list for a network.
func (k Keeper) SetNativeToken(ctx sdk.Context, networkDescriptor types.NetworkIdentity, token string) {
	store := ctx.KVStore(k.storeKey)
	key := networkDescriptor.GetNativeTokenPrefix()
	nativeToken := types.NativeToken{NativeToken: token}
	store.Set(key, k.cdc.MustMarshalBinaryBare(&nativeToken))
}

// GetNativeToken return validator list
func (k Keeper) GetNativeToken(ctx sdk.Context, networkIdentity types.NetworkIdentity) string {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetNativeTokenPrefix()

	store.Has(key)

	bz := store.Get(key)
	nativeToken := &types.NativeToken{}
	k.cdc.MustUnmarshalBinaryBare(bz, nativeToken)
	return nativeToken.NativeToken
}
