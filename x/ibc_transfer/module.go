package ibc_transfer

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	ibc_transfer_types "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"
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
	_ porttypes.IBCModule   = AppModule{}
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
	CosmosAppModuleBasic{}.RegisterRESTRoutes(ctx, rtr)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	CosmosAppModuleBasic{}.RegisterGRPCGatewayRoutes(clientCtx, mux)
}

// GetTxCmd returns the root tx command for the dispensation module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	// Append local TX cmd to this if required
	return CosmosAppModuleBasic{}.GetTxCmd()
}

// GetQueryCmd returns no root query command for the dispensation module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	// Append local TX cmd to this if required
	return CosmosAppModuleBasic{}.GetQueryCmd()
}

//____________________________________________________________________________

// AppModule implements an application module for the dispensation module.
type AppModule struct {
	AppModuleBasic
	cosmosAppModule CosmosAppModule
	keeper          Keeper
}

func (am AppModule) OnChanOpenInit(ctx sdk.Context, order types.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty types.Counterparty, version string) error {
	return am.cosmosAppModule.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (am AppModule) OnChanOpenTry(ctx sdk.Context, order types.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty types.Counterparty, version, counterpartyVersion string) error {
	return am.cosmosAppModule.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version, counterpartyVersion)
}

func (am AppModule) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyVersion string) error {
	return am.cosmosAppModule.OnChanOpenAck(ctx, portID, channelID, counterpartyVersion)
}

func (am AppModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return am.cosmosAppModule.OnChanOpenConfirm(ctx, portID, channelID)
}

func (am AppModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return am.cosmosAppModule.OnChanCloseInit(ctx, portID, channelID)
}

func (am AppModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return am.cosmosAppModule.OnChanOpenConfirm(ctx, portID, channelID)
}

func (am AppModule) OnRecvPacket(ctx sdk.Context, packet types.Packet) (*sdk.Result, []byte, error) {
	return OnRecvPacketWhiteListed(am.keeper, ctx, packet)
}

func (am AppModule) OnAcknowledgementPacket(ctx sdk.Context, packet types.Packet, acknowledgement []byte) (*sdk.Result, error) {
	return am.cosmosAppModule.OnAcknowledgementPacket(ctx, packet, acknowledgement)
}

func (am AppModule) OnTimeoutPacket(ctx sdk.Context, packet types.Packet) (*sdk.Result, error) {
	return am.cosmosAppModule.OnTimeoutPacket(ctx, packet)
}

// Call this function from app.go
func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic:  AppModuleBasic{},
		cosmosAppModule: NewCosmosAppModule(keeper),
		keeper:          keeper,
	}
}

// IBC does not support a legacy querier
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	am.cosmosAppModule.RegisterServices(cfg)
}

// Name returns the dispensation module's name.
func (AppModule) Name() string {
	return ibc_transfer_types.ModuleName
}

// RegisterInvariants registers the dispensation module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the dispensation module.
func (am AppModule) Route() sdk.Route {
	return am.cosmosAppModule.Route()
}

// QuerierRoute returns the dispensation module's querier route name.
func (AppModule) QuerierRoute() string {
	return ibc_transfer_types.QuerierRoute
}

// InitGenesis performs genesis initialization for the dispensation module. It returns
// no validator updates
func (am AppModule) InitGenesis(ctx sdk.Context, codec codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	return am.cosmosAppModule.InitGenesis(ctx, codec, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the dispensation
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, codec codec.JSONMarshaler) json.RawMessage {
	return am.cosmosAppModule.ExportGenesis(ctx, codec)
}

// BeginBlock returns the begin blocker for the dispensation module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	am.cosmosAppModule.BeginBlock(ctx, req)
}

// EndBlock returns the end blocker for the dispensation module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.cosmosAppModule.EndBlock(ctx, req)
}

// OnRecvPacketWhiteListed overrides the default implementation to add whitelisting functionality
