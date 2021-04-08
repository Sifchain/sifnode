package app

import (
	"github.com/Sifchain/sifnode/x/dispensation"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

func SetupHandlers(app *SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("changePoolFormula", GetPoolChangeFunc(app))
	app.UpgradeKeeper.SetUpgradeHandler("release-20210401000000", func(ctx sdk.Context, plan upgrade.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("release-20210407042000", GetAddDispensation(app))
	SetState(app)
}

func SetState(app *SifchainApp) {
	app.SetStoreLoader(bam.StoreLoaderWithUpgrade(&types.StoreUpgrades{
		Added: []string{dispensation.ModuleName},
	}))
}
