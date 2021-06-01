package dispensation_test

//func TestNewHandler_CreateDistribution(t *testing.T) {
//	app, ctx := test.CreateTestApp(false)
//	keeper := app.DispensationKeeper
//	handler := dispensation.NewHandler(keeper)
//	recipients := 3000
//	inputList := test.CreatInputList(2, "15000000000000000000000")
//	outputList := test.CreatOutputList(recipients, "10000000000000000000")
//	err := bank.ValidateInputsOutputs(inputList, outputList)
//	assert.NoError(t, err)
//	for _, in := range inputList {
//		_, err := keeper.GetBankKeeper().AddCoins(ctx, in.Address, in.Coins)
//		assert.NoError(t, err)
//	}
//	msgAirdrop := types.NewMsgDistribution(sdk.AccAddress{}, "AR1", types.Airdrop, inputList, outputList)
//	res, err := handler(ctx, msgAirdrop)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	dr := keeper.GetRecordsForNameAll(ctx, "AR1")
//	assert.Len(t, dr, recipients)
//}

//func TestNewHandler_CreateClaim(t *testing.T) {
//	app, ctx := test.CreateTestApp(false)
//	keeper := app.DispensationKeeper
//	handler := dispensation.NewHandler(keeper)
//	address := sdk.AccAddress(crypto.AddressHash([]byte("User1")))
//	msgClaim := types.NewMsgCreateClaim(address, types.ValidatorSubsidy)
//	res, err := handler(ctx, msgClaim)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//
//	cl, err := keeper.GetClaim(ctx, address.String(), types.ValidatorSubsidy)
//	require.NoError(t, err)
//	assert.False(t, cl.Locked)
//	assert.Equal(t, cl.UserAddress.String(), address.String())
//}
