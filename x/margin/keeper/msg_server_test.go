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

func TestKeeper_OpenLong(t *testing.T) {
	table := []struct {
		name          string
		msgOpenLong   types.MsgOpenLong
		poolAsset     string
		token         string
		poolEnabled   bool
		fundedAccount bool
		err           error
		errString     error
	}{
		{
			name: "pool does not exist",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "same collateral and native asset but pool does not exist",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "xxx",
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "same collateral and native asset but pool exists",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "xxx",
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "rowan",
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(types.ErrMTPDisabled, "rowan"),
		},
		{
			name: "margin enabled but denom does not exist",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
			},
			poolAsset:   "xxx",
			token:       "somethingelse",
			poolEnabled: true,
			err:         tokenregistrytypes.ErrNotFound,
		},
		{
			name: "wrong address",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "xxx",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
			},
			poolAsset:   "xxx",
			token:       "xxx",
			poolEnabled: true,
			errString:   errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name: "insufficient funds",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				CollateralAsset:  "xxx",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
			},
			poolAsset:   "xxx",
			token:       "xxx",
			poolEnabled: true,
			errString:   errors.New("user does not have enough balance of the required coin"),
		},
		{
			name: "account funded",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "xxx",
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
				nativeAsset := tt.msgOpenLong.CollateralAsset
				externalAsset := clptypes.Asset{Symbol: tt.msgOpenLong.BorrowAsset}

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
				address = tt.msgOpenLong.Signer
			}

			msg := tt.msgOpenLong
			msg.Signer = address

			_, got := msgServer.OpenLong(sdk.WrapSDKContext(ctx), &msg)

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

func TestKeeper_CloseLong(t *testing.T) {
	table := []struct {
		name           string
		msgCloseLong   types.MsgCloseLong
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
			msgCloseLong: types.MsgCloseLong{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
			},
			poolAsset:      "rowan",
			token:          "somethingelse",
			overrideSigner: "otheraddress",
			errString:      types.ErrMTPDoesNotExist,
		},
		{
			name: "pool does not exist",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "same collateral and native asset but pool does not exist",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "xxx",
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
			},
			poolAsset: "rowan",
			token:     "somethingelse",
			errString: sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name: "denom does not exist",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
			},
			poolAsset:   "xxx",
			token:       "somethingelse",
			poolEnabled: true,
			err:         tokenregistrytypes.ErrNotFound,
		},
		{
			name: "wrong address",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "xxx",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
			},
			poolAsset:   "xxx",
			token:       "xxx",
			poolEnabled: true,
			errString:   errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name: "insufficient funds",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				CollateralAsset: "xxx",
				BorrowAsset:     "xxx",
			},
			poolAsset:   "xxx",
			token:       "xxx",
			poolEnabled: true,
			errString:   errors.New("0xxx is smaller than 1000xxx: insufficient funds"),
		},
		{
			name: "account funded",
			msgCloseLong: types.MsgCloseLong{
				Signer:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				CollateralAsset: "rowan",
				BorrowAsset:     "xxx",
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
				nativeAsset := tt.msgCloseLong.CollateralAsset
				externalAsset := clptypes.Asset{Symbol: tt.msgCloseLong.BorrowAsset}

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
				address = tt.msgCloseLong.Signer
			}

			msg := tt.msgCloseLong
			msg.Signer = address

			var signer string = msg.Signer
			if tt.overrideSigner != "" {
				signer = tt.overrideSigner
			}

			addMTPKey(t, ctx, app, marginKeeper, msg.CollateralAsset, msg.BorrowAsset, signer)

			_, got := msgServer.CloseLong(sdk.WrapSDKContext(ctx), &msg)

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

func TestKeeper_OpenCloseLong(t *testing.T) {
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

			msgOpenLong := types.MsgOpenLong{
				Signer:           signer.String(),
				CollateralAsset:  nativeAsset,
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      tt.externalAsset,
			}
			msgCloseLong := types.MsgCloseLong{
				Signer:          signer.String(),
				CollateralAsset: nativeAsset,
				BorrowAsset:     tt.externalAsset,
			}
			fmt.Println(pool)
			_, openLongError := msgServer.OpenLong(sdk.WrapSDKContext(ctx), &msgOpenLong)
			require.Nil(t, openLongError)

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(99999999999000))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			openLongExpectedMTP := types.MTP{
				Address:          signer.String(),
				CollateralAsset:  nativeAsset,
				CollateralAmount: sdk.NewUint(1000),
				LiabilitiesP:     sdk.NewUint(1000),
				LiabilitiesI:     sdk.ZeroUint(),
				CustodyAsset:     tt.externalAsset,
				CustodyAmount:    sdk.NewUint(4000),
				Leverage:         sdk.NewUint(1),
				MtpHealth:        sdk.NewDecWithPrec(1, 1),
			}

			openLongMTP, _ := marginKeeper.GetMTP(ctx, nativeAsset, tt.externalAsset, signer.String())

			fmt.Println(openLongMTP)

			require.Equal(t, openLongExpectedMTP, openLongMTP)

			openLongExpectedPool := clptypes.Pool{
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

			openLongPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)
			fmt.Println(openLongPool)
			require.Equal(t, openLongExpectedPool, openLongPool)

			_, closeLongError := msgServer.CloseLong(sdk.WrapSDKContext(ctx), &msgCloseLong)
			require.Nil(t, closeLongError)

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(100000000006800))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			closeLongExpectedPool := clptypes.Pool{
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

			closeLongPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)
			require.Equal(t, closeLongExpectedPool, closeLongPool)
		})
	}
}

func TestKeeper_EC(t *testing.T) {
	type Chunk struct {
		chunk                                    sdk.Uint
		signerNativeAssetBalanceAfterOpenLong    sdk.Uint
		signerExternalAssetBalanceAfterOpenLong  sdk.Uint
		signerNativeAssetBalanceAfterCloseLong   sdk.Uint
		signerExternalAssetBalanceAfterCloseLong sdk.Uint
		poolNativeAssetBalanceAfterOpenLong      sdk.Uint
		poolExternalAssetBalanceAfterOpenLong    sdk.Uint
		poolHealthAfterOpenLong                  sdk.Dec
		poolNativeAssetBalanceAfterCloseLong     sdk.Uint
		poolExternalAssetBalanceAfterCloseLong   sdk.Uint
		poolHealthAfterCloseLong                 sdk.Dec
		mtpCustodyAmount                         sdk.Uint
		mtpHealth                                sdk.Dec
		openLongErrorString                      error
		openLongError                            error
		closeLongErrorString                     error
		closeLongError                           error
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
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(31),
							mtpHealth:                                sdk.NewDecWithPrec(162001036806635562, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(116926),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(13),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(680094924560566755, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(680094924560566755, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							closeLongErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(10),
					chunks: []Chunk{
						{
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000036994),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(7),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(63006),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(10),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(3),
							mtpHealth:                                sdk.NewDecWithPrec(164397974616952719, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(1),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(682091950568188387, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(682091950568188387, 18),
							mtpCustodyAmount:                         sdk.NewUint(9),
							mtpHealth:                                sdk.NewDecWithPrec(353577237340327734, 18),
							closeLongErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(1),
					chunks: []Chunk{
						{
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(1),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(108000),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(1),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.ZeroUint(),
							mtpHealth:                                sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(1),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(682091950568188387, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(682091950568188387, 18),
							mtpCustodyAmount:                         sdk.NewUint(9),
							mtpHealth:                                sdk.NewDecWithPrec(353577237340327734, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000037598),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69444),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(62402),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100000),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(30556),
							mtpHealth:                                sdk.NewDecWithPrec(163049681237873180, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999982598),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999982598),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(117402),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(13101),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(680978178907437269, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(680978178907437269, 18),
							mtpCustodyAmount:                         sdk.NewUint(86899),
							mtpHealth:                                sdk.NewDecWithPrec(355899519859193208, 18),
							closeLongErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(200000),
					chunks: []Chunk{
						{
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000037597),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(138889),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(62403),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(200000),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(61111),
							mtpHealth:                                sdk.NewDecWithPrec(163049681237873180, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999982597),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(117403),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(26203),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(680980029349837300, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(680980029349837300, 18),
							mtpCustodyAmount:                         sdk.NewUint(173797),
							mtpHealth:                                sdk.NewDecWithPrec(355899519859193208, 18),
							closeLongErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(1),
					chunks: []Chunk{
						{
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(1),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(108000),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(1),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.ZeroUint(),
							mtpHealth:                                sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(1),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(682091950568188387, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(682091950568188387, 18),
							mtpCustodyAmount:                         sdk.NewUint(9),
							mtpHealth:                                sdk.NewDecWithPrec(353577237340327734, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
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
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000037602),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(3472),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(62398),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(5000),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(1528),
							mtpHealth:                                sdk.NewDecWithPrec(163039047851960545, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999982602),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999935000),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(117398),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(655),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(680970776923166162, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(680970776923166162, 18),
							mtpCustodyAmount:                         sdk.NewUint(4345),
							mtpHealth:                                sdk.NewDecWithPrec(355906428964312292, 18),
							closeLongErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
					},
				},
				{
					X_A: sdk.NewUint(10000),
					Y_A: sdk.NewUint(100),
					chunks: []Chunk{
						{
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999999000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000003807),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(11000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(6193),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(31),
							mtpHealth:                                sdk.NewDecWithPrec(161995788109509153, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999998307),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(11693),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(13),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(680102367242482406, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(680102367242482406, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356633380884450785, 18),
							closeLongErrorString:                     errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
					},
				},
				{
					X_A: sdk.NewUint(100000),
					Y_A: sdk.NewUint(1),
					chunks: []Chunk{
						{
							chunk:                                    sdk.NewUint(10),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999990000),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(1),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(108000),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(1),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.ZeroUint(),
							mtpHealth:                                sdk.NewDecWithPrec(500000000000000000, 18),
						},
						{
							chunk:                                    sdk.NewUint(55),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999981994),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(118006),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(1),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(682091950568188387, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(116926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(682091950568188387, 18),
							mtpCustodyAmount:                         sdk.NewUint(9),
							mtpHealth:                                sdk.NewDecWithPrec(353577237340327734, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
						},
						{
							chunk:                                    sdk.NewUint(99),
							signerNativeAssetBalanceAfterOpenLong:    sdk.NewUint(99999999983074),
							signerExternalAssetBalanceAfterOpenLong:  sdk.NewUint(1000000000000000),
							signerNativeAssetBalanceAfterCloseLong:   sdk.NewUint(100000000038074),
							signerExternalAssetBalanceAfterCloseLong: sdk.NewUint(1000000000000000),
							poolNativeAssetBalanceAfterOpenLong:      sdk.NewUint(110000),
							poolExternalAssetBalanceAfterOpenLong:    sdk.NewUint(69),
							poolHealthAfterOpenLong:                  sdk.NewDecWithPrec(916666666666666667, 18),
							poolNativeAssetBalanceAfterCloseLong:     sdk.NewUint(61926),
							poolExternalAssetBalanceAfterCloseLong:   sdk.NewUint(100),
							poolHealthAfterCloseLong:                 sdk.NewDecWithPrec(916666666666666667, 18),
							mtpCustodyAmount:                         sdk.NewUint(87),
							mtpHealth:                                sdk.NewDecWithPrec(356640318512226279, 18),
							openLongErrorString:                      errors.New("not enough received asset tokens to swap"),
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
					msgOpenLong := types.MsgOpenLong{
						Signer:           signer.String(),
						CollateralAsset:  nativeAsset,
						CollateralAmount: testItem.X_A.Mul(chunkItem.chunk).Quo(sdk.NewUint(100)),
						BorrowAsset:      ec.externalAsset,
					}
					msgCloseLong := types.MsgCloseLong{
						Signer:          signer.String(),
						CollateralAsset: nativeAsset,
						BorrowAsset:     ec.externalAsset,
					}
					_, openLongError := msgServer.OpenLong(sdk.WrapSDKContext(ctx), &msgOpenLong)
					if chunkItem.openLongErrorString != nil {
						require.EqualError(t, openLongError, chunkItem.openLongErrorString.Error())
						return
					} else if chunkItem.openLongError != nil {
						require.ErrorIs(t, openLongError, chunkItem.openLongError)
						return
					} else {
						require.NoError(t, openLongError)
					}

					require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(chunkItem.signerNativeAssetBalanceAfterOpenLong)), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
					require.Equal(t, sdk.NewCoin(ec.externalAsset, sdk.Int(chunkItem.signerExternalAssetBalanceAfterOpenLong)), app.BankKeeper.GetBalance(ctx, signer, ec.externalAsset))

					openLongExpectedMTP := types.MTP{
						Address:          signer.String(),
						CollateralAsset:  nativeAsset,
						CollateralAmount: msgOpenLong.CollateralAmount,
						LiabilitiesP:     msgOpenLong.CollateralAmount,
						LiabilitiesI:     sdk.ZeroUint(),
						CustodyAsset:     ec.externalAsset,
						CustodyAmount:    chunkItem.mtpCustodyAmount,
						Leverage:         sdk.NewUint(1),
						MtpHealth:        chunkItem.mtpHealth,
					}
					openLongMTP, _ := marginKeeper.GetMTP(ctx, nativeAsset, ec.externalAsset, signer.String())
					require.Equal(t, openLongExpectedMTP, openLongMTP)

					openLongExpectedPool := clptypes.Pool{
						ExternalAsset:        &externalAsset,
						NativeAssetBalance:   chunkItem.poolNativeAssetBalanceAfterOpenLong,
						ExternalAssetBalance: chunkItem.poolExternalAssetBalanceAfterOpenLong,
						NativeCustody:        sdk.ZeroUint(),
						ExternalCustody:      chunkItem.mtpCustodyAmount,
						NativeLiabilities:    msgOpenLong.CollateralAmount,
						ExternalLiabilities:  sdk.ZeroUint(),
						PoolUnits:            sdk.ZeroUint(),
						Health:               chunkItem.poolHealthAfterOpenLong,
						InterestRate:         sdk.NewDecWithPrec(1, 1),
					}
					openLongPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, ec.externalAsset)
					require.Equal(t, openLongExpectedPool, openLongPool)

					_, closeLongError := msgServer.CloseLong(sdk.WrapSDKContext(ctx), &msgCloseLong)
					if chunkItem.closeLongErrorString != nil {
						require.EqualError(t, closeLongError, chunkItem.closeLongErrorString.Error())
						return
					} else if chunkItem.closeLongError != nil {
						require.ErrorIs(t, closeLongError, chunkItem.closeLongError)
						return
					} else {
						require.NoError(t, closeLongError)
					}

					require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(chunkItem.signerNativeAssetBalanceAfterCloseLong)), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
					require.Equal(t, sdk.NewCoin(ec.externalAsset, sdk.Int(chunkItem.signerExternalAssetBalanceAfterCloseLong)), app.BankKeeper.GetBalance(ctx, signer, ec.externalAsset))

					closeLongExpectedPool := clptypes.Pool{
						ExternalAsset:        &externalAsset,
						NativeAssetBalance:   chunkItem.poolNativeAssetBalanceAfterCloseLong,
						ExternalAssetBalance: chunkItem.poolExternalAssetBalanceAfterCloseLong,
						NativeCustody:        sdk.ZeroUint(),
						ExternalCustody:      sdk.ZeroUint(),
						NativeLiabilities:    sdk.ZeroUint(),
						ExternalLiabilities:  sdk.ZeroUint(),
						PoolUnits:            sdk.ZeroUint(),
						Health:               chunkItem.poolHealthAfterCloseLong,
						InterestRate:         sdk.NewDecWithPrec(1, 1),
					}
					closeLongPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, ec.externalAsset)
					require.Equal(t, closeLongExpectedPool, closeLongPool)
				})
			}
		}
	}
}
