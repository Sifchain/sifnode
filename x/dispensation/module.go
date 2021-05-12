package dispensation

import (
	"encoding/json"
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation/client/cli"
	"github.com/Sifchain/sifnode/x/dispensation/client/rest"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/dispensation/keeper"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the dispensation module.
type AppModuleBasic struct{}

// Name returns the clp module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the dispensation module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the dispensation
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the dispensation module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the dispensation module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the dispensation module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the dispensation module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the dispensation module.
type AppModule struct {
	AppModuleBasic

	keeper       keeper.Keeper
	bankKeeper   types.BankKeeper
	supplyKeeper types.SupplyKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		bankKeeper:     bankKeeper,
		supplyKeeper:   supplyKeeper,
	}
}

// Name returns the dispensation module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the dispensation module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the dispensation module.
func (AppModule) Route() string {
	return RouterKey
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
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the dispensation module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the dispensation
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the dispensation module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the dispensation module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
