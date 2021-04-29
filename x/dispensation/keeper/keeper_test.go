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
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String())
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForNameAll(ctx, name)
	assert.Len(t, list, 3)
}

func TestKeeper_GetRecordsForRecipient(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.GenerateOutputList("1000000000")
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1)
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		_, err = keeper.GetDistributionRecord(ctx, name, rec.Address.String())
		assert.NoError(t, err)
	}
	list := keeper.GetRecordsForRecipient(ctx, outList[0].Address)
	assert.Len(t, list, 1)
}

func TestKeeper_GetClaimsByType(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	numberOfClaims := 1000
	claimList := test.CreateClaimsList(numberOfClaims, types.ValidatorSubsidy)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	claimList = test.CreateClaimsList(numberOfClaims, types.LiquidityMining)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	fetchList := keeper.GetClaimsByType(ctx, types.ValidatorSubsidy)
	assert.Len(t, fetchList, numberOfClaims)
	fetchList = keeper.GetClaims(ctx)
	assert.Len(t, fetchList, numberOfClaims+numberOfClaims)
}

func TestKeeper_GetClaims(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	numberOfClaims := 1000
	claimList := test.CreateClaimsList(numberOfClaims, types.ValidatorSubsidy)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
		assert.True(t, keeper.ExistsClaim(ctx, claim.UserAddress.String(), claim.UserClaimType))
	}
	fetchList := keeper.GetClaims(ctx)
	assert.Len(t, fetchList, numberOfClaims)
}

func TestKeeper_DeleteClaim(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	numberOfClaims := 1000
	claimList := test.CreateClaimsList(numberOfClaims, types.ValidatorSubsidy)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	fetchList := keeper.GetClaims(ctx)
	assert.Len(t, fetchList, numberOfClaims)
	for _, claim := range claimList {
		keeper.DeleteClaim(ctx, claim.UserAddress.String(), claim.UserClaimType)
	}
	fetchList = keeper.GetClaims(ctx)
	assert.Len(t, fetchList, 0)

}
