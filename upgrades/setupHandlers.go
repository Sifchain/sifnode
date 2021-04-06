package upgrades

import (
	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/dispensation"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

func SetupHandlers(app *app.SifchainApp) {
	app.UpgradeKeeper.SetUpgradeHandler("changePoolFormula", GetPoolChangeFunc(app))
	app.UpgradeKeeper.SetUpgradeHandler("release-20210401000000", func(ctx sdk.Context, plan upgrade.Plan) {})
	app.UpgradeKeeper.SetUpgradeHandler("release-20210501000000-AddDispensation", GetAddDispensation(app))
	SetState(app)
}

func SetState(app *app.SifchainApp) {
	app.SetStoreLoader(bam.StoreLoaderWithUpgrade(&types.StoreUpgrades{
		Added: []string{dispensation.ModuleName},
	}))
}
