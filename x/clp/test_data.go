package clp

import (
	"encoding/json"
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func ReadPoolData() types.Pools {
	var assets []string
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fp := filepath.Join(mydir, "pools.json")
	file, err := filepath.Abs(fp)
	if err != nil {
		panic(err)
	}
	input, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(input, &assets)
	if err != nil {
		panic(err)
	}
	var pools types.Pools
	for _, asset := range assets {
		pools = append(pools, types.Pool{
			ExternalAsset:                 &types.Asset{Symbol: asset},
			NativeAssetBalance:            sdk.Uint{},
			ExternalAssetBalance:          sdk.Uint{},
			PoolUnits:                     sdk.Uint{},
			SwapPriceNative:               nil,
			SwapPriceExternal:             nil,
			RewardPeriodNativeDistributed: sdk.Uint{},
		})
	}
	return pools
}

func TestingONLY_CreateAccounts(keeper keeper.Keeper, ctx sdk.Context) {
	SetParams(keeper, ctx)
	rowanAmount := "1000"
	amount, ok := sdk.NewIntFromString(rowanAmount)
	if !ok {
		panic("Unable to generate rowan amount")
	}
	coin := sdk.NewCoins(sdk.NewCoin("rowan", amount), sdk.NewCoin("ceth", amount), sdk.NewCoin("cusdc", amount), sdk.NewCoin("cwbtc", amount))
	bigAmount, ok := sdk.NewIntFromString("9999999999999999990000000000000000000000000000000000")
	if !ok {
		panic("unable to get big rowan")
	}
	pools := ReadPoolData()
	poolMap := map[string]sdk.Uint{}

	for _, pool := range pools {
		poolMap[pool.ExternalAsset.Symbol] = sdk.ZeroUint()
	}
	for i := 0; i < 8000; i++ {
		address := sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i))))
		err := keeper.GetBankKeeper().MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("rowan", bigAmount)))
		if err != nil {
			panic(err)
		}
		err = keeper.GetBankKeeper().MintCoins(ctx, types.ModuleName, coin)
		if err != nil {
			panic(err)
		}
		err = keeper.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, coin)
		if err != nil {
			panic(err)
		}
		CreateLiquidityProvidersForAddress(address, pools, keeper, ctx, poolMap)
	}
	for externalAsset, amount := range poolMap {
		err := keeper.SetPool(ctx, &types.Pool{
			ExternalAsset:                 &types.Asset{Symbol: externalAsset},
			NativeAssetBalance:            sdk.NewUintFromString("999000000000000000000000000000000"),
			ExternalAssetBalance:          amount,
			PoolUnits:                     sdk.NewUintFromString("1000000000000000000000"),
			SwapPriceNative:               nil,
			SwapPriceExternal:             nil,
			RewardPeriodNativeDistributed: sdk.Uint{},
		})
		if err != nil {
			panic(err)
		}
	}
}

func CreateLiquidityProvidersForAddress(address sdk.AccAddress, pools types.Pools, keeper keeper.Keeper, ctx sdk.Context, poolMap map[string]sdk.Uint) {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(pools) - 1
	for i := 0; i < 10; i++ {
		index := rand.Intn(max-min+1) + min
		randomPool := pools[index]
		keeper.CreateLiquidityProvider(ctx, &types.Asset{Symbol: randomPool.ExternalAsset.Symbol}, sdk.NewUintFromString("1000000000000000000000"), address)
		poolMap[randomPool.ExternalAsset.Symbol] = poolMap[randomPool.ExternalAsset.Symbol].Add(sdk.NewUintFromString("1000000000000000000000"))
	}
}

func SetParams(keeper keeper.Keeper, ctx sdk.Context) {
	keeper.SetPmtpRateParams(ctx, types.PmtpRateParams{
		PmtpPeriodBlockRate:    sdk.OneDec(),
		PmtpCurrentRunningRate: sdk.OneDec(),
	})
	keeper.SetParams(ctx, types.Params{
		MinCreatePoolThreshold: 100,
	})
	keeper.SetPmtpParams(ctx, &types.PmtpParams{
		PmtpPeriodGovernanceRate: sdk.OneDec(),
		PmtpPeriodEpochLength:    1,
		PmtpPeriodStartBlock:     1,
		PmtpPeriodEndBlock:       2,
	})
	allocation := sdk.NewUintFromString("100000000000000000000000000000")
	defautMultiplier := sdk.NewDec(1)
	keeper.SetRewardParams(ctx, &types.RewardParams{
		LiquidityRemovalLockPeriod:   0,
		LiquidityRemovalCancelPeriod: 2,
		RewardPeriodStartTime:        "",
		RewardPeriods: []*types.RewardPeriod{{
			RewardPeriodId:                "1",
			RewardPeriodStartBlock:        1,
			RewardPeriodEndBlock:          10000000,
			RewardPeriodAllocation:        &allocation,
			RewardPeriodPoolMultipliers:   nil,
			RewardPeriodDefaultMultiplier: &defautMultiplier,
			RewardPeriodDistribute:        true,
			RewardPeriodMod:               1,
		}},
	})
	liquidityProtectionParam := keeper.GetLiquidityProtectionParams(ctx)
	liquidityProtectionParam.MaxRowanLiquidityThreshold = sdk.ZeroUint()
	keeper.SetLiquidityProtectionParams(ctx, liquidityProtectionParam)
	keeper.SetProviderDistributionParams(ctx, &types.ProviderDistributionParams{
		DistributionPeriods: nil,
	})

}
