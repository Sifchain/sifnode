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
		name                         string
		nativeAssetBalance           sdk.Uint
		externalAssetBalance         sdk.Uint
		sentAmount                   sdk.Uint
		fromAsset                    types.Asset
		toAsset                      types.Asset
		pmtpCurrentRunningRate       sdk.Dec
		errString                    error
		expectedSwapResult           sdk.Uint
		expectedLiquidityFee         sdk.Uint
		expectedPriceImpact          sdk.Uint
		expectedExternalAssetBalance sdk.Uint
		expectedNativeAssetBalance   sdk.Uint
	}{
		{
			name:                         "test1",
			nativeAssetBalance:           sdk.NewUint(1000),
			externalAssetBalance:         sdk.NewUint(877),
			sentAmount:                   sdk.NewUint(50),
			fromAsset:                    types.GetSettlementAsset(),
			toAsset:                      types.NewAsset("eth"),
			pmtpCurrentRunningRate:       sdk.NewDec(0),
			expectedSwapResult:           sdk.NewUint(41),
			expectedLiquidityFee:         sdk.NewUint(0),
			expectedPriceImpact:          sdk.ZeroUint(),
			expectedExternalAssetBalance: sdk.NewUint(836),
			expectedNativeAssetBalance:   sdk.NewUint(1050),
		},
		{
			name:                         "test2",
			nativeAssetBalance:           sdk.NewUint(10000000),
			externalAssetBalance:         sdk.NewUint(8770000),
			sentAmount:                   sdk.NewUint(50000),
			fromAsset:                    types.GetSettlementAsset(),
			toAsset:                      types.NewAsset("eth"),
			pmtpCurrentRunningRate:       sdk.NewDec(0),
			expectedSwapResult:           sdk.NewUint(43500),
			expectedLiquidityFee:         sdk.NewUint(130),
			expectedPriceImpact:          sdk.ZeroUint(),
			expectedExternalAssetBalance: sdk.NewUint(8726500),
			expectedNativeAssetBalance:   sdk.NewUint(10050000),
		},
		{
			name:                         "test3",
			nativeAssetBalance:           sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance:         sdk.NewUint(2674623482959),
			sentAmount:                   sdk.NewUint(200000000),
			toAsset:                      types.GetSettlementAsset(),
			fromAsset:                    types.NewAsset("cusdt"),
			pmtpCurrentRunningRate:       sdk.NewDec(0),
			expectedSwapResult:           sdk.NewUintFromString("11704434254784015637542"),
			expectedLiquidityFee:         sdk.NewUintFromString("35218959643281892590"),
			expectedPriceImpact:          sdk.ZeroUint(),
			expectedExternalAssetBalance: sdk.NewUint(2674823482959),
			expectedNativeAssetBalance:   sdk.NewUintFromString("156995796064471436224541544"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClp(false)
			poolUnits := sdk.NewUint(2000) //don't care
			pool := types.NewPool(&tc.toAsset, tc.nativeAssetBalance, tc.externalAssetBalance, poolUnits)

			marginEnabled := app.ClpKeeper.GetMarginKeeper().IsPoolEnabled(ctx, pool.String())
			swapResult, liquidityFee, priceImpact, pool, err := clpkeeper.SwapOne(tc.fromAsset, tc.sentAmount, tc.toAsset, pool, tc.pmtpCurrentRunningRate, sdk.NewDecWithPrec(3, 3), marginEnabled)

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
	X, Y, toRowan := pool.ExtractValues(asset)

	assert.Equal(t, X, sdk.NewUint(998))
	assert.Equal(t, Y, sdk.NewUint(998))
	assert.Equal(t, toRowan, false)
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
	marginEnabled := app.ClpKeeper.GetMarginKeeper().IsPoolEnabled(ctx, pool.String())
	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, *pool, sdk.OneDec(), sdk.NewDecWithPrec(3, 3), marginEnabled)
	assert.Equal(t, "1", swapResult.String())
}

func TestKeeper_GetSwapFee_PmtpParams(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	pool := types.Pool{
		NativeAssetBalance:   sdk.NewUint(10),
		ExternalAssetBalance: sdk.NewUint(100),
	}
	asset := types.Asset{}

	marginEnabled := app.ClpKeeper.GetMarginKeeper().IsPoolEnabled(ctx, pool.String())
	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, pool, sdk.NewDec(100), sdk.NewDecWithPrec(3, 3), marginEnabled)

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

func TestKeeper_CalculatePoolUnits(t *testing.T) {
	testcases := []struct {
		name                 string
		oldPoolUnits         sdk.Uint
		nativeAssetBalance   sdk.Uint
		externalAssetBalance sdk.Uint
		nativeAssetAmount    sdk.Uint
		externalAssetAmount  sdk.Uint
		externalDecimals     uint8
		poolUnits            sdk.Uint
		lpunits              sdk.Uint
		err                  error
		errString            error
		panicErr             string
	}{
		{
			name:                 "tx amount too low throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.ZeroUint(),
			externalAssetAmount:  sdk.ZeroUint(),
			externalDecimals:     18,
			errString:            errors.New("Tx amount is too low"),
		},
		{
			name:                 "tx amount too low with no adjustment throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.ZeroUint(),
			externalAssetAmount:  sdk.ZeroUint(),
			externalDecimals:     18,
			errString:            errors.New("Tx amount is too low"),
		},
		{
			name:                 "insufficient native funds throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.ZeroUint(),
			externalAssetAmount:  sdk.OneUint(),
			externalDecimals:     18,
			errString:            errors.New("0: insufficient funds"),
		},
		{
			name:                 "insufficient external funds throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.ZeroUint(),
			externalDecimals:     18,
			errString:            errors.New("0: insufficient funds"),
		},
		{
			name:                 "as native asset balance zero then returns native asset amount",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			externalDecimals:     18,
			poolUnits:            sdk.OneUint(),
			lpunits:              sdk.OneUint(),
		},
		{
			name:                 "successful",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			externalDecimals:     18,
			poolUnits:            sdk.ZeroUint(),
			lpunits:              sdk.ZeroUint(),
		},
		{
			name:                 "fail asymmetric",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(10000),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			externalDecimals:     18,
			poolUnits:            sdk.ZeroUint(),
			lpunits:              sdk.ZeroUint(),
			errString:            errors.New("Cannot add liquidity asymmetrically"),
		},
		{
			name:                 "successful",
			oldPoolUnits:         sdk.NewUint(1),
			nativeAssetBalance:   sdk.NewUint(1),
			externalAssetBalance: sdk.NewUint(1),
			nativeAssetAmount:    sdk.NewUint(1),
			externalAssetAmount:  sdk.NewUint(1),
			externalDecimals:     18,
			poolUnits:            sdk.NewUint(2),
			lpunits:              sdk.NewUint(1),
		},
		{
			name:                 "successful no slip",
			oldPoolUnits:         sdk.NewUint(1099511627776), //2**40
			nativeAssetBalance:   sdk.NewUint(1099511627776),
			externalAssetBalance: sdk.NewUint(1099511627776),
			nativeAssetAmount:    sdk.NewUint(1099511627776),
			externalAssetAmount:  sdk.NewUint(1099511627776),
			externalDecimals:     18,
			poolUnits:            sdk.NewUint(2199023255552),
			lpunits:              sdk.NewUint(1099511627776),
		},
		{
			name:                 "no asymmetric",
			oldPoolUnits:         sdk.NewUint(1099511627776), //2**40
			nativeAssetBalance:   sdk.NewUint(1048576),
			externalAssetBalance: sdk.NewUint(1024123),
			nativeAssetAmount:    sdk.NewUint(999),
			externalAssetAmount:  sdk.NewUint(111),
			externalDecimals:     18,
			poolUnits:            sdk.NewUintFromString("1100094484982"),
			lpunits:              sdk.NewUintFromString("582857206"),
			errString:            errors.New("Cannot add liquidity asymmetrically"),
		},
		{
			name:                 "successful - very big",
			oldPoolUnits:         sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"), //2**200
			nativeAssetBalance:   sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			externalAssetBalance: sdk.NewUintFromString("1606938044258990275541962092341162602522202993782792835301376"),
			nativeAssetAmount:    sdk.NewUint(1099511627776), // 2**40
			externalAssetAmount:  sdk.NewUint(1099511627776),
			externalDecimals:     18,
			poolUnits:            sdk.NewUintFromString("1606938044258990275541962092341162602522202993783892346929152"),
			lpunits:              sdk.NewUint(1099511627776),
		},
		{
			name:                 "failure - asymmetric",
			oldPoolUnits:         sdk.NewUintFromString("23662660550457383692937954"),
			nativeAssetBalance:   sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance: sdk.NewUint(2674623482959),
			nativeAssetAmount:    sdk.NewUint(0),
			externalAssetAmount:  sdk.NewUint(200000000),
			externalDecimals:     18,
			errString:            errors.New("Cannot add liquidity with asymmetric ratio"),
		},
		{
			name:                 "opportunist scenario - fails trivially due to div zero",
			oldPoolUnits:         sdk.NewUintFromString("23662660550457383692937954"),
			nativeAssetBalance:   sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance: sdk.NewUint(2674623482959),
			nativeAssetAmount:    sdk.NewUint(0),
			externalAssetAmount:  sdk.NewUint(200000000),
			externalDecimals:     6,
			errString:            errors.New("Cannot add liquidity with asymmetric ratio"),
		},
		{
			name:                 "opportunist scenario with one native asset - avoids div zero trivial fail",
			oldPoolUnits:         sdk.NewUintFromString("23662660550457383692937954"),
			nativeAssetBalance:   sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance: sdk.NewUint(2674623482959),
			nativeAssetAmount:    sdk.NewUint(1),
			externalAssetAmount:  sdk.NewUint(200000000),
			externalDecimals:     6,
			errString:            errors.New("Cannot add liquidity with asymmetric ratio"),
		},
		{
			name:                 "success",
			oldPoolUnits:         sdk.NewUintFromString("23662660550457383692937954"),
			nativeAssetBalance:   sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance: sdk.NewUint(2674623482959),
			nativeAssetAmount:    sdk.NewUintFromString("4000000000000000000"),
			externalAssetAmount:  sdk.NewUint(68140),
			externalDecimals:     6,
			poolUnits:            sdk.NewUintFromString("23662661153298835875523384"),
			lpunits:              sdk.NewUintFromString("602841452182585430"),
		},
		{
			// Same test as above but with external asset amount just below top limit
			name:                 "success (normalized) ratios diff = 0.000000000000000499",
			oldPoolUnits:         sdk.NewUintFromString("23662660550457383692937954"),
			nativeAssetBalance:   sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance: sdk.NewUint(2674623482959),
			nativeAssetAmount:    sdk.NewUintFromString("4000000000000000000"),
			externalAssetAmount:  sdk.NewUint(70140),
			externalDecimals:     6,
			poolUnits:            sdk.NewUintFromString("23662661162145935094484778"),
			lpunits:              sdk.NewUintFromString("611688551401546824"),
		},
		{
			// Same test as above but with external asset amount just above top limit
			name:                 "failure (normalized) ratios diff = 0.000000000000000500",
			oldPoolUnits:         sdk.NewUintFromString("23662660550457383692937954"),
			nativeAssetBalance:   sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance: sdk.NewUint(2674623482959),
			nativeAssetAmount:    sdk.NewUintFromString("4000000000000000000"),
			externalAssetAmount:  sdk.NewUint(70141),
			externalDecimals:     6,
			errString:            errors.New("Cannot add liquidity with asymmetric ratio"),
		},
		{
			// Same test as above but with external asset amount just above bottom limit
			name:                 "success (normalized) ratios diff = 0.000000000000000499",
			oldPoolUnits:         sdk.NewUintFromString("23662660550457383692937954"),
			nativeAssetBalance:   sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance: sdk.NewUint(2674623482959),
			nativeAssetAmount:    sdk.NewUintFromString("4000000000000000000"),
			externalAssetAmount:  sdk.NewUint(66141),
			externalDecimals:     6,
			poolUnits:            sdk.NewUintFromString("23662661144456159305055227"),
			lpunits:              sdk.NewUintFromString("593998775612117273"),
		},
		{
			// Same test as above but with external asset amount just below bottom limit
			name:                 "failure (normalized) ratios diff = 0.000000000000000500",
			oldPoolUnits:         sdk.NewUintFromString("23662660550457383692937954"),
			nativeAssetBalance:   sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetBalance: sdk.NewUint(2674623482959),
			nativeAssetAmount:    sdk.NewUintFromString("4000000000000000000"),
			externalAssetAmount:  sdk.NewUint(66140),
			externalDecimals:     6,
			errString:            errors.New("Cannot add liquidity with asymmetric ratio"),
		},
	}

	symmetryThreshold := sdk.NewDecWithPrec(1, 4)
	ratioThreshold := sdk.NewDecWithPrec(5, 4)
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.panicErr != "" {
				// nolint:errcheck
				require.PanicsWithError(t, tc.panicErr, func() {
					clpkeeper.CalculatePoolUnits(
						tc.oldPoolUnits,
						tc.nativeAssetBalance,
						tc.externalAssetBalance,
						tc.nativeAssetAmount,
						tc.externalAssetAmount,
						tc.externalDecimals,
						symmetryThreshold,
						ratioThreshold,
					)
				})
				return
			}

			poolUnits, lpunits, err := clpkeeper.CalculatePoolUnits(
				tc.oldPoolUnits,
				tc.nativeAssetBalance,
				tc.externalAssetBalance,
				tc.nativeAssetAmount,
				tc.externalAssetAmount,
				tc.externalDecimals,
				symmetryThreshold,
				ratioThreshold,
			)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.poolUnits.String(), poolUnits.String()) // compare strings so that the expected amounts can be read from the failure message
			require.Equal(t, tc.lpunits.String(), lpunits.String())
		})
	}
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
		tc := tc
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

func TestKeeper_CalcLiquidityFee(t *testing.T) {
	testcases := []struct {
		name                string
		toRowan             bool
		adjustExternalToken bool
		normalizationFactor sdk.Dec
		X, x, Y, fee        sdk.Uint
		err                 error
		errString           error
	}{
		{
			name:    "success",
			toRowan: true,
			X:       sdk.NewUint(0),
			x:       sdk.NewUint(0),
			Y:       sdk.NewUint(1),
			fee:     sdk.NewUint(0),
		},
		{
			name:    "success",
			toRowan: true,
			X:       sdk.NewUint(1),
			x:       sdk.NewUint(1),
			Y:       sdk.NewUint(1),
			fee:     sdk.NewUint(0),
		},
		{
			name:    "success",
			toRowan: true,
			X:       sdk.NewUint(1),
			x:       sdk.NewUint(1),
			Y:       sdk.NewUint(4),
			fee:     sdk.NewUint(0),
		},
		{
			name:    "success",
			toRowan: true,
			X:       sdk.NewUint(2),
			x:       sdk.NewUint(2),
			Y:       sdk.NewUint(16),
			fee:     sdk.NewUint(0),
		},
		{
			name:    "success",
			toRowan: true,
			X:       sdk.NewUint(1054677676764),
			x:       sdk.NewUint(2567655449999),
			Y:       sdk.NewUint(1099511627776),
			fee:     sdk.NewUint(2292288459),
		},
		{
			name:    "success",
			toRowan: true,
			X:       sdk.NewUintFromString("20300000000000000000000000000000000000000000000000000000000000000000000000"),
			x:       sdk.NewUintFromString("10000000000000000658000000000000000000000000000000000000000000000000000000"),
			Y:       sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000021344"),
			fee:     sdk.NewUintFromString("9706853038245001398600547299925090963413469041036560829464056119132287"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fee := clpkeeper.CalcLiquidityFee(tc.toRowan, tc.X, tc.x, tc.Y, sdk.NewDecWithPrec(3, 3), sdk.NewDecWithPrec(2, 2))
			require.Equal(t, tc.fee.String(), fee.String()) // compare strings so that the expected amounts can be read from the failure message
		})
	}
}

func TestKeeper_CalcSwapResult(t *testing.T) {
	testcases := []struct {
		name                   string
		toRowan                bool
		X, x, Y, y             sdk.Uint
		pmtpCurrentRunningRate sdk.Dec
		swapFeeRate            sdk.Dec
		err                    error
		errString              error
	}{
		{
			name:                   "adjust external token with rowan",
			toRowan:                true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(1),
			y:                      sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(1),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "adjust external token without rowan",
			toRowan:                false,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(1),
			y:                      sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(1),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "x=0, X=0, Y=0",
			toRowan:                true,
			X:                      sdk.NewUint(0),
			x:                      sdk.NewUint(0),
			Y:                      sdk.NewUint(0),
			y:                      sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(0),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "x=1, X=1, Y=1",
			toRowan:                true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(1),
			y:                      sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(0),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "x=1, X=1, Y=4",
			toRowan:                true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(4),
			y:                      sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(0),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "x=1, X=1, Y=4, nf=10",
			toRowan:                true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(4),
			y:                      sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(0),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "x=23, X=42, Y=1337",
			toRowan:                true,
			X:                      sdk.NewUint(42),
			x:                      sdk.NewUint(23),
			Y:                      sdk.NewUint(1337),
			y:                      sdk.NewUint(471),
			pmtpCurrentRunningRate: sdk.NewDec(0),
			swapFeeRate:            sdk.NewDecWithPrec(3, 3),
		},
		{
			name:                   "swap fee rate = 0.01",
			toRowan:                false,
			X:                      sdk.NewUint(1999800619938006200),
			x:                      sdk.NewUint(200000000000000),
			Y:                      sdk.NewUint(2000200000000000000),
			y:                      sdk.NewUint(198019738620019),
			pmtpCurrentRunningRate: sdk.NewDec(0),
			swapFeeRate:            sdk.NewDecWithPrec(1, 2),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			y := clpkeeper.CalcSwapResult(tc.toRowan, tc.X, tc.x, tc.Y, tc.pmtpCurrentRunningRate, tc.swapFeeRate)

			require.Equal(t, tc.y.String(), y.String()) // compare strings so that the expected amounts can be read from the failure message
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
		tc := tc
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
		tc := tc
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pool := types.Pool{
				NativeAssetBalance:   tc.nativeAssetBalance,
				ExternalAssetBalance: tc.externalAssetBalance,
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pool := types.Pool{
				NativeAssetBalance:   tc.nativeAssetBalance,
				ExternalAssetBalance: tc.externalAssetBalance,
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

func TestKeeper_CalculateRatioDiff(t *testing.T) {

	testcases := []struct {
		name       string
		A, R, a, r *big.Int
		expected   sdk.Dec
		errString  error
	}{
		{
			name:     "symmetric",
			A:        big.NewInt(20),
			R:        big.NewInt(10),
			a:        big.NewInt(8),
			r:        big.NewInt(4),
			expected: sdk.MustNewDecFromStr("0.000000000000000000"),
		},
		{
			name:     "not symmetric",
			A:        big.NewInt(20),
			R:        big.NewInt(10),
			a:        big.NewInt(16),
			r:        big.NewInt(4),
			expected: sdk.MustNewDecFromStr("2.000000000000000000"),
		},
		{
			name:     "not symmetric",
			A:        big.NewInt(501),
			R:        big.NewInt(100),
			a:        big.NewInt(5),
			r:        big.NewInt(1),
			expected: sdk.MustNewDecFromStr("0.010000000000000000"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			ratio, err := clpkeeper.CalculateRatioDiff(tc.A, tc.R, tc.a, tc.r)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}

			require.NoError(t, err)

			ratioDec, _ := clpkeeper.RatToDec(&ratio)

			require.Equal(t, tc.expected.String(), ratioDec.String())
		})
	}
}

func TestKeeper_CalcRowanSpotPrice(t *testing.T) {
	testcases := []struct {
		name                          string
		rowanBalance, externalBalance sdk.Uint
		pmtpCurrentRunningRate        sdk.Dec
		expectedSpotPrice             sdk.Dec
		expectedError                 error
	}{
		{
			name:                   "success simple",
			rowanBalance:           sdk.NewUint(1),
			externalBalance:        sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(1),
			expectedSpotPrice:      sdk.MustNewDecFromStr("2"),
		},
		{
			name:                   "success small",
			rowanBalance:           sdk.NewUint(1000000000123),
			externalBalance:        sdk.NewUint(20000000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.4"),
			expectedSpotPrice:      sdk.MustNewDecFromStr("0.000047999999994096"),
		},

		{
			name:                   "success",
			rowanBalance:           sdk.NewUint(1000),
			externalBalance:        sdk.NewUint(2000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.4"),
			expectedSpotPrice:      sdk.MustNewDecFromStr("4.8"),
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pool := types.Pool{
				NativeAssetBalance:   tc.rowanBalance,
				ExternalAssetBalance: tc.externalBalance,
			}

			spotPrice, err := clpkeeper.CalcRowanSpotPrice(&pool, tc.pmtpCurrentRunningRate)
			if tc.expectedError != nil {
				require.EqualError(t, tc.expectedError, err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedSpotPrice, spotPrice)
		})
	}
}

func TestKeeper_CalcRowanValue(t *testing.T) {
	testcases := []struct {
		name                          string
		rowanBalance, externalBalance sdk.Uint
		rowanAmount                   sdk.Uint
		pmtpCurrentRunningRate        sdk.Dec
		expectedValue                 sdk.Uint
		expectedError                 error
	}{
		{
			name:                   "success simple",
			rowanBalance:           sdk.NewUint(1),
			externalBalance:        sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(1),
			rowanAmount:            sdk.NewUint(100),
			expectedValue:          sdk.NewUint(200),
		},
		{
			name:                   "success zero",
			rowanBalance:           sdk.NewUint(1000000000123),
			externalBalance:        sdk.NewUint(20000000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.4"),
			rowanAmount:            sdk.NewUint(100),
			expectedValue:          sdk.NewUint(0),
		},
		{
			name:                   "success",
			rowanBalance:           sdk.NewUint(1000),
			externalBalance:        sdk.NewUint(2000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.4"),
			rowanAmount:            sdk.NewUint(100),
			expectedValue:          sdk.NewUint(480),
		},
		{
			name:                   "fail - rowan balance zero",
			rowanBalance:           sdk.NewUint(0),
			externalBalance:        sdk.NewUint(2000),
			pmtpCurrentRunningRate: sdk.MustNewDecFromStr("1.4"),
			rowanAmount:            sdk.NewUint(100),
			expectedError:          errors.New("amount is invalid"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pool := types.Pool{
				NativeAssetBalance:   tc.rowanBalance,
				ExternalAssetBalance: tc.externalBalance,
			}

			rowanValue, err := clpkeeper.CalcRowanValue(&pool, tc.pmtpCurrentRunningRate, tc.rowanAmount)
			if tc.expectedError != nil {
				require.EqualError(t, tc.expectedError, err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedValue.String(), rowanValue.String())
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
			swapResult:             sdk.NewUint(181),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(817),
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
				ExternalAssetBalance:          sdk.NewUint(1098),
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
				NativeAssetBalance:            sdk.NewUint(1098),
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
				NativeAssetBalance:            sdk.NewUint(1098),
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
				NativeAssetBalance:            sdk.NewUint(1098),
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
				NativeAssetBalance:            sdk.NewUint(1098),
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
				NativeAssetBalance:            sdk.NewUint(1098),
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
				NativeAssetBalance:            sdk.NewUint(1098),
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
				NativeAssetBalance:            sdk.NewUint(1098),
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
			swapResult:             sdk.NewUint(154),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(844),
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
			swapResult:             sdk.NewUint(163),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(835),
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
			swapResult:             sdk.NewUint(172),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(826),
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
			swapResult:             sdk.NewUint(181),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(817),
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
			swapResult:             sdk.NewUint(271),
			liquidityFee:           sdk.NewUint(0),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(727),
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
			swapResult:             sdk.NewUint(362),
			liquidityFee:           sdk.NewUint(1),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(636),
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
			swapResult:             sdk.NewUint(453),
			liquidityFee:           sdk.NewUint(1),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(545),
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
			swapResult:             sdk.NewUint(543),
			liquidityFee:           sdk.NewUint(1),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(455),
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
			swapResult:             sdk.NewUint(634),
			liquidityFee:           sdk.NewUint(1),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(364),
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
			swapResult:             sdk.NewUint(724),
			liquidityFee:           sdk.NewUint(2),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(274),
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
			swapResult:             sdk.NewUint(815),
			liquidityFee:           sdk.NewUint(2),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(183),
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
			swapResult:             sdk.NewUint(906),
			liquidityFee:           sdk.NewUint(2),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(92),
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
			swapResult:             sdk.NewUint(996),
			liquidityFee:           sdk.NewUint(2),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(2),
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
				Health:                        sdk.ZeroDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
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

			marginEnabled := app.ClpKeeper.GetMarginKeeper().IsPoolEnabled(ctx, pool.String())
			swapResult, liquidityFee, priceImpact, newPool, err := clpkeeper.SwapOne(from, swapAmount, to, pool, tc.pmtpCurrentRunningRate, sdk.NewDecWithPrec(3, 3), marginEnabled)

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
