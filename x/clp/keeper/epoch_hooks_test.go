package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// AfterEpochEnd
func TestAfterEpochEnd(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper

	asset := types.NewAsset("cusdc")

	// Create a RewardsBucket with a denom and amount
	initialAmount := sdk.NewInt(100)
	rewardsBucket := types.RewardsBucket{
		Denom:  asset.Symbol,
		Amount: initialAmount,
	}
	clpKeeper.SetRewardsBucket(ctx, rewardsBucket)

	// mint coins to the module account
	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin(asset.Symbol, initialAmount),
	))
	require.NoError(t, err)

	// check module balance
	moduleBalance := app.BankKeeper.GetBalance(ctx, types.GetCLPModuleAddress(), asset.Symbol)
	require.Equal(t, initialAmount, moduleBalance.Amount)

	params := clpKeeper.GetRewardsParams(ctx)

	// Create a liquidity provider
	lpAddress, err := sdk.AccAddressFromBech32("sif1azpar20ck9lpys89r8x7zc8yu0qzgvtp48ng5v")
	require.NoError(t, err)
	lp := types.NewLiquidityProvider(&asset, sdk.NewUint(1), lpAddress, ctx.BlockHeight()-int64(params.RewardsLockPeriod)-1)
	clpKeeper.SetLiquidityProvider(ctx, &lp)

	// check account balance
	preBalance := app.BankKeeper.GetBalance(ctx, lpAddress, asset.Symbol)
	require.Equal(t, sdk.ZeroInt(), preBalance.Amount)

	clpKeeper.AfterEpochEnd(ctx, params.RewardsEpochIdentifier, 1)

	// check account balance
	postBalance := app.BankKeeper.GetBalance(ctx, lpAddress, asset.Symbol)
	require.Equal(t, preBalance.Add(sdk.NewCoin(asset.Symbol, initialAmount)), postBalance)
}
