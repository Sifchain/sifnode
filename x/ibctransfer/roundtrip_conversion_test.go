package ibctransfer_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	test "github.com/Sifchain/sifnode/x/ethbridge/test"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper/testhelpers"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestExportImportConversionEquality(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	maxUInt64 := uint64(18446744073709551615)
	rowanEntry := tokenregistrytypes.RegistryEntry{
		Decimals:             18,
		Denom:                "rowan",
		BaseDenom:            "rowan",
		IbcCounterpartyDenom: "microrowan",
	}
	microRowanEntry := tokenregistrytypes.RegistryEntry{
		Decimals:  10,
		Denom:     "microrowan",
		BaseDenom: "microrowan",
		UnitDenom: "rowan",
	}
	app.TokenRegistryKeeper.SetToken(ctx, &rowanEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &microRowanEntry)
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	rEntry := app.TokenRegistryKeeper.GetDenom(registry, "rowan")
	require.NotNil(t, rEntry)
	mrEntry := app.TokenRegistryKeeper.GetDenom(registry, "microrowan")
	require.NotNil(t, mrEntry)
	msg := &transfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(maxUInt64))}
	outgoingDeduction, outgoingAddition := helpers.ConvertCoinsForTransfer(msg, rEntry, mrEntry)
	mrEntryUnit := app.TokenRegistryKeeper.GetDenom(registry, mrEntry.UnitDenom)
	require.NotNil(t, mrEntryUnit)
	diff := uint64(mrEntryUnit.Decimals - mrEntry.Decimals)
	convAmount := helpers.ConvertIncomingCoins(184467440737, diff)
	incomingDeduction := sdk.NewCoin("microrowan", sdk.NewIntFromUint64(184467440737))
	incomingAddition := sdk.NewCoin("rowan", convAmount)
	require.Greater(t, incomingAddition.Amount.String(), incomingDeduction.Amount.String())
	require.Equal(t, outgoingDeduction, incomingAddition)
	require.Equal(t, outgoingAddition, incomingDeduction)
}

func TestMultihopTransfer(t *testing.T) {
	sifapp.SetConfig(false)
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	addrs, _ := test.CreateTestAddrs(3)
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
		Amount:   uint64(123456789123456789),
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
	sdkSentDenom, _ := testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, sdk.NewCoin("ibc/4BFA1CE7B80A9A830F8E164495276CCD9E9B5424951749ED92F80B394E8C91C8", sdk.NewIntFromUint64(123456789123456789)), second, "transfer", "channel-12")
	recvTokenPacket = transfertypes.FungibleTokenPacketData{
		Denom:    sdkSentDenom,
		Amount:   uint64(123456789123456789),
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
	sdkSentDenom, _ = testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, sdk.NewCoin("ibc/ED52642E49540BE90488C9027BEA1C1AFA2BD296A548D99CA20EDAEF8F3BB5B9", sdk.NewIntFromUint64(123456789123456789)), third, "transfer", "channel-66")
	recvTokenPacket = transfertypes.FungibleTokenPacketData{
		Denom:    sdkSentDenom,
		Amount:   uint64(123456789123456789),
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
