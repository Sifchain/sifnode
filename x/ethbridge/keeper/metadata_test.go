package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const invalidDenom = "A Nonexistant Denom Hash"

var testMetadata = types.TokenMetadata{
	Decimals:          15,
	Name:              "Test Token",
	NetworkDescriptor: oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN,
	Symbol:            "TT",
	TokenAddress:      "0x0123456789ABCDEF",
}

func TestGetAddTokenMetadata(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	expected := types.TokenMetadata{}
	result := keeper.GetTokenMetadata(ctx, invalidDenom)
	require.Equal(t, expected, result)
	resultDenom := keeper.AddTokenMetadata(ctx, testMetadata)
	expectedDenom := types.GetDenomHash(
		testMetadata.NetworkDescriptor,
		testMetadata.TokenAddress,
		testMetadata.Decimals,
		testMetadata.Name,
		testMetadata.Symbol,
	)
	require.Equal(t, expectedDenom, resultDenom)
	result = keeper.GetTokenMetadata(ctx, resultDenom)
	require.Equal(t, testMetadata, result)
}

func TestExistsTokenMetadata(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	expected := false
	result := keeper.ExistsTokenMetadata(ctx, invalidDenom)
	require.Equal(t, expected, result)
	denom := keeper.AddTokenMetadata(ctx, testMetadata)
	expected = true
	result = keeper.ExistsTokenMetadata(ctx, denom)
	require.Equal(t, expected, result)
}

func TestAddDeleteTokenMetadata(t *testing.T) {
	ctx, keeper, _, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	testData := testMetadata
	testData.Name = "ibc/" + testData.Name
	expectedDenom := ""
	denom := keeper.AddIBCTokenMetadata(ctx, testData, cosmosSender)
	require.Equal(t, expectedDenom, denom)
	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	testData = testMetadata
	denom = keeper.AddTokenMetadata(ctx, testData)
	expected := false
	result := keeper.DeleteTokenMetadata(ctx, cosmosSender, denom)
	require.Equal(t, expected, result)
	testData.Name = "ibc/" + testData.Name
	denom = keeper.AddIBCTokenMetadata(ctx, testData, cosmosSender)
	expected = true
	result = keeper.DeleteTokenMetadata(ctx, cosmosSender, denom)
	require.Equal(t, expected, result)
}
