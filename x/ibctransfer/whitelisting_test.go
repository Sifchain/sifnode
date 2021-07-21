package ibctransfer

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	whitelistkeeper "github.com/Sifchain/sifnode/x/whitelist/keeper"
)

func TestIsRecvPacketAllowed(t *testing.T) {
	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "foochainid"}, false, nil)
	packet := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}
	data := transfertypes.FungibleTokenPacketData{
		Denom: "transfer/channel-0/atom",
	}

	enc := test.MakeTestEncodingConfig()
	wl := whitelistkeeper.NewKeeper(enc.Marshaler, sdk.NewKVStoreKey(""))
	got := isRecvPacketAllowed(ctx, wl, packet, data)
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
