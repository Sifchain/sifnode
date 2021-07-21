package ibctransfer

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer"
	sdktransferkeeper "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	whitelisttypes "github.com/Sifchain/sifnode/x/whitelist/types"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	_ porttypes.IBCModule   = AppModule{}
)

// AppModuleBasic defines the basic application module.
type AppModuleBasic struct {
	cosmosAppModule transfer.AppModule
}

func (am AppModuleBasic) Name() string {
	return am.cosmosAppModule.Name()
}

func (am AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.cosmosAppModule.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (am AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	am.cosmosAppModule.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the module.
func (am AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return am.cosmosAppModule.DefaultGenesis(cdc)
}

// ValidateGenesis performs genesis state validation for the module.
func (am AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	return am.cosmosAppModule.ValidateGenesis(cdc, config, bz)
}

// RegisterRESTRoutes registers the REST routes for the module.
func (am AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	am.cosmosAppModule.RegisterRESTRoutes(ctx, rtr)
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	am.cosmosAppModule.RegisterGRPCGatewayRoutes(clientCtx, mux)
}

// GetTxCmd returns the root tx command for the module.
func (am AppModuleBasic) GetTxCmd() *cobra.Command {
	// Append local TX cmd to this if required
	return am.cosmosAppModule.GetTxCmd()
}

// GetQueryCmd returns no root query command for the module.
func (am AppModuleBasic) GetQueryCmd() *cobra.Command {
	// Append local TX cmd to this if required
	return am.cosmosAppModule.GetQueryCmd()
}

//____________________________________________________________________________

// AppModule implements an application module for the dispensation module.
type AppModule struct {
	AppModuleBasic
	whitelistKeeper whitelisttypes.Keeper
	cdc             codec.BinaryMarshaler
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
	return OnRecvPacketWhiteListed(ctx, am.cosmosAppModule, am.whitelistKeeper, packet)
}

func (am AppModule) OnAcknowledgementPacket(ctx sdk.Context, packet types.Packet, acknowledgement []byte) (*sdk.Result, error) {
	return am.cosmosAppModule.OnAcknowledgementPacket(ctx, packet, acknowledgement)
}

func (am AppModule) OnTimeoutPacket(ctx sdk.Context, packet types.Packet) (*sdk.Result, error) {
	return am.cosmosAppModule.OnTimeoutPacket(ctx, packet)
}

func NewAppModule(keeper sdktransferkeeper.Keeper, whitelistKeeper whitelisttypes.Keeper, cdc codec.BinaryMarshaler) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cosmosAppModule: transfer.NewAppModule(keeper),
		},
		whitelistKeeper: whitelistKeeper,
		cdc:             cdc,
	}
}

// IBC does not support a legacy querier
func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return am.cosmosAppModule.LegacyQuerierHandler(amino)
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	am.cosmosAppModule.RegisterServices(cfg)
}

// Name returns the dispensation module's name.
func (am AppModule) Name() string {
	return am.cosmosAppModule.Name()
}

// RegisterInvariants registers the dispensation module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.cosmosAppModule.RegisterInvariants(ir)
}

// Route returns the message routing key for the dispensation module.
func (am AppModule) Route() sdk.Route {
	return am.cosmosAppModule.Route()
}

// QuerierRoute returns the dispensation module's querier route name.
func (am AppModule) QuerierRoute() string {
	return am.cosmosAppModule.QuerierRoute()
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
