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
	_, err = app.TokenRegistryKeeper.GetEntry(registry, pool.ExternalAsset.Symbol)
	assert.NoError(t, err)
	// asymmetry is positive
	_, _, _, swapAmount := clpkeeper.CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(), wBasis.String(), asymmetry)
	swapResult, liquidityFee, priceImpact, _, err := clpkeeper.SwapOne(types.GetSettlementAsset(), swapAmount, asset, *pool, sdk.OneDec())
	assert.NoError(t, err)
	assert.Equal(t, "19", swapResult.String())
	assert.Equal(t, "978", liquidityFee.String())
	assert.Equal(t, "0", priceImpact.String())
}

func TestKeeper_SwapOneFromGenesis(t *testing.T) {
	const address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"

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
			swapResult:             sdk.NewUint(19),
			liquidityFee:           sdk.NewUint(978),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(100598),
				ExternalAssetBalance:          sdk.NewUint(979),
				PoolUnits:                     sdk.NewUint(1),
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
				ExternalAsset:        &types.Asset{Symbol: "eth"},
				NativeAssetBalance:   sdk.NewUint(1098),
				ExternalAssetBalance: sdk.NewUint(833),
				PoolUnits:            sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(82),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(916),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(90),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(908),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(115),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(883),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(123),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(875),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(148),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(850),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(156),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(842),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(247),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(751),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(495),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(503),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(660),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(338),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(743),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(255),
				PoolUnits:                     sdk.NewUint(998),
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
			swapResult:             sdk.NewUint(908),
			liquidityFee:           sdk.NewUint(8),
			priceImpact:            sdk.ZeroUint(),
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1098),
				ExternalAssetBalance:          sdk.NewUint(90),
				PoolUnits:                     sdk.NewUint(998),
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
					AdminAccounts: test.GetAdmins(tc.address),
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

			require.Equal(t, pool, types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: tc.poolAsset},
				NativeAssetBalance:            tc.nativeAssetAmount,
				ExternalAssetBalance:          tc.externalAssetAmount,
				PoolUnits:                     tc.poolUnits,
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
			swapResult, liquidityFee, priceImpact, newPool, err := clpkeeper.SwapOne(
				from,
				swapAmount,
				to,
				pool,
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
			require.Equal(t, swapResult, tc.swapResult, "swapResult")
			require.Equal(t, liquidityFee, tc.liquidityFee)
			require.Equal(t, priceImpact, tc.priceImpact)
			require.Equal(t, newPool, tc.expectedPool)
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
	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, *pool, sdk.OneDec())
	assert.Equal(t, "2", swapResult.String())
}

func TestKeeper_GetSwapFee_PmtpParams(t *testing.T) {
	pool := types.Pool{
		NativeAssetBalance:   sdk.NewUint(10),
		ExternalAssetBalance: sdk.NewUint(100),
	}
	asset := types.Asset{}

	swapResult := clpkeeper.GetSwapFee(sdk.NewUint(1), asset, pool, sdk.NewDec(100))

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
			errString:            errors.New("Tx amount is too low"),
		},
		{
			name:                 "tx amount too low with no adjustment throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.ZeroUint(),
			externalAssetAmount:  sdk.ZeroUint(),
			errString:            errors.New("Tx amount is too low"),
		},
		{
			name:                 "insufficient native funds throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.ZeroUint(),
			externalAssetAmount:  sdk.OneUint(),
			errString:            errors.New("0: insufficient funds"),
		},
		{
			name:                 "insufficient external funds throws error",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.NewUint(100),
			externalAssetBalance: sdk.ZeroUint(),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.ZeroUint(),
			errString:            errors.New("0: insufficient funds"),
		},
		{
			name:                 "as native asset balance zero then returns native asset amount",
			oldPoolUnits:         sdk.ZeroUint(),
			nativeAssetBalance:   sdk.ZeroUint(),
			externalAssetBalance: sdk.NewUint(100),
			nativeAssetAmount:    sdk.OneUint(),
			externalAssetAmount:  sdk.OneUint(),
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
			poolUnits:            sdk.NewUint(2),
			lpunits:              sdk.NewUint(1),
			panicErr:             "fail to convert 1606938044258990275541962092341162602522202993782792835301376 to cosmos.Dec: decimal out of range; bitLen: got 260, max 256",
		},
	}

	symmetryThreshold := sdk.NewDecWithPrec(1, 4)
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
						symmetryThreshold,
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
				symmetryThreshold,
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
		X, x, Y, fee        sdk.Uint
		err                 error
		errString           error
	}{
		{
			name: "success",
			X:    sdk.NewUint(0),
			x:    sdk.NewUint(0),
			Y:    sdk.NewUint(1),
			fee:  sdk.NewUint(0),
		},
		{
			name: "success",
			X:    sdk.NewUint(1),
			x:    sdk.NewUint(1),
			Y:    sdk.NewUint(1),
			fee:  sdk.NewUint(0),
		},
		{
			name: "success",
			X:    sdk.NewUint(1),
			x:    sdk.NewUint(1),
			Y:    sdk.NewUint(4),
			fee:  sdk.NewUint(1),
		},
		{
			name: "success",
			X:    sdk.NewUint(2),
			x:    sdk.NewUint(2),
			Y:    sdk.NewUint(16),
			fee:  sdk.NewUint(4),
		},
		{
			name: "success",
			X:    sdk.NewUint(1054677676764),
			x:    sdk.NewUint(2567655449999),
			Y:    sdk.NewUint(1099511627776),
			fee:  sdk.NewUint(552454535440),
		},
		{
			name: "success",
			X:    sdk.NewUintFromString("20300000000000000000000000000000000000000000000000000000000000000000000000"),
			x:    sdk.NewUintFromString("10000000000000000658000000000000000000000000000000000000000000000000000000"),
			Y:    sdk.NewUintFromString("10000000000000000000000000000000000000000000000000000000000000000000021344"),
			fee:  sdk.NewUintFromString("1089217832674356640599131638158097447402363655799918705091874559386226334"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fee := clpkeeper.CalcLiquidityFee(tc.X, tc.x, tc.Y)
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
		},
		{
			name:                   "adjust external token without rowan",
			toRowan:                false,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(1),
			y:                      sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(1),
		},
		{
			name:                   "x=0, X=0, Y=0",
			toRowan:                true,
			X:                      sdk.NewUint(0),
			x:                      sdk.NewUint(0),
			Y:                      sdk.NewUint(0),
			y:                      sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(0),
		},
		{
			name:                   "x=1, X=1, Y=1",
			toRowan:                true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(1),
			y:                      sdk.NewUint(0),
			pmtpCurrentRunningRate: sdk.NewDec(0),
		},
		{
			name:                   "x=1, X=1, Y=4",
			toRowan:                true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(4),
			y:                      sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(0),
		},
		{
			name:                   "x=1, X=1, Y=4, nf=10",
			toRowan:                true,
			X:                      sdk.NewUint(1),
			x:                      sdk.NewUint(1),
			Y:                      sdk.NewUint(4),
			y:                      sdk.NewUint(1),
			pmtpCurrentRunningRate: sdk.NewDec(0),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			y := clpkeeper.CalcSwapResult(tc.toRowan, tc.X, tc.x, tc.Y, tc.pmtpCurrentRunningRate)

			require.Equal(t, tc.y.String(), y.String()) // compare strings so that the expected amounts can be read from the failure message
		})
	}
}
