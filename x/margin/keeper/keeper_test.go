package keeper_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sifchain/sifnode/app"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_Errors(t *testing.T) {
	_, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper
	assert.NotNil(t, marginKeeper)
}

func TestKeeper_SetMTP(t *testing.T) {
	t.Run("missed defining asset", func(t *testing.T) {
		ctx, _, marginKeeper := initKeeper(t)
		mtp := types.MTP{}
		err := marginKeeper.SetMTP(ctx, &mtp)
		assert.EqualError(t, err, "no asset specified: mtp invalid")
	})
	t.Run("define asset but no address", func(t *testing.T) {
		ctx, _, marginKeeper := initKeeper(t)
		mtp := types.MTP{CollateralAsset: "xxx"}
		err := marginKeeper.SetMTP(ctx, &mtp)
		assert.EqualError(t, err, "no address specified: mtp invalid")
	})
	t.Run("define asset and address", func(t *testing.T) {
		ctx, _, marginKeeper := initKeeper(t)
		mtp := types.MTP{CollateralAsset: "xxx", Address: "xxx"}
		err := marginKeeper.SetMTP(ctx, &mtp)
		assert.NoError(t, err)
	})
}

func TestKeeper_GetMTP(t *testing.T) {
	t.Run("get MTP from a store key that exists", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		want := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "xxx")
		got, err := marginKeeper.GetMTP(ctx, "ceth", "xxx", want.Address)

		if err != nil {
			t.Error("got an error but didn't want one")
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("fails when store keys does not exist", func(t *testing.T) {
		ctx, _, marginKeeper := initKeeper(t)
		marginKeeper.GetMTP(ctx, "ceth", "xxx", "xxx")
	})
}

func TestKeeper_GetMTPIterator(t *testing.T) {
	ctx, app, marginKeeper := initKeeper(t)
	want := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "xxx")
	iterator := marginKeeper.GetMTPIterator(ctx)
	bytesValue := iterator.Value()
	var got types.MTP
	types.ModuleCdc.MustUnmarshal(bytesValue, &got)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestKeeper_GetMTPs(t *testing.T) {
	ctx, app, marginKeeper := initKeeper(t)
	key1 := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "key1")
	key2 := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "key2")
	want := []*types.MTP{&key1, &key2}
	got := marginKeeper.GetMTPs(ctx)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestKeeper_GetMTPsForAsset(t *testing.T) {
	ctx, app, marginKeeper := initKeeper(t)
	key1 := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "key1")
	key2 := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "key2")
	want := []*types.MTP{&key1, &key2}
	got := marginKeeper.GetMTPsForAsset(ctx, "ceth")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestKeeper_GetAssetsForMTP(t *testing.T) {
	ctx, app, marginKeeper := initKeeper(t)
	addr, _ := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")

	addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", addr.String())
	addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", addr.String())
	want := []string{"ceth"}
	got := marginKeeper.GetAssetsForMTP(ctx, addr)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestKeeper_DestroyMTP(t *testing.T) {
	t.Run("key does not exist", func(t *testing.T) {
		ctx, _, marginKeeper := initKeeper(t)
		got := marginKeeper.DestroyMTP(ctx, "ceth", "xxx", "xxx")

		assertError(t, got, types.ErrMTPDoesNotExist)
	})
	t.Run("key exists", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "ceth", "xxx", "xxx")
		got := marginKeeper.DestroyMTP(ctx, "ceth", "xxx", mtp.Address)

		assertNoError(t, got)
	})
}

func TestKeeper_ClpKeeper(t *testing.T) {
	_, _, marginKeeper := initKeeper(t)
	marginKeeper.ClpKeeper()
}

func TestKeeper_BankKeeper(t *testing.T) {
	_, _, marginKeeper := initKeeper(t)
	marginKeeper.BankKeeper()
}

func TestKeeper_GetLeverageParam(t *testing.T) {
	ctx, _, marginKeeper := initKeeper(t)
	marginKeeper.GetLeverageParam(ctx)
}

func TestKeeper_CustodySwap(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			_, got := marginKeeper.CustodySwap(ctx, pool, tt.to, tt.sentAmount)

			if tt.err == nil {
				assertNoError(t, got)
			} else {
				assertError(t, got, tt.err)
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
		borrowAmount     sdk.Uint
		leverage         sdk.Uint
		health           sdk.Dec
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
			borrowAmount:     sdk.NewUint(1000),
			leverage:         sdk.NewUint(1),
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
			borrowAmount:     sdk.NewUint(1000),
			leverage:         sdk.NewUint(1),
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
			borrowAmount:     sdk.NewUint(1000),
			leverage:         sdk.NewUint(1),
			health:           sdk.NewDec(1),
			err:              tokenregistrytypes.ErrNotFound,
		},
		{
			name:             "not enough received asset tokens to swap",
			denom:            "rowan",
			decimals:         18,
			to:               "rowan",
			address:          "xxx",
			collateralAmount: sdk.NewUint(1000000000000000),
			borrowAmount:     sdk.NewUint(1000000000000000),
			leverage:         sdk.NewUint(1),
			health:           sdk.NewDec(1),
			err:              clptypes.ErrNotEnoughAssetTokens,
		},
		{
			name:             "invalid address",
			denom:            "rowan",
			decimals:         9,
			to:               "xxx",
			address:          "xxx",
			collateralAmount: sdk.NewUint(10000),
			borrowAmount:     sdk.NewUint(1000),
			leverage:         sdk.NewUint(1),
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
			borrowAmount:     sdk.NewUint(1000),
			leverage:         sdk.NewUint(1),
			health:           sdk.NewDec(1),
			errString:        errors.New("0xxx is smaller than 10000xxx: insufficient funds"),
		},
		{
			name:             "insufficient funds",
			denom:            "rowan",
			decimals:         9,
			to:               "rowan",
			address:          "sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v",
			collateralAmount: sdk.NewUint(1000),
			borrowAmount:     sdk.NewUint(1000),
			leverage:         sdk.NewUint(1),
			health:           sdk.NewDec(1),
			errString:        errors.New("0xxx is smaller than 10000xxx: insufficient funds"),
		},
	}

	for _, tt := range borrowTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			mtp := addMTPKey(t, ctx, app, marginKeeper, tt.to, "xxx", tt.address)

			got := marginKeeper.Borrow(ctx, tt.to, tt.collateralAmount, tt.borrowAmount, mtp, pool, tt.leverage)

			if tt.errString != nil {
				assertErrorString(t, got, tt.errString)
			} else if tt.err != nil {
				assertError(t, got, tt.err)
			} else {
				assertNoError(t, got)
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
	assert.Nil(t, err)
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
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			mtp := addMTPKey(t, ctx, app, marginKeeper, tt.to, "xxx", "xxx")
			mtp.CustodyAmount = tt.custodyAmount
			mtp.LiabilitiesP = tt.liabilitiesP
			mtp.CollateralAmount = tt.collateralAmount
			mtp.LiabilitiesI = tt.liabilitiesI

			_, got := marginKeeper.UpdateMTPHealth(ctx, mtp, pool)

			if tt.errString != nil {
				assertErrorString(t, got, tt.errString)
			} else if tt.err != nil {
				assertError(t, got, tt.err)
			} else {
				assertNoError(t, got)
			}
		})
	}
}

func TestKeeper_TestInCustody(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}
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

	t.Run("settlement asset and mtp asset is equal", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "rowan", "xxx", "xxx")

		got := marginKeeper.TakeInCustody(ctx, mtp, pool)

		assertNoError(t, got)
	})

	t.Run("settlement asset and mtp asset is not equal", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "notrowan", "xxx", "xxx")

		got := marginKeeper.TakeInCustody(ctx, mtp, pool)

		assertNoError(t, got)
	})
}

func TestKeeper_TestOutCustody(t *testing.T) {
	asset := clptypes.Asset{Symbol: "rowan"}
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

	t.Run("settlement asset and mtp asset is equal", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "rowan", "xxx", "xxx")

		got := marginKeeper.TakeOutCustody(ctx, mtp, pool)

		assertNoError(t, got)
	})

	t.Run("settlement asset and mtp asset is not equal", func(t *testing.T) {
		ctx, app, marginKeeper := initKeeper(t)
		mtp := addMTPKey(t, ctx, app, marginKeeper, "notrowan", "xxx", "xxx")

		got := marginKeeper.TakeOutCustody(ctx, mtp, pool)

		assertNoError(t, got)
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
			name:             "cannot affort principle liability",
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
			name:             "can affort both",
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
			errString:        errors.New("0xxx is smaller than 1000xxx: insufficient funds"),
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
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			mtp := addMTPKey(t, ctx, app, marginKeeper, tt.to, "xxx", tt.address)
			mtp.CustodyAmount = tt.custodyAmount
			mtp.LiabilitiesP = tt.liabilitiesP
			mtp.CollateralAmount = tt.collateralAmount
			mtp.LiabilitiesI = tt.liabilitiesI
			if tt.overrideAddress != "" {
				mtp.Address = tt.overrideAddress
			}

			got := marginKeeper.Repay(ctx, mtp, pool, tt.repayAmount)

			if tt.errString != nil {
				assertErrorString(t, got, tt.errString)
			} else if tt.err != nil {
				assertError(t, got, tt.err)
			} else {
				assertNoError(t, got)
			}
		})
	}
}

func TestKeeper_UpdateMTPInterestLiabilities(t *testing.T) {
	ctx, app, marginKeeper := initKeeper(t)

	mtp := addMTPKey(t, ctx, app, marginKeeper, "rowan", "xxx", "xxx")

	got := marginKeeper.UpdateMTPInterestLiabilities(ctx, &mtp, sdk.NewDec(1.0))
	assert.Nil(t, got)
}

func TestKeeper_InterestRateComputation(t *testing.T) {
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

	interestRateComputationTests := []struct {
		name      string
		denom     string
		decimals  int64
		err       error
		errString error
	}{
		{
			name:     "denom not registered",
			denom:    "unregistred_denom",
			decimals: 18,
			err:      tokenregistrytypes.ErrNotFound,
		},
	}

	for _, tt := range interestRateComputationTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, app, marginKeeper := initKeeper(t)

			app.TokenRegistryKeeper.SetToken(ctx, &tokenregistrytypes.RegistryEntry{
				Denom:       tt.denom,
				Decimals:    tt.decimals,
				Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			})

			_, got := marginKeeper.InterestRateComputation(ctx, pool)

			if tt.errString != nil {
				assertErrorString(t, got, tt.errString)
			} else if tt.err != nil {
				assertError(t, got, tt.err)
			} else {
				assertNoError(t, got)
			}
		})
	}
}

func initKeeper(t testing.TB) (sdk.Context, *app.SifchainApp, types.Keeper) {
	ctx, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper
	assert.NotNil(t, marginKeeper)
	return ctx, app, marginKeeper
}
func addMTPKey(t testing.TB, ctx sdk.Context, app *app.SifchainApp, marginKeeper types.Keeper, collateralAsset string, custodyAsset string, address string) types.MTP {
	storeKey := app.GetKey(types.StoreKey)
	store := ctx.KVStore(storeKey)
	key := types.GetMTPKey(collateralAsset, custodyAsset, address)

	newMTP := types.MTP{
		Address:          address,
		CollateralAsset:  collateralAsset,
		LiabilitiesP:     sdk.NewUint(1000),
		LiabilitiesI:     sdk.NewUint(1000),
		CollateralAmount: sdk.NewUint(1000),
		CustodyAsset:     custodyAsset,
		CustodyAmount:    sdk.NewUint(1000),
		Leverage:         sdk.NewUint(10),
		MtpHealth:        sdk.NewDec(20)}
	store.Set(key, types.ModuleCdc.MustMarshal(&newMTP))

	return newMTP
}
func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
func assertErrorString(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}
	if got.Error() != want.Error() {
		t.Errorf("got %q, want %q", got, want)
	}
}
func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}
