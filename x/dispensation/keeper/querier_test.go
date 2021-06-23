package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/app"
	dispensationkeeper "github.com/Sifchain/sifnode/x/dispensation/keeper"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
)

func GenerateQueryData(app *app.SifchainApp, ctx sdk.Context, name string, outList []bank.Output) {
	keeper := app.DispensationKeeper
	for i := 0; i < 10; i++ {
		name := uuid.New().String()
		distribution := types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name)
		_ = keeper.SetDistribution(ctx, distribution)
	}

	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), int64(-1))
		_ = keeper.SetDistributionRecord(ctx, record)
	}

}

func TestQueryRecordsName(t *testing.T) {
	sifapp, ctx := test.CreateTestApp(false)
	name := uuid.New().String()
	outList := test.CreatOutputList(3, "1000000000")
	GenerateQueryData(sifapp, ctx, name, outList)
	keeper := sifapp.DispensationKeeper
	querier := dispensationkeeper.NewLegacyQuerier(keeper)
	queryRecName := types.QueryRecordsByDistributionNameRequest{
		DistributionName: name,
		Status:           types.DistributionStatus_DISTRIBUTION_STATUS_PENDING,
	}
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	qp, errRes := sifapp.LegacyAmino().MarshalJSON(&queryRecName)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	res, err := querier(ctx, []string{types.QueryRecordsByDistrName}, query)
	require.NoError(t, err)
	var dr types.DistributionRecords
	err = sifapp.LegacyAmino().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr.DistributionRecords, 3)
}

func TestQueryRecordsAddr(t *testing.T) {
	sifapp, ctx := test.CreateTestApp(false)
	name := uuid.New().String()
	outList := test.CreatOutputList(3, "1000000000")
	GenerateQueryData(sifapp, ctx, name, outList)
	keeper := sifapp.DispensationKeeper
	querier := dispensationkeeper.NewLegacyQuerier(keeper)
	quereyRecName := types.QueryRecordsByRecipientAddrRequest{
		Address: outList[0].Address,
	}
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	qp, errRes := sifapp.LegacyAmino().MarshalJSON(&quereyRecName)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	res, err := querier(ctx, []string{types.QueryRecordsByRecipient}, query)
	assert.NoError(t, err)
	var dr types.DistributionRecords
	err = sifapp.LegacyAmino().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr.DistributionRecords, 1)
}

func TestQueryAllDistributions(t *testing.T) {
	sifapp, ctx := test.CreateTestApp(false)
	name := uuid.New().String()
	outList := test.CreatOutputList(3, "1000000000")
	GenerateQueryData(sifapp, ctx, name, outList)
	keeper := sifapp.DispensationKeeper
	querier := dispensationkeeper.NewLegacyQuerier(keeper)
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	query.Path = ""
	query.Data = nil
	res, err := querier(ctx, []string{types.QueryAllDistributions}, query)
	assert.NoError(t, err)
	var dr types.Distributions
	err = sifapp.LegacyAmino().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr.Distributions, 10)
}

func TestQueryClaims(t *testing.T) {
	testApp, ctx := test.CreateTestApp(false)
	keeper := testApp.DispensationKeeper
	claimsVS := test.CreateClaimsList(1000, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	for _, claim := range claimsVS {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	claimsLM := test.CreateClaimsList(1000, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING)
	for _, claim := range claimsLM {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	// Query by type ValidatorSubsidy
	queryData := types.QueryClaimsByTypeRequest{UserClaimType: types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY}
	qp, errRes := testApp.LegacyAmino().MarshalJSON(&queryData)
	require.NoError(t, errRes)
	query := abci.RequestQuery{
		Path: "",
		Data: qp,
	}
	querier := dispensationkeeper.NewLegacyQuerier(keeper)
	res, err := querier(ctx, []string{types.QueryClaimsByType}, query)
	assert.NoError(t, err)
	var dr types.QueryClaimsResponse
	err = testApp.LegacyAmino().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr.Claims, 1000)
}
