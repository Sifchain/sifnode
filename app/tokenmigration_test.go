package app_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBankMigration(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	addrs, _ := test.CreateTestAddrs(1000)
	// Create Peggy 1 balances for all denoms in map
	tokenMap := sifapp.ReadTokenMapJSON()
	for peggy1, _ := range tokenMap {
		err := app.BankKeeper.MintCoins(ctx, ethbridge.ModuleName, sdk.NewCoins(sdk.NewCoin(peggy1, sdk.NewInt(1))))
		require.NoError(t, err)
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ethbridge.ModuleName, addrs[0], sdk.NewCoins(sdk.NewCoin(peggy1, sdk.NewInt(1))))
		require.NoError(t, err)
	}
	sifapp.MigrateBalance(ctx, tokenMap, app.BankKeeper)
	for peggy1, peggy2 := range tokenMap {
		if peggy1 != peggy2 {
			coin := app.BankKeeper.GetBalance(ctx, addrs[0], peggy1)
			assert.True(t, coin.IsZero(), "coin %s", coin.String())
			supply := app.BankKeeper.GetSupply(ctx, peggy1)
			assert.True(t, supply.IsZero(), "supply %s", supply.String())
		}
		coin := app.BankKeeper.GetBalance(ctx, addrs[0], peggy2)
		require.True(t, coin.IsEqual(sdk.NewCoin(peggy2, sdk.NewInt(1))), "coin %s", coin.String())
		supply := app.BankKeeper.GetSupply(ctx, peggy2)
		assert.True(t, coin.IsEqual(sdk.NewCoin(peggy2, sdk.NewInt(1))), "supply %s", supply.String())
	}
}
