package dispensation_test

import (
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExportGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(1000, "1000000000")
	claimList := test.CreateClaimsList(10000, types.ValidatorSubsidy)
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.Airdrop, name, sdk.AccAddress{}))
		assert.NoError(t, err)
	}
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)
	assert.Len(t, genState.DistributionRecords, 1000)
	assert.Len(t, genState.Distributions, 1)
	assert.Len(t, genState.Claims, 10000)
}

func TestInitGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	app2, ctx2 := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	keeper2 := app2.DispensationKeeper
	outList := test.CreatOutputList(1000, "1000000000")
	claimList := test.CreateClaimsList(10000, types.ValidatorSubsidy)
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.Airdrop, name, sdk.AccAddress{}))
		assert.NoError(t, err)
	}
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)
	assert.Len(t, keeper2.GetDistributions(ctx2), 0)
	assert.Len(t, keeper2.GetRecords(ctx2), 0)
	assert.Len(t, keeper2.GetClaims(ctx2), 0)
	dispensation.InitGenesis(ctx2, keeper2, genState)
	assert.Len(t, keeper2.GetDistributions(ctx2), 1)
	assert.Len(t, keeper2.GetRecords(ctx2), 1000)
	assert.Len(t, keeper2.GetClaims(ctx2), 10000)
}

func TestValidateGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outList := test.CreatOutputList(1000, "1000000000")
	claimList := test.CreateClaimsList(10000, types.ValidatorSubsidy)
	name := uuid.New().String()
	for _, rec := range outList {
		record := types.NewDistributionRecord(name, types.Airdrop, rec.Address, rec.Coins, ctx.BlockHeight(), -1, sdk.AccAddress{})
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.Airdrop, name, sdk.AccAddress{}))
		assert.NoError(t, err)
	}
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)
	assert.Len(t, genState.DistributionRecords, 1000)
	assert.Len(t, genState.Distributions, 1)
	assert.Len(t, genState.Claims, 10000)
	assert.NoError(t, dispensation.ValidateGenesis(genState))
}
