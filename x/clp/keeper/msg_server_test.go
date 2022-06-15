package keeper_test

import (
	"errors"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
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

			_, err := msgServer.DecommissionPool(sdk.WrapSDKContext(ctx), tc.msg)

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
		msg                    *types.MsgSwap
		err                    error
		errString              error
	}{
		{
			name:          "sent asset token not supported",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
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
			name:          "received asset token not supported",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			poolAsset:     "eth",
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
			name:          "external asset permission denied",
			createBalance: false,
			createPool:    false,
			createLPs:     false,
			poolAsset:     "eth",
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
			name:                 "native asset permission denied",
			createBalance:        false,
			createPool:           false,
			createLPs:            false,
			poolAsset:            "eth",
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
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
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(1),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("Unable to swap, received amount is below expected"),
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
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(100),
				MinReceivingAmount: sdk.NewUint(1),
			},
			errString: errors.New("0rowan is smaller than 41rowan: insufficient funds: Unable to swap"),
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

			_, err := msgServer.Swap(sdk.WrapSDKContext(ctx), tc.msg)

			//if tc.errString != nil {
			//	require.EqualError(t, err, tc.errString.Error())
			//	return
			//}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			//require.NoError(t, err)
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
					AdminAccounts: test.GetAdmins(tc.address),
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

			_, err := msgServer.RemoveLiquidity(sdk.WrapSDKContext(ctx), tc.msg)

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
					AdminAccounts: test.GetAdmins(tc.address),
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
			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			if tc.expectedEvents != nil {
				require.ElementsMatch(t, tc.expectedEvents, ctx.EventManager().Events())
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

			_, err := msgServer.AddLiquidity(sdk.WrapSDKContext(ctx), tc.msg)

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
		trGs := &tokenregistrytypes.GenesisState{
			AdminAccounts: test.GetAdmins(admin),
			Registry:      nil,
		}
		bz, _ := app.AppCodec().MarshalJSON(trGs)
		genesisState["tokenregistry"] = bz

		return genesisState
	})
	msgServer := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

	_, err := msgServer.AddProviderDistributionPeriod(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err)

	var periods []*types.ProviderDistributionPeriod
	validPeriod := types.ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 10, DistributionPeriodBlockRate: sdk.NewDecWithPrec(1, 2)}
	wrongPeriod := types.ProviderDistributionPeriod{DistributionPeriodStartBlock: 10, DistributionPeriodEndBlock: 9, DistributionPeriodBlockRate: sdk.NewDecWithPrec(1, 2)}

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
	// check events fired
	require.Equal(t, len(ctx.EventManager().Events()), 2)

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
