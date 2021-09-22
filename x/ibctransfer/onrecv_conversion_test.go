package ibctransfer_test

import (
	"testing"

	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func TestShouldConvertIncomingCoins(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	returningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}
	ibcDenom := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-0/ueth",
	}
	nonIBCDenom := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-0/cusdt",
	}
	unitDenomEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ceth",
		Decimals:      18,
		UnitDenom:     "ceth",
	}
	ibcRegistryEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ueth",
		Decimals:      10,
		UnitDenom:     "ceth",
	}
	nonIBCRegistryEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "cusdt",
		Decimals:      6,
	}
	app.TokenRegistryKeeper.SetToken(ctx, &unitDenomEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &ibcRegistryEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &nonIBCRegistryEntry)
	registry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)
	entry1 := app.TokenRegistryKeeper.GetDenom(registry, "ceth")
	require.NotNil(t, entry1)
	entry2 := app.TokenRegistryKeeper.GetDenom(registry, "ueth")
	require.NotNil(t, entry2)
	entry3 := app.TokenRegistryKeeper.GetDenom(registry, "cusdt")
	require.NotNil(t, entry3)
	incomingDeduction, incomingAddition := ibctransfer.GetConvForIncomingCoins(ctx, app.TokenRegistryKeeper, returningTransferPacket, ibcDenom)
	require.NotNil(t, incomingDeduction)
	require.NotNil(t, incomingAddition)
	incomingDeduction, incomingAddition = ibctransfer.GetConvForIncomingCoins(ctx, app.TokenRegistryKeeper, returningTransferPacket, nonIBCDenom)
	require.Nil(t, incomingDeduction)
	require.Nil(t, incomingAddition)
}

func TestGetConvForIncomingCoins(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	returningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}
	ibcDenom := transfertypes.FungibleTokenPacketData{
		Denom:  "transfer/channel-0/ueth",
		Amount: 1000000000000,
	}
	ibcRegistryEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ueth",
		Decimals:      10,
		UnitDenom:     "ceth",
	}
	unitDenomEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ceth",
		Decimals:      18,
	}
	app.TokenRegistryKeeper.SetToken(ctx, &unitDenomEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &ibcRegistryEntry)
	registry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)
	entry1 := app.TokenRegistryKeeper.GetDenom(registry, "ueth")
	require.NotNil(t, entry1)
	entry2 := app.TokenRegistryKeeper.GetDenom(registry, "ceth")
	require.NotNil(t, entry2)
	gotIBCToken, gotConvToken := ibctransfer.GetConvForIncomingCoins(ctx, app.TokenRegistryKeeper, returningTransferPacket, ibcDenom)
	intAmount, _ := sdk.NewIntFromString("100000000000000000000")
	require.Equal(t, *gotIBCToken, sdk.NewCoin("ueth", sdk.NewInt(1000000000000)))
	require.Equal(t, *gotConvToken, sdk.NewCoin("ceth", intAmount))
}
