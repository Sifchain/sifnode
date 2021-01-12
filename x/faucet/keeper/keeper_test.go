package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/faucet/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetWithdrawnAmountInEpoch(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	user := test.GenerateAddress("")
	faucetRequestAmount := sdk.NewInt(1000)
	keeper := app.FaucetKeeper
	err := keeper.SetWithdrawnAmountInEpoch(ctx, user.String(), faucetRequestAmount, types.FaucetToken)
	assert.NoError(t, err)
}

func TestKeeper_GetWithdrawnAmountInEpoch(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	user := test.GenerateAddress("")
	faucetRequestAmount := sdk.NewInt(1000)
	keeper := app.FaucetKeeper
	err := keeper.SetWithdrawnAmountInEpoch(ctx, user.String(), faucetRequestAmount, types.FaucetToken)
	assert.NoError(t, err)
	amt, err := keeper.GetWithdrawnAmountInEpoch(ctx, user.String(), types.FaucetToken)
	assert.NoError(t, err)
	assert.Equal(t, amt.String(), faucetRequestAmount.String())
}

func TestKeeper_StartNextEpoch(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	user := test.GenerateAddress("")
	faucetRequestAmount := sdk.NewInt(1000)
	keeper := app.FaucetKeeper
	err := keeper.SetWithdrawnAmountInEpoch(ctx, user.String(), faucetRequestAmount, types.FaucetToken)
	assert.NoError(t, err)
	keeper.StartNextEpoch(ctx)
	amt, err := keeper.GetWithdrawnAmountInEpoch(ctx, user.String(), types.FaucetToken)
	assert.NoError(t, err)
	assert.Equal(t, sdk.ZeroInt().String(), amt.String())
}

func TestKeeper_CanRequest(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	user := test.GenerateAddress("")
	faucetRequestAmount, ok := sdk.NewIntFromString(types.MaxWithdrawAmountPerEpoch)
	assert.True(t, ok)
	keeper := app.FaucetKeeper
	unitCoins := sdk.Coins{sdk.NewCoin(types.FaucetToken, sdk.NewInt(1))}
	ok, err := keeper.CanRequest(ctx, user.String(), unitCoins)
	assert.True(t, ok)
	assert.NoError(t, err)
	err = keeper.SetWithdrawnAmountInEpoch(ctx, user.String(), faucetRequestAmount, types.FaucetToken)
	assert.NoError(t, err)
	ok, err = keeper.CanRequest(ctx, user.String(), unitCoins)
	assert.False(t, ok)
	assert.Error(t, err)
}

func TestKeeper_ExecuteRequest(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	user := test.GenerateAddress("")
	faucetRequestAmount, ok := sdk.NewIntFromString(types.MaxWithdrawAmountPerEpoch)
	assert.True(t, ok)
	keeper := app.FaucetKeeper
	coins := sdk.Coins{sdk.NewCoin(types.FaucetToken, faucetRequestAmount)}
	ok, err := keeper.CanRequest(ctx, user.String(), coins)
	assert.True(t, ok)
	assert.NoError(t, err)
	ok, err = keeper.ExecuteRequest(ctx, user.String(), coins)
	assert.True(t, ok)
	assert.NoError(t, err)
	ok, err = keeper.CanRequest(ctx, user.String(), coins)
	assert.False(t, ok)
	assert.Error(t, err)
}
