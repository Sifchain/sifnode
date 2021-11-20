package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	tkrtypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/grpc"
)

type Migrator struct {
	keeper      types.TokenRegistryKeeper
	queryServer grpc.Server
}

func NewMigrator(keeper types.TokenRegistryKeeper, queryServer grpc.Server) Migrator {
	return Migrator{keeper: keeper, queryServer: queryServer}
}

func (m Migrator) MigrateToVer2(ctx sdk.Context) error {
	registry := m.keeper.GetRegistry(ctx)
	for _, entry := range registry.Entries {
		if entry.Decimals > 9 && m.keeper.CheckEntryPermissions(entry, []tkrtypes.Permission{tkrtypes.Permission_CLP, tkrtypes.Permission_IBCEXPORT}) {
			entry.Permissions = append(entry.Permissions, tkrtypes.Permission_IBCIMPORT)
			entry.IbcCounterpartyDenom = ""
		}
	}
	m.keeper.SetRegistry(ctx, registry)
	return nil
}
