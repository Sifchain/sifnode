package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func Test_PmtpFloatCalculations(t *testing.T) {
	pmtpPeriodGovernanceRate := sdk.MustNewDecFromStr("667577.111234534525628462")
	numEpochsInPolicyPeriod := 10
	numBlocksInPolicyPeriod := 100
	pmtpPeriodBlockRate := (sdk.NewDec(1).Add(pmtpPeriodGovernanceRate))
	pow := float64(numEpochsInPolicyPeriod) / float64(numBlocksInPolicyPeriod)
	require.Equal(t, pow, float64(0.1))
	require.Equal(t, pmtpPeriodBlockRate, sdk.MustNewDecFromStr("667578.111234534525628462"))
	s := 1.232322323223435445
	//f:=strconv.FormatFloat(s, 'E', -1, 64)
	f := fmt.Sprintf("%v", s)
	bpow := sdk.MustNewDecFromStr(f).BigInt()
	//bpow := sdk.MustNewDecFromStr(fmt.Sprintf("%f", 1.232322323223435445)).BigInt()
	//bpow := big.NewFloat(1.232322323223435445)
	bbr := pmtpPeriodBlockRate.BigInt()
	fbr := new(big.Float).SetInt(bbr)
	fpow := new(big.Float).SetInt(bpow)
	fbr = fbr.Quo(fbr, big.NewFloat(1000000000000000000.00))
	fpow = fpow.Quo(fpow, big.NewFloat(1000000000000000000.00))
	//
	//value := math.Pow(bint,bpow)
	//dec := sdk.MustNewDecFromStr(fmt.Sprintf("%f", value)).MustFloat64(
	require.Equal(t, fbr.String(), "667578.1112")
	require.Equal(t, fpow.String(), "1.232322323")

	fss := fmt.Sprintf("%.18f", 0.00000013751833967123872)
	require.Equal(t, fss, "0.000000137518339671")
}
func TestKeeper_PolicyRun(t *testing.T) {
	SwapPriceNative := sdk.MustNewDecFromStr("1.996005992010000000")
	SwapPriceExternal := sdk.MustNewDecFromStr("0.499001498002500000")

	testcases := []struct {
		name                   string
		createBalance          bool
		createPool             bool
		createLPs              bool
		poolAsset              string
		address                string
		nativeBalance          sdk.Int
		externalBalance        sdk.Int
		nativeAssetAmount      sdk.Uint
		externalAssetAmount    sdk.Uint
		poolUnits              sdk.Uint
		poolAssetPermissions   []tokenregistrytypes.Permission
		nativeAssetPermissions []tokenregistrytypes.Permission
		expectedPool           types.Pool
		err                    error
		errString              error
	}{
		{
			name:                   "default",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			expectedPool: types.Pool{
				ExternalAsset:        &types.Asset{Symbol: "eth"},
				NativeAssetBalance:   sdk.NewUint(1000),
				ExternalAssetBalance: sdk.NewUint(1000),
				PoolUnits:            sdk.NewUint(1000),
				SwapPriceNative:      &SwapPriceNative,
				SwapPriceExternal:    &SwapPriceExternal,
			},
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
					clpGs.Params = types.Params{
						MinCreatePoolThreshold:   100,
						PmtpPeriodGovernanceRate: sdk.OneDec(),
						PmtpPeriodEpochLength:    1,
						PmtpPeriodStartBlock:     1,
						PmtpPeriodEndBlock:       2,
					}
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			app.ClpKeeper.SetPmtpCurrentRunningRate(ctx, sdk.NewDec(1))

			err := app.ClpKeeper.PolicyRun(ctx, sdk.NewDec(1))

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			pool, _ := app.ClpKeeper.GetPool(ctx, tc.poolAsset)

			require.Equal(t, pool, tc.expectedPool)
		})
	}
}
