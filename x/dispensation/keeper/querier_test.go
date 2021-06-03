package keeper_test

import (
	"github.com/Sifchain/sifnode/simapp"
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func GenerateQueryData(app *simapp.SimApp, ctx sdk.Context, name string, outList []bank.Output) {
	keeper := app.DispensationKeeper
	for i := 0; i < 10; i++ {
		name := uuid.New().String()
		distribution := types.NewDistribution(types.Airdrop, name, sdk.AccAddress{})
		_ = keeper.SetDistribution(ctx, distribution)
	}

	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		_ = keeper.SetDistributionRecord(ctx, record)
	}

}

func TestQueryRecordsName(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	name := uuid.New().String()
	outList := test.GenerateOutputList("1000000000")
	GenerateQueryData(app, ctx, name, outList)
	keeper := app.DispensationKeeper
	querier := dispensation.NewQuerier(keeper)
	quereyRecName := types.QueryRecordsByDistributionName{
		DistributionName: name,
	}
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	qp, errRes := app.Codec().MarshalJSON(quereyRecName)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	res, err := querier(ctx, []string{types.QueryRecordsByDistrName}, query)
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
		Address: outList[0].Address,
	}
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	qp, errRes := app.Codec().MarshalJSON(quereyRecName)
	require.NoError(t, errRes)
	query.Path = ""
	query.Data = qp
	res, err := querier(ctx, []string{types.QueryRecordsByRecipient}, query)
	assert.NoError(t, err)
	var dr types.DistributionRecords
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
	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	query.Path = ""
	query.Data = nil
	res, err := querier(ctx, []string{types.QueryAllDistributions}, query)
	assert.NoError(t, err)
	var dr types.Distributions
	err = keeper.Codec().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr, 10)
}

func TestQueryClaims(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	claimsVS := test.CreateClaimsList(1000, types.ValidatorSubsidy)
	for _, claim := range claimsVS {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	claimsLM := test.CreateClaimsList(1000, types.LiquidityMining)
	for _, claim := range claimsLM {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	// Query by type ValidatorSubsidy
	queryData := types.QueryUserClaims{UserClaimType: types.ValidatorSubsidy}
	qp, errRes := app.Codec().MarshalJSON(queryData)
	require.NoError(t, errRes)
	query := abci.RequestQuery{
		Path: "",
		Data: qp,
	}
	querier := dispensation.NewQuerier(keeper)
	res, err := querier(ctx, []string{types.QueryClaimsByType}, query)
	assert.NoError(t, err)
	var dr []types.UserClaim
	err = keeper.Codec().UnmarshalJSON(res, &dr)
	assert.NoError(t, err)
	assert.Len(t, dr, 1000)
}
