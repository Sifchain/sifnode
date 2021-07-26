package types_test

import (
	"github.com/Sifchain/sifnode/x/whitelist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMsgUpdateWhitelist_ValidateBasic(t *testing.T) {
	admin := sdk.AccAddress("addr1_______________")
	tests := []struct {
		name      string
		msg       types.MsgRegister
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Valid Test",
			msg: types.MsgRegister{
				From:     admin.String(),
				Denom:    "TestDenom",
				Decimals: 18,
			},
			assertion: assert.NoError,
		},
		{
			name: "Denom Missing",
			msg: types.MsgRegister{
				From:     admin.String(),
				Denom:    "",
				Decimals: 18,
			},
			assertion: assert.Error,
		},
		{
			name: "Negative Decimals",
			msg: types.MsgRegister{
				From:     admin.String(),
				Denom:    "TestDenom",
				Decimals: -1,
			},
			assertion: assert.Error,
		},
		{
			name: "Empty from",
			msg: types.MsgRegister{
				From:     "",
				Denom:    "TestDenom",
				Decimals: 0,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid Denom",
			msg: types.MsgRegister{
				From:     admin.String(),
				Denom:    "Test%%%$$%%Denom",
				Decimals: 0,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.msg.ValidateBasic(), "")
		})
	}
}
