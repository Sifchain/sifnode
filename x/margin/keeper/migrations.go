package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper}
}

func (m Migrator) MigrateToVer2(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)

	params.RowanCollateralEnabled = false

	m.keeper.SetParams(ctx, &params)

	return nil
}
