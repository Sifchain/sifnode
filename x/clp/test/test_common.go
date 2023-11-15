package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/clp/types"
)

// Constants for test scripts only .
const (
	AddressKey1 = "A58856F0FD53BF058B4909A21AEC019107BA6"
	AddressKey2 = "A58856F0FD53BF058B4909A21AEC019107BA7"
	AddressKey3 = "A58856F0FD53BF058B4909A21AEC019107BA9"
)

// CreateTestApp returns context and app with params set on account keeper
func CreateTestApp(isCheckTx bool) (*sifapp.SifchainApp, sdk.Context) {
	app := sifapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)
	_ = sifapp.AddTestAddrs(app, ctx, 6, initTokens)
	return app, ctx
}
func CreateTestAppClp(isCheckTx bool) (sdk.Context, *sifapp.SifchainApp) {
	return CreateTestAppClpWithBlacklist(isCheckTx, []sdk.AccAddress{})
}

func CreateTestAppClpWithBlacklist(isCheckTx bool, blacklist []sdk.AccAddress) (sdk.Context, *sifapp.SifchainApp) {
	sifapp.SetConfig(false)
	app := sifapp.SetupWithBlacklist(isCheckTx, blacklist)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.TokenRegistryKeeper.SetRegistry(ctx, tokenregistrytypes.Registry{
		Entries: []*tokenregistrytypes.RegistryEntry{
			{Denom: "ceth", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			{Denom: "cdash", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			{Denom: "eth", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			{Denom: "cacoin", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			{Denom: "dash", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			{Denom: "atom", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			{Denom: "cusdc", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			{Denom: "rowan"},
		},
	})
	app.ClpKeeper.SetPmtpRateParams(ctx, types.PmtpRateParams{
		PmtpPeriodBlockRate:    sdk.OneDec(),
		PmtpCurrentRunningRate: sdk.OneDec(),
	})
	app.ClpKeeper.SetParams(ctx, types.Params{
		MinCreatePoolThreshold: 100,
	})
	app.ClpKeeper.SetPmtpParams(ctx, &types.PmtpParams{
		PmtpPeriodGovernanceRate: sdk.OneDec(),
		PmtpPeriodEpochLength:    1,
		PmtpPeriodStartBlock:     1,
		PmtpPeriodEndBlock:       2,
	})
	app.ClpKeeper.SetRewardParams(ctx, &types.RewardParams{
		LiquidityRemovalLockPeriod:   0, // 0 blocks
		LiquidityRemovalCancelPeriod: 2, // 2 blocks
		RewardPeriodStartTime:        "",
		RewardPeriods:                nil,
		RewardsLockPeriod:            12 * 60 * 24 * 14, // 14 days,
	})
	liquidityProtectionParam := app.ClpKeeper.GetLiquidityProtectionParams(ctx)
	liquidityProtectionParam.MaxRowanLiquidityThreshold = sdk.ZeroUint()
	app.ClpKeeper.SetLiquidityProtectionParams(ctx, liquidityProtectionParam)
	app.ClpKeeper.SetProviderDistributionParams(ctx, &types.ProviderDistributionParams{
		DistributionPeriods: nil,
	})
	return ctx, app
}

func CreateTestAppClpFromGenesis(isCheckTx bool, genesisTransformer func(*sifapp.SifchainApp, sifapp.GenesisState) sifapp.GenesisState) (sdk.Context, *sifapp.SifchainApp) {
	sifapp.SetConfig(false)
	app := sifapp.SetupFromGenesis(isCheckTx, genesisTransformer)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})

	app.ClpKeeper.SetPmtpRateParams(ctx, types.PmtpRateParams{
		PmtpPeriodBlockRate:    sdk.OneDec(),
		PmtpCurrentRunningRate: sdk.OneDec(),
	})
	app.ClpKeeper.SetParams(ctx, types.Params{
		MinCreatePoolThreshold: 100,
	})
	app.ClpKeeper.SetPmtpParams(ctx, &types.PmtpParams{
		PmtpPeriodGovernanceRate: sdk.OneDec(),
		PmtpPeriodEpochLength:    1,
		PmtpPeriodStartBlock:     1,
		PmtpPeriodEndBlock:       2,
	})
	app.ClpKeeper.SetRewardParams(ctx, &types.RewardParams{
		LiquidityRemovalLockPeriod:   0,
		LiquidityRemovalCancelPeriod: 2,
		RewardPeriodStartTime:        "",
		RewardPeriods:                nil,
	})
	app.ClpKeeper.SetProviderDistributionParams(ctx, &types.ProviderDistributionParams{
		DistributionPeriods: nil,
	})
	return ctx, app
}

func GenerateRandomPool(numberOfPools int) []types.Pool {
	var poolList []types.Pool
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfPools; i++ {
		// initialize global pseudo random generator
		externalToken := tokens[rand.Intn(len(tokens))]
		externalAsset := types.NewAsset(TrimFirstRune(externalToken))
		pool := types.NewPool(&externalAsset, sdk.NewUint(1000), sdk.NewUint(100), sdk.NewUint(1))
		poolList = append(poolList, pool)
	}
	return poolList
}

func GenerateRandomLPWithUnitsAndAsset(poolUnitss []uint64, asset types.Asset) []*types.LiquidityProvider {
	lpList := make([]*types.LiquidityProvider, len(poolUnitss))
	for i, poolUnits := range poolUnitss {
		address := GenerateAddress2(fmt.Sprintf("%d%d%d%d", i, i, i, i))
		lp := types.NewLiquidityProvider(&asset, sdk.NewUint(poolUnits), address, 0)
		lpList[i] = &lp
	}

	return lpList
}

func GenerateRandomLPWithUnits(poolUnitss []uint64) []*types.LiquidityProvider {
	lpList := make([]*types.LiquidityProvider, len(poolUnitss))
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx"}

	rand.Seed(time.Now().Unix())

	for i, poolUnits := range poolUnitss {
		externalToken := tokens[rand.Intn(len(tokens))]
		asset := types.NewAsset(TrimFirstRune(externalToken))
		address := GenerateAddress(fmt.Sprintf("%d", i))
		lp := types.NewLiquidityProvider(&asset, sdk.NewUint(poolUnits), address, 0)
		lpList[i] = &lp
	}

	return lpList
}
func genTokens(n int) []string {
	var runes = []rune("abcdefghijklmnopqrstuvwxyz")
	set := make(map[string]bool, n)

	for len(set) != n {
		token := make([]rune, 6)
		for i := range token {
			token[i] = runes[rand.Intn(len(runes))]
		}
		set[string(token)] = true
	}

	var strings = make([]string, n)
	i := 0
	for str := range set {
		strings[i] = str
		i++
	}

	return strings
}

func GeneratePoolsSetLPs(keeper clpkeeper.Keeper, ctx sdk.Context, nPools, nLPs int) []*types.Pool {
	tokens := genTokens(nPools)

	rand.Seed(time.Now().Unix())
	poolList := make([]*types.Pool, nPools)
	for i := 0; i < nPools; i++ {
		externalToken := tokens[i]
		externalAsset := types.NewAsset(TrimFirstRune(externalToken))

		poolUnits := make([]uint64, nLPs)
		totalPoolUnits := sdk.ZeroUint()
		for i := 0; i < nLPs; i++ {
			val := uint64(rand.Int31())
			poolUnits[i] = val
			totalPoolUnits = totalPoolUnits.Add(sdk.NewUint(val))
		}

		lps := GenerateRandomLPWithUnitsAndAsset(poolUnits, externalAsset)
		for _, lp := range lps {
			keeper.SetLiquidityProvider(ctx, lp)
		}

		pool := types.NewPool(&externalAsset, sdk.NewUint(100000000000*uint64(i+1)), sdk.NewUint(100*uint64(i+1)), totalPoolUnits)
		err := keeper.SetPool(ctx, &pool)
		if err != nil {
			panic(err)
		}

		poolList[i] = &pool
	}

	return poolList
}

func GenerateRandomLP(numberOfLp int) []*types.LiquidityProvider {
	var lpList []*types.LiquidityProvider
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfLp; i++ {
		externalToken := tokens[rand.Intn(len(tokens))]
		asset := types.NewAsset(TrimFirstRune(externalToken))
		lpAddress, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
		if err != nil {
			panic(err)
		}
		lp := types.NewLiquidityProvider(&asset, sdk.NewUint(1), lpAddress, 0)
		lpList = append(lpList, &lp)
	}
	return lpList
}

func GeneratePoolsAndLPs(keeper clpkeeper.Keeper, ctx sdk.Context, tokens []string) ([]types.Pool, []types.LiquidityProvider) {
	var poolList []types.Pool
	var lpList []types.LiquidityProvider
	for i := 0; i < len(tokens); i++ {
		externalToken := tokens[i]
		externalAsset := types.NewAsset(TrimFirstRune(externalToken))
		pool := types.NewPool(&externalAsset, sdk.NewUint(1000*uint64(i+1)), sdk.NewUint(100*uint64(i+1)), sdk.NewUint(1))
		err := keeper.SetPool(ctx, &pool)
		if err != nil {
			panic(err)
		}
		poolList = append(poolList, pool)
		lpAddress, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
		if err != nil {
			panic(err)
		}
		lp := types.NewLiquidityProvider(&externalAsset, sdk.NewUint(1), lpAddress, 0)
		keeper.SetLiquidityProvider(ctx, &lp)
		lpList = append(lpList, lp)
	}
	return poolList, lpList
}

func TrimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return strings.ToLower(s[i:])
}

func GenerateAddress2(key string) sdk.AccAddress {
	if key == "" {
		key = AddressKey1
	}
	var buffer bytes.Buffer
	buffer.WriteString(key)

	return genAddressInternal(buffer)
}

func GenerateAddress(key string) sdk.AccAddress {
	if key == "" {
		key = AddressKey1
	}
	var buffer bytes.Buffer
	buffer.WriteString(key)
	buffer.WriteString(strconv.Itoa(100))

	return genAddressInternal(buffer)
}

func genAddressInternal(buffer bytes.Buffer) sdk.AccAddress {
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

func GeneratePoolsFromFile(app *sifapp.SifchainApp, keeper clpkeeper.Keeper, ctx sdk.Context) []*types.Pool {
	var poolList types.PoolsRes

	file, err := filepath.Abs("test/pools_input.json")
	if err != nil {
		panic(err)
	}
	input, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(input, &poolList)
	if err != nil {
		panic(err)
	}
	// Set all pools
	for _, pool := range poolList.Pools {
		err := keeper.SetPool(ctx, pool)
		if err != nil {
			panic(err)
		}
		err = app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(
			sdk.NewCoin("rowan", sdk.NewIntFromBigInt(pool.NativeAssetBalance.BigInt())),
			sdk.NewCoin(pool.ExternalAsset.Symbol, sdk.NewIntFromBigInt(pool.ExternalAssetBalance.BigInt())),
		))
		if err != nil {
			panic(err)
		}
	}
	return poolList.Pools
}
