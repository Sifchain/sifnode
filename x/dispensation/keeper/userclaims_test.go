package keeper_test

import (
<<<<<<< HEAD
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/stretchr/testify/assert"
=======
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/stretchr/testify/assert"
	"testing"
>>>>>>> develop
)

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
