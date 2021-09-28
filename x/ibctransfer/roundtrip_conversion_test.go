package ibctransfer_test

import (
	"testing"

	tokenregistrytest "github.com/Sifchain/sifnode/x/tokenregistry/test"

	"github.com/Sifchain/sifnode/x/ibctransfer/helpers"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
)

func TestExportImportConversionEquality(t *testing.T) {
	app, ctx, _ := tokenregistrytest.CreateTestApp(false)
	maxUInt64 := uint64(18446744073709551615)
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
	outgoingDeduction, outgoingAddition := helpers.ConvertCoinsForTransfer(msg, rEntry, mrEntry)
	mrEntryUnit := app.TokenRegistryKeeper.GetDenom(registry, mrEntry.UnitDenom)
	require.NotNil(t, mrEntryUnit)
	diff := uint64(mrEntryUnit.Decimals - mrEntry.Decimals)
	convAmount := helpers.ConvertIncomingCoins(ctx, app.TokenRegistryKeeper, 184467440737, diff)
	incomingDeduction := sdk.NewCoin("microrowan", sdk.NewIntFromUint64(184467440737))
	incomingAddition := sdk.NewCoin("rowan", convAmount)
	require.Greater(t, incomingAddition.Amount.String(), incomingDeduction.Amount.String())
	require.Equal(t, outgoingDeduction, incomingAddition)
	require.Equal(t, outgoingAddition, incomingDeduction)
}
