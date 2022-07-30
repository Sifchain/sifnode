package clp

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

func ReadPoolData() types.Pools {

	var pools types.Pools
	for _, asset := range Assets {
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
			RewardPeriodDistribute:        false,
			RewardPeriodMod:               100,
		}},
	})
	liquidityProtectionParam := keeper.GetLiquidityProtectionParams(ctx)
	liquidityProtectionParam.MaxRowanLiquidityThreshold = sdk.ZeroUint()
	keeper.SetLiquidityProtectionParams(ctx, liquidityProtectionParam)
	keeper.SetProviderDistributionParams(ctx, &types.ProviderDistributionParams{
		DistributionPeriods: []*types.ProviderDistributionPeriod{
			{
				DistributionPeriodBlockRate:  sdk.MustNewDecFromStr("0.0000034"),
				DistributionPeriodStartBlock: 1,
				DistributionPeriodEndBlock:   10000000,
				DistributionPeriodMod:        100,
			},
		},
	})

}

var Assets = []string{"rowan",
	"cusdt",
	"cusdc",
	"ccro",
	"cwbtc",
	"ceth",
	"cdai",
	"cyfi",
	"czrx",
	"cwscrt",
	"cwfil",
	"cuni",
	"cuma",
	"ctusd",
	"csxp",
	"csushi",
	"csusd",
	"csrm",
	"csnx",
	"csand",
	"crune",
	"creef",
	"cogn",
	"cocean",
	"cmana",
	"clrc",
	"clon",
	"clink",
	"ciotx",
	"cgrt",
	"cftm",
	"cesd",
	"cenj",
	"ccream",
	"ccomp",
	"ccocos",
	"cbond",
	"cbnt",
	"cbat",
	"cband",
	"cbal",
	"cant",
	"caave",
	"c1inch",
	"cleash",
	"cshib",
	"ctidal",
	"cpaid",
	"crndr",
	"cconv",
	"crfuel",
	"cakro",
	"cb20",
	"ctshp",
	"clina",
	"cdaofi",
	"ckeep",
	"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
	"ibc/6D717BFF5537D129035BAB39F593D638BA258A9F8D86FB7ECCEAB05B6950CC3E",
	"ibc/21CB41565FCA19AB6613EE06B0D56E588E0DC3E53FF94BA499BB9635794A1A35",
	"crly",
	"ibc/D87BC708A791246AA683D514C273736F07579CBD56C9CA79B7823F9A01C16270",
	"ibc/11DFDFADE34DCE439BA732EBA5CD8AA804A544BA1ECC0882856289FAF01FE53F",
	"ibc/B21954812E6E642ADC0B5ACB233E02A634BF137C572575BF80F7C0CC3DB2E74D",
	"ibc/2CC6F10253D563A7C238096BA63D060F7F356E37D5176E517034B8F730DB4AB6",
	"caxs",
	"cdfyn",
	"cdnxc",
	"cdon",
	"cern",
	"cfrax",
	"cfxs",
	"ckft",
	"cmatic",
	"cmetis",
	"cpols",
	"csaito",
	"ctoke",
	"czcn",
	"czcx",
	"cust",
	"cbtsg",
	"cquick",
	"cldo",
	"crail",
	"cpond",
	"cdino",
	"cufo",
	"ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA",
	"ibc/C5C8682EB9AA1313EF1B12C991ADCDA465B80C05733BFB2972E2005E01BCE459",
	"ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
	"cratom",
	"cfis",
	"ibc/17F5C77854734CFE1301E6067AA42CDF62DAF836E4467C635E6DB407853C6082",
	"ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D",
	"ibc/ACA7D0100794F39DF3FF0C5E31638B24737321C24F32C2C486A24C78DD8F2029",
	"ibc/7B8A3357032F3DB000ACFF3B2C9F8E77B932F21004FC93B5A8F77DE24161A573",
	"coh",
	"ibc/7876FB1D317D993F1F54185DF6E405C7FE070B71E3A53AE0CEA5A86AC878EB7A",
	"ccsms",
	"clgcy",
	"ibc/3313DFB885C0C0EBE85E307A529985AFF7CA82239D404329BDF294E357FBC73A",
	"cmc",
	"cinj",
	"cpush",
	"cgala",
	"cosqth",
	"cnewo",
	"cuos",
	"cxft",
	"ibc/F20C4E30E4202C11FE009D6D58B2FF212C99084CB6F767287A51A93EFD960086",
	"ibc/57BB0CFF9782730595988FD330AA41605B0628E11507BABC1207B830A23493B9",
	"ibc/345D30E8ED06B47FC538ED131D99D16126F07CD6F8B35DE96AAF4C1E445AF466",
	"ibc/E46B030074825C99488BC57FD2DA711B0650FEF2BD24B61C228BBE3BCD73E69E",
	"ibc/7B1E1EFA6808065DA759354B6F21433156F4BF5DF2CF96DCBBC91738683748AF",
	"ibc/84506C652F91EA3742B9E00C4240BB039466DBAC48BD12872D2C1BA3FCFCA31E",
	"ibc/B650115F83DF4CA83E406A0ABDCE0BC284DC0B382DEFF634321D256FA8AFE2B9",
	"ibc/C8D8DAB01D770335E61A09D5468FBD6AEA080794AB4B866CCAFB7AD85DD270FB",
	"ibc/41139CF1224ADAF97D1E1466815F50E9BBF19A8C311B8331A333035DA938A5CF",
	"ccudos",
	"ibc/902CFB7D533886C25315A4602EB1938968565215A770E6E2EBA0842FC14A62C9",
	"ibc/37AE3DD9177BAD68DC2E39BD43FB1E70C3BB719FAEBAB42F0D03132A2E23A7BF",
}
