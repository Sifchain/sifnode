package keeper_test

import (
	"testing"

	keepertest "github.com/Sifchain/sifnode/testutil/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAddLiquidityToRewardsBucket(t *testing.T) {
	keeper, ctx, bankKeeper := keepertest.ClpKeeper(t)

	signer := "sif1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3sxxeku"
	amount := sdk.NewCoins(sdk.NewInt64Coin("atom", 100))
	msg := types.NewMsgAddLiquidityToRewardsBucketRequest(signer, amount)

	addr, err := sdk.AccAddressFromBech32(signer)
	require.NoError(t, err)

	// Mock expectations
	bankKeeper.EXPECT().
		HasBalance(ctx, addr, msg.Amount[0]).
		Return(true).
		Times(1)

	bankKeeper.EXPECT().
		SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(msg.Amount...)).
		Return(nil).
		Times(1)

	// Call the method
	_, err = keeper.AddLiquidityToRewardsBucket(ctx, msg.Signer, msg.Amount)
	require.NoError(t, err)

	// check if rewards bucket is created
	rewardsBucket, found := keeper.GetRewardsBucket(ctx, msg.Amount[0].Denom)
	require.True(t, found)

	// check if rewards bucket has correct amount
	require.Equal(t, msg.Amount[0].Amount, rewardsBucket.Amount)
}

func TestAddLiquidityToRewardsBucket_BalanceNotAvailable(t *testing.T) {
	keeper, ctx, bankKeeper := keepertest.ClpKeeper(t)

	signer := "sif1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3sxxeku"
	amount := sdk.NewCoins(sdk.NewInt64Coin("atom", 100))
	msg := types.NewMsgAddLiquidityToRewardsBucketRequest(signer, amount)

	addr, err := sdk.AccAddressFromBech32(signer)
	require.NoError(t, err)

	// Mock expectations for HasBalance to return false
	bankKeeper.EXPECT().
		HasBalance(ctx, addr, msg.Amount[0]).
		Return(false).
		Times(1)

	// Call the method and expect an error
	_, err = keeper.AddLiquidityToRewardsBucket(ctx, msg.Signer, msg.Amount)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrBalanceNotAvailable)

	// check if rewards bucket is created
	rewardsBucket, found := keeper.GetRewardsBucket(ctx, msg.Amount[0].Denom)
	require.False(t, found)
	require.Equal(t, types.RewardsBucket{}, rewardsBucket)
}

func TestAddLiquidityToRewardsBucket_MultipleCoins(t *testing.T) {
	keeper, ctx, bankKeeper := keepertest.ClpKeeper(t)

	signer := "sif1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3sxxeku"
	amount := sdk.NewCoins(
		sdk.NewInt64Coin("atom", 100),
		sdk.NewInt64Coin("rowan", 100),
	)
	msg := types.NewMsgAddLiquidityToRewardsBucketRequest(signer, amount)

	addr, err := sdk.AccAddressFromBech32(signer)
	require.NoError(t, err)

	// Mock expectations
	bankKeeper.EXPECT().
		HasBalance(ctx, addr, msg.Amount[0]).
		Return(true).
		Times(1)

	bankKeeper.EXPECT().
		HasBalance(ctx, addr, msg.Amount[1]).
		Return(true).
		Times(1)

	bankKeeper.EXPECT().
		SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(msg.Amount...)).
		Return(nil).
		Times(1)

	// Call the method
	_, err = keeper.AddLiquidityToRewardsBucket(ctx, msg.Signer, msg.Amount)
	require.NoError(t, err)

	// check if rewards bucket is created
	rewardsBucket, found := keeper.GetRewardsBucket(ctx, msg.Amount[0].Denom)
	require.True(t, found)

	// check if rewards bucket has correct amount and denom
	require.Equal(t, msg.Amount[0].Amount, rewardsBucket.Amount)
	require.Equal(t, msg.Amount[0].Denom, rewardsBucket.Denom)

	// check if rewards bucket is created
	rewardsBucket, found = keeper.GetRewardsBucket(ctx, msg.Amount[1].Denom)
	require.True(t, found)

	// check if rewards bucket has correct amount
	require.Equal(t, msg.Amount[1].Amount, rewardsBucket.Amount)
	require.Equal(t, msg.Amount[1].Denom, rewardsBucket.Denom)
}
