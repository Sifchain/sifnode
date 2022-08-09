//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	sifapp "github.com/Sifchain/sifnode/app"
	clptest "github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_Errors(t *testing.T) {
	_, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper
	require.NotNil(t, marginKeeper)
}

func TestKeeper_SetMTP(t *testing.T) {
	table := []struct {
		name      string
		mtp       types.MTP
		err       error
		errString error
	}{
		{
			name:      "missed defining asset",
			mtp:       types.MTP{},
			errString: errors.New("no asset specified: mtp invalid"),
		},
		{
			name:      "define asset but no address",
			mtp:       types.MTP{CollateralAsset: "xxx"},
			errString: errors.New("no address specified: mtp invalid"),
		},
		{
			name:      "define asset and address but no position",
			mtp:       types.MTP{CollateralAsset: "xxx", Address: "xxx"},
			errString: errors.New("no position specified: mtp invalid"),
		},
		{
			name: "define asset and address with long position",
			mtp:  types.MTP{CollateralAsset: "xxx", Address: "xxx", Position: types.Position_LONG},
		},
		{
			name: "define asset and address with short position",
			mtp:  types.MTP{CollateralAsset: "xxx", Address: "xxx", Position: types.Position_SHORT},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, marginKeeper := initKeeper(t)
			got := marginKeeper.SetMTP(ctx, &tt.mtp)

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

func TestKeeper_GetMTP(t *testing.T) {
	t.Run("get MTP from a store key that exists", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		want := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "xxx", types.Position_LONG, 1)
		got, err := marginKeeper.GetMTP(ctx, want.Address, 1)

		require.NoError(t, err)
		require.Equal(t, got, want)
	})

	t.Run("fails when store keys does not exist", func(t *testing.T) {
		ctx, _, marginKeeper := initKeeper(t)
		_, got := marginKeeper.GetMTP(ctx, "xxx", 0)

		require.ErrorIs(t, got, types.ErrMTPDoesNotExist)
	})
}

func TestKeeper_GetMTPIterator(t *testing.T) {
	ctx, app, marginKeeper := initKeeper(t)
	want := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "xxx", types.Position_LONG, 1)
	iterator := marginKeeper.GetMTPIterator(ctx)
	bytesValue := iterator.Value()
	var got types.MTP
	types.ModuleCdc.MustUnmarshal(bytesValue, &got)

	require.Equal(t, got, want)
}

func TestKeeper_GetMTPs(t *testing.T) {
	ctx, app, marginKeeper := initKeeper(t)
	key1 := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "key1", types.Position_LONG, 1)
	key2 := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "key2", types.Position_LONG, 1)
	want := []*types.MTP{&key1, &key2}
	got, _, err := marginKeeper.GetMTPs(ctx, nil)
	require.NoError(t, err)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestKeeper_DestroyMTP(t *testing.T) {
	t.Run("key does not exist", func(t *testing.T) {
		ctx, _, marginKeeper := initKeeper(t)
		got := marginKeeper.DestroyMTP(ctx, "xxx", 1)

		require.ErrorIs(t, got, types.ErrMTPDoesNotExist)
	})
	t.Run("key exists", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "xxx", types.Position_LONG, 1)
		got := marginKeeper.DestroyMTP(ctx, mtp.Address, 1)

		require.NoError(t, got)
	})
}

func TestKeeper_ClpKeeper(t *testing.T) {
	_, _, marginKeeper := initKeeper(t)
	clp := marginKeeper.ClpKeeper()
	require.NotNil(t, clp)
}

func TestKeeper_BankKeeper(t *testing.T) {
	_, _, marginKeeper := initKeeper(t)
	bank := marginKeeper.BankKeeper()
	require.NotNil(t, bank)
}

func TestKeeper_CLPSwap(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}
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

	custodySwapTests := []struct {
		name       string
		denom      string
		decimals   int64
		to         string
		sentAmount sdk.Uint
		err        error
	}{
		{
			name:       "denom not registered",
			denom:      "unregistred_denom",
			decimals:   18,
			to:         "xxx",
			sentAmount: sdk.NewUint(0),
			err:        tokenregistrytypes.ErrNotFound,
		},
		{
			name:       "invalid sent amount",
			denom:      "rowan",
			decimals:   18,
			to:         "xxx",
			sentAmount: sdk.NewUint(0),
			err:        nil,
		},
		{
			name:       "no token adjustment and non-rowan target asset",
			denom:      "rowan",
			decimals:   18,
			to:         "xxx",
			sentAmount: sdk.NewUint(10000000000000),
			err:        clptypes.ErrNotEnoughAssetTokens,
		},
		{
			name:       "no token adjustment and rowan target asset",
			denom:      "rowan",
			decimals:   18,
			to:         "rowan",
			sentAmount: sdk.NewUint(10000000000000),
			err:        clptypes.ErrNotEnoughAssetTokens,
		},
		{
			name:       "token adjustment and non-rowan target asset",
			denom:      "rowan",
			decimals:   9,
			to:         "xxx",
			sentAmount: sdk.NewUint(1000),
			err:        nil,
		},
		{
			name:       "token adjustment and rowan target asset",
			denom:      "rowan",
			decimals:   9,
			to:         "rowan",
			sentAmount: sdk.NewUint(1000000000000),
			err:        clptypes.ErrNotEnoughAssetTokens,
		},
	}

	for _, tt := range custodySwapTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			_, got := marginKeeper.CLPSwap(ctx, tt.sentAmount, tt.to, pool)

			if tt.err == nil {
				require.NoError(t, got)
			} else {
				require.ErrorIs(t, got, tt.err)
			}
		})
	}
}

func TestKeeper_Borrow(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}
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

	borrowTests := []struct {
		name             string
		denom            string
		decimals         int64
		to               string
		address          string
		collateralAmount sdk.Uint
		custodyAmount    sdk.Uint
		leverage         sdk.Dec
		health           sdk.Dec
		fundedAccount    bool
		err              error
		errString        error
	}{
		{
			name:             "wrong address",
			denom:            "unregistered_denom",
			decimals:         18,
			to:               "rowan",
			address:          "xxx",
			collateralAmount: sdk.NewUint(10000),
			custodyAmount:    sdk.NewUint(1000),
			leverage:         sdk.NewDec(1),
			health:           sdk.NewDec(1),
			errString:        errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name:             "not enough fund",
			denom:            "unregistered_denom",
			decimals:         18,
			to:               "rowan",
			address:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAmount: sdk.NewUint(10000),
			custodyAmount:    sdk.NewUint(1000),
			leverage:         sdk.NewDec(1),
			health:           sdk.NewDec(1),
			errString:        errors.New("user does not have enough balance of the required coin"),
		},
		{
			name:             "denom not registered",
			denom:            "unregistered_denom",
			decimals:         18,
			to:               "rowan",
			address:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAmount: sdk.NewUint(10000),
			custodyAmount:    sdk.NewUint(1000),
			leverage:         sdk.NewDec(1),
			health:           sdk.NewDec(1),
			fundedAccount:    true,
			err:              tokenregistrytypes.ErrNotFound,
		},
		{
			name:             "not enough balance required to swap",
			denom:            "rowan",
			decimals:         18,
			to:               "rowan",
			address:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAmount: sdk.NewUint(1000000000000000),
			custodyAmount:    sdk.NewUint(1000000000000000),
			leverage:         sdk.NewDec(1),
			health:           sdk.NewDec(1),
			errString:        errors.New("user does not have enough balance of the required coin"),
		},
		{
			name:             "invalid address",
			denom:            "rowan",
			decimals:         9,
			to:               "xxx",
			address:          "xxx",
			collateralAmount: sdk.NewUint(10000),
			custodyAmount:    sdk.NewUint(1000),
			leverage:         sdk.NewDec(1),
			health:           sdk.NewDec(1),
			errString:        errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name:             "insufficient funds",
			denom:            "rowan",
			decimals:         9,
			to:               "xxx",
			address:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAmount: sdk.NewUint(10000),
			custodyAmount:    sdk.NewUint(1000),
			leverage:         sdk.NewDec(1),
			health:           sdk.NewDec(1),
			errString:        errors.New("user does not have enough balance of the required coin"),
		},
		{
			name:             "account funded",
			denom:            "rowan",
			decimals:         9,
			to:               "rowan",
			address:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAmount: sdk.NewUint(1000),
			custodyAmount:    sdk.NewUint(1000),
			leverage:         sdk.NewDec(1),
			health:           sdk.NewDec(1),
			fundedAccount:    true,
			err:              nil,
		},
	}

	for _, tt := range borrowTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			var address string

			if tt.fundedAccount {
				signer := clptest.GenerateAddress(clptest.AddressKey1)
				address = signer.String()
				nativeCoin := sdk.NewCoin(clptypes.NativeSymbol, sdk.Int(sdk.NewUintFromString("10000")))
				err := sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, signer, sdk.NewCoins(nativeCoin))
				require.Nil(t, err)
			} else {
				address = tt.address
			}

			mtp := addMTPKey(t, ctx, app, marginKeeper, tt.to, "xxx", address, types.Position_LONG, 1)

			got := marginKeeper.Borrow(ctx, tt.to, tt.collateralAmount, tt.custodyAmount, &mtp, &pool, tt.leverage)

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

func TestKeeper_UpdatePoolHealth(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}
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

	ctx, _, marginKeeper := initKeeper(t)

	err := marginKeeper.UpdatePoolHealth(ctx, &pool)
	require.Nil(t, err)
}

func TestKeeper_UpdateMTPHealth(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}
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

	updateMTPHealthTests := []struct {
		name             string
		denom            string
		decimals         int64
		to               string
		collateralAmount sdk.Uint
		custodyAmount    sdk.Uint
		liabilitiesP     sdk.Uint
		liabilitiesI     sdk.Uint
		health           sdk.Dec
		err              error
		errString        error
	}{
		{
			name:             "denom not registered",
			denom:            "unregistred_denom",
			decimals:         18,
			to:               "xxx",
			collateralAmount: sdk.NewUint(1000),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(1000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			err:              tokenregistrytypes.ErrNotFound,
		},
		{
			name:             "not enough received asset tokens to swap",
			denom:            "rowan",
			decimals:         18,
			to:               "rowan",
			collateralAmount: sdk.NewUint(1000),
			custodyAmount:    sdk.NewUint(10000000000),
			liabilitiesP:     sdk.NewUint(1000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			err:              clptypes.ErrNotEnoughAssetTokens,
		},
		{
			name:             "swap with same asset",
			denom:            "rowan",
			decimals:         18,
			to:               "rowan",
			collateralAmount: sdk.NewUint(1000),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(1000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			err:              nil,
		},
		{
			name:             "swap with different asset",
			denom:            "rowan",
			decimals:         9,
			to:               "xxx",
			collateralAmount: sdk.NewUint(1000),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(1000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			err:              nil,
		},
		{
			name:             "insufficient liabilities funds",
			denom:            "rowan",
			decimals:         18,
			to:               "xxx",
			collateralAmount: sdk.NewUint(1000),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(10000000000000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			err:              clptypes.ErrNotEnoughAssetTokens,
		},
		{
			name:             "mtp invalid",
			denom:            "rowan",
			decimals:         18,
			to:               "xxx",
			collateralAmount: sdk.NewUint(0),
			custodyAmount:    sdk.NewUint(0),
			liabilitiesP:     sdk.NewUint(0),
			liabilitiesI:     sdk.NewUint(0),
			health:           sdk.NewDec(1),
			err:              types.ErrMTPInvalid,
		},
	}

	for _, tt := range updateMTPHealthTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			mtp := addMTPKey(t, ctx, app, marginKeeper, tt.to, "xxx", "xxx", types.Position_LONG, 1)
			mtp.CustodyAmount = tt.custodyAmount
			mtp.LiabilitiesP = tt.liabilitiesP
			mtp.CollateralAmount = tt.collateralAmount
			mtp.LiabilitiesI = tt.liabilitiesI

			_, got := marginKeeper.UpdateMTPHealth(ctx, mtp, pool)

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

func TestKeeper_TakeInCustody(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}

	t.Run("settlement asset and mtp asset is equal", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "rowan", "xxx", "xxx", types.Position_LONG, 1)

		pool := clptypes.Pool{
			ExternalAsset:        &asset,
			NativeAssetBalance:   sdk.NewUint(1000),
			NativeLiabilities:    sdk.NewUint(1000),
			ExternalCustody:      sdk.NewUint(1000),
			ExternalAssetBalance: sdk.NewUint(1000),
			ExternalLiabilities:  sdk.NewUint(1000),
			NativeCustody:        sdk.NewUint(1000),
			PoolUnits:            sdk.NewUint(1),
			Health:               sdk.NewDec(1),
		}

		got := marginKeeper.TakeInCustody(ctx, mtp, &pool)

		require.NoError(t, got)
	})

	t.Run("settlement asset and mtp asset is not equal", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "notrowan", "xxx", "xxx", types.Position_LONG, 1)

		pool := clptypes.Pool{
			ExternalAsset:        &asset,
			NativeAssetBalance:   sdk.NewUint(1000),
			NativeLiabilities:    sdk.NewUint(1000),
			ExternalCustody:      sdk.NewUint(1000),
			ExternalAssetBalance: sdk.NewUint(1000),
			ExternalLiabilities:  sdk.NewUint(1000),
			NativeCustody:        sdk.NewUint(1000),
			PoolUnits:            sdk.NewUint(1),
			Health:               sdk.NewDec(1),
		}

		got := marginKeeper.TakeInCustody(ctx, mtp, &pool)

		require.NoError(t, got)
	})
}

func TestKeeper_TakeOutCustody(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}

	t.Run("settlement asset and mtp asset is equal", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "rowan", "xxx", "xxx", types.Position_LONG, 1)

		pool := clptypes.Pool{
			ExternalAsset:        &asset,
			NativeAssetBalance:   sdk.NewUint(1000),
			NativeLiabilities:    sdk.NewUint(1000),
			ExternalCustody:      sdk.NewUint(1000),
			ExternalAssetBalance: sdk.NewUint(1000),
			ExternalLiabilities:  sdk.NewUint(1000),
			NativeCustody:        sdk.NewUint(1000),
			PoolUnits:            sdk.NewUint(1),
			Health:               sdk.NewDec(1),
		}

		got := marginKeeper.TakeOutCustody(ctx, mtp, &pool)

		require.NoError(t, got)
	})

	t.Run("settlement asset and mtp asset is not equal", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "notrowan", "xxx", "xxx", types.Position_LONG, 1)

		pool := clptypes.Pool{
			ExternalAsset:        &asset,
			NativeAssetBalance:   sdk.NewUint(1000),
			NativeLiabilities:    sdk.NewUint(1000),
			ExternalCustody:      sdk.NewUint(1000),
			ExternalAssetBalance: sdk.NewUint(1000),
			ExternalLiabilities:  sdk.NewUint(1000),
			NativeCustody:        sdk.NewUint(1000),
			PoolUnits:            sdk.NewUint(1),
			Health:               sdk.NewDec(1),
		}

		got := marginKeeper.TakeOutCustody(ctx, mtp, &pool)

		require.NoError(t, got)
	})
}

func TestKeeper_Repay(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}
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

	repayTests := []struct {
		name             string
		denom            string
		decimals         int64
		to               string
		address          string
		collateralAmount sdk.Uint
		custodyAmount    sdk.Uint
		liabilitiesP     sdk.Uint
		liabilitiesI     sdk.Uint
		health           sdk.Dec
		repayAmount      sdk.Uint
		overrideAddress  string
		err              error
		errString        error
	}{
		{
			name:             "denom not registered",
			denom:            "unregistred_denom",
			decimals:         18,
			to:               "xxx",
			address:          "xxx",
			collateralAmount: sdk.NewUint(1000),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(1000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			repayAmount:      sdk.NewUint(1),
			err:              tokenregistrytypes.ErrNotFound,
		},
		{
			name:             "cannot afford principle liability",
			denom:            "rowan",
			decimals:         18,
			to:               "xxx",
			address:          "xxx",
			collateralAmount: sdk.NewUint(0),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(1000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			repayAmount:      sdk.NewUint(0),
			err:              nil,
		},
		{
			name:             "v principle libarity; x excess liability",
			denom:            "rowan",
			decimals:         18,
			to:               "xxx",
			address:          "xxx",
			collateralAmount: sdk.NewUint(0),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(0),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			repayAmount:      sdk.NewUint(0),
			err:              nil,
		},
		{
			name:             "can afford both",
			denom:            "rowan",
			decimals:         18,
			to:               "xxx",
			address:          "xxx",
			collateralAmount: sdk.NewUint(0),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(0),
			liabilitiesI:     sdk.NewUint(0),
			health:           sdk.NewDec(1),
			repayAmount:      sdk.NewUint(0),
			err:              nil,
		},
		{
			name:             "non zero return amount + fails because of wrong address",
			denom:            "rowan",
			decimals:         18,
			to:               "xxx",
			address:          "xxx",
			collateralAmount: sdk.NewUint(0),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(0),
			liabilitiesI:     sdk.NewUint(0),
			health:           sdk.NewDec(1),
			repayAmount:      sdk.NewUint(1000),
			errString:        errors.New("decoding bech32 failed: invalid bech32 string length 3"),
		},
		{
			name:             "non zero return amount",
			denom:            "rowan",
			decimals:         18,
			to:               "xxx",
			address:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAmount: sdk.NewUint(0),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(0),
			liabilitiesI:     sdk.NewUint(0),
			health:           sdk.NewDec(1),
			repayAmount:      sdk.NewUint(1000),
		},
		{
			name:             "collateral and native assets are equal",
			denom:            "rowan",
			decimals:         18,
			to:               "rowan",
			address:          "xxx",
			collateralAmount: sdk.NewUint(0),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(1000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			repayAmount:      sdk.NewUint(0),
			err:              nil,
		},
		{
			name:             "mtp not found",
			denom:            "rowan",
			decimals:         18,
			to:               "rowan",
			address:          "xxx",
			overrideAddress:  "yyy",
			collateralAmount: sdk.NewUint(0),
			custodyAmount:    sdk.NewUint(1000),
			liabilitiesP:     sdk.NewUint(1000),
			liabilitiesI:     sdk.NewUint(1000),
			health:           sdk.NewDec(1),
			repayAmount:      sdk.NewUint(0),
			err:              types.ErrMTPDoesNotExist,
		},
	}

	for _, tt := range repayTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			mtp := addMTPKey(t, ctx, app, marginKeeper, tt.to, "xxx", tt.address, types.Position_LONG, 1)
			mtp.CustodyAmount = tt.custodyAmount
			mtp.LiabilitiesP = tt.liabilitiesP
			mtp.CollateralAmount = tt.collateralAmount
			mtp.LiabilitiesI = tt.liabilitiesI
			if tt.overrideAddress != "" {
				mtp.Address = tt.overrideAddress
			}

			got := marginKeeper.Repay(ctx, &mtp, pool, tt.repayAmount, false)

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

func TestKeeper_UpdateMTPInterestLiabilities(t *testing.T) {
	ctx, app, marginKeeper := initKeeper(t)

	mtp := addMTPKey(t, ctx, app, marginKeeper, "rowan", "xxx", "xxx", types.Position_LONG, 1)

	got := marginKeeper.UpdateMTPInterestLiabilities(ctx, &mtp, sdk.NewDec(1.0))
	require.Nil(t, got)
}

func TestKeeper_UpdateMTPInterestLiabilitiesOverflow(t *testing.T) {
	ctx, _, marginKeeper := initKeeper(t)

	mtp := types.MTP{
		Address:          "sif123",
		CollateralAsset:  "rowan",
		CollateralAmount: sdk.Uint{},
		LiabilitiesP:     sdk.NewUintFromString("100"),
		LiabilitiesI:     sdk.NewUintFromString("45231284858326638837332416019018714005183587760015845327913118753091066265500"),
		CustodyAsset:     "ceth",
		CustodyAmount:    sdk.Uint{},
		Leverage:         sdk.Dec{},
		MtpHealth:        sdk.Dec{},
		Position:         types.Position_LONG,
		Id:               1,
	}

	got := marginKeeper.UpdateMTPInterestLiabilities(ctx, &mtp, sdk.NewDec(3.0))
	require.Nil(t, got)
}

func TestKeeper_InterestRateComputation(t *testing.T) {
	interestRateComputationTests := []struct {
		name                 string
		denom                string
		decimals             int64
		interestRate         sdk.Dec
		interestRateIncrease sdk.Dec
		interestRateDecrease sdk.Dec
		interestRateMax      sdk.Dec
		err                  error
		errString            error
	}{
		{
			name:                 "interest rate change lesser than decrease and increase",
			denom:                "unregistred_denom",
			decimals:             18,
			interestRate:         sdk.NewDec(1),
			interestRateIncrease: sdk.NewDec(1),
			interestRateDecrease: sdk.NewDec(1),
			interestRateMax:      sdk.NewDec(5),
			err:                  nil,
		},
		{
			name:                 "interest rate change greater than increase",
			denom:                "unregistred_denom",
			decimals:             18,
			interestRate:         sdk.NewDec(10),
			interestRateIncrease: sdk.NewDec(1),
			interestRateDecrease: sdk.NewDec(1),
			interestRateMax:      sdk.NewDec(5),
			err:                  nil,
		},
		{
			name:                 "interest rate greater than rate max",
			denom:                "unregistred_denom",
			decimals:             18,
			interestRate:         sdk.NewDec(10),
			interestRateIncrease: sdk.NewDec(1),
			interestRateDecrease: sdk.NewDec(1),
			interestRateMax:      sdk.NewDec(0),
			err:                  nil,
		},
	}

	for _, tt := range interestRateComputationTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			data := types.GenesisState{Params: &types.Params{
				LeverageMax:          sdk.NewDec(10),
				InterestRateMax:      tt.interestRateMax,
				InterestRateMin:      sdk.NewDec(1),
				InterestRateIncrease: tt.interestRateIncrease,
				InterestRateDecrease: tt.interestRateDecrease,
				HealthGainFactor:     sdk.NewDec(1),
				EpochLength:          1,
				ForceCloseThreshold:  sdk.NewDec(1), //TODO get real default
			}}
			marginKeeper.InitGenesis(ctx, data)

			asset := clptypes.Asset{Symbol: "rowan"}
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
				InterestRate:         tt.interestRate,
			}

			_, got := marginKeeper.InterestRateComputation(ctx, pool)

			t.Logf("got %v", got)

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

func TestSQBeginBlocker(t *testing.T) {
	ctx, _, marginKeeper := initKeeper(t)
	params := marginKeeper.GetParams(ctx)
	params.PoolOpenThreshold = sdk.NewDec(1)
	marginKeeper.SetParams(ctx, &params)
	pool := clptypes.Pool{
		ExternalAsset: &clptypes.Asset{Symbol: "ceth"},
		Health:        sdk.NewDec(2),
	}

	marginKeeper.TrackSQBeginBlock(ctx, &pool)
	require.Equal(t, uint64(0), marginKeeper.GetSQBeginBlock(ctx, &pool))

	pool.Health = sdk.NewDecWithPrec(1, 1)
	marginKeeper.TrackSQBeginBlock(ctx.WithBlockHeight(1), &pool)
	require.Equal(t, uint64(1), marginKeeper.GetSQBeginBlock(ctx, &pool))

	pool.Health = sdk.NewDec(2)
	marginKeeper.TrackSQBeginBlock(ctx.WithBlockHeight(2), &pool)
	require.Equal(t, uint64(0), marginKeeper.GetSQBeginBlock(ctx, &pool))
}

func initKeeper(t testing.TB) (sdk.Context, *sifapp.SifchainApp, types.Keeper) {
	ctx, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper
	require.NotNil(t, marginKeeper)
	return ctx, app, marginKeeper
}
func addMTPKey(t testing.TB, ctx sdk.Context, app *sifapp.SifchainApp, marginKeeper types.Keeper, collateralAsset string, custodyAsset string, address string, position types.Position, id uint64) types.MTP {
	storeKey := app.GetKey(types.StoreKey)
	store := ctx.KVStore(storeKey)
	key := types.GetMTPKey(address, id)

	newMTP := types.MTP{
		Id:               id,
		Address:          address,
		CollateralAsset:  collateralAsset,
		LiabilitiesP:     sdk.NewUint(1000),
		LiabilitiesI:     sdk.NewUint(1000),
		CollateralAmount: sdk.NewUint(1000),
		CustodyAsset:     custodyAsset,
		CustodyAmount:    sdk.NewUint(1000),
		Leverage:         sdk.NewDec(10),
		MtpHealth:        sdk.NewDec(20),
		Position:         position}

	store.Set(key, types.ModuleCdc.MustMarshal(&newMTP))

	return newMTP
}
