package ibctransfer_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"

	sifapp "github.com/Sifchain/sifnode/app"
	test2 "github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ibctransfer"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper"
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func TestOnAcknowledgementMaybeConvert_Source(t *testing.T) {
	sifapp.SetConfig(false)
	addrs, _ := test2.CreateTestAddrs(2)

	rowanToken := tokenregistrytypes.RegistryEntry{
		Denom:                "rowan",
		IbcCounterPartyDenom: "xrowan",
		Decimals:             18,
	}

	xrowanToken := tokenregistrytypes.RegistryEntry{
		Denom:     "xrowan",
		UnitDenom: "rowan",
		Decimals:  10,
	}

	// successAck := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	errorAck := channeltypes.NewErrorAcknowledgement("failed packet transfer")

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
		goCtx           context.Context
		msg             *types.MsgTransfer
		destChannel     string
		transferToken   tokenregistrytypes.RegistryEntry
		packetToken     tokenregistrytypes.RegistryEntry
		acknowledgement channeltypes.Acknowledgement
	}

	tests := []struct {
		name   string
		args   args
		err    error
		events sdk.Events
	}{
		{
			name: "Ack err sender is source, causes refund - success",
			args: args{
				context.Background(),
				msgSourceTransfer,
				"channel-1",
				rowanToken,
				xrowanToken,
				errorAck,
			},
			err:    nil,
			events: []sdk.Event{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app, ctx, _ := test.CreateTestApp(false)

			// Setup the send conversion before testing ACK.
			tokenDeduction, tokensConverted := keeper.ConvertCoinsForTransfer(sdk.WrapSDKContext(ctx), tt.args.msg, tt.args.transferToken, tt.args.packetToken)

			initCoins := sdk.NewCoins(tokenDeduction)
			sender, err := sdk.AccAddressFromBech32(tt.args.msg.Sender)
			require.NoError(t, err)

			err = app.BankKeeper.AddCoins(ctx, sender, initCoins)
			require.NoError(t, err)

			err = keeper.PrepareToSendConvertedCoins(sdk.WrapSDKContext(ctx), tt.args.msg, tokenDeduction, tokensConverted, app.BankKeeper)
			require.NoError(t, err)

			// Simulate send with SDK stub.
			sdkSentDenom, err := sendStub(ctx, app, tokensConverted, sender, tt.args.msg.SourcePort, tt.args.msg.SourceChannel)

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

			_, err = ibctransfer.OnAcknowledgementMaybeConvert(ctx, app.TransferKeeper, app.TokenRegistryKeeper, app.BankKeeper, packet, app.AppCodec().MustMarshalJSON(&tt.args.acknowledgement))
			require.ErrorIs(t, err, tt.err)
			// Assert events have recorded what happened.
		})
	}
}

func TestOnAcknowledgementMaybeConvert_Sink(t *testing.T) {
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

	// successAck := channeltypes.NewResultAcknowledgement([]byte{byte(1)})
	errorAck := channeltypes.NewErrorAcknowledgement("failed packet transfer")

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
		goCtx           context.Context
		msg             *types.MsgTransfer
		transferToken   tokenregistrytypes.RegistryEntry
		acknowledgement channeltypes.Acknowledgement
	}

	tests := []struct {
		name   string
		args   args
		err    error
		events sdk.Events
	}{
		{
			name: "Ack err sender is sink, causes refund - success",
			args: args{
				context.Background(),
				msgSinkTransfer,
				atomToken,
				errorAck,
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
			sdkSentDenom, err := sendStub(ctx, app, tt.args.msg.Token, sender, tt.args.msg.SourcePort, tt.args.msg.SourceChannel)

			// Test Ack.
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

			_, err = ibctransfer.OnAcknowledgementMaybeConvert(ctx, app.TransferKeeper, app.TokenRegistryKeeper, app.BankKeeper, ackPacket, app.AppCodec().MustMarshalJSON(&tt.args.acknowledgement))
			require.ErrorIs(t, err, tt.err)
			// Assert events have recorded what happened.
		})
	}
}

func sendStub(ctx sdk.Context, app *sifapp.SifchainApp, token sdk.Coin, sender sdk.AccAddress, sourcePort, sourceChannel string) (string, error) {
	// deconstruct the token denomination into the denomination trace info
	// to determine if the sender is the source chain
	fullDenomPath := token.Denom
	var err error
	if strings.HasPrefix(token.Denom, "ibc/") {
		fullDenomPath, err = app.TransferKeeper.DenomPathFromHash(ctx, token.Denom)
		if err != nil {
			return "", err
		}
	}

	if types.SenderChainIsSource(sourcePort, sourceChannel, fullDenomPath) {
		// create the escrow address for the tokens
		escrowAddress := types.GetEscrowAddress(sourcePort, sourceChannel)

		// escrow source tokens. It fails if balance insufficient.
		if err := app.BankKeeper.SendCoins(
			ctx, sender, escrowAddress, sdk.NewCoins(token),
		); err != nil {
			return "", err
		}

	} else {
		// transfer the coins to the module account and burn them
		if err := app.BankKeeper.SendCoinsFromAccountToModule(
			ctx, sender, types.ModuleName, sdk.NewCoins(token),
		); err != nil {
			return "", err
		}

		if err := app.BankKeeper.BurnCoins(
			ctx, types.ModuleName, sdk.NewCoins(token),
		); err != nil {
			// NOTE: should not happen as the module account was
			// retrieved on the step above and it has enough balace
			// to burn.
			panic(fmt.Sprintf("cannot burn coins after a successful send to a module account: %v", err))
		}
	}

	return fullDenomPath, nil
}
