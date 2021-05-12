package keeper_test

//func TestKeeper_AccumulateDrops(t *testing.T) {
//	app, ctx := test.CreateTestApp(false)
//	keeper := app.DispensationKeeper
//	inputList := test.GenerateInputList("15000000000000000000")
//	//outputList := test.GenerateOutputList("10000000000000000000")
//
//	for _, in := range inputList {
//		_, err := keeper.GetBankKeeper().AddCoins(ctx, in.Address, in.Coins)
//		assert.NoError(t, err)
//	}
//	err := keeper.AccumulateDrops(ctx, inputList)
//	assert.NoError(t, err)
//	moduleBalance, _ := sdk.NewIntFromString("30000000000000000000")
//	assert.True(t, keeper.GetBankKeeper().HasCoins(ctx, types.GetDistributionModuleAddress(), sdk.Coins{sdk.NewCoin("rowan", moduleBalance)}))
//
//}
//
//func TestKeeper_CreateAndDistributeDrops(t *testing.T) {
//	app, ctx := test.CreateTestApp(false)
//	keeper := app.DispensationKeeper
//	inputList := test.GenerateInputList("15000000000000000000")
//	outputList := test.GenerateOutputList("10000000000000000000")
//	for _, in := range inputList {
//		_, err := keeper.GetBankKeeper().AddCoins(ctx, in.Address, in.Coins)
//		assert.NoError(t, err)
//	}
//	err := keeper.AccumulateDrops(ctx, inputList)
//	assert.NoError(t, err)
//	moduleBalance, _ := sdk.NewIntFromString("30000000000000000000")
//	assert.True(t, keeper.GetBankKeeper().HasCoins(ctx, types.GetDistributionModuleAddress(), sdk.Coins{sdk.NewCoin("rowan", moduleBalance)}))
//
//	err = keeper.CreateDrops(ctx, outputList, "ar1")
//	assert.NoError(t, err)
//}
//
//func TestKeeper_VerifyDistribution(t *testing.T) {
//	app, ctx := test.CreateTestApp(false)
//	keeper := app.DispensationKeeper
//	err := keeper.VerifyDistribution(ctx, "AR1", types.Airdrop)
//	assert.NoError(t, err)
//	err = keeper.VerifyDistribution(ctx, "AR1", types.Airdrop)
//	assert.Error(t, err)
//}
