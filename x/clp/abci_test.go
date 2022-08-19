package clp_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEndBlocker(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	_ = test.GeneratePoolsFromFile(app.ClpKeeper, ctx)
	SetRewardParams(app.ClpKeeper, ctx)

	_ = clp.EndBlocker(ctx, app.ClpKeeper)

	pooldash, err := app.ClpKeeper.GetPool(ctx, "cdash")
	assert.NoError(t, err)
	poolceth, err := app.ClpKeeper.GetPool(ctx, "ceth")
	assert.NoError(t, err)
	assert.True(t, poolceth.NativeAssetBalance.GT(pooldash.NativeAssetBalance))

}

func SetRewardParams(keeper keeper.Keeper, ctx sdk.Context) {
	multiplierDec1 := sdk.MustNewDecFromStr("0.5")
	multiplierDec2 := sdk.MustNewDecFromStr("1.5")
	allocations := sdk.NewUintFromString("2000000000000000000")
	keeper.SetRewardParams(ctx, &types.RewardParams{
		LiquidityRemovalLockPeriod:   0,
		LiquidityRemovalCancelPeriod: 2,
		RewardPeriodStartTime:        "",
		RewardPeriods: []*types.RewardPeriod{{
			RewardPeriodId:         "1",
			RewardPeriodStartBlock: 0,
			RewardPeriodEndBlock:   2,
			RewardPeriodAllocation: &allocations,
			RewardPeriodPoolMultipliers: []*types.PoolMultiplier{{
				PoolMultiplierAsset: "cdash",
				Multiplier:          &multiplierDec1,
			},
				{
					PoolMultiplierAsset: "ceth",
					Multiplier:          &multiplierDec2,
				},
			},
		}},
	})
}

func TestBeginBlocker(t *testing.T) {
	testcases := []struct { //nolint
		name                           string
		createBalance                  bool
		createPool                     bool
		createLPs                      bool
		poolAsset                      string
		address                        string
		nativeBalance                  sdk.Int
		externalBalance                sdk.Int
		nativeAssetAmount              sdk.Uint
		externalAssetAmount            sdk.Uint
		poolUnits                      sdk.Uint
		poolAssetPermissions           []tokenregistrytypes.Permission
		nativeAssetPermissions         []tokenregistrytypes.Permission
		params                         types.Params
		epoch                          types.PmtpEpoch
		err                            error
		errString                      error
		panicErr                       string
		maxRowanLiquidityThreshold     sdk.Uint
		currentRowanLiquidityThreshold sdk.Uint
		liquidityProtectionEpochLength uint64
		liquidityProtectionIsActive    bool
		expectedRunningThresholdEnd    sdk.Uint
	}{
		{
			name:                   "current height equals to start block",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			params: types.Params{
				MinCreatePoolThreshold: types.DefaultMinCreatePoolThreshold,
			},
			epoch: types.PmtpEpoch{
				EpochCounter: 0,
				BlockCounter: 0,
			},
		},
		{
			name:                   "current height equals or greater than start block",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			params: types.Params{
				MinCreatePoolThreshold: types.DefaultMinCreatePoolThreshold,
			},
			epoch: types.PmtpEpoch{
				EpochCounter: 10,
				BlockCounter: 0,
			},
		},
		{
			name:                   "last block counter set to zero",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			params: types.Params{
				MinCreatePoolThreshold: types.DefaultMinCreatePoolThreshold,
			},
			epoch: types.PmtpEpoch{
				EpochCounter: 10,
				BlockCounter: 0,
			},
		},
		{
			name:                   "throws panic error",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:      sdk.NewUint(0),
			externalAssetAmount:    sdk.NewUint(0),
			poolUnits:              sdk.NewUint(0),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			params: types.Params{
				MinCreatePoolThreshold: types.DefaultMinCreatePoolThreshold,
			},
			epoch: types.PmtpEpoch{
				EpochCounter: 10,
				BlockCounter: 10,
			},
			// panicErr: "not enough received asset tokens to swap",
		},
		{
			name:                   "liquidity protection correct replenishment",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			params: types.Params{
				MinCreatePoolThreshold: types.DefaultMinCreatePoolThreshold,
			},
			epoch: types.PmtpEpoch{
				EpochCounter: 10,
				BlockCounter: 0,
			},
			liquidityProtectionIsActive:    true,
			maxRowanLiquidityThreshold:     sdk.NewUint(100),
			currentRowanLiquidityThreshold: sdk.NewUint(80),
			liquidityProtectionEpochLength: 10,
			expectedRunningThresholdEnd:    sdk.NewUint(90),
		},
		{
			name:                   "liquidity protection correct replenishment hit maximum",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			params: types.Params{
				MinCreatePoolThreshold: types.DefaultMinCreatePoolThreshold,
			},
			epoch: types.PmtpEpoch{
				EpochCounter: 10,
				BlockCounter: 0,
			},
			liquidityProtectionIsActive:    true,
			maxRowanLiquidityThreshold:     sdk.NewUint(100),
			currentRowanLiquidityThreshold: sdk.NewUint(95),
			liquidityProtectionEpochLength: 10,
			expectedRunningThresholdEnd:    sdk.NewUint(100),
		},
		{
			name:                   "liquidity protection maximum max liquidity threshold",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			params: types.Params{
				MinCreatePoolThreshold: types.DefaultMinCreatePoolThreshold,
			},
			epoch: types.PmtpEpoch{
				EpochCounter: 10,
				BlockCounter: 0,
			},
			liquidityProtectionIsActive:    true,
			maxRowanLiquidityThreshold:     sdk.NewUintFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
			currentRowanLiquidityThreshold: sdk.NewUintFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
			liquidityProtectionEpochLength: 10,
			expectedRunningThresholdEnd:    sdk.NewUintFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				trGs := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: tc.poolAsset, BaseDenom: tc.poolAsset, Decimals: 18, Permissions: tc.poolAssetPermissions},
							{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: tc.nativeAssetPermissions},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(trGs)
				genesisState["tokenregistry"] = bz

				if tc.createBalance {
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
				}

				if tc.createPool {
					pools := []*types.Pool{
						{
							ExternalAsset:        &types.Asset{Symbol: tc.poolAsset},
							NativeAssetBalance:   tc.nativeAssetAmount,
							ExternalAssetBalance: tc.externalAssetAmount,
							PoolUnits:            tc.poolUnits,
						},
					}
					clpGs := types.DefaultGenesisState()
					if tc.createLPs {
						lps := []*types.LiquidityProvider{
							{
								Asset:                    &types.Asset{Symbol: tc.poolAsset},
								LiquidityProviderAddress: tc.address,
								LiquidityProviderUnits:   tc.nativeAssetAmount,
							},
						}
						clpGs.LiquidityProviders = append(clpGs.LiquidityProviders, lps...)
					}
					clpGs.Params = tc.params
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			app.ClpKeeper.SetParams(ctx, tc.params)
			app.ClpKeeper.SetPmtpRateParams(ctx, types.PmtpRateParams{
				PmtpPeriodBlockRate:    sdk.OneDec(),
				PmtpCurrentRunningRate: sdk.OneDec(),
			})
			app.ClpKeeper.SetPmtpEpoch(ctx, tc.epoch)

			liquidityProtectionParam := app.ClpKeeper.GetLiquidityProtectionParams(ctx)
			liquidityProtectionParam.MaxRowanLiquidityThreshold = tc.maxRowanLiquidityThreshold
			liquidityProtectionParam.IsActive = tc.liquidityProtectionIsActive
			liquidityProtectionParam.EpochLength = tc.liquidityProtectionEpochLength
			app.ClpKeeper.SetLiquidityProtectionParams(ctx, liquidityProtectionParam)
			app.ClpKeeper.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, tc.currentRowanLiquidityThreshold)

			if tc.panicErr != "" {
				// nolint:errcheck
				require.PanicsWithError(t, tc.panicErr, func() {
					clp.BeginBlocker(ctx, app.ClpKeeper)
				})
				return
			}

			clp.BeginBlocker(ctx, app.ClpKeeper)

			if tc.liquidityProtectionIsActive {
				runningThreshold := app.ClpKeeper.GetLiquidityProtectionRateParams(ctx).CurrentRowanLiquidityThreshold
				require.Equal(t, tc.expectedRunningThresholdEnd.String(), runningThreshold.String())
			}
		})
	}
}
