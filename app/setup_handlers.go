package app

import (
	tokenRegistryTypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const upgradeName = "0.9.14"

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan types.Plan, vm m.VersionMap) (m.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + upgradeName)
		if plan.Name == "0.9.14" {
			registry := app.TokenRegistryKeeper.GetRegistry(ctx)
			for _, entry := range registry.Entries {
				if entry.Decimals > 9 && app.TokenRegistryKeeper.CheckEntryPermissions(entry, []tokenRegistryTypes.Permission{tokenRegistryTypes.Permission_CLP, tokenRegistryTypes.Permission_IBCEXPORT}) {
					entry.Permissions = append(entry.Permissions, tokenRegistryTypes.Permission_IBCIMPORT)
					entry.IbcCounterpartyDenom = ""
				}
			}
			app.TokenRegistryKeeper.SetRegistry(ctx, registry)
		}
		return vm, nil
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == upgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authz.ModuleName, feegrant.ModuleName, "vesting", "crisis"},
		}

		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
