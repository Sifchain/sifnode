package ibc

import (
	"encoding/json"
	"github.com/Sifchain/sifnode/x/ibc/client/cli"
	"github.com/Sifchain/sifnode/x/ibc/client/rest"

	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the dispensation module.
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return CosmosAppModuleBasic{}.Name()
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	CosmosAppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	CosmosAppModuleBasic{}.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the dispensation module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return CosmosAppModuleBasic{}.DefaultGenesis(cdc)
}

// ValidateGenesis performs genesis state validation for the dispensation module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	return CosmosAppModuleBasic{}.ValidateGenesis(cdc, config, bz)
}

// RegisterRESTRoutes registers the REST routes for the dispensation module.
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	//types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

// GetTxCmd returns the root tx command for the dispensation module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the dispensation module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

//____________________________________________________________________________

// AppModule implements an application module for the dispensation module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule CosmosAppModule
	keeper          Keeper
}

// Call this function from app.go
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: NewCosmosAppModule(&keeper),
		keeper:          keeper,
	}
}

// IBC does not support a legacy querier
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	CosmosAppModule{}.RegisterServices(cfg)
}

// Name returns the dispensation module's name.
func (AppModule) Name() string {
	return ibchost.ModuleName
}

// RegisterInvariants registers the dispensation module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the dispensation module.
func (am AppModule) Route() sdk.Route {
	return CosmosAppModule{}.Route()
}

// QuerierRoute returns the dispensation module's querier route name.
func (AppModule) QuerierRoute() string {
	return ibchost.QuerierRoute
}

// InitGenesis performs genesis initialization for the dispensation module. It returns
// no validator updates
func (am AppModule) InitGenesis(ctx sdk.Context, codec codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	return CosmosAppModule{}.InitGenesis(ctx, codec, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the dispensation
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, codec codec.JSONMarshaler) json.RawMessage {
	return CosmosAppModule{}.ExportGenesis(ctx, codec)
}

// BeginBlock returns the begin blocker for the dispensation module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	CosmosAppModule{}.BeginBlock(ctx, req)
}

// EndBlock returns the end blocker for the dispensation module. It returns no validator
// updates.
func (AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return CosmosAppModule{}.EndBlock(ctx, req)
}
