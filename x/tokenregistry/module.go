package tokenregistry

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sifchain/sifnode/x/tokenregistry/client/cli"
	"github.com/Sifchain/sifnode/x/tokenregistry/client/rest"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/tokenregistry/handler"
	"github.com/Sifchain/sifnode/x/tokenregistry/keeper"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

var (
	ModuleName                       = types.ModuleName
	_          module.AppModule      = AppModule{}
	_          module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module.
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes.
func (b AppModuleBasic) DefaultGenesis(marshaler codec.JSONMarshaler) json.RawMessage {
	return marshaler.MustMarshalJSON(&types.GenesisState{})
}

// ValidateGenesis performs genesis state validation.
func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, _ sdkclient.TxEncodingConfig, message json.RawMessage) error {
	var data types.GenesisState

	if message != nil {
		err := marshaler.UnmarshalJSON(message, &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
		}
	}

	return nil
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx sdkclient.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}
}

// RegisterRESTRoutes registers the REST routes.
func (b AppModuleBasic) RegisterRESTRoutes(ctx sdkclient.Context, router *mux.Router) {
	rest.RegisterRESTRoutes(ctx, router)
}

// GetTxCmd returns the root tx command.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions.
type AppModuleSimulation struct{}

// AppModule implements an application module.
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation

	// BankKeeper types.BankKeeper
	Keeper types.Keeper
	Codec  *codec.Marshaler
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.Keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.Keeper))
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	// bankKeeper types.BankKeeper,
	keeper types.Keeper,
	cdc *codec.Marshaler) AppModule {

	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},

		// BankKeeper:    bankKeeper,
		Keeper: keeper,
		Codec:  cdc,
	}
}

// Name returns the module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
}

// Route returns the message routing key for the module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, am.NewHandler())
}

// NewHandler returns an sdk.Handler for the module.
func (am AppModule) NewHandler() sdk.Handler {
	return handler.NewHandler(am.Keeper)
}

// QuerierRoute returns the module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// Deprecated: LegacyQuerierHandler use RegisterServices
func (am AppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier {
	return keeper.NewLegacyQuerier(am.Keeper)
}

// InitGenesis performs genesis initialization. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, marshaler codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	marshaler.MustUnmarshalJSON(data, &genesisState)

	return am.Keeper.InitGenesis(ctx, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	return marshaler.MustMarshalJSON(am.Keeper.ExportGenesis(ctx))
}

// BeginBlock returns the begin blocker.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker. It returns no validator
// updates.
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
