package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_GetRecordsForNamePrefixed(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
		assert.NoError(t, err)
	}
	outList2 := test.CreatOutputList(7, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
		assert.NoError(t, err)

	}
	assert.Len(t, keeper.GetRecordsForName(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 3)
	assert.Len(t, keeper.GetRecordsForName(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED).DistributionRecords, 7)
}

func TestKeeper_GetPendingRecordsLimited(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(1000, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
		assert.NoError(t, err)
	}
	outList2 := test.CreatOutputList(7, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
		assert.NoError(t, err)
	}
	assert.Len(t, keeper.GetRecordsLimited(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 10)
	assert.Len(t, keeper.GetRecordsForName(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 1000)
}
func TestKeeper_GetPendingRecordsLimitedMultipleDistributions(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(2, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
		assert.NoError(t, err)
	}
	name = uuid.New().String()
	outList2 := test.CreatOutputList(3, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
		assert.NoError(t, err)
	}
	assert.Len(t, keeper.GetRecordsLimited(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 5)
	assert.Len(t, keeper.GetRecordsForName(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 3)
}

func TestKeeper_GetRecordsForName(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, name, rec.Address, rec.Coins, ctx.BlockHeight(), int64(-1))
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForName(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	assert.Len(t, list.DistributionRecords, 3)
}

func TestKeeper_GetRecordsForRecipient(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, name, rec.Address, rec.Coins, ctx.BlockHeight(), int64(-1))
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForRecipient(ctx, outList[0].Address)
	assert.Len(t, list.DistributionRecords, 1)
}
