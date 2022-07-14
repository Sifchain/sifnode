package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetDistributionRecord_unableToSetRecord(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper

	dr := types.DistributionRecord{
		DistributionStatus:          types.DistributionStatus_DISTRIBUTION_STATUS_FAILED,
		DistributionType:            types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED,
		DistributionName:            "",
		RecipientAddress:            "",
		Coins:                       sdk.Coins{},
		DistributionStartHeight:     int64(0),
		DistributionCompletedHeight: int64(0),
		AuthorizedRunner:            types.AttributeKeyDistributionRunner,
	}
	bool := dr.Validate()
	assert.False(t, bool)
	err := keeper.SetDistributionRecord(ctx, dr)
	assert.Error(t, err)

}

func TestKeeper_GetDistributionRecord(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	key := types.GetDistributionRecordKey(types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, uuid.New().String(), sdk.AccAddress("addr1_______________").String(), types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED)
	res := keeper.Exists(ctx, key)
	assert.False(t, res)
	_, err := keeper.GetDistributionRecord(ctx, uuid.New().String(), sdk.AccAddress("addr1_______________").String(), types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED)
	assert.Error(t, err)
}

func TestKeeper_GetDistributionRecordsIterator_Default(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	status := types.DistributionStatus_DISTRIBUTION_STATUS_UNSPECIFIED
	res := keeper.GetDistributionRecordsIterator(ctx, status)
	assert.Nil(t, res)

}

func TestKeeper_DeleteDistributionRecord_recordNotExist(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	key := types.GetDistributionRecordKey(types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, uuid.New().String(), sdk.AccAddress("addr1_______________").String(), types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED)
	res := keeper.Exists(ctx, key)
	assert.False(t, res)
	err := keeper.DeleteDistributionRecord(ctx, uuid.New().String(), sdk.AccAddress("addr1_______________").String(), types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED)
	assert.Error(t, err)
}

func TestKeeper_GetRecordsForNamePrefixed(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType)
		assert.NoError(t, err)
	}
	outList2 := test.CreatOutputList(7, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, record.DistributionType)
		assert.NoError(t, err)

	}
	assert.Len(t, keeper.GetRecordsForNameAndStatus(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 3)
	assert.Len(t, keeper.GetRecordsForNameAndStatus(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED).DistributionRecords, 7)
}

func TestKeeper_GetPendingRecordsLimited(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(1000, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType)
		assert.NoError(t, err)
	}
	outList2 := test.CreatOutputList(7, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, record.DistributionType)
		assert.NoError(t, err)
	}
	assert.Len(t, keeper.GetLimitedRecordsForStatus(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 20)
	assert.Len(t, keeper.GetRecordsForNameAndStatus(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 1000)
}
func TestKeeper_GetPendingRecordsLimitedMultipleDistributions(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList1 := test.CreatOutputList(2, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList1 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType)
		assert.NoError(t, err)
	}
	name = uuid.New().String()
	outList2 := test.CreatOutputList(3, "1000000000")
	for _, rec := range outList2 {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType)
		assert.NoError(t, err)
	}
	assert.Len(t, keeper.GetLimitedRecordsForStatus(ctx, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 5)
	assert.Len(t, keeper.GetRecordsForNameAndStatus(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords, 3)
}

func TestKeeper_GetRecordsForName(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForNameAndStatus(ctx, name, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	assert.Len(t, list.DistributionRecords, 3)
}

func TestKeeper_GetRecordsForRecipient(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForRecipient(ctx, outList[0].Address)
	assert.Len(t, list.DistributionRecords, 1)
}

func TestKeeper_GetRecordsForRecipient_StatusCompleted(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForRecipient(ctx, outList[0].Address)
	assert.Len(t, list.DistributionRecords, 1)

}

func TestKeeper_GetRecordsForRecipient_StatusFailed(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForRecipient(ctx, outList[0].Address)
	assert.Len(t, list.DistributionRecords, 1)

}

func TestKeeper_GetRecords(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecords(ctx)
	assert.Len(t, list.DistributionRecords, 3)

}

func TestKeeper_GetRecords_StatusCompleted(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecords(ctx)
	assert.Len(t, list.DistributionRecords, 3)

}

func TestKeeper_GetRecords_StatusFailed(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecords(ctx)
	assert.Len(t, list.DistributionRecords, 3)

}
func TestKeeper_GetRecords_fromName(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForName(ctx, name)
	assert.Len(t, list.DistributionRecords, 3)
}

func TestKeeper_GetRecords_fromName_StatusCompleted(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForName(ctx, name)
	assert.Len(t, list.DistributionRecords, 3)
}
func TestKeeper_GetRecords_fromName_StatusFailed(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(3, "1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address, types.DistributionStatus_DISTRIBUTION_STATUS_FAILED, record.DistributionType)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForName(ctx, name)
	assert.Len(t, list.DistributionRecords, 3)
}

func TestKeeper_ChangeRecordStatus_ErrorInSettingRecord(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	dr := types.DistributionRecord{
		DistributionStatus:          types.DistributionStatus_DISTRIBUTION_STATUS_FAILED,
		DistributionType:            types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED,
		DistributionName:            "",
		RecipientAddress:            "",
		Coins:                       sdk.Coins{},
		DistributionStartHeight:     int64(0),
		DistributionCompletedHeight: int64(0),
		AuthorizedRunner:            types.AttributeKeyDistributionRunner,
	}
	height := int64(1)
	newstatus := types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED

	err := keeper.ChangeRecordStatus(ctx, dr, height, newstatus)
	assert.Error(t, err)
}
