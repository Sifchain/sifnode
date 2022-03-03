package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"
)

const releaseVersion = "0.11.0"

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, plan types.Plan, vm m.VersionMap) (m.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)
		app.IBCKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())
		/*
				The exact APR depends on the total Bonded Rowan , and can thus fluctuate a little .

				- Inflation Percentage Required = APR * BondRatio
					Where
			        BondRatio = ( Total Bonded Rowan/ Total Supply Rowan)

				- Calculations for APR 300 % , assuming the max APR to be 350 and min APR to be 250
				    - 300% → 41.78
				    - 350 % → 48.74
				    - 250 % → 34.81
		*/
		minter := minttypes.Minter{
			Inflation:        sdk.MustNewDecFromStr("0.417800000000000000"),
			AnnualProvisions: sdk.ZeroDec(),
		}
		app.MintKeeper.SetMinter(ctx, minter)
		params := app.MintKeeper.GetParams(ctx)
		params.InflationMax = sdk.MustNewDecFromStr("0.487400000000000000")
		params.InflationMin = sdk.MustNewDecFromStr("0.348100000000000000")
		app.MintKeeper.SetParams(ctx, params)
		return app.mm.RunMigrations(ctx, app.configurator, vm)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == releaseVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
