package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const upgradeNameV096 = "0.9.6"

func SetupHandlers(app *SifchainApp) {
	SetupHandlersForMint(app)
}

func SetupHandlersForMint(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("0.9.6", func(ctx sdk.Context, plan types.Plan) {
		app.Logger().Info("Running upgrade handler for " + upgradeNameV096 + " with new store " + minttypes.StoreKey)
		// Install initial params and minter for mint module.
		mintGenesis := minttypes.DefaultGenesisState()
		// Replace default MintDenom with staking bond denom.
		mintGenesis.Params.MintDenom = app.StakingKeeper.GetParams(ctx).BondDenom
		mint.InitGenesis(ctx, app.MintKeeper, app.AccountKeeper, mintGenesis)

	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeNameV096 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{minttypes.StoreKey},
		}

		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
