package app

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge"

	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

//go:embed denom_mapping_peggy1_to_peggy2.json
var tokenMap string

// GetTokenMigrationFunc return a function to migrate token denom
func GetTokenMigrationFunc(app *SifchainApp) func(ctx sdk.Context, plan upgradetypes.Plan) {
	return func(ctx sdk.Context, plan upgradetypes.Plan) {
		ctx.Logger().Info("Starting to execute token migration for balance pool and liquidity")

		ExportAppState("changePoolFormula", app, ctx)

		tokenMap := ReadTokenMapJSON()

		MigrateBalance(ctx, tokenMap, app.BankKeeper)
		migratePool(ctx, tokenMap, app.ClpKeeper)
		migrateLiquidity(ctx, tokenMap, app.ClpKeeper)
	}
}

func ReadTokenMapJSON() map[string]string {
	jsonData := map[string]string{}
	err := json.Unmarshal([]byte(tokenMap), &jsonData)
	if err != nil {
		panic("fail to parse token migration file")
	}
	return jsonData
}

func MigrateBalance(ctx sdk.Context, tokenMap map[string]string, bankKeeper bankkeeper.Keeper) {
	addresses := []sdk.AccAddress{}
	coins := []sdk.Coin{}

	bankKeeper.IterateAllBalances(ctx, func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {
		addresses = append(addresses, address)
		coins = append(coins, coin)
		return false
	})

	for index, address := range addresses {

		coin := coins[index]
		amount := coin.Amount

		// set the balance for new denom
		if newDenom, ok := tokenMap[coin.Denom]; ok {
			if newDenom == coin.Denom {
				continue
			}

			// send old coins to module
			err := bankKeeper.SendCoinsFromAccountToModule(ctx, address, ethbridge.ModuleName, sdk.NewCoins(coin))
			if err != nil {
				panic(err)
			}

			err = bankKeeper.BurnCoins(ctx, ethbridge.ModuleName, sdk.NewCoins(coin))
			if err != nil {
				panic(err)
			}

			newCoins := sdk.NewCoins(sdk.NewCoin(newDenom, amount))

			// can't set balance directly, we mint to module, then transfer to address
			err = bankKeeper.MintCoins(ctx, ethbridge.ModuleName, newCoins)
			if err != nil {
				panic("failed to mint coins during token migration")
			}

			err = bankKeeper.SendCoinsFromModuleToAccount(ctx, ethbridge.ModuleName, address, newCoins)
			if err != nil {
				panic("failed to set balance during token migration")
			}
		} else {
			ctx.Logger().Error("new denom for %s not found\n", coin.Denom)
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
