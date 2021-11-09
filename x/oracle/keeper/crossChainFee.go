package keeper

import (
	"fmt"
	"go.uber.org/zap"

	ethbridgeTypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetCrossChainFee set the crosschain fee for a network.
func (k Keeper) SetCrossChainFee(ctx sdk.Context,
	networkIdentity types.NetworkIdentity,
	token string,
	gas, lockCost, burnCost, firstLockDoublePeggyCost sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetCrossChainFeePrefix()
	crossChainFee := types.CrossChainFeeConfig{
		FeeCurrency:              token,
		FeeCurrencyGas:           gas,
		MinimumLockCost:          lockCost,
		MinimumBurnCost:          burnCost,
		FirstLockDoublePeggyCost: firstLockDoublePeggyCost,
	}
	store.Set(key, k.cdc.MustMarshalBinaryBare(&crossChainFee))
}

// GetCrossChainFeeConfig return crosschain fee config
func (k Keeper) GetCrossChainFeeConfig(ctx sdk.Context, networkIdentity types.NetworkIdentity) (types.CrossChainFeeConfig, error) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetCrossChainFeePrefix()

	if !store.Has(key) {
		return types.CrossChainFeeConfig{}, fmt.Errorf("%s%s", "crosschain fee not set for ", networkIdentity.NetworkDescriptor.String())
	}

	bz := store.Get(key)
	crossChainFeeConfig := &types.CrossChainFeeConfig{}
	k.cdc.MustUnmarshalBinaryBare(bz, crossChainFeeConfig)

	ctx.Logger().Debug(ethbridgeTypes.PeggyTestMarker, "kind", "GetCrossChainFeeConfig", "crossChainFeeConfig", zap.Reflect("crossChainFeeConfig", crossChainFeeConfig))

	return *crossChainFeeConfig, nil
}

// GetCrossChainFee return crosschain fee
func (k Keeper) GetCrossChainFee(ctx sdk.Context, networkIdentity types.NetworkIdentity) (string, error) {
	crossChainFeeConfig, err := k.GetCrossChainFeeConfig(ctx, networkIdentity)
	if err != nil {
		return "", err
	}

	return crossChainFeeConfig.FeeCurrency, nil
}
