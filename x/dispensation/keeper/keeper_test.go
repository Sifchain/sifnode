package keeper_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_GetDistributions(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	for i := 0; i < 10; i++ {
		name := uuid.New().String()
		distribution := types.NewDistribution(types.Airdrop, name)
		assert.True(t, distribution.Validate())
		err := keeper.SetDistribution(ctx, distribution)
		assert.NoError(t, err)
		_, err = keeper.GetDistribution(ctx, name, types.Airdrop)
		assert.NoError(t, err)
	}
	list := keeper.GetDistributions(ctx)
	assert.Len(t, list, 10)
}

func TestKeeper_GetRecordsForName(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.GenerateOutputList("1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.Pending, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), types.Pending)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForName(ctx, name, types.Pending)
	assert.Len(t, list, 3)
}

func TestKeeper_GetRecordsForRecipient(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.GenerateOutputList("1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.Pending, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), types.Pending)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForRecipient(ctx, outList[0].Address)
	assert.Len(t, list, 1)
}

func TestKeeper_GetRecordsForNamePrefixed(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.Pending, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), types.Pending)
		assert.NoError(t, err)
	}
	outList2 := test.CreatOutputList(7, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.Completed, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), types.Completed)
		assert.NoError(t, err)

	}
	assert.Len(t, keeper.GetRecordsForName(ctx, name, types.Pending), 3)
	assert.Len(t, keeper.GetRecordsForName(ctx, name, types.Completed), 7)
}

func TestKeeper_GetPendingRecordsLimited(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(1000, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.Pending, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), types.Pending)
		assert.NoError(t, err)
	}
	outList2 := test.CreatOutputList(7, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.Completed, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), types.Completed)
		assert.NoError(t, err)
	}
	assert.Len(t, keeper.GetRecordsLimited(ctx, types.Pending), 10)
	assert.Len(t, keeper.GetRecordsForName(ctx, name, types.Pending), 1000)
}
func TestKeeper_GetPendingRecordsLimitedMultipleDistributions(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(2, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.Pending, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), types.Pending)
		assert.NoError(t, err)
	}
	name = uuid.New().String()
	outList2 := test.CreatOutputList(3, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.Pending, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), types.Pending)
		assert.NoError(t, err)
	}
	assert.Len(t, keeper.GetRecordsLimited(ctx, types.Pending), 5)
	assert.Len(t, keeper.GetRecordsForName(ctx, name, types.Pending), 3)
}
