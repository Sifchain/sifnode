package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	transferTypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper interface {
	IsAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress) bool
	SetAdminAccount(ctx sdk.Context, adminAccount sdk.AccAddress)
	IsDenomWhitelisted(ctx sdk.Context, denom string) bool
	GetDenom(ctx sdk.Context, denom string) RegistryEntry
	SetToken(ctx sdk.Context, entry *RegistryEntry)
	RemoveToken(ctx sdk.Context, denom string)
	InitGenesis(ctx sdk.Context, state GenesisState) []abci.ValidatorUpdate
	ExportGenesis(ctx sdk.Context) *GenesisState
	GetDenomWhitelist(ctx sdk.Context) Registry
}

type SDKTransferKeeper interface {
	OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data transferTypes.FungibleTokenPacketData) error
	OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data transferTypes.FungibleTokenPacketData, ack channeltypes.Acknowledgement) error
	OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data transferTypes.FungibleTokenPacketData) error
}
