package ibctransfer_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestShouldConvertIncomingCoins(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	unitDenomEntry := tokenregistrytypes.RegistryEntry{
		Denom:     "ceth",
		Decimals:  18,
		UnitDenom: "ceth",
	}
	ibcRegistryEntry := tokenregistrytypes.RegistryEntry{
		Denom:     "ueth",
		Decimals:  10,
		UnitDenom: "ceth",
	}
	nonIBCRegistryEntry := tokenregistrytypes.RegistryEntry{
		Denom:    "cusdt",
		Decimals: 6,
	}
	app.TokenRegistryKeeper.SetToken(ctx, &unitDenomEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &ibcRegistryEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &nonIBCRegistryEntry)
	registry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)
	entry1 := app.TokenRegistryKeeper.GetDenom(registry, "ueth")
	require.NotNil(t, entry1)
	entry1c := app.TokenRegistryKeeper.GetDenom(registry, entry1.UnitDenom)
	require.NotNil(t, entry1c)
	diff := uint64(entry1c.Decimals - entry1.Decimals)
	convAmount := helpers.ConvertIncomingCoins(ctx, app.TokenRegistryKeeper, 1000000000000, diff)
	incomingDeduction := sdk.NewCoin("ueth", sdk.NewIntFromUint64(1000000000000))
	incomingAddition := sdk.NewCoin("ceth", convAmount)
	require.NotNil(t, incomingDeduction)
	require.NotNil(t, incomingAddition)
	require.Equal(t, incomingDeduction.Denom, "ueth")
	require.Equal(t, incomingDeduction.Amount.String(), "1000000000000")
	require.Equal(t, incomingAddition.Denom, "ceth")
	require.Equal(t, incomingAddition.Amount.String(), "100000000000000000000")
	entry2 := app.TokenRegistryKeeper.GetDenom(registry, "cusdt")
	require.NotNil(t, entry2)
	entry2c := app.TokenRegistryKeeper.GetDenom(registry, entry2.UnitDenom)
	require.Nil(t, entry2c)
}

func TestGetConvForIncomingCoins(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	ibcRegistryEntry := tokenregistrytypes.RegistryEntry{
		Denom:     "ueth",
		Decimals:  10,
		UnitDenom: "ceth",
	}
	unitDenomEntry := tokenregistrytypes.RegistryEntry{
		Denom:    "ceth",
		Decimals: 18,
	}
	app.TokenRegistryKeeper.SetToken(ctx, &unitDenomEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &ibcRegistryEntry)
	registry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)
	entry1 := app.TokenRegistryKeeper.GetDenom(registry, "ueth")
	require.NotNil(t, entry1)
	entry2 := app.TokenRegistryKeeper.GetDenom(registry, "ceth")
	require.NotNil(t, entry2)
	entry1c := app.TokenRegistryKeeper.GetDenom(registry, entry1.UnitDenom)
	require.NotNil(t, entry1c)
	diff := uint64(entry1c.Decimals - entry1.Decimals)
	convAmount := helpers.ConvertIncomingCoins(ctx, app.TokenRegistryKeeper, 1000000000000, diff)
	incomingDeduction := sdk.NewCoin("ueth", sdk.NewIntFromUint64(1000000000000))
	incomingAddition := sdk.NewCoin("ceth", convAmount)
	intAmount, _ := sdk.NewIntFromString("100000000000000000000")
	require.Equal(t, incomingDeduction, sdk.NewCoin("ueth", sdk.NewInt(1000000000000)))
	require.Equal(t, incomingAddition, sdk.NewCoin("ceth", intAmount))
}

func TestIsRecvPacketAllowed(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	transferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}
	returningDenom := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-0/rowan",
	}
	whitelistedDenom := transfertypes.FungibleTokenPacketData{
		Denom: "atom",
	}
	disallowedDenom := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-66/atom",
	}
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:                "atom",
		Decimals:             6,
		IbcCounterpartyDenom: "",
		Permissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:                "ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3",
		Decimals:             6,
		IbcCounterpartyDenom: "",
		Permissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:                "ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2",
		Decimals:             6,
		IbcCounterpartyDenom: "",
		Permissions:          []tokenregistrytypes.Permission{},
	})
	registry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)
	entry1 := app.TokenRegistryKeeper.GetDenom(registry, "ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3")
	require.NotNil(t, entry1)
	permitted1 := app.TokenRegistryKeeper.CheckDenomPermissions(entry1, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT})
	require.Equal(t, permitted1, true)
	got := ibctransfer.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, whitelistedDenom, entry1)
	require.Equal(t, got, true)
	entry2 := app.TokenRegistryKeeper.GetDenom(registry, "ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2")
	require.NotNil(t, entry2)
	permitted2 := app.TokenRegistryKeeper.CheckDenomPermissions(entry2, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT})
	require.Equal(t, permitted2, false)
	got = ibctransfer.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, disallowedDenom, entry2)
	require.Equal(t, got, false)
	entry3 := app.TokenRegistryKeeper.GetDenom(registry, "rowan")
	require.Nil(t, entry3)
	got = ibctransfer.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, returningDenom, entry3)
	require.Equal(t, got, true)
}

func TestIsRecvPacketReturning(t *testing.T) {
	packet := channeltypes.Packet{
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
	}
	returningData := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-0/atom",
	}
	nonReturningData := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-11/atom",
	}
	got := transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), returningData.Denom)
	require.Equal(t, got, true)
	got = transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), nonReturningData.Denom)
	require.Equal(t, got, false)
}
