package types_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
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
				From: admin.String(),
				Entry: &types.RegistryEntry{
					Denom:    "TestDenom",
					Decimals: 18,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Denom Missing",
			msg: types.MsgRegister{
				From: admin.String(),
				Entry: &types.RegistryEntry{
					Denom:    "",
					Decimals: 18,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Negative Decimals",
			msg: types.MsgRegister{
				From: admin.String(),
				Entry: &types.RegistryEntry{
					Denom:    "TestDenom",
					Decimals: -1,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Empty from",
			msg: types.MsgRegister{
				From: "",
				Entry: &types.RegistryEntry{
					Denom:    "TestDenom",
					Decimals: 0,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid Denom",
			msg: types.MsgRegister{
				From: admin.String(),
				Entry: &types.RegistryEntry{
					Denom:    "Test%%%$$%%Denom",
					Decimals: 0,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.msg.ValidateBasic(), "")
		})
	}
}

func TestTokenMetadataAddRequest_ValidateBasic(t *testing.T) {
	admin := sdk.AccAddress("addr1_______________")
	tests := []struct {
		name      string
		msg       types.TokenMetadataAddRequest
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Valid Test",
			msg: types.TokenMetadataAddRequest{
				CosmosSender: admin.String(),
				Metadata: &types.TokenMetadata{
					Name:         "TestDenom",
					Symbol:       "TestSymbol",
					Decimals:     18,
					TokenAddress: "0x0123456789012345678901234567890123456789",
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Sender is empty",
			msg: types.TokenMetadataAddRequest{
				CosmosSender: "",
				Metadata: &types.TokenMetadata{
					Name:         "TestDenom",
					Symbol:       "TestSymbol",
					Decimals:     18,
					TokenAddress: "0x0123456789012345678901234567890123456789",
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Token name is empty",
			msg: types.TokenMetadataAddRequest{
				CosmosSender: admin.String(),
				Metadata: &types.TokenMetadata{
					Name:         "",
					Symbol:       "TestSymbol",
					Decimals:     18,
					TokenAddress: "0x0123456789012345678901234567890123456789",
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Token symbol is empty",
			msg: types.TokenMetadataAddRequest{
				CosmosSender: admin.String(),
				Metadata: &types.TokenMetadata{
					Name:         "TestDenom",
					Symbol:       "",
					Decimals:     18,
					TokenAddress: "0x0123456789012345678901234567890123456789",
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid token address",
			msg: types.TokenMetadataAddRequest{
				CosmosSender: admin.String(),
				Metadata: &types.TokenMetadata{
					Name:         "TestDenom",
					Symbol:       "TestSymbol",
					Decimals:     18,
					TokenAddress: "0x012",
				},
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid network descriptor",
			msg: types.TokenMetadataAddRequest{
				CosmosSender: admin.String(),
				Metadata: &types.TokenMetadata{
					Name:              "TestDenom",
					Symbol:            "TestSymbol",
					Decimals:          18,
					TokenAddress:      "0x0123456789012345678901234567890123456789",
					NetworkDescriptor: 8888,
				},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, tt.msg.ValidateBasic(), "")
		})
	}
}
