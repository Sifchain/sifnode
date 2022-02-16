package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/keeper"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_Logger(t *testing.T) {
	app, ctx := test.CreateTestApp(false)

	key := sdk.NewKVStoreKey("rowan")
	cdc := codec.BinaryCodec(app.AppCodec())
	accountKeeper := types.AccountKeeper(app.AccountKeeper)
	bankkeeper := types.BankKeeper(app.BankKeeper)
	ps := paramtypes.Subspace{}
	result := keeper.NewKeeper(cdc, key, bankkeeper, accountKeeper, ps)
	res := result.Logger(ctx)
	assert.Equal(t, res, result.Logger(ctx))

}

func TestKeeper_Codec(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	t.Log(ctx)
	key := sdk.NewKVStoreKey("rowan")
	cdc := codec.BinaryCodec(app.AppCodec())
	accountKeeper := types.AccountKeeper(app.AccountKeeper)
	bankkeeper := types.BankKeeper(app.BankKeeper)
	ps := paramtypes.Subspace{}
	result := keeper.NewKeeper(cdc, key, bankkeeper, accountKeeper, ps)

	res := result.Codec()
	assert.Equal(t, res, result.Codec())
}

func TestKeeper_GetAccountKeeper(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	t.Log(ctx)
	key := sdk.NewKVStoreKey("rowan")
	cdc := codec.BinaryCodec(app.AppCodec())
	accountKeeper := types.AccountKeeper(app.AccountKeeper)
	bankkeeper := types.BankKeeper(app.BankKeeper)
	ps := paramtypes.Subspace{}
	result := keeper.NewKeeper(cdc, key, bankkeeper, accountKeeper, ps)

	res := result.GetAccountKeeper()
	t.Log(res)

}

func TestKeeper_HasCoin(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	t.Log(ctx)
	key := sdk.NewKVStoreKey("rowan")
	cdc := codec.BinaryCodec(app.AppCodec())
	accountKeeper := types.AccountKeeper(app.AccountKeeper)
	bankkeeper := types.BankKeeper(app.BankKeeper)
	ps := paramtypes.Subspace{}
	result := keeper.NewKeeper(cdc, key, bankkeeper, accountKeeper, ps)
	user := sdk.AccAddress("addr1_____")
	res := result.HasCoins(ctx, user, sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(1000000))))
	t.Log(res)

}
