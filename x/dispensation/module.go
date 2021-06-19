package dispensation

import (
	"encoding/json"
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/client/rest"

	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/Sifchain/sifnode/x/dispensation/client/cli"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/dispensation/keeper"
	"github.com/Sifchain/sifnode/x/dispensation/types"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the dispensation module.
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the dispensation module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// RegisterCodec registers the dispensation module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	RegisterCodec(cdc)
}

// ValidateGenesis performs genesis state validation for the dispensation module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}
	return ValidateGenesis(data)
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
	return cli.GetQueryCmd(StoreKey)
}

//____________________________________________________________________________

// AppModule implements an application module for the dispensation module.
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
}

func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return keeper.NewLegacyQuerier(am.keeper)
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	querier := keeper.NewQuerier(am.keeper)
	types.RegisterQueryServer(cfg.QueryServer(), querier)
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		bankKeeper:     bankKeeper,
		accountKeeper:  accountKeeper,
	}
}

// Name returns the dispensation module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the dispensation module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the dispensation module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// NewHandler returns an sdk.Handler for the dispensation module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the dispensation module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the dispensation module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewLegacyQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the dispensation module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, codec codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	codec.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the dispensation
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, codec codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return codec.MustMarshalJSON(&gs)
}

// BeginBlock returns the begin blocker for the dispensation module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the dispensation module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
