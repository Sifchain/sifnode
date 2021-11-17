package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
)

const upgradeName = "0.9.14"

const distributionVersion = "0.9.14"
const distributionAddress = "sif1ct2s3t8u2kffjpaekhtngzv6yc4vm97xajqyl3"
const distributionAmount = "200000000000000000000000000" // 200m

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan types.Plan, vm m.VersionMap) (m.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + upgradeName)
		if plan.Name == distributionVersion {
			mintAmount, ok := sdk.NewIntFromString(distributionAmount)
			if !ok {
				panic("failed to get mint amount")
			}
			mintCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), mintAmount))
			err := app.BankKeeper.MintCoins(ctx, dispensationtypes.ModuleName, mintCoins)
			if err != nil {
				panic(err)
			}
			address, err := sdk.AccAddressFromBech32(distributionAddress)
			if err != nil {
				panic(err)
			}
			err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, dispensationtypes.ModuleName, address, mintCoins)
			if err != nil {
				panic(err)
			}
		}
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{},
		}

		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}