package keeper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

//nolint:lll
const (
	TestResponseJSON = "{\"id\":\"300x7B95B6EC7EbD73572298cEf32Bb54FA408207359\",\"status\":{\"text\":1},\"claims\":[{\"ethereum_chain_id\":\"3\",\"bridge_contract_address\":\"0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB\",\"symbol\":\"eth\",\"token_contract_address\":\"0x0000000000000000000000000000000000000000\",\"ethereum_sender\":\"0x7B95B6EC7EbD73572298cEf32Bb54FA408207359\",\"cosmos_receiver\":\"cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv\",\"validator_address\":\"cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk\",\"amount\":\"10\",\"claim_type\":1}]}"
)

func TestNewQuerier(t *testing.T) {
	ctx, keeper, _, _, _, encCfg, _ := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewLegacyQuerier(keeper, encCfg.Amino)

	//Test wrong paths
	bz, err := querier(ctx, []string{"other"}, query)
	require.Error(t, err)
	require.Nil(t, bz)
}

func TestQueryEthProphecy(t *testing.T) {
	ctx, keeper, _, _, _, encCfg, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	valAddress := validatorAddresses[0]
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)
	testBridgeContractAddress := types.NewEthereumAddress(types.TestBridgeContractAddress)
	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)

	initialEthBridgeClaim := types.CreateTestEthClaim(
		t, testBridgeContractAddress, testTokenContractAddress, valAddress,
		testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK)
	oracleClaim, _ := types.CreateOracleClaimFromEthClaim(initialEthBridgeClaim)
	_, err := keeper.oracleKeeper.ProcessClaim(ctx, oracleClaim)
	require.NoError(t, err)

	testResponse := types.CreateTestQueryEthProphecyResponse(t, valAddress, types.ClaimType_CLAIM_TYPE_LOCK)

	//Test query String()
	testJSON, err := encCfg.Amino.MarshalJSON(testResponse)
	require.NoError(t, err)
	require.Equal(t, TestResponseJSON, string(testJSON))

	bz, err2 := encCfg.Amino.MarshalJSON(types.NewQueryEthProphecyRequest(
		types.TestEthereumChainID, testBridgeContractAddress, types.TestNonce,
		types.TestCoinsSymbol, testTokenContractAddress, testEthereumAddress))
	require.Nil(t, err2)

	query := abci.RequestQuery{
		Path: "/custom/ethbridge/prophecies",
		Data: bz,
	}

	//Test query
	res, err3 := legacyQueryEthProphecy(ctx, encCfg.Amino, query, keeper)
	require.Nil(t, err3)

	var ethProphecyResp types.QueryEthProphecyResponse
	err4 := encCfg.Amino.UnmarshalJSON(res, &ethProphecyResp)
	require.Nil(t, err4)
	require.True(t, reflect.DeepEqual(ethProphecyResp, testResponse))

	// Test error with bad request
	query.Data = bz[:len(bz)-1]

	_, err5 := legacyQueryEthProphecy(ctx, encCfg.Amino, query, keeper)
	require.NotNil(t, err5)

	// Test error with nonexistent request
	badEthereumAddress := types.NewEthereumAddress("badEthereumAddress")

	bz2, err6 := encCfg.Amino.MarshalJSON(types.NewQueryEthProphecyRequest(
		types.TestEthereumChainID, testBridgeContractAddress, 12,
		types.TestCoinsSymbol, testTokenContractAddress, badEthereumAddress))
	require.Nil(t, err6)

	query2 := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz2,
	}

	_, err7 := legacyQueryEthProphecy(ctx, encCfg.Amino, query2, keeper)
	require.NotNil(t, err7)

	// Test error with empty address
	emptyEthereumAddress := types.NewEthereumAddress("")

	bz3, err8 := encCfg.Amino.MarshalJSON(
		types.NewQueryEthProphecyRequest(
			types.TestEthereumChainID, testBridgeContractAddress, 12,
			types.TestCoinsSymbol, testTokenContractAddress, emptyEthereumAddress))

	require.Nil(t, err8)

	query3 := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz3,
	}

	_, err9 := legacyQueryEthProphecy(ctx, encCfg.Amino, query3, keeper)
	require.NotNil(t, err9)
}
