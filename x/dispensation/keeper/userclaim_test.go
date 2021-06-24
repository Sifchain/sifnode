package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/stretchr/testify/assert"
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
	assert.Len(t, fetchList, numberOfClaims)
	fetchList = keeper.GetClaims(ctx)
	assert.Len(t, fetchList, numberOfClaims+numberOfClaims)
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
	assert.Len(t, fetchList, numberOfClaims)
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
	assert.Len(t, fetchList, numberOfClaims)
	for _, claim := range claimList {
		keeper.DeleteClaim(ctx, claim.UserAddress, claim.UserClaimType)
	}
	fetchList = keeper.GetClaims(ctx)
	assert.Len(t, fetchList, 0)
}

func TestKeeper_LockClaim(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	numberOfClaims := 1
	claimList := test.CreateClaimsList(numberOfClaims, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	fetchList := keeper.GetClaims(ctx)
	assert.Len(t, fetchList, numberOfClaims)
	for _, claim := range claimList {
		err := keeper.LockClaim(ctx, claim.UserAddress, claim.UserClaimType)
		assert.NoError(t, err)
	}
	lockedList := keeper.GetClaims(ctx)
	for _, claim := range lockedList {
		assert.True(t, claim.Locked)
	}
}
