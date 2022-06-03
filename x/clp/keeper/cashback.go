package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCashbackBlockRate(ctx sdk.Context) sdk.Dec {
	return sdk.NewDecWithPrec(1, 2)
	//k.GetCashbackParams(ctx).BlockRate
}

func (k Keeper) GetCashbackFinalBlock(ctx sdk.Context) int64 {
	return 42
	//k.GetCashbackParams(ctx).BlockRate
}

func CalcCashbackAmount(rowanCashedback sdk.Dec, totalPoolUnits, providerPoolUnits sdk.Uint) sdk.Uint {
	//provider_percentage = provider_units / total_pool_units
	totalPoolUnitsDec := sdk.NewDecFromBigInt(totalPoolUnits.BigInt())
	providerPercentage := sdk.NewDecFromBigInt(providerPoolUnits.BigInt()).Quo(totalPoolUnitsDec)

	//provider_rowan = provider_percentage * rowan_cashbacked
	providerRowan := providerPercentage.Mul(rowanCashedback)

	return sdk.Uint(providerRowan.RoundInt())
}

func (k Keeper) payOutLPs(ctx sdk.Context, rowanCashbacked sdk.Dec, totalPoolUnits sdk.Uint, lp *types.LiquidityProvider) error {
	providerRowan := CalcCashbackAmount(rowanCashbacked, totalPoolUnits, lp.LiquidityProviderUnits)

	//TransferCoinsFromPool(pool, provider_rowan, provider_address)
	coin := sdk.NewCoin(types.NativeSymbol, sdk.Int(providerRowan))
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(lp.LiquidityProviderAddress), sdk.NewCoins(coin))
}

func (k Keeper) doCashback(ctx sdk.Context, pools []*types.Pool) error {
	blockRate := k.GetCashbackBlockRate(ctx)
	for _, pool := range pools {
		lps, err := k.GetAllLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
		if err != nil {
			// Ignore and continue for the rest of the pools?
			continue
			//return err
		}

		//	rowan_cashbacked = r_block * pool_depth_rowan
		rowanCashbacked := blockRate.Mul(sdk.NewDecFromBigInt(pool.NativeAssetBalance.BigInt()))
		for _, lp := range lps {
			err = k.payOutLPs(ctx, rowanCashbacked, pool.PoolUnits, lp)
			if err != nil {
				// Ignore and continue for the rest of the pools?
				//return sdkerrors.Wrap(types.ErrUnableToAddBalance, err.Error())
			}
		}
	}

	return nil
}

func (k Keeper) CashbackPolicyRun(ctx sdk.Context) error {
	currentHeight := ctx.BlockHeight()
	finalBlock := k.GetCashbackFinalBlock(ctx)

	if currentHeight > finalBlock {
		// Log
		return nil
	}

	allPools := k.GetPools(ctx)
	return k.doCashback(ctx, allPools)
}