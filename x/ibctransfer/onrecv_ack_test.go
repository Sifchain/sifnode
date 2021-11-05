package ibctransfer_test

// TODO : Check if we need this . Refunds should now be controlled by the SDK refund
//func TestExecConvForRefundCoins(t *testing.T) {
//	app, ctx, _ := test.CreateTestApp(false)
//	addrs, _ := test2.CreateTestAddrs(2)
//	packet := channeltypes.Packet{
//		SourcePort:         "transfer",
//		SourceChannel:      "channel-0",
//		DestinationPort:    "transfer",
//		DestinationChannel: "channel-1",
//	}
//	returningData := types.FungibleTokenPacketData{
//		Denom:  "transfer/channel-0/ueth",
//		Sender: addrs[0].String(),
//	}
//	nonReturningData := types.FungibleTokenPacketData{
//		Denom:  "transfer/channel-1/ueth",
//		Sender: addrs[0].String(),
//	}
//	ibcRegistryEntry := tokenregistrytypes.RegistryEntry{
//		Denom:     "ueth",
//		Decimals:  10,
//		UnitDenom: "ceth",
//	}
//	ibcRegistryEntry2 := tokenregistrytypes.RegistryEntry{
//		Denom:       "ibc/C1061B25E69D71E96BED65B5652168F41927316D07D6B417A3A9774F94A4CB7A",
//		Decimals:    10,
//		UnitDenom:   "ceth",
//		Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_IBCIMPORT},
//	}
//	unitDenomEntry := tokenregistrytypes.RegistryEntry{
//		Denom:    "ceth",
//		Decimals: 18,
//	}
//	app.TokenRegistryKeeper.SetToken(ctx, &unitDenomEntry)
//	app.TokenRegistryKeeper.SetToken(ctx, &ibcRegistryEntry)
//	app.TokenRegistryKeeper.SetToken(ctx, &ibcRegistryEntry2)
//	mintedDenom := helpers.GetMintedDenomFromPacket(packet, returningData)
//	registry := app.TokenRegistryKeeper.GetRegistry(ctx)
//	mintedDenomEntry, err := app.TokenRegistryKeeper.GetEntry(registry, mintedDenom)
//	require.NoError(t, err)
//	allowed := helpers.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, packet, returningData, mintedDenomEntry)
//	require.Equal(t, allowed, true)
//	convertToDenomEntry, err := app.TokenRegistryKeeper.GetEntry(registry, mintedDenomEntry.UnitDenom)
//	require.NoError(t, err)
//	err = helpers.ExecConvForRefundCoins(ctx, app.BankKeeper, mintedDenomEntry, convertToDenomEntry, packet, returningData)
//	require.NoError(t, err)
//	mintedDenom = helpers.GetMintedDenomFromPacket(packet, nonReturningData)
//	mintedDenomEntry, err = app.TokenRegistryKeeper.GetEntry(registry, mintedDenom)
//	require.NoError(t, err)
//	allowed = helpers.IsRecvPacketAllowed(ctx, app.TokenRegistryKeeper, packet, nonReturningData, mintedDenomEntry)
//	require.Equal(t, allowed, true)
//	convertToDenomEntry, err = app.TokenRegistryKeeper.GetEntry(registry, mintedDenomEntry.UnitDenom)
//	require.NoError(t, err)
//	err = helpers.ExecConvForRefundCoins(ctx, app.BankKeeper, mintedDenomEntry, convertToDenomEntry, packet, nonReturningData)
//	require.NoError(t, err)
//}
