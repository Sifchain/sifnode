package keeper

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper}
}

// func (m Migrator) MigrateToVer2(ctx sdk.Context) error {
// 	m.keeper.SetParams(ctx, types.DefaultGenesis().Params)

// 	return nil
// }