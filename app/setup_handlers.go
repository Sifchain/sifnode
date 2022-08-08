package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	m "github.com/cosmos/cosmos-sdk/types/module"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const releaseVersion = "0.15.0"

var minCommissionRate = sdk.NewDecWithPrec(5, 2) //5%

func SetupHandlers(app *SifchainApp) {

	app.UpgradeKeeper.SetUpgradeHandler(releaseVersion, func(ctx sdk.Context, plan types.Plan, vm m.VersionMap) (m.VersionMap, error) {
		app.Logger().Info("Running upgrade handler for " + releaseVersion)

		validators := app.StakingKeeper.GetAllValidators(ctx)

		for _, v := range validators {
			if v.Commission.Rate.LT(minCommissionRate) {
				comm, err := MustUpdateValidatorCommission(
					ctx, v, minCommissionRate)
				if err != nil {
					panic(err)
				}
				v.Commission = comm

				// call the before-modification hook since we're about to update the commission
				app.StakingKeeper.BeforeValidatorModified(ctx, v.GetOperator())

				app.StakingKeeper.SetValidator(ctx, v)
			}
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

// MustUpdateValidatorCommission updates a validator's commission rate,
// ignoring the max change rate.
func MustUpdateValidatorCommission(ctx sdk.Context,
	validator stakingtypes.Validator, newRate sdk.Dec) (stakingtypes.Commission, error) {
	commission := validator.Commission
	blockTime := ctx.BlockHeader().Time

	commission.Rate = newRate
	commission.UpdateTime = blockTime

	if validator.Commission.MaxRate.LT(newRate) {
		commission.MaxRate = newRate
	}

	return commission, nil
}
