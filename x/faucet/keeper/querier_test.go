package keeper_test

import (
	"github.com/Sifchain/sifnode/x/clp/test"
	keeper2 "github.com/Sifchain/sifnode/x/faucet/keeper"
	"github.com/Sifchain/sifnode/x/faucet/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestKeeper_GetBalance(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	faucetfundingAmount := sdk.NewInt(100000000000)
	keeper := app.FaucetKeeper
	faucetFundingCoins := sdk.Coins{sdk.NewCoin("rowan", faucetfundingAmount)}
	c, err := keeper.GetBankKeeper().AddCoins(ctx, types.GetFaucetModuleAddress(), faucetFundingCoins)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	q := keeper2.NewQuerier(keeper)
	assert.NotNil(t, q)
	balance, err := q(ctx, []string{types.QueryBalance}, abci.RequestQuery{})
	assert.NoError(t, err)
	var coins sdk.Coins
	app.Codec().MustUnmarshalJSON(balance, &coins)
	assert.Equal(t, coins.String(), faucetFundingCoins.String())
}
