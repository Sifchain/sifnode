package clp_test

import (
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndBlocker(t *testing.T) {
	ctx, app := test.CreateTestAppClp(false)
	_ = test.GeneratePoolsFromFile(app.ClpKeeper, ctx)
	SetRewardParams(app.ClpKeeper, ctx)

	_ = clp.EndBlocker(ctx, app.ClpKeeper)

	pooldash, err := app.ClpKeeper.GetPool(ctx, "cdash")
	assert.NoError(t, err)
	poolceth, err := app.ClpKeeper.GetPool(ctx, "ceth")
	assert.NoError(t, err)
	assert.True(t, poolceth.NativeAssetBalance.GT(pooldash.NativeAssetBalance))

}

func SetRewardParams(keeper keeper.Keeper, ctx sdk.Context) {
	multiplierDec1 := sdk.MustNewDecFromStr("0.5")
	multiplierDec2 := sdk.MustNewDecFromStr("1.5")
	allocations := sdk.NewUintFromString("2000000000000000000")
	keeper.SetRewardParams(ctx, &types.RewardParams{
		LiquidityRemovalLockPeriod:   0,
		LiquidityRemovalCancelPeriod: 2,
		RewardPeriods: []*types.RewardPeriod{{
			Id:         "1",
			StartBlock: 0,
			EndBlock:   2,
			Allocation: &allocations,
			Multipliers: []*types.PoolMultiplier{{
				Asset:      "cdash",
				Multiplier: &multiplierDec1,
			},
				{
					Asset:      "ceth",
					Multiplier: &multiplierDec2,
				},
			},
		}},
	})
}
