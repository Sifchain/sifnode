package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	bridgekeeper "github.com/Sifchain/sifnode/x/ethbridge/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		migratePeggedToken(ctx, tokenMap, app.EthbridgeKeeper)
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
	pools := poolKeeper.GetPools(ctx)

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
	liquidity := poolKeeper.GetLiquidityProviders(ctx)
	// at first check all old denom mapped
	for _, value := range liquidity {
		token := value.Asset.Symbol
		if _, ok := tokenMap[token]; ok {
		} else {
			panic(fmt.Sprintf("new denom for %s not found\n", token))
		}
	}

	for _, value := range liquidity {
		token := value.Asset.Symbol
		if newDenom, ok := tokenMap[token]; ok {
			poolKeeper.DestroyLiquidityProvider(ctx, token, value.LiquidityProviderAddress)
			value.Asset.Symbol = newDenom
			poolKeeper.SetLiquidityProvider(ctx, value)
		}
	}
}

func migratePeggedToken(ctx sdk.Context, tokenMap map[string]string, bridgeKeeper bridgekeeper.Keeper) {
	tokens := bridgeKeeper.GetPeggyToken(ctx).Tokens
	newTokens := []string{}
	for _, token := range tokens {
		if value, ok := tokenMap[token]; ok {
			newTokens = append(newTokens, value)
		} else {
			panic(fmt.Sprintf("new denom for %s not found\n", token))
		}
	}
	bridgeKeeper.SetPeggyToken(ctx, newTokens)
}
