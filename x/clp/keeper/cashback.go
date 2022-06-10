package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CashbackMap map[string]sdk.Uint

func (k Keeper) CashbackPolicyRun(ctx sdk.Context) {
	cashbackMap := k.doCashback(ctx)
	for lpAddress, cashbackRowan := range cashbackMap {
		address, err := sdk.AccAddressFromBech32(lpAddress)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Liquidity provider address %s error %s", lpAddress, err.Error()))
			continue
		}

		err = k.transferCashback(ctx, address, cashbackRowan)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Paying out liquidity provider %s error %s", address, err.Error()))
		}
	}
}

func (k Keeper) doCashback(ctx sdk.Context) CashbackMap {
	blockHeight := ctx.BlockHeight()
	params := k.GetCashbackParams(ctx)
	if params == nil {
		return make(CashbackMap)
	}

	period := FindActiveCashbackPeriod(blockHeight, params.CashbackPeriods)
	if period == nil {
		return make(CashbackMap)
	}

	allPools := k.GetPools(ctx)
	return k.CollectCashbacks(ctx, allPools, period.CashbackPeriodBlockRate)
}

func (k Keeper) transferCashback(ctx sdk.Context, providerAddress sdk.AccAddress, providerRowan sdk.Uint) error {
	//TransferCoinsFromPool(pool, provider_rowan, provider_address)
	coin := sdk.NewCoin(types.NativeSymbol, sdk.Int(providerRowan))
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, providerAddress, sdk.NewCoins(coin))
}

func FindActiveCashbackPeriod(currentHeight int64, periods []*types.CashbackPeriod) *types.CashbackPeriod {
	for _, period := range periods {
		if isActivePeriod(currentHeight, period.CashbackPeriodStartBlock, period.CashbackPeriodEndBlock) {
			return period
		}
	}

	return nil
}

func isActivePeriod(current, start, end int64) bool {
	return current >= start && current <= end
}

func (k Keeper) CollectCashbacks(ctx sdk.Context, pools []*types.Pool, blockRate sdk.Dec) CashbackMap {
	m := make(CashbackMap)

	for _, pool := range pools {
		lps, err := k.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Getting liquidity providers for asset %s error %s", pool.ExternalAsset.Symbol, err.Error()))
			continue
		}

		CollectCashback(sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()),
			blockRate, pool.PoolUnits, lps, m)
	}

	return m
}

func CollectCashback(poolDepthRowan, blockRate sdk.Dec, poolUnits sdk.Uint, lps []*types.LiquidityProvider, cbm CashbackMap) {
	//	rowan_cashbacked = r_block * pool_depth_rowan
	rowanCashbacked := blockRate.Mul(poolDepthRowan)
	for _, lp := range lps {
		providerRowan := CalcCashbackAmount(rowanCashbacked, poolUnits, lp.LiquidityProviderUnits)
		rowanSoFar := cbm[lp.LiquidityProviderAddress]
		if rowanSoFar == (sdk.Uint{}) { // sdk.Uint{} seems to be the default value instead of zero... lol
			rowanSoFar = sdk.ZeroUint()
		}
		cbm[lp.LiquidityProviderAddress] = rowanSoFar.Add(providerRowan)
	}
}

func CalcCashbackAmount(rowanCashedback sdk.Dec, totalPoolUnits, providerPoolUnits sdk.Uint) sdk.Uint {
	//provider_percentage = provider_units / total_pool_units
	providerPercentage := sdk.NewDecFromBigInt(providerPoolUnits.BigInt()).Quo(sdk.NewDecFromBigInt(totalPoolUnits.BigInt()))

	//provider_rowan = provider_percentage * rowan_cashbacked
	providerRowan := providerPercentage.Mul(rowanCashedback)

	return sdk.Uint(providerRowan.RoundInt())
}

func (k Keeper) SetCashbackParams(ctx sdk.Context, params *types.CashbackParams) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CashbackParamsPrefix, k.cdc.MustMarshal(params))
}

func (k Keeper) GetCashbackParams(ctx sdk.Context) *types.CashbackParams {
	params := types.CashbackParams{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CashbackParamsPrefix)
	k.cdc.MustUnmarshal(bz, &params)

	return &params
}
