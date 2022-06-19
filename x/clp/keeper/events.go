package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func CreateEventMsg(signer string) sdk.Event {
	return sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, signer))
}

func CreateEventBlockHeight(ctx sdk.Context, eventType string, attribute sdk.Attribute) sdk.Event {
	return sdk.NewEvent(
		eventType,
		attribute,
		sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	)
}
