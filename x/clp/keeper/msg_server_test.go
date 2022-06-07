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
		decimals               int64
		address                string
		nativeBalance          sdk.Int
		externalBalance        sdk.Int
		nativeAssetAmount      sdk.Uint
		externalAssetAmount    sdk.Uint
		poolUnits              sdk.Uint
		nativeBalanceEnd       sdk.Int
		externalBalanceEnd     sdk.Int
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
			decimals:      18,
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
			decimals:      18,
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
			decimals:             18,
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
			decimals:               18,
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			nativeBalanceEnd:       sdk.NewInt(10000),
			externalBalanceEnd:     sdk.NewInt(10000),
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
			name:                   "success",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			decimals:               18,
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			nativeBalanceEnd:       sdk.NewInt(10082),
			externalBalanceEnd:     sdk.NewInt(9900),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "eth"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(100),
				MinReceivingAmount: sdk.NewUint(1),
			},
		},
		{
			name:                   "correct amount when external asset has eighteen decimals",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "decvar",
			decimals:               18,
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.NewInt(10),
			externalBalance:        sdk.NewInt(10),
			nativeAssetAmount:      sdk.NewUint(400),
			externalAssetAmount:    sdk.NewUint(100),
			poolUnits:              sdk.NewUint(100),
			nativeBalanceEnd:       sdk.NewInt(13),
			externalBalanceEnd:     sdk.NewInt(9),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "decvar"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(1),
				MinReceivingAmount: sdk.NewUint(0),
			},
		},
		{
			name:                   "correct amount when external asset has nineteen decimals",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "decvar",
			decimals:               19,
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.NewInt(10),
			externalBalance:        sdk.NewInt(10),
			nativeAssetAmount:      sdk.NewUint(400),
			externalAssetAmount:    sdk.NewUint(100),
			poolUnits:              sdk.NewUint(100),
			nativeBalanceEnd:       sdk.NewInt(13),
			externalBalanceEnd:     sdk.NewInt(9),
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgSwap{
				Signer:             "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				SentAsset:          &types.Asset{Symbol: "decvar"},
				ReceivedAsset:      &types.Asset{Symbol: "rowan"},
				SentAmount:         sdk.NewUint(1),
				MinReceivingAmount: sdk.NewUint(0),
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
							{Denom: tc.poolAsset, BaseDenom: tc.poolAsset, Decimals: tc.decimals, Permissions: tc.poolAssetPermissions},
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

			app.ClpKeeper.SetPmtpCurrentRunningRate(ctx, sdk.NewDec(0))

			msgServer := clpkeeper.NewMsgServerImpl(app.ClpKeeper)

			_, err := msgServer.Swap(sdk.WrapSDKContext(ctx), tc.msg)

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			sentAssetBalanceRequest := banktypes.QueryBalanceRequest{
				Address: tc.address,
				Denom:   tc.msg.SentAsset.Symbol,
			}

			sentAssetBalanceResponse, err := app.BankKeeper.Balance(sdk.WrapSDKContext(ctx), &sentAssetBalanceRequest)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			sentAssetBalance := sentAssetBalanceResponse.Balance.Amount

			require.Equal(t, tc.externalBalanceEnd.String(), sentAssetBalance.String())

			nativeAssetBalanceRequest := banktypes.QueryBalanceRequest{
				Address: tc.address,
				Denom:   "rowan",
			}

			nativeAssetBalanceResponse, err := app.BankKeeper.Balance(sdk.WrapSDKContext(ctx), &nativeAssetBalanceRequest)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			nativeAssetBalance := nativeAssetBalanceResponse.Balance.Amount

			require.Equal(t, tc.nativeBalanceEnd.String(), nativeAssetBalance.String())

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

			_, err := msgServer.CreatePool(sdk.WrapSDKContext(ctx), tc.msg)

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
		externalDecimals       int64
		address                string
		nativeBalance          sdk.Int
		externalBalance        sdk.Int
		nativeAssetAmount      sdk.Uint
		externalAssetAmount    sdk.Uint
		poolUnits              sdk.Uint
		poolAssetPermissions   []tokenregistrytypes.Permission
		nativeAssetPermissions []tokenregistrytypes.Permission
		msg                    *types.MsgAddLiquidity
		expectedPoolUnits      sdk.Uint
		expectedLPUnits        sdk.Uint
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
			expectedPoolUnits: sdk.NewUintFromString("1000000000000001000"),
			expectedLPUnits:   sdk.NewUintFromString("1000000000000000000"),
		},
		{
			name:                 "native asset is 0",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:      sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:    sdk.ZeroUint(),
			externalAssetAmount:  sdk.NewUint(1000),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.ZeroUint(),
				ExternalAssetAmount: sdk.NewUintFromString(types.PoolThrehold),
			},
			errString: errors.New("0: insufficient funds"),
		},
		{
			name:                 "external asset is 0",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:      sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:    sdk.NewUint(1000),
			externalAssetAmount:  sdk.ZeroUint(),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.NewUintFromString(types.PoolThrehold),
				ExternalAssetAmount: sdk.ZeroUint(),
			},
			errString: errors.New("0: insufficient funds"),
		},
		{
			name:                 "external and native asset is 0",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "eth",
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			externalBalance:      sdk.Int(sdk.NewUintFromString(types.PoolThrehold)),
			nativeAssetAmount:    sdk.ZeroUint(),
			externalAssetAmount:  sdk.ZeroUint(),
			poolUnits:            sdk.NewUint(1000),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "eth"},
				NativeAssetAmount:   sdk.ZeroUint(),
				ExternalAssetAmount: sdk.ZeroUint(),
			},
			errString: errors.New("Tx amount is too low"),
		},
		{
			name:                 "success",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "cusdc",
			externalDecimals:     6,
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString("4000000000000000000")),
			externalBalance:      sdk.Int(sdk.NewUint(68140)),
			nativeAssetAmount:    sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetAmount:  sdk.NewUint(2674623482959),
			poolUnits:            sdk.NewUintFromString("23662660550457383692937954"),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "cusdc"},
				NativeAssetAmount:   sdk.NewUintFromString("4000000000000000000"),
				ExternalAssetAmount: sdk.NewUint(68140),
			},
			expectedPoolUnits: sdk.NewUintFromString("23662661153298835875523384"),
			expectedLPUnits:   sdk.NewUintFromString("602841452182585430"),
		},
		{
			// Same test as above but with external asset amount just below top limit
			name:                 "success (normalized) ratios diff = 0.000000000000000499",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "cusdc",
			externalDecimals:     6,
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString("4000000000000000000")),
			externalBalance:      sdk.Int(sdk.NewUint(70140)),
			nativeAssetAmount:    sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetAmount:  sdk.NewUint(2674623482959),
			poolUnits:            sdk.NewUintFromString("23662660550457383692937954"),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "cusdc"},
				NativeAssetAmount:   sdk.NewUintFromString("4000000000000000000"),
				ExternalAssetAmount: sdk.NewUint(70140),
			},
			expectedPoolUnits: sdk.NewUintFromString("23662661162145935094484778"),
			expectedLPUnits:   sdk.NewUintFromString("611688551401546824"),
		},
		{
			// Same test as above but with external asset amount just above top limit
			name:                 "failure (normalized) ratios diff = 0.000000000000000500",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "cusdc",
			externalDecimals:     6,
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString("4000000000000000000")),
			externalBalance:      sdk.Int(sdk.NewUint(70141)),
			nativeAssetAmount:    sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetAmount:  sdk.NewUint(2674623482959),
			poolUnits:            sdk.NewUintFromString("23662660550457383692937954"),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "cusdc"},
				NativeAssetAmount:   sdk.NewUintFromString("4000000000000000000"),
				ExternalAssetAmount: sdk.NewUint(70141),
			},
			errString: errors.New("Cannot add liquidity with asymmetric ratio"),
		},
		{
			// Same test as above but with external asset amount just above bottom limit
			name:                 "success (normalized) ratios diff = 0.000000000000000499",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "cusdc",
			externalDecimals:     6,
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString("4000000000000000000")),
			externalBalance:      sdk.Int(sdk.NewUint(66141)),
			nativeAssetAmount:    sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetAmount:  sdk.NewUint(2674623482959),
			poolUnits:            sdk.NewUintFromString("23662660550457383692937954"),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "cusdc"},
				NativeAssetAmount:   sdk.NewUintFromString("4000000000000000000"),
				ExternalAssetAmount: sdk.NewUint(66141),
			},
			expectedPoolUnits: sdk.NewUintFromString("23662661144456159305055227"),
			expectedLPUnits:   sdk.NewUintFromString("593998775612117273"),
		},
		{
			// Same test as above but with external asset amount just below bottom limit
			name:                 "failure (normalized) ratios diff = 0.000000000000000500",
			createBalance:        true,
			createPool:           true,
			createLPs:            true,
			poolAsset:            "cusdc",
			externalDecimals:     6,
			address:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:        sdk.Int(sdk.NewUintFromString("4000000000000000000")),
			externalBalance:      sdk.Int(sdk.NewUint(66140)),
			nativeAssetAmount:    sdk.NewUintFromString("157007500498726220240179086"),
			externalAssetAmount:  sdk.NewUint(2674623482959),
			poolUnits:            sdk.NewUintFromString("23662660550457383692937954"),
			poolAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			msg: &types.MsgAddLiquidity{
				Signer:              "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				ExternalAsset:       &types.Asset{Symbol: "cusdc"},
				NativeAssetAmount:   sdk.NewUintFromString("4000000000000000000"),
				ExternalAssetAmount: sdk.NewUint(66140),
			},
			errString: errors.New("Cannot add liquidity with asymmetric ratio"),
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
							{Denom: tc.poolAsset, BaseDenom: tc.poolAsset, Decimals: tc.externalDecimals, Permissions: tc.poolAssetPermissions},
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
								LiquidityProviderUnits:   sdk.NewUint(0),
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

			lp, _ := app.ClpKeeper.GetLiquidityProvider(ctx, tc.poolAsset, tc.address)
			pool, _ := app.ClpKeeper.GetPool(ctx, tc.poolAsset)

			require.Equal(t, tc.expectedPoolUnits.String(), pool.PoolUnits.String()) // compare strings so that the expected amounts can be read from the failure message
			require.Equal(t, tc.expectedLPUnits.String(), lp.LiquidityProviderUnits.String())
		})
	}
}
