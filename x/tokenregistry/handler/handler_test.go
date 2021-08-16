package handler_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/tokenregistry/handler"
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleRegister(t *testing.T) {
	app, ctx, admin := test.CreateTestApp(false)
	h := handler.NewHandler(app.TokenRegistryKeeper)
	tests := []struct {
		name           string
		msg            types.MsgRegister
		errorAssertion assert.ErrorAssertionFunc
		valueAssertion require.ValueAssertionFunc
	}{
		{
			name: "Successful Registration",
			msg: types.MsgRegister{
				From: admin,
				Entry: &types.RegistryEntry{
					Denom:       "TestDenom",
					DisplayName: "Test Denom",
					Decimals:    18,
				},
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				d := app.TokenRegistryKeeper.GetDenom(ctx, "TestDenom")
				require.Equal(t, "Test Denom", d.DisplayName)
			},
		},
		{
			name: "Successful IBC Registration",
			msg: types.MsgRegister{
				From: admin,
				Entry: &types.RegistryEntry{
					Denom:       "TestDenomIBC",
					DisplayName: "Test Denom IBC",
					Decimals:    18,
					IbcDenom:    "Test Denom IBC",
					IbcDecimals: 10,
				},
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				d := app.TokenRegistryKeeper.GetDenom(ctx, "TestDenomIBC")
				require.Equal(t, "Test Denom IBC", d.DisplayName)
			},
		},
		{
			name: "Successful Registration Converted",
			msg: types.MsgRegister{
				From: admin,
				Entry: &types.RegistryEntry{
					Denom:       "TestDenomIBC2",
					DisplayName: "Test Denom IBC 2",
					Decimals:    8,
					IbcDenom:    "Test Denom IBC 2",
					IbcDecimals: 10,
				},
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				d := app.TokenRegistryKeeper.GetDenom(ctx, "TestDenomIBC2")
				require.Equal(t, "Test Denom IBC 2", d.DisplayName)
			},
		},
		{
			name: "Non Admin Account",
			msg: types.MsgRegister{
				From: sdk.AccAddress("addr2_______________").String(),
				Entry: &types.RegistryEntry{
					Denom:    "TestDenom",
					Decimals: 18,
				},
			},
			errorAssertion: assert.Error,
			valueAssertion: require.Nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res, err := h(ctx, &tt.msg)
			tt.errorAssertion(t, err)
			tt.valueAssertion(t, res)
		})
	}
}

func TestHandleDeregister(t *testing.T) {
	app, ctx, admin := test.CreateTestApp(false)

	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "tokenToRemove",
		Decimals:      18,
	})

	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ibcTokenToRemove",
		Decimals:      18,
		IbcDecimals:   10,
		IbcDenom:      "ibcTokenToRemove",
	})

	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		IsWhitelisted: true,
		Denom:         "ibcTokenToRemove2",
		Decimals:      8,
		IbcDecimals:   10,
		IbcDenom:      "ibcTokenToRemove2",
	})

	require.True(t, app.TokenRegistryKeeper.IsDenomWhitelisted(ctx, "tokenToRemove"))
	require.True(t, app.TokenRegistryKeeper.IsDenomWhitelisted(ctx, "ibcTokenToRemove"))
	require.True(t, app.TokenRegistryKeeper.IsDenomWhitelisted(ctx, "ibcTokenToRemove2"))
	require.Equal(t, 3, len(app.TokenRegistryKeeper.GetDenomWhitelist(ctx).Entries))

	h := handler.NewHandler(app.TokenRegistryKeeper)

	tests := []struct {
		name           string
		msg            types.MsgDeregister
		errorAssertion assert.ErrorAssertionFunc
		valueAssertion require.ValueAssertionFunc
	}{
		{
			name: "Successful De-registration",
			msg: types.MsgDeregister{
				From:  admin,
				Denom: "tokenToRemove",
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				require.False(t, app.TokenRegistryKeeper.IsDenomWhitelisted(ctx, "tokenToRemove"))
				require.Equal(t, 2, len(app.TokenRegistryKeeper.GetDenomWhitelist(ctx).Entries))
			},
		},
		{
			name: "Successful IBC De-registration",
			msg: types.MsgDeregister{
				From:  admin,
				Denom: "ibcTokenToRemove",
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				require.False(t, app.TokenRegistryKeeper.IsDenomWhitelisted(ctx, "ibcTokenToRemove"))
				require.Equal(t, 1, len(app.TokenRegistryKeeper.GetDenomWhitelist(ctx).Entries))
			},
		},
		{
			name: "Successful IBC De-registration 2",
			msg: types.MsgDeregister{
				From:  admin,
				Denom: "ibcTokenToRemove2",
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				require.False(t, app.TokenRegistryKeeper.IsDenomWhitelisted(ctx, "ibcTokenToRemove2"))
				require.Empty(t, app.TokenRegistryKeeper.GetDenomWhitelist(ctx).Entries)
			},
		},
		{
			name: "Non Admin Account",
			msg: types.MsgDeregister{
				From: sdk.AccAddress("addr2_______________").String(),
			},
			errorAssertion: assert.Error,
			valueAssertion: require.Nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res, err := h(ctx, &tt.msg)
			tt.errorAssertion(t, err)
			tt.valueAssertion(t, res)
		})
	}
}
