package faucet

import (
	"encoding/json"

	"github.com/Sifchain/sifnode/x/faucet/client/cli"
	"github.com/Sifchain/sifnode/x/faucet/client/rest"
	"github.com/Sifchain/sifnode/x/faucet/keeper"
	"github.com/Sifchain/sifnode/x/faucet/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the faucet module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the faucet module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the faucet module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the faucet
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the faucet module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the faucet module.
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	if profile == TESTNET {
		rest.RegisterRoutes(ctx, rtr)
	}
}

// GetTxCmd returns the root tx command for the faucet module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	if profile == TESTNET {
		return cli.GetTxCmd(cdc)
	}
	return nil
}

// GetQueryCmd returns no root query command for the faucet module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	if profile == TESTNET {
		return cli.GetQueryCmd(types.StoreKey, cdc)
	}
	return nil
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the faucet module.
type AppModuleSimulation struct{}

// AppModule implements an application module for the faucet module.

//TODO Verify if we can remove supplykeeper and bankkeeper from this struct ,and access it only through keeper.Get{..}() methods
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation

	keeper       keeper.Keeper
	supplyKeeper types.SupplyKeeper
	bankKeeper   types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper, bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},
		keeper:              keeper,
		supplyKeeper:        supplyKeeper,
		bankKeeper:          bankKeeper,
	}
}

// Name returns the faucet module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the faucet module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// Route returns the message routing key for the faucet module.
func (AppModule) Route() string {
	return types.RouterKey
}

// NewHandler returns an sdk.Handler for the faucet module.
func (am AppModule) NewHandler() sdk.Handler {
	if profile == TESTNET {
		return NewHandler(am.keeper)
	}
	return nil
}

// QuerierRoute returns the faucet module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the faucet module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the faucet module. It returns
// no validator updates.
// TODO add functionality for Init genesis , Would need to verify if the faucet balance is automatically initialized by the supply module
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the faucet
// module.
// TODO add functionality for export genesis , We would need to keep track of how much an address withdrew to prevent spam
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the faucet module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock returns the end blocker for the faucet module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
