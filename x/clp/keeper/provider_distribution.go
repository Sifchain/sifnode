package keeper

import (
	"fmt"
	"strconv"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DistributionTuple struct {
	Amount          sdk.Uint
	ProviderAddress sdk.AccAddress
}

type PoolMap map[*types.Pool]([]DistributionTuple)

func (k Keeper) ProviderDistributionPolicyRun(ctx sdk.Context) {
	poolMap := k.doProviderDistribution(ctx)
	for pool, tuples := range poolMap {
		for _, tuple := range tuples {
			err := k.TransferProviderDistribution(ctx, pool, &tuple)
			if err != nil {
				k.Logger(ctx).Error(fmt.Sprintf("Paying out liquidity provider %s error %s", tuple.ProviderAddress, err.Error()))
			}
		}
	}
}

func (k Keeper) doProviderDistribution(ctx sdk.Context) PoolMap {
	blockHeight := ctx.BlockHeight()
	params := k.GetProviderDistributionParams(ctx)
	if params == nil {
		return make(PoolMap)
	}

	period := FindProviderDistributionPeriod(blockHeight, params.DistributionPeriods)
	if period == nil {
		return make(PoolMap)
	}

	allPools := k.GetPools(ctx)
	return k.CollectProviderDistributions(ctx, allPools, period.DistributionPeriodBlockRate)
}

func (k Keeper) TransferProviderDistribution(ctx sdk.Context, pool *types.Pool, tuple *DistributionTuple) error {
	//TransferCoinsFromPool(pool, provider_rowan, provider_address)
	err := k.SendRowanFromPool(ctx, pool, tuple.Amount, tuple.ProviderAddress)
	if err != nil {
		// TODO fire failure event
		return err
	}

	fireDistributionEvent(ctx, tuple.Amount, tuple.ProviderAddress)

	return nil
}

func fireDistributionEvent(ctx sdk.Context, amount sdk.Uint, to sdk.Address) {
	coin := sdk.NewCoin(types.NativeSymbol, sdk.NewIntFromBigInt(amount.BigInt()))
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

func isActivePeriod(current int64, start, end uint64) bool {
	return current >= int64(start) && current <= int64(end)
}

func (k Keeper) CollectProviderDistributions(ctx sdk.Context, pools []*types.Pool, blockRate sdk.Dec) PoolMap {
	poolMap := make(PoolMap, len(pools))

	for _, pool := range pools {
		lps, err := k.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Getting liquidity providers for asset %s error %s", pool.ExternalAsset.Symbol, err.Error()))
			continue
		}

		tuples := CollectProviderDistribution(sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()),
			blockRate, pool.PoolUnits, lps)
		poolMap[pool] = tuples
	}

	return poolMap
}

func CollectProviderDistribution(poolDepthRowan, blockRate sdk.Dec, poolUnits sdk.Uint, lps []*types.LiquidityProvider) []DistributionTuple {
	tuples := make([]DistributionTuple, len(lps))

	//	rowan_provider_distribution = r_block * pool_depth_rowan
	rowanPd := blockRate.Mul(poolDepthRowan)
	for _, lp := range lps {
		address, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
		if err != nil {
			// TODO: collect and return
			//k.Logger(ctx).Error(fmt.Sprintf("Liquidity provider address %s error %s", lp.LiquidityProviderAddress, err.Error()))
			continue
		}

		providerRowan := CalcProviderDistributionAmount(rowanPd, poolUnits, lp.LiquidityProviderUnits)
		tuples = append(tuples, DistributionTuple{Amount: providerRowan, ProviderAddress: address})
	}

	return tuples
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

func (k Keeper) IsDistributionBlock(ctx sdk.Context) bool {
	blockHeight := ctx.BlockHeight()
	params := k.GetProviderDistributionParams(ctx)
	period := FindProviderDistributionPeriod(blockHeight, params.DistributionPeriods)
	if period == nil {
		return false
	}

	startHeight := period.DistributionPeriodStartBlock
	mod := period.DistributionPeriodMod

	return IsDistributionBlockPure(blockHeight, startHeight, mod)
}

// do the thing every mod blocks starting at startHeight
func IsDistributionBlockPure(blockHeight int64, startHeight, mod uint64) bool {
	return (blockHeight-int64(startHeight))%int64(mod) == 0
}
