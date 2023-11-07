package clp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/clp/client/cli"
	"github.com/Sifchain/sifnode/x/clp/client/rest"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the clp module.
type AppModuleBasic struct{}

// Name returns the clp module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the clp module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) { //nolint
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the clp
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the clp module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the clp module.
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the clp module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic("Failed to register GRPC gateway routes.")
	}
}

// GetTxCmd returns the root tx command for the clp module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the clp module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

//____________________________________________________________________________

// AppModule implements an application module for the clp module.
type AppModule struct {
	AppModuleBasic

	keeper     keeper.Keeper
	bankKeeper types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, bankKeeper types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		bankKeeper:     bankKeeper,
	}
}

// Name returns the clp module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the clp module invariants.
func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
	keeper.RegisterInvariants(registry, am.keeper)
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// QuerierRoute returns the clp module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier { //nolint
	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	querier := keeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), querier)
	m := keeper.NewMigrator(am.keeper)
	err := cfg.RegisterMigration(types.ModuleName, 1, m.MigrateToVer2)
	if err != nil {
		panic(err)
	}
	err = cfg.RegisterMigration(types.ModuleName, 2, m.MigrateToVer3)
	if err != nil {
		panic(err)
	}
	err = cfg.RegisterMigration(types.ModuleName, 3, m.MigrateToVer4)
	if err != nil {
		panic(err)
	}
	err = cfg.RegisterMigration(types.ModuleName, 4, m.MigrateToVer5)
	if err != nil {
		panic(err)
	}
}

// InitGenesis performs genesis initialization for the clp module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the clp
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)

	return cdc.MustMarshalJSON(&gs)
}

func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock returns the end blocker for the clp module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

func (AppModule) ConsensusVersion() uint64 { return 5 }
