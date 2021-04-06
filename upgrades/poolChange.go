package upgrades

import (
	"encoding/json"
	"fmt"
	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/clp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"io/ioutil"
)

func GetPoolChangeFunc(app *app.SifchainApp) func(ctx sdk.Context, plan upgrade.Plan) {
	return func(ctx sdk.Context, plan upgrade.Plan) {
		ctx.Logger().Info("Starting to execute upgrade plan for pool re-balance")

		ExportAppState("changePoolFormula", app, ctx)

		allPools := app.ClpKeeper.GetPools(ctx)
		lps := clp.LiquidityProviders{}
		poolList := clp.Pools{}
		hasError := false
		for _, pool := range allPools {
			lpList := app.ClpKeeper.GetLiquidityProvidersForAsset(ctx, pool.ExternalAsset)
			temp := sdk.ZeroUint()
			tempExternal := sdk.ZeroUint()
			tempNative := sdk.ZeroUint()
			for _, lp := range lpList {
				withdrawNativeAssetAmount, withdrawExternalAssetAmount, _, _ := clp.CalculateWithdrawal(pool.PoolUnits, pool.NativeAssetBalance.String(),
					pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(), sdk.NewUint(clp.MaxWbasis).String(), sdk.NewInt(0))
				newLpUnits, lpUnits, err := clp.CalculatePoolUnits(pool.ExternalAsset.Symbol, temp, tempNative, tempExternal,
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
			if !app.ClpKeeper.ValidatePool(pool) {
				hasError = true
				ctx.Logger().Error(fmt.Sprintf("Invalid | Pool %s ", pool.String()))
				break
			}
			poolList = append(poolList, pool)
		}
		// If we have error dont set state
		if hasError {
			ctx.Logger().Error("Failed to execute upgrade plan for pool re-balance")
		}
		// If we have no errors , Set state .
		if !hasError {
			for _, pool := range poolList {
				_ = app.ClpKeeper.SetPool(ctx, pool)
			}
			for _, l := range lps {
				app.ClpKeeper.SetLiquidityProvider(ctx, l)
			}
		}
	}
}

func ExportAppState(name string, app *app.SifchainApp, ctx sdk.Context) {
	appState, vallist, err := app.ExportAppStateAndValidators(true, []string{})
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("failed to export app state: %s", err))
		return
	}
	appStateJSON, err := app.Cdc.MarshalJSON(appState)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("failed to marshal application genesis state: %s", err.Error()))
		return
	}
	valList, err := json.MarshalIndent(vallist, "", " ")
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("failed to marshal application genesis state: %s", err.Error()))
	}

	err = ioutil.WriteFile(fmt.Sprintf("%v-state.json", name), appStateJSON, 0600)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("failed to write state to file: %s", err.Error()))
	}
	err = ioutil.WriteFile(fmt.Sprintf("%v-validator.json", name), valList, 0600)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("failed to write Validator List to file: %s", err.Error()))
	}
}
