package keeper_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_GetRecordsForName(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.GenerateOutputList("1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
		record.DistributionStatus = types.Pending
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), record.DistributionType.String())
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForNameAll(ctx, name)
	assert.Len(t, list, 3)
}

func TestKeeper_GetRecordsForNames(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.GenerateOutputList("1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.LiquidityMining, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
		record.DistributionStatus = types.Pending
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), record.DistributionType.String())
		assert.NoError(t, err)
	}
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.ValidatorSubsidy, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
		record.DistributionStatus = types.Pending
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), record.DistributionType.String())
		assert.NoError(t, err)
	}
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.ValidatorSubsidy, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
		record.DistributionStatus = types.Pending
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), record.DistributionType.String())
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForNameAll(ctx, name)
	assert.Len(t, list, 6)
	list = keeper.GetRecordsForNameAndType(ctx, name, types.LiquidityMining)
	assert.Len(t, list, 3)
	list = keeper.GetRecordsForNameAndType(ctx, name, types.ValidatorSubsidy)
	assert.Len(t, list, 3)
}

func TestKeeper_GetRecordsForRecipient(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.GenerateOutputList("1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
		record.DistributionStatus = types.Pending
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String(), record.DistributionType.String())
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForRecipient(ctx, outList[0].Address)
	assert.Len(t, list, 1)
}
