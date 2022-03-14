package keeper_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func TestKeeper_SwapOne(t *testing.T) {
	testcases :=
		[]struct {
			name                string
			createToken         bool
			denom               string
			decimals            int64
			asset               types.Asset
			fundAccount         bool
			nativeBalance       sdk.Uint
			externalBalance     sdk.Uint
			createPool          bool
			nativeAssetAmount   sdk.Uint
			externalAssetAmount sdk.Uint
			addLiquidity        bool
			wBasis              sdk.Int
			asymmetry           sdk.Int
			swapResult          sdk.Uint
			liquidityFee        sdk.Uint
			priceImpact         sdk.Uint
			expPanic            bool
			expPanicMsg         string
			err                 error
			errString           error
		}{
			// {
			// 	name:                "token missing throws error",
			// 	createToken:         false,
			// 	denom:               "xxx",
			// 	decimals:            18,
			// 	asset:               types.Asset{Symbol: "xxx"},
			// 	fundAccount:         true,
			// 	nativeBalance:       sdk.NewUint(10000),
			// 	externalBalance:     sdk.NewUint(10000),
			// 	createPool:          true,
			// 	nativeAssetAmount:   sdk.NewUint(998),
			// 	externalAssetAmount: sdk.NewUint(998),
			// 	addLiquidity:        true,
			// 	wBasis:              sdk.NewInt(1000),
			// 	asymmetry:           sdk.NewInt(10000),
			// 	swapResult:          sdk.NewUint(20),
			// 	liquidityFee:        sdk.NewUint(978),
			// 	priceImpact:         sdk.NewUint(0),
			// 	expPanic:            true,
			// 	expPanicMsg:         "invalid memory address or nil pointer dereference",
			// },
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
				require.NoError(t, err)
			}

			var lp *types.LiquidityProvider

			if tc.addLiquidity {
				msg := types.NewMsgAddLiquidity(signer, tc.asset, tc.nativeAssetAmount, tc.externalAssetAmount)
				clpKeeper.CreateLiquidityProvider(ctx, &tc.asset, sdk.OneUint(), signer)
				lp, err = clpKeeper.AddLiquidity(ctx, &msg, *pool, sdk.OneUint(), sdk.NewUint(998))
				require.NoError(t, err)
			}

			normalizationFactor, adjustExternalToken := app.ClpKeeper.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
			_, _, _, swapAmount := clpkeeper.CalculateWithdrawal(pool.PoolUnits,
				pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(), tc.wBasis.String(), tc.asymmetry)
			if tc.expPanic {
				require.PanicsWithValue(t, tc.expPanicMsg, func() {
					clpkeeper.SwapOne(types.GetSettlementAsset(), swapAmount, tc.asset, *pool, normalizationFactor, adjustExternalToken, sdk.OneDec())
				})
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
