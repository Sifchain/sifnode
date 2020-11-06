package keeper_test

import (
	"fmt"
	"github.com/Sifchain/sifnode/simapp"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/rand"
	"testing"
	"time"
)

func GenerateRandomPool(numberOfPools int) []types.Pool {
	var poolList []types.Pool
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfPools; i++ {
		// initialize global pseudo random generator
		externalToken := tokens[rand.Intn(len(tokens))]
		externalAsset := types.NewAsset("ROWAN", "c"+"ROWAN"+externalToken, externalToken)
		pool, err := types.NewPool(externalAsset, 1000, 100, 1)
		if err != nil {
			fmt.Println("Error Generating new pool :", err)
		}
		poolList = append(poolList, pool)
	}
	return poolList
}

// returns context and app with params set on account keeper
func CreateTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	// UNDONE: is this needed?
	initTokens := sdk.TokensFromConsensusPower(1000)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	// UNDONE: is this needed?
	_ = simapp.AddTestAddrs(app, ctx, 6, initTokens)

	return app, ctx
}

func TestKeeper_SetPoolIntegration(t *testing.T) {

	pool := GenerateRandomPool(1)[0]
	app, ctx := CreateTestApp(false)
	keeper := app.ClpKeeper
	//ctx, keeper := CreateTestInputDefault(t, false, 1000)
	keeper.SetPool(ctx, pool)
	getpool, err := keeper.GetPool(ctx, pool.ExternalAsset.Ticker)
	assert.NoError(t, err, "Error in get pool")
	assert.Equal(t, getpool, pool)
	assert.Equal(t, keeper.ExistsPool(ctx, pool.ExternalAsset.Ticker), true)
}


