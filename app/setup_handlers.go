package app

import (
	"fmt"

	kpr "github.com/Sifchain/sifnode/x/clp/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	marginkeeper "github.com/Sifchain/sifnode/x/margin/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const releaseVersion = "1.0-beta.12"

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, plan types.Plan, vm m.VersionMap) (m.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)
		// This is part of the scheduled process , directly doing state transitions here instead to migrating consensus version
		// The following functions fix the state of sifnode caused by unexpected swap behaviour triggered by margin logic.
		fixAtomPool(ctx, app.ClpKeeper)
		closeMtp(ctx, app.MarginKeeper)
		return app.mm.RunMigrations(ctx, app.configurator, vm)
	})

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
	if upgradeInfo.Name == releaseVersion && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			// Added: []string{},
		}
		// Use upgrade store loader for the initial loading of all stores when app starts,
		// it checks if version == upgradeHeight and applies store upgrades before loading the stores,
		// so that new stores start with the correct version (the current height of chain),
		// instead the default which is the latest version that store last committed i.e 0 for new stores.
		app.SetStoreLoader(types.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func fixAtomPool(ctx sdk.Context, k kpr.Keeper) {
	atomIbcHash := "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"
	// Get Rowan Balance from CLP Module
	clpModuleTotalNativeBalance := k.GetBankKeeper().GetBalance(ctx, clptypes.GetCLPModuleAddress(), clptypes.GetSettlementAsset().Symbol)
	// Get Atom Balance from CLP Module
	clpModulebalanceAtom := k.GetBankKeeper().GetBalance(ctx, clptypes.GetCLPModuleAddress(), atomIbcHash)

	// Get Uint amount from coin
	clpModuleTotalNativeBalanceUint := sdk.NewUintFromString(clpModuleTotalNativeBalance.Amount.String())
	clpModulebalanceAtomUint := sdk.NewUintFromString(clpModulebalanceAtom.Amount.String())

	// Get Atom Pool
	atomPool, err := k.GetPool(ctx, atomIbcHash)
	if err != nil {
		panic(fmt.Sprintf("Error getting pool %s | %s", atomIbcHash, err.Error()))
	}

	// Calculate total native balance of all pools
	pools := k.GetPools(ctx)
	poolTotalNativeBalance := sdk.ZeroUint()
	for _, pool := range pools {
		if pool.ExternalAsset.Symbol != atomIbcHash {
			poolTotalNativeBalance = poolTotalNativeBalance.Add(pool.NativeAssetBalance)
		}
	}
	// Set Atom pool back
	atomPool.ExternalAssetBalance = clpModulebalanceAtomUint
	atomPool.NativeAssetBalance = clpModuleTotalNativeBalanceUint.Sub(poolTotalNativeBalance)
	atomPool.ExternalLiabilities = sdk.ZeroUint()
	atomPool.NativeLiabilities = sdk.ZeroUint()
	atomPool.ExternalCustody = sdk.ZeroUint()
	atomPool.NativeCustody = sdk.ZeroUint()
	err = k.SetPool(ctx, &atomPool)
	if err != nil {
		panic(fmt.Sprintf("Error setting pool %s | %s", atomPool.String(), err.Error()))
	}
}

func closeMtp(ctx sdk.Context, k marginkeeper.Keeper) {
	mtps := k.GetAllMTPS(ctx)
	for _, mtp := range mtps {
		err := k.DestroyMTP(ctx, mtp.Address, mtp.Id)
		if err != nil {
			panic(fmt.Sprintf("Error closing MTP | %s", err.Error()))
		}
	}
}
