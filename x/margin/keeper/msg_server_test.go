package keeper_test

import (
	"errors"
	"fmt"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	clptest "github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/keeper"
	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
			name: "margin enabled but denom does not exist",
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
			err:         tokenregistrytypes.ErrNotFound,
		},
		{
			name: "wrong address",
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
			errString:   errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name: "insufficient funds",
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
			errString:   errors.New("user does not have enough balance of the required coin"),
		},
		{
			name: "account funded",
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
			err:           nil,
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

/*
func TestKeeper_Close(t *testing.T) {
	table := []struct {
		name           string
		msgClose       types.MsgClose
		poolAsset      string
		token          string
		poolEnabled    bool
		fundedAccount  bool
		overrideSigner string
		err            error
		errString      error
	}{
		{
			name: "mtp does not exist",
			msgClose: types.MsgClose{
				Signer:          "xxx",
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
				Signer:          "xxx",
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
				Signer:          "xxx",
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "denom does not exist",
			msgClose: types.MsgClose{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:   "xxx",
			token:       "somethingelse",
			poolEnabled: true,
			err:         tokenregistrytypes.ErrNotFound,
		},
		{
			name: "wrong address",
			msgClose: types.MsgClose{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:   "xxx",
			token:       "xxx",
			poolEnabled: true,
			errString:   errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name: "insufficient funds",
			msgClose: types.MsgClose{
				Signer:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
				Position:        types.Position_LONG,
			},
			poolAsset:   "xxx",
			token:       "xxx",
			poolEnabled: true,
			errString:   errors.New("0xxx is smaller than 1000xxx: insufficient funds"),
		},
		{
			name: "account funded",
			msgClose: types.MsgClose{
				Signer:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
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

			marginKeeper.ClpKeeper().SetPool(ctx, &pool)

			var address string

			if tt.fundedAccount {
				nativeAsset := tt.msgClose.CollateralAsset
				externalAsset := clptypes.Asset{Symbol: tt.msgClose.BorrowAsset}

				nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(1000000000000)))
				externalCoin := sdk.NewCoin(externalAsset.Symbol, sdk.Int(sdk.NewUint(1000000000000)))
				err := app.BankKeeper.MintCoins(ctx, clptypes.ModuleName, sdk.NewCoins(nativeCoin, externalCoin))
				require.Nil(t, err)

				nativeCoin = sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(10000)))
				externalCoin = sdk.NewCoin(externalAsset.Symbol, sdk.Int(sdk.NewUint(10000)))

				_signer := clptest.GenerateAddress(clptest.AddressKey1)
				address = _signer.String()
				err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, _signer, sdk.NewCoins(nativeCoin))
				require.Nil(t, err)
				marginKeeper.BankKeeper().SendCoinsFromAccountToModule(ctx, _signer, types.ModuleName, sdk.NewCoins(nativeCoin))
			} else {
				address = tt.msgClose.Signer
			}

			msg := tt.msgClose
			msg.Signer = address

			var signer string = msg.Signer
			if tt.overrideSigner != "" {
				signer = tt.overrideSigner
			}

			addMTPKey(t, ctx, app, marginKeeper, msg.CollateralAsset, msg.BorrowAsset, signer, types.Position_LONG)

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
}*/

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
			err:           nil,
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
				LeverageMax:          sdk.NewUint(1),
				InterestRateMax:      sdk.NewDec(1),
				InterestRateMin:      sdk.ZeroDec(),
				InterestRateIncrease: sdk.NewDecWithPrec(1, 1),
				InterestRateDecrease: sdk.NewDecWithPrec(1, 1),
				HealthGainFactor:     sdk.NewDecWithPrec(1, 2),
				EpochLength:          0,
				ForceCloseThreshold:  sdk.ZeroDec(),
			}
			expectedGenesis := types.GenesisState{Params: &params}
			marginKeeper.InitGenesis(ctx, expectedGenesis)
			genesis := marginKeeper.ExportGenesis(ctx)
			require.Equal(t, expectedGenesis, *genesis)

			msgServer := keeper.NewMsgServerImpl(marginKeeper)

			nativeAsset := clptypes.NativeSymbol
			externalAsset := clptypes.Asset{Symbol: tt.externalAsset}

			pool := clptypes.Pool{
				ExternalAsset:        &externalAsset,
				NativeAssetBalance:   sdk.NewUint(1000000000000),
				ExternalAssetBalance: sdk.NewUint(1000000000000),
				NativeCustody:        sdk.ZeroUint(),
				ExternalCustody:      sdk.ZeroUint(),
				NativeLiabilities:    sdk.ZeroUint(),
				ExternalLiabilities:  sdk.ZeroUint(),
				PoolUnits:            sdk.ZeroUint(),
				Health:               sdk.ZeroDec(),
				InterestRate:         sdk.NewDecWithPrec(1, 1),
			}

			marginKeeper.SetEnabledPools(ctx, []string{tt.externalAsset})
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
			}
			msgClose := types.MsgClose{
				Signer: signer.String(),
				Id:     1,
			}
			fmt.Println(pool)
			_, openError := msgServer.Open(sdk.WrapSDKContext(ctx), &msgOpen)
			require.Nil(t, openError)

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(99999999999000))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			openExpectedMTP := types.MTP{
				Address:          signer.String(),
				CollateralAsset:  nativeAsset,
				CollateralAmount: sdk.NewUint(1000),
				LiabilitiesP:     sdk.NewUint(1000),
				LiabilitiesI:     sdk.ZeroUint(),
				CustodyAsset:     tt.externalAsset,
				CustodyAmount:    sdk.NewUint(4000),
				Leverage:         sdk.NewUint(1),
				MtpHealth:        sdk.NewDecWithPrec(1, 1),
				Position:         types.Position_LONG,
			}

			openMTP, _ := marginKeeper.GetMTP(ctx, signer.String(), 1)

			fmt.Println(openMTP)

			require.Equal(t, openExpectedMTP, openMTP)

			openExpectedPool := clptypes.Pool{
				ExternalAsset:        &externalAsset,
				NativeAssetBalance:   sdk.NewUint(1000000001000),
				ExternalAssetBalance: sdk.NewUint(999999996000),
				NativeCustody:        sdk.ZeroUint(),
				ExternalCustody:      sdk.NewUint(4000),
				NativeLiabilities:    sdk.NewUint(1000),
				ExternalLiabilities:  sdk.ZeroUint(),
				PoolUnits:            sdk.ZeroUint(),
				Health:               sdk.NewDecWithPrec(999999999000000002, 18),
				InterestRate:         sdk.NewDecWithPrec(1, 1),
			}

			openPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)
			fmt.Println(openPool)
			require.Equal(t, openExpectedPool, openPool)

			_, closeError := msgServer.Close(sdk.WrapSDKContext(ctx), &msgClose)
			require.Nil(t, closeError)

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(100000000006800))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			closeExpectedPool := clptypes.Pool{
				ExternalAsset:        &externalAsset,
				NativeAssetBalance:   sdk.NewUint(999999993200),
				ExternalAssetBalance: sdk.NewUint(1000000000000),
				NativeCustody:        sdk.ZeroUint(),
				ExternalCustody:      sdk.ZeroUint(),
				NativeLiabilities:    sdk.ZeroUint(),
				ExternalLiabilities:  sdk.ZeroUint(),
				PoolUnits:            sdk.ZeroUint(),
				Health:               sdk.NewDecWithPrec(999999999000000002, 18),
				InterestRate:         sdk.NewDecWithPrec(1, 1),
			}

			closePool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)
			require.Equal(t, closeExpectedPool, closePool)
		})
	}
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
		X_A    sdk.Uint
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(31),
							mtpHealth:                            sdk.NewDecWithPrec(162001036806635562, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(116926),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(13),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(680094924560566755, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(680094924560566755, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							closeErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000036994),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(7),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(63006),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(10),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(3),
							mtpHealth:                            sdk.NewDecWithPrec(164397974616952719, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(682091950568188387, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(682091950568188387, 18),
							mtpCustodyAmount:                     sdk.NewUint(9),
							mtpHealth:                            sdk.NewDecWithPrec(353577237340327734, 18),
							closeErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(108000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.ZeroUint(),
							mtpHealth:                            sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(682091950568188387, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(682091950568188387, 18),
							mtpCustodyAmount:                     sdk.NewUint(9),
							mtpHealth:                            sdk.NewDecWithPrec(353577237340327734, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000037598),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69444),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(62402),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(30556),
							mtpHealth:                            sdk.NewDecWithPrec(163049681237873180, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999982598),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999982598),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(117402),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(13101),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(680978178907437269, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(680978178907437269, 18),
							mtpCustodyAmount:                     sdk.NewUint(86899),
							mtpHealth:                            sdk.NewDecWithPrec(355899519859193208, 18),
							closeErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000037597),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(138889),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(62403),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(200000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(61111),
							mtpHealth:                            sdk.NewDecWithPrec(163049681237873180, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999982597),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(117403),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(26203),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(680980029349837300, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(680980029349837300, 18),
							mtpCustodyAmount:                     sdk.NewUint(173797),
							mtpHealth:                            sdk.NewDecWithPrec(355899519859193208, 18),
							closeErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(108000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.ZeroUint(),
							mtpHealth:                            sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(682091950568188387, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(682091950568188387, 18),
							mtpCustodyAmount:                     sdk.NewUint(9),
							mtpHealth:                            sdk.NewDecWithPrec(353577237340327734, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000037602),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(3472),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(62398),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(5000),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(1528),
							mtpHealth:                            sdk.NewDecWithPrec(163039047851960545, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999982602),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(117398),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(655),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(680970776923166162, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(680970776923166162, 18),
							mtpCustodyAmount:                     sdk.NewUint(4345),
							mtpHealth:                            sdk.NewDecWithPrec(355906428964312292, 18),
							closeErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000003807),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(11000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(6193),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(31),
							mtpHealth:                            sdk.NewDecWithPrec(161995788109509153, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999998307),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(11693),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(13),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(680102367242482406, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(680102367242482406, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356633380884450785, 18),
							closeErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(108000),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(1),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.ZeroUint(),
							mtpHealth:                            sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(1),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(682091950568188387, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(682091950568188387, 18),
							mtpCustodyAmount:                     sdk.NewUint(9),
							mtpHealth:                            sdk.NewDecWithPrec(353577237340327734, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpen:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpen:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterClose:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterClose: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpen:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpen:    sdk.NewUint(69),
							poolHealthAfterOpen:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterClose:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterClose:   sdk.NewUint(100),
							poolHealthAfterClose:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                     sdk.NewUint(87),
							mtpHealth:                            sdk.NewDecWithPrec(356640318512226279, 18),
							openErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
					},
				},
			},
		},
	}

	for _, ec := range table {
		ec := ec
		for _, testItem := range ec.tests {
			testItem := testItem

			ctx, app := test.CreateTestAppMargin(false)
			marginKeeper := app.MarginKeeper

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       ec.externalAsset,
				Decimals:    18,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			params := types.Params{
				LeverageMax:          sdk.NewUint(1),
				InterestRateMax:      sdk.NewDec(1),
				InterestRateMin:      sdk.ZeroDec(),
				InterestRateIncrease: sdk.NewDecWithPrec(1, 1),
				InterestRateDecrease: sdk.NewDecWithPrec(1, 1),
				HealthGainFactor:     sdk.NewDecWithPrec(1, 2),
				EpochLength:          0,
				ForceCloseThreshold:  sdk.ZeroDec(),
			}
			expectedGenesis := types.GenesisState{Params: &params}
			marginKeeper.InitGenesis(ctx, expectedGenesis)
			genesis := marginKeeper.ExportGenesis(ctx)
			require.Equal(t, expectedGenesis, *genesis)

			msgServer := keeper.NewMsgServerImpl(marginKeeper)

			nativeAsset := clptypes.NativeSymbol
			externalAsset := clptypes.Asset{Symbol: ec.externalAsset}

			pool := clptypes.Pool{
				ExternalAsset:        &externalAsset,
				NativeAssetBalance:   testItem.X_A,
				ExternalAssetBalance: testItem.Y_A,
				NativeCustody:        sdk.ZeroUint(),
				ExternalCustody:      sdk.ZeroUint(),
				NativeLiabilities:    sdk.ZeroUint(),
				ExternalLiabilities:  sdk.ZeroUint(),
				PoolUnits:            sdk.ZeroUint(),
				Health:               sdk.ZeroDec(),
				InterestRate:         sdk.NewDecWithPrec(1, 1),
			}

			marginKeeper.SetEnabledPools(ctx, []string{ec.externalAsset})
			marginKeeper.ClpKeeper().SetPool(ctx, &pool)

			nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(1000000000000)))
			externalCoin := sdk.NewCoin(ec.externalAsset, sdk.Int(sdk.NewUint(1000000000000)))
			err := app.BankKeeper.MintCoins(ctx, clptypes.ModuleName, sdk.NewCoins(nativeCoin, externalCoin))
			require.Nil(t, err)

			clpAccount := app.AccountKeeper.GetModuleAccount(ctx, clptypes.ModuleName)

			require.Equal(t, app.BankKeeper.GetBalance(ctx, clpAccount.GetAddress(), nativeAsset), nativeCoin)
			require.Equal(t, app.BankKeeper.GetBalance(ctx, clpAccount.GetAddress(), ec.externalAsset), externalCoin)

			signer := clptest.GenerateAddress(clptest.AddressKey1)
			nativeCoin = sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(100000000000000)))
			externalCoin = sdk.NewCoin(ec.externalAsset, sdk.Int(sdk.NewUint(1000000000000000)))
			err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(nativeCoin, externalCoin))
			require.Nil(t, err)

			require.Equal(t, app.BankKeeper.GetBalance(ctx, signer, nativeAsset), nativeCoin)
			require.Equal(t, app.BankKeeper.GetBalance(ctx, signer, ec.externalAsset), externalCoin)

			for _, chunkItem := range testItem.chunks {
				chunkItem := chunkItem
				name := fmt.Sprintf("%v, X_A=%v, Y_A=%v, delta x=%v%%", ec.name, testItem.X_A, testItem.Y_A, chunkItem.chunk)
				t.Run(name, func(t *testing.T) {
					msgOpen := types.MsgOpen{
						Signer:           signer.String(),
						CollateralAsset:  nativeAsset,
						CollateralAmount: testItem.X_A.Mul(chunkItem.chunk).Quo(sdk.NewUint(100)),
						BorrowAsset:      ec.externalAsset,
						Position:         types.Position_LONG,
					}
					msgClose := types.MsgClose{
						Signer: signer.String(),
						Id:     1,
					}
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

					require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(chunkItem.signerNativeAssetBalanceAfterOpen)), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
					require.Equal(t, sdk.NewCoin(ec.externalAsset, sdk.Int(chunkItem.signerExternalAssetBalanceAfterOpen)), app.BankKeeper.GetBalance(ctx, signer, ec.externalAsset))

					openExpectedMTP := types.MTP{
						Address:          signer.String(),
						CollateralAsset:  nativeAsset,
						CollateralAmount: msgOpen.CollateralAmount,
						LiabilitiesP:     msgOpen.CollateralAmount,
						LiabilitiesI:     sdk.ZeroUint(),
						CustodyAsset:     ec.externalAsset,
						CustodyAmount:    chunkItem.mtpCustodyAmount,
						Leverage:         sdk.NewUint(1),
						MtpHealth:        chunkItem.mtpHealth,
						Position:         types.Position_LONG,
					}
					openMTP, _ := marginKeeper.GetMTP(ctx, signer.String(), 1)
					require.Equal(t, openExpectedMTP, openMTP)

					openExpectedPool := clptypes.Pool{
						ExternalAsset:        &externalAsset,
						NativeAssetBalance:   chunkItem.poolNativeAssetBalanceAfterOpen,
						ExternalAssetBalance: chunkItem.poolExternalAssetBalanceAfterOpen,
						NativeCustody:        sdk.ZeroUint(),
						ExternalCustody:      chunkItem.mtpCustodyAmount,
						NativeLiabilities:    msgOpen.CollateralAmount,
						ExternalLiabilities:  sdk.ZeroUint(),
						PoolUnits:            sdk.ZeroUint(),
						Health:               chunkItem.poolHealthAfterOpen,
						InterestRate:         sdk.NewDecWithPrec(1, 1),
					}
					openPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, ec.externalAsset)
					require.Equal(t, openExpectedPool, openPool)

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

					require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(chunkItem.signerNativeAssetBalanceAfterClose)), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
					require.Equal(t, sdk.NewCoin(ec.externalAsset, sdk.Int(chunkItem.signerExternalAssetBalanceAfterClose)), app.BankKeeper.GetBalance(ctx, signer, ec.externalAsset))

					closeExpectedPool := clptypes.Pool{
						ExternalAsset:        &externalAsset,
						NativeAssetBalance:   chunkItem.poolNativeAssetBalanceAfterClose,
						ExternalAssetBalance: chunkItem.poolExternalAssetBalanceAfterClose,
						NativeCustody:        sdk.ZeroUint(),
						ExternalCustody:      sdk.ZeroUint(),
						NativeLiabilities:    sdk.ZeroUint(),
						ExternalLiabilities:  sdk.ZeroUint(),
						PoolUnits:            sdk.ZeroUint(),
						Health:               chunkItem.poolHealthAfterClose,
						InterestRate:         sdk.NewDecWithPrec(1, 1),
					}
					closePool, _ := marginKeeper.ClpKeeper().GetPool(ctx, ec.externalAsset)
					require.Equal(t, closeExpectedPool, closePool)
				})
			}
		}
	}
}
