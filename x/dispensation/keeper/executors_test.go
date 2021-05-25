package keeper_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_AccumulateDrops(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	rowanAmount := "15000000000000000000"
	inputList := test.GenerateInputList(rowanAmount)
	//outputList := test.GenerateOutputList("10000000000000000000")
	for _, in := range inputList {
		_, err := keeper.GetBankKeeper().AddCoins(ctx, in.Address, in.Coins)
		assert.NoError(t, err)
	}
	err := keeper.AccumulateDrops(ctx, inputList[0].Address, inputList[0].Coins)
	assert.NoError(t, err)
	moduleBalance, _ := sdk.NewIntFromString(rowanAmount)
	assert.True(t, keeper.GetBankKeeper().HasCoins(ctx, types.GetDistributionModuleAddress(), sdk.Coins{sdk.NewCoin("rowan", moduleBalance)}))

}

func TestKeeper_CreateAndDistributeDrops(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	inputList := test.GenerateInputList("30000000000000000000")
	outputList := test.GenerateOutputList("10000000000000000000")
	for _, in := range inputList {
		_, err := keeper.GetBankKeeper().AddCoins(ctx, in.Address, in.Coins)
		assert.NoError(t, err)
	}
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	err = keeper.AccumulateDrops(ctx, inputList[0].Address, totalCoins)
	assert.NoError(t, err)
	moduleBalance, _ := sdk.NewIntFromString("15000000000000000000")
	assert.True(t, keeper.GetBankKeeper().HasCoins(ctx, types.GetDistributionModuleAddress(), sdk.Coins{sdk.NewCoin("rowan", moduleBalance)}))
	distributionName := "ar1"
	err = keeper.CreateDrops(ctx, outputList, distributionName, types.Airdrop)
	assert.NoError(t, err)

	err = keeper.DistributeDrops(ctx, 1)
	assert.NoError(t, err)
	completedRecords := keeper.GetRecordsForNameCompleted(ctx, distributionName)
	assert.Equal(t, 3, len(completedRecords))
}

func TestKeeper_VerifyDistribution(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	err := keeper.VerifyAndSetDistribution(ctx, "AR1", types.Airdrop)
	assert.NoError(t, err)
	err = keeper.VerifyAndSetDistribution(ctx, "AR1", types.Airdrop)
	assert.Error(t, err)
}
