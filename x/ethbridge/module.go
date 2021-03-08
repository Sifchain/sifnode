package ethbridge

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/ethbridge/client"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the ethbridge module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the ethbridge module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the ethbridge module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the ethbridge
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the ethbridge module.
func (AppModuleBasic) ValidateGenesis(_ json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the ethbridge module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	client.RegisterRESTRoutes(ctx, rtr, StoreKey)
}

// GetTxCmd returns the root tx command for the ethbridge module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return client.GetTxCmd(StoreKey, cdc)
}

// GetQueryCmd returns no root query command for the ethbridge module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return client.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the ethbridge module.
type AppModuleSimulation struct{}

// AppModule implements an application module for the ethbridge module.
type AppModule struct {
	AppModuleBasic

	OracleKeeper  types.OracleKeeper
	BankKeeper    types.BankKeeper
	AccountKeeper types.AccountKeeper
	BridgeKeeper  Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	oracleKeeper types.OracleKeeper, bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper, bridgeKeeper Keeper,
	cdc *codec.Codec) AppModule {

	return AppModule{
		AppModuleBasic: AppModuleBasic{},

		OracleKeeper:  oracleKeeper,
		BankKeeper:    bankKeeper,
		AccountKeeper: accountKeeper,
		BridgeKeeper:  bridgeKeeper,
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
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the ethbridge module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.AccountKeeper, am.BridgeKeeper)
}

// QuerierRoute returns the ethbridge module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the ethbridge module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.OracleKeeper)
}

// InitGenesis performs genesis initialization for the ethbridge module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, _ json.RawMessage) []abci.ValidatorUpdate {
	bridgeAccount := authtypes.NewEmptyModuleAccount(ModuleName, authtypes.Burner, authtypes.Minter)
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes for the ethbridge
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}

// BeginBlock returns the begin blocker for the ethbridge module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the ethbridge module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
