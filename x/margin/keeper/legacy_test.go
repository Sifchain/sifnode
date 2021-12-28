package keeper_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Sifchain/sifnode/x/margin/keeper"
	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestKeeper_NewLegacyQuerier(t *testing.T) {
	ctx, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	got := keeper.NewLegacyQuerier(marginKeeper)

	require.NotNil(t, got)
	require.Equal(t, reflect.TypeOf(got).String(), "func(types.Context, []string, types.RequestQuery) ([]uint8, error)")

	_, err := got(ctx, []string{"xxx"}, abci.RequestQuery{})

	require.ErrorIs(t, err, sdkerrors.Wrap(types.ErrUnknownRequest, "unknown request"))
}

func TestKeeper_NewLegacyHandler(t *testing.T) {
	ctx, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	handler := keeper.NewLegacyHandler(marginKeeper)

	require.NotNil(t, handler)
	require.Equal(t, reflect.TypeOf(handler).String(), "types.Handler")

	var (
		msgOpenLong  sdk.Msg = &types.MsgOpenLong{}
		msgCloseLong sdk.Msg = &types.MsgCloseLong{}
		msgOther     sdk.Msg
	)

	newLegacyHandlerTests := []struct {
		name      string
		msg       sdk.Msg
		err       error
		errString error
	}{
		{
			name:      "msg open long",
			msg:       msgOpenLong,
			errString: errors.New(": pool does not exist"),
		},
		{
			name:      "msg close long",
			msg:       msgCloseLong,
			errString: errors.New("mtp not found"),
		},
		{
			name: "msg other",
			msg:  msgOther,
			err:  sdkerrors.Wrap(types.ErrUnknownRequest, "unknown request"),
		},
	}

	for _, tt := range newLegacyHandlerTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, got := handler(ctx, tt.msg)

			if tt.errString != nil {
				require.EqualError(t, got, tt.errString.Error())
			} else if tt.err == nil {
				require.NoError(t, got)
			} else {
				require.ErrorIs(t, got, tt.err)
			}
		})
	}
}
