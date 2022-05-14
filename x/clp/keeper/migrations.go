package keeper

import (
	"fmt"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	tkrKeeper "github.com/Sifchain/sifnode/x/tokenregistry/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) MigrateToVer3(ctx sdk.Context) error {
	pools := m.keeper.GetPools(ctx)
	migrationMap := tkrKeeper.GetDenomMigrationMap()
	for _, pool := range pools {
		if migrationExtAsset, found := migrationMap[pool.ExternalAsset.Symbol]; found {
			peggy2Denom := types.GetDenom(migrationExtAsset.NetworkDescriptor, migrationExtAsset.TokenContractAddress)
			peggy2Pool := *pool
			peggy2Pool.ExternalAsset = &clptypes.Asset{Symbol: peggy2Denom}
			err := m.keeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
			if err != nil {
				return err
			}
			err = m.keeper.SetPool(ctx, &peggy2Pool)
			if err != nil {
				return err
			}
		} else {
			return clptypes.ErrMigrationFailed.Wrap(fmt.Sprintf("Unable to migrate pool : %s ", pool.String()))
		}
	}
	// TODO : Add migration for Rewards after Branch is rebased with develop
	// TODO : Modify consensus version of CLP after Branch is rebased with develop
	return nil
}
