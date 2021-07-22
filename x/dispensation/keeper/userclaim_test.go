package keeper_test

import (
	"testing"
	"time"

	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
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
	assert.Len(t, claimList, numberOfClaims)
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
	assert.Len(t, claimList, numberOfClaims)
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

func TestKeeper_FailGetUserClaimIfNotCreated(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	claimList := keeper.GetClaims(ctx)
	assert.Len(t, claimList.UserClaims, 0)
	address := sdk.AccAddress(crypto.AddressHash([]byte("User1")))
	_, err := keeper.GetClaim(ctx, address.String(), types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	assert.Error(t, err)
}

func TestKeeper_GetClaim(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	address := sdk.AccAddress(crypto.AddressHash([]byte("User1")))
	claimType := types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY
	claim, err := types.NewUserClaim(address.String(), claimType, time.Now())
	assert.NoError(t, err)
	assert.Equal(t, claim.UserClaimType, claimType)
	assert.Equal(t, claim.UserAddress, address.String())
	_, err = keeper.GetClaim(ctx, claim.UserAddress, claimType)
	assert.Error(t, err)
	err = keeper.SetClaim(ctx, claim)
	assert.NoError(t, err)
	cl, err := keeper.GetClaim(ctx, claim.UserAddress, claimType)
	assert.NoError(t, err)
	assert.Equal(t, cl.UserClaimType, claimType)
	assert.Equal(t, cl.UserAddress, claim.UserAddress)
}
