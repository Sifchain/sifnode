package dispensation_test

import (
	types2 "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_BeginBlocker(t *testing.T) {

	// Basic Setup
	app, ctx := test.CreateTestApp(false)
	ecoPoolAddress, err := sdk.AccAddressFromBech32(types.EcoPool)
	if err != nil {
		ctx.Logger().Error("Unable to get address")
		return
	}
	rowan := types2.GetSettlementAsset().Symbol
	expectedMintAmount, ok := sdk.NewIntFromString(types.MintAmountPerBlock)
	assert.True(t, ok)
	totalMintAmount, ok := sdk.NewIntFromString(types.MaxMintAmount)
	assert.True(t, ok)
	expectedBlocks := totalMintAmount.Quo(expectedMintAmount).Int64()
	if !totalMintAmount.Mod(expectedMintAmount).IsZero() {
		expectedBlocks++
	}
	ogMintAmount := expectedMintAmount

	// Starting Balance of Ecopool is 0
	assert.Equal(t, sdk.ZeroInt(), app.DispensationKeeper.GetBankKeeper().GetBalance(ctx, ecoPoolAddress, rowan).Amount)

	// Verify starting counter
	controller, found := app.DispensationKeeper.GetMintController(ctx)
	assert.True(t, found)
	require.Equal(t, sdk.ZeroInt(), controller.TotalCounter.Amount)

	// Simulate Blocks
	for i := int64(0); i < expectedBlocks; i++ {
		dispensation.BeginBlocker(ctx, app.DispensationKeeper)
		// Asserting for non-last block
		if i < expectedBlocks-1 {
			assert.Equal(t, expectedMintAmount, app.DispensationKeeper.GetBankKeeper().GetBalance(ctx, ecoPoolAddress, rowan).Amount)
			// Assertion for last block
		} else {
			assert.Equal(t, totalMintAmount, app.DispensationKeeper.GetBankKeeper().GetBalance(ctx, ecoPoolAddress, rowan).Amount)
		}
		expectedMintAmount = expectedMintAmount.Add(ogMintAmount)
	}
	// Asserting Token Mint Condition
	assert.False(t, app.DispensationKeeper.TokensCanBeMinted(ctx))

	// Checking BeginBlocker After failed Conditional
	dispensation.BeginBlocker(ctx, app.DispensationKeeper)
	assert.Equal(t, app.DispensationKeeper.GetBankKeeper().GetBalance(ctx, ecoPoolAddress, rowan).Amount, totalMintAmount)

}
