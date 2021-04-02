package ethbridge

import (
	"encoding/json"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/Sifchain/sifnode/x/ethbridge/client"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the ethbridge module.
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	panic("implement me")
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the ethbridge module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the ethbridge module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the ethbridge
// module.
func (b AppModuleBasic) DefaultGenesis(marshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the ethbridge module.
func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config sdkclient.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(c sdkclient.Context, serveMux *runtime.ServeMux) {
	panic("implement me")
}

// RegisterRESTRoutes registers the REST routes for the ethbridge module.
func (b AppModuleBasic) RegisterRESTRoutes(ctx sdkclient.Context, router *mux.Router) {
	client.RegisterRESTRoutes(ctx, router, StoreKey)
}

// GetTxCmd returns the root tx command for the ethbridge module.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return client.GetTxCmd(StoreKey)
}

// GetQueryCmd returns no root query command for the ethbridge module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return client.GetQueryCmd(StoreKey)
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the ethbridge module.
type AppModuleSimulation struct{}

// AppModule implements an application module for the ethbridge module.
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation

	OracleKeeper  types.OracleKeeper
	SupplyKeeper  types.SupplyKeeper
	AccountKeeper types.AccountKeeper
	BridgeKeeper  Keeper
	Codec         *codec.Marshaler
}

func (am AppModule) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	panic("implement me")
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	panic("implement me")
}

func (am AppModule) DefaultGenesis(marshaler codec.JSONMarshaler) json.RawMessage {
	panic("implement me")
}

func (am AppModule) ValidateGenesis(marshaler codec.JSONMarshaler, config sdkclient.TxEncodingConfig, message json.RawMessage) error {
	panic("implement me")
}

func (am AppModule) RegisterRESTRoutes(c sdkclient.Context, router *mux.Router) {
	panic("implement me")
}

func (am AppModule) RegisterGRPCGatewayRoutes(c sdkclient.Context, serveMux *runtime.ServeMux) {
	panic("implement me")
}

func (am AppModule) GetTxCmd() *cobra.Command {
	panic("implement me")
}

func (am AppModule) GetQueryCmd() *cobra.Command {
	panic("implement me")
}

func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	panic("implement me")
}

func (am AppModule) RegisterServices(configurator module.Configurator) {
	panic("implement me")
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	oracleKeeper types.OracleKeeper, supplyKeeper types.SupplyKeeper,
	accountKeeper types.AccountKeeper, bridgeKeeper Keeper,
	cdc *codec.Marshaler) AppModule {

	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},

		OracleKeeper:  oracleKeeper,
		SupplyKeeper:  supplyKeeper,
		AccountKeeper: accountKeeper,
		BridgeKeeper:  bridgeKeeper,
		Codec:         cdc,
	}
}

// Name returns the ethbridge module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the ethbridge module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// Route returns the message routing key for the ethbridge module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(RouterKey, am.NewHandler())
}

// NewHandler returns an sdk.Handler for the ethbridge module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.AccountKeeper, am.BridgeKeeper, am.Codec)
}

// QuerierRoute returns the ethbridge module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the ethbridge module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.OracleKeeper, am.Codec)
}

// InitGenesis performs genesis initialization for the ethbridge module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, marshaler codec.JSONMarshaler, _ json.RawMessage) []abci.ValidatorUpdate {
	bridgeAccount := authtypes.NewEmptyModuleAccount(ModuleName, authtypes.Burner, authtypes.Minter)
	am.SupplyKeeper.SetModuleAccount(ctx, bridgeAccount)
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the ethbridge
// module.
func (am AppModule) ExportGenesis(s sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}

// BeginBlock returns the begin blocker for the ethbridge module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the ethbridge module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
