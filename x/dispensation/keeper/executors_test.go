package keeper_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
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
	outputAmount := "10000000000000000000"
	inputList := test.GenerateInputList("90000000000000000000")
	outputList := test.GenerateOutputList(outputAmount)
	for _, in := range inputList {
		_, err := keeper.GetBankKeeper().AddCoins(ctx, in.Address, in.Coins)
		assert.NoError(t, err)
	}
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	totalCoins = totalCoins.Add(totalCoins...).Add(totalCoins...)
	err = keeper.AccumulateDrops(ctx, inputList[0].Address, totalCoins)
	assert.NoError(t, err)
	moduleBalance, _ := sdk.NewIntFromString("15000000000000000000")
	assert.True(t, keeper.GetBankKeeper().HasCoins(ctx, types.GetDistributionModuleAddress(), sdk.Coins{sdk.NewCoin("rowan", moduleBalance)}))
	distributionName := "ar1"
	runner := sdk.AccAddress{}
	err = keeper.CreateDrops(ctx, outputList, distributionName, types.Airdrop, runner)
	assert.NoError(t, err)
	err = keeper.CreateDrops(ctx, outputList, distributionName, types.LiquidityMining, runner)
	assert.NoError(t, err)
	err = keeper.CreateDrops(ctx, outputList, distributionName, types.LiquidityMining, runner)
	assert.NoError(t, err)

	_, err = keeper.DistributeDrops(ctx, 1, distributionName, runner, types.Airdrop)
	assert.NoError(t, err)
	_, err = keeper.DistributeDrops(ctx, 1, distributionName, runner, types.LiquidityMining)
	assert.NoError(t, err)
	completedRecords := keeper.GetRecordsForNameCompleted(ctx, distributionName)
	assert.Equal(t, 6, len(completedRecords))
	recordsLM := keeper.GetRecordsForNameAndType(ctx, distributionName, types.LiquidityMining)
	assert.Equal(t, 3, len(recordsLM))
	recordsAD := keeper.GetRecordsForNameAndType(ctx, distributionName, types.Airdrop)
	assert.Equal(t, 3, len(recordsAD))
	doubleOutputAmount, ok := sdk.NewIntFromString("20000000000000000000")
	assert.True(t, ok)
	for i := 0; i < len(outputList); i++ {
		assert.True(t, recordsLM[i].Coins.IsAllGT(recordsAD[i].Coins))
		assert.True(t, recordsLM[i].Coins.AmountOf("rowan").Equal(doubleOutputAmount))
	}
}

func TestKeeper_CreateAndDistributeDrops_AddressError(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	outputAmount := "10000000000000000000"

	outputList := test.GenerateOutputList(outputAmount)

	newCoins, _ := sdk.NewIntFromString("900000")
	newAddressCoins := sdk.NewCoins(sdk.NewCoin("rowan", newCoins))
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)

	// New address added after total Coins is calculated
	// This means total coins is less than that the amount required for distribution
	outputList = append(outputList, bank.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		newAddressCoins))
	outputList = append(outputList, bank.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		newAddressCoins))
	fundingAddress := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	_, err = keeper.GetBankKeeper().AddCoins(ctx, fundingAddress, totalCoins)
	assert.NoError(t, err)
	err = keeper.AccumulateDrops(ctx, fundingAddress, totalCoins)
	assert.NoError(t, err)
	distributionName := "ar1"
	runner := sdk.AccAddress{}
	err = keeper.CreateDrops(ctx, outputList, distributionName, types.LiquidityMining, runner)
	assert.NoError(t, err)

	_, err = keeper.DistributeDrops(ctx, 1, distributionName, runner, types.LiquidityMining)
	assert.NoError(t, err)
	completedRecords := keeper.GetRecordsForNameCompleted(ctx, distributionName)
	assert.GreaterOrEqual(t, len(completedRecords), 3)
	failedRecords := keeper.GetRecordsForNameAllFailed(ctx, distributionName)
	assert.LessOrEqual(t, len(failedRecords), 2)
}

func TestKeeper_VerifyDistribution(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	err := keeper.VerifyAndSetDistribution(ctx, "AR1", types.Airdrop, sdk.AccAddress{})
	assert.NoError(t, err)
	err = keeper.VerifyAndSetDistribution(ctx, "AR1", types.Airdrop, sdk.AccAddress{})
	assert.Error(t, err)
}
