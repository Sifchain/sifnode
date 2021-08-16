package ibctransfer

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	whitelistmocks "github.com/Sifchain/sifnode/x/tokenregistry/types/mock"
)

func TestIsRecvPacketAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "foochainid"}, false, nil)

	returningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}

	nonReturningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}

	whitelistedDenom := transfertypes.FungibleTokenPacketData{
		// When sender chain is the source,
		// it simply sends the base denom without path prefix
		Denom: "atom",
	}

	disallowedDenom := transfertypes.FungibleTokenPacketData{
		// If atom has a prefix when coming in,
		// it has an extra hop in the path received in ibc transfer OnRecvPacket().
		Denom: "transfer/channel-66/atom",
	}

	wl := whitelistmocks.NewMockKeeper(ctrl)

	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3").
		Return(true)
	got := isRecvPacketAllowed(ctx, wl, nonReturningTransferPacket, whitelistedDenom)
	require.Equal(t, got, true)

	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2").
		Return(false)
	got = isRecvPacketAllowed(ctx, wl, nonReturningTransferPacket, disallowedDenom)
	require.Equal(t, got, false)

	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/A916425D9C00464330F8B333711C4A51AA8CF0141392E7E250371EC6D4285BF2").
		Return(true)
	got = isRecvPacketAllowed(ctx, wl, returningTransferPacket, disallowedDenom)
	require.Equal(t, got, true)

	wl.EXPECT().
		IsDenomWhitelisted(ctx,
			"ibc/44F0BAC50DDD0C83DAC9CEFCCC770C12F700C0D1F024ED27B8A3EE9DD949BAD3").
		Return(true)
	got = isRecvPacketAllowed(ctx, wl, returningTransferPacket, whitelistedDenom)
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

	got := IsRecvPacketReturning(packet, returningData)
	require.Equal(t, got, true)

	got = IsRecvPacketReturning(packet, nonReturningData)
	require.Equal(t, got, false)
}

func TestCheckRecvConvert(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "foochainid"}, false, nil)

	returningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}

	ibcDenom := transfertypes.FungibleTokenPacketData{
		// When sender chain is the source,
		// it simply sends the base denom without path prefix
		Denom: "transfer/channel-0/ueth",
	}

	ibcRegistryEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ceth",
		Decimals:      18,
		IbcDenom:      "ueth",
		IbcDecimals:   10,
	}

	nonIBCDenom := transfertypes.FungibleTokenPacketData{
		// When sender chain is the source,
		// it simply sends the base denom without path prefix
		Denom: "transfer/channel-0/cusdt",
	}

	nonIBCRegistryEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "cusdt",
		Decimals:      6,
		IbcDenom:      "",
		IbcDecimals:   0,
	}

	wl := whitelistmocks.NewMockKeeper(ctrl)

	wl.EXPECT().GetRegistryEntry(ctx, "ueth").Return(ibcRegistryEntry)
	got := shouldConvertDecimals(ctx, wl, returningTransferPacket, ibcDenom)
	require.Equal(t, got, true)

	wl.EXPECT().GetRegistryEntry(ctx, "cusdt").Return(nonIBCRegistryEntry)
	got = shouldConvertDecimals(ctx, wl, returningTransferPacket, nonIBCDenom)
	require.Equal(t, got, false)
}

func TestConvertRecvDenom(t *testing.T) {
	ctrl := gomock.NewController(t)

	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "foochainid"}, false, nil)

	returningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}

	ibcDenom := transfertypes.FungibleTokenPacketData{
		// When sender chain is the source,
		// it simply sends the base denom without path prefix
		Denom:  "transfer/channel-0/ueth",
		Amount: 1000000000000,
	}

	ibcRegistryEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ceth",
		Decimals:      18,
		IbcDenom:      "ueth",
		IbcDecimals:   10,
	}

	// nonIBCDenom := transfertypes.FungibleTokenPacketData{
	// 	// When sender chain is the source,
	// 	// it simply sends the base denom without path prefix
	// 	Denom:  "transfer/channel-0/cusdt",
	// 	Amount: 100000000,
	// }

	// nonIBCRegistryEntry := tokenregistrytypes.RegistryEntry{
	// 	IsWhitelisted: true,
	// 	Denom:         "cusdt",
	// 	Decimals:      6,
	// 	IbcDenom:      "",
	// 	IbcDecimals:   0,
	// }

	wl := whitelistmocks.NewMockKeeper(ctrl)

	wl.EXPECT().GetRegistryEntry(ctx, "ueth").Return(ibcRegistryEntry)
	gotIBCToken, gotConvToken := convertDecimals(ctx, wl, returningTransferPacket, ibcDenom)
	intAmount, _ := sdk.NewIntFromString("100000000000000000000")
	require.Equal(t, gotIBCToken, sdk.NewCoin("ueth", sdk.NewInt(1000000000000)))
	require.Equal(t, gotConvToken, sdk.NewCoin("ceth", intAmount))

	// wl.EXPECT().GetRegistryEntry(ctx, "cusdt").Return(nonIBCRegistryEntry)
	// got = checkRecvConvert(ctx, wl, returningTransferPacket, nonIBCDenom)
	// require.Equal(t, got, false)
}
