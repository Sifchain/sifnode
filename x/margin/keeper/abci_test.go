package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/Sifchain/sifnode/x/margin/types"
)

func TestKeeper_BeginBlocker(t *testing.T) {
	t.Skip()
	t.Run("default", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper

		params := types.Params{
			LeverageMax:          sdk.NewDec(10),
			InterestRateMax:      sdk.NewDec(5),
			InterestRateMin:      sdk.NewDec(1),
			InterestRateIncrease: sdk.NewDec(1),
			InterestRateDecrease: sdk.NewDec(1),
			HealthGainFactor:     sdk.NewDec(1),
			EpochLength:          1,
		}
		want := types.GenesisState{Params: &params}
		marginKeeper.InitGenesis(ctx, want)

		marginKeeper.BeginBlocker(ctx)
	})
}
