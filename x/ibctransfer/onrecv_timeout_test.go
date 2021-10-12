package ibctransfer_test

import (
	"context"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	test2 "github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ibctransfer"
	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper/testhelpers"
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestOnTimeoutPacketConvert_Source(t *testing.T) {
	sifapp.SetConfig(false)
	addrs, _ := test2.CreateTestAddrs(2)
	rowanToken := tokenregistrytypes.RegistryEntry{
		Denom:                "rowan",
		IbcCounterpartyDenom: "xrowan",
		Decimals:             18,
	}
	xrowanToken := tokenregistrytypes.RegistryEntry{
		Denom:     "xrowan",
		UnitDenom: "rowan",
		Decimals:  10,
	}
	msgSourceTransfer := types.NewMsgTransfer(
		"transfer",
		"channel-0",
		sdk.NewCoin("rowan", sdk.NewInt(123456789123456789)),
		addrs[0],
		addrs[1].String(),
		clienttypes.NewHeight(0, 0),
		0,
	)
	type args struct {
		goCtx         context.Context
		msg           *types.MsgTransfer
		destChannel   string
		transferToken tokenregistrytypes.RegistryEntry
		packetToken   tokenregistrytypes.RegistryEntry
	}
	tests := []struct {
		name   string
		args   args
		err    error
		events sdk.Events
	}{
		{
			name: "Timeout sender is source, causes refund - success",
			args: args{
				context.Background(),
				msgSourceTransfer,
				"channel-1",
				rowanToken,
				xrowanToken,
			},
			err:    nil,
			events: []sdk.Event{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app, ctx, _ := test.CreateTestApp(false)
			app.TokenRegistryKeeper.SetToken(ctx, &tt.args.transferToken)
			app.TokenRegistryKeeper.SetToken(ctx, &tt.args.packetToken)
			// Setup the send conversion before testing timeout.
			tokenDeduction, tokensConverted := helpers.ConvertCoinsForTransfer(tt.args.msg, &tt.args.transferToken, &tt.args.packetToken)
			initCoins := sdk.NewCoins(tt.args.msg.Token)
			sender, err := sdk.AccAddressFromBech32(tt.args.msg.Sender)
			require.NoError(t, err)
			err = app.BankKeeper.AddCoins(ctx, sender, initCoins)
			require.NoError(t, err)
			err = helpers.PrepareToSendConvertedCoins(sdk.WrapSDKContext(ctx), tt.args.msg, tokenDeduction, tokensConverted, app.BankKeeper)
			require.NoError(t, err)
			require.Equal(t, tokensConverted.String(), app.BankKeeper.GetBalance(ctx, sender, tokensConverted.Denom).String())
			require.Equal(t, tt.args.msg.Token.Sub(tokenDeduction).String(), app.BankKeeper.GetBalance(ctx, sender, tt.args.msg.Token.Denom).String())
			// Simulate send with SDK stub.
			sdkSentDenom, err := testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, tokensConverted, sender, tt.args.msg.SourcePort, tt.args.msg.SourceChannel)
			require.NoError(t, err)
			require.Equal(t, tt.args.msg.Token.Sub(tokenDeduction).String(), app.BankKeeper.GetBalance(ctx, sender, tt.args.msg.Token.Denom).String())
			require.Equal(t, "0"+tokensConverted.Denom, app.BankKeeper.GetBalance(ctx, sender, tokensConverted.Denom).String())
			// Test Ack.
			packet := channeltypes.Packet{
				SourceChannel:      "channel-0",
				SourcePort:         "transfer",
				DestinationChannel: "channel-1",
				DestinationPort:    "transfer",
				Data: app.AppCodec().MustMarshalJSON(&types.FungibleTokenPacketData{
					Denom:    sdkSentDenom,
					Amount:   tokensConverted.Amount.Uint64(),
					Sender:   tt.args.msg.Sender,
					Receiver: tt.args.msg.Receiver,
				}),
			}
			_, err = ibctransfer.OnTimeoutMaybeConvert(ctx, app.TransferKeeper, app.TokenRegistryKeeper, app.BankKeeper, packet)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.args.msg.Token.String(), app.BankKeeper.GetBalance(ctx, sender, tt.args.msg.Token.Denom).String())
		})
	}
}

func TestOnTimeoutPacketConvert_Sink(t *testing.T) {
	sifapp.SetConfig(false)
	addrs, _ := test2.CreateTestAddrs(2)
	denomTrace := types.DenomTrace{
		// A token coming from source will have this chain's source channel prepended when this chain generates hash.
		Path:      "transfer/channel-0",
		BaseDenom: "uatom",
	}
	atomToken := tokenregistrytypes.RegistryEntry{
		Denom:     denomTrace.IBCDenom(),
		BaseDenom: "uatom",
		Decimals:  6,
	}
	msgSinkTransfer := types.NewMsgTransfer(
		"transfer",
		"channel-0", // Sent from this chain back to source
		sdk.NewCoin(atomToken.Denom, sdk.NewIntFromUint64(123456789123456789)),
		addrs[0],
		addrs[1].String(),
		clienttypes.NewHeight(0, 0),
		0,
	)
	type args struct {
		goCtx         context.Context
		msg           *types.MsgTransfer
		transferToken tokenregistrytypes.RegistryEntry
	}
	tests := []struct {
		name   string
		args   args
		err    error
		events sdk.Events
	}{
		{
			name: "Timeout sender is sink, causes refund - success",
			args: args{
				context.Background(),
				msgSinkTransfer,
				atomToken,
			},
			err:    nil,
			events: []sdk.Event{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app, ctx, _ := test.CreateTestApp(false)
			recvTokenPacket := types.FungibleTokenPacketData{
				Denom:  atomToken.BaseDenom,
				Amount: tt.args.msg.Token.Amount.Uint64(),
				Sender: tt.args.msg.Receiver,
				// Fund the addr that will do a send later.
				Receiver: tt.args.msg.Sender,
			}
			recvPacket := channeltypes.Packet{
				SourceChannel:      "channel-1",
				SourcePort:         "transfer",
				DestinationChannel: "channel-0",
				DestinationPort:    "transfer",
				Data:               app.AppCodec().MustMarshalJSON(&recvTokenPacket),
			}
			// Simulate OnRecv so that IBC hash is stored in transfer keeper and can be,
			// converted to denom trace during processing ack.
			err := app.TransferKeeper.OnRecvPacket(ctx, recvPacket, recvTokenPacket)
			require.NoError(t, err)
			sender, err := sdk.AccAddressFromBech32(tt.args.msg.Sender)
			require.NoError(t, err)
			// Simulate send from this chain, with SDK stub.
			sdkSentDenom, err := testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, tt.args.msg.Token, sender, tt.args.msg.SourcePort, tt.args.msg.SourceChannel)
			require.NoError(t, err)
			ackPacket := channeltypes.Packet{
				SourceChannel:      "channel-0",
				SourcePort:         "transfer",
				DestinationChannel: "channel-1",
				DestinationPort:    "transfer",
				Data: app.AppCodec().MustMarshalJSON(&types.FungibleTokenPacketData{
					Denom:    sdkSentDenom,
					Amount:   tt.args.msg.Token.Amount.Uint64(),
					Sender:   tt.args.msg.Sender,
					Receiver: tt.args.msg.Receiver,
				}),
			}
			_, err = ibctransfer.OnTimeoutMaybeConvert(ctx, app.TransferKeeper, app.TokenRegistryKeeper, app.BankKeeper, ackPacket)
			require.ErrorIs(t, err, tt.err)
			require.Equal(t, tt.args.msg.Token.String(), app.BankKeeper.GetBalance(ctx, sender, tt.args.msg.Token.Denom).String())
		})
	}
}

func TestOnTimeoutMaybeConvert(t *testing.T) {
	app, ctx, _ := test.CreateTestApp(false)
	addrs, _ := test2.CreateTestAddrs(2)
	rowanToken := tokenregistrytypes.RegistryEntry{
		Denom:                "rowan",
		IbcCounterpartyDenom: "xrowan",
		Decimals:             18,
	}
	xrowanToken := tokenregistrytypes.RegistryEntry{
		Denom:     "xrowan",
		UnitDenom: "rowan",
		Decimals:  10,
	}
	app.TokenRegistryKeeper.SetToken(ctx, &rowanToken)
	app.TokenRegistryKeeper.SetToken(ctx, &xrowanToken)
	rowan := sdk.NewCoin(rowanToken.Denom, sdk.NewInt(123456789123456789))
	msgSourceTransfer := types.NewMsgTransfer(
		"transfer",
		"channel-0",
		rowan,
		addrs[0],
		addrs[1].String(),
		clienttypes.NewHeight(0, 0),
		0,
	)
	initCoins := sdk.NewCoins(rowan)
	sender, err := sdk.AccAddressFromBech32(addrs[0].String())
	require.NoError(t, err)
	err = app.BankKeeper.AddCoins(ctx, sender, initCoins)
	require.NoError(t, err)
	tokenDeduction, tokensConverted := helpers.ConvertCoinsForTransfer(msgSourceTransfer, &rowanToken, &xrowanToken)
	err = helpers.PrepareToSendConvertedCoins(sdk.WrapSDKContext(ctx), msgSourceTransfer, tokenDeduction, tokensConverted, app.BankKeeper)
	require.NoError(t, err)
	sentDenom, err := testhelpers.SendStub(ctx, app.TransferKeeper, app.BankKeeper, tokensConverted, sender, "transfer", "channel-0")
	require.NoError(t, err)
	require.Equal(t, "0", app.BankKeeper.GetBalance(ctx, sender, sentDenom).Amount.String())
	timeoutPacket := channeltypes.Packet{
		SourceChannel:      "channel-0",
		SourcePort:         "transfer",
		DestinationChannel: "channel-1",
		DestinationPort:    "transfer",
		Data: app.AppCodec().MustMarshalJSON(&types.FungibleTokenPacketData{
			Denom:    sentDenom,
			Amount:   uint64(1234567891),
			Sender:   addrs[0].String(),
			Receiver: addrs[1].String(),
		}),
	}
	_, err = ibctransfer.OnTimeoutMaybeConvert(ctx, app.TransferKeeper, app.TokenRegistryKeeper, app.BankKeeper, timeoutPacket)
	require.NoError(t, err)
	require.Equal(t, rowan.String(), app.BankKeeper.GetBalance(ctx, sender, rowan.Denom).String())
}
