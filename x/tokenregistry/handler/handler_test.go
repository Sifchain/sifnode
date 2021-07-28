package handler_test

import (
	"github.com/Sifchain/sifnode/x/tokenregistry/handler"
	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
		Denom:         "removeMe",
	})

	require.True(t, app.TokenRegistryKeeper.IsDenomWhitelisted(ctx, "removeMe"))

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
				Denom: "removeMe",
			},
			errorAssertion: assert.NoError,
			valueAssertion: func(t require.TestingT, res interface{}, i ...interface{}) {
				require.False(t, app.TokenRegistryKeeper.IsDenomWhitelisted(ctx, "removeMe"))
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
