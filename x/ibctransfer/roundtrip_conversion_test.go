package ibctransfer_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	whitelistmocks "github.com/Sifchain/sifnode/x/tokenregistry/types/mock"
)

func TestExportImportConversionEquality(t *testing.T) {
	ctrl := gomock.NewController(t)
	wl := whitelistmocks.NewMockKeeper(ctrl)
	ctx := sdk.NewContext(nil, tmproto.Header{ChainID: "foochainid"}, false, nil)

	maxUInt64 := uint64(18446744073709551615)
	microRowanEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Decimals:      10,
		Denom:         "microrowan",
		BaseDenom:     "microrowan",
		UnitDenom:     "rowan",
	}
	rowanEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted:        true,
		Decimals:             18,
		Denom:                "rowan",
		BaseDenom:            "rowan",
		IbcCounterPartyDenom: "microrowan",
	}

	wl.EXPECT().GetDenom(ctx, "microrowan").Return(microRowanEntry)

	msg := &transfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(maxUInt64))}
	outgoingDeduction, outgoingAddition := keeper.ConvertCoinsForTransfer(context.Background(), msg, rowanEntry, microRowanEntry)

	returningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}

	tokenPacket := transfertypes.FungibleTokenPacketData{
		// When sender chain is the source,
		// it simply sends the base denom without path prefix
		Denom:  "transfer/channel-0/microrowan",
		Amount: 184467440737,
	}

	wl.EXPECT().GetDenom(ctx, "microrowan").Return(microRowanEntry)
	wl.EXPECT().GetDenom(ctx, "rowan").Return(rowanEntry)

	incomingDeduction, incomingAddition := ibctransfer.GetConvForIncomingCoins(ctx, wl, returningTransferPacket, tokenPacket)
	require.Greater(t, incomingAddition.Amount.String(), incomingDeduction.Amount.String())
	require.Equal(t, outgoingDeduction, incomingAddition)
	require.Equal(t, outgoingAddition, incomingDeduction)
}
