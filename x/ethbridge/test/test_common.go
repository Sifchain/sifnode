package test

import (
	"bytes"

	"github.com/Sifchain/sifnode/simapp"
	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"math/rand"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/x/supply"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	AddressKey1 = "A58856F0FD53BF058B4909A21AEC019107BA6"
)

//// returns context and app with params set on account keeper
func CreateTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	_ = simapp.AddTestAddrs(app, ctx, 6, initTokens)
	return app, ctx
}

func CreateTestAppEthBridge(isCheckTx bool) (sdk.Context, keeper.Keeper) {
	ctx, app := GetSimApp(isCheckTx)
	return ctx, app.EthBridgeKeeper
}

func GetSimApp(isCheckTx bool) (sdk.Context, *simapp.SimApp) {
	app, ctx := CreateTestApp(isCheckTx)
	return ctx, app
}

func GenerateRandomTokens(numberOfTokens int) []string {
	var tokenList []string
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx", "cacoin", "cbcoin", "ccoin", "cdcoin"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfTokens; i++ {
		// initialize global pseudo random generator
		randToken := tokens[rand.Intn(len(tokens))]

		tokenList = append(tokenList, randToken)
	}
	return tokenList
}

func GenerateAddress(key string) sdk.AccAddress {
	if key == "" {
		key = AddressKey1
	}
	var buffer bytes.Buffer
	buffer.WriteString(key)
	buffer.WriteString(strconv.Itoa(100))
	res, _ := sdk.AccAddressFromHex(buffer.String())
	bech := res.String()
	addr := buffer.String()
	res, err := sdk.AccAddressFromHex(addr)

	if err != nil {
		panic(err)
	}

	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}
	return res
}
