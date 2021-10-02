package keeper_test

import (
	"context"
	"math"
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper"
	scibctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	scibctransfermocks "github.com/Sifchain/sifnode/x/ibctransfer/types/mocks"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

/* Test that when a conversion is needed the right amounts are converted before sending to underlying SDK Transfer. */
func TestMsgServer_Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	bankKeeper := scibctransfermocks.NewMockBankKeeper(ctrl)
	msgSrv := scibctransfermocks.NewMockMsgServer(ctrl)
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	addrs, _ := test.CreateTestAddrs(2)
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:                "rowan",
		Decimals:             18,
		IbcCounterpartyDenom: "xrowan",
		Permissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCEXPORT},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:       "xrowan",
		Decimals:    10,
		UnitDenom:   "rowan",
		Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCEXPORT},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:       "misconfigured",
		Decimals:    18,
		Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCEXPORT},
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:                "ceth",
		Decimals:             18,
		IbcCounterpartyDenom: "ceth",
		Permissions:          []tokenregistrytypes.Permission{},
	})
	rowanAmount, ok := sdk.NewIntFromString("1234567891123456789")
	require.True(t, ok)
	rowanAmountEscrowed, ok := sdk.NewIntFromString("1234567891100000000")
	require.True(t, ok)
	xrowanAmount, ok := sdk.NewIntFromString("12345678911")
	require.True(t, ok)
	packetOverflowAmount := sdk.NewIntFromUint64(math.MaxUint64).Add(sdk.NewInt(1))
	rowanSmallest, ok := sdk.NewIntFromString("183456789")
	require.True(t, ok)
	rowanTooSmall, ok := sdk.NewIntFromString("12345678")
	require.True(t, ok)
	tooLargeToSend, ok := sdk.NewIntFromString("940000000000000000000000000")
	require.True(t, ok)
	tooLargeToSendAs, ok := sdk.NewIntFromString("9400000000000000000")
	require.True(t, ok)
	tooLargeToSend2, ok := sdk.NewIntFromString("8940000000000000000000000000")
	require.True(t, ok)
	tooLargeToSendAs2, ok := sdk.NewIntFromString("89400000000000000000")
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
			setupBankKeeperCalls: func() {
				bankKeeper.EXPECT().SendCoins(gomock.Any(), addrs[0], scibctransfertypes.GetEscrowAddress("transfer", "channel-0"), sdk.NewCoins(sdk.NewCoin("rowan", rowanAmountEscrowed))).Return(nil)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), scibctransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("xrowan", xrowanAmount))).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), scibctransfertypes.ModuleName, addrs[0], sdk.NewCoins(sdk.NewCoin("xrowan", xrowanAmount))).Return(nil)
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
				bankKeeper.EXPECT().SendCoins(gomock.Any(), addrs[0], scibctransfertypes.GetEscrowAddress("transfer", "channel-0"), sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(100000000)))).Return(nil)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), scibctransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("xrowan", sdk.NewInt(1)))).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), scibctransfertypes.ModuleName, addrs[0], sdk.NewCoins(sdk.NewCoin("xrowan", sdk.NewInt(1)))).Return(nil)
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
		{
			name:       "transfer denom without ibc export permission",
			err:        tokenregistrytypes.ErrPermissionDenied,
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("ceth", sdk.NewInt(1)),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupBankKeeperCalls: func() {},
			setupMsgServerCalls:  func() {},
		},
		{
			name:       "transfer denom is not whitelisted",
			err:        tokenregistrytypes.ErrPermissionDenied,
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("caave", sdk.NewInt(1)),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupBankKeeperCalls: func() {},
			setupMsgServerCalls:  func() {},
		},
		{
			name:       "transfer denom alias with unit denom set in registry",
			err:        tokenregistrytypes.ErrPermissionDenied,
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("xrowan", sdk.NewInt(1)),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupBankKeeperCalls: func() {},
			setupMsgServerCalls:  func() {},
		},
		{
			name:       "transfer amount too large to send without conversion",
			err:        scibctransfertypes.ErrAmountTooLargeToSend,
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("misconfigured", packetOverflowAmount),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupBankKeeperCalls: func() {},
			setupMsgServerCalls:  func() {},
		},
		{
			name:       "transfer amount too large to transfer edge case 1",
			err:        scibctransfertypes.ErrAmountTooLargeToSend,
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("rowan", tooLargeToSend),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupMsgServerCalls: func() {
				msgSrv.EXPECT().Transfer(gomock.Any(), &sdktransfertypes.MsgTransfer{
					SourcePort:       "transfer",
					SourceChannel:    "channel-0",
					Token:            sdk.NewCoin("xrowan", tooLargeToSendAs),
					Sender:           addrs[0].String(),
					Receiver:         addrs[1].String(),
					TimeoutHeight:    clienttypes.NewHeight(0, 0),
					TimeoutTimestamp: 0,
				})
			},
			setupBankKeeperCalls: func() {
				bankKeeper.EXPECT().SendCoins(gomock.Any(), addrs[0], scibctransfertypes.GetEscrowAddress("transfer", "channel-0"), sdk.NewCoins(sdk.NewCoin("rowan", tooLargeToSend))).Return(nil)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), scibctransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("xrowan", tooLargeToSendAs))).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), scibctransfertypes.ModuleName, addrs[0], sdk.NewCoins(sdk.NewCoin("xrowan", tooLargeToSendAs))).Return(nil)
			},
		},
		{
			name:       "transfer amount too large to transfer edge case 2",
			err:        scibctransfertypes.ErrAmountTooLargeToSend,
			bankKeeper: bankKeeper,
			msgSrv:     msgSrv,
			msg: sdktransfertypes.NewMsgTransfer(
				"transfer",
				"channel-0",
				sdk.NewCoin("rowan", tooLargeToSend2),
				addrs[0],
				addrs[1].String(),
				clienttypes.NewHeight(0, 0),
				0,
			),
			setupMsgServerCalls: func() {
				msgSrv.EXPECT().Transfer(gomock.Any(), &sdktransfertypes.MsgTransfer{
					SourcePort:       "transfer",
					SourceChannel:    "channel-0",
					Token:            sdk.NewCoin("xrowan", tooLargeToSendAs2),
					Sender:           addrs[0].String(),
					Receiver:         addrs[1].String(),
					TimeoutHeight:    clienttypes.NewHeight(0, 0),
					TimeoutTimestamp: 0,
				})
			},
			setupBankKeeperCalls: func() {
				bankKeeper.EXPECT().SendCoins(gomock.Any(), addrs[0], scibctransfertypes.GetEscrowAddress("transfer", "channel-0"), sdk.NewCoins(sdk.NewCoin("rowan", tooLargeToSend2))).Return(nil)
				bankKeeper.EXPECT().MintCoins(gomock.Any(), scibctransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("xrowan", tooLargeToSendAs2))).Return(nil)
				bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), scibctransfertypes.ModuleName, addrs[0], sdk.NewCoins(sdk.NewCoin("xrowan", tooLargeToSendAs2))).Return(nil)
			},
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMsgServerCalls()
			tc.setupBankKeeperCalls()
			srv := keeper.NewMsgServerImpl(tc.msgSrv, tc.bankKeeper, app.TokenRegistryKeeper)
			_, err := srv.Transfer(sdk.WrapSDKContext(ctx), tc.msg)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestConvertCoins(t *testing.T) {
	ctx := context.Background()
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
	decimal12Entry := tokenregistrytypes.RegistryEntry{
		Decimals:             12,
		Denom:                "twelve",
		BaseDenom:            "twelve",
		IbcCounterpartyDenom: "microtwelve",
	}
	decimal10Entry := tokenregistrytypes.RegistryEntry{
		Decimals:  10,
		Denom:     "microtwelce",
		BaseDenom: "microtwelve",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tokenDeduction, tokensConverted := helpers.ConvertCoinsForTransfer(tt.args.msg, &tt.args.registryEntry, &tt.args.sendAsEntry)
			require.Equal(t, tt.tokensConverted, tokensConverted)
			require.Equal(t, tt.tokenDeduction, tokenDeduction)
		})
	}
}

func TestPrepareToSendConvertedCoins(t *testing.T) {
	maxUInt64 := uint64(18446744073709551615)
	app, appCtx, admin := tokenregistrytest.CreateTestApp(false)
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
	decimal12Entry := tokenregistrytypes.RegistryEntry{
		Decimals:             12,
		Denom:                "twelve",
		BaseDenom:            "twelve",
		IbcCounterpartyDenom: "microtwelve",
	}
	decimal10Entry := tokenregistrytypes.RegistryEntry{
		Decimals:  10,
		Denom:     "microtwelce",
		BaseDenom: "microtwelve",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tokenDeduction, tokensConverted := helpers.ConvertCoinsForTransfer(tt.args.msg, &tt.args.registryEntry, &tt.args.sendAsEntry)
			require.Equal(t, tt.tokensConverted, tokensConverted)
			require.Equal(t, tt.tokenDeduction, tokenDeduction)
			initCoins := sdk.NewCoins(tokenDeduction)
			sender, err := sdk.AccAddressFromBech32(admin)
			require.NoError(t, err)
			err = app.BankKeeper.AddCoins(appCtx, sender, initCoins)
			require.NoError(t, err)
			err = helpers.PrepareToSendConvertedCoins(sdk.WrapSDKContext(appCtx), tt.args.msg, tokenDeduction, tokensConverted, app.BankKeeper)
			require.NoError(t, err)
		})
	}
}
