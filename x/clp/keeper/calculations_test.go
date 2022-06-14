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
			name:                   "successful swap with single pool units",
			poolAsset:              "eth",
			address:                address,
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(998),
			externalAssetAmount:    sdk.NewUint(998),
			poolUnits:              sdk.OneUint(),
			calculateWithdraw:      true,
			wBasis:                 sdk.NewInt(1000),
			asymmetry:              sdk.NewInt(10000),
			pmtpCurrentRunningRate: sdk.OneDec(),
			swapResult:             sdk.NewUint(20),
			liquidityFee:           sdk.NewUint(978),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(100598),
				ExternalAssetBalance:          sdk.NewUint(978),
				PoolUnits:                     sdk.NewUint(1),
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
			swapResult:             sdk.NewUint(165),
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
			swapResult:             sdk.NewUint(165),
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
			swapResult:             sdk.NewUint(41),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(957),
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
			swapResult:             sdk.NewUint(83),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(915),
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
			swapResult:             sdk.NewUint(91),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(907),
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
			swapResult:             sdk.NewUint(99),
			liquidityFee:           sdk.NewUint(8),
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
			swapResult:             sdk.NewUint(107),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(891),
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
			swapResult:             sdk.NewUint(116),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(882),
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
			swapResult:             sdk.NewUint(124),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(874),
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
			swapResult:             sdk.NewUint(132),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(866),
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
			swapResult:             sdk.NewUint(140),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(858),
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
			swapResult:             sdk.NewUint(149),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(849),
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
			swapResult:             sdk.NewUint(157),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(841),
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
			swapResult:             sdk.NewUint(165),
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
			swapResult:             sdk.NewUint(248),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(750),
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
			swapResult:             sdk.NewUint(330),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(668),
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
			swapResult:             sdk.NewUint(413),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(585),
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
			swapResult:             sdk.NewUint(496),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(502),
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
			swapResult:             sdk.NewUint(578),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(420),
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
			swapResult:             sdk.NewUint(661),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(337),
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
			swapResult:             sdk.NewUint(744),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(254),
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
			swapResult:             sdk.NewUint(826),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(172),
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
			swapResult:             sdk.NewUint(909),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(89),
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

			var normalizationFactor sdk.Dec
			var adjustExternalToken bool
			var swapAmount sdk.Uint

			if tc.calculateWithdraw {
				normalizationFactor, adjustExternalToken, _ = app.ClpKeeper.GetNormalizationFactorFromAsset(ctx, *pool.ExternalAsset)
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

			from := tc.from
			if from == (types.Asset{}) {
				from = types.GetSettlementAsset()
			}
			to := tc.to
			if to == (types.Asset{}) {
				to = types.Asset{Symbol: tc.poolAsset}
			}
			swapResult, liquidityFee, priceImpact, newPool, err := clpkeeper.SwapOne(
				from,
				swapAmount,
				to,
				pool,
				normalizationFactor,
				adjustExternalToken,
				tc.pmtpCurrentRunningRate,
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

func TestKeeper_GetSwapFee_PmtpParams(t *testing.T) {
	pool := types.Pool{
		NativeAssetBalance:   sdk.NewUint(10),
		ExternalAssetBalance: sdk.NewUint(100),
	}
	asset := types.Asset{}
	normalizationFactor := sdk.NewDec(1)
	adjustExternalToken := false

	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, pool, normalizationFactor, adjustExternalToken, sdk.NewDec(100))

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
		normalizationFactor  sdk.Dec
		adjustExternalToken  bool
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
			normalizationFactor:  sdk.ZeroDec(),
			adjustExternalToken:  true,
			errString:            errors.New("Tx amount is too low"),
		},
		{
			name:                 "tx amount too low with no adjustment throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.ZeroUint(),
			externalAssetAmount:  sdk.ZeroUint(),
			normalizationFactor:  sdk.ZeroDec(),
			adjustExternalToken:  false,
			errString:            errors.New("Tx amount is too low"),
		},
		{
			name:                 "insufficient native funds throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			normalizationFactor:  sdk.ZeroDec(),
			adjustExternalToken:  false,
			errString:            errors.New("0: insufficient funds"),
		},
		{
			name:                 "insufficient external funds throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.ZeroUint(),
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
			errString:            errors.New("0: insufficient funds"),
		},
		{
			name:                 "as native asset balance zero then returns native asset amount",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
			poolUnits:            sdk.OneUint(),
			lpunits:              sdk.OneUint(),
		},
		{
			name:                 "fail to convert oldPoolUnits to Dec",
			oldPoolUnits:         sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000000000"),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
			poolUnits:            sdk.OneUint(),
			lpunits:              sdk.OneUint(),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "fail to convert nativeAssetBalance to Dec",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000000000"),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
			poolUnits:            sdk.OneUint(),
			lpunits:              sdk.OneUint(),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "fail to convert externalAssetBalance to Dec",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000000000"),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
			poolUnits:            sdk.OneUint(),
			lpunits:              sdk.OneUint(),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "fail to convert nativeAssetAmount to Dec",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000000000"),
			externalAssetAmount:  sdk.OneUint(),
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
			poolUnits:            sdk.OneUint(),
			lpunits:              sdk.OneUint(),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "fail to convert externalAssetAmount to Dec",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000000000"),
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
			poolUnits:            sdk.OneUint(),
			lpunits:              sdk.OneUint(),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "successful",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
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
			normalizationFactor:  sdk.OneDec(),
			adjustExternalToken:  false,
			poolUnits:            sdk.ZeroUint(),
			lpunits:              sdk.ZeroUint(),
			errString:            errors.New("Cannot add liquidity asymmetrically"),
		},
	}

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
						tc.normalizationFactor,
						tc.adjustExternalToken,
						sdk.NewDecWithPrec(1, 4),
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
				tc.normalizationFactor,
				tc.adjustExternalToken,
				sdk.NewDecWithPrec(1, 4),
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
			require.Equal(t, poolUnits, tc.poolUnits)
			require.Equal(t, lpunits, tc.lpunits)
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
			nativeAssetBalance:   "10000000000000000000000000000000000000000000000000000000000000000000000000",
			externalAssetBalance: "1",
			lpUnits:              "1",
			wBasisPoints:         "1",
			asymmetry:            sdk.NewInt(1),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "fail to convert externalAssetBalance to Dec",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "1",
			externalAssetBalance: "10000000000000000000000000000000000000000000000000000000000000000000000000",
			lpUnits:              "1",
			wBasisPoints:         "1",
			asymmetry:            sdk.NewInt(1),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "fail to convert lpUnits to Dec",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "1",
			externalAssetBalance: "1",
			lpUnits:              "10000000000000000000000000000000000000000000000000000000000000000000000000",
			wBasisPoints:         "1",
			asymmetry:            sdk.NewInt(1),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "fail to convert wBasisPoints to Dec",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "1",
			externalAssetBalance: "1",
			lpUnits:              "1",
			wBasisPoints:         "10000000000000000000000000000000000000000000000000000000000000000000000000",
			asymmetry:            sdk.NewInt(1),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
		{
			name:                 "fail to convert asymmetry to Dec",
			poolUnits:            sdk.NewUint(1),
			nativeAssetBalance:   "1",
			externalAssetBalance: "1",
			lpUnits:              "1",
			wBasisPoints:         "1",
			asymmetry:            sdk.Int(sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000000000")),
			panicErr:             "fail to convert 10000000000000000000000000000000000000000000000000000000000000000000000000 to cosmos.Dec: decimal out of range; bitLen: got 303, max 256",
		},
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
		X, x, Y             sdk.Uint
		err                 error
		errString           error
	}{
		{
			name:                "Y zero",
			toRowan:             true,
			normalizationFactor: sdk.NewDec(1),
			adjustExternalToken: true,
			X:                   sdk.NewUint(1),
			x:                   sdk.NewUint(1),
			Y:                   sdk.NewUint(0),
		},
		{
			name:                "adjust external token with rowan",
			toRowan:             true,
			normalizationFactor: sdk.NewDec(1),
			adjustExternalToken: true,
			X:                   sdk.NewUint(1),
			x:                   sdk.NewUint(1),
			Y:                   sdk.NewUint(1),
		},
		{
			name:                "adjust external token without rowan",
			toRowan:             false,
			normalizationFactor: sdk.NewDec(1),
			adjustExternalToken: true,
			X:                   sdk.NewUint(1),
			x:                   sdk.NewUint(1),
			Y:                   sdk.NewUint(1),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := clpkeeper.CalcLiquidityFee(tc.toRowan, tc.normalizationFactor, tc.adjustExternalToken, tc.X, tc.x, tc.Y)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestKeeper_CalcSwapResult(t *testing.T) {
	testcases := []struct {
		name                   string
		toRowan                bool
		adjustExternalToken    bool
		X, x, Y                sdk.Uint
		normalizationFactor    sdk.Dec
		pmtpCurrentRunningRate sdk.Dec
		err                    error
		errString              error
	}{
		{
			name:                   "adjust external token with rowan",
			toRowan:                true,
			normalizationFactor:    sdk.NewDec(1),
			adjustExternalToken:    true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(1),
		},
		{
			name:                   "adjust external token without rowan",
			toRowan:                false,
			normalizationFactor:    sdk.NewDec(1),
			adjustExternalToken:    true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(1),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := clpkeeper.CalcSwapResult(tc.toRowan, tc.normalizationFactor, tc.adjustExternalToken, tc.X, tc.x, tc.Y, tc.pmtpCurrentRunningRate)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
