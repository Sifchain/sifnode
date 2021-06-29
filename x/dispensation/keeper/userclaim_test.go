package keeper_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_GetClaimsByType(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	numberOfClaims := 1000
	claimList := test.CreateClaimsList(numberOfClaims, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	claimList = test.CreateClaimsList(numberOfClaims, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	fetchList := keeper.GetClaimsByType(ctx, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	assert.Len(t, fetchList.UserClaims, numberOfClaims)
	fetchList = keeper.GetClaims(ctx)
	assert.Len(t, fetchList.UserClaims, numberOfClaims+numberOfClaims)
}

func TestKeeper_GetClaims(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	numberOfClaims := 1000
	claimList := test.CreateClaimsList(numberOfClaims, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
		assert.True(t, keeper.ExistsClaim(ctx, claim.UserAddress, claim.UserClaimType))
	}
	fetchList := keeper.GetClaims(ctx)
	assert.Len(t, fetchList.UserClaims, numberOfClaims)
}

func TestKeeper_DeleteClaim(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	numberOfClaims := 1000
	claimList := test.CreateClaimsList(numberOfClaims, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	fetchList := keeper.GetClaims(ctx)
	assert.Len(t, fetchList.UserClaims, numberOfClaims)
	for _, claim := range claimList {
		keeper.DeleteClaim(ctx, claim.UserAddress, claim.UserClaimType)
	}
	fetchList = keeper.GetClaims(ctx)
	assert.Len(t, fetchList.UserClaims, 0)
}
