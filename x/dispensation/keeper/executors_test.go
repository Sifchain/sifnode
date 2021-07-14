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
	inputList := test.CreateInputList(3, rowanAmount)
	//outputList := test.GenerateOutputList("10000000000000000000")
	for _, in := range inputList {
		address, err := sdk.AccAddressFromBech32(in.Address)
		assert.NoError(t, err)
		err = keeper.GetBankKeeper().AddCoins(ctx, address, in.Coins)
		assert.NoError(t, err)
	}
	distributor := inputList[0]
	err := keeper.AccumulateDrops(ctx, distributor.Address, distributor.Coins)
	assert.NoError(t, err)
	moduleBalance, _ := sdk.NewIntFromString(rowanAmount)
	assert.True(t, keeper.HasCoins(ctx, types.GetDistributionModuleAddress(), sdk.Coins{sdk.NewCoin("rowan", moduleBalance)}))

}

func TestKeeper_CreateAndDistributeDrops(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outputAmount := "10000000000000000000"
	inputList := test.CreateInputList(3, "90000000000000000000")
	outputList := test.CreatOutputList(3, outputAmount)
	for _, in := range inputList {
		address, err := sdk.AccAddressFromBech32(in.Address)
		assert.NoError(t, err)
		err = keeper.GetBankKeeper().AddCoins(ctx, address, in.Coins)
		assert.NoError(t, err)
	}
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	totalCoins = totalCoins.Add(totalCoins...).Add(totalCoins...)
	err = keeper.AccumulateDrops(ctx, inputList[0].Address, totalCoins)
	assert.NoError(t, err)
	moduleBalance, _ := sdk.NewIntFromString("15000000000000000000")
	assert.True(t, keeper.HasCoins(ctx, types.GetDistributionModuleAddress(), sdk.Coins{sdk.NewCoin("rowan", moduleBalance)}))
	distributionName := "ar1"
	runner := ""
	err = keeper.CreateDrops(ctx, outputList, distributionName, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, runner)
	assert.NoError(t, err)
	err = keeper.CreateDrops(ctx, outputList, distributionName, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, runner)
	assert.NoError(t, err)
	err = keeper.CreateDrops(ctx, outputList, distributionName, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, runner)
	assert.NoError(t, err)

	_, err = keeper.DistributeDrops(ctx, 1, distributionName, runner, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
	assert.NoError(t, err)
	_, err = keeper.DistributeDrops(ctx, 1, distributionName, runner, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING)
	assert.NoError(t, err)
	completedRecords := keeper.GetRecordsForNameAndStatus(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	assert.Equal(t, 6, len(completedRecords.DistributionRecords))
	recordsLM := keeper.GetRecordsForNameStatusAndType(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING)
	assert.Equal(t, 3, len(recordsLM.DistributionRecords))
	recordsAD := keeper.GetRecordsForNameStatusAndType(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
	assert.Equal(t, 3, len(recordsAD.DistributionRecords))
	doubleOutputAmount, ok := sdk.NewIntFromString("20000000000000000000")
	assert.True(t, ok)
	for i := 0; i < len(outputList); i++ {
		assert.True(t, recordsLM.DistributionRecords[i].Coins.IsAllGT(recordsAD.DistributionRecords[i].Coins))
		assert.True(t, recordsLM.DistributionRecords[i].Coins.AmountOf("rowan").Equal(doubleOutputAmount))
	}
}

func TestKeeper_VerifyDistribution(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	err := keeper.VerifyAndSetDistribution(ctx, "AR1", types.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
	assert.NoError(t, err)
	err = keeper.VerifyAndSetDistribution(ctx, "AR1", types.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
	assert.Error(t, err)
}
