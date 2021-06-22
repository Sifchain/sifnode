package keeper_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	ethbridgekeeper "github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

//nolint:lll
const (
	TestResponseJSON = "{\"prophecy_id\":\"W\\ufffdM(\\ufffd\\u0026\\ufffd\\ufffd%\\u0006\\ufffd\\u000b\\ufffd\\u0014\\ufffd\\ufffd\\ufffd\\ufffd\\ufffd\\u000c\\ufffd\\ufffd\\u001a\\u0017\\ufffd\\ufffd:@]\\ufffdy\\ufffd\",\"status\":1,\"claim_validators\":[\"cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk\"]}"
)

func TestNewQuerier(t *testing.T) {
	ctx, keeper, _, _, _, encCfg, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

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
	ctx, keeper, _, _, oracleKeeper, encCfg, validatorAddresses := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	valAddress := validatorAddresses[0]
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)
	testBridgeContractAddress := types.NewEthereumAddress(types.TestBridgeContractAddress)
	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)

	initialEthBridgeClaim := types.CreateTestEthClaim(
		t, testBridgeContractAddress, testTokenContractAddress, valAddress,
		testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK)

	_, err := oracleKeeper.ProcessClaim(ctx, initialEthBridgeClaim.GetProphecyID(), initialEthBridgeClaim.ValidatorAddress)
	require.NoError(t, err)

	testResponse := types.CreateTestQueryEthProphecyResponse(t, valAddress, types.ClaimType_CLAIM_TYPE_LOCK)

	//Test query String()
	testJSON, err := encCfg.Amino.MarshalJSON(testResponse)
	require.NoError(t, err)
	require.Equal(t, TestResponseJSON, string(testJSON))

	req := types.NewQueryEthProphecyRequest(initialEthBridgeClaim.GetProphecyID())
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
	// badEthereumAddress := types.NewEthereumAddress("badEthereumAddress")

	bz2, err6 := encCfg.Amino.MarshalJSON(types.NewQueryEthProphecyRequest(types.TestProphecyID))
	require.Nil(t, err6)

	query2 := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz2,
	}

	_, err7 := querier(ctx, []string{types.QueryEthProphecy}, query2)
	require.NotNil(t, err7)

	// Test error with empty address
	// emptyEthereumAddress := types.NewEthereumAddress("")

	bz3, err8 := encCfg.Amino.MarshalJSON(
		types.NewQueryEthProphecyRequest(types.TestProphecyID))

	require.Nil(t, err8)

	query3 := abci.RequestQuery{
		Path: "/custom/oracle/prophecies",
		Data: bz3,
	}

	_, err9 := querier(ctx, []string{types.QueryEthProphecy}, query3)
	require.NotNil(t, err9)
}
