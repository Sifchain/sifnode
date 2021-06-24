package dispensation_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExportGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(1000, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name))
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)
	assert.Len(t, genState.DistributionRecords.DistributionRecords, 1000)
	assert.Len(t, genState.Distributions.Distributions, 1)
}

func TestInitGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	app2, ctx2 := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	keeper2 := app2.DispensationKeeper
	outList := test.CreatOutputList(1000, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name))
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)
	assert.Len(t, keeper2.GetDistributions(ctx2).Distributions, 0)
	assert.Len(t, keeper2.GetRecords(ctx2).DistributionRecords, 0)
	dispensation.InitGenesis(ctx2, keeper2, genState)
	assert.Len(t, keeper2.GetDistributions(ctx2).Distributions, 1)
	assert.Len(t, keeper2.GetRecords(ctx2).DistributionRecords, 1000)
}

func TestValidateGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(1000, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name))
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)
	assert.Len(t, genState.DistributionRecords.DistributionRecords, 1000)
	assert.Len(t, genState.Distributions.Distributions, 1)
	assert.NoError(t, dispensation.ValidateGenesis(genState))
}
