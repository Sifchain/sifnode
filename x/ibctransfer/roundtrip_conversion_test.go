package ibctransfer_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper/testhelpers"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestMultihopTransfer(t *testing.T) {
	sifapp.SetConfig(false)
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	addrs, _ := test.CreateTestAddrs(3)
	amount := "123456789123456789"
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
	intAmount, ok := sdk.NewIntFromString(amount)
	require.True(t, ok)
	require.Equal(t, "123456789123456789", app.BankKeeper.GetBalance(ctx, second, "ibc/4BFA1CE7B80A9A830F8E164495276CCD9E9B5424951749ED92F80B394E8C91C8").Amount.String())
	sdkSentDenom, _ := testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, sdk.NewCoin("ibc/4BFA1CE7B80A9A830F8E164495276CCD9E9B5424951749ED92F80B394E8C91C8", intAmount), second, "transfer", "channel-12")
	require.Equal(t, "transfer/channel-11/uphoton", sdkSentDenom)
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
	sdkSentDenom, _ = testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, sdk.NewCoin("ibc/ED52642E49540BE90488C9027BEA1C1AFA2BD296A548D99CA20EDAEF8F3BB5B9", intAmount), third, "transfer", "channel-66")
	require.Equal(t, "transfer/channel-66/transfer/channel-11/uphoton", sdkSentDenom)
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
