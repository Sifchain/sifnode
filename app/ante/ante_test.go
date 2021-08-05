package ante_test

import (
	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/app/ante"
	dispensationtypes "github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"testing"
)

func TestReduceGasPriceDecorator_AnteHandle(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000)
	addrs := sifapp.AddTestAddrs(app, ctx, 6, initTokens)

	decorator := ante.ReduceGasPriceDecorator{}

	highGasPrice := sdk.DecCoin{
		Denom:  "rowan",
		Amount: sdk.MustNewDecFromStr("0.5"),
	}

	loweredGasPrice := sdk.DecCoin{
		Denom:  "rowan",
		Amount: sdk.MustNewDecFromStr("0.00000005"),
	}

	ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(highGasPrice))
	dispensationCreateMsg := dispensationtypes.NewMsgCreateDistribution(addrs[0], dispensationtypes.DistributionType_DISTRIBUTION_TYPE_AIRDROP, []banktypes.Output{}, "")
	dispensationRunMsg := dispensationtypes.NewMsgRunDistribution(addrs[0].String(), "airdrop", dispensationtypes.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
	otherMsg := banktypes.NewMsgSend(addrs[0], addrs[1], sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewIntFromUint64(100))))
	// next doesn't accept err, it is only called if decorator does not return error, it passes ctx to decorator caller
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) { return ctx, nil }

	tt := []struct {
		name             string
		ctx              sdk.Context
		msgs             []sdk.Msg
		expectedGasPrice sdk.DecCoin
		err              bool
	}{
		{"no messages", ctx, []sdk.Msg{}, highGasPrice, false},
		{"dispensation create", ctx, []sdk.Msg{&dispensationCreateMsg}, loweredGasPrice, false},
		{"dispensation create with extra msg", ctx, []sdk.Msg{&dispensationCreateMsg, otherMsg}, highGasPrice, true},
		{"dispensation run", ctx, []sdk.Msg{&dispensationRunMsg}, loweredGasPrice, false},
		{"dispensation run with extra msg", ctx, []sdk.Msg{&dispensationRunMsg, otherMsg}, highGasPrice, true},
		{"other message without dispensation", ctx, []sdk.Msg{otherMsg}, highGasPrice, false},
		{"other messages without dispensation", ctx, []sdk.Msg{otherMsg, otherMsg}, highGasPrice, false},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tx := legacytx.StdTx{
				Msgs:          tc.msgs,
				Fee:           legacytx.StdFee{},
				Signatures:    nil,
				Memo:          "",
				TimeoutHeight: 0,
			}

			newCtx, err := decorator.AnteHandle(ctx, tx, false, next)
			require.Equal(t, err != nil, tc.err)
			require.NotNil(t, newCtx)
			require.Equal(t, tc.expectedGasPrice.String(), newCtx.MinGasPrices().String())
		})
	}
}
