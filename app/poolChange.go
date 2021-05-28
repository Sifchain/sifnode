package app

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func GetPoolChangeFunc(app *SifchainApp) func(ctx sdk.Context, plan upgradetypes.Plan) {
	return func(ctx sdk.Context, plan upgradetypes.Plan) {
		ctx.Logger().Info("Starting to execute upgrade plan for pool re-balance")

		ExportAppState("changePoolFormula", app, ctx)

		allPools := app.ClpKeeper.GetPools(ctx)
		lps := types.LiquidityProviders{}
		poolList := types.Pools{}
		hasError := false
		for _, pool := range allPools {
			lpList := app.ClpKeeper.GetLiquidityProvidersForAsset(ctx, *pool.ExternalAsset)
			temp := sdk.ZeroUint()
			tempExternal := sdk.ZeroUint()
			tempNative := sdk.ZeroUint()
			for _, lp := range lpList {
				withdrawNativeAssetAmount, withdrawExternalAssetAmount, _, _ := keeper.CalculateWithdrawal(pool.PoolUnits, pool.NativeAssetBalance.String(),
					pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(), sdk.NewUint(types.MaxWbasis).String(), sdk.NewInt(0))
				newLpUnits, lpUnits, err := keeper.CalculatePoolUnits(pool.ExternalAsset.Symbol, temp, tempNative, tempExternal,
					withdrawNativeAssetAmount, withdrawExternalAssetAmount)
				if err != nil {
					hasError = true
					ctx.Logger().Error(fmt.Sprintf("failed to calculate pool units for | Pool : %s | LP %s ", pool.String(), lp.String()))
					break
				}
				lp.LiquidityProviderUnits = lpUnits
				if !lp.Validate() {
					hasError = true
					ctx.Logger().Error(fmt.Sprintf("Invalid | LP %s ", lp.String()))
					break
				}
				lps = append(lps, lp)
				tempExternal = tempExternal.Add(withdrawExternalAssetAmount)
				tempNative = tempNative.Add(withdrawNativeAssetAmount)
				temp = newLpUnits
			}
			pool.PoolUnits = temp
			if !app.ClpKeeper.ValidatePool(*pool) {
				hasError = true
				ctx.Logger().Error(fmt.Sprintf("Invalid | Pool %s ", pool.String()))
				break
			}
			poolList = append(poolList, *pool)
		}
		// If we have error dont set state
		if hasError {
			ctx.Logger().Error("Failed to execute upgrade plan for pool re-balance")
		}
		// If we have no errors , Set state .
		if !hasError {
			for i := range poolList {
				pool := poolList[i]
				_ = app.ClpKeeper.SetPool(ctx, &pool)
			}
			for i := range lps {
				l := lps[i]
				app.ClpKeeper.SetLiquidityProvider(ctx, &l)
			}
		}
	}
}
