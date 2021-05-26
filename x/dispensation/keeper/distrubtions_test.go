package keeper_test

//func TestKeeper_GetDistributions(t *testing.T) {
//	app, ctx := test.CreateTestApp(false)
//	keeper := app.DispensationKeeper
//	for i := 0; i < 10; i++ {
//		name := uuid.New().String()
//		distribution := types.NewDistribution(types.Airdrop, name)
//		err := keeper.SetDistribution(ctx, distribution)
//		assert.NoError(t, err)
//		res, err := keeper.GetDistribution(ctx, name, types.Airdrop)
//		assert.NoError(t, err)
//		assert.Equal(t, res.String(), distribution.String())
//	}
//	list := keeper.GetDistributions(ctx)
//	assert.Len(t, list, 10)
//}
