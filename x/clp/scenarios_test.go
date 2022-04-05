package clp_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestScenarios(t *testing.T) {
	type ExpectedStates []struct {
		Height            int64                `json:"height,omitempty"`
		Pool              types.Pool           `json:"pool,omitempty"`
		SwapPriceNative   sdk.Dec              `json:"swap_price_native,omitempty"`
		SwapPriceExternal sdk.Dec              `json:"swap_price_external,omitempty"`
		PmtpRateParams    types.PmtpRateParams `json:"pmtp_rate_params,omitempty"`
	}

	type Scenarios []struct {
		Name                   string                          `json:"name,omitempty"`
		CreateBalance          bool                            `json:"create_balance,omitempty"`
		CreatePool             bool                            `json:"create_pool,omitempty"`
		CreateLPs              bool                            `json:"create_lps,omitempty"`
		PoolAsset              string                          `json:"pool_asset,omitempty"`
		Address                string                          `json:"address,omitempty"`
		NativeBalance          sdk.Int                         `json:"native_balance,omitempty"`
		ExternalBalance        sdk.Int                         `json:"external_balance,omitempty"`
		NativeAssetAmount      sdk.Uint                        `json:"native_asset_amount,omitempty"`
		ExternalAssetAmount    sdk.Uint                        `json:"external_asset_amount,omitempty"`
		PoolUnits              sdk.Uint                        `json:"pool_units,omitempty"`
		PoolAssetDecimals      int64                           `json:"pool_asset_decimals,omitempty"`
		PoolAssetPermissions   []tokenregistrytypes.Permission `json:"pool_asset_permissions,omitempty"`
		NativeAssetPermissions []tokenregistrytypes.Permission `json:"native_asset_permissions,omitempty"`
		Params                 types.Params                    `json:"params,omitempty"`
		Epoch                  types.PmtpEpoch                 `json:"epoch,omitempty"`
		ExpectedStates         ExpectedStates                  `json:"expected_states,omitempty"`
		Err                    error                           `json:"err,omitempty"`
		ErrString              error                           `json:"err_string,omitempty"`
	}

	file, err := ioutil.ReadFile("scenarios.json")
	require.Nil(t, err, "some error occured while reading file. Error: %s", err)
	var scenarios Scenarios
	err = json.Unmarshal([]byte(file), &scenarios)
	require.Nil(t, err, "error occured during unmarshalling. Error: %s", err)

	for _, tc := range scenarios {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				trGs := &tokenregistrytypes.GenesisState{
					AdminAccount: tc.Address,
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: tc.PoolAsset, BaseDenom: tc.PoolAsset, Decimals: tc.PoolAssetDecimals, Permissions: tc.PoolAssetPermissions},
							{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: tc.NativeAssetPermissions},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(trGs)
				genesisState["tokenregistry"] = bz

				if tc.CreateBalance {
					balances := []banktypes.Balance{
						{
							Address: tc.Address,
							Coins: sdk.Coins{
								sdk.NewCoin(tc.PoolAsset, tc.ExternalBalance),
								sdk.NewCoin("rowan", tc.NativeBalance),
							},
						},
					}

					bankGs := banktypes.DefaultGenesisState()
					bankGs.Balances = append(bankGs.Balances, balances...)
					bz, _ = app.AppCodec().MarshalJSON(bankGs)
					genesisState["bank"] = bz
				}

				if tc.CreatePool {
					pools := []*types.Pool{
						{
							ExternalAsset:        &types.Asset{Symbol: tc.PoolAsset},
							NativeAssetBalance:   tc.NativeAssetAmount,
							ExternalAssetBalance: tc.ExternalAssetAmount,
							PoolUnits:            tc.PoolUnits,
						},
					}
					clpGs := types.DefaultGenesisState()
					if tc.CreateLPs {
						lps := []*types.LiquidityProvider{
							{
								Asset:                    &types.Asset{Symbol: tc.PoolAsset},
								LiquidityProviderAddress: tc.Address,
								LiquidityProviderUnits:   tc.NativeAssetAmount,
							},
						}
						clpGs.LiquidityProviders = append(clpGs.LiquidityProviders, lps...)
					}
					clpGs.Params = tc.Params
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.Address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			app.ClpKeeper.SetParams(ctx, tc.Params)
			app.ClpKeeper.SetPmtpRateParams(ctx, types.PmtpRateParams{
				PmtpPeriodBlockRate:    sdk.ZeroDec(),
				PmtpCurrentRunningRate: sdk.ZeroDec(),
				PmtpInterPolicyRate:    sdk.ZeroDec(),
			})
			app.ClpKeeper.SetPmtpEpoch(ctx, tc.Epoch)

			for i := 0; i < len(tc.ExpectedStates); i++ {
				for j := 0; j < 1000000; j++ {
					ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

					if tc.ExpectedStates[i].Height == ctx.BlockHeight() {
						clp.BeginBlocker(ctx, app.ClpKeeper)
						got, _ := app.ClpKeeper.GetPool(ctx, tc.PoolAsset)

						tc.ExpectedStates[i].Pool.SwapPriceNative = &tc.ExpectedStates[i].SwapPriceNative
						tc.ExpectedStates[i].Pool.SwapPriceExternal = &tc.ExpectedStates[i].SwapPriceExternal

						require.Equal(t, tc.ExpectedStates[i].Height, ctx.BlockHeight())
						require.Equal(t, tc.ExpectedStates[i].Pool, got)
						require.Equal(t, tc.ExpectedStates[i].PmtpRateParams, app.ClpKeeper.GetPmtpRateParams(ctx))

						break
					}
				}
			}
		})
	}
}
