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

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ibctransfer"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper/testhelpers"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"
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
		IbcCounterpartyDenom: "microrowan",
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

func TestMultihopTransfer(t *testing.T) {
	sifapp.SetConfig(false)
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	addrs, _ := test.CreateTestAddrs(3)
	amount := uint64(123456789123456789)
	photonToken := tokenregistrytypes.RegistryEntry{
		Denom:     "ibc/4BFA1CE7B80A9A830F8E164495276CCD9E9B5424951749ED92F80B394E8C91C8",
		BaseDenom: "uphoton",
		Decimals:  6,
	}
	app.TokenRegistryKeeper.SetToken(ctx, &photonToken)
	first, err := sdk.AccAddressFromBech32(addrs[0].String())
	require.NoError(t, err)
	second, err := sdk.AccAddressFromBech32(addrs[1].String())
	require.NoError(t, err)
	third, err := sdk.AccAddressFromBech32(addrs[2].String())
	require.NoError(t, err)
	recvTokenPacket := transfertypes.FungibleTokenPacketData{
		Denom:    "uphoton",
		Amount:   amount,
		Sender:   first.String(),
		Receiver: second.String(),
	}
	recvPacket := channeltypes.Packet{
		SourceChannel:      "channel-27",
		SourcePort:         "transfer",
		DestinationChannel: "channel-11",
		DestinationPort:    "transfer",
		Data:               app.AppCodec().MustMarshalJSON(&recvTokenPacket),
	}
	err = app.TransferKeeper.OnRecvPacket(ctx, recvPacket, recvTokenPacket)
	require.NoError(t, err)
	require.Equal(t, "0", app.BankKeeper.GetBalance(ctx, first, "uphoton").Amount.String())
	require.Equal(t, "123456789123456789", app.BankKeeper.GetBalance(ctx, second, "ibc/4BFA1CE7B80A9A830F8E164495276CCD9E9B5424951749ED92F80B394E8C91C8").Amount.String())
	sdkSentDenom, _ := testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, sdk.NewCoin("ibc/4BFA1CE7B80A9A830F8E164495276CCD9E9B5424951749ED92F80B394E8C91C8", sdk.NewIntFromUint64(amount)), second, "transfer", "channel-12")
	recvTokenPacket = transfertypes.FungibleTokenPacketData{
		Denom:    sdkSentDenom,
		Amount:   amount,
		Sender:   second.String(),
		Receiver: third.String(),
	}
	recvPacket = channeltypes.Packet{
		SourceChannel:      "channel-12",
		SourcePort:         "transfer",
		DestinationChannel: "channel-66",
		DestinationPort:    "transfer",
		Data:               app.AppCodec().MustMarshalJSON(&recvTokenPacket),
	}
	err = app.TransferKeeper.OnRecvPacket(ctx, recvPacket, recvTokenPacket)
	require.NoError(t, err)
	require.Equal(t, "0", app.BankKeeper.GetBalance(ctx, second, "ibc/4BFA1CE7B80A9A830F8E164495276CCD9E9B5424951749ED92F80B394E8C91C8").Amount.String())
	require.Equal(t, "123456789123456789", app.BankKeeper.GetBalance(ctx, third, "ibc/ED52642E49540BE90488C9027BEA1C1AFA2BD296A548D99CA20EDAEF8F3BB5B9").Amount.String())
	sdkSentDenom, _ = testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, sdk.NewCoin("ibc/ED52642E49540BE90488C9027BEA1C1AFA2BD296A548D99CA20EDAEF8F3BB5B9", sdk.NewIntFromUint64(amount)), third, "transfer", "channel-66")
	recvTokenPacket = transfertypes.FungibleTokenPacketData{
		Denom:    sdkSentDenom,
		Amount:   amount,
		Sender:   third.String(),
		Receiver: second.String(),
	}
	recvPacket = channeltypes.Packet{
		SourceChannel:      "channel-66",
		SourcePort:         "transfer",
		DestinationChannel: "channel-12",
		DestinationPort:    "transfer",
		Data:               app.AppCodec().MustMarshalJSON(&recvTokenPacket),
	}
	err = app.TransferKeeper.OnRecvPacket(ctx, recvPacket, recvTokenPacket)
	require.NoError(t, err)
	require.Equal(t, "0", app.BankKeeper.GetBalance(ctx, third, "ibc/ED52642E49540BE90488C9027BEA1C1AFA2BD296A548D99CA20EDAEF8F3BB5B9").Amount.String())
	require.Equal(t, "123456789123456789", app.BankKeeper.GetBalance(ctx, second, "ibc/4BFA1CE7B80A9A830F8E164495276CCD9E9B5424951749ED92F80B394E8C91C8").Amount.String())
}
