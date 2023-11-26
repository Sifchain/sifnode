package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestKeeper_CheckBalances(t *testing.T) {
	nativeAmount, _ := sdk.NewIntFromString("999999000000000000000000000")
	externalAmount, _ := sdk.NewIntFromString("500000000000000000000000")
	const address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"

	ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
		balances := []banktypes.Balance{
			{
				Address: address,
				Coins: sdk.Coins{
					sdk.NewCoin("catk", externalAmount),
					sdk.NewCoin("cbtk", externalAmount),
					sdk.NewCoin("cdash", externalAmount),
					sdk.NewCoin("ceth", externalAmount),
					sdk.NewCoin("clink", externalAmount),
					sdk.NewCoin("rowan", nativeAmount),
				},
			},
		}
		gs := banktypes.DefaultGenesisState()
		gs.Balances = append(gs.Balances, balances...)
		bz, _ := app.AppCodec().MarshalJSON(gs)

		genesisState["bank"] = bz

		return genesisState
	})

	accAddress, _ := sdk.AccAddressFromBech32(address)

	balances := app.BankKeeper.GetAllBalances(ctx, accAddress)
	require.Contains(t, balances, sdk.Coin{
		Denom: "catk", Amount: externalAmount,
	})
	require.Contains(t, balances, sdk.Coin{
		Denom: "ceth", Amount: externalAmount,
	})
	require.Contains(t, balances, sdk.Coin{
		Denom: "clink", Amount: externalAmount,
	})
}

func TestKeeper_SwapOne(t *testing.T) {
	testcases := []struct {
		name                 string
		nativeAssetBalance   sdk.Uint
		externalAssetBalance sdk.Uint
		nativeCustody        sdk.Uint
		externalCustody      sdk.Uint
		nativeLiabilities    sdk.Uint
		externalLiabilities,
		sentAmount sdk.Uint
		fromAsset                    types.Asset
		toAsset                      types.Asset
		pmtpCurrentRunningRate       sdk.Dec
		swapFeeRate                  sdk.Dec
		errString                    error
		expectedSwapResult           sdk.Uint
		expectedLiquidityFee         sdk.Uint
		expectedPriceImpact          sdk.Uint
		expectedExternalAssetBalance sdk.Uint
		expectedNativeAssetBalance   sdk.Uint
	}{
		{
			name:                         "real world numbers",
			nativeAssetBalance:           sdk.NewUint(10000000),
			externalAssetBalance:         sdk.NewUint(8770000),
			nativeCustody:                sdk.ZeroUint(),
			externalCustody:              sdk.ZeroUint(),
			nativeLiabilities:            sdk.ZeroUint(),
			externalLiabilities:          sdk.ZeroUint(),
			sentAmount:                   sdk.NewUint(50000),
			fromAsset:                    types.GetSettlementAsset(),
			toAsset:                      types.NewAsset("eth"),
			pmtpCurrentRunningRate:       sdk.NewDec(0),
			swapFeeRate:                  sdk.NewDecWithPrec(3, 3),
			expectedSwapResult:           sdk.NewUint(43501),
			expectedLiquidityFee:         sdk.NewUint(130),
			expectedPriceImpact:          sdk.ZeroUint(),
			expectedExternalAssetBalance: sdk.NewUint(8726499),
			expectedNativeAssetBalance:   sdk.NewUint(10050000),
		},
		{
			name:                         "big numbers",
			nativeAssetBalance:           sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance:         sdk.NewUint(2674623482959),
			nativeCustody:                sdk.ZeroUint(),
			externalCustody:              sdk.ZeroUint(),
			nativeLiabilities:            sdk.ZeroUint(),
			externalLiabilities:          sdk.ZeroUint(),
			sentAmount:                   sdk.NewUint(200000000),
			toAsset:                      types.GetSettlementAsset(),
			fromAsset:                    types.NewAsset("cusdt"),
			pmtpCurrentRunningRate:       sdk.NewDec(0),
			swapFeeRate:                  sdk.NewDecWithPrec(3, 3),
			expectedSwapResult:           sdk.NewUintFromString("11704434254784015637542"),
			expectedLiquidityFee:         sdk.NewUintFromString("35218959643281892590"),
			expectedPriceImpact:          sdk.ZeroUint(),
			expectedExternalAssetBalance: sdk.NewUint(2674823482959),
			expectedNativeAssetBalance:   sdk.NewUintFromString("156995796064471436224541544"),
		},
		{
			name:                         "margin enabled",
			nativeAssetBalance:           sdk.NewUint(10000000),
			externalAssetBalance:         sdk.NewUint(8770000),
			nativeCustody:                sdk.NewUint(10000),
			externalCustody:              sdk.ZeroUint(),
			nativeLiabilities:            sdk.ZeroUint(),
			externalLiabilities:          sdk.NewUint(10000),
			sentAmount:                   sdk.NewUint(50000),
			fromAsset:                    types.GetSettlementAsset(),
			toAsset:                      types.NewAsset("eth"),
			pmtpCurrentRunningRate:       sdk.NewDec(0),
			swapFeeRate:                  sdk.NewDecWithPrec(3, 3),
			expectedSwapResult:           sdk.NewUint(43550),
			expectedLiquidityFee:         sdk.NewUint(131),
			expectedPriceImpact:          sdk.ZeroUint(),
			expectedExternalAssetBalance: sdk.NewUint(8726450),
			expectedNativeAssetBalance:   sdk.NewUint(10050000),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			//ctx, app := test.CreateTestAppClp(false)
			poolUnits := sdk.NewUint(2000) //don't care
			pool := types.NewPool(&tc.toAsset, tc.nativeAssetBalance, tc.externalAssetBalance, poolUnits)
			pool.NativeCustody = tc.nativeCustody
			pool.ExternalCustody = tc.externalCustody
			pool.NativeLiabilities = tc.nativeLiabilities
			pool.ExternalLiabilities = tc.externalLiabilities

			swapResult, liquidityFee, priceImpact, pool, err := clpkeeper.SwapOne(tc.fromAsset, tc.sentAmount, tc.toAsset, pool, tc.pmtpCurrentRunningRate, tc.swapFeeRate)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}

			assert.NoError(t, err)
			require.Equal(t, tc.expectedSwapResult.String(), swapResult.String())
			require.Equal(t, tc.expectedLiquidityFee.String(), liquidityFee.String())
			require.Equal(t, tc.expectedPriceImpact.String(), priceImpact.String())
			require.Equal(t, tc.expectedExternalAssetBalance.String(), pool.ExternalAssetBalance.String())
			require.Equal(t, tc.expectedNativeAssetBalance.String(), pool.NativeAssetBalance.String())
		})
	}

}

func TestKeeper_ExtractValuesFromPool(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, _ := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	X, Y, toRowan, from := pool.ExtractValues(asset)

	assert.Equal(t, X, sdk.NewUint(998))
	assert.Equal(t, Y, sdk.NewUint(998))
	assert.Equal(t, toRowan, false)
	assert.Equal(t, types.GetSettlementAsset(), from)
}

func TestKeeper_GetSwapFee(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, _ := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	swapFeeRate := sdk.NewDecWithPrec(3, 3)
	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, *pool, sdk.OneDec(), swapFeeRate)
	assert.Equal(t, "1", swapResult.String())
}

func TestKeeper_GetSwapFee_PmtpParams(t *testing.T) {
	pool := types.Pool{
		NativeAssetBalance:   sdk.NewUint(10),
		ExternalAssetBalance: sdk.NewUint(100),
		NativeLiabilities:    sdk.ZeroUint(),
		NativeCustody:        sdk.ZeroUint(),
		ExternalLiabilities:  sdk.ZeroUint(),
		ExternalCustody:      sdk.ZeroUint(),
	}
	asset := types.Asset{}

	swapFeeRate := sdk.NewDecWithPrec(3, 3)

	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, pool, sdk.NewDec(100), swapFeeRate)

	require.Equal(t, swapResult, sdk.ZeroUint())
}

func TestKeeper_CalculateAssetsForLP(t *testing.T) {
	_, app, ctx := createTestInput()
	keeper := app.ClpKeeper
	tokens := []string{"cada", "cbch", "cbnb", "cbtc", "ceos", "ceth", "ctrx", "cusdt"}
	pools, lpList := test.GeneratePoolsAndLPs(keeper, ctx, tokens)
	native, external, _, _ := clpkeeper.CalculateAllAssetsForLP(pools[0], lpList[0])
	assert.Equal(t, "100", external.String())
	assert.Equal(t, "1000", native.String())
}

func TestKeeper_CalculateWithdrawal(t *testing.T) {
	testcases := []struct {
		name                 string
		poolUnits            sdk.Uint
		nativeAssetBalance   string
		externalAssetBalance string
		lpUnits              string
		wBasisPoints         string
		asymmetry            sdk.Int
		panicErr             string
	}{
		{
			name:                 "fail to convert nativeAssetBalance to Dec",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			externalAssetBalance: "1",
			lpUnits:              "1",
			wBasisPoints:         "1",
			asymmetry:            sdk.NewInt(1),
			panicErr:             "fail to convert 100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal '100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000' out of range; bitLen: got 545, max 315",
		},
		{
			name:                 "fail to convert externalAssetBalance to Dec",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "1",
			externalAssetBalance: "100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			lpUnits:              "1",
			wBasisPoints:         "1",
			asymmetry:            sdk.NewInt(1),
			panicErr:             "fail to convert 100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal '100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000' out of range; bitLen: got 545, max 315",
		},
		{
			name:                 "fail to convert lpUnits to Dec",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "1",
			externalAssetBalance: "1",
			lpUnits:              "100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			wBasisPoints:         "1",
			asymmetry:            sdk.NewInt(1),
			panicErr:             "fail to convert 100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal '100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000' out of range; bitLen: got 545, max 315",
		},
		{
			name:                 "fail to convert wBasisPoints to Dec",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "1",
			externalAssetBalance: "1",
			lpUnits:              "1",
			wBasisPoints:         "100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			asymmetry:            sdk.NewInt(1),
			panicErr:             "fail to convert 100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal '100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000' out of range; bitLen: got 545, max 315",
		},
		//The panic is from sdk.NewUintFromString and not CalculateWithdrawal
		//{
		//	name:                 "fail to convert asymmetry to INT",
		//	poolUnits:            sdk.NewUint(1),
		//	nativeAssetBalance:   "1",
		//	externalAssetBalance: "1",
		//	lpUnits:              "1",
		//	wBasisPoints:         "1",
		//	asymmetry:            sdk.Int(sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
		//	panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal '10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000' out of range; bitLen: got 293, max 256",
		//},
		{
			name:                 "asymmetric value negative",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "1",
			externalAssetBalance: "1",
			lpUnits:              "1",
			wBasisPoints:         "1",
			asymmetry:            sdk.NewInt(-1000),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panicErr != "" {
				require.PanicsWithError(t, tc.panicErr, func() {
					clpkeeper.CalculateWithdrawal(tc.poolUnits, tc.nativeAssetBalance, tc.externalAssetBalance, tc.lpUnits, tc.wBasisPoints, tc.asymmetry)
				})
				return
			}

			w, x, y, z := clpkeeper.CalculateWithdrawal(tc.poolUnits, tc.nativeAssetBalance, tc.externalAssetBalance, tc.lpUnits, tc.wBasisPoints, tc.asymmetry)

			require.NotNil(t, w)
			require.NotNil(t, x)
			require.NotNil(t, y)
			require.NotNil(t, z)
		})
	}
}

func TestKeeper_CalcSwapResult(t *testing.T) {
	testcases := []struct {
		name                    string
		toRowan                 bool
		X, x, Y, y, expectedFee sdk.Uint
		pmtpCurrentRunningRate  sdk.Dec
		swapFeeRate             sdk.Dec
		err                     error
		errString               error
	}{
		{
			name:                   "one side of pool empty",
			toRowan:                true,
			X:                      sdk.NewUint(0),
			x:                      sdk.NewUint(12),
			Y:                      sdk.NewUint(12),
			y:                      sdk.NewUint(0),
			expectedFee:            sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(2),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "swap amount zero",
			toRowan:                true,
			X:                      sdk.NewUint(117),
			x:                      sdk.NewUint(0),
			Y:                      sdk.NewUint(12),
			y:                      sdk.NewUint(0),
			expectedFee:            sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(2),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "real world amounts, buy rowan",
			toRowan:                true,
			X:                      sdk.NewUint(1999800619938006200),
			x:                      sdk.NewUint(200000000000000),
			Y:                      sdk.NewUint(2000200000000000000),
			y:                      sdk.NewUint(66473292728673),
			expectedFee:            sdk.NewUint(200019938000),
			pmtpCurrentRunningRate: sdk.NewDec(2),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "real world amounts, sell rowan",
			toRowan:                false,
			X:                      sdk.NewUint(1999800619938006200),
			x:                      sdk.NewUint(200000000000000),
			Y:                      sdk.NewUint(2000200000000000000),
			y:                      sdk.NewUint(598259634558057),
			expectedFee:            sdk.NewUint(1800179442000),
			pmtpCurrentRunningRate: sdk.NewDec(2),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "big numbers",
			toRowan:                true,
			X:                      sdk.NewUintFromString("20300000000000000000000000000000000000000000000000000000000000000000000000"),
			x:                      sdk.NewUintFromString("10000000000000000658000000000000000000000000000000000000000000000000000000"),
			Y:                      sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000021344"),
			y:                      sdk.NewUintFromString("1096809680968096858032537841242869111592632578510191129990575247754487593"),
			expectedFee:            sdk.NewUintFromString("3300330033003300475524186081974530927560579473952430682017779080504977"),
			pmtpCurrentRunningRate: sdk.NewDec(2),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			y, fee := clpkeeper.CalcSwapResult(tc.toRowan, tc.X, tc.x, tc.Y, tc.pmtpCurrentRunningRate, tc.swapFeeRate)

			require.Equal(t, tc.y.String(), y.String()) // compare strings so that the expected amounts can be read from the failure message
			require.Equal(t, tc.expectedFee.String(), fee.String())
		})
	}
}

func getFirstArg(a *big.Int, b bool) *big.Int {
	return a
}

func TestKeeper_CalcDenomChangeMultiplier(t *testing.T) {
	testcases := []struct {
		name      string
		decimalsX uint8
		decimalsY uint8
		expected  big.Rat
	}{
		{
			name:      "zero values",
			decimalsX: 0,
			decimalsY: 0,
			expected:  *big.NewRat(1, 1),
		},
		{
			name:      "equal values",
			decimalsX: 5,
			decimalsY: 5,
			expected:  *big.NewRat(1, 1),
		},
		{
			name:      "zero X",
			decimalsX: 0,
			decimalsY: 2,
			expected:  *big.NewRat(1, 100),
		},
		{
			name:      "zero Y",
			decimalsX: 2,
			decimalsY: 0,
			expected:  *big.NewRat(100, 1),
		},
		{
			name:      "small numbers",
			decimalsX: 18,
			decimalsY: 14,
			expected:  *big.NewRat(10000, 1),
		},
		{
			name:      "small numbers",
			decimalsX: 14,
			decimalsY: 18,
			expected:  *big.NewRat(1, 10000),
		},
		{
			name:      "big X, small Y",
			decimalsX: 255,
			decimalsY: 0,
			expected:  *big.NewRat(1, 1).SetInt(big.NewInt(1).Exp(big.NewInt(10), big.NewInt(255), nil)),
		},
		{
			name:      "small X, big Y",
			decimalsX: 0,
			decimalsY: 255,
			expected:  *big.NewRat(1, 1).SetFrac(big.NewInt(1), big.NewInt(1).Exp(big.NewInt(10), big.NewInt(255), nil)),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			y := clpkeeper.CalcDenomChangeMultiplier(tc.decimalsX, tc.decimalsY)

			require.Equal(t, tc.expected.String(), y.String()) // compare strings so that the expected amounts can be read from the failure message
		})
	}
}

// nolint
func TestKeeper_CalcSpotPriceX(t *testing.T) {

	testcases := []struct {
		name                   string
		X                      sdk.Uint
		Y                      sdk.Uint
		decimalsX              uint8
		decimalsY              uint8
		pmtpCurrentRunningRate sdk.Dec
		isXNative              bool
		expected               sdk.Dec
		errString              error
	}{
		{
			name:                   "fail when X = 0",
			X:                      sdk.ZeroUint(),
			Y:                      sdk.OneUint(),
			decimalsX:              10,
			decimalsY:              80,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			isXNative:              true,
			errString:              errors.New("amount is invalid"),
		},
		{
			name:                   "success when Y = 0",
			X:                      sdk.OneUint(),
			Y:                      sdk.ZeroUint(),
			decimalsX:              10,
			decimalsY:              80,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			isXNative:              true,
			expected:               sdk.NewDec(0),
		},
		{
			name:                   "success small values",
			X:                      sdk.OneUint(),
			Y:                      sdk.OneUint(),
			decimalsX:              18,
			decimalsY:              18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			isXNative:              true,
			expected:               sdk.NewDec(1),
		},
		{
			name:                   "success mid values",
			X:                      sdk.NewUint(12345678),
			Y:                      sdk.NewUint(67890123),
			decimalsX:              18,
			decimalsY:              18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			isXNative:              true,
			expected:               sdk.MustNewDecFromStr("5.499100413926233941"),
		},
		{
			name:                   "success mid values with PMTP",
			X:                      sdk.NewUint(12345678),
			Y:                      sdk.NewUint(67890123),
			decimalsX:              18,
			decimalsY:              18,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			isXNative:              true,
			expected:               sdk.MustNewDecFromStr("10.998200827852467883"),
		},
		{
			name:                   "success mid values with PMTP and decimals",
			X:                      sdk.NewUint(12345678),
			Y:                      sdk.NewUint(67890123),
			decimalsX:              16,
			decimalsY:              18,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			isXNative:              true,
			expected:               sdk.MustNewDecFromStr("0.109982008278524678"),
		},
		{
			name:                   "success big numbers",
			X:                      sdk.OneUint(),
			Y:                      sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"), //2**200
			decimalsX:              18,
			decimalsY:              18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			isXNative:              true,
			expected:               sdk.NewDecFromBigIntWithPrec(getFirstArg(big.NewInt(1).SetString("1606938044258990275541962092341162602522202993782792835301376000000000000000000", 10)), 18),
		},
		{
			name:                   "failure big decimals",
			X:                      sdk.NewUint(100),
			Y:                      sdk.NewUint(100),
			decimalsX:              255,
			decimalsY:              0,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			isXNative:              true,
			errString:              errors.New("decimal out of range; bitLen: got 907, max 315"),
		},
		{
			name:                   "success big decimals, small answer",
			X:                      sdk.NewUint(100),
			Y:                      sdk.NewUint(100),
			decimalsX:              0,
			decimalsY:              255,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			isXNative:              true,
			expected:               sdk.MustNewDecFromStr("0.000000000000000000"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			price, err := clpkeeper.CalcSpotPriceX(tc.X, tc.Y, tc.decimalsX, tc.decimalsY, tc.pmtpCurrentRunningRate, tc.isXNative)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, price)
		})
	}
}

func TestKeeper_CalcSpotPriceNative(t *testing.T) {

	testcases := []struct {
		name                   string
		nativeAssetBalance     sdk.Uint
		externalAssetBalance   sdk.Uint
		decimalsExternal       uint8
		pmtpCurrentRunningRate sdk.Dec
		expected               sdk.Dec
		errString              error
	}{
		{
			name:                   "fail when native balance = 0",
			nativeAssetBalance:     sdk.ZeroUint(),
			externalAssetBalance:   sdk.OneUint(),
			decimalsExternal:       80,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			errString:              errors.New("amount is invalid"),
		},
		{
			name:                   "success when external balance = 0",
			nativeAssetBalance:     sdk.OneUint(),
			externalAssetBalance:   sdk.ZeroUint(),
			decimalsExternal:       10,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			expected:               sdk.NewDec(0),
		},
		{
			name:                   "success small values",
			nativeAssetBalance:     sdk.OneUint(),
			externalAssetBalance:   sdk.OneUint(),
			decimalsExternal:       18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.NewDec(1),
		},
		{
			name:                   "success mid values",
			nativeAssetBalance:     sdk.NewUint(12345678),
			externalAssetBalance:   sdk.NewUint(67890123),
			decimalsExternal:       18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.MustNewDecFromStr("5.499100413926233941"),
		},
		{
			name:                   "success mid values with PMTP",
			nativeAssetBalance:     sdk.NewUint(12345678),
			externalAssetBalance:   sdk.NewUint(67890123),
			decimalsExternal:       18,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			expected:               sdk.MustNewDecFromStr("10.998200827852467883"),
		},
		{
			name:                   "success mid values with PMTP and decimals",
			nativeAssetBalance:     sdk.NewUint(12345678),
			externalAssetBalance:   sdk.NewUint(67890123),
			decimalsExternal:       16,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			expected:               sdk.MustNewDecFromStr("1099.820082785246788390"),
		},
		{
			name:                   "success big numbers",
			nativeAssetBalance:     sdk.OneUint(),
			externalAssetBalance:   sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"), //2**200
			decimalsExternal:       18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.NewDecFromBigIntWithPrec(getFirstArg(big.NewInt(1).SetString("1606938044258990275541962092341162602522202993782792835301376000000000000000000", 10)), 18),
		},
		{
			name:                   "success big decimals",
			nativeAssetBalance:     sdk.NewUint(100),
			externalAssetBalance:   sdk.NewUint(100),
			decimalsExternal:       255,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.MustNewDecFromStr("0.000000000000000000"),
		},
		{
			name:                   "success small decimals",
			nativeAssetBalance:     sdk.NewUint(100),
			externalAssetBalance:   sdk.NewUint(100),
			decimalsExternal:       0,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.MustNewDecFromStr("1000000000000000000.000000000000000000"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			pool := types.Pool{
				NativeAssetBalance:   tc.nativeAssetBalance,
				ExternalAssetBalance: tc.externalAssetBalance,
				NativeLiabilities:    sdk.ZeroUint(),
				NativeCustody:        sdk.ZeroUint(),
				ExternalLiabilities:  sdk.ZeroUint(),
				ExternalCustody:      sdk.ZeroUint(),
			}

			price, err := clpkeeper.CalcSpotPriceNative(&pool, tc.decimalsExternal, tc.pmtpCurrentRunningRate)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, price)
		})
	}
}

func TestKeeper_CalcSpotPriceExternal(t *testing.T) {

	testcases := []struct {
		name                   string
		nativeAssetBalance     sdk.Uint
		externalAssetBalance   sdk.Uint
		decimalsExternal       uint8
		pmtpCurrentRunningRate sdk.Dec
		expected               sdk.Dec
		errString              error
	}{
		{
			name:                   "success when native balance = 0",
			nativeAssetBalance:     sdk.ZeroUint(),
			externalAssetBalance:   sdk.OneUint(),
			decimalsExternal:       80,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			expected:               sdk.NewDec(0),
		},
		{
			name:                   "fail when external balance = 0",
			nativeAssetBalance:     sdk.OneUint(),
			externalAssetBalance:   sdk.ZeroUint(),
			decimalsExternal:       10,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			errString:              errors.New("amount is invalid"),
		},
		{
			name:                   "success small values",
			nativeAssetBalance:     sdk.OneUint(),
			externalAssetBalance:   sdk.OneUint(),
			decimalsExternal:       18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.NewDec(1),
		},
		{
			name:                   "success mid values",
			nativeAssetBalance:     sdk.NewUint(12345678),
			externalAssetBalance:   sdk.NewUint(67890123),
			decimalsExternal:       18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.MustNewDecFromStr("0.181847925065624052"),
		},
		{
			name:                   "success mid values with PMTP",
			nativeAssetBalance:     sdk.NewUint(12345678),
			externalAssetBalance:   sdk.NewUint(67890123),
			decimalsExternal:       18,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			expected:               sdk.MustNewDecFromStr("0.090923962532812026"),
		},
		{
			name:                   "success mid values with PMTP and decimals",
			nativeAssetBalance:     sdk.NewUint(12345678),
			externalAssetBalance:   sdk.NewUint(67890123),
			decimalsExternal:       16,
			pmtpCurrentRunningRate: sdk.NewDec(1),
			expected:               sdk.MustNewDecFromStr("0.000909239625328120"),
		},
		{
			name:                   "success big numbers",
			nativeAssetBalance:     sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"), //2**200
			externalAssetBalance:   sdk.OneUint(),
			decimalsExternal:       18,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.NewDecFromBigIntWithPrec(getFirstArg(big.NewInt(1).SetString("1606938044258990275541962092341162602522202993782792835301376000000000000000000", 10)), 18),
		},
		{
			name:                   "failure big decimals",
			nativeAssetBalance:     sdk.NewUint(100),
			externalAssetBalance:   sdk.NewUint(100),
			decimalsExternal:       255,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			errString:              errors.New("decimal out of range; bitLen: got 848, max 315"),
		},
		{
			name:                   "success small decimals",
			nativeAssetBalance:     sdk.NewUint(100),
			externalAssetBalance:   sdk.NewUint(100),
			decimalsExternal:       0,
			pmtpCurrentRunningRate: sdk.NewDec(0),
			expected:               sdk.MustNewDecFromStr("0.000000000000000001"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			pool := types.Pool{
				NativeAssetBalance:   tc.nativeAssetBalance,
				ExternalAssetBalance: tc.externalAssetBalance,
				NativeLiabilities:    sdk.ZeroUint(),
				NativeCustody:        sdk.ZeroUint(),
				ExternalLiabilities:  sdk.ZeroUint(),
				ExternalCustody:      sdk.ZeroUint(),
			}

			price, err := clpkeeper.CalcSpotPriceExternal(&pool, tc.decimalsExternal, tc.pmtpCurrentRunningRate)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, price)
		})
	}
}

func TestKeeper_CalcRowanSpotPrice(t *testing.T) {
	testcases := []struct {
		name                          string
		rowanBalance, externalBalance sdk.Uint
		pmtpCurrentRunningRate        sdk.Dec
		expectedPrice                 sdk.Dec
		expectedError                 error
	}{
		{
			name:                   "success simple",
			rowanBalance:           sdk.NewUint(1),
			externalBalance:        sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(1),
			expectedPrice:          sdk.MustNewDecFromStr("2"),
		},
		{
			name:                   "success small",
			rowanBalance:           sdk.NewUint(1000000000123),
			externalBalance:        sdk.NewUint(20000000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.4"),
			expectedPrice:          sdk.MustNewDecFromStr("0.000047999999994096"),
		},

		{
			name:                   "success",
			rowanBalance:           sdk.NewUint(1000),
			externalBalance:        sdk.NewUint(2000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.4"),
			expectedPrice:          sdk.MustNewDecFromStr("4.8"),
		},
		{
			name:                   "fail - rowan balance zero",
			rowanBalance:           sdk.NewUint(0),
			externalBalance:        sdk.NewUint(2000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.4"),
			expectedError:          errors.New("amount is invalid"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			pool := types.Pool{
				NativeAssetBalance:   tc.rowanBalance,
				ExternalAssetBalance: tc.externalBalance,
				NativeLiabilities:    sdk.ZeroUint(),
				NativeCustody:        sdk.ZeroUint(),
				ExternalLiabilities:  sdk.ZeroUint(),
				ExternalCustody:      sdk.ZeroUint(),
			}

			calcPrice, err := clpkeeper.CalcRowanSpotPrice(&pool, tc.pmtpCurrentRunningRate)
			if tc.expectedError != nil {
				require.EqualError(t, tc.expectedError, err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedPrice, calcPrice)
		})
	}
}

func TestKeeper_CalcRowanValue(t *testing.T) {
	testcases := []struct {
		name          string
		rowanAmount   sdk.Uint
		price         sdk.Dec
		expectedValue sdk.Uint
	}{
		{
			name:          "success simple",
			rowanAmount:   sdk.NewUint(100),
			price:         sdk.NewDecWithPrec(232, 2),
			expectedValue: sdk.NewUint(232),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rowanValue := clpkeeper.CalcRowanValue(tc.rowanAmount, tc.price)
			require.Equal(t, tc.expectedValue.String(), rowanValue.String())
		})
	}
}

// // Used only to generate expected results for TestKeeper_CalculateExternalSwapAmountAsymmetricRat
// // Useful to keep around if more test cases are needed in future
// func TestKeeper_GenerateCalculateExternalSwapAmountAsymmetricRatTestCases(t *testing.T) {
// 	testcases := []struct {
// 		Y, X, y, x, f, r float64
// 		expectedValue    float64
// 	}{
// 		{
// 			Y: 100000,
// 			X: 100000,
// 			y: 2000,
// 			x: 8000,
// 			f: 0.003,
// 			r: 0.01,
// 		},
// 		{
// 			Y: 3456789887,
// 			X: 1244516357,
// 			y: 2000,
// 			x: 99887776,
// 			f: 0.003,
// 			r: 0.01,
// 		},
// 		{
// 			Y: 157007500498726220240179086,
// 			X: 2674623482959,
// 			y: 0,
// 			x: 200000000,
// 			f: 0.003,
// 			r: 0.01,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		Y := tc.Y
// 		X := tc.X
// 		y := tc.y
// 		x := tc.x
// 		f := tc.f
// 		r := tc.r
// 		expected := math.Abs((math.Sqrt(Y*(-1*(x+X))*(-1*f*f*x*Y-f*f*X*Y-2*f*r*x*Y+4*f*r*X*y+2*f*r*X*Y+4*f*X*y+4*f*X*Y-r*r*x*Y-r*r*X*Y-4*r*X*y-4*r*X*Y-4*X*y-4*X*Y)) + f*x*Y + f*X*Y + r*x*Y - 2*r*X*y - r*X*Y - 2*X*y - 2*X*Y) / (2 * (r + 1) * (y + Y)))
// 		fmt.Println(expected)
// 	}

// }
func TestKeeper_CalculateExternalSwapAmountAsymmetricRat(t *testing.T) {
	testcases := []struct {
		name             string
		Y, X, y, x, f, r *big.Rat
		expectedValue    sdk.Dec
	}{
		{
			name:          "test1",
			Y:             big.NewRat(100000, 1),
			X:             big.NewRat(100000, 1),
			y:             big.NewRat(2000, 1),
			x:             big.NewRat(8000, 1),
			f:             big.NewRat(3, 1000), // 0.003
			r:             big.NewRat(1, 100),  // 0.01
			expectedValue: sdk.MustNewDecFromStr("2918.476067753834206950"),
		},
		{
			name:          "test2",
			Y:             big.NewRat(3456789887, 1),
			X:             big.NewRat(1244516357, 1),
			y:             big.NewRat(2000, 1),
			x:             big.NewRat(99887776, 1),
			f:             big.NewRat(3, 1000), // 0.003
			r:             big.NewRat(1, 100),  // 0.01
			expectedValue: sdk.MustNewDecFromStr("49309453.001511112834211406"),
		},
		{
			name:          "test3",
			Y:             MustRatFromString("157007500498726220240179086"),
			X:             big.NewRat(2674623482959, 1),
			y:             big.NewRat(0, 1),
			x:             big.NewRat(200000000, 1),
			f:             big.NewRat(3, 1000), // 0.003
			r:             big.NewRat(1, 100),  // 0.01
			expectedValue: sdk.MustNewDecFromStr("100645875.768947133021515445"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := clpkeeper.CalculateExternalSwapAmountAsymmetricRat(tc.Y, tc.X, tc.y, tc.x, tc.f, tc.r)
			got, _ := clpkeeper.RatToDec(&res)

			require.Equal(t, tc.expectedValue.String(), got.String())
		})
	}

}

// // Used only to generate expected results for TestKeeper_CalculateNativeSwapAmountAsymmetricRat
// // Useful to keep around if more test cases are needed in future
// func TestKeeper_GenerateCalculateNativeSwapAmountAsymmetricRatTestCases(t *testing.T) {
// 	testcases := []struct {
// 		Y, X, y, x, f, r float64
// 		expectedValue    float64
// 	}{
// 		{
// 			Y: 100000,
// 			X: 100000,
// 			y: 8000,
// 			x: 2000,
// 			f: 0.003,
// 			r: 0.01,
// 		},
// 		{
// 			Y: 3456789887,
// 			X: 1244516357,
// 			y: 99887776,
// 			x: 2000,
// 			f: 0.003,
// 			r: 0.01,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		Y := tc.Y
// 		X := tc.X
// 		y := tc.y
// 		x := tc.x
// 		f := tc.f
// 		r := tc.r
// 		expected := math.Abs((math.Sqrt(math.Pow((-1*f*r*X*y-f*r*X*Y-f*X*y-f*X*Y+r*X*y+r*X*Y+2*x*Y+2*X*Y), 2)-4*(x+X)*(x*Y*Y-X*y*Y)) + f*r*X*y + f*r*X*Y + f*X*y + f*X*Y - r*X*y - r*X*Y - 2*x*Y - 2*X*Y) / (2 * (x + X)))
// 		fmt.Println(expected)
// 	}

// }
func TestKeeper_CalculateNativeSwapAmountAsymmetricRat(t *testing.T) {
	testcases := []struct {
		name             string
		Y, X, y, x, f, r *big.Rat
		expectedValue    sdk.Dec
	}{
		{
			name:          "test1",
			Y:             big.NewRat(100000, 1),
			X:             big.NewRat(100000, 1),
			y:             big.NewRat(8000, 1),
			x:             big.NewRat(2000, 1),
			f:             big.NewRat(3, 1000), // 0.003
			r:             big.NewRat(1, 100),  // 0.01
			expectedValue: sdk.MustNewDecFromStr("2888.791254901960784313"),
		},
		{
			name:          "test2",
			Y:             big.NewRat(3456789887, 1),
			X:             big.NewRat(1244516357, 1),
			y:             big.NewRat(99887776, 1),
			x:             big.NewRat(2000, 1),
			f:             big.NewRat(3, 1000), // 0.003
			r:             big.NewRat(1, 100),  // 0.01
			expectedValue: sdk.MustNewDecFromStr("49410724.289235769454274911"),
		},
		{
			//NOTE: cannot be confirmed with the float64 model above since that runs out of precision.
			// However the expectedValue is about half the value of y, which is as expected
			name:          "test3",
			Y:             MustRatFromString("157007500498726220240179086"),
			X:             big.NewRat(2674623482959, 1),
			y:             big.NewRat(200000000, 1),
			x:             big.NewRat(0, 1),
			f:             big.NewRat(3, 1000), // 0.003
			r:             big.NewRat(1, 100),  // 0.01
			expectedValue: sdk.MustNewDecFromStr("99652710.304588509013984918"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := clpkeeper.CalculateNativeSwapAmountAsymmetricRat(tc.Y, tc.X, tc.y, tc.x, tc.f, tc.r)
			got, _ := clpkeeper.RatToDec(&res)

			require.Equal(t, tc.expectedValue.String(), got.String())

		})
	}

}

func MustRatFromString(x string) *big.Rat {
	res, success := big.NewRat(1, 1).SetString("157007500498726220240179086")
	if success == false {
		panic("Could not create rat from string")
	}
	return res
}

func TestKeeper_GetLiquidityAddSymmetryType(t *testing.T) {
	testcases := []struct {
		name          string
		X, x, Y, y    sdk.Uint
		expectedValue int
	}{
		{
			name:          "one side of the pool empty",
			X:             sdk.ZeroUint(),
			x:             sdk.NewUint(11200),
			Y:             sdk.NewUint(100),
			y:             sdk.NewUint(100),
			expectedValue: clpkeeper.ErrorEmptyPool,
		},
		{
			name:          "nothing added",
			X:             sdk.NewUint(11200),
			x:             sdk.ZeroUint(),
			Y:             sdk.NewUint(1000),
			y:             sdk.ZeroUint(),
			expectedValue: clpkeeper.ErrorNothingAdded,
		},
		{
			name:          "negative symmetry - x zero",
			X:             sdk.NewUint(11200),
			x:             sdk.ZeroUint(),
			Y:             sdk.NewUint(1000),
			y:             sdk.NewUint(100),
			expectedValue: clpkeeper.NeedMoreX,
		},
		{
			name:          "negative symmetry - x > 0",
			X:             sdk.NewUint(11200),
			x:             sdk.NewUint(15),
			Y:             sdk.NewUint(1000),
			y:             sdk.NewUint(100),
			expectedValue: clpkeeper.NeedMoreX,
		},
		{
			name:          "symmetric",
			X:             sdk.NewUint(11200),
			x:             sdk.NewUint(1120),
			Y:             sdk.NewUint(1000),
			y:             sdk.NewUint(100),
			expectedValue: clpkeeper.Symmetric,
		},
		{
			name:          "positive symmetry",
			X:             sdk.NewUint(11200),
			x:             sdk.NewUint(100),
			Y:             sdk.NewUint(1000),
			y:             sdk.NewUint(5),
			expectedValue: clpkeeper.NeedMoreY,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := clpkeeper.GetLiquidityAddSymmetryState(tc.X, tc.x, tc.Y, tc.y)

			require.Equal(t, tc.expectedValue, res)
		})
	}
}

func TestKeeper_CalculatePoolUnits(t *testing.T) {
	testcases := []struct {
		name                  string
		oldPoolUnits          sdk.Uint
		nativeAssetBalance    sdk.Uint
		externalAssetBalance  sdk.Uint
		nativeAssetAmount     sdk.Uint
		externalAssetAmount   sdk.Uint
		sellNativeSwapFeeRate sdk.Dec
		buyNativeSwapFeeRate  sdk.Dec
		expectedPoolUnits     sdk.Uint
		expectedLPunits       sdk.Uint
		expectedSwapStatus    int
		expectedSwapAmount    sdk.Uint
		expectedError         error
	}{
		{
			name:                  "empty pool",
			oldPoolUnits:          sdk.ZeroUint(),
			nativeAssetBalance:    sdk.ZeroUint(),
			externalAssetBalance:  sdk.ZeroUint(),
			nativeAssetAmount:     sdk.NewUint(100),
			externalAssetAmount:   sdk.NewUint(90),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUint(100),
			expectedLPunits:       sdk.NewUint(100),
			expectedSwapStatus:    clpkeeper.NoSwap,
		},
		{
			name:                  "empty pool - no external asset added",
			oldPoolUnits:          sdk.ZeroUint(),
			nativeAssetBalance:    sdk.ZeroUint(),
			externalAssetBalance:  sdk.ZeroUint(),
			nativeAssetAmount:     sdk.NewUint(100),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			externalAssetAmount:   sdk.ZeroUint(),
			expectedError:         errors.New("amount is invalid"),
		},
		{
			name:                  "add nothing",
			oldPoolUnits:          sdk.NewUint(1000),
			nativeAssetBalance:    sdk.NewUint(12327),
			externalAssetBalance:  sdk.NewUint(132233),
			nativeAssetAmount:     sdk.ZeroUint(),
			externalAssetAmount:   sdk.ZeroUint(),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUint(1000),
			expectedLPunits:       sdk.ZeroUint(),
			expectedSwapStatus:    clpkeeper.NoSwap,
		},
		{
			name:                  "positive symmetry - zero native",
			oldPoolUnits:          sdk.NewUint(7656454334323412),
			nativeAssetBalance:    sdk.NewUint(16767626535600),
			externalAssetBalance:  sdk.NewUint(2345454545400),
			nativeAssetAmount:     sdk.ZeroUint(),
			externalAssetAmount:   sdk.NewUint(4556664545),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUint(7663887695258361),
			expectedLPunits:       sdk.NewUint(7433360934949),
			expectedSwapStatus:    clpkeeper.BuyNative,
			expectedSwapAmount:    sdk.NewUint(2277340758),
		},
		{
			name:                  "symmetric",
			oldPoolUnits:          sdk.NewUint(7656454334323412),
			nativeAssetBalance:    sdk.NewUint(16767626535600),
			externalAssetBalance:  sdk.NewUint(2345454545400),
			nativeAssetAmount:     sdk.NewUint(167676265356),
			externalAssetAmount:   sdk.NewUint(23454545454),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUint(7733018877666646),
			expectedLPunits:       sdk.NewUint(76564543343234),
			expectedSwapStatus:    clpkeeper.NoSwap,
			expectedSwapAmount:    sdk.ZeroUint(),
		},
		{
			name:                  "negative symmetry - zero external",
			oldPoolUnits:          sdk.NewUint(7656454334323412),
			nativeAssetBalance:    sdk.NewUint(16767626535600),
			externalAssetBalance:  sdk.NewUint(2345454545400),
			nativeAssetAmount:     sdk.NewUint(167676265356),
			externalAssetAmount:   sdk.ZeroUint(),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUint(7694639456903696),
			expectedLPunits:       sdk.NewUint(38185122580284),
			expectedSwapStatus:    clpkeeper.SellNative,
			expectedSwapAmount:    sdk.NewUint(83633781363),
		},
		{
			name:                  "positive symmetry - non zero external",
			oldPoolUnits:          sdk.NewUint(7656454334323412),
			nativeAssetBalance:    sdk.NewUint(16767626535600),
			externalAssetBalance:  sdk.NewUint(2345454545400),
			nativeAssetAmount:     sdk.NewUint(167676265356),
			externalAssetAmount:   sdk.NewUint(46798998888),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUint(7771026137435008),
			expectedLPunits:       sdk.NewUint(114571803111596),
			expectedSwapStatus:    clpkeeper.BuyNative,
			expectedSwapAmount:    sdk.NewUint(11528907497),
		},
		{
			name:                  "very big - positive symmetry",
			oldPoolUnits:          sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"), //2**200
			nativeAssetBalance:    sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			externalAssetBalance:  sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			nativeAssetAmount:     sdk.NewUint(0),
			externalAssetAmount:   sdk.NewUint(1099511627776), // 2**40
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUintFromString("1606938044258990275541962092341162602522202993783342563626098"),
			expectedLPunits:       sdk.NewUint(549728324722),
			expectedSwapStatus:    clpkeeper.BuyNative,
			expectedSwapAmount:    sdk.NewUint(549783303053),
		},
		{
			name:                  "very big - symmetric",
			oldPoolUnits:          sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"), //2**200
			nativeAssetBalance:    sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			externalAssetBalance:  sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			nativeAssetAmount:     sdk.NewUint(1099511627776), // 2**40
			externalAssetAmount:   sdk.NewUint(1099511627776),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUintFromString("1606938044258990275541962092341162602522202993783892346929152"),
			expectedLPunits:       sdk.NewUint(1099511627776),
			expectedSwapStatus:    clpkeeper.NoSwap,
			expectedSwapAmount:    sdk.ZeroUint(),
		},
		{
			name:                  "very big - negative symmetry",
			oldPoolUnits:          sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"), //2**200
			nativeAssetBalance:    sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			externalAssetBalance:  sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			nativeAssetAmount:     sdk.NewUint(1099511627776), // 2**40
			externalAssetAmount:   sdk.ZeroUint(),
			sellNativeSwapFeeRate: sdk.NewDecWithPrec(1, 4),
			buyNativeSwapFeeRate:  sdk.NewDecWithPrec(1, 4),
			expectedPoolUnits:     sdk.NewUintFromString("1606938044258990275541962092341162602522202993783342563626098"),
			expectedLPunits:       sdk.NewUint(549728324722),
			expectedSwapStatus:    clpkeeper.SellNative,
			expectedSwapAmount:    sdk.NewUint(549783303053),
		},
		{
			name:                  "swap fee rates = 1, zero external asset",
			oldPoolUnits:          sdk.NewUint(7656454334323412),
			nativeAssetBalance:    sdk.NewUint(16767626535600),
			externalAssetBalance:  sdk.NewUint(2345454545400),
			nativeAssetAmount:     sdk.NewUint(167676265356),
			externalAssetAmount:   sdk.ZeroUint(),
			sellNativeSwapFeeRate: sdk.OneDec(),
			buyNativeSwapFeeRate:  sdk.OneDec(),
			expectedPoolUnits:     sdk.NewUint(7656454334323412),
			expectedLPunits:       sdk.NewUint(0),
			expectedSwapStatus:    clpkeeper.SellNative,
			expectedSwapAmount:    sdk.NewUint(167676265356),
		},
		{
			name:                  "swap fee rates = 0, zero external asset",
			oldPoolUnits:          sdk.NewUint(7656454334323412),
			nativeAssetBalance:    sdk.NewUint(16767626535600),
			externalAssetBalance:  sdk.NewUint(2345454545400),
			nativeAssetAmount:     sdk.NewUint(167676265356),
			externalAssetAmount:   sdk.ZeroUint(),
			sellNativeSwapFeeRate: sdk.ZeroDec(),
			buyNativeSwapFeeRate:  sdk.ZeroDec(),
			expectedPoolUnits:     sdk.NewUint(7694641375874505),
			expectedLPunits:       sdk.NewUint(38187041551093),
			expectedSwapStatus:    clpkeeper.SellNative,
			expectedSwapAmount:    sdk.NewUint(83629578818),
		},
	}

	pmtpCurrentRunningRate := sdk.ZeroDec()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			poolUnits, lpunits, swapStatus, swapAmount, err := clpkeeper.CalculatePoolUnits(
				tc.oldPoolUnits,
				tc.nativeAssetBalance,
				tc.externalAssetBalance,
				tc.nativeAssetAmount,
				tc.externalAssetAmount,
				tc.sellNativeSwapFeeRate,
				tc.buyNativeSwapFeeRate,
				pmtpCurrentRunningRate,
			)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.expectedPoolUnits.String(), poolUnits.String()) // compare strings so that the expected amounts can be read from the failure message
			require.Equal(t, tc.expectedLPunits.String(), lpunits.String())
			require.Equal(t, tc.expectedSwapStatus, swapStatus)
			require.Equal(t, tc.expectedSwapAmount.String(), swapAmount.String())

		})
	}
}

func TestKeeper_CalculatePoolUnitsSymmetric(t *testing.T) {
	testcases := []struct {
		name              string
		X, x, P           sdk.Uint
		expectedPoolUnits sdk.Uint
		expectedLPUnits   sdk.Uint
	}{
		{
			name:              "test 1",
			X:                 sdk.NewUint(167676265356),
			x:                 sdk.NewUint(5120000099),
			P:                 sdk.NewUint(112323227872),
			expectedPoolUnits: sdk.NewUint(115753021209),
			expectedLPUnits:   sdk.NewUint(3429793337),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			poolUnits, lpUnits := clpkeeper.CalculatePoolUnitsSymmetric(tc.X, tc.x, tc.P)

			require.Equal(t, tc.expectedPoolUnits.String(), poolUnits.String())
			require.Equal(t, tc.expectedLPUnits.String(), lpUnits.String())
		})
	}
}

func TestKeeper_SwapOneFromGenesis(t *testing.T) {
	const address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	SwapPriceNative := sdk.ZeroDec()
	SwapPriceExternal := sdk.ZeroDec()

	testcases := []struct {
		name                   string
		poolAsset              string
		address                string
		calculateWithdraw      bool
		adjustExternalToken    bool
		nativeBalance          sdk.Int
		externalBalance        sdk.Int
		wBasis                 sdk.Int
		asymmetry              sdk.Int
		nativeAssetAmount      sdk.Uint
		externalAssetAmount    sdk.Uint
		poolUnits              sdk.Uint
		swapAmount             sdk.Uint
		swapResult             sdk.Uint
		liquidityFee           sdk.Uint
		priceImpact            sdk.Uint
		normalizationFactor    sdk.Dec
		pmtpCurrentRunningRate sdk.Dec
		from                   types.Asset
		to                     types.Asset
		expectedPool           types.Pool
		err                    error
		errString              error
	}{
		{
			name:                   "successful swap with equal amount of pool units",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.OneDec(),
			swapResult:             sdk.NewUint(180),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(818),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "failed swap with empty pool",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(0),
			externalAssetAmount:    sdk.NewUint(0),
			poolUnits:              sdk.NewUint(0),
			calculateWithdraw:      false,
			normalizationFactor:    sdk.NewDec(0),
			adjustExternalToken:    true,
			swapAmount:             sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.OneDec(),
			swapResult:             sdk.NewUint(166),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(833),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
			errString: errors.New("not enough received asset tokens to swap"),
		},
		{
			name:                   "successful swap by inversing from/to assets",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			from:                   types.Asset{Symbol: "eth"},
			to:                     types.Asset{Symbol: "rowan"},
			pmtpCurrentRunningRate: sdk.OneDec(),
			swapResult:             sdk.NewUint(45),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(953),
				ExternalAssetBalance:          sdk.NewUint(1097),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.0"),
			swapResult:             sdk.NewUint(90),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(908),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.1",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.1"),
			swapResult:             sdk.NewUint(99),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(899),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.2",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.2"),
			swapResult:             sdk.NewUint(108),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(890),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.3",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.3"),
			swapResult:             sdk.NewUint(117),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(881),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.4",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.4"),
			swapResult:             sdk.NewUint(126),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(872),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.5",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.5"),
			swapResult:             sdk.NewUint(135),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(863),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.6",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.6"),
			swapResult:             sdk.NewUint(144),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(854),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.7",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.7"),
			swapResult:             sdk.NewUint(153),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(845),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.8",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.8"),
			swapResult:             sdk.NewUint(162),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(836),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 0.9",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("0.9"),
			swapResult:             sdk.NewUint(171),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(827),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 1.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.0"),
			swapResult:             sdk.NewUint(180),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(818),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 2.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("2.0"),
			swapResult:             sdk.NewUint(270),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(728),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 3.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("3.0"),
			swapResult:             sdk.NewUint(359),
			liquidityFee:           sdk.NewUint(1),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(639),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 4.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("4.0"),
			swapResult:             sdk.NewUint(449),
			liquidityFee:           sdk.NewUint(1),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(549),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 5.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("5.0"),
			swapResult:             sdk.NewUint(539),
			liquidityFee:           sdk.NewUint(1),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(459),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 6.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("6.0"),
			swapResult:             sdk.NewUint(629),
			liquidityFee:           sdk.NewUint(1),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(369),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 7.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("7.0"),
			swapResult:             sdk.NewUint(718),
			liquidityFee:           sdk.NewUint(2),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(280),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 8.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("8.0"),
			swapResult:             sdk.NewUint(808),
			liquidityFee:           sdk.NewUint(2),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(190),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 9.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("9.0"),
			swapResult:             sdk.NewUint(898),
			liquidityFee:           sdk.NewUint(2),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(100),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "successful swap with pmtp current running rate value at 10.0",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("10.0"),
			swapResult:             sdk.NewUint(988),
			liquidityFee:           sdk.NewUint(2),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1097),
				ExternalAssetBalance:          sdk.NewUint(10),
				PoolUnits:                     sdk.NewUint(998),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			},
		},
		{
			name:                   "failed swap with bigger pmtp current running rate value",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.NewDec(20),
			errString:              errors.New("not enough received asset tokens to swap"),
		},
		{
			name:                   "failed swap with bigger pmtp current running rate value",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.NewUint(998),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.NewDec(20),
			errString:              errors.New("not enough received asset tokens to swap"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				trGs := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: tc.poolAsset, BaseDenom: tc.poolAsset, Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
							{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(trGs)
				genesisState["tokenregistry"] = bz

				balances := []banktypes.Balance{
					{
						Address: tc.address,
						Coins: sdk.Coins{
							sdk.NewCoin(tc.poolAsset, tc.externalBalance),
							sdk.NewCoin("rowan", tc.nativeBalance),
						},
					},
				}
				bankGs := banktypes.DefaultGenesisState()
				bankGs.Balances = append(bankGs.Balances, balances...)
				bz, _ = app.AppCodec().MarshalJSON(bankGs)
				genesisState["bank"] = bz

				pools := []*types.Pool{
					{
						ExternalAsset:                 &types.Asset{Symbol: tc.poolAsset},
						NativeAssetBalance:            tc.nativeAssetAmount,
						ExternalAssetBalance:          tc.externalAssetAmount,
						PoolUnits:                     tc.poolUnits,
						NativeCustody:                 sdk.ZeroUint(),
						ExternalCustody:               sdk.ZeroUint(),
						NativeLiabilities:             sdk.ZeroUint(),
						ExternalLiabilities:           sdk.ZeroUint(),
						Health:                        sdk.ZeroDec(),
						InterestRate:                  sdk.NewDecWithPrec(1, 1),
						SwapPriceNative:               &SwapPriceNative,
						SwapPriceExternal:             &SwapPriceExternal,
						RewardPeriodNativeDistributed: sdk.ZeroUint(),
						RewardAmountExternal:          sdk.ZeroUint(),
					},
				}
				lps := []*types.LiquidityProvider{
					{
						Asset:                    &types.Asset{Symbol: tc.poolAsset},
						LiquidityProviderAddress: tc.address,
						LiquidityProviderUnits:   tc.nativeAssetAmount,
					},
				}
				clpGs := types.DefaultGenesisState()
				clpGs.Params = types.Params{
					MinCreatePoolThreshold: 100,
				}
				clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
				clpGs.PoolList = append(clpGs.PoolList, pools...)
				clpGs.LiquidityProviders = append(clpGs.LiquidityProviders, lps...)
				bz, _ = app.AppCodec().MarshalJSON(clpGs)
				genesisState["clp"] = bz

				return genesisState
			})

			pool, _ := app.ClpKeeper.GetPool(ctx, tc.poolAsset)
			lp, _ := app.ClpKeeper.GetLiquidityProvider(ctx, tc.poolAsset, tc.address)

			SwapPriceNative := sdk.ZeroDec()
			SwapPriceExternal := sdk.ZeroDec()

			require.Equal(t, pool, types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: tc.poolAsset},
				NativeAssetBalance:            tc.nativeAssetAmount,
				ExternalAssetBalance:          tc.externalAssetAmount,
				PoolUnits:                     tc.poolUnits,
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			})

			var swapAmount sdk.Uint

			if tc.calculateWithdraw {
				_, _, _, swapAmount = clpkeeper.CalculateWithdrawal(
					pool.PoolUnits,
					pool.NativeAssetBalance.String(),
					pool.ExternalAssetBalance.String(),
					lp.LiquidityProviderUnits.String(),
					tc.wBasis.String(),
					tc.asymmetry,
				)
			} else {
				swapAmount = tc.swapAmount
			}

			from := tc.from
			if from == (types.Asset{}) {
				from = types.GetSettlementAsset()
			}
			to := tc.to
			if to == (types.Asset{}) {
				to = types.Asset{Symbol: tc.poolAsset}
			}

			swapFeeRate := sdk.NewDecWithPrec(3, 3)

			swapResult, liquidityFee, priceImpact, newPool, err := clpkeeper.SwapOne(from, swapAmount, to, pool, tc.pmtpCurrentRunningRate, swapFeeRate)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.swapResult.String(), swapResult.String(), "swapResult")
			require.Equal(t, tc.liquidityFee.String(), liquidityFee.String())
			require.Equal(t, tc.priceImpact.String(), priceImpact.String())
			require.Equal(t, tc.expectedPool.String(), newPool.String())
		})
	}
}
