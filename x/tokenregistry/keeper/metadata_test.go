package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

const invalidDenom = "A Nonexistant Denom Hash"

var testTokenMetadata = types.TokenMetadata{
	Decimals:          15,
	Name:              "Test Token",
	NetworkDescriptor: oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN,
	Symbol:            "TT",
	TokenAddress:      "0x0123456789ABCDEF",
}

func TestGetAddTokenMetadata(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	expected := types.TokenMetadata{}
	result, _ := keeper.GetTokenMetadata(ctx, invalidDenom)
	require.Equal(t, expected, result)

	expectedDenom := ethbridgetypes.GetDenomHash(
		testTokenMetadata.NetworkDescriptor,
		testTokenMetadata.TokenAddress,
		testTokenMetadata.Decimals,
		testTokenMetadata.Name,
		testTokenMetadata.Symbol,
	)

	entry := types.RegistryEntry{
		Denom:         expectedDenom,
		DisplaySymbol: testTokenMetadata.Symbol,
		Decimals:      testTokenMetadata.Decimals,
		IsWhitelisted: true,
		Permissions:   []types.Permission{types.Permission_PERMISSION_CLP},
	}
	keeper.GetTokenRegistryKeeper().SetToken(ctx, &entry)

	resultDenom := keeper.AddTokenMetadata(ctx, testTokenMetadata)

	require.Equal(t, expectedDenom, resultDenom)
	result, _ = keeper.GetTokenMetadata(ctx, resultDenom)
	require.Equal(t, testTokenMetadata, result)
}
