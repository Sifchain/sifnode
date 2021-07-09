package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetNativeToken set the native token for a network.
func (k Keeper) SetNativeToken(ctx sdk.Context, networkIdentity types.NetworkIdentity, token string, gas, lockCost, burnCost sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetNativeTokenPrefix()
	nativeToken := types.NativeTokenConfig{
		NativeToken:     token,
		NativeGas:       gas,
		MinimumLockCost: lockCost,
		MinimumBurnCost: burnCost,
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(&nativeToken))
}

// GetNativeTokenConfig return native token config
func (k Keeper) GetNativeTokenConfig(ctx sdk.Context, networkIdentity types.NetworkIdentity) (types.NativeTokenConfig, error) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetNativeTokenPrefix()

	if !store.Has(key) {
		return types.NativeTokenConfig{}, fmt.Errorf("%s%s", "native token not set for ", networkIdentity.NetworkDescriptor.String())
	}

	bz := store.Get(key)
	nativeTokenConfig := &types.NativeTokenConfig{}
	k.cdc.MustUnmarshalBinaryBare(bz, nativeTokenConfig)
	return *nativeTokenConfig, nil
}

// GetNativeToken return native token
func (k Keeper) GetNativeToken(ctx sdk.Context, networkIdentity types.NetworkIdentity) (string, error) {
	nativeTokenConfig, err := k.GetNativeTokenConfig(ctx, networkIdentity)
	if err != nil {
		return "", err
	}

	return nativeTokenConfig.NativeToken, nil
}
