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
		Denom:                "rowan",
		IsWhitelisted:        true,
		Decimals:             18,
		IbcCounterPartyDenom: "microrowan",
	})
	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:         "microrowan",
		IsWhitelisted: true,
		Decimals:      10,
		UnitDenom:     "rowan",
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
			args:            args{ctx, &ibctransfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(maxUInt64))}, rowanEntry, microRowanEntry},
			tokenDeduction:  sdk.NewCoin("rowan", sdk.NewIntFromUint64(18446744073700000000)),
			tokensConverted: sdk.NewCoin("microrowan", sdk.NewIntFromUint64(184467440737)),
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
