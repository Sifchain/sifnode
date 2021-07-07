package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetNativeToken set the validator list for a network.
func (k Keeper) SetNativeToken(ctx sdk.Context, networkIdentity types.NetworkIdentity, token string) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetNativeTokenPrefix()
	fmt.Printf("+++ SetNativeToken key is %v\n", key)
	nativeToken := types.NativeToken{NativeToken: token}
	store.Set(key, k.cdc.MustMarshalBinaryBare(&nativeToken))
}

// GetNativeToken return validator list
func (k Keeper) GetNativeToken(ctx sdk.Context, networkIdentity types.NetworkIdentity) (string, error) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetNativeTokenPrefix()

	if !store.Has(key) {
		return "", fmt.Errorf("%s%s", "native token not set for ", networkIdentity.NetworkDescriptor.String())
	}

	bz := store.Get(key)
	nativeToken := &types.NativeToken{}
	k.cdc.MustUnmarshalBinaryBare(bz, nativeToken)
	return nativeToken.NativeToken, nil
}
