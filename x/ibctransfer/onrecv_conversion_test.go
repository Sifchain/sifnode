package ibctransfer_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	whitelistmocks "github.com/Sifchain/sifnode/x/tokenregistry/types/mock"
)

func TestShouldConvertIncomingCoins(t *testing.T) {
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
		Denom:         "ueth",
		Decimals:      10,
		UnitDenom:     "ceth",
	}

	unitDenomEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ceth",
		Decimals:      18,
		UnitDenom:     "ceth",
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
	}

	wl := whitelistmocks.NewMockKeeper(ctrl)

	wl.EXPECT().GetDenom(ctx, "ueth").Return(ibcRegistryEntry)
	wl.EXPECT().GetDenom(ctx, "ceth").Return(unitDenomEntry)
	got := ibctransfer.ShouldConvertIncomingCoins(ctx, wl, returningTransferPacket, ibcDenom)
	require.Equal(t, got, true)

	wl.EXPECT().GetDenom(ctx, "cusdt").Return(nonIBCRegistryEntry)
	got = ibctransfer.ShouldConvertIncomingCoins(ctx, wl, returningTransferPacket, nonIBCDenom)
	require.Equal(t, got, false)
}

func TestGetConvForIncomingCoins(t *testing.T) {
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
		Denom:         "ueth",
		Decimals:      10,
		UnitDenom:     "ceth",
	}

	unitDenomEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ceth",
		Decimals:      18,
	}

	wl := whitelistmocks.NewMockKeeper(ctrl)

	wl.EXPECT().GetDenom(ctx, "ueth").Return(ibcRegistryEntry)
	wl.EXPECT().GetDenom(ctx, "ceth").Return(unitDenomEntry)
	gotIBCToken, gotConvToken := ibctransfer.GetConvForIncomingCoins(ctx, wl, returningTransferPacket, ibcDenom)
	intAmount, _ := sdk.NewIntFromString("100000000000000000000")
	require.Equal(t, gotIBCToken, sdk.NewCoin("ueth", sdk.NewInt(1000000000000)))
	require.Equal(t, gotConvToken, sdk.NewCoin("ceth", intAmount))
}
