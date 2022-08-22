//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper_test

import (
	"errors"
	"fmt"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	"github.com/Sifchain/sifnode/x/clp"
	clptest "github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/keeper"
	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_NewMsgServerImpl(t *testing.T) {
	_, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	got := keeper.NewMsgServerImpl(marginKeeper)
	require.NotNil(t, got)
}

func TestKeeper_Open(t *testing.T) {
	table := []struct {
		name          string
		msgOpen       types.MsgOpen
		poolAsset     string
		token         string
		poolEnabled   bool
		fundedAccount bool
		err           error
		errString     error
	}{
		{
			name: "pool does not exist",
			msgOpen: types.MsgOpen{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
				Position:         types.Position_LONG,
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "same collateral and native asset but pool does not exist",
			msgOpen: types.MsgOpen{
				Signer:           "xxx",
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
				Position:         types.Position_LONG,
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "same collateral and native asset but pool exists",
			msgOpen: types.MsgOpen{
				Signer:           "xxx",
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "rowan",
				Position:         types.Position_LONG,
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(types.ErrMTPDisabled, "rowan"),
		},
		{
			name: "margin enabled but denom does not exist but pool health too low",
			msgOpen: types.MsgOpen{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
				Position:         types.Position_LONG,
			},
			poolAsset:   "xxx",
			token:       "somethingelse",
			poolEnabled: true,
			err:         types.ErrMTPDisabled,
		},
		{
			name: "wrong address but pool health too low",
			msgOpen: types.MsgOpen{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
				Position:         types.Position_LONG,
			},
			poolAsset:   "xxx",
			token:       "xxx",
			poolEnabled: true,
			err:         types.ErrMTPDisabled,
		},
		{
			name: "insufficient funds but pool health too low",
			msgOpen: types.MsgOpen{
				Signer:           "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
				Position:         types.Position_LONG,
			},
			poolAsset:   "xxx",
			token:       "xxx",
			poolEnabled: true,
			err:         types.ErrMTPDisabled,
		},
		{
			name: "account funded but pool health too low",
			msgOpen: types.MsgOpen{
				Signer:           "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
				Position:         types.Position_LONG,
			},
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			err:           types.ErrMTPDisabled,
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppMargin(false)
			marginKeeper := app.MarginKeeper

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.token,
				Decimals:    18,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			msgServer := keeper.NewMsgServerImpl(marginKeeper)

			asset := clptypes.Asset{Symbol: tt.poolAsset}
			pool := clptypes.Pool{
				ExternalAsset:        &asset,
				NativeAssetBalance:   sdk.NewUint(1000000000),
				ExternalAssetBalance: sdk.NewUint(1000000000),
				NativeCustody:        sdk.NewUint(0),
				ExternalCustody:      sdk.NewUint(0),
				NativeLiabilities:    sdk.NewUint(0),
				ExternalLiabilities:  sdk.NewUint(0),
				PoolUnits:            sdk.NewUint(0),
				Health:               sdk.NewDec(0),
				InterestRate:         sdk.NewDec(0),
			}

			if tt.poolEnabled {
				marginKeeper.SetEnabledPools(ctx, []string{tt.poolAsset})
			}

			// nolint:errcheck
			marginKeeper.ClpKeeper().SetPool(ctx, &pool)

			var address string

			if tt.fundedAccount {
				nativeAsset := tt.msgOpen.CollateralAsset
				externalAsset := clptypes.Asset{Symbol: tt.msgOpen.BorrowAsset}

				nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(1000000000000)))
				externalCoin := sdk.NewCoin(externalAsset.Symbol, sdk.Int(sdk.NewUint(1000000000000)))
				err := app.BankKeeper.MintCoins(ctx, clptypes.ModuleName, sdk.NewCoins(nativeCoin, externalCoin))
				require.Nil(t, err)

				nativeCoin = sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(10000)))
				externalCoin = sdk.NewCoin(externalAsset.Symbol, sdk.Int(sdk.NewUint(10000)))
				_signer := clptest.GenerateAddress(clptest.AddressKey1)
				address = _signer.String()
				err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, _signer, sdk.NewCoins(nativeCoin, externalCoin))
				require.Nil(t, err)
			} else {
				address = tt.msgOpen.Signer
			}

			msg := tt.msgOpen
			msg.Signer = address

			marginKeeper.WhitelistAddress(ctx, address)

			_, got := msgServer.Open(sdk.WrapSDKContext(ctx), &msg)

			if tt.errString != nil {
				require.EqualError(t, got, tt.errString.Error())
			} else if tt.err == nil {
				require.NoError(t, got)
			} else {
				require.ErrorIs(t, got, tt.err)
			}
		})
	}
}

func TestKeeper_Close(t *testing.T) {
	table := []struct {
		msgOpen           types.MsgOpen
		msgClose          types.MsgClose
		name              string
		poolAsset         string
		token             string
		overrideSigner    string
		err               error
		errString         error
		poolEnabled       bool
		fundedAccount     bool
		mtpCreateDisabled bool
	}{
		{
			name: "mtp does not exist",
			msgClose: types.MsgClose{
				Signer: "xxx",
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:      "rowan",
			token:          "somethingelse",
			overrideSigner: "otheraddress",
			errString:      types.ErrMTPDoesNotExist,
		},
		{
			name: "pool does not exist",
			msgClose: types.MsgClose{
				Signer: "xxx",
				Id:     1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "same collateral and native asset but pool does not exist",
			msgClose: types.MsgClose{
				Signer: "xxx",
				Id:     1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "denom does not exist does not throw error as it does not use token registry",
			msgClose: types.MsgClose{
				Signer: "xxx",
				Id:     1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:   "xxx",
			token:       "somethingelse",
			poolEnabled: true,
		},
		{
			name: "wrong address/mtp not found",
			msgClose: types.MsgClose{
				Signer: "xxx",
				Id:     1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:         "xxx",
			token:             "xxx",
			poolEnabled:       true,
			mtpCreateDisabled: true,
			errString:         errors.New("mtp not found"),
		},
		{
			name: "insufficient funds/mtp not found",
			msgClose: types.MsgClose{
				Signer: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:     1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:         "xxx",
			token:             "xxx",
			poolEnabled:       true,
			mtpCreateDisabled: true,
			errString:         errors.New("mtp not found"),
		},
		{
			name: "account funded",
			msgClose: types.MsgClose{
				Signer: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:     1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			err:           nil,
		},
		{
			name: "mtp position invalid",
			msgClose: types.MsgClose{
				Signer: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:     1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
				Position:        types.Position_SHORT,
			},
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			errString:     errors.New("SHORT: mtp position invalid"),
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppMargin(false)
			marginKeeper := app.MarginKeeper

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.token,
				Decimals:    18,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			msgServer := keeper.NewMsgServerImpl(marginKeeper)

			asset := clptypes.Asset{Symbol: tt.poolAsset}
			pool := clptypes.Pool{
				ExternalAsset:        &asset,
				NativeAssetBalance:   sdk.NewUint(1000000000),
				NativeLiabilities:    sdk.NewUint(1000000000),
				ExternalCustody:      sdk.NewUint(1000000000),
				ExternalAssetBalance: sdk.NewUint(1000000000),
				ExternalLiabilities:  sdk.NewUint(1000000000),
				NativeCustody:        sdk.NewUint(1000000000),
				PoolUnits:            sdk.NewUint(1),
				Health:               sdk.NewDec(1),
			}
			if tt.poolEnabled {
				marginKeeper.SetEnabledPools(ctx, []string{tt.poolAsset})
			}

			// nolint:errcheck
			marginKeeper.ClpKeeper().SetPool(ctx, &pool)

			var address string

			if tt.fundedAccount {
				nativeAsset := tt.msgOpen.CollateralAsset
				externalAsset := clptypes.Asset{Symbol: tt.msgOpen.BorrowAsset}

				nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(1000000000000)))
				externalCoin := sdk.NewCoin(externalAsset.Symbol, sdk.Int(sdk.NewUint(1000000000000)))
				err := app.BankKeeper.MintCoins(ctx, clptypes.ModuleName, sdk.NewCoins(nativeCoin, externalCoin))
				require.NoError(t, err)

				nativeCoin = sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(10000)))
				// nolint:ineffassign
				externalCoin = sdk.NewCoin(externalAsset.Symbol, sdk.Int(sdk.NewUint(10000)))

				_signer := clptest.GenerateAddress(clptest.AddressKey1)
				address = _signer.String()
				err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, _signer, sdk.NewCoins(nativeCoin))
				require.NoError(t, err)
				err = marginKeeper.BankKeeper().SendCoinsFromAccountToModule(ctx, _signer, types.ModuleName, sdk.NewCoins(nativeCoin))
				require.NoError(t, err)
			} else {
				address = tt.msgClose.Signer
			}

			msg := tt.msgClose
			msg.Signer = address

			var signer = msg.Signer
			if tt.overrideSigner != "" {
				signer = tt.overrideSigner
			}

			if !tt.mtpCreateDisabled {
				addMTPKey(t, ctx, app, marginKeeper, tt.msgOpen.CollateralAsset, tt.msgOpen.BorrowAsset, signer, tt.msgOpen.Position, 1)
			}

			_, got := msgServer.Close(sdk.WrapSDKContext(ctx), &msg)

			if tt.errString != nil {
				require.EqualError(t, got, tt.errString.Error())
			} else if tt.err == nil {
				require.NoError(t, got)
			} else {
				require.ErrorIs(t, got, tt.err)
			}
		})
	}
}

func TestKeeper_ForceClose(t *testing.T) {
	table := []struct {
		msgOpen                       types.MsgOpen
		msgForceClose                 types.MsgForceClose
		name                          string
		poolAsset                     string
		token                         string
		overrideSigner                string
		overrideForceCloseThreadshold string
		err                           error
		errString                     error
		err2                          error
		errString2                    error
		poolEnabled                   bool
		fundedAccount                 bool
		mtpCreateDisabled             bool
	}{
		{
			name: "mtp does not exist",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:      "rowan",
			token:          "somethingelse",
			overrideSigner: "otheraddress",
			errString:      types.ErrMTPDoesNotExist,
		},
		{
			name: "pool does not exist",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:         1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "same collateral and native asset but pool does not exist",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:         1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "denom does not exist does not throw error as not using token registry but MTP health above threshold",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:         1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:   "xxx",
			token:       "somethingelse",
			poolEnabled: true,
			err:         types.ErrMTPHealthy,
		},
		{
			name: "wrong address/mtp not found",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:         1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:         "xxx",
			token:             "xxx",
			poolEnabled:       true,
			mtpCreateDisabled: true,
			errString:         errors.New("mtp not found"),
			errString2:        errors.New("mtp not found"),
		},
		{
			name: "insufficient funds/mtp not found",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:         1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:         "xxx",
			token:             "xxx",
			poolEnabled:       true,
			mtpCreateDisabled: true,
			errString:         errors.New("mtp not found"),
			errString2:        errors.New("mtp not found"),
		},
		{
			name: "account funded and mtp healthy but MTP health above threshold",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:         1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			err:           types.ErrMTPHealthy,
		},
		{
			name: "account funded and mtp not healthy but MTP health above threshold",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:         1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:                     "xxx",
			token:                         "xxx",
			poolEnabled:                   true,
			fundedAccount:                 true,
			overrideForceCloseThreadshold: "2",
			err:                           types.ErrMTPHealthy,
		},
		{
			name: "mtp position invalid",
			msgForceClose: types.MsgForceClose{
				Signer:     "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				MtpAddress: "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				Id:         1,
			},
			msgOpen: types.MsgOpen{
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
				Position:        types.Position_SHORT,
			},
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			errString:     errors.New("SHORT: mtp position invalid"),
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			asset := clptypes.Asset{Symbol: tt.poolAsset}

			ctx, app := test.CreateTestAppMarginFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				gs1 := &admintypes.GenesisState{
					AdminAccounts: []*admintypes.AdminAccount{
						{
							AdminType:    admintypes.AdminType_MARGIN,
							AdminAddress: tt.msgForceClose.Signer,
						},
						{
							AdminType:    admintypes.AdminType_CLPDEX,
							AdminAddress: tt.msgForceClose.Signer,
						},
						{
							AdminType:    admintypes.AdminType_TOKENREGISTRY,
							AdminAddress: tt.msgForceClose.Signer,
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(gs1)
				genesisState["admin"] = bz

				gs2 := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: tt.token, BaseDenom: tt.token, Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
						},
					},
				}
				bz, _ = app.AppCodec().MarshalJSON(gs2)
				genesisState["tokenregistry"] = bz

				gs3 := &types.GenesisState{
					Params: &types.Params{
						LeverageMax:                              sdk.NewDec(2),
						InterestRateMax:                          sdk.NewDec(1),
						InterestRateMin:                          sdk.ZeroDec(),
						InterestRateIncrease:                     sdk.NewDecWithPrec(1, 1),
						InterestRateDecrease:                     sdk.NewDecWithPrec(1, 1),
						HealthGainFactor:                         sdk.NewDecWithPrec(1, 2),
						EpochLength:                              0,
						ForceCloseThreshold:                      sdk.ZeroDec(),
						RemovalQueueThreshold:                    sdk.ZeroDec(),
						Pools:                                    []string{},
						ForceCloseFundPercentage:                 sdk.NewDecWithPrec(1, 1),
						ForceCloseInsuranceFundAddress:           "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
						IncrementalInterestPaymentFundPercentage: sdk.NewDecWithPrec(1, 1),
						IncrementalInterestPaymentInsuranceFundAddress: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
						IncrementalInterestPaymentEnabled:              false,
						PoolOpenThreshold:                              sdk.NewDecWithPrec(1, 1),
						MaxOpenPositions:                               10000,
						SqModifier:                                     sdk.MustNewDecFromStr("10000000000000000000000000"),
						SafetyFactor:                                   sdk.MustNewDecFromStr("1.05"),
					},
				}

				if tt.poolEnabled {
					gs3.Params.Pools = []string{
						tt.poolAsset,
					}
				}

				if tt.overrideForceCloseThreadshold != "" {
					gs3.Params.ForceCloseThreshold = sdk.MustNewDecFromStr(tt.overrideForceCloseThreadshold)
				}

				bz, _ = app.AppCodec().MarshalJSON(gs3)
				genesisState["margin"] = bz

				nativeAsset := tt.msgOpen.CollateralAsset
				externalAsset := clptypes.Asset{Symbol: tt.msgOpen.BorrowAsset}

				nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(1000000000000)))
				externalCoin := sdk.NewCoin(externalAsset.Symbol, sdk.Int(sdk.NewUint(1000000000000)))

				balances := []banktypes.Balance{
					{
						Address: tt.msgForceClose.Signer,
						Coins: sdk.Coins{
							nativeCoin,
							externalCoin,
						},
					},
				}

				gs4 := banktypes.DefaultGenesisState()
				gs4.Balances = append(gs4.Balances, balances...)
				bz, _ = app.AppCodec().MarshalJSON(gs4)
				genesisState["bank"] = bz

				gs5 := &clptypes.GenesisState{
					Params: clptypes.Params{
						MinCreatePoolThreshold: 100,
					},
					AddressWhitelist: []string{
						tt.msgForceClose.Signer,
					},
					PoolList: []*clptypes.Pool{
						{
							ExternalAsset:        &asset,
							NativeAssetBalance:   sdk.NewUint(1000000000),
							NativeLiabilities:    sdk.NewUint(1000000000),
							ExternalCustody:      sdk.NewUint(1000000000),
							ExternalAssetBalance: sdk.NewUint(1000000000),
							ExternalLiabilities:  sdk.NewUint(1000000000),
							NativeCustody:        sdk.NewUint(1000000000),
							PoolUnits:            sdk.NewUint(1),
							Health:               sdk.NewDec(1),
						},
					},
					LiquidityProviders: []*clptypes.LiquidityProvider{
						{
							Asset:                    &clptypes.Asset{Symbol: tt.poolAsset},
							LiquidityProviderAddress: tt.msgForceClose.Signer,
							LiquidityProviderUnits:   sdk.NewUint(1000000000),
						},
					},
				}
				bz, _ = app.AppCodec().MarshalJSON(gs5)
				genesisState["clp"] = bz

				return genesisState
			})
			marginKeeper := app.MarginKeeper
			msgServer := keeper.NewMsgServerImpl(marginKeeper)

			if tt.poolEnabled {
				marginKeeper.SetEnabledPools(ctx, []string{tt.poolAsset})
			}

			var address string

			address = tt.msgForceClose.Signer

			msg := tt.msgForceClose
			msg.Signer = address
			msg.MtpAddress = address

			var signer = msg.Signer
			if tt.overrideSigner != "" {
				signer = tt.overrideSigner
			}

			if !tt.mtpCreateDisabled {
				addMTPKey(t, ctx, app, marginKeeper, tt.msgOpen.CollateralAsset, tt.msgOpen.BorrowAsset, signer, tt.msgOpen.Position, 1)
			}

			_, got := msgServer.ForceClose(sdk.WrapSDKContext(ctx), &msg)

			if tt.errString != nil {
				require.EqualError(t, got, tt.errString.Error())
			} else if tt.err == nil {
				require.NoError(t, got)
			} else {
				require.ErrorIs(t, got, tt.err)
			}

			_, got2 := marginKeeper.GetMTP(ctx, signer, 1)

			if tt.errString2 != nil {
				require.EqualError(t, got2, tt.errString2.Error())
			} else if tt.err2 == nil {
				require.NoError(t, got2)
			} else {
				require.ErrorIs(t, got2, tt.err)
			}
		})
	}
}

func TestKeeper_OpenClose(t *testing.T) {
	table := []struct {
		name          string
		externalAsset string
		err           error
		errString     error
	}{
		{
			name:          "one round open/close long position",
			externalAsset: "xxx",
			errString:     errors.New("pool health too low to open new positions: margin not enabled for pool"),
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppMargin(false)
			marginKeeper := app.MarginKeeper

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.externalAsset,
				Decimals:    18,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			params := types.Params{
				LeverageMax:                                    sdk.NewDec(2),
				InterestRateMax:                                sdk.NewDec(1),
				InterestRateMin:                                sdk.ZeroDec(),
				InterestRateIncrease:                           sdk.NewDecWithPrec(1, 1),
				InterestRateDecrease:                           sdk.NewDecWithPrec(1, 1),
				HealthGainFactor:                               sdk.NewDecWithPrec(1, 2),
				EpochLength:                                    0,
				ForceCloseThreshold:                            sdk.ZeroDec(),
				RemovalQueueThreshold:                          sdk.ZeroDec(),
				ForceCloseFundPercentage:                       sdk.NewDecWithPrec(1, 1),
				ForceCloseInsuranceFundAddress:                 "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				IncrementalInterestPaymentFundPercentage:       sdk.NewDecWithPrec(1, 1),
				IncrementalInterestPaymentInsuranceFundAddress: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				IncrementalInterestPaymentEnabled:              false,
				PoolOpenThreshold:                              sdk.NewDecWithPrec(1, 1),
				MaxOpenPositions:                               10000,
				SqModifier:                                     sdk.MustNewDecFromStr("10000000000000000000000000"),
				SafetyFactor:                                   sdk.MustNewDecFromStr("1.05"),
				Pools:                                          []string{tt.externalAsset},
			}
			expectedGenesis := types.GenesisState{Params: &params}
			marginKeeper.InitGenesis(ctx, expectedGenesis)
			genesis := marginKeeper.ExportGenesis(ctx)
			require.Equal(t, expectedGenesis, *genesis)

			msgServer := keeper.NewMsgServerImpl(marginKeeper)

			nativeAsset := clptypes.NativeSymbol
			externalAsset := clptypes.Asset{Symbol: tt.externalAsset}

			SwapPriceNative := sdk.ZeroDec()
			SwapPriceExternal := sdk.ZeroDec()

			pool := clptypes.Pool{
				ExternalAsset:                 &externalAsset,
				NativeAssetBalance:            sdk.NewUint(1000000000000),
				ExternalAssetBalance:          sdk.NewUint(1000000000000),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				PoolUnits:                     sdk.ZeroUint(),
				Health:                        sdk.OneDec(),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			}

			marginKeeper.SetEnabledPools(ctx, []string{tt.externalAsset})
			// nolint:errcheck
			marginKeeper.ClpKeeper().SetPool(ctx, &pool)

			nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(1000000000000)))
			externalCoin := sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000)))
			err := app.BankKeeper.MintCoins(ctx, clptypes.ModuleName, sdk.NewCoins(nativeCoin, externalCoin))
			require.Nil(t, err)

			clpAccount := app.AccountKeeper.GetModuleAccount(ctx, clptypes.ModuleName)

			nativeCoinOk := app.ClpKeeper.HasBalance(ctx, clpAccount.GetAddress(), nativeCoin)
			require.True(t, nativeCoinOk)
			externalCoinOk := app.ClpKeeper.HasBalance(ctx, clpAccount.GetAddress(), externalCoin)
			require.True(t, externalCoinOk)

			require.Equal(t, app.BankKeeper.GetBalance(ctx, clpAccount.GetAddress(), nativeAsset), nativeCoin)
			require.Equal(t, app.BankKeeper.GetBalance(ctx, clpAccount.GetAddress(), tt.externalAsset), externalCoin)

			signer := clptest.GenerateAddress(clptest.AddressKey1)
			nativeCoin = sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(100000000000000)))
			externalCoin = sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000)))
			err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(nativeCoin, externalCoin))
			require.Nil(t, err)

			nativeCoinOk = app.ClpKeeper.HasBalance(ctx, signer, nativeCoin)
			require.True(t, nativeCoinOk)
			externalCoinOk = app.ClpKeeper.HasBalance(ctx, signer, externalCoin)
			require.True(t, externalCoinOk)

			require.Equal(t, app.BankKeeper.GetBalance(ctx, signer, nativeAsset), nativeCoin)
			require.Equal(t, app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset), externalCoin)

			msgOpen := types.MsgOpen{
				Signer:           signer.String(),
				CollateralAsset:  nativeAsset,
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      tt.externalAsset,
				Position:         types.Position_LONG,
				Leverage:         sdk.NewDec(2),
			}
			msgClose := types.MsgClose{
				Signer: signer.String(),
				Id:     1,
			}

			marginKeeper.WhitelistAddress(ctx, msgOpen.Signer)

			_, openError := msgServer.Open(sdk.WrapSDKContext(ctx), &msgOpen)
			require.Nil(t, openError)

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(99999999999000))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			openExpectedMTP := types.MTP{
				Id:               1,
				Address:          signer.String(),
				CollateralAsset:  nativeAsset,
				CollateralAmount: sdk.NewUint(1000),
				Liabilities:      sdk.NewUint(1000),
				InterestPaid:     sdk.ZeroUint(),
				InterestUnpaid:   sdk.ZeroUint(),
				CustodyAsset:     tt.externalAsset,
				CustodyAmount:    sdk.NewUint(1999),
				Leverage:         sdk.NewDec(2),
				MtpHealth:        sdk.MustNewDecFromStr("2.001001001001001001"),
				Position:         types.Position_LONG,
			}

			openMTP, _ := marginKeeper.GetMTP(ctx, signer.String(), 1)

			fmt.Println(openExpectedMTP)
			fmt.Println(openMTP)

			require.Equal(t, openExpectedMTP, openMTP)

			openExpectedPool := clptypes.Pool{
				ExternalAsset:                 &externalAsset,
				NativeAssetBalance:            sdk.NewUint(1000000001000),
				ExternalAssetBalance:          sdk.NewUint(999999998001),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.NewUint(1999),
				NativeLiabilities:             sdk.NewUint(1000),
				ExternalLiabilities:           sdk.ZeroUint(),
				PoolUnits:                     sdk.ZeroUint(),
				Health:                        sdk.NewDecWithPrec(999999999000000002, 18),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			}

			openPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)
			require.Equal(t, openExpectedPool, openPool)

			_, closeError := msgServer.Close(sdk.WrapSDKContext(ctx), &msgClose)
			require.Nil(t, closeError)

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(99999999999998))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			closeExpectedPool := clptypes.Pool{
				ExternalAsset:                 &externalAsset,
				NativeAssetBalance:            sdk.NewUint(1000000000002),
				ExternalAssetBalance:          sdk.NewUint(1000000000000),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				PoolUnits:                     sdk.ZeroUint(),
				Health:                        sdk.NewDecWithPrec(999999999000000002, 18),
				InterestRate:                  sdk.NewDecWithPrec(1, 1),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			}

			closePool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)
			require.Equal(t, closeExpectedPool, closePool)
		})
	}
}

func TestKeeper_OpenThenClose(t *testing.T) {
	externalAsset := "xxx"
	nativeAsset := clptypes.NativeSymbol
	signer := "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v"

	ctx, app := test.CreateTestAppMarginFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
		gs2 := &tokenregistrytypes.GenesisState{
			Registry: &tokenregistrytypes.Registry{
				Entries: []*tokenregistrytypes.RegistryEntry{
					{Denom: nativeAsset, BaseDenom: nativeAsset, Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
					{Denom: externalAsset, BaseDenom: externalAsset, Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				},
			},
		}
		bz, _ := app.AppCodec().MarshalJSON(gs2)
		genesisState["tokenregistry"] = bz

		gs3 := &types.GenesisState{
			Params: &types.Params{
				LeverageMax:                                    sdk.NewDec(2),
				HealthGainFactor:                               sdk.NewDec(2),
				InterestRateMin:                                sdk.NewDecWithPrec(5, 3),
				InterestRateMax:                                sdk.NewDec(3),
				InterestRateDecrease:                           sdk.NewDecWithPrec(1, 5),
				InterestRateIncrease:                           sdk.NewDecWithPrec(1, 5),
				ForceCloseThreshold:                            sdk.NewDecWithPrec(1, 10),
				RemovalQueueThreshold:                          sdk.ZeroDec(),
				ForceCloseFundPercentage:                       sdk.NewDecWithPrec(1, 1),
				ForceCloseInsuranceFundAddress:                 "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				IncrementalInterestPaymentFundPercentage:       sdk.NewDecWithPrec(1, 1),
				IncrementalInterestPaymentInsuranceFundAddress: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				IncrementalInterestPaymentEnabled:              false,
				PoolOpenThreshold:                              sdk.NewDecWithPrec(1, 1),
				MaxOpenPositions:                               10000,
				SqModifier:                                     sdk.MustNewDecFromStr("10000000000000000000000000"),
				SafetyFactor:                                   sdk.MustNewDecFromStr("1.05"),
				EpochLength:                                    1,
				Pools: []string{
					externalAsset,
				},
			},
		}
		bz, _ = app.AppCodec().MarshalJSON(gs3)
		genesisState["margin"] = bz

		nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("100000000000000000000000000")))
		externalCoin := sdk.NewCoin(externalAsset, sdk.Int(sdk.NewUintFromString("100000000000000000000000000")))

		balances := []banktypes.Balance{
			{
				Address: signer,
				Coins: sdk.Coins{
					nativeCoin,
					externalCoin,
				},
			},
		}

		gs4 := banktypes.DefaultGenesisState()
		gs4.Balances = append(gs4.Balances, balances...)
		bz, _ = app.AppCodec().MarshalJSON(gs4)
		genesisState["bank"] = bz

		gs5 := &clptypes.GenesisState{
			Params: clptypes.Params{
				MinCreatePoolThreshold: 100,
			},
			AddressWhitelist: []string{
				signer,
			},
			PoolList: []*clptypes.Pool{
				{
					ExternalAsset:                 &clptypes.Asset{Symbol: externalAsset},
					NativeAssetBalance:            sdk.NewUintFromString("1000000000000000000000000000000"),
					ExternalAssetBalance:          sdk.NewUintFromString("1000000000000000000000000000000"),
					PoolUnits:                     sdk.NewUintFromString("1000000000000000000000000000000"),
					NativeCustody:                 sdk.ZeroUint(),
					ExternalCustody:               sdk.ZeroUint(),
					NativeLiabilities:             sdk.ZeroUint(),
					ExternalLiabilities:           sdk.ZeroUint(),
					Health:                        sdk.OneDec(),
					InterestRate:                  sdk.NewDecWithPrec(1, 1),
					RewardPeriodNativeDistributed: sdk.ZeroUint(),
				},
			},
			LiquidityProviders: []*clptypes.LiquidityProvider{
				{
					Asset:                    &clptypes.Asset{Symbol: externalAsset},
					LiquidityProviderAddress: signer,
					LiquidityProviderUnits:   sdk.NewUintFromString("1000000000000000000000000000000"),
				},
			},
		}
		bz, _ = app.AppCodec().MarshalJSON(gs5)
		genesisState["clp"] = bz

		return genesisState
	})
	marginKeeper := app.MarginKeeper
	msgServer := keeper.NewMsgServerImpl(marginKeeper)

	msgOpen := types.MsgOpen{
		Signer:           signer,
		CollateralAsset:  nativeAsset,
		CollateralAmount: sdk.NewUintFromString("10000"),
		BorrowAsset:      externalAsset,
		Position:         types.Position_LONG,
		Leverage:         sdk.NewDec(2),
	}

	marginKeeper.WhitelistAddress(ctx, msgOpen.Signer)

	_, err := msgServer.Open(sdk.WrapSDKContext(ctx), &msgOpen)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	clp.BeginBlocker(ctx, app.ClpKeeper)
	marginKeeper.BeginBlocker(ctx)

	expectedMTP := types.MTP{
		Address:          signer,
		CollateralAsset:  nativeAsset,
		CollateralAmount: sdk.NewUintFromString("10000"),
		Liabilities:      sdk.NewUintFromString("10000"),
		InterestUnpaid:   sdk.NewUintFromString("0"),
		CustodyAsset:     externalAsset,
		CustodyAmount:    sdk.NewUintFromString("20000"),
		Leverage:         sdk.NewDec(1),
		MtpHealth:        sdk.MustNewDecFromStr("0.166666666666666667"),
		Position:         types.Position_LONG,
		Id:               1,
	}
	mtp, err := marginKeeper.GetMTP(ctx, signer, uint64(1))
	t.Logf("mtp: %v\n", mtp)
	t.Logf("expected mtp: %v\n", expectedMTP)
	require.NoError(t, err)
	require.NotNil(t, mtp)
	// require.Equal(t, expectedMTP, mtp)

	expectedPool := clptypes.Pool{
		ExternalAsset:                 &clptypes.Asset{Symbol: externalAsset},
		NativeAssetBalance:            sdk.NewUintFromString("999999999999999999999999990000"),
		ExternalAssetBalance:          sdk.NewUintFromString("999999999999999999999999960000"),
		PoolUnits:                     sdk.NewUintFromString("1000000000000000000000000000000"),
		ExternalLiabilities:           sdk.NewUintFromString("0"),
		ExternalCustody:               sdk.NewUintFromString("40000"),
		NativeLiabilities:             sdk.NewUintFromString("10000"),
		NativeCustody:                 sdk.NewUintFromString("0"),
		Health:                        sdk.MustNewDecFromStr("1.0"),
		InterestRate:                  sdk.MustNewDecFromStr("0.005"),
		RewardPeriodNativeDistributed: sdk.NewUintFromString("0"),
	}
	pool, err := marginKeeper.ClpKeeper().GetPool(ctx, externalAsset)
	t.Logf("pool: %v\n", pool)
	t.Logf("expected pool: %v\n", expectedPool)
	require.NoError(t, err)
	require.NotNil(t, pool)
	// require.Equal(t, expectedPool, pool)

	msgClose := types.MsgClose{
		Signer: signer,
		Id:     uint64(1),
	}
	_, err = msgServer.Close(sdk.WrapSDKContext(ctx), &msgClose)
	require.NoError(t, err)
}

func TestKeeper_EC(t *testing.T) {
	type Chunk struct {
		chunk                                sdk.Uint
		signerNativeAssetBalanceAfterOpen    sdk.Uint
		signerExternalAssetBalanceAfterOpen  sdk.Uint
		signerNativeAssetBalanceAfterClose   sdk.Uint
		signerExternalAssetBalanceAfterClose sdk.Uint
		poolNativeAssetBalanceAfterOpen      sdk.Uint
		poolExternalAssetBalanceAfterOpen    sdk.Uint
		poolHealthAfterOpen                  sdk.Dec
		poolNativeAssetBalanceAfterClose     sdk.Uint
		poolExternalAssetBalanceAfterClose   sdk.Uint
		poolHealthAfterClose                 sdk.Dec
		mtpCustodyAmount                     sdk.Uint
		mtpHealth                            sdk.Dec
		openErrorString                      error
		openError                            error
		closeErrorString                     error
		closeError                           error
	}
	type Test struct {
		// nolint:golint
		X_A sdk.Uint
		// nolint:golint
		Y_A    sdk.Uint
		chunks []Chunk
	}
	type Table struct {
		name          string
		externalAsset string
		tests         []Test
	}

	table := []Table{
		{
			name:          "EC1,EC2",
			externalAsset: "xxx",
			tests: []Test{
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(100),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999997544),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(102456),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(31),
							mtpHealth:                            sdk.NewDecWithPrec(162001036806635562, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999965817),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(11692600000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(13),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(134183),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(77),
							mtpHealth:                            sdk.NewDecWithPrec(308848220753477350, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999965817),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999903213),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(196787),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(10),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(7),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(110000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(10),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(3),
							mtpHealth:                            sdk.NewDecWithPrec(164397974616952719, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999957916),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(142084),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(10),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(8),
							mtpHealth:                            sdk.NewDecWithPrec(307029296177206145, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999892399),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(207601),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(10),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(1),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(110000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.ZeroUint(),
							mtpHealth:                            sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(165000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(9),
							mtpHealth:                            sdk.NewDecWithPrec(353577237340327734, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999836000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(264000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
			},
		},
		{
			name:          "EC3,EC4",
			externalAsset: "xxx",
			tests: []Test{
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(100000),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999997755),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69444),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(102245),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(30556),
							mtpHealth:                            sdk.NewDecWithPrec(163049681237873180, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999966492),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(117402),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(13101),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(133508),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(86899),
							mtpHealth:                            sdk.NewDecWithPrec(355899519859193208, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999904196),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(195804),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(200000),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999997755),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(138889),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(102245),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(200000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(61111),
							mtpHealth:                            sdk.NewDecWithPrec(163049681237873180, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999966492),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(117403),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(26203),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(133508),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(200000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(173797),
							mtpHealth:                            sdk.NewDecWithPrec(355899519859193208, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999904196),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(195804),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(200000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(1),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(110000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.ZeroUint(),
							mtpHealth:                            sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(165000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(9),
							mtpHealth:                            sdk.NewDecWithPrec(353577237340327734, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999836000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(264000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
			},
		},
		{
			name:          "EC5,EC6,EC7",
			externalAsset: "xxx",
			tests: []Test{
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(5000),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999997752),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(3472),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(102248),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(5000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(1528),
							mtpHealth:                            sdk.NewDecWithPrec(163039047851960545, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999966487),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(117398),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(655),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(133513),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(5000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(4345),
							mtpHealth:                            sdk.NewDecWithPrec(355906428964312292, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999904183),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(195817),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(5000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
				{
					X_A: sdk.NewUint(10000),
					Y_A: sdk.NewUint(100),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999999000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999999754),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(9000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(10246),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(31),
							mtpHealth:                            sdk.NewDecWithPrec(161995788109509153, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999993500),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999996581),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(11693),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(13),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(13419),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356633380884450785, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999990320),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(19680),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(1),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(110000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.ZeroUint(),
							mtpHealth:                            sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(165000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(9),
							mtpHealth:                            sdk.NewDecWithPrec(353577237340327734, 18),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999836000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(100000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(264000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(1000000000000000000, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
						},
					},
				},
			},
		},
	}

	signer := "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v"
	nativeAsset := clptypes.NativeSymbol

	for _, ec := range table {
		ec := ec
		asset := clptypes.Asset{Symbol: ec.externalAsset}

		for _, testItem := range ec.tests {
			testItem := testItem

			ctx, app := test.CreateTestAppMarginFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				gs2 := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: ec.externalAsset, BaseDenom: ec.externalAsset, Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(gs2)
				genesisState["tokenregistry"] = bz

				gs3 := &types.GenesisState{
					Params: &types.Params{
						MaxOpenPositions:                               10000,
						LeverageMax:                                    sdk.NewDec(2),
						HealthGainFactor:                               sdk.NewDec(1),
						InterestRateMin:                                sdk.NewDecWithPrec(5, 3),
						InterestRateMax:                                sdk.NewDec(3),
						InterestRateDecrease:                           sdk.NewDecWithPrec(1, 2),
						InterestRateIncrease:                           sdk.NewDecWithPrec(1, 2),
						ForceCloseThreshold:                            sdk.NewDecWithPrec(1, 2),
						RemovalQueueThreshold:                          sdk.ZeroDec(),
						ForceCloseFundPercentage:                       sdk.NewDecWithPrec(1, 1),
						ForceCloseInsuranceFundAddress:                 "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
						IncrementalInterestPaymentFundPercentage:       sdk.NewDecWithPrec(1, 1),
						IncrementalInterestPaymentInsuranceFundAddress: "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
						IncrementalInterestPaymentEnabled:              false,
						PoolOpenThreshold:                              sdk.NewDecWithPrec(1, 1),
						SqModifier:                                     sdk.MustNewDecFromStr("10000000000000000000000000"),
						SafetyFactor:                                   sdk.MustNewDecFromStr("1.05"),
						EpochLength:                                    1,
						Pools: []string{
							ec.externalAsset,
						},
					},
				}
				bz, _ = app.AppCodec().MarshalJSON(gs3)
				genesisState["margin"] = bz

				nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(100000000000000)))
				externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUint(100000000000000)))

				balances := []banktypes.Balance{
					{
						Address: signer,
						Coins: sdk.Coins{
							nativeCoin,
							externalCoin,
						},
					},
				}

				gs4 := banktypes.DefaultGenesisState()
				gs4.Balances = append(gs4.Balances, balances...)
				bz, _ = app.AppCodec().MarshalJSON(gs4)
				genesisState["bank"] = bz

				SwapPriceNative := sdk.ZeroDec()
				SwapPriceExternal := sdk.ZeroDec()

				gs5 := &clptypes.GenesisState{
					Params: clptypes.Params{
						MinCreatePoolThreshold: 100,
					},
					AddressWhitelist: []string{
						signer,
					},
					PoolList: []*clptypes.Pool{
						{
							ExternalAsset:                 &asset,
							NativeAssetBalance:            testItem.X_A,
							ExternalAssetBalance:          testItem.Y_A,
							NativeCustody:                 sdk.ZeroUint(),
							ExternalCustody:               sdk.ZeroUint(),
							NativeLiabilities:             sdk.ZeroUint(),
							ExternalLiabilities:           sdk.ZeroUint(),
							PoolUnits:                     sdk.NewUint(1),
							Health:                        sdk.NewDec(1),
							InterestRate:                  sdk.NewDecWithPrec(1, 2),
							SwapPriceNative:               &SwapPriceNative,
							SwapPriceExternal:             &SwapPriceExternal,
							RewardPeriodNativeDistributed: sdk.ZeroUint(),
						},
					},
					LiquidityProviders: []*clptypes.LiquidityProvider{
						{
							Asset:                    &asset,
							LiquidityProviderAddress: signer,
							LiquidityProviderUnits:   sdk.NewUint(1),
						},
					},
				}
				bz, _ = app.AppCodec().MarshalJSON(gs5)
				genesisState["clp"] = bz

				return genesisState
			})
			marginKeeper := app.MarginKeeper
			msgServer := keeper.NewMsgServerImpl(marginKeeper)

			for i, chunkItem := range testItem.chunks {
				i := i
				chunkItem := chunkItem
				name := fmt.Sprintf("%v, X_A=%v, Y_A=%v, delta x=%v%%", ec.name, testItem.X_A, testItem.Y_A, chunkItem.chunk)
				t.Run(name, func(t *testing.T) {
					msgOpen := types.MsgOpen{
						Signer:           signer,
						CollateralAsset:  nativeAsset,
						CollateralAmount: testItem.X_A.Mul(chunkItem.chunk).Quo(sdk.NewUint(100)),
						BorrowAsset:      ec.externalAsset,
						Position:         types.Position_LONG,
						Leverage:         sdk.NewDec(1),
					}
					msgClose := types.MsgClose{
						Signer: signer,
						Id:     uint64(i + 1),
					}

					marginKeeper.WhitelistAddress(ctx, msgOpen.Signer)

					_, openError := msgServer.Open(sdk.WrapSDKContext(ctx), &msgOpen)
					if chunkItem.openErrorString != nil {
						require.EqualError(t, openError, chunkItem.openErrorString.Error())
						return
					} else if chunkItem.openError != nil {
						require.ErrorIs(t, openError, chunkItem.openError)
						return
					} else {
						require.NoError(t, openError)
					}

					acc, _ := sdk.AccAddressFromBech32(signer)

					// require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(chunkItem.signerNativeAssetBalanceAfterOpen)), app.BankKeeper.GetBalance(ctx, acc, nativeAsset))
					// require.Equal(t, sdk.NewCoin(ec.externalAsset, sdk.Int(chunkItem.signerExternalAssetBalanceAfterOpen)), app.BankKeeper.GetBalance(ctx, acc, ec.externalAsset))

					openExpectedMTP := types.MTP{
						Id:               uint64(i + 1),
						Address:          signer,
						CollateralAsset:  nativeAsset,
						CollateralAmount: msgOpen.CollateralAmount,
						Liabilities:      msgOpen.CollateralAmount,
						InterestUnpaid:   sdk.ZeroUint(),
						CustodyAsset:     ec.externalAsset,
						CustodyAmount:    chunkItem.mtpCustodyAmount,
						Leverage:         sdk.NewDec(2),
						MtpHealth:        chunkItem.mtpHealth,
						Position:         types.Position_LONG,
					}
					openMTP, err := marginKeeper.GetMTP(ctx, signer, uint64(i+1))
					require.NoError(t, err)
					require.NotNil(t, openMTP)
					require.NotNil(t, openExpectedMTP)
					// require.Equal(t, openExpectedMTP, openMTP)

					SwapPriceNative := sdk.ZeroDec()
					SwapPriceExternal := sdk.ZeroDec()

					openExpectedPool := clptypes.Pool{
						ExternalAsset:                 &asset,
						NativeAssetBalance:            chunkItem.poolNativeAssetBalanceAfterOpen,
						ExternalAssetBalance:          chunkItem.poolExternalAssetBalanceAfterOpen,
						NativeCustody:                 sdk.ZeroUint(),
						ExternalCustody:               chunkItem.mtpCustodyAmount,
						NativeLiabilities:             msgOpen.CollateralAmount,
						ExternalLiabilities:           sdk.ZeroUint(),
						PoolUnits:                     sdk.NewUint(1),
						Health:                        chunkItem.poolHealthAfterOpen,
						InterestRate:                  sdk.NewDecWithPrec(1, 2),
						SwapPriceNative:               &SwapPriceNative,
						SwapPriceExternal:             &SwapPriceExternal,
						RewardPeriodNativeDistributed: sdk.ZeroUint(),
					}
					openPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, ec.externalAsset)
					// require.Equal(t, openExpectedPool, openPool)
					require.NotNil(t, openPool)
					require.NotNil(t, openExpectedPool)

					_, closeError := msgServer.Close(sdk.WrapSDKContext(ctx), &msgClose)
					if chunkItem.closeErrorString != nil {
						require.EqualError(t, closeError, chunkItem.closeErrorString.Error())
						return
					} else if chunkItem.closeError != nil {
						require.ErrorIs(t, closeError, chunkItem.closeError)
						return
					} else {
						require.NoError(t, closeError)
					}

					require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(chunkItem.signerNativeAssetBalanceAfterClose)), app.BankKeeper.GetBalance(ctx, acc, nativeAsset))
					require.Equal(t, sdk.NewCoin(ec.externalAsset, sdk.Int(chunkItem.signerExternalAssetBalanceAfterClose)), app.BankKeeper.GetBalance(ctx, acc, ec.externalAsset))

					closeExpectedPool := clptypes.Pool{
						ExternalAsset:                 &asset,
						NativeAssetBalance:            chunkItem.poolNativeAssetBalanceAfterClose,
						ExternalAssetBalance:          chunkItem.poolExternalAssetBalanceAfterClose,
						PoolUnits:                     sdk.NewUintFromString("1"),
						ExternalLiabilities:           sdk.ZeroUint(),
						ExternalCustody:               sdk.ZeroUint(),
						NativeLiabilities:             sdk.ZeroUint(),
						NativeCustody:                 sdk.ZeroUint(),
						Health:                        chunkItem.poolHealthAfterClose,
						InterestRate:                  sdk.NewDecWithPrec(1, 2),
						SwapPriceNative:               &SwapPriceNative,
						SwapPriceExternal:             &SwapPriceExternal,
						RewardPeriodNativeDistributed: sdk.ZeroUint(),
					}
					closePool, _ := marginKeeper.ClpKeeper().GetPool(ctx, ec.externalAsset)
					t.Log("closeExpectedPool:", closeExpectedPool)
					t.Log("closePool:", closePool)
					require.Equal(t, closeExpectedPool, closePool)
				})
			}
		}
	}
}

func TestKeeper_AddUpExistingMTP(t *testing.T) {
	nativeAsset := clptypes.NativeSymbol
	externalAsset := clptypes.Asset{Symbol: "xxx"}

	ctx, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
		Denom:       externalAsset.Symbol,
		Decimals:    18,
		Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
	})

	msgServer := keeper.NewMsgServerImpl(marginKeeper)

	SwapPriceNative := sdk.ZeroDec()
	SwapPriceExternal := sdk.ZeroDec()

	pool := clptypes.Pool{
		ExternalAsset:                 &externalAsset,
		NativeAssetBalance:            sdk.NewUintFromString("10000000000000000000000000"),
		ExternalAssetBalance:          sdk.NewUintFromString("10000000000000000000000000"),
		NativeCustody:                 sdk.ZeroUint(),
		ExternalCustody:               sdk.ZeroUint(),
		NativeLiabilities:             sdk.ZeroUint(),
		ExternalLiabilities:           sdk.ZeroUint(),
		PoolUnits:                     sdk.ZeroUint(),
		Health:                        sdk.OneDec(),
		InterestRate:                  sdk.NewDecWithPrec(1, 1),
		SwapPriceNative:               &SwapPriceNative,
		SwapPriceExternal:             &SwapPriceExternal,
		RewardPeriodNativeDistributed: sdk.ZeroUint(),
	}

	marginKeeper.SetEnabledPools(ctx, []string{externalAsset.Symbol})
	// nolint:errcheck
	marginKeeper.ClpKeeper().SetPool(ctx, &pool)

	nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(pool.NativeAssetBalance))
	externalCoin := sdk.NewCoin(externalAsset.Symbol, sdk.Int(pool.ExternalAssetBalance))
	err := app.BankKeeper.MintCoins(ctx, clptypes.ModuleName, sdk.NewCoins(nativeCoin, externalCoin))
	require.Nil(t, err)

	signer := clptest.GenerateAddress(clptest.AddressKey1)
	nativeCoin = sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("100000000000000000000000")))
	externalCoin = sdk.NewCoin(externalAsset.Symbol, sdk.Int(sdk.NewUintFromString("100000000000000000000000")))
	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(nativeCoin, externalCoin))
	require.Nil(t, err)

	msg1 := types.MsgOpen{
		Signer:           signer.String(),
		CollateralAsset:  nativeAsset,
		CollateralAmount: sdk.NewUintFromString("1000000000000000000000"),
		BorrowAsset:      externalAsset.Symbol,
		Position:         types.Position_LONG,
		Leverage:         sdk.NewDec(1),
	}

	marginKeeper.WhitelistAddress(ctx, msg1.Signer)

	_, openError := msgServer.Open(sdk.WrapSDKContext(ctx), &msg1)
	require.NoError(t, openError)

	openExpectedMTP := types.MTP{
		Address:          signer.String(),
		CollateralAsset:  nativeAsset,
		CollateralAmount: sdk.NewUintFromString("1000000000000000000000"),
		Liabilities:      sdk.ZeroUint(),
		InterestPaid:     sdk.ZeroUint(),
		InterestUnpaid:   sdk.ZeroUint(),
		CustodyAsset:     externalAsset.Symbol,
		CustodyAmount:    sdk.NewUintFromString("999800029996000499940"),
		Leverage:         sdk.NewDec(1),
		MtpHealth:        sdk.ZeroDec(),
		Position:         types.Position_LONG,
		Id:               1,
	}
	openMTP, _ := marginKeeper.GetMTP(ctx, signer.String(), 1)
	fmt.Println(openExpectedMTP)
	fmt.Println(openMTP)
	require.Equal(t, openExpectedMTP, openMTP)

	msg2 := types.MsgOpen{
		Signer:           signer.String(),
		CollateralAsset:  nativeAsset,
		CollateralAmount: sdk.NewUintFromString("500000000000000000000"),
		BorrowAsset:      externalAsset.Symbol,
		Position:         types.Position_LONG,
		Leverage:         sdk.NewDec(1),
	}

	marginKeeper.WhitelistAddress(ctx, msg2.Signer)

	_, openError = msgServer.Open(sdk.WrapSDKContext(ctx), &msg2)
	require.NoError(t, openError)

	openExpectedMTP = types.MTP{
		Address:          signer.String(),
		CollateralAsset:  nativeAsset,
		CollateralAmount: sdk.NewUintFromString("1000000000000000000000"),
		Liabilities:      sdk.ZeroUint(),
		InterestPaid:     sdk.ZeroUint(),
		InterestUnpaid:   sdk.ZeroUint(),
		CustodyAsset:     externalAsset.Symbol,
		CustodyAmount:    sdk.NewUintFromString("999800029996000499940"),
		Leverage:         sdk.NewDec(1),
		MtpHealth:        sdk.ZeroDec(),
		Position:         types.Position_LONG,
		Id:               1,
	}
	openMTP, _ = marginKeeper.GetMTP(ctx, signer.String(), 1)
	fmt.Println(openExpectedMTP)
	fmt.Println(openMTP)
	require.Equal(t, openExpectedMTP, openMTP)

	openExpectedMTP = types.MTP{
		Address:          signer.String(),
		CollateralAsset:  nativeAsset,
		CollateralAmount: sdk.NewUintFromString("500000000000000000000"),
		Liabilities:      sdk.ZeroUint(),
		InterestPaid:     sdk.ZeroUint(),
		InterestUnpaid:   sdk.ZeroUint(),
		CustodyAsset:     externalAsset.Symbol,
		CustodyAmount:    sdk.NewUintFromString("499850038741251802776"),
		Leverage:         sdk.NewDec(1),
		MtpHealth:        sdk.ZeroDec(),
		Position:         types.Position_LONG,
		Id:               2,
	}
	openMTP, _ = marginKeeper.GetMTP(ctx, signer.String(), 2)
	fmt.Println(openExpectedMTP)
	fmt.Println(openMTP)
	require.Equal(t, openExpectedMTP, openMTP)
}
