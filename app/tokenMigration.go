package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// GetTokenMigrationFunc return a function to migrate token denom
func GetTokenMigrationFunc(app *SifchainApp) func(ctx sdk.Context, plan upgradetypes.Plan) {
	return func(ctx sdk.Context, plan upgradetypes.Plan) {
		ctx.Logger().Info("Starting to execute token migration for balance pool and liquidity")

		ExportAppState("changePoolFormula", app, ctx)

		tokenMap := readTokenMapJSON()

		migrateBalance(ctx, tokenMap, app.BankKeeper)
		migratePool(ctx, tokenMap, app.ClpKeeper)
		migrateLiquidity(ctx, tokenMap, app.ClpKeeper)
	}
}

func readTokenMapJSON() map[string]string {
	data, err := ioutil.ReadFile("./token_migration.json")
	if err != nil {
		panic(err)
	}

	jsonData := map[string]string{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		panic("fail to parse token migration file")
	}
	return jsonData
}

func getAll(addresses *[]sdk.AccAddress, coins *[]sdk.Coin) func(address sdk.AccAddress, coin sdk.Coin) bool {
	return func(address sdk.AccAddress, coin sdk.Coin) bool {
		*addresses = append(*addresses, address)
		*coins = append(*coins, coin)
		return true
	}
}

func migrateBalance(ctx sdk.Context, tokenMap map[string]string, bankKeeper bankkeeper.Keeper) {
	addresses := []sdk.AccAddress{}
	coins := []sdk.Coin{}

	bankKeeper.IterateAllBalances(ctx, getAll(&addresses, &coins))

	for index, address := range addresses {

		coin := coins[index]
		amount := coin.Amount
		// clear the balance for old denom
		coin.Amount = sdk.NewInt(0)
		err := bankKeeper.SetBalance(ctx, address, coin)
		if err != nil {
			panic("failed to set balance during token migration")
		}

		// set the balance for new denom
		if value, ok := tokenMap[coin.Denom]; ok {
			coin = sdk.NewCoin(value, amount)
			err = bankKeeper.SetBalance(ctx, address, coin)
			if err != nil {
				panic("failed to set balance during token migration")
			}
		} else {
			panic(fmt.Sprintf("new denom for %s not found\n", coin.Denom))
		}
	}

}

func migratePool(ctx sdk.Context, tokenMap map[string]string, poolKeeper clpkeeper.Keeper) {
	pools, _, err := poolKeeper.GetPoolsPaginated(ctx, &query.PageRequest{Limit: math.MaxUint64})
	if err != nil {
		panic("failed to get pools during token migration")
	}

	// at first check all old denom mapped
	for _, value := range pools {
		token := value.ExternalAsset.Symbol
		if _, ok := tokenMap[token]; ok {
		} else {
			panic(fmt.Sprintf("new denom for %s not found\n", token))
		}
	}

	for _, value := range pools {
		token := value.ExternalAsset.Symbol
		if newDenom, ok := tokenMap[token]; ok {
			err := poolKeeper.DestroyPool(ctx, token)
			if err != nil {
				panic("failed to destroy pool during token migration")
			}
			value.ExternalAsset.Symbol = newDenom
			err = poolKeeper.SetPool(ctx, value)
			if err != nil {
				panic("failed to set pool during token migration")
			}

		}
	}
}

func migrateLiquidity(ctx sdk.Context, tokenMap map[string]string, poolKeeper clpkeeper.Keeper) {
	iterator := poolKeeper.GetLiquidityProviderIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {

		key := string(iterator.Key())
		_, _, err := types.ParsePoolKey(key)

		if err != nil {
			panic(err.Error())
		}
	}

	iterator = poolKeeper.GetLiquidityProviderIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		key := string(iterator.Key())
		symbol, lp, _ := types.ParsePoolKey(key)

		if newDenom, ok := tokenMap[symbol]; ok {
			poolKeeper.DestroyLiquidityProvider(ctx, newDenom, lp)
			poolKeeper.SetRawLiquidityProvider(ctx, newDenom, lp, iterator.Value())
		}
	}
}
