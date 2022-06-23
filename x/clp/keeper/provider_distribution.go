package keeper

import (
	"fmt"
	"strconv"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProviderDistributionMap map[string]sdk.Uint

func (k Keeper) ProviderDistributionPolicyRun(ctx sdk.Context) {
	pdm := k.doProviderDistribution(ctx)
	k.TranferCoins(ctx, &pdm)
}

func (k Keeper) TranferCoins(ctx sdk.Context, pdm *ProviderDistributionMap) {
	for lpAddress, pdRowan := range *pdm {
		address, err := sdk.AccAddressFromBech32(lpAddress)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Liquidity provider address %s error %s", lpAddress, err.Error()))
			continue
		}

		err = k.TransferProviderDistribution(ctx, address, pdRowan)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Paying out liquidity provider %s error %s", address, err.Error()))
		}
	}
}

func (k Keeper) doProviderDistribution(ctx sdk.Context) ProviderDistributionMap {
	blockHeight := ctx.BlockHeight()
	params := k.GetProviderDistributionParams(ctx)
	if params == nil {
		return make(ProviderDistributionMap)
	}

	period := FindProviderDistributionPeriod(blockHeight, params.DistributionPeriods)
	if period == nil {
		return make(ProviderDistributionMap)
	}

	allPools := k.GetPools(ctx)
	return k.CollectProviderDistributions(ctx, allPools, period.DistributionPeriodBlockRate)
}

func (k Keeper) TransferProviderDistribution(ctx sdk.Context, providerAddress sdk.AccAddress, providerRowan sdk.Uint) error {
	//TransferCoinsFromPool(pool, provider_rowan, provider_address)
	coin := sdk.NewCoin(types.NativeSymbol, sdk.Int(providerRowan))
	fireDistributionEvent(ctx, coin, providerAddress)

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, providerAddress, sdk.NewCoins(coin))
}

func fireDistributionEvent(ctx sdk.Context, coin sdk.Coin, to sdk.Address) {
	distribtionEvent := sdk.NewEvent(
		types.EventTypeProviderDistributionDistribution,
		sdk.NewAttribute(types.AttributeProbiverDistributionAmount, coin.String()),
		sdk.NewAttribute(types.AttributeProbiverDistributionReceiver, to.String()),
		sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	)

	ctx.EventManager().EmitEvents(sdk.Events{distribtionEvent})
}

func FindProviderDistributionPeriod(currentHeight int64, periods []*types.ProviderDistributionPeriod) *types.ProviderDistributionPeriod {
	for _, period := range periods {
		if isActivePeriod(currentHeight, period.DistributionPeriodStartBlock, period.DistributionPeriodEndBlock) {
			return period
		}
	}

	return nil
}

func isActivePeriod(current, start, end int64) bool {
	return current >= start && current <= end
}

func (k Keeper) CollectProviderDistributions(ctx sdk.Context, pools []*types.Pool, blockRate sdk.Dec) ProviderDistributionMap {
	m := make(ProviderDistributionMap)

	for _, pool := range pools {
		lps, err := k.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Getting liquidity providers for asset %s error %s", pool.ExternalAsset.Symbol, err.Error()))
			continue
		}

		CollectProviderDistribution(sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()),
			blockRate, pool.PoolUnits, lps, m)
	}

	return m
}

func CollectProviderDistribution(poolDepthRowan, blockRate sdk.Dec, poolUnits sdk.Uint, lps []*types.LiquidityProvider, cbm ProviderDistributionMap) {
	//	rowan_provider_distribution = r_block * pool_depth_rowan
	rowanPd := blockRate.Mul(poolDepthRowan)
	for _, lp := range lps {
		providerRowan := CalcProviderDistributionAmount(rowanPd, poolUnits, lp.LiquidityProviderUnits)
		rowanSoFar := cbm[lp.LiquidityProviderAddress]
		if rowanSoFar == (sdk.Uint{}) { // sdk.Uint{} seems to be the default value instead of zero...
			rowanSoFar = sdk.ZeroUint()
		}
		cbm[lp.LiquidityProviderAddress] = rowanSoFar.Add(providerRowan)
	}
}

func CalcProviderDistributionAmount(rowanProviderDistribution sdk.Dec, totalPoolUnits, providerPoolUnits sdk.Uint) sdk.Uint {
	//provider_percentage = provider_units / total_pool_units
	providerPercentage := sdk.NewDecFromBigInt(providerPoolUnits.BigInt()).Quo(sdk.NewDecFromBigInt(totalPoolUnits.BigInt()))

	//provider_rowan = provider_percentage * rowan_provider_distribution
	providerRowan := providerPercentage.Mul(rowanProviderDistribution)

	return sdk.Uint(providerRowan.RoundInt())
}

func (k Keeper) SetProviderDistributionParams(ctx sdk.Context, params *types.ProviderDistributionParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProviderDistributionParamsPrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetProviderDistributionParams(ctx sdk.Context) *types.ProviderDistributionParams {
	params := types.ProviderDistributionParams{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProviderDistributionParamsPrefix)
	k.cdc.MustUnmarshal(bz, &params)

	return &params
}
