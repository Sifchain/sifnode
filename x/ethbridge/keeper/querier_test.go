package keeper_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	ethbridgekeeper "github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

//nolint:lll
const (
	TestResponseJSON = "{\"id\":\"100x7B95B6EC7EbD73572298cEf32Bb54FA408207359\",\"status\":{\"text\":1},\"claims\":[{\"network_id\":1,\"bridge_contract_address\":\"0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB\",\"symbol\":\"eth\",\"token_contract_address\":\"0x0000000000000000000000000000000000000000\",\"ethereum_sender\":\"0x7B95B6EC7EbD73572298cEf32Bb54FA408207359\",\"cosmos_receiver\":\"cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv\",\"validator_address\":\"cosmosvaloper1353a4uac03etdylz86tyq9ssm3x2704j3a9n7n\",\"amount\":\"10\",\"claim_type\":2}]}"

	networkID = 1
)

func TestNewQuerier(t *testing.T) {
	ctx, keeper, _, _, _, encCfg, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := ethbridgekeeper.NewLegacyQuerier(keeper, encCfg.Amino)

	//Test wrong paths
	bz, err := querier(ctx, []string{"other"}, query)
	require.Error(t, err)
	require.Nil(t, bz)
}

func TestQueryEthProphecy(t *testing.T) {
	ctx, keeper, _, _, oracleKeeper, encCfg, _, validatorAddresses := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	valAddress := validatorAddresses[0]
	NewTestResponseJSON := strings.Replace(TestResponseJSON, "cosmosvaloper1353a4uac03etdylz86tyq9ssm3x2704j3a9n7n", valAddress.String(), -1)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)
	testBridgeContractAddress := types.NewEthereumAddress(types.TestBridgeContractAddress)
	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)

	initialEthBridgeClaim := types.CreateTestEthClaim(
		t, testBridgeContractAddress, testTokenContractAddress, valAddress,
		testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK)
	oracleClaim, _ := types.CreateOracleClaimFromEthClaim(initialEthBridgeClaim)
	_, err := oracleKeeper.ProcessClaim(ctx, networkID, oracleClaim)
	require.NoError(t, err)

	testResponse := types.CreateTestQueryEthProphecyResponse(t, valAddress, types.ClaimType_CLAIM_TYPE_LOCK)

	//Test query String()
	testJSON, err := encCfg.Amino.MarshalJSON(testResponse)
	require.NoError(t, err)
	require.Equal(t, NewTestResponseJSON, string(testJSON))

	req := types.NewQueryEthProphecyRequest(
		types.TestNetworkID, testBridgeContractAddress, types.TestNonce,
		types.TestCoinsSymbol, testTokenContractAddress, testEthereumAddress)
	bz, err2 := encCfg.Amino.MarshalJSON(req)
	require.Nil(t, err2)

	query := abci.RequestQuery{
		Path: "/custom/ethbridge/prophecies",
		Data: bz,
	}

	//Test query
	querier := ethbridgekeeper.NewLegacyQuerier(keeper, encCfg.Amino)
	res, err3 := querier(ctx, []string{types.QueryEthProphecy}, query)
	require.Nil(t, err3)

	var ethProphecyResp types.QueryEthProphecyResponse
	err4 := encCfg.Amino.UnmarshalJSON(res, &ethProphecyResp)
	require.Nil(t, err4)
	require.True(t, reflect.DeepEqual(ethProphecyResp, testResponse))

	// Test error with bad request
	query.Data = bz[:len(bz)-1]

	_, err5 := querier(ctx, []string{types.QueryEthProphecy}, query)
	require.NotNil(t, err5)

	// Test error with nonexistent request
	badEthereumAddress := types.NewEthereumAddress("badEthereumAddress")

	bz2, err6 := encCfg.Amino.MarshalJSON(types.NewQueryEthProphecyRequest(
		types.TestNetworkID, testBridgeContractAddress, 12,
		types.TestCoinsSymbol, testTokenContractAddress, badEthereumAddress))
	require.Nil(t, err6)

	query2 := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz2,
	}

	_, err7 := querier(ctx, []string{types.QueryEthProphecy}, query2)
	require.NotNil(t, err7)

	// Test error with empty address
	emptyEthereumAddress := types.NewEthereumAddress("")

	bz3, err8 := encCfg.Amino.MarshalJSON(
		types.NewQueryEthProphecyRequest(
			types.TestNetworkID, testBridgeContractAddress, 12,
			types.TestCoinsSymbol, testTokenContractAddress, emptyEthereumAddress))

	require.Nil(t, err8)

	query3 := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz3,
	}

	_, err9 := querier(ctx, []string{types.QueryEthProphecy}, query3)
	require.NotNil(t, err9)
}
