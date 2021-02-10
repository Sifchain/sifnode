package test

import (
	"bytes"
	"fmt"
	"github.com/Sifchain/sifnode/simapp"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/Sifchain/sifnode/x/clp/types"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Constants for test scripts only .
//
const (
	AddressKey1 = "A58856F0FD53BF058B4909A21AEC019107BA6"
	AddressKey2 = "A58856F0FD53BF058B4909A21AEC019107BA7"
	AddressKey3 = "A58856F0FD53BF058B4909A21AEC019107BA9"
)

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // distr
	return cdc
}

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

func CreateTestAppClp(isCheckTx bool) (sdk.Context, keeper.Keeper) {
	ctx, app := GetSimApp(isCheckTx)
	return ctx, app.ClpKeeper
}

func GetSimApp(isCheckTx bool) (sdk.Context, *simapp.SimApp) {
	app, ctx := CreateTestApp(isCheckTx)
	return ctx, app

}

func GenerateRandomPool(numberOfPools int) []types.Pool {
	var poolList []types.Pool
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfPools; i++ {
		// initialize global pseudo random generator
		externalToken := tokens[rand.Intn(len(tokens))]
		externalAsset := types.NewAsset(trimFirstRune(externalToken))
		pool, err := types.NewPool(externalAsset, sdk.NewUint(1000), sdk.NewUint(100), sdk.NewUint(1))
		if err != nil {
			fmt.Println("Error Generating new pool :", err)
		}
		poolList = append(poolList, pool)
	}
	return poolList
}

func GenerateRandomLP(numberOfLp int) []types.LiquidityProvider {
	var lpList []types.LiquidityProvider
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfLp; i++ {
		externalToken := tokens[rand.Intn(len(tokens))]
		asset := types.NewAsset(trimFirstRune(externalToken))
		lpAddess, _ := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
		lp := types.NewLiquidityProvider(asset, sdk.NewUint(1), lpAddess)
		lpList = append(lpList, lp)
	}
	return lpList
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return strings.ToLower(s[i:])
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

func GenerateWhitelistAddress(key string) sdk.AccAddress {
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
