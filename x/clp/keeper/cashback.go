package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CashbackPolicyRun(ctx sdk.Context) error {
	currentHeight := ctx.BlockHeight()
	period := k.findValidCashbackPeriod(ctx, currentHeight)
	if period == nil {
		return nil
	}

	allPools := k.GetPools(ctx)
	return k.doCashback(ctx, allPools, period.CashbackPeriodBlockRate)
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

func (k Keeper) doCashback(ctx sdk.Context, pools []*types.Pool, blockRate sdk.Dec) error {
	for _, pool := range pools {
		lps, err := k.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("Getting liquidity providers for asset %s error %s", pool.ExternalAsset.Symbol, err.Error()))
			continue
		}

		//	rowan_cashbacked = r_block * pool_depth_rowan
		rowanCashbacked := blockRate.Mul(sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()))
		for _, lp := range lps {
			err = k.payOutLPs(ctx, rowanCashbacked, pool.PoolUnits, lp)
			if err != nil {
				k.Logger(ctx).Error(fmt.Sprintf("Paying out liquidity provider %s for asset %s error %s", lp.LiquidityProviderAddress, pool.ExternalAsset.Symbol, err.Error()))
			}
		}
	}

	return nil
}

func (k Keeper) payOutLPs(ctx sdk.Context, rowanCashbacked sdk.Dec, totalPoolUnits sdk.Uint, lp *types.LiquidityProvider) error {
	address, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	if err != nil {
		return err
	}

	providerRowan := CalcCashbackAmount(rowanCashbacked, totalPoolUnits, lp.LiquidityProviderUnits)

	//TransferCoinsFromPool(pool, provider_rowan, provider_address)
	coin := sdk.NewCoin(types.NativeSymbol, sdk.Int(providerRowan))
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, sdk.NewCoins(coin))
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
