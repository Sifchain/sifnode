package ibctransfer_test

import (
	"context"
	"testing"

	"github.com/Sifchain/sifnode/x/ibctransfer/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func TestOnTimeoutPacketConvert(t *testing.T) {
	t.Skip()
	// TODO fix packet
	maxUInt64 := uint64(18446744073709551615)
	app, appCtx, admin := tokenregistrytest.CreateTestApp(false)

	packet := channeltypes.Packet{
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
	}

	rowanEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted:        true,
		Decimals:             18,
		Denom:                "rowan",
		BaseDenom:            "rowan",
		IbcCounterPartyDenom: "microrowan",
	}

	microRowanEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Decimals:      10,
		Denom:         "microrowan",
		BaseDenom:     "microrowan",
		UnitDenom:     "rowan",
	}

	decimal12Entry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted:        true,
		Decimals:             12,
		Denom:                "twelve",
		BaseDenom:            "twelve",
		IbcCounterPartyDenom: "microtwelve",
	}

	decimal10Entry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Decimals:      10,
		Denom:         "microtwelce",
		BaseDenom:     "microtwelve",
	}

	type args struct {
		goCtx         context.Context
		msg           *ibctransfertypes.MsgTransfer
		registryEntry tokenregistrytypes.RegistryEntry
		sendAsEntry   tokenregistrytypes.RegistryEntry
	}
	tests := []struct {
		name            string
		args            args
		tokenDeduction  sdk.Coin
		tokensConverted sdk.Coin
	}{
		{
			args:            args{sdk.WrapSDKContext(appCtx), &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(maxUInt64)), Sender: admin}, rowanEntry, microRowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(18446744073700000000)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(184467440737)),
		},
		{
			args:            args{sdk.WrapSDKContext(appCtx), &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(1000000)), Sender: admin}, rowanEntry, microRowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(0)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(0)),
		},
		{
			args:            args{sdk.WrapSDKContext(appCtx), &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("twelve", sdk.NewIntFromUint64(10000000000000)), Sender: admin}, decimal12Entry, decimal10Entry},
			tokenDeduction:  sdk.NewCoin("twelve", sdk.NewIntFromUint64(10000000000000)),
			tokensConverted: sdk.NewCoin("microtwelve", sdk.NewIntFromUint64(100000000000)),
		},
		{
			args:            args{sdk.WrapSDKContext(appCtx), &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("twelve", sdk.NewIntFromUint64(1)), Sender: admin}, decimal12Entry, decimal10Entry},
			tokenDeduction:  sdk.NewCoin("twelve", sdk.NewIntFromUint64(0)),
			tokensConverted: sdk.NewCoin("microtwelve", sdk.NewIntFromUint64(0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenDeduction, tokensConverted := keeper.ConvertCoinsForTransfer(tt.args.goCtx, tt.args.msg, tt.args.registryEntry, tt.args.sendAsEntry)
			require.Equal(t, tt.tokensConverted, tokensConverted)
			require.Equal(t, tt.tokenDeduction, tokenDeduction)

			initCoins := sdk.NewCoins(tokenDeduction)
			sender, err := sdk.AccAddressFromBech32(admin)
			require.NoError(t, err)

			err = app.BankKeeper.AddCoins(appCtx, sender, initCoins)
			require.NoError(t, err)

			err = keeper.PrepareToSendConvertedCoins(sdk.WrapSDKContext(appCtx), tt.args.msg, tokenDeduction, tokensConverted, app.BankKeeper)
			require.NoError(t, err)

			res, err := ibctransfer.OnTimeoutMaybeConvert(appCtx, app.TransferKeeper, app.TokenRegistryKeeper, app.BankKeeper, packet)
			require.NoError(t, err)
			require.NotNil(t, res)
			require.NotEmpty(t, res.Events)
		})
	}
}
