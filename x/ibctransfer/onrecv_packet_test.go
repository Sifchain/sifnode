package ibctransfer_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"

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
	got := helpers.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, whitelistedDenom, entry1)
	require.Equal(t, got, true)
	entry2 := app.TokenRegistryKeeper.GetDenom(registry, "ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2")
	require.NotNil(t, entry2)
	permitted2 := app.TokenRegistryKeeper.CheckDenomPermissions(entry2, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT})
	require.Equal(t, permitted2, false)
	got = helpers.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, disallowedDenom, entry2)
	require.Equal(t, got, false)
	entry3 := app.TokenRegistryKeeper.GetDenom(registry, "rowan")
	require.Nil(t, entry3)
	got = helpers.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, returningDenom, entry3)
	require.Equal(t, got, true)
}

func TestExecConvForIncomingCoins(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	addrs, _ := test.CreateTestAddrs(2)
	packet := channeltypes.Packet{
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
	}
	returningData := transfertypes.FungibleTokenPacketData{
		Denom:    "transfer/channel-0/ueth",
		Receiver: addrs[0].String(),
	}
	nonReturningData := transfertypes.FungibleTokenPacketData{
		Denom:    "transfer/channel-1/ueth",
		Receiver: addrs[0].String(),
	}
	ibcRegistryEntry := tokenregistrytypes.RegistryEntry{
		Denom:     "ueth",
		Decimals:  10,
		UnitDenom: "ceth",
	}
	ibcRegistryEntry2 := tokenregistrytypes.RegistryEntry{
		Denom:       "ibc/C1061B25E69D71E96BED65B5652168F41927316D07D6B417A3A9774F94A4CB7A",
		Decimals:    10,
		UnitDenom:   "ceth",
		Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT},
	}
	unitDenomEntry := tokenregistrytypes.RegistryEntry{
		Denom:    "ceth",
		Decimals: 18,
	}
	app.TokenRegistryKeeper.SetToken(ctx, &unitDenomEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &ibcRegistryEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &ibcRegistryEntry2)
	mintedDenom := helpers.GetMintedDenomFromPacket(packet, returningData)
	registry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)
	mintedDenomEntry := app.TokenRegistryKeeper.GetDenom(registry, mintedDenom)
	require.NotNil(t, mintedDenomEntry)
	allowed := helpers.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, packet, returningData, mintedDenomEntry)
	require.Equal(t, allowed, true)
	convertToDenomEntry := app.TokenRegistryKeeper.GetDenom(registry, mintedDenomEntry.UnitDenom)
	require.NotNil(t, convertToDenomEntry)
	err := helpers.ExecConvForIncomingCoins(ctx, app.BankKeeper, app.TokenRegistryKeeper, mintedDenomEntry, convertToDenomEntry, packet, returningData)
	require.NoError(t, err)
	mintedDenom = helpers.GetMintedDenomFromPacket(packet, nonReturningData)
	mintedDenomEntry = app.TokenRegistryKeeper.GetDenom(registry, mintedDenom)
	require.NotNil(t, mintedDenomEntry)
	allowed = helpers.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, packet, nonReturningData, mintedDenomEntry)
	require.Equal(t, allowed, true)
	convertToDenomEntry = app.TokenRegistryKeeper.GetDenom(registry, mintedDenomEntry.UnitDenom)
	require.NotNil(t, convertToDenomEntry)
	err = helpers.ExecConvForIncomingCoins(ctx, app.BankKeeper, app.TokenRegistryKeeper, mintedDenomEntry, convertToDenomEntry, packet, nonReturningData)
	require.NoError(t, err)
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

func TestGetMintedDenomFromPacket(t *testing.T) {
	packet := channeltypes.Packet{
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
	}
	returningData := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-0/atom",
	}
	returningData2 := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-0/rowan",
	}
	nonReturningData := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-11/atom",
	}
	nonReturningData2 := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-11/rowan",
	}
	nonReturningData3 := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-1/rowan",
	}
	nonReturningData4 := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-1/atom",
	}
	returningDenom := helpers.GetMintedDenomFromPacket(packet, returningData)
	returningDenom2 := helpers.GetMintedDenomFromPacket(packet, returningData2)
	nonReturningDenom := helpers.GetMintedDenomFromPacket(packet, nonReturningData)
	nonReturningDenom2 := helpers.GetMintedDenomFromPacket(packet, nonReturningData2)
	nonReturningDenom3 := helpers.GetMintedDenomFromPacket(packet, nonReturningData3)
	nonReturningDenom4 := helpers.GetMintedDenomFromPacket(packet, nonReturningData4)
	require.Equal(t, "atom", returningDenom)
	require.Equal(t, "rowan", returningDenom2)
	require.Equal(t, "ibc/611BB1D7CBB019DBA91690697697B3CB56335EBCFDD4573B9A11A34A20802940", nonReturningDenom)
	require.Equal(t, "ibc/666BF1A729F1F7FFF31563C704C303B0B73DC00A8B4C9072894239006227B96A", nonReturningDenom2)
	require.Equal(t, "ibc/6ABBE597A317EA31C9D1522D4DC4C5BF2EC8815A5B276713EE11EEDF2FA79012", nonReturningDenom3)
	require.Equal(t, "ibc/6D0449781D39534D032041B75F6C32DB251650083F7AC79C3975FFB7CDF7727F", nonReturningDenom4)
}
