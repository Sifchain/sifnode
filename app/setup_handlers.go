package app

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const releaseVersion = "0.13.2"

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, plan types.Plan, vm m.VersionMap) (m.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)
		amt, ok := sdk.NewIntFromString("150000000000000000000000000")
		if !ok {
			panic("error converting mint amount")
		}
		mintCoins := sdk.NewCoins(sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, amt))
		err := app.BankKeeper.MintCoins(ctx, dispensationtypes.ModuleName, mintCoins)
		if err != nil {
			panic(err)
		}
		address, err := sdk.AccAddressFromBech32("sif1ct2s3t8u2kffjpaekhtngzv6yc4vm97xajqyl3")
		if err != nil {
			panic(err)
		}
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, dispensationtypes.ModuleName, address, mintCoins)
		if err != nil {
			panic(err)
		}

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
