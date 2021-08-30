package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper"
	scibctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	scibctransfermocks "github.com/Sifchain/sifnode/x/ibctransfer/types/mocks"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func TestMsgServer_Transfer(t *testing.T) {
	/* Test that when a conversion is needed the right amounts are converted before sending to underlying SDK Transfer.
	 */
	ctrl := gomock.NewController(t)
	bankKeeper := scibctransfermocks.NewMockBankKeeper(ctrl)
	msgSrv := scibctransfermocks.NewMockMsgServer(ctrl)

	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	addrs, _ := test.CreateTestAddrs(2)

	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:                "rowan",
		IsWhitelisted:        true,
		Decimals:             18,
		IbcCounterPartyDenom: "xrowan",
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:         "xrowan",
		IsWhitelisted: true,
		Decimals:      10,
		UnitDenom:     "rowan",
	})

	rowanAmount, ok := sdk.NewIntFromString("1234567891123456789")
	require.True(t, ok)
	rowanAmountEscrowed, ok := sdk.NewIntFromString("1234567891100000000")
	require.True(t, ok)
	xrowanAmount, ok := sdk.NewIntFromString("12345678911")
	require.True(t, ok)

	rowanSmallest, ok := sdk.NewIntFromString("183456789")
	require.True(t, ok)

	rowanTooSmall, ok := sdk.NewIntFromString("12345678")
	require.True(t, ok)

	tt := []struct {
		name                 string
		err                  error
		bankKeeper           scibctransfertypes.BankKeeper
		msgSrv               scibctransfertypes.MsgServer
		msg                  *sdktransfertypes.MsgTransfer
		setupMsgServerCalls  func()
		setupBankKeeperCalls func()
	}{
		{
			name:       "transfer rowan with conversion",
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("rowan", rowanAmount),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupMsgServerCalls: func() {
				msgSrv.EXPECT().Transfer(gomock.Any(), &sdktransfertypes.MsgTransfer{
					SourcePort:       "transfer",
					SourceChannel:    "channel-0",
					Token:            sdk.NewCoin("xrowan", xrowanAmount),
					Sender:           addrs[0].String(),
					Receiver:         addrs[1].String(),
					TimeoutHeight:    clienttypes.NewHeight(0, 0),
					TimeoutTimestamp: 0,
				})
			},
			// TODO: Mint into the scibctransfer module account instead of SDK module account.
			setupBankKeeperCalls: func() {
				bankKeeper.EXPECT().SendCoins(gomock.Any(), addrs[0], sdktransfertypes.GetEscrowAddress("transfer", "channel-0"), sdk.NewCoins(sdk.NewCoin("rowan", rowanAmountEscrowed))).Return(nil)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), sdktransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("xrowan", xrowanAmount))).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), sdktransfertypes.ModuleName, addrs[0], sdk.NewCoins(sdk.NewCoin("xrowan", xrowanAmount))).Return(nil)
			},
		},
		{
			name:       "transfer smallest rowan without rounding",
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("rowan", rowanSmallest),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupMsgServerCalls: func() {
				msgSrv.EXPECT().Transfer(gomock.Any(), &sdktransfertypes.MsgTransfer{
					SourcePort:       "transfer",
					SourceChannel:    "channel-0",
					Token:            sdk.NewCoin("xrowan", sdk.NewInt(1)),
					Sender:           addrs[0].String(),
					Receiver:         addrs[1].String(),
					TimeoutHeight:    clienttypes.NewHeight(0, 0),
					TimeoutTimestamp: 0,
				})
			},
			setupBankKeeperCalls: func() {
				bankKeeper.EXPECT().SendCoins(gomock.Any(), addrs[0], sdktransfertypes.GetEscrowAddress("transfer", "channel-0"), sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(100000000)))).Return(nil)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), sdktransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("xrowan", sdk.NewInt(1)))).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), sdktransfertypes.ModuleName, addrs[0], sdk.NewCoins(sdk.NewCoin("xrowan", sdk.NewInt(1)))).Return(nil)
			},
		},
		{
			name:       "transfer amount too small for conversion",
			err:        scibctransfertypes.ErrAmountTooLowToConvert,
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("rowan", rowanTooSmall),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupBankKeeperCalls: func() {},
			setupMsgServerCalls:  func() {},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMsgServerCalls()
			tc.setupBankKeeperCalls()

			srv := keeper.NewMsgServerImpl(tc.msgSrv, tc.bankKeeper, app.TokenRegistryKeeper)
			_, err := srv.Transfer(sdk.WrapSDKContext(ctx), tc.msg)
			require.ErrorIs(t, tc.err, err)
		})
	}
}

func TestConvertCoins(t *testing.T) {
	ctx := context.Background()

	maxUInt64 := uint64(18446744073709551615)

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
		msg           *sdktransfertypes.MsgTransfer
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
			args:            args{ctx, &sdktransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(maxUInt64))}, rowanEntry, microRowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(18446744073700000000)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(184467440737)),
		},
		{
			args:            args{ctx, &sdktransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(1000000))}, rowanEntry, microRowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(0)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(0)),
		},
		{
			args:            args{ctx, &sdktransfertypes.MsgTransfer{Token: sdk.NewCoin("twelve", sdk.NewIntFromUint64(10000000000000))}, decimal12Entry, decimal10Entry},
			tokenDeduction:  sdk.NewCoin("twelve", sdk.NewIntFromUint64(10000000000000)),
			tokensConverted: sdk.NewCoin("microtwelve", sdk.NewIntFromUint64(100000000000)),
		},
		{
			args:            args{ctx, &sdktransfertypes.MsgTransfer{Token: sdk.NewCoin("twelve", sdk.NewIntFromUint64(1))}, decimal12Entry, decimal10Entry},
			tokenDeduction:  sdk.NewCoin("twelve", sdk.NewIntFromUint64(0)),
			tokensConverted: sdk.NewCoin("microtwelve", sdk.NewIntFromUint64(0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenDeduction, tokensConverted := keeper.ConvertCoinsForTransfer(tt.args.goCtx, tt.args.msg, tt.args.registryEntry, tt.args.sendAsEntry)
			require.Equal(t, tt.tokensConverted, tokensConverted)
			require.Equal(t, tt.tokenDeduction, tokenDeduction)
		})
	}
}

func TestPrepareToSendConvertedCoins(t *testing.T) {
	maxUInt64 := uint64(18446744073709551615)
	app, appCtx, admin := tokenregistrytest.CreateTestApp(false)

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
		msg           *sdktransfertypes.MsgTransfer
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
			args:            args{sdk.WrapSDKContext(appCtx), &sdktransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(maxUInt64)), Sender: admin}, rowanEntry, microRowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(18446744073700000000)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(184467440737)),
		},
		{
			args:            args{sdk.WrapSDKContext(appCtx), &sdktransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(1000000)), Sender: admin}, rowanEntry, microRowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(0)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(0)),
		},
		{
			args:            args{sdk.WrapSDKContext(appCtx), &sdktransfertypes.MsgTransfer{Token: sdk.NewCoin("twelve", sdk.NewIntFromUint64(10000000000000)), Sender: admin}, decimal12Entry, decimal10Entry},
			tokenDeduction:  sdk.NewCoin("twelve", sdk.NewIntFromUint64(10000000000000)),
			tokensConverted: sdk.NewCoin("microtwelve", sdk.NewIntFromUint64(100000000000)),
		},
		{
			args:            args{sdk.WrapSDKContext(appCtx), &sdktransfertypes.MsgTransfer{Token: sdk.NewCoin("twelve", sdk.NewIntFromUint64(1)), Sender: admin}, decimal12Entry, decimal10Entry},
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
			// TODO: Assert amounts
		})
	}
}
