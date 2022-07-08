package keeper_test

import (
	"errors"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	admintest "github.com/Sifchain/sifnode/x/admin/test"
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func TestMsgServer_DecommissionPool(t *testing.T) {
	testcases := []struct {
		name                string
		createBalance       bool
		createPool          bool
		createLPs           bool
		poolAsset           string
		address             string
		nativeBalance       sdk.Int
		externalBalance     sdk.Int
		nativeAssetAmount   sdk.Uint
		externalAssetAmount sdk.Uint
		poolUnits           sdk.Uint
		msg                 *types.MsgDecommissionPool
		expectedEvents      []sdk.Event
		err                 error
		errString           error
	}{
		{
			name:          "pool does not exist",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			msg: &types.MsgDecommissionPool{
				Signer: "xxx",
				Symbol: "xxx",
			},
			errString: errors.New("pool does not exist"),
		},
		{
			name:                "wrong address",
			createBalance:       false,
			createPool:          true,
			createLPs:           false,
			poolAsset:           "eth",
			address:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:       sdk.NewInt(10000),
			externalBalance:     sdk.NewInt(10000),
			nativeAssetAmount:   sdk.NewUint(1000),
			externalAssetAmount: sdk.NewUint(1000),
			poolUnits:           sdk.NewUint(1000),
			msg: &types.MsgDecommissionPool{
				Signer: "xxx",
				Symbol: "eth",
			},
			errString: errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name:                "invalid user",
			createBalance:       false,
			createPool:          true,
			createLPs:           false,
			poolAsset:           "eth",
			address:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:       sdk.NewInt(10000),
			externalBalance:     sdk.NewInt(10000),
			nativeAssetAmount:   sdk.NewUint(1000),
			externalAssetAmount: sdk.NewUint(1000),
			poolUnits:           sdk.NewUint(1000),
			msg: &types.MsgDecommissionPool{
				Signer: "sif1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8zzt2x5",
				Symbol: "eth",
			},
			errString: errors.New("user does not have permission to decommission pool: invalid"),
		},
		{
			name:                "balance too high",
			createBalance:       true,
			createPool:          true,
			createLPs:           true,
			poolAsset:           "eth",
			address:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:       sdk.NewInt(10000),
			externalBalance:     sdk.NewInt(10000),
			nativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
			externalAssetAmount: sdk.NewUint(1000),
			poolUnits:           sdk.NewUint(1000),
			msg: &types.MsgDecommissionPool{
				Signer: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				Symbol: "eth",
			},
			errString: errors.New("Pool Balance too high to be decommissioned"),
		},
		{
			name:                "liquidity provider does not exist",
			createBalance:       true,
			createPool:          true,
			createLPs:           false,
			poolAsset:           "eth",
			address:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:       sdk.NewInt(10000),
			externalBalance:     sdk.NewInt(10000),
			nativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
			externalAssetAmount: sdk.NewUint(1000),
			poolUnits:           sdk.NewUint(1000),
			msg: &types.MsgDecommissionPool{
				Signer: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				Symbol: "eth",
			},
			errString: errors.New("Pool Balance too high to be decommissioned"),
		},
		{
			name:                "insufficient funds",
			createBalance:       true,
			createPool:          true,
			createLPs:           true,
			poolAsset:           "eth",
			address:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:       sdk.NewInt(10000),
			externalBalance:     sdk.NewInt(10000),
			nativeAssetAmount:   sdk.NewUint(1000),
			externalAssetAmount: sdk.NewUint(1000),
			poolUnits:           sdk.NewUint(1000),
			msg: &types.MsgDecommissionPool{
				Signer: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				Symbol: "eth",
			},
			errString: errors.New("0eth is smaller than 1000eth: insufficient funds: unable to add balance: Unable to add liquidity provider"),
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
						MinCreatePoolThreshold: 100,
					}
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			msgServer := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

			// We clear the EventManager before every call as Events accumulate throughout calls
			ctx = ctx.WithEventManager(sdk.NewEventManager())
			_, err := msgServer.DecommissionPool(sdk.WrapSDKContext(ctx), tc.msg)
			if tc.expectedEvents != nil {
				checkEvents(t, tc.expectedEvents, ctx.EventManager().Events())
			}

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

func TestMsgServer_Swap(t *testing.T) {
	testcases := []struct {
		name                            string
		createBalance                   bool
		createPool                      bool
		createLPs                       bool
		poolAsset                       string
		address                         string
		nativeBalance                   sdk.Int
		externalBalance                 sdk.Int
		nativeAssetAmount               sdk.Uint
		externalAssetAmount             sdk.Uint
		poolUnits                       sdk.Uint
		poolAssetPermissions            []tokenregistrytypes.Permission
		nativeAssetPermissions          []tokenregistrytypes.Permission
		currentRowanLiquidityThreshold  sdk.Uint
		expectedRunningThresholdEnd     sdk.Uint
		maxRowanLiquidityThresholdAsset string
		maxRowanLiquidityThreshold      sdk.Uint
		expectedEvents                  []sdk.Event
		msg                             *types.MsgSwap
		err                             error
		errString                       error
	}{
		{
			name:                            "sent asset token not supported",
			createBalance:                   false,
			createPool:                      false,
			createLPs:                       false,
			currentRowanLiquidityThreshold:  sdk.NewUint(1000),
			maxRowanLiquidityThresholdAsset: "rowan",
			msg: &types.MsgSwap{
				Signer:             "xxx",
				SentAsset:          &types.Asset{Symbol: "xxx"},
				ReceivedAsset:      &types.Asset{Symbol: "xxx"},
				SentAmount:         sdk.NewUint(1),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("Token not supported by sifchain"),
		},
		{
			name:                            "received asset token not supported",
			createBalance:                   false,
			createPool:                      false,
			createLPs:                       false,
			poolAsset:                       "eth",
			currentRowanLiquidityThreshold:  sdk.NewUint(1000),
			maxRowanLiquidityThresholdAsset: "rowan",
			msg: &types.MsgSwap{
				Signer:             "xxx",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "xxx"},
				SentAmount:         sdk.NewUint(1),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("Token not supported by sifchain"),
		},
		{
			name:                            "external asset permission denied",
			createBalance:                   false,
			createPool:                      false,
			createLPs:                       false,
			poolAsset:                       "eth",
			currentRowanLiquidityThreshold:  sdk.NewUint(1000),
			maxRowanLiquidityThresholdAsset: "rowan",
			msg: &types.MsgSwap{
				Signer:             "xxx",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(1),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("permission denied for denom"),
		},
		{
			name:                            "native asset permission denied",
			createBalance:                   false,
			createPool:                      false,
			createLPs:                       false,
			poolAsset:                       "eth",
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(1000),
			maxRowanLiquidityThresholdAsset: "rowan",
			msg: &types.MsgSwap{
				Signer:             "xxx",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(1),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("permission denied for denom"),
		},
		{
			name:                            "received amount below expected",
			createBalance:                   true,
			createPool:                      true,
			createLPs:                       true,
			poolAsset:                       "eth",
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:                   sdk.NewInt(10000),
			externalBalance:                 sdk.NewInt(10000),
			nativeAssetAmount:               sdk.NewUint(1000),
			externalAssetAmount:             sdk.NewUint(1000),
			poolUnits:                       sdk.NewUint(1000),
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(1000),
			maxRowanLiquidityThresholdAsset: "rowan",
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(1),
				MinReceivingAmount: sdk.NewUint(1),
			},
			expectedEvents: []sdk.Event{sdk.NewEvent("coin_spent",
				sdk.NewAttribute("spender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				sdk.NewAttribute("amount", "1eth"),
			),
				sdk.NewEvent("coin_received",
					sdk.NewAttribute("receiver", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("amount", "1eth"),
				),
				sdk.NewEvent("transfer",
					sdk.NewAttribute("recipient", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
					sdk.NewAttribute("amount", "1eth"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
				sdk.NewEvent("swap_failed",
					sdk.NewAttribute("swap_amount", "0"),
					sdk.NewAttribute("min_threshold", "1"),
					sdk.NewAttribute("in_pool", "external_asset:<symbol:\"eth\" > native_asset_balance:\"1000\" external_asset_balance:\"1000\" pool_units:\"1000\" reward_period_native_distributed:\"0\" "),
					sdk.NewAttribute("out_pool", "external_asset:<symbol:\"eth\" > native_asset_balance:\"1000\" external_asset_balance:\"1000\" pool_units:\"1000\" reward_period_native_distributed:\"0\" "),
					sdk.NewAttribute("height", "0"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("module", "clp"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
			},
			errString: errors.New("Unable to swap, received amount is below expected"),
		},
		{
			name:                            "success",
			createBalance:                   true,
			createPool:                      true,
			createLPs:                       true,
			poolAsset:                       "eth",
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:                   sdk.NewInt(10000),
			externalBalance:                 sdk.NewInt(10000),
			nativeAssetAmount:               sdk.NewUint(1000),
			externalAssetAmount:             sdk.NewUint(1000),
			poolUnits:                       sdk.NewUint(1000),
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(1000),
			expectedRunningThresholdEnd:     sdk.NewUint(1041),
			maxRowanLiquidityThresholdAsset: "rowan",
			maxRowanLiquidityThreshold:      sdk.NewUint(2000),
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(100),
				MinReceivingAmount: sdk.NewUint(1),
			},
		},
		{
			name:                            "sell native token over threshold",
			poolAsset:                       "eth",
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(1000),
			maxRowanLiquidityThresholdAsset: "rowan",
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "rowan"},
				ReceivedAsset:      &types.Asset{Symbol: "eth"},
				SentAmount:         sdk.NewUintFromString("10000000000000000000000"),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("Unable to swap, reached maximum rowan liquidity threshold"),
		},
		{
			name:                            "sell native token over threshold - zero threshold",
			poolAsset:                       "eth",
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(0),
			maxRowanLiquidityThresholdAsset: "rowan",
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "rowan"},
				ReceivedAsset:      &types.Asset{Symbol: "eth"},
				SentAmount:         sdk.NewUintFromString("10000000000000000000000"),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("Unable to swap, reached maximum rowan liquidity threshold"),
		},
		{
			name:                            "sell native token just over threshold",
			createBalance:                   true,
			createPool:                      true,
			createLPs:                       true,
			poolAsset:                       "eth",
			nativeBalance:                   sdk.NewInt(10000),
			externalBalance:                 sdk.NewInt(10000),
			nativeAssetAmount:               sdk.NewUint(100000),
			externalAssetAmount:             sdk.NewUint(200000),
			poolUnits:                       sdk.NewUint(100000),
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(1000),
			maxRowanLiquidityThresholdAsset: "eth",
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "rowan"},
				ReceivedAsset:      &types.Asset{Symbol: "eth"},
				SentAmount:         sdk.NewUint(251),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("Unable to swap, reached maximum rowan liquidity threshold"),
		},
		{
			name:                            "sell native token just below threshold",
			createBalance:                   true,
			createPool:                      true,
			createLPs:                       true,
			poolAsset:                       "eth",
			nativeBalance:                   sdk.NewInt(10000),
			externalBalance:                 sdk.NewInt(10000),
			nativeAssetAmount:               sdk.NewUint(100000),
			externalAssetAmount:             sdk.NewUint(200000),
			poolUnits:                       sdk.NewUint(100000),
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(4000),
			expectedRunningThresholdEnd:     sdk.NewUint(3000),
			maxRowanLiquidityThresholdAsset: "eth",
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "rowan"},
				ReceivedAsset:      &types.Asset{Symbol: "eth"},
				SentAmount:         sdk.NewUint(250),
				MinReceivingAmount: sdk.NewUint(1),
			},
		},
		{
			name:                            "buy rowan, threshold in eth",
			createBalance:                   true,
			createPool:                      true,
			createLPs:                       true,
			poolAsset:                       "eth",
			nativeBalance:                   sdk.NewInt(10000),
			externalBalance:                 sdk.NewInt(10000),
			nativeAssetAmount:               sdk.NewUint(100000),
			externalAssetAmount:             sdk.NewUint(200000),
			poolUnits:                       sdk.NewUint(100000),
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(200),
			expectedRunningThresholdEnd:     sdk.NewUint(300),
			maxRowanLiquidityThresholdAsset: "eth",
			maxRowanLiquidityThreshold:      sdk.NewUint(1000),
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(100),
				MinReceivingAmount: sdk.NewUint(1),
			},
		},
		{
			name:                            "buy rowan, threshold in eth, low max liquidity threshold",
			createBalance:                   true,
			createPool:                      true,
			createLPs:                       true,
			poolAsset:                       "eth",
			nativeBalance:                   sdk.NewInt(10000),
			externalBalance:                 sdk.NewInt(10000),
			nativeAssetAmount:               sdk.NewUint(100000),
			externalAssetAmount:             sdk.NewUint(200000),
			poolUnits:                       sdk.NewUint(100000),
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUint(200),
			expectedRunningThresholdEnd:     sdk.NewUint(250),
			maxRowanLiquidityThresholdAsset: "eth",
			maxRowanLiquidityThreshold:      sdk.NewUint(250),
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(100),
				MinReceivingAmount: sdk.NewUint(1),
			},
		},
		{
			name:                            "buy rowan, threshold in eth, maximum max liquidity threshold",
			createBalance:                   true,
			createPool:                      true,
			createLPs:                       true,
			poolAsset:                       "eth",
			nativeBalance:                   sdk.NewInt(10000),
			externalBalance:                 sdk.NewInt(10000),
			nativeAssetAmount:               sdk.NewUint(100000),
			externalAssetAmount:             sdk.NewUint(200000),
			poolUnits:                       sdk.NewUint(100000),
			address:                         "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			poolAssetPermissions:            []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions:          []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			currentRowanLiquidityThreshold:  sdk.NewUintFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
			expectedRunningThresholdEnd:     sdk.NewUintFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
			maxRowanLiquidityThresholdAsset: "eth",
			maxRowanLiquidityThreshold:      sdk.NewUintFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(100),
				MinReceivingAmount: sdk.NewUint(1),
			},

			expectedEvents: []sdk.Event{sdk.NewEvent("coin_spent",
				sdk.NewAttribute("spender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				sdk.NewAttribute("amount", "100eth"),
			),
				sdk.NewEvent("coin_received",
					sdk.NewAttribute("receiver", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("amount", "100eth"),
				),
				sdk.NewEvent("transfer",
					sdk.NewAttribute("recipient", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
					sdk.NewAttribute("amount", "100eth"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
				sdk.NewEvent("coin_spent",
					sdk.NewAttribute("spender", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("amount", "25rowan"),
				),
				sdk.NewEvent("coin_received",
					sdk.NewAttribute("receiver", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
					sdk.NewAttribute("amount", "25rowan"),
				),
				sdk.NewEvent("transfer",
					sdk.NewAttribute("recipient", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
					sdk.NewAttribute("sender", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("amount", "25rowan"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("sender", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
				),
				sdk.NewEvent("swap_successful",
					sdk.NewAttribute("swap_amount", "25"),
					sdk.NewAttribute("liquidity_fee", "0"),
					sdk.NewAttribute("price_impact", "0"),
					sdk.NewAttribute("in_pool", "external_asset:<symbol:\"eth\" > native_asset_balance:\"100000\" external_asset_balance:\"200000\" pool_units:\"100000\" reward_period_native_distributed:\"0\" "),
					sdk.NewAttribute("out_pool", "external_asset:<symbol:\"eth\" > native_asset_balance:\"100000\" external_asset_balance:\"200000\" pool_units:\"100000\" reward_period_native_distributed:\"0\" "),
					sdk.NewAttribute("pmtp_block_rate", "1.000000000000000000"),
					sdk.NewAttribute("pmtp_current_running_rate", "1.000000000000000000"),
					sdk.NewAttribute("height", "0"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("module", "clp"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
			},
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
					externalCLP, _ := sdk.NewIntFromString(tc.externalAssetAmount.String())
					nativeCLP, _ := sdk.NewIntFromString(tc.nativeAssetAmount.String())
					clpAddrs := app.AccountKeeper.GetModuleAddress("clp").String()

					balances := []banktypes.Balance{
						{
							Address: tc.address,
							Coins: sdk.Coins{
								sdk.NewCoin(tc.poolAsset, tc.externalBalance),
								sdk.NewCoin("rowan", tc.nativeBalance),
							},
						},
						{
							Address: clpAddrs,
							Coins: sdk.Coins{
								sdk.NewCoin(tc.poolAsset, externalCLP),
								sdk.NewCoin("rowan", nativeCLP),
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
						MinCreatePoolThreshold: 100,
					}
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			app.ClpKeeper.SetPmtpCurrentRunningRate(ctx, sdk.NewDec(1))

			liquidityProtectionParam := app.ClpKeeper.GetLiquidityProtectionParams(ctx)
			liquidityProtectionParam.MaxRowanLiquidityThresholdAsset = tc.maxRowanLiquidityThresholdAsset
			liquidityProtectionParam.MaxRowanLiquidityThreshold = tc.maxRowanLiquidityThreshold
			liquidityProtectionParam.IsActive = true
			app.ClpKeeper.SetLiquidityProtectionParams(ctx, liquidityProtectionParam)
			app.ClpKeeper.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, tc.currentRowanLiquidityThreshold)

			msgServer := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

			// We clear the EventManager before every call as Events accumulate throughout calls
			ctx = ctx.WithEventManager(sdk.NewEventManager())
			_, err := msgServer.Swap(sdk.WrapSDKContext(ctx), tc.msg)
			if tc.expectedEvents != nil {
				checkEvents(t, tc.expectedEvents, ctx.EventManager().Events())
			}

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			runningThreshold := app.ClpKeeper.GetLiquidityProtectionRateParams(ctx).CurrentRowanLiquidityThreshold
			require.Equal(t, tc.expectedRunningThresholdEnd.String(), runningThreshold.String())

		})
	}
}

func TestMsgServer_RemoveLiquidity(t *testing.T) {
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
		msg                    *types.MsgRemoveLiquidity
		expectedEvents         []sdk.Event
		err                    error
		errString              error
	}{
		{
			name:          "sent asset token not supported",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			msg: &types.MsgRemoveLiquidity{
				Signer:        "xxx",
				ExternalAsset: &types.Asset{Symbol: "xxx"},
				WBasisPoints:  sdk.NewInt(1),
				Asymmetry:     sdk.NewInt(1),
			},
			errString: errors.New("Token not supported by sifchain"),
		},
		{
			name:          "received asset token not supported",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			poolAsset:     "eth",
			msg: &types.MsgRemoveLiquidity{
				Signer:        "xxx",
				ExternalAsset: &types.Asset{Symbol: "xxx"},
				WBasisPoints:  sdk.NewInt(1),
				Asymmetry:     sdk.NewInt(1),
			},
			errString: errors.New("Token not supported by sifchain"),
		},
		{
			name:          "external asset permission denied",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			poolAsset:     "eth",
			msg: &types.MsgRemoveLiquidity{
				Signer:        "xxx",
				ExternalAsset: &types.Asset{Symbol: "eth"},
				WBasisPoints:  sdk.NewInt(1),
				Asymmetry:     sdk.NewInt(1),
			},
			errString: errors.New("permission denied for denom"),
		},
		{
			name:                 "pool does not exist",
			createBalance:        false,
			createPool:           false,
			createLPs:            false,
			poolAsset:            "xxx",
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgRemoveLiquidity{
				Signer:        "xxx",
				ExternalAsset: &types.Asset{Symbol: "xxx"},
				WBasisPoints:  sdk.NewInt(1),
				Asymmetry:     sdk.NewInt(1),
			},
			errString: errors.New("pool does not exist"),
		},
		{
			name:                 "no lp",
			createBalance:        true,
			createPool:           true,
			createLPs:            false,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.NewInt(10000),
			externalBalance:      sdk.NewInt(10000),
			nativeAssetAmount:    sdk.NewUint(1000),
			externalAssetAmount:  sdk.NewUint(1000),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgRemoveLiquidity{
				Signer:        "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset: &types.Asset{Symbol: "eth"},
				WBasisPoints:  sdk.NewInt(1),
				Asymmetry:     sdk.NewInt(1),
			},
			errString: errors.New("liquidity Provider does not exist"),
		},
		{
			name:                   "received amount below expected",
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
			msg: &types.MsgRemoveLiquidity{
				Signer:        "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset: &types.Asset{Symbol: "eth"},
				WBasisPoints:  sdk.NewInt(1),
				Asymmetry:     sdk.NewInt(0),
			},
			expectedEvents: []sdk.Event{sdk.NewEvent("removed_liquidity",
				sdk.NewAttribute("liquidity_provider", "asset:<symbol:\"eth\" > liquidity_provider_units:\"1000\" liquidity_provider_address:\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\" "),
				sdk.NewAttribute("liquidity_units", "0"),
				sdk.NewAttribute("pmtp_block_rate", "1.000000000000000000"),
				sdk.NewAttribute("pmtp_current_running_rate", "1.000000000000000000"),
				sdk.NewAttribute("height", "0"),
			),
				sdk.NewEvent("message",
					sdk.NewAttribute("module", "clp"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
			},
		},
		{
			name:                   "received amount below expected",
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
			msg: &types.MsgRemoveLiquidity{
				Signer:        "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset: &types.Asset{Symbol: "eth"},
				WBasisPoints:  sdk.NewInt(1),
				Asymmetry:     sdk.NewInt(1),
			},
			err: types.ErrAsymmetricRemove,
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
					clpGs.Params = types.Params{
						MinCreatePoolThreshold: 100,
					}
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			app.ClpKeeper.SetPmtpCurrentRunningRate(ctx, sdk.NewDec(1))

			msgServer := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

			// We clear the EventManager before every call as Events accumulate throughout calls
			ctx = ctx.WithEventManager(sdk.NewEventManager())
			_, err := msgServer.RemoveLiquidity(sdk.WrapSDKContext(ctx), tc.msg)
			if tc.expectedEvents != nil {
				checkEvents(t, tc.expectedEvents, ctx.EventManager().Events())
			}

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

func TestMsgServer_CreatePool(t *testing.T) {
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
		msg                    *types.MsgCreatePool
		expectedEvents         []sdk.Event
		err                    error
		errString              error
	}{
		{
			name:          "total amount too low",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			msg: &types.MsgCreatePool{
				Signer:              "xxx",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUint(1000),
				ExternalAssetAmount: sdk.NewUint(1000),
			},
			errString: errors.New("total amount is less than minimum threshold"),
		},
		{
			name:          "external asset token not supported",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			msg: &types.MsgCreatePool{
				Signer:              "xxx",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("Token not supported by sifchain"),
		},
		{
			name:          "external asset permission denied",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			poolAsset:     "eth",
			msg: &types.MsgCreatePool{
				Signer:              "xxx",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("permission denied for denom"),
		},
		{
			name:                 "pool already exists",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.NewInt(10000),
			externalBalance:      sdk.NewInt(10000),
			nativeAssetAmount:    sdk.NewUint(1000),
			externalAssetAmount:  sdk.NewUint(1000),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgCreatePool{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("Unable to create pool"),
		},
		{
			name:                 "user does have enough balance of required coin",
			createBalance:        true,
			createPool:           false,
			createLPs:            false,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.NewInt(10000),
			externalBalance:      sdk.NewInt(10000),
			nativeAssetAmount:    sdk.NewUint(1000),
			externalAssetAmount:  sdk.NewUint(1000),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgCreatePool{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("user does not have enough balance of the required coin: Unable to set pool"),
		},
		{
			name:                 "successful",
			createBalance:        true,
			createPool:           false,
			createLPs:            false,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:      sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:    sdk.NewUint(1000),
			externalAssetAmount:  sdk.NewUint(1000),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgCreatePool{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			expectedEvents: []sdk.Event{sdk.NewEvent("coin_spent",
				sdk.NewAttribute("spender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				sdk.NewAttribute("amount", "1000000000000000000eth,1000000000000000000rowan"),
			),
				sdk.NewEvent("coin_received",
					sdk.NewAttribute("receiver", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("amount", "1000000000000000000eth,1000000000000000000rowan"),
				),
				sdk.NewEvent("transfer",
					sdk.NewAttribute("recipient", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
					sdk.NewAttribute("amount", "1000000000000000000eth,1000000000000000000rowan"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
				sdk.NewEvent("created_new_pool",
					sdk.NewAttribute("pool", "external_asset:<symbol:\"eth\" > native_asset_balance:\"1000000000000000000\" external_asset_balance:\"1000000000000000000\" pool_units:\"1000000000000000000\" reward_period_native_distributed:\"0\" "),
					sdk.NewAttribute("height", "0"),
				),
				sdk.NewEvent("created_new_liquidity_provider",
					sdk.NewAttribute("liquidity_provider", "asset:<symbol:\"eth\" > liquidity_provider_units:\"1000000000000000000\" liquidity_provider_address:\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\" "),
					sdk.NewAttribute("height", "0"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("module", "clp"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
			},
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
					clpGs.Params = types.Params{
						MinCreatePoolThreshold: 100,
					}
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			app.ClpKeeper.SetPmtpCurrentRunningRate(ctx, sdk.NewDec(1))

			msgServer := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

			// We clear the EventManager before every call as Events accumulate throughout calls
			ctx = ctx.WithEventManager(sdk.NewEventManager())
			_, err := msgServer.CreatePool(sdk.WrapSDKContext(ctx), tc.msg)
			if tc.expectedEvents != nil {
				checkEvents(t, tc.expectedEvents, ctx.EventManager().Events())
			}

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

func TestMsgServer_AddLiquidity(t *testing.T) {
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
		msg                    *types.MsgAddLiquidity
		expectedEvents         []sdk.Event
		err                    error
		errString              error
	}{
		{
			name:          "external asset token not supported",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			msg: &types.MsgAddLiquidity{
				Signer:              "xxx",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("Token not supported by sifchain"),
		},
		{
			name:          "external asset permission denied",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			poolAsset:     "eth",
			msg: &types.MsgAddLiquidity{
				Signer:              "xxx",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("permission denied for denom"),
		},
		{
			name:                 "pool does not exist",
			createBalance:        true,
			createPool:           false,
			createLPs:            false,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.NewInt(10000),
			externalBalance:      sdk.NewInt(10000),
			nativeAssetAmount:    sdk.NewUint(1000),
			externalAssetAmount:  sdk.NewUint(1000),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("pool does not exist"),
		},
		{
			name:                 "user does have enough balance of required coin",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.NewInt(10000),
			externalBalance:      sdk.NewInt(10000),
			nativeAssetAmount:    sdk.NewUint(1000),
			externalAssetAmount:  sdk.NewUint(1000),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("user does not have enough balance of the required coin: Unable to add liquidity"),
		},
		{
			name:                 "successful",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:      sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:    sdk.NewUint(1000),
			externalAssetAmount:  sdk.NewUint(1000),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			expectedEvents: []sdk.Event{sdk.NewEvent("coin_spent",
				sdk.NewAttribute("spender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				sdk.NewAttribute("amount", "1000000000000000000eth,1000000000000000000rowan"),
			),
				sdk.NewEvent("coin_received",
					sdk.NewAttribute("receiver", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("amount", "1000000000000000000eth,1000000000000000000rowan"),
				),
				sdk.NewEvent("transfer",
					sdk.NewAttribute("recipient", "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
					sdk.NewAttribute("amount", "1000000000000000000eth,1000000000000000000rowan"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
				sdk.NewEvent("added_liquidity",
					sdk.NewAttribute("liquidity_provider", "asset:<symbol:\"eth\" > liquidity_provider_units:\"1000000000000001000\" liquidity_provider_address:\"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd\" "),
					sdk.NewAttribute("liquidity_units", "1000000000000000000"),
					sdk.NewAttribute("height", "0"),
				),
				sdk.NewEvent("message",
					sdk.NewAttribute("module", "clp"),
					sdk.NewAttribute("sender", "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"),
				),
			},
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
					clpGs.Params = types.Params{
						MinCreatePoolThreshold: 100,
					}
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			app.ClpKeeper.SetPmtpCurrentRunningRate(ctx, sdk.NewDec(1))

			msgServer := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

			// We clear the EventManager before every call as Events accumulate throughout calls
			ctx = ctx.WithEventManager(sdk.NewEventManager())
			_, err := msgServer.AddLiquidity(sdk.WrapSDKContext(ctx), tc.msg)
			if tc.expectedEvents != nil {
				checkEvents(t, tc.expectedEvents, ctx.EventManager().Events())
			}

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

func TestMsgServer_AddProviderDistribution(t *testing.T) {
	admin := "sif1gy2ne7m62uer4h5s4e7xlfq7aeem5zpwx6nu9q"
	nonAdmin := "sif1gy2ne7m62uer4h5s4e7xlfq7aeem5zpwx6nu9r"
	ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
		adminGs := &admintypes.GenesisState{
			AdminAccounts: admintest.GetAdmins(admin),
		}
		bz, _ := app.AppCodec().MarshalJSON(adminGs)
		genesisState["admin"] = bz
		trGs := &tokenregistrytypes.GenesisState{
			Registry: nil,
		}
		bz, _ = app.AppCodec().MarshalJSON(trGs)
		genesisState["tokenregistry"] = bz

		return genesisState
	})
	msgServer := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

	_, err := msgServer.AddProviderDistributionPeriod(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)

	var periods []*types.ProviderDistributionPeriod
	validPeriod := types.ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 10, DistributionPeriodBlockRate: sdk.NewDecWithPrec(1, 2), DistributionPeriodMod: 1}
	wrongPeriod := types.ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 9, DistributionPeriodBlockRate: sdk.NewDecWithPrec(1, 2), DistributionPeriodMod: 1}

	periods = append(periods, &wrongPeriod)
	msg := types.MsgAddProviderDistributionPeriodRequest{Signer: admin, DistributionPeriods: periods}
	_, err = msgServer.AddProviderDistributionPeriod(sdk.WrapSDKContext(ctx), &msg)
	require.Error(t, err)
	// check events didn't fire
	require.Equal(t, len(ctx.EventManager().Events()), 0)

	periods[0] = &validPeriod
	msg = types.MsgAddProviderDistributionPeriodRequest{Signer: admin, DistributionPeriods: periods}
	_, err = msgServer.AddProviderDistributionPeriod(sdk.WrapSDKContext(ctx), &msg)
	require.NoError(t, err)
	// check correct events fired
	require.Equal(t, len(ctx.EventManager().Events()), 2)
	expectedEvents := []sdk.Event{sdk.NewEvent("lppd_new_policy",
		sdk.NewAttribute("lppd_params", "distribution_periods:<distribution_period_block_rate:\"10000000000000000\" distribution_period_start_block:10 distribution_period_end_block:10 distribution_period_mod:1 > "),
		sdk.NewAttribute("height", "0"),
	),
		sdk.NewEvent("message",
			sdk.NewAttribute("module", "clp"),
			sdk.NewAttribute("sender", "sif1gy2ne7m62uer4h5s4e7xlfq7aeem5zpwx6nu9q"),
		),
	}
	checkEvents(t, expectedEvents, ctx.EventManager().Events())

	// non admin acc
	msg = types.MsgAddProviderDistributionPeriodRequest{Signer: nonAdmin, DistributionPeriods: periods}
	_, err = msgServer.AddProviderDistributionPeriod(sdk.WrapSDKContext(ctx), &msg)
	require.Error(t, err)
	// check no additional events fired
	require.Equal(t, len(ctx.EventManager().Events()), 2)

	cbp := app.ClpKeeper.GetProviderDistributionParams(ctx)
	require.NotNil(t, cbp)
	require.Equal(t, 1, len(cbp.DistributionPeriods))
	require.Equal(t, *cbp.DistributionPeriods[0], validPeriod)
}

func checkEvents(t *testing.T, expectedEvents, events []sdk.Event) {
	require.ElementsMatch(t, expectedEvents, events)
	// For more readable debugging on failure, replace above with the below
	// and add prints were needed
	//
	//for i, event := range ctx.EventManager().Events() {
	//	require.Equal(t, tc.expectedEvents[i].Type, event.Type)
	//	expectedAttributes := tc.expectedEvents[i].Attributes
	//	require.Equal(t, len(expectedAttributes), len(event.Attributes))
	//	for j, attr := range event.Attributes {
	//		require.Equal(t, expectedAttributes[j].String(), attr.String())
	//	}
	//}
}
