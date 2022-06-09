package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CashbackMap map[string]sdk.Uint

func (k Keeper) CashbackPolicyRun(ctx sdk.Context) {
	cashbackMap := k.DoCashback(ctx)
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

func (k Keeper) DoCashback(ctx sdk.Context) CashbackMap {
	currentHeight := ctx.BlockHeight()
	period := k.findValidCashbackPeriod(ctx, currentHeight)
	if period == nil {
		return make(CashbackMap)
	}

	allPools := k.GetPools(ctx)
	return k.collectCashbacks(ctx, allPools, period.CashbackPeriodBlockRate)
}

func (k Keeper) transferCashback(ctx sdk.Context, providerAddress sdk.AccAddress, providerRowan sdk.Uint) error {
	//TransferCoinsFromPool(pool, provider_rowan, provider_address)
	coin := sdk.NewCoin(types.NativeSymbol, sdk.Int(providerRowan))
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, providerAddress, sdk.NewCoins(coin))
}

func (k Keeper) findValidCashbackPeriod(ctx sdk.Context, currentHeight int64) *types.CashbackPeriod {
	params := k.GetCashbackParams(ctx)
	for _, period := range params.CashbackPeriods {
		if isActivePeriod(currentHeight, period.CashbackPeriodStartBlock, period.CashbackPeriodEndBlock) {
			return period
		}
	}

	return nil
}

func isActivePeriod(current, start, end int64) bool {
	return start >= current && end <= current
}

func (k Keeper) collectCashbacks(ctx sdk.Context, pools []*types.Pool, blockRate sdk.Dec) CashbackMap {
	m := make(CashbackMap)

	for _, pool := range pools {
		lps, err := k.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Getting liquidity providers for asset %s error %s", pool.ExternalAsset.Symbol, err.Error()))
			continue
		}

		//	rowan_cashbacked = r_block * pool_depth_rowan
		rowanCashbacked := blockRate.Mul(sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()))
		for _, lp := range lps {
			providerRowan := CalcCashbackAmount(rowanCashbacked, pool.PoolUnits, lp.LiquidityProviderUnits)
			rowanSoFar := m[lp.LiquidityProviderAddress]
			m[lp.LiquidityProviderAddress] = rowanSoFar.Add(providerRowan)
		}
	}

	return m
}

func CalcCashbackAmount(rowanCashedback sdk.Dec, totalPoolUnits, providerPoolUnits sdk.Uint) sdk.Uint {
	//provider_percentage = provider_units / total_pool_units
	totalPoolUnitsDec := sdk.NewDecFromBigInt(totalPoolUnits.BigInt())
	providerPercentage := sdk.NewDecFromBigInt(providerPoolUnits.BigInt()).Quo(totalPoolUnitsDec)

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
