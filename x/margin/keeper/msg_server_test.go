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
				BorrowAsset:      "rowan",
			},
			poolAsset:     "rowan",
			token:         "rowan",
			poolEnabled:   true,
			fundedAccount: true,
			err:           nil,
		},
		{
			name: "account funded",
			msgOpenLong: types.MsgOpenLong{
				Signer:           "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
				CollateralAsset:  "rowan",
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      "rowan",
			},
			poolAsset:     "rowan",
			token:         "rowan",
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
				_signer := clptest.GenerateAddress(clptest.AddressKey1)
				address = _signer.String()
				nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(sdk.NewUintFromString("10000")))
				err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, _signer, sdk.NewCoins(nativeCoin))
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
				BorrowAsset:     "rowan",
			},
			poolAsset:     "rowan",
			token:         "rowan",
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
				_signer := clptest.GenerateAddress(clptest.AddressKey1)
				address = _signer.String()
				nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(sdk.NewUintFromString("10000")))
				err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, _signer, sdk.NewCoins(nativeCoin))
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
				InterestRateMax:      sdk.ZeroDec(),
				InterestRateMin:      sdk.ZeroDec(),
				InterestRateIncrease: sdk.ZeroDec(),
				InterestRateDecrease: sdk.ZeroDec(),
				HealthGainFactor:     sdk.ZeroDec(),
				EpochLength:          0,
			}
			expectedGenesis := types.GenesisState{Params: &params}

			genesis := marginKeeper.ExportGenesis(ctx)
			require.Equal(t, expectedGenesis, *genesis)

			msgServer := keeper.NewMsgServerImpl(marginKeeper)

			nativeAsset := clptypes.NativeSymbol
			externalAsset := clptypes.Asset{Symbol: tt.externalAsset}

			pool := clptypes.Pool{
				ExternalAsset:        &externalAsset,
				NativeAssetBalance:   sdk.NewUint(1000000000000),
				ExternalAssetBalance: sdk.NewUint(10000000),
				NativeCustody:        sdk.NewUint(0),
				ExternalCustody:      sdk.NewUint(0),
				NativeLiabilities:    sdk.NewUint(0),
				ExternalLiabilities:  sdk.NewUint(0),
				PoolUnits:            sdk.NewUint(0),
				Health:               sdk.NewDec(0),
				InterestRate:         sdk.NewDec(0),
			}

			marginKeeper.SetEnabledPools(ctx, []string{tt.externalAsset})
			marginKeeper.ClpKeeper().SetPool(ctx, &pool)

			signer := clptest.GenerateAddress(clptest.AddressKey1)
			nativeCoin := sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(100000000000000)))
			externalCoin := sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000)))
			err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(nativeCoin, externalCoin))
			require.Nil(t, err)

			nativeCoinOk := app.ClpKeeper.HasBalance(ctx, signer, nativeCoin)
			require.True(t, nativeCoinOk)
			externalCoinOk := app.ClpKeeper.HasBalance(ctx, signer, externalCoin)
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

			_, openLongError := msgServer.OpenLong(sdk.WrapSDKContext(ctx), &msgOpenLong)
			require.Nil(t, openLongError)

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(99999999999000))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			openLongExpectedMTP := types.MTP{
				Address:          signer.String(),
				CollateralAsset:  nativeAsset,
				CollateralAmount: sdk.NewUint(2000),
				LiabilitiesP:     sdk.NewUint(1000),
				LiabilitiesI:     sdk.ZeroUint(),
				CustodyAsset:     tt.externalAsset,
				CustodyAmount:    sdk.ZeroUint(),
				Leverage:         sdk.NewUint(1),
				MtpHealth:        sdk.ZeroDec(),
			}

			openLongMTP, _ := marginKeeper.GetMTP(ctx, nativeAsset, tt.externalAsset, signer.String())

			fmt.Print(openLongMTP)

			require.Equal(t, openLongExpectedMTP, openLongMTP)

			openLongExpectedPool := clptypes.Pool{
				ExternalAsset:        &externalAsset,
				NativeAssetBalance:   sdk.NewUint(1000000000000),
				ExternalAssetBalance: sdk.NewUint(10000000),
				NativeCustody:        sdk.NewUint(0),
				ExternalCustody:      sdk.NewUint(0),
				NativeLiabilities:    sdk.NewUint(0),
				ExternalLiabilities:  sdk.NewUint(0),
				PoolUnits:            sdk.NewUint(0),
				Health:               sdk.NewDec(1),
				InterestRate:         sdk.NewDec(0),
			}

			openLongPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)
			require.Equal(t, openLongExpectedPool, openLongPool)

			_, closeLongError := msgServer.CloseLong(sdk.WrapSDKContext(ctx), &msgCloseLong)
			// require.Nil(t, closeLongError)
			require.EqualError(t, closeLongError, "0rowan is smaller than 1000rowan: insufficient funds")

			require.Equal(t, sdk.NewCoin(nativeAsset, sdk.Int(sdk.NewUint(99999999999000))), app.BankKeeper.GetBalance(ctx, signer, nativeAsset))
			require.Equal(t, sdk.NewCoin(tt.externalAsset, sdk.Int(sdk.NewUint(1000000000000000))), app.BankKeeper.GetBalance(ctx, signer, tt.externalAsset))

			closeLongExpectedPool := clptypes.Pool{
				ExternalAsset:        &externalAsset,
				NativeAssetBalance:   sdk.NewUint(1000000000000),
				ExternalAssetBalance: sdk.NewUint(10000000),
				NativeCustody:        sdk.NewUint(0),
				ExternalCustody:      sdk.NewUint(0),
				NativeLiabilities:    sdk.NewUint(0),
				ExternalLiabilities:  sdk.NewUint(0),
				PoolUnits:            sdk.NewUint(0),
				Health:               sdk.NewDec(1),
				InterestRate:         sdk.NewDec(0),
			}

			closeLongPool, _ := marginKeeper.ClpKeeper().GetPool(ctx, tt.externalAsset)
			require.Equal(t, closeLongExpectedPool, closeLongPool)
		})
	}
}
