package keeper_test

import (
	"errors"
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
	openLongTests := []struct {
		name             string
		signer           string
		collateralAsset  string
		collateralAmount sdk.Uint
		borrowAsset      string
		poolAsset        string
		token            string
		marginEnabled    bool
		fundedAccount    bool
		err              error
		errString        error
	}{
		{
			name:             "pool does not exist",
			signer:           "xxx",
			collateralAsset:  "xxx",
			collateralAmount: sdk.NewUint(1000),
			borrowAsset:      "xxx",
			poolAsset:        "rowan",
			token:            "somethingelse",
			errString:        sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name:             "same collateral and native asset but pool does not exist",
			signer:           "xxx",
			collateralAsset:  "rowan",
			collateralAmount: sdk.NewUint(1000),
			borrowAsset:      "xxx",
			poolAsset:        "rowan",
			token:            "somethingelse",
			errString:        sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name:             "same collateral and native asset but pool exists",
			signer:           "xxx",
			collateralAsset:  "rowan",
			collateralAmount: sdk.NewUint(1000),
			borrowAsset:      "rowan",
			poolAsset:        "rowan",
			token:            "somethingelse",
			errString:        sdkerrors.Wrap(types.ErrMTPDisabled, "rowan"),
		},
		{
			name:             "margin enabled but denom does not exist",
			signer:           "xxx",
			collateralAsset:  "xxx",
			collateralAmount: sdk.NewUint(1000),
			borrowAsset:      "xxx",
			poolAsset:        "xxx",
			token:            "somethingelse",
			marginEnabled:    true,
			err:              tokenregistrytypes.ErrNotFound,
		},
		{
			name:             "wrong address",
			signer:           "xxx",
			collateralAsset:  "xxx",
			collateralAmount: sdk.NewUint(1000),
			borrowAsset:      "xxx",
			poolAsset:        "xxx",
			token:            "xxx",
			marginEnabled:    true,
			errString:        errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name:             "insufficient funds",
			signer:           "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAsset:  "xxx",
			collateralAmount: sdk.NewUint(1000),
			borrowAsset:      "xxx",
			poolAsset:        "xxx",
			token:            "xxx",
			marginEnabled:    true,
			errString:        errors.New("user does not have enough balance of the required coin"),
		},
		{
			name:             "account funded",
			signer:           "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAsset:  "rowan",
			collateralAmount: sdk.NewUint(1000),
			borrowAsset:      "rowan",
			poolAsset:        "rowan",
			token:            "rowan",
			marginEnabled:    true,
			fundedAccount:    true,
			err:              nil,
		},
	}

	for _, tt := range openLongTests {
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

			if tt.marginEnabled {
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
				address = tt.signer
			}

			msg := types.MsgOpenLong{
				Signer:           address,
				CollateralAsset:  tt.collateralAsset,
				CollateralAmount: sdk.NewUint(1000),
				BorrowAsset:      tt.borrowAsset,
			}

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
	openLongTests := []struct {
		name            string
		signer          string
		collateralAsset string
		borrowAsset     string
		poolAsset       string
		token           string
		id              uint64
		marginEnabled   bool
		fundedAccount   bool
		overrideSigner  string
		err             error
		errString       error
	}{
		{
			name:            "mtp does not exist",
			signer:          "xxx",
			collateralAsset: "xxx",
			borrowAsset:     "xxx",
			poolAsset:       "rowan",
			id:              2,
			token:           "somethingelse",
			overrideSigner:  "otheraddress",
			errString:       types.ErrMTPDoesNotExist,
		},
		{
			name:            "pool does not exist",
			signer:          "xxx",
			collateralAsset: "xxx",
			borrowAsset:     "xxx",
			poolAsset:       "rowan",
			id:              1,
			token:           "somethingelse",
			errString:       sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name:            "same collateral and native asset but pool does not exist",
			signer:          "xxx",
			collateralAsset: "rowan",
			borrowAsset:     "xxx",
			poolAsset:       "rowan",
			id:              1,
			token:           "somethingelse",
			errString:       sdkerrors.Wrap(clptypes.ErrPoolDoesNotExist, "xxx"),
		},
		{
			name:            "denom does not exist",
			signer:          "xxx",
			collateralAsset: "xxx",
			borrowAsset:     "xxx",
			poolAsset:       "xxx",
			id:              1,
			token:           "somethingelse",
			marginEnabled:   true,
			err:             tokenregistrytypes.ErrNotFound,
		},
		{
			name:            "wrong address",
			signer:          "xxx",
			collateralAsset: "xxx",
			borrowAsset:     "xxx",
			poolAsset:       "xxx",
			token:           "xxx",
			id:              1,
			marginEnabled:   true,
			errString:       errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name:            "insufficient funds",
			signer:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAsset: "xxx",
			borrowAsset:     "xxx",
			poolAsset:       "xxx",
			id:              1,
			token:           "xxx",
			marginEnabled:   true,
			errString:       errors.New("0xxx is smaller than 1000xxx: insufficient funds"),
		},
		{
			name:            "account funded",
			signer:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAsset: "rowan",
			borrowAsset:     "rowan",
			poolAsset:       "rowan",
			token:           "rowan",
			id:              1,
			marginEnabled:   true,
			fundedAccount:   true,
			err:             nil,
		},
	}

	for _, tt := range openLongTests {
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
			if tt.marginEnabled {
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
				address = tt.signer
			}

			msg := types.MsgCloseLong{
				Signer: address,
				Id:     tt.id,
			}

			var signer string = msg.Signer
			if tt.overrideSigner != "" {
				signer = tt.overrideSigner
			}

			addMTPKey(t, ctx, app, marginKeeper, tt.collateralAsset, tt.borrowAsset, signer, 1)

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
