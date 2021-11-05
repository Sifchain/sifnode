package ibctransfer

// OnTimeoutMaybeConvert potentially does a conversion from the denom that was sent out,
// back to the original unit_denom.
// Todo : Check with @Ivan why this file is name onrecv .

//func OnTimeoutMaybeConvert(
//	ctx sdk.Context,
//	sdkTransferKeeper sctransfertypes.SDKTransferKeeper,
//	whitelistKeeper tokenregistrytypes.Keeper,
//	bankKeeper transfertypes.BankKeeper,
//	packet channeltypes.Packet,
//	relayer sdk.AccAddress,
//) error {
//	var data transfertypes.FungibleTokenPacketData
//	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
//		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
//	}
//	// refund tokens
//	if err := sdkTransferKeeper.OnTimeoutPacket(ctx, packet, data); err != nil {
//		return err
//	}
//	ctx.EventManager().EmitEvent(
//		sdk.NewEvent(
//			transfertypes.EventTypeTimeout,
//			sdk.NewAttribute(sdk.AttributeKeyModule, transfertypes.ModuleName),
//			sdk.NewAttribute(transfertypes.AttributeKeyRefundReceiver, data.Sender),
//			sdk.NewAttribute(transfertypes.AttributeKeyRefundDenom, data.Denom),
//			sdk.NewAttribute(transfertypes.AttributeKeyRefundAmount, fmt.Sprintf("%d", data.Amount)),
//		),
//	)
//	denom := data.Denom
//	registry := whitelistKeeper.GetRegistry(ctx)
//	denomEntry, err := whitelistKeeper.GetEntry(registry, denom)
//	if err != nil {
//		ctx.Logger().Error(err.Error())
//	} else if denomEntry.Decimals > 0 && denomEntry.UnitDenom != "" {
//		convertToDenomEntry, err := whitelistKeeper.GetEntry(registry, denomEntry.UnitDenom)
//		if err == nil && convertToDenomEntry.Decimals > denomEntry.Decimals {
//			err := helpers.ExecConvForRefundCoins(ctx, bankKeeper, denomEntry, convertToDenomEntry, packet, data)
//			if err != nil {
//				return  err
//			}
//			return nil
//		}
//	}
//	return nil
//}
