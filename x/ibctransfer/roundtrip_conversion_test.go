package ibctransfer_test

import (
	"testing"

	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"

	"github.com/Sifchain/sifnode/x/ibctransfer"
	"github.com/Sifchain/sifnode/x/ibctransfer/keeper"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func TestExportImportConversionEquality(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	maxUInt64 := uint64(18446744073709551615)
	returningTransferPacket := channeltypes.Packet{
		Sequence:           0,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               nil,
	}
	tokenPacket := transfertypes.FungibleTokenPacketData{
		Denom:  "transfer/channel-0/microrowan",
		Amount: 184467440737,
	}
	rowanEntry := tokenregistrytypes.RegistryEntry{
		Decimals:             18,
		Denom:                "rowan",
		BaseDenom:            "rowan",
		IbcCounterpartyDenom: "microrowan",
	}
	microRowanEntry := tokenregistrytypes.RegistryEntry{
		Decimals:  10,
		Denom:     "microrowan",
		BaseDenom: "microrowan",
		UnitDenom: "rowan",
	}
	app.TokenRegistryKeeper.SetToken(ctx, &rowanEntry)
	app.TokenRegistryKeeper.SetToken(ctx, &microRowanEntry)
	registry := app.TokenRegistryKeeper.GetDenomWhitelist(ctx)
	rEntry := app.TokenRegistryKeeper.GetDenom(registry, "rowan")
	require.NotNil(t, rEntry)
	mrEntry := app.TokenRegistryKeeper.GetDenom(registry, "microrowan")
	require.NotNil(t, mrEntry)
	msg := &transfertypes.MsgTransfer{Token: sdk.NewCoin("rowan", sdk.NewIntFromUint64(maxUInt64))}
	outgoingDeduction, outgoingAddition := keeper.ConvertCoinsForTransfer(msg, rEntry, mrEntry)
	incomingDeduction, incomingAddition := ibctransfer.GetConvForIncomingCoins(ctx, app.TokenRegistryKeeper, returningTransferPacket, tokenPacket)
	require.Greater(t, (*incomingAddition).Amount.String(), (*incomingDeduction).Amount.String())
	require.Equal(t, outgoingDeduction, *incomingAddition)
	require.Equal(t, outgoingAddition, *incomingDeduction)
}
