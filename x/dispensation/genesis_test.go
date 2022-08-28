package dispensation_test

import (
	"testing"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
)

func TestExportGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper

	mintController := types.MintController{
		TotalCounter: sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.NewInt(1234)),
	}
	keeper.SetMintController(ctx, mintController)
	outList := test.CreatOutputList(1000, "1000000000")
	claimList := test.CreateClaimsList(10000, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	name := uuid.New().String()
	authorizedRunner := sdk.AccAddress(crypto.AddressHash([]byte("Runner")))
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, authorizedRunner.String()))
		assert.NoError(t, err)
	}
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)
	assert.Equal(t, *genState.MintController, mintController)
	assert.Len(t, genState.DistributionRecords.DistributionRecords, 1000)
	assert.Len(t, genState.Distributions.Distributions, 1)
	assert.Len(t, genState.Claims.UserClaims, 10000)
}

func TestInitGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	app2, ctx2 := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	keeper2 := app2.DispensationKeeper

	zeroMintController := types.MintController{
		TotalCounter: sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.NewInt(0)),
	}

	mintController := types.MintController{
		TotalCounter: sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.NewInt(1234)),
	}
	keeper.SetMintController(ctx, mintController)

	outList := test.CreatOutputList(1000, "1000000000")
	claimList := test.CreateClaimsList(10000, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	name := uuid.New().String()
	authorizedRunner := sdk.AccAddress(crypto.AddressHash([]byte("Runner")))
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, "")
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, authorizedRunner.String()))
		assert.NoError(t, err)
	}
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)

	exportedMintController, exists := keeper2.GetMintController(ctx2)
	// MintController's amount set as zero when
	assert.Equal(t, exists, true)
	assert.Equal(t, zeroMintController, exportedMintController)
	assert.Len(t, keeper2.GetDistributions(ctx2).Distributions, 0)
	assert.Len(t, keeper2.GetRecords(ctx2).DistributionRecords, 0)
	assert.Len(t, keeper2.GetClaims(ctx2).UserClaims, 0)

	dispensation.InitGenesis(ctx2, keeper2, genState)
	exportedMintController, exists = keeper2.GetMintController(ctx2)
	assert.Equal(t, exists, true)
	assert.Equal(t, *genState.MintController, exportedMintController)

	assert.Equal(t, *genState.MintController, mintController)
	assert.Len(t, keeper2.GetDistributions(ctx2).Distributions, 1)
	assert.Len(t, keeper2.GetRecords(ctx2).DistributionRecords, 1000)
	assert.Len(t, keeper2.GetClaims(ctx2).UserClaims, 10000)
}

func TestValidateGenesis(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	mintController := types.MintController{
		TotalCounter: sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, sdk.NewInt(1234)),
	}
	keeper.SetMintController(ctx, mintController)
	outList := test.CreatOutputList(1000, "1000000000")
	claimList := test.CreateClaimsList(10000, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	name := uuid.New().String()
	authorizedRunner := sdk.AccAddress(crypto.AddressHash([]byte("Runner")))
	for _, rec := range outList {
		record := types.NewDistributionRecord(types.DistributionStatus_DISTRIBUTION_STATUS_PENDING, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, rec.Address, rec.Coins, ctx.BlockHeight(), -1, authorizedRunner.String())
		err := keeper.SetDistributionRecord(ctx, record)
		assert.NoError(t, err)
		err = keeper.SetDistribution(ctx, types.NewDistribution(types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, name, authorizedRunner.String()))
		assert.NoError(t, err)
	}
	for _, claim := range claimList {
		err := keeper.SetClaim(ctx, claim)
		assert.NoError(t, err)
	}
	genState := dispensation.ExportGenesis(ctx, keeper)
	assert.Equal(t, *genState.MintController, mintController)
	assert.Len(t, genState.DistributionRecords.DistributionRecords, 1000)
	assert.Len(t, genState.Distributions.Distributions, 1)
	assert.Len(t, genState.Claims.UserClaims, 10000)
	assert.NoError(t, dispensation.ValidateGenesis(genState))
}
