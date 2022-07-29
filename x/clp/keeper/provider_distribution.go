package keeper

import (
	"encoding/json"
	"strconv"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolRowanMap map[*types.Pool]sdk.Uint
type LpRowanMap map[string]sdk.Uint
type LpPoolMap map[string][]LPPool

func (k Keeper) ProviderDistributionPolicyRun(ctx sdk.Context) {
	a, b, c := k.doProviderDistribution(ctx)
	k.TransferProviderDistribution(ctx, a, b, c)
}

func (k Keeper) doProviderDistribution(ctx sdk.Context) (PoolRowanMap, LpRowanMap, LpPoolMap) {
	blockHeight := ctx.BlockHeight()
	params := k.GetProviderDistributionParams(ctx)
	if params == nil {
		return make(PoolRowanMap), make(LpRowanMap), make(LpPoolMap)
	}

	period := FindProviderDistributionPeriod(blockHeight, params.DistributionPeriods)
	if period == nil {
		return make(PoolRowanMap), make(LpRowanMap), make(LpPoolMap)
	}

	allPools := k.GetPools(ctx)
	return k.CollectProviderDistributions(ctx, allPools, period.DistributionPeriodBlockRate)
}

func (k Keeper) TransferProviderDistribution(ctx sdk.Context, poolRowanMap PoolRowanMap, lpRowanMap LpRowanMap, lpPoolMap LpPoolMap) {
	k.TransferProviderDistributionGeneric(ctx, poolRowanMap, lpRowanMap, lpPoolMap, "lppd/liquidity_provider_payout_error", "lppd/distribution")

	for pool, sub := range poolRowanMap {
		// will never fail
		k.RemoveRowanFromPool(ctx, pool, sub) // nolint:errcheck
	}
}

func (k Keeper) TransferProviderDistributionGeneric(ctx sdk.Context, poolRowanMap PoolRowanMap, lpRowanMap LpRowanMap, lpPoolMap LpPoolMap, typeStr string, successEventType string) {
	for lpAddress, totalRowan := range lpRowanMap {
		addr, _ := sdk.AccAddressFromBech32(lpAddress) // We know this can't fail as we previously filtered out invalid strings
		coin := sdk.NewCoin(types.NativeSymbol, sdk.NewIntFromBigInt(totalRowan.BigInt()))

		//TransferCoinsFromPool(pool, provider_rowan, provider_address)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(coin))
		if err != nil {
			fireLPPayoutErrorEvent(ctx, addr, typeStr, err)

			for _, lpPool := range lpPoolMap[lpAddress] {
				poolRowanMap[lpPool.Pool] = poolRowanMap[lpPool.Pool].Sub(lpPool.Amount)
			}
		} else {
			fireDistributeSuccessEvent(ctx, lpAddress, lpPoolMap[lpAddress], totalRowan, successEventType)
		}
	}
}

func fireDistributeSuccessEvent(ctx sdk.Context, lpAddress string, pools []LPPool, totalDistributed sdk.Uint, typeStr string) {
	data := PrintPools(pools)
	successEvent := sdk.NewEvent(
		typeStr,
		sdk.NewAttribute("recipient", lpAddress),
		sdk.NewAttribute("total_amount", totalDistributed.String()),
		sdk.NewAttribute("amounts", data),
	)

	ctx.EventManager().EmitEvents(sdk.Events{successEvent})
}

type FormattedPool struct {
	Pool   string   `json:"pool"`
	Amount sdk.Uint `json:"amount"`
}

func PrintPools(pools []LPPool) string {
	var formattedPools = make([]FormattedPool, len(pools))

	for i, pool := range pools {
		formattedPools[i] = FormattedPool{Pool: pool.Pool.ExternalAsset.Symbol, Amount: pool.Amount}
	}

	data, _ := json.Marshal(formattedPools) // as used, this should never return an error
	return string(data)
}

func fireLPPayoutErrorEvent(ctx sdk.Context, address sdk.AccAddress, typeStr string, err error) {
	failureEvent := sdk.NewEvent(
		typeStr,
		sdk.NewAttribute("liquidity_provider", address.String()),
		sdk.NewAttribute(types.AttributeKeyError, err.Error()),
		sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	)

	ctx.EventManager().EmitEvents(sdk.Events{failureEvent})
}

//nolint
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

func (k Keeper) CollectProviderDistributions(ctx sdk.Context, pools []*types.Pool, blockRate sdk.Dec) (PoolRowanMap, LpRowanMap, LpPoolMap) {
	poolRowanMap := make(PoolRowanMap, len(pools))
	lpMap := make(LpRowanMap, 0)
	lpPoolMap := make(LpPoolMap, 0)

	partitions, err := k.GetAllLiquidityProvidersPartitions(ctx)
	if err != nil {
		fireLPPGetLPsErrorEvent(ctx, err)
	}

	for _, pool := range pools {
		lps, exists := partitions[*pool.ExternalAsset]
		if !exists { // TODO: fire event
			continue
		}
		lpsFiltered := FilterValidLiquidityProviders(ctx, lps)
		rowanToDistribute := CollectProviderDistribution(ctx, pool, sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()),
			blockRate, pool.PoolUnits, lpsFiltered, lpMap, lpPoolMap)
		poolRowanMap[pool] = rowanToDistribute
	}

	return poolRowanMap, lpMap, lpPoolMap
}

type ValidLiquidityProvider struct {
	Address sdk.AccAddress
	LP      *types.LiquidityProvider
}

type LPPool struct {
	Pool   *types.Pool
	Amount sdk.Uint
}

func PoolRowanMapToLPPools(poolRowanMap PoolRowanMap) []LPPool {
	arr := make([]LPPool, 0, len(poolRowanMap))

	for pool, coins := range poolRowanMap {
		arr = append(arr, LPPool{Pool: pool, Amount: coins})
	}

	return arr
}

func FilterValidLiquidityProviders(ctx sdk.Context, lps []*types.LiquidityProvider) []ValidLiquidityProvider {
	var valid []ValidLiquidityProvider //nolint

	for _, lp := range lps {
		address, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
		if err != nil {
			//k.Logger(ctx).Error(fmt.Sprintf("Liquidity provider address %s error %s", lp.LiquidityProviderAddress, err.Error()))
			fireLPAddressErrorEvent(ctx, lp.LiquidityProviderAddress, err)
			continue
		}

		valid = append(valid, ValidLiquidityProvider{Address: address, LP: lp})
	}

	return valid
}

func fireLPPGetLPsErrorEvent(ctx sdk.Context, err error) {
	failureEvent := sdk.NewEvent(
		"lppd/get_liquidity_providers_error",
		sdk.NewAttribute(types.AttributeKeyError, err.Error()),
		sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	)

	ctx.EventManager().EmitEvents(sdk.Events{failureEvent})
}

func CollectProviderDistribution(ctx sdk.Context, pool *types.Pool, poolDepthRowan, blockRate sdk.Dec, poolUnits sdk.Uint, lps []ValidLiquidityProvider, globalLpRowanMap LpRowanMap, globalLpPoolMap LpPoolMap) sdk.Uint {
	totalRowanDistribute := sdk.ZeroUint()

	//	rowan_provider_distribution = r_block * pool_depth_rowan
	rowanPd := blockRate.Mul(poolDepthRowan)
	rowanPdUint := sdk.NewUintFromBigInt(rowanPd.RoundInt().BigInt())
	for _, lp := range lps {
		providerRowan := CalcProviderDistributionAmount(rowanPd, poolUnits, lp.LP.LiquidityProviderUnits)
		totalRowanDistribute = totalRowanDistribute.Add(providerRowan)

		// TODO: find a proper solution
		if totalRowanDistribute.GT(rowanPdUint) {
			providerRowan = rowanPdUint.Sub(totalRowanDistribute.Sub(providerRowan))
			totalRowanDistribute = rowanPdUint
		}

		addr := lp.Address.String()

		globalLpRowan := globalLpRowanMap[addr]
		if globalLpRowan == (sdk.Uint{}) {
			globalLpRowan = sdk.ZeroUint()
		}
		globalLpRowanMap[addr] = globalLpRowan.Add(providerRowan)

		elem := LPPool{Pool: pool, Amount: providerRowan}
		globalLpPool := globalLpPoolMap[addr]
		if globalLpPool == nil {
			arr := []LPPool{elem}
			globalLpPoolMap[addr] = arr
		} else {
			globalLpPool = append(globalLpPool, elem)
			globalLpPoolMap[addr] = globalLpPool
		}
	}

	return totalRowanDistribute
}

func fireLPAddressErrorEvent(ctx sdk.Context, address string, err error) {
	failureEvent := sdk.NewEvent(
		"lppd/liquidity_provider_address_error",
		sdk.NewAttribute("liquidity_provider", address),
		sdk.NewAttribute(types.AttributeKeyError, err.Error()),
		sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	)

	ctx.EventManager().EmitEvents(sdk.Events{failureEvent})
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
