package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) MigrateToVer2(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)

	params.ForceCloseFundPercentage = sdk.NewDecWithPrec(1, 1)
	params.InsuranceFundAddress = ""
	params.SqModifier = sdk.MustNewDecFromStr("10000000000000000000000000")
	params.SafetyFactor = sdk.MustNewDecFromStr("1.05")

	m.keeper.SetParams(ctx, &params)

	return nil
}
