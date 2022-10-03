package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/instrumentation"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetCrossChainFee set the crosschain fee for a network.
func (k Keeper) SetCrossChainFee(ctx sdk.Context,
	networkIdentity types.NetworkIdentity,
	token string,
	gas, lockCost, burnCost, firstBurnDoublePeggyCost sdk.Int) {

	crossChainFee := types.CrossChainFeeConfig{
		FeeCurrency:              token,
		FeeCurrencyGas:           gas,
		MinimumLockCost:          lockCost,
		MinimumBurnCost:          burnCost,
		FirstBurnDoublePeggyCost: firstBurnDoublePeggyCost,
	}

	k.SetCrossChainFeeObj(ctx, networkIdentity, &crossChainFee)
}

// SetCrossChainFeeObj set the crosschain fee object for a network.
func (k Keeper) SetCrossChainFeeObj(ctx sdk.Context,
	networkIdentity types.NetworkIdentity,
	crossChainFee *types.CrossChainFeeConfig) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetCrossChainFeePrefix(k.cdc)
	store.Set(key, k.cdc.MustMarshal(crossChainFee))
}

// GetCrossChainFeeConfig return crosschain fee config
func (k Keeper) GetCrossChainFeeConfig(ctx sdk.Context, networkIdentity types.NetworkIdentity) (types.CrossChainFeeConfig, error) {
	store := ctx.KVStore(k.storeKey)
	key := networkIdentity.GetCrossChainFeePrefix(k.cdc)

	if !store.Has(key) {
		return types.CrossChainFeeConfig{}, fmt.Errorf("%s%s", "crosschain fee not set for ", networkIdentity.NetworkDescriptor.String())
	}

	bz := store.Get(key)
	crossChainFeeConfig := &types.CrossChainFeeConfig{}
	k.cdc.MustUnmarshal(bz, crossChainFeeConfig)

	instrumentation.PeggyCheckpoint(ctx.Logger(), instrumentation.GetCrossChainFeeConfig, "crossChainFeeConfig", zap.Reflect("crossChainFeeConfig", crossChainFeeConfig))

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

// GetAllCrossChainFeeConfig get all fee configs for all network descriptors
func (k Keeper) GetAllCrossChainFeeConfig(ctx sdk.Context) []*types.GenesisCrossChainFeeConfig {
	configs := make([]*types.GenesisCrossChainFeeConfig, 0)
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.CrossChainFeePrefix)

	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()

		var config types.GenesisCrossChainFeeConfig

		network_identity, err := types.GetFromPrefix(k.cdc, key, types.CrossChainFeePrefix)
		if err != nil {
			panic(err)
		}
		// k.cdc.MustUnmarshal(value, &config)

		var crossChainFeeConfig types.CrossChainFeeConfig
		k.cdc.MustUnmarshal(value, &crossChainFeeConfig)

		config.NetworkDescriptor = network_identity.NetworkDescriptor
		config.CrossChainFee = &crossChainFeeConfig

		configs = append(configs, &config)
	}
	return configs
}
