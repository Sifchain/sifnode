package keeper_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_GetDistributions(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	for i := 0; i < 10; i++ {
		name := uuid.New().String()
		distribution := types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name)
		err := keeper.SetDistribution(ctx, distribution)
		assert.NoError(t, err)
		_, err = keeper.GetDistribution(ctx, name, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
		assert.NoError(t, err)
	}
	list := keeper.GetDistributions(ctx)
	assert.Len(t, list.Distributions, 10)
}

func TestKeeper_GetRecordsForName(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.GenerateOutputList("1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, rec.Address, rec.Coins, sdk.NewInt(ctx.BlockHeight()), sdk.NewInt(-1))
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForNameAll(ctx, name)
	assert.Len(t, list.DistributionRecords, 3)
}

func TestKeeper_GetRecordsForRecipient(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.GenerateOutputList("1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, rec.Address, rec.Coins, sdk.NewInt(ctx.BlockHeight()), sdk.NewInt(-1))
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address)
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForRecipient(ctx, sdk.AccAddress(outList[0].Address))
	assert.Len(t, list.DistributionRecords, 1)
}
