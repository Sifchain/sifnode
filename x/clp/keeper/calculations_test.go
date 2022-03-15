package keeper_test

import (
	"errors"
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
	address := "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"

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
	ctx, app := test.CreateTestAppClp(false)
	signer := test.GenerateAddress(test.AddressKey1)
	//Parameters for create pool
	nativeAssetAmount := sdk.NewUintFromString("998")
	externalAssetAmount := sdk.NewUintFromString("998")
	asset := types.NewAsset("eth")
	externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(10000)))
	nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(sdk.NewUint(10000)))
	wBasis := sdk.NewInt(1000)
	asymmetry := sdk.NewInt(10000)
	err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
	require.NoError(t, err)
	msgCreatePool := types.NewMsgCreatePool(signer, asset, nativeAssetAmount, externalAssetAmount)
	// Create Pool
	pool, err := app.ClpKeeper.CreatePool(ctx, sdk.NewUint(1), &msgCreatePool)
	assert.NoError(t, err)
	msg := types.NewMsgAddLiquidity(signer, asset, nativeAssetAmount, externalAssetAmount)
	app.ClpKeeper.CreateLiquidityProvider(ctx, &asset, sdk.NewUint(1), signer)
	lp, err := app.ClpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.NewUint(1), sdk.NewUint(998))
	assert.NoError(t, err)
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := app.TokenRegistryKeeper.GetEntry(registry, pool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	// asymmetry is positive
	normalizationFactor, adjustExternalToken := app.ClpKeeper.GetNormalizationFactor(eAsset.Decimals)
	_, _, _, swapAmount := clpkeeper.CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(), wBasis.String(), asymmetry)
	swapResult, liquidityFee, priceImpact, _, err := clpkeeper.SwapOne(types.GetSettlementAsset(), swapAmount, asset, *pool, normalizationFactor, adjustExternalToken, sdk.OneDec())
	assert.NoError(t, err)
	assert.Equal(t, swapResult.String(), "20")
	assert.Equal(t, liquidityFee.String(), "978")
	assert.Equal(t, priceImpact.String(), "0")
}

func TestKeeper_SwapOneFromGenesis(t *testing.T) {
	testcases := []struct {
		name                string
		poolAsset           string
		address             string
		nativeBalance       sdk.Int
		externalBalance     sdk.Int
		nativeAssetAmount   sdk.Uint
		externalAssetAmount sdk.Uint
		poolUnits           sdk.Uint
		calculateWithdraw   bool
		normalizationFactor sdk.Dec
		adjustExternalToken bool
		swapAmount          sdk.Uint
		wBasis              sdk.Int
		asymmetry           sdk.Int
		swapResult          sdk.Uint
		liquidityFee        sdk.Uint
		priceImpact         sdk.Uint
		expectedPool        types.Pool
		err                 error
		errString           error
	}{
		{
			name:                "swap with single pool units",
			poolAsset:           "eth",
			address:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:       sdk.NewInt(10000),
			externalBalance:     sdk.NewInt(10000),
			nativeAssetAmount:   sdk.NewUint(998),
			externalAssetAmount: sdk.NewUint(998),
			poolUnits:           sdk.OneUint(),
			calculateWithdraw:   true,
			wBasis:              sdk.NewInt(1000),
			asymmetry:           sdk.NewInt(10000),
			swapResult:          sdk.NewUint(20),
			liquidityFee:        sdk.NewUint(978),
			priceImpact:         sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:        &types.Asset{Symbol: "eth"},
				NativeAssetBalance:   sdk.NewUint(100598),
				ExternalAssetBalance: sdk.NewUint(978),
				PoolUnits:            sdk.NewUint(1),
			},
		},
		{
			name:                "swap with equal amount of pool units",
			poolAsset:           "eth",
			address:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:       sdk.NewInt(10000),
			externalBalance:     sdk.NewInt(10000),
			nativeAssetAmount:   sdk.NewUint(998),
			externalAssetAmount: sdk.NewUint(998),
			poolUnits:           sdk.NewUint(998),
			calculateWithdraw:   true,
			wBasis:              sdk.NewInt(1000),
			asymmetry:           sdk.NewInt(10000),
			swapResult:          sdk.NewUint(165),
			liquidityFee:        sdk.NewUint(8),
			priceImpact:         sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:        &types.Asset{Symbol: "eth"},
				NativeAssetBalance:   sdk.NewUint(1098),
				ExternalAssetBalance: sdk.NewUint(833),
				PoolUnits:            sdk.NewUint(998),
			},
		},
		{
			name:                "swap with empty pool",
			poolAsset:           "eth",
			address:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:       sdk.NewInt(10000),
			externalBalance:     sdk.NewInt(10000),
			nativeAssetAmount:   sdk.NewUint(0),
			externalAssetAmount: sdk.NewUint(0),
			poolUnits:           sdk.NewUint(0),
			calculateWithdraw:   false,
			normalizationFactor: sdk.NewDec(0),
			adjustExternalToken: true,
			swapAmount:          sdk.NewUint(0),
			swapResult:          sdk.NewUint(165),
			liquidityFee:        sdk.NewUint(8),
			priceImpact:         sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:        &types.Asset{Symbol: "eth"},
				NativeAssetBalance:   sdk.NewUint(1098),
				ExternalAssetBalance: sdk.NewUint(833),
				PoolUnits:            sdk.NewUint(998),
			},
			errString: errors.New("not enough received asset tokens to swap"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				trGs := &tokenregistrytypes.GenesisState{
					AdminAccount: tc.address,
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
						ExternalAsset:        &types.Asset{Symbol: tc.poolAsset},
						NativeAssetBalance:   tc.nativeAssetAmount,
						ExternalAssetBalance: tc.externalAssetAmount,
						PoolUnits:            tc.poolUnits,
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
					MinCreatePoolThreshold:   100,
					PmtpPeriodGovernanceRate: sdk.OneDec(),
					PmtpPeriodEpochLength:    1,
					PmtpPeriodStartBlock:     1,
					PmtpPeriodEndBlock:       2,
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

			require.Equal(t, pool, types.Pool{
				ExternalAsset:        &types.Asset{Symbol: tc.poolAsset},
				NativeAssetBalance:   tc.nativeAssetAmount,
				ExternalAssetBalance: tc.externalAssetAmount,
				PoolUnits:            tc.poolUnits,
			})

			var normalizationFactor sdk.Dec
			var adjustExternalToken bool
			var swapAmount sdk.Uint

			if tc.calculateWithdraw {
				normalizationFactor, adjustExternalToken = app.ClpKeeper.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
				_, _, _, swapAmount = clpkeeper.CalculateWithdrawal(
					pool.PoolUnits,
					pool.NativeAssetBalance.String(),
					pool.ExternalAssetBalance.String(),
					lp.LiquidityProviderUnits.String(),
					tc.wBasis.String(),
					tc.asymmetry,
				)
			} else {
				normalizationFactor = tc.normalizationFactor
				adjustExternalToken = tc.adjustExternalToken
				swapAmount = tc.swapAmount
			}

			swapResult, liquidityFee, priceImpact, newPool, err := clpkeeper.SwapOne(
				types.GetSettlementAsset(),
				swapAmount,
				types.Asset{Symbol: tc.poolAsset},
				pool,
				normalizationFactor,
				adjustExternalToken,
				sdk.OneDec(),
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
			require.Equal(t, swapResult, tc.swapResult)
			require.Equal(t, liquidityFee, tc.liquidityFee)
			require.Equal(t, priceImpact, tc.priceImpact)
			require.Equal(t, newPool, tc.expectedPool)
		})
	}
}

func TestKeeper_SwapOneNew(t *testing.T) {
	testcases :=
		[]struct {
			name                  string
			createToken           bool
			denom                 string
			decimals              int64
			asset                 types.Asset
			fundAccount           bool
			nativeBalance         sdk.Uint
			externalBalance       sdk.Uint
			createPool            bool
			nativeAssetAmount     sdk.Uint
			externalAssetAmount   sdk.Uint
			addLiquidity          bool
			wBasis                sdk.Int
			asymmetry             sdk.Int
			swapResult            sdk.Uint
			liquidityFee          sdk.Uint
			priceImpact           sdk.Uint
			expPanic              bool
			expPanicMsg           string
			createPoolErr         error
			createPoolErrString   error
			addLiquidityErr       error
			addLiquidityErrString error
			err                   error
			errString             error
		}{
			{
				name:                "missing token throws error",
				createToken:         false,
				denom:               "xxx",
				decimals:            18,
				asset:               types.Asset{Symbol: "xxx"},
				fundAccount:         true,
				nativeBalance:       sdk.NewUint(10000),
				externalBalance:     sdk.NewUint(10000),
				createPool:          true,
				nativeAssetAmount:   sdk.NewUint(998),
				externalAssetAmount: sdk.NewUint(998),
				addLiquidity:        true,
				wBasis:              sdk.NewInt(1000),
				asymmetry:           sdk.NewInt(10000),
				swapResult:          sdk.NewUint(20),
				liquidityFee:        sdk.NewUint(978),
				priceImpact:         sdk.NewUint(0),
				expPanic:            true,
			},
			{
				name:                "successful swap one",
				createToken:         true,
				denom:               "xxx",
				decimals:            18,
				asset:               types.Asset{Symbol: "eth"},
				fundAccount:         true,
				nativeBalance:       sdk.NewUint(10000),
				externalBalance:     sdk.NewUint(10000),
				createPool:          true,
				nativeAssetAmount:   sdk.NewUint(998),
				externalAssetAmount: sdk.NewUint(998),
				addLiquidity:        true,
				wBasis:              sdk.NewInt(1000),
				asymmetry:           sdk.NewInt(10000),
				swapResult:          sdk.NewUint(20),
				liquidityFee:        sdk.NewUint(978),
				priceImpact:         sdk.NewUint(0),
			},
		}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClp(false)
			clpKeeper := app.ClpKeeper
			signer := test.GenerateAddress(test.AddressKey1)

			var err error

			if tc.createToken {
				app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
					Denom:       tc.denom,
					Decimals:    tc.decimals,
					Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
				})
			}

			if tc.fundAccount {
				externalCoin := sdk.NewCoin(tc.asset.Symbol, sdk.Int(tc.externalBalance))
				nativeCoin := sdk.NewCoin(types.NativeSymbol, sdk.Int(tc.nativeBalance))
				err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(externalCoin, nativeCoin))
				require.NoError(t, err)
			}

			var pool *types.Pool

			if tc.createPool {
				msgCreatePool := types.NewMsgCreatePool(signer, tc.asset, tc.nativeAssetAmount, tc.externalAssetAmount)
				pool, err = clpKeeper.CreatePool(ctx, sdk.OneUint(), &msgCreatePool)

				if tc.createPoolErrString != nil {
					require.EqualError(t, err, tc.createPoolErrString.Error())
				} else if tc.createPoolErr == nil {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, tc.createPoolErr)
				}
			}

			var lp *types.LiquidityProvider

			if tc.addLiquidity {
				msg := types.NewMsgAddLiquidity(signer, tc.asset, tc.nativeAssetAmount, tc.externalAssetAmount)
				clpKeeper.CreateLiquidityProvider(ctx, &tc.asset, sdk.OneUint(), signer)
				lp, err = clpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.OneUint(), sdk.NewUint(998))

				if tc.addLiquidityErrString != nil {
					require.EqualError(t, err, tc.addLiquidityErrString.Error())
				} else if tc.err == nil {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, tc.addLiquidityErr)
				}
			}

			normalizationFactor, adjustExternalToken := app.ClpKeeper.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
			_, _, _, swapAmount := clpkeeper.CalculateWithdrawal(pool.PoolUnits,
				pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(), tc.wBasis.String(), tc.asymmetry)
			if tc.expPanic {
				if tc.expPanicMsg != "" {
					require.PanicsWithValue(t, tc.expPanicMsg, func() {
						clpkeeper.SwapOne(types.GetSettlementAsset(), swapAmount, tc.asset, *pool, normalizationFactor, adjustExternalToken, sdk.OneDec())
					})
				} else {
					require.Panics(t, func() {
						clpkeeper.SwapOne(types.GetSettlementAsset(), swapAmount, tc.asset, *pool, normalizationFactor, adjustExternalToken, sdk.OneDec())
					})
				}
			} else {
				swapResult, liquidityFee, priceImpact, _, err := clpkeeper.SwapOne(types.GetSettlementAsset(), swapAmount, tc.asset, *pool, normalizationFactor, adjustExternalToken, sdk.OneDec())
				require.NoError(t, err)
				require.Equal(t, swapResult, tc.swapResult)
				require.Equal(t, liquidityFee, tc.liquidityFee)
				require.Equal(t, priceImpact, tc.priceImpact)

				if tc.errString != nil {
					require.EqualError(t, err, tc.errString.Error())
				} else if tc.err == nil {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, tc.err)
				}
			}
		})
	}
}

func TestKeeper_SetInputs(t *testing.T) {
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
	X, x, Y, toRowan := clpkeeper.SetInputs(sdk.NewUint(1), asset, *pool)
	assert.Equal(t, X, sdk.NewUint(998))
	assert.Equal(t, x, sdk.NewUint(1))
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
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	eAsset, _ := app.TokenRegistryKeeper.GetEntry(registry, pool.ExternalAsset.Symbol)
	normalizationFactor, adjustExternalToken := app.ClpKeeper.GetNormalizationFactor(eAsset.Decimals)
	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, *pool, normalizationFactor, adjustExternalToken, sdk.OneDec())
	assert.Equal(t, "2", swapResult.String())
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
