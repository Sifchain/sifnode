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
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	require.Empty(t, registry)
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
					Denom:    "TestDenom",
					Decimals: 18,
				},
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				registry = app.TokenRegistryKeeper.GetRegistry(ctx)
				require.Len(t, registry.Entries, 1)
				d := app.TokenRegistryKeeper.GetDenom(registry, "TestDenom")
				require.Equal(t, "TestDenom", d.Denom)
			},
		},
		{
			name: "Successful IBC Registration",
			msg: types.MsgRegister{
				From: admin,
				Entry: &types.RegistryEntry{
					Denom:    "TestDenomIBC",
					Decimals: 18,
				},
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				registry = app.TokenRegistryKeeper.GetRegistry(ctx)
				require.Len(t, registry.Entries, 2)
				d := app.TokenRegistryKeeper.GetDenom(registry, "TestDenomIBC")
				require.Equal(t, "TestDenomIBC", d.Denom)
			},
		},
		{
			name: "Successful IBC Registration 2",
			msg: types.MsgRegister{
				From: admin,
				Entry: &types.RegistryEntry{
					Denom:    "TestDenomIBC2",
					Decimals: 8,
				},
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				registry = app.TokenRegistryKeeper.GetRegistry(ctx)
				require.Len(t, registry.Entries, 3)
				d := app.TokenRegistryKeeper.GetDenom(registry, "TestDenomIBC2")
				require.Equal(t, "TestDenomIBC2", d.Denom)
			},
		},
		{
			name: "Registration Ignore Duplicate Tokens",
			msg: types.MsgRegister{
				From: admin,
				Entry: &types.RegistryEntry{
					Denom:    "TestDenomIBC",
					Decimals: 18,
				},
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				// Denom whitelist size is still 3, duplicate denoms are ignored.
				registry = app.TokenRegistryKeeper.GetRegistry(ctx)
				require.Len(t, registry.Entries, 3)
				d := app.TokenRegistryKeeper.GetDenom(registry, "TestDenomIBC")
				require.Equal(t, "TestDenomIBC", d.Denom)
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
		Denom:    "tokenToRemove",
		Decimals: 18,
	})
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:    "ibcTokenToRemove",
		Decimals: 18,
	})
	app.TokenRegistryKeeper.SetToken(ctx, &types.RegistryEntry{
		Denom:    "ibcTokenToRemove2",
		Decimals: 8,
	})
	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
	require.Len(t, registry.Entries, 3)
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
				registry = app.TokenRegistryKeeper.GetRegistry(ctx)
				require.Len(t, registry.Entries, 2)
				tokenToRemove := app.TokenRegistryKeeper.GetDenom(registry, "tokenToRemove")
				require.Nil(t, tokenToRemove)
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
				registry = app.TokenRegistryKeeper.GetRegistry(ctx)
				require.Len(t, registry.Entries, 1)
				ibcTokenToRemove := app.TokenRegistryKeeper.GetDenom(registry, "ibcTokenToRemove")
				require.Nil(t, ibcTokenToRemove)
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
				registry = app.TokenRegistryKeeper.GetRegistry(ctx)
				require.Empty(t, registry.Entries)
				ibcTokenToRemove2 := app.TokenRegistryKeeper.GetDenom(registry, "ibcTokenToRemove2")
				require.Nil(t, ibcTokenToRemove2)
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
