package ibctransfer_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

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
		IsWhitelisted:        true,
		Decimals:             6,
		IbcCounterpartyDenom: "",
		Permissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:                "ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3",
		IsWhitelisted:        true,
		Decimals:             6,
		IbcCounterpartyDenom: "",
		Permissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:                "ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2",
		IsWhitelisted:        true,
		Decimals:             6,
		IbcCounterpartyDenom: "",
		Permissions:          []tokenregistrytypes.Permission{},
	})
	registry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)
	entry1 := app.TokenRegistryKeeper.GetDenom(registry, "ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3")
	require.NotNil(t, entry1)
	permitted1 := app.TokenRegistryKeeper.CheckDenomPermissions(entry1, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT})
	require.Equal(t, permitted1, true)
	got := ibctransfer.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, whitelistedDenom)
	require.Equal(t, got, true)
	entry2 := app.TokenRegistryKeeper.GetDenom(registry, "ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2")
	require.NotNil(t, entry2)
	permitted2 := app.TokenRegistryKeeper.CheckDenomPermissions(entry2, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT})
	require.Equal(t, permitted2, false)
	got = ibctransfer.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, disallowedDenom)
	require.Equal(t, got, false)
	entry3 := app.TokenRegistryKeeper.GetDenom(registry, "rowan")
	require.Nil(t, entry3)
	got = ibctransfer.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, transferPacket, returningDenom)
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
