package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ibctransfer/keeper"
	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func TestMsgServer_Transfer(t *testing.T) {
	t.Skip()
	app, ctx, admin := tokenregistrytest.CreateTestApp(false)
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:         "rowan",
		IsWhitelisted: true,
		Decimals:      18,
		IbcDenom:      "microrowan",
		IbcDecimals:   10,
	})
	// TODO: Need to create channel if possible.
	// TODO: Setup funded addresses.
	srv := keeper.NewMsgServerImpl(app.TransferKeeper, app.BankKeeper, app.TokenRegistryKeeper)
	_, err := srv.Transfer(sdk.WrapSDKContext(ctx), &ibctransfertypes.MsgTransfer{
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
		Token: sdk.Coin{
			Denom:  "rowan",
			Amount: sdk.NewInt(int64(1000000)),
		},
		Sender:           admin,
		Receiver:         "",
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: 0,
	})
	require.NoError(t, err)
}

func TestConvertCoins(t *testing.T) {
	ctx := context.Background()

	maxUInt64 := uint64(18446744073709551615)

	rowanEntry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Decimals:      18,
		Denom:         "rowan",
		BaseDenom:     "rowan",
		IbcDenom:      "microrowan",
		IbcDecimals:   10,
	}

	decimal12Entry := tokenregistrytypes.RegistryEntry{
		IsWhitelisted: true,
		Decimals:      12,
		Denom:         "twelve",
		BaseDenom:     "twelve",
		IbcDenom:      "microtwelve",
		IbcDecimals:   10,
	}

	type args struct {
		goCtx         context.Context
		msg           *ibctransfertypes.MsgTransfer
		registryEntry tokenregistrytypes.RegistryEntry
	}
	tests := []struct {
		name            string
		args            args
		tokenDeduction  sdk.Coin
		tokensConverted sdk.Coin
	}{
		{
			args:            args{ctx, &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(maxUInt64))}, rowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(18446744073700000000)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(184467440737)),
		},
		{
			args:            args{ctx, &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(1000000))}, rowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(0)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(0)),
		},
		{
			args:            args{ctx, &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("twelve", sdk.NewIntFromUint64(10000000000000))}, decimal12Entry},
			tokenDeduction:  sdk.NewCoin("twelve", sdk.NewIntFromUint64(10000000000000)),
			tokensConverted: sdk.NewCoin("microtwelve", sdk.NewIntFromUint64(100000000000)),
		},
		{
			args:            args{ctx, &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("twelve", sdk.NewIntFromUint64(1))}, decimal12Entry},
			tokenDeduction:  sdk.NewCoin("twelve", sdk.NewIntFromUint64(0)),
			tokensConverted: sdk.NewCoin("microtwelve", sdk.NewIntFromUint64(0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenDeduction, tokensConverted := keeper.ConvertCoinsForTransfer(tt.args.goCtx, tt.args.msg, tt.args.registryEntry)
			require.Equal(t, tt.tokensConverted, tokensConverted)
			require.Equal(t, tt.tokenDeduction, tokenDeduction)
		})
	}
}
