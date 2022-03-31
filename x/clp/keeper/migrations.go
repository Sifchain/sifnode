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
	pools := m.keeper.GetPools(ctx)
	// compute swap prices for each pool
	for _, pool := range pools {
		spe := sdk.ZeroDec()
		spn := sdk.ZeroDec()
		pool.SwapPriceExternal = &spe
		pool.SwapPriceNative = &spn
		m.keeper.SetPool(ctx, pool)
	}
	return nil
}
