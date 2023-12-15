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
			// errString:   errors.New("external balance mismatch in pool xxx (module: 0 != pool: 1000001000): Balance of module account check failed"),
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
			// errString:     errors.New("external balance mismatch in pool xxx (module: 1000000000000 != pool: 1000001000): Balance of module account check failed"),
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
				// nolint:staticcheck,ineffassign
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
				addMTPKey(t, ctx, app, marginKeeper, tt.msgOpen.CollateralAsset, tt.msgOpen.BorrowAsset, signer, tt.msgOpen.Position, 1, sdk.NewDec(20))
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
		msgOpen           types.MsgOpen
		msgForceClose     types.MsgForceClose
		health            sdk.Dec
		name              string
		poolAsset         string
		token             string
		overrideSigner    string
		err               error
		errString         error
		err2              error
		errString2        error
		poolEnabled       bool
		fundedAccount     bool
		mtpCreateDisabled bool
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
			health:         sdk.NewDec(20),
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
			health:    sdk.NewDec(20),
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
			health:    sdk.NewDec(20),
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
			health:      sdk.NewDec(20),
			poolAsset:   "xxx",
			token:       "somethingelse",
			poolEnabled: true,
			// errString:   errors.New("external balance mismatch in pool xxx (module: 0 != pool: 1000001000): Balance of module account check failed"),
			err2: types.ErrMTPDoesNotExist,
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
			health:            sdk.NewDec(20),
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
			health:            sdk.NewDec(20),
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
			health:        sdk.NewDec(20),
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			// errString:     errors.New("external balance mismatch in pool xxx (module: 0 != pool: 1000001000): Balance of module account check failed"),
			err2: types.ErrMTPDoesNotExist,
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
			health:        sdk.NewDec(20),
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			// errString:     errors.New("external balance mismatch in pool xxx (module: 0 != pool: 1000001000): Balance of module account check failed"),
			err2: types.ErrMTPDoesNotExist,
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
			health:        sdk.NewDec(20),
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			errString:     errors.New("SHORT: mtp position invalid"),
		},
		{
			name: "admin closure does not check health",
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
			health:        sdk.NewDecWithPrec(1, 2),
			poolAsset:     "xxx",
			token:         "xxx",
			poolEnabled:   true,
			fundedAccount: true,
			// errString:     errors.New("external balance mismatch in pool xxx (module: 0 != pool: 1000001000): Balance of module account check failed"),
			err2: types.ErrMTPDoesNotExist,
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
						RemovalQueueThreshold:                    sdk.ZeroDec(),
						Pools:                                    []string{},
						ForceCloseFundPercentage:                 sdk.NewDecWithPrec(1, 1),
						ForceCloseFundAddress:                    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
						IncrementalInterestPaymentFundPercentage: sdk.NewDecWithPrec(1, 1),
						IncrementalInterestPaymentFundAddress:    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
						IncrementalInterestPaymentEnabled:        false,
						PoolOpenThreshold:                        sdk.NewDecWithPrec(1, 1),
						MaxOpenPositions:                         10000,
						SqModifier:                               sdk.MustNewDecFromStr("10000000000000000000000000"),
						SafetyFactor:                             sdk.MustNewDecFromStr("1.05"),
						RowanCollateralEnabled:                   true,
					},
				}

				if tt.poolEnabled {
					gs3.Params.Pools = []string{
						tt.poolAsset,
					}
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
				addMTPKey(t, ctx, app, marginKeeper, tt.msgOpen.CollateralAsset, tt.msgOpen.BorrowAsset, signer, tt.msgOpen.Position, 1, sdk.NewDec(20))
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
				require.ErrorIs(t, got2, tt.err2)
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
			// errString:     errors.New("pool health too low to open new positions: margin not enabled for pool"),
			// errString: errors.New("rowan: using rowan as collateral asset is not allowed"),
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppMargin(false)
			marginKeeper := app.MarginKeeper

			app.ClpKeeper.SetSwapFeeParams(ctx, clptypes.GetDefaultSwapFeeParams())

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.externalAsset,
				Decimals:    18,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			params := types.Params{
				LeverageMax:                              sdk.NewDec(2),
				InterestRateMax:                          sdk.NewDec(1),
				InterestRateMin:                          sdk.NewDecWithPrec(1, 1),
				InterestRateIncrease:                     sdk.NewDecWithPrec(1, 1),
				InterestRateDecrease:                     sdk.NewDecWithPrec(1, 1),
				HealthGainFactor:                         sdk.NewDecWithPrec(1, 2),
				EpochLength:                              0,
				RemovalQueueThreshold:                    sdk.ZeroDec(),
				ForceCloseFundPercentage:                 sdk.NewDecWithPrec(1, 1),
				ForceCloseFundAddress:                    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				IncrementalInterestPaymentFundPercentage: sdk.NewDecWithPrec(1, 1),
				IncrementalInterestPaymentFundAddress:    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				IncrementalInterestPaymentEnabled:        false,
				PoolOpenThreshold:                        sdk.NewDecWithPrec(1, 1),
				MaxOpenPositions:                         10000,
				SqModifier:                               sdk.MustNewDecFromStr("10000000000000000000000000"),
				SafetyFactor:                             sdk.MustNewDecFromStr("1.05"),
				Pools:                                    []string{tt.externalAsset},
				RowanCollateralEnabled:                   true,
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
				NativeAssetBalance:            sdk.NewUintFromString("100000000000000000000000000"),
				ExternalAssetBalance:          sdk.NewUintFromString("100000000000000000000000000"),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
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

			nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("100000000000000000000000000")))
			externalCoin := sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUintFromString("100000000000000000000000000")))
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
			nativeCoin = sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("2000000000000000000000")))
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
				CollateralAmount: sdk.NewUintFromString("1000000000000000000000"),
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

			if tt.errString != nil {
				require.EqualError(t, openError, tt.errString.Error())
			} else if tt.err == nil {
				require.NoError(t, openError)
			} else {
				require.ErrorIs(t, openError, tt.err)
			}

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("1000000000000000000000"))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			openExpectedMTP := types.MTP{
				Id:                       1,
				Address:                  signer.String(),
				CollateralAsset:          nativeAsset,
				CollateralAmount:         sdk.NewUintFromString("1000000000000000000000"),
				Liabilities:              sdk.NewUintFromString("1000000000000000000000"),
				InterestPaidCollateral:   sdk.ZeroUint(),
				InterestPaidCustody:      sdk.ZeroUint(),
				InterestUnpaidCollateral: sdk.ZeroUint(),
				CustodyAsset:             tt.externalAsset,
				CustodyAmount:            sdk.NewUintFromString("1993960120797584048319"),
				Leverage:                 sdk.NewDec(2),
				MtpHealth:                sdk.MustNewDecFromStr("1.987938601732246814"),
				Position:                 types.Position_LONG,
			}

			openMTP, _ := marginKeeper.GetMTP(ctx, signer.String(), 1)

			fmt.Println(openExpectedMTP)
			fmt.Println(openMTP)

			require.Equal(t, openExpectedMTP, openMTP)

			openExpectedPool := clptypes.Pool{
				ExternalAsset:                 &externalAsset,
				NativeAssetBalance:            sdk.NewUintFromString("100001000000000000000000000"),
				ExternalAssetBalance:          sdk.NewUintFromString("99998006039879202415951681"),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.NewUintFromString("1993960120797584048319"),
				NativeLiabilities:             sdk.NewUintFromString("1000000000000000000000"),
				ExternalLiabilities:           sdk.ZeroUint(),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
				PoolUnits:                     sdk.ZeroUint(),
				Health:                        sdk.MustNewDecFromStr("0.999990000199996000"),
				InterestRate:                  sdk.MustNewDecFromStr("0.100000000000000000"),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			}

			openPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)

			fmt.Println(openExpectedPool)
			fmt.Println(openPool)

			require.Equal(t, openExpectedPool, openPool)

			_, closeError := msgServer.Close(sdk.WrapSDKContext(ctx), &msgClose)
			require.Nil(t, closeError)

			fmt.Println("native balance", app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			fmt.Println("external balance", app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("1987978360504281458998"))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			closeExpectedPool := clptypes.Pool{
				ExternalAsset:                 &externalAsset,
				NativeAssetBalance:            sdk.NewUintFromString("100000012021639495718541002"),
				ExternalAssetBalance:          sdk.NewUintFromString("100000000000000000000000000"),
				NativeCustody:                 sdk.ZeroUint(),
				ExternalCustody:               sdk.ZeroUint(),
				NativeLiabilities:             sdk.ZeroUint(),
				ExternalLiabilities:           sdk.ZeroUint(),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
				PoolUnits:                     sdk.ZeroUint(),
				Health:                        sdk.MustNewDecFromStr("0.999990000199996000"),
				InterestRate:                  sdk.MustNewDecFromStr("0.100000000000000000"),
				SwapPriceNative:               &SwapPriceNative,
				SwapPriceExternal:             &SwapPriceExternal,
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				RewardAmountExternal:          sdk.ZeroUint(),
			}

			closePool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)

			fmt.Println(closeExpectedPool)
			fmt.Println(closePool)

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
				LeverageMax:                              sdk.NewDec(2),
				HealthGainFactor:                         sdk.NewDec(2),
				InterestRateMin:                          sdk.NewDecWithPrec(5, 3),
				InterestRateMax:                          sdk.NewDec(3),
				InterestRateDecrease:                     sdk.NewDecWithPrec(1, 5),
				InterestRateIncrease:                     sdk.NewDecWithPrec(1, 5),
				RemovalQueueThreshold:                    sdk.ZeroDec(),
				ForceCloseFundPercentage:                 sdk.NewDecWithPrec(1, 1),
				ForceCloseFundAddress:                    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				IncrementalInterestPaymentFundPercentage: sdk.NewDecWithPrec(1, 1),
				IncrementalInterestPaymentFundAddress:    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
				IncrementalInterestPaymentEnabled:        false,
				PoolOpenThreshold:                        sdk.NewDecWithPrec(1, 1),
				MaxOpenPositions:                         10000,
				SqModifier:                               sdk.MustNewDecFromStr("10000000000000000000000000"),
				SafetyFactor:                             sdk.MustNewDecFromStr("1.05"),
				EpochLength:                              1,
				RowanCollateralEnabled:                   true,
				Pools: []string{
					externalAsset,
				},
			},
		}
		bz, _ = app.AppCodec().MarshalJSON(gs3)
		genesisState["margin"] = bz

		balances := []banktypes.Balance{
			{
				Address: "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85",
				Coins: sdk.Coins{
					sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("1000000000000000000000000000000"))),
					sdk.NewCoin(externalAsset, sdk.Int(sdk.NewUintFromString("1000000000000000000000000000000"))),
				},
			},
			{
				Address: signer,
				Coins: sdk.Coins{
					sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("100000000000000000000000000"))),
					sdk.NewCoin(externalAsset, sdk.Int(sdk.NewUintFromString("100000000000000000000000000"))),
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

	app.ClpKeeper.SetSwapFeeParams(ctx, clptypes.GetDefaultSwapFeeParams())

	msgOpen := types.MsgOpen{
		Signer:           signer,
		CollateralAsset:  externalAsset,
		CollateralAmount: sdk.NewUintFromString("10000"),
		BorrowAsset:      nativeAsset,
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
		Address:                  signer,
		CollateralAsset:          nativeAsset,
		CollateralAmount:         sdk.NewUintFromString("10000"),
		Liabilities:              sdk.NewUintFromString("10000"),
		InterestUnpaidCollateral: sdk.NewUintFromString("0"),
		CustodyAsset:             externalAsset,
		CustodyAmount:            sdk.NewUintFromString("20000"),
		Leverage:                 sdk.NewDec(1),
		MtpHealth:                sdk.MustNewDecFromStr("0.166666666666666667"),
		Position:                 types.Position_LONG,
		Id:                       1,
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
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("100000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776760883"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223239117"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("16616666666666666666"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148063"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628697475"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534732253"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("51158456267039810375"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422559"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628697475"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128776080"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("55427768026625260782"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
						},
					},
				},
				{
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("10000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776755604"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223244396"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("10000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("1661666666666666666"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148063"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628692637"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534735246"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("10000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("5115845626703981037"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422559"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628692637"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128778014"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("10000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("5542776802662526078"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
						},
					},
				},
				{
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("1000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776702820"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223297180"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("166166666666666666"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148059"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628644251"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534751280"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("511584562670398103"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422558"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628644251"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128714848"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("554277680266252607"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
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
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("100000000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776761468"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223238532"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("16616666666666666666666"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148063"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628698012"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534732101"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("51158456267039810375815"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422559"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628698012"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128776509"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("55427768026625260782599"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
						},
					},
				},
				{
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("200000000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776761469"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223238531"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("200000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("33233333333333333333333"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148063"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628698013"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534732100"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("200000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("102316912534079620751631"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422559"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628698013"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128776508"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("200000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("110855536053250521565198"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
						},
					},
				},
				{
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("1000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776702820"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223297180"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("166166666666666666"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148059"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628644251"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534751280"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("511584562670398103"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422558"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628644251"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128714848"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("554277680266252607"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
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
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("5000000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776761463"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223238537"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("5000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("830833333333333333333"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148063"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628698008"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534732096"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("5000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("2557922813351990518790"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422559"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628698008"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128776504"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("5000000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("2771388401331263039130"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
						},
					},
				},
				{
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("100000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776760883"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223239117"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("16616666666666666666"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148063"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628697475"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534732253"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("51158456267039810375"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422559"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628697475"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128776080"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("100000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("55427768026625260782"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
						},
					},
				},
				{
					X_A: sdk.NewUintFromString("100000000000000000000000"),
					Y_A: sdk.NewUintFromString("1000000000000000000"),
					chunks: []Chunk{
						{
							chunk:                                sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999997047460340145776702820"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("102952539659854223297180"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.916666666666666667"),
							mtpCustodyAmount:                     sdk.NewUintFromString("166166666666666666"),
							mtpHealth:                            sdk.MustNewDecFromStr("1.420621695012148059"),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999940626838645133628644251"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("158150007529031534751280"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.743438032763090219"),
							mtpCustodyAmount:                     sdk.NewUintFromString("511584562670398103"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.640333598224422558"),
							openErrorString:                      errors.New("110000000000000000000000: borrowed amount is higher than pool depth"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(100000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUintFromString("99999999841626838645133628644251"),
							signerExternalAssetBalanceAfterClose: sdk.NewUintFromString("100000000000000000000000000000000"),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(90000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(1000000000000000000, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUintFromString("249578272844105128714848"),
							poolExternalAssetBalanceAfterClose:   sdk.NewUintFromString("1000000000000000000"),
							poolHealthAfterClose:                 sdk.MustNewDecFromStr("0.722027241591648642"),
							mtpCustodyAmount:                     sdk.NewUintFromString("554277680266252607"),
							mtpHealth:                            sdk.MustNewDecFromStr("0.567973352996916171"),
							openErrorString:                      errors.New("198000000000000000000000: borrowed amount is higher than pool depth"),
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
						MaxOpenPositions:                         10000,
						LeverageMax:                              sdk.NewDec(2),
						HealthGainFactor:                         sdk.NewDec(1),
						InterestRateMin:                          sdk.NewDecWithPrec(5, 3),
						InterestRateMax:                          sdk.NewDec(3),
						InterestRateDecrease:                     sdk.NewDecWithPrec(1, 2),
						InterestRateIncrease:                     sdk.NewDecWithPrec(1, 2),
						RemovalQueueThreshold:                    sdk.ZeroDec(),
						ForceCloseFundPercentage:                 sdk.NewDecWithPrec(1, 1),
						ForceCloseFundAddress:                    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
						IncrementalInterestPaymentFundPercentage: sdk.NewDecWithPrec(1, 1),
						IncrementalInterestPaymentFundAddress:    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
						IncrementalInterestPaymentEnabled:        false,
						PoolOpenThreshold:                        sdk.NewDecWithPrec(1, 1),
						SqModifier:                               sdk.MustNewDecFromStr("10000000000000000000000000"),
						SafetyFactor:                             sdk.MustNewDecFromStr("1.05"),
						EpochLength:                              1,
						RowanCollateralEnabled:                   true,
						Pools: []string{
							ec.externalAsset,
						},
					},
				}
				bz, _ = app.AppCodec().MarshalJSON(gs3)
				genesisState["margin"] = bz

				nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUintFromString("100000000000000000000000000000000")))
				externalCoin := sdk.NewCoin(asset.Symbol, sdk.Int(sdk.NewUintFromString("100000000000000000000000000000000")))

				balances := []banktypes.Balance{
					{
						Address: "sif1pjm228rsgwqf23arkx7lm9ypkyma7mzr3y2n85",
						Coins: sdk.Coins{
							sdk.NewCoin(nativeAsset, sdk.Int(testItem.X_A)),
							sdk.NewCoin(asset.Symbol, sdk.Int(testItem.Y_A)),
						},
					},
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
							UnsettledExternalLiabilities:  sdk.ZeroUint(),
							UnsettledNativeLiabilities:    sdk.ZeroUint(),
							BlockInterestExternal:         sdk.ZeroUint(),
							BlockInterestNative:           sdk.ZeroUint(),
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

			app.ClpKeeper.SetSwapFeeParams(ctx, clptypes.GetDefaultSwapFeeParams())

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
						Leverage:         sdk.NewDec(2),
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
						Address:                  signer,
						CollateralAsset:          nativeAsset,
						CollateralAmount:         msgOpen.CollateralAmount,
						Liabilities:              msgOpen.CollateralAmount,
						InterestPaidCollateral:   sdk.ZeroUint(),
						InterestPaidCustody:      sdk.ZeroUint(),
						InterestUnpaidCollateral: sdk.ZeroUint(),
						CustodyAsset:             ec.externalAsset,
						CustodyAmount:            chunkItem.mtpCustodyAmount,
						Leverage:                 sdk.NewDec(2),
						MtpHealth:                chunkItem.mtpHealth,
						Position:                 types.Position_LONG,
						Id:                       uint64(i + 1),
					}
					openMTP, err := marginKeeper.GetMTP(ctx, signer, uint64(i+1))
					require.NoError(t, err)
					require.NotNil(t, openMTP)
					require.NotNil(t, openExpectedMTP)
					require.Equal(t, openExpectedMTP, openMTP)

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
						UnsettledExternalLiabilities:  sdk.ZeroUint(),
						UnsettledNativeLiabilities:    sdk.ZeroUint(),
						BlockInterestExternal:         sdk.ZeroUint(),
						BlockInterestNative:           sdk.ZeroUint(),
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
						UnsettledExternalLiabilities:  sdk.ZeroUint(),
						UnsettledNativeLiabilities:    sdk.ZeroUint(),
						BlockInterestExternal:         sdk.ZeroUint(),
						BlockInterestNative:           sdk.ZeroUint(),
						Health:                        chunkItem.poolHealthAfterClose,
						InterestRate:                  sdk.NewDecWithPrec(1, 2),
						SwapPriceNative:               &SwapPriceNative,
						SwapPriceExternal:             &SwapPriceExternal,
						RewardPeriodNativeDistributed: sdk.ZeroUint(),
						RewardAmountExternal:          sdk.ZeroUint(),
					}
					closePool, _ := marginKeeper.ClpKeeper().GetPool(ctx, ec.externalAsset)
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
		CollateralAsset:  externalAsset.Symbol,
		CollateralAmount: sdk.NewUintFromString("1000000000000000000000"),
		BorrowAsset:      nativeAsset,
		Position:         types.Position_LONG,
		Leverage:         sdk.NewDec(2),
	}

	marginKeeper.WhitelistAddress(ctx, msg1.Signer)

	_, openError := msgServer.Open(sdk.WrapSDKContext(ctx), &msg1)
	require.NoError(t, openError)

	openExpectedMTP := types.MTP{
		Address:                  signer.String(),
		CollateralAsset:          externalAsset.Symbol,
		CollateralAmount:         sdk.NewUintFromString("1000000000000000000000"),
		Liabilities:              sdk.NewUintFromString("1000000000000000000000"),
		InterestPaidCollateral:   sdk.ZeroUint(),
		InterestPaidCustody:      sdk.ZeroUint(),
		InterestUnpaidCollateral: sdk.ZeroUint(),
		CustodyAsset:             nativeAsset,
		CustodyAmount:            sdk.NewUintFromString("1993601279744051189762"),
		Leverage:                 sdk.NewDec(2),
		MtpHealth:                sdk.MustNewDecFromStr("1.987224302613536154"),
		Position:                 types.Position_LONG,
		Id:                       1,
	}
	openMTP, _ := marginKeeper.GetMTP(ctx, signer.String(), 1)
	fmt.Println(openExpectedMTP)
	fmt.Println(openMTP)
	require.Equal(t, openExpectedMTP, openMTP)

	msg2 := types.MsgOpen{
		Signer:           signer.String(),
		CollateralAsset:  externalAsset.Symbol,
		CollateralAmount: sdk.NewUintFromString("500000000000000000000"),
		BorrowAsset:      nativeAsset,
		Position:         types.Position_LONG,
		Leverage:         sdk.NewDec(2),
	}

	marginKeeper.WhitelistAddress(ctx, msg2.Signer)

	_, openError = msgServer.Open(sdk.WrapSDKContext(ctx), &msg2)
	require.NoError(t, openError)

	openExpectedMTP = types.MTP{
		Address:                  signer.String(),
		CollateralAsset:          externalAsset.Symbol,
		CollateralAmount:         sdk.NewUintFromString("1000000000000000000000"),
		Liabilities:              sdk.NewUintFromString("1000000000000000000000"),
		InterestPaidCollateral:   sdk.ZeroUint(),
		InterestPaidCustody:      sdk.ZeroUint(),
		InterestUnpaidCollateral: sdk.ZeroUint(),
		CustodyAsset:             nativeAsset,
		CustodyAmount:            sdk.NewUintFromString("1993601279744051189762"),
		Leverage:                 sdk.NewDec(2),
		MtpHealth:                sdk.MustNewDecFromStr("1.987224302613536154"),
		Position:                 types.Position_LONG,
		Id:                       1,
	}
	openMTP, _ = marginKeeper.GetMTP(ctx, signer.String(), 1)
	fmt.Println(openExpectedMTP)
	fmt.Println(openMTP)
	require.Equal(t, openExpectedMTP, openMTP)

	openExpectedMTP = types.MTP{
		Address:                  signer.String(),
		CollateralAsset:          externalAsset.Symbol,
		CollateralAmount:         sdk.NewUintFromString("500000000000000000000"),
		Liabilities:              sdk.NewUintFromString("500000000000000000000"),
		InterestPaidCollateral:   sdk.ZeroUint(),
		InterestPaidCustody:      sdk.ZeroUint(),
		InterestUnpaidCollateral: sdk.ZeroUint(),
		CustodyAsset:             nativeAsset,
		CustodyAmount:            sdk.NewUintFromString("996502287266229649201"),
		Leverage:                 sdk.NewDec(2),
		MtpHealth:                sdk.MustNewDecFromStr("1.987621151425775118"),
		Position:                 types.Position_LONG,
		Id:                       2,
	}
	openMTP, _ = marginKeeper.GetMTP(ctx, signer.String(), 2)
	fmt.Println(openExpectedMTP)
	fmt.Println(openMTP)
	require.Equal(t, openExpectedMTP, openMTP)
}
