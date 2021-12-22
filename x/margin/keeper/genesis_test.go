package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/Sifchain/sifnode/x/margin/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestKeeper_ExportGenesis(t *testing.T) {
	t.Run("no prior import of genesis", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper
		require.NotNil(t, marginKeeper)
		state := marginKeeper.ExportGenesis(ctx)
		require.NotNil(t, state)
		require.Nil(t, state.Params)
	})

	t.Run("prior import of genesis then export", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper
		require.NotNil(t, marginKeeper)

		params := types.Params{
			LeverageMax:          sdk.NewUint(10),
			InterestRateMax:      sdk.NewDec(5),
			InterestRateMin:      sdk.NewDec(1),
			InterestRateIncrease: sdk.NewDec(1),
			InterestRateDecrease: sdk.NewDec(1),
			HealthGainFactor:     sdk.NewDec(1),
			EpochLength:          1,
		}
		want := types.GenesisState{Params: &params}
		marginKeeper.InitGenesis(ctx, want)

		got := marginKeeper.ExportGenesis(ctx)

		require.Equal(t, *got, want)
	})
}

func TestKeeper_InitGenesis(t *testing.T) {
	t.Run("params with empty fields", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper
		require.NotNil(t, marginKeeper)

		params := types.Params{}
		want := types.GenesisState{Params: &params}
		validatorUpdate := marginKeeper.InitGenesis(ctx, want)

		require.Equal(t, validatorUpdate, []abci.ValidatorUpdate{})

		got := marginKeeper.ExportGenesis(ctx)

		require.Equal(t, *got, want)
	})

	t.Run("params fields set", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper
		require.NotNil(t, marginKeeper)

		params := types.Params{
			LeverageMax:          sdk.NewUint(10),
			InterestRateMax:      sdk.NewDec(5),
			InterestRateMin:      sdk.NewDec(1),
			InterestRateIncrease: sdk.NewDec(1),
			InterestRateDecrease: sdk.NewDec(1),
			HealthGainFactor:     sdk.NewDec(1),
			EpochLength:          1,
		}
		want := types.GenesisState{Params: &params}
		validatorUpdate := marginKeeper.InitGenesis(ctx, want)

		require.Equal(t, validatorUpdate, []abci.ValidatorUpdate{})

		got := marginKeeper.ExportGenesis(ctx)

		require.Equal(t, *got, want)
	})
}
