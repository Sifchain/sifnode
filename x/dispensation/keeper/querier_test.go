package keeper_test

import (
<<<<<<< HEAD
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
=======
	"github.com/Sifchain/sifnode/simapp"
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
>>>>>>> develop
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
<<<<<<< HEAD

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
=======
	"testing"
)

func GenerateQueryData(app *simapp.SimApp, ctx sdk.Context, name string, outList []bank.Output) {
	keeper := app.DispensationKeeper
	for i := 0; i < 10; i++ {
		name := uuid.New().String()
		distribution := types.NewDistribution(types.Airdrop, name, sdk.AccAddress{})
>>>>>>> develop
		_ = keeper.SetDistribution(ctx, distribution)
	}

	for _, rec := range outList {
<<<<<<< HEAD
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), int64(-1))
=======
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
>>>>>>> develop
		_ = keeper.SetDistributionRecord(ctx, record)
	}

}

func TestQueryRecordsName(t *testing.T) {
<<<<<<< HEAD
	sifapp, ctx := test.CreateTestApp(false)
	name := uuid.New().String()
	outList := test.CreatOutputList(3, "1000000000")
	GenerateQueryData(sifapp, ctx, name, outList)
	keeper := sifapp.DispensationKeeper
	querier := dispensationkeeper.NewLegacyQuerier(keeper)
	queryRecName := types.QueryRecordsByDistributionNameRequest{
		DistributionName: name,
		Status:           types.DistributionStatus_DISTRIBUTION_STATUS_PENDING,
=======
	app, ctx := test.CreateTestApp(false)
	name := uuid.New().String()
	outList := test.GenerateOutputList("1000000000")
	GenerateQueryData(app, ctx, name, outList)
	keeper := app.DispensationKeeper
	querier := dispensation.NewQuerier(keeper)
	quereyRecName := types.QueryRecordsByDistributionName{
		DistributionName: name,
>>>>>>> develop
	}
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
<<<<<<< HEAD
	qp, errRes := sifapp.LegacyAmino().MarshalJSON(&queryRecName)
=======
	qp, errRes := app.Codec().MarshalJSON(quereyRecName)
>>>>>>> develop
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	res, err := querier(ctx, []string{types.QueryRecordsByDistrName}, query)
<<<<<<< HEAD
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
=======
	assert.NoError(t, err)
	var dr types.DistributionRecords
	err = keeper.Codec().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr, 3)
}

func TestQueryRecordsAddr(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	name := uuid.New().String()
	outList := test.GenerateOutputList("1000000000")
	GenerateQueryData(app, ctx, name, outList)
	keeper := app.DispensationKeeper
	querier := dispensation.NewQuerier(keeper)
	quereyRecName := types.QueryRecordsByRecipientAddr{
>>>>>>> develop
		Address: outList[0].Address,
	}
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
<<<<<<< HEAD
	qp, errRes := sifapp.LegacyAmino().MarshalJSON(&quereyRecName)
=======
	qp, errRes := app.Codec().MarshalJSON(quereyRecName)
>>>>>>> develop
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	res, err := querier(ctx, []string{types.QueryRecordsByRecipient}, query)
	assert.NoError(t, err)
	var dr types.DistributionRecords
<<<<<<< HEAD
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
=======
	err = keeper.Codec().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr, 1)
}

func TestQueryAllDistributions(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	name := uuid.New().String()
	outList := test.GenerateOutputList("1000000000")
	GenerateQueryData(app, ctx, name, outList)
	keeper := app.DispensationKeeper
	querier := dispensation.NewQuerier(keeper)
>>>>>>> develop
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	query.Path = ""
	query.Data = nil
	res, err := querier(ctx, []string{types.QueryAllDistributions}, query)
	assert.NoError(t, err)
	var dr types.Distributions
<<<<<<< HEAD
	err = sifapp.LegacyAmino().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr.Distributions, 10)
}

func TestQueryClaims(t *testing.T) {
	testApp, ctx := test.CreateTestApp(false)
	keeper := testApp.DispensationKeeper
	claimsVS := test.CreateClaimsList(1000, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
=======
	err = keeper.Codec().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr, 10)
}

func TestQueryClaims(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	claimsVS := test.CreateClaimsList(1000, types.ValidatorSubsidy)
>>>>>>> develop
	for _, claim := range claimsVS {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
<<<<<<< HEAD
	claimsLM := test.CreateClaimsList(1000, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING)
=======
	claimsLM := test.CreateClaimsList(1000, types.LiquidityMining)
>>>>>>> develop
	for _, claim := range claimsLM {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	// Query by type ValidatorSubsidy
<<<<<<< HEAD
	queryData := types.QueryClaimsByTypeRequest{UserClaimType: types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY}
	qp, errRes := testApp.LegacyAmino().MarshalJSON(&queryData)
=======
	queryData := types.QueryUserClaims{UserClaimType: types.ValidatorSubsidy}
	qp, errRes := app.Codec().MarshalJSON(queryData)
>>>>>>> develop
	require.NoError(t, errRes)
	query := abci.RequestQuery{
		Path: "",
		Data: qp,
	}
<<<<<<< HEAD
	querier := dispensationkeeper.NewLegacyQuerier(keeper)
	res, err := querier(ctx, []string{types.QueryClaimsByType}, query)
	assert.NoError(t, err)
	var dr types.QueryClaimsResponse
	err = testApp.LegacyAmino().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr.Claims, 1000)
=======
	querier := dispensation.NewQuerier(keeper)
	res, err := querier(ctx, []string{types.QueryClaimsByType}, query)
	assert.NoError(t, err)
	var dr []types.UserClaim
	err = keeper.Codec().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr, 1000)
>>>>>>> develop
}
