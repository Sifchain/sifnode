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
	t.Run("no prior import of genesis, default param settings", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper
		require.NotNil(t, marginKeeper)
		state := marginKeeper.ExportGenesis(ctx)
		require.NotNil(t, state)

		require.Equal(t, state.Params.LeverageMax, marginKeeper.GetMaxLeverageParam(ctx))
		require.Equal(t, state.Params.InterestRateMax, marginKeeper.GetInterestRateMax(ctx))
		require.Equal(t, state.Params.InterestRateMin, marginKeeper.GetInterestRateMin(ctx))
		require.Equal(t, state.Params.InterestRateIncrease, marginKeeper.GetInterestRateIncrease(ctx))
		require.Equal(t, state.Params.InterestRateDecrease, marginKeeper.GetInterestRateDecrease(ctx))
		require.Equal(t, state.Params.HealthGainFactor, marginKeeper.GetHealthGainFactor(ctx))
		require.Equal(t, state.Params.EpochLength, marginKeeper.GetEpochLength(ctx))
	})

	t.Run("prior import of genesis then export", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper
		require.NotNil(t, marginKeeper)

		params := types.Params{
			LeverageMax:                              sdk.NewDec(2),
			HealthGainFactor:                         sdk.NewDec(1),
			InterestRateMin:                          sdk.NewDecWithPrec(5, 3),
			InterestRateMax:                          sdk.NewDec(3),
			InterestRateDecrease:                     sdk.NewDecWithPrec(1, 1),
			InterestRateIncrease:                     sdk.NewDecWithPrec(1, 1),
			ForceCloseFundPercentage:                 sdk.NewDecWithPrec(1, 1),
			ForceCloseFundAddress:                    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			IncrementalInterestPaymentFundPercentage: sdk.NewDecWithPrec(1, 1),
			IncrementalInterestPaymentFundAddress:    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			PoolOpenThreshold:                        sdk.NewDecWithPrec(1, 1),
			RemovalQueueThreshold:                    sdk.NewDecWithPrec(1, 1),
			EpochLength:                              1,
			MaxOpenPositions:                         10000,
			SqModifier:                               sdk.MustNewDecFromStr("10000000000000000000000000"),
			SafetyFactor:                             sdk.MustNewDecFromStr("1.05"),
		}
		want := types.GenesisState{Params: &params}
		marginKeeper.InitGenesis(ctx, want)

		got := marginKeeper.ExportGenesis(ctx)

		require.Equal(t, *got, want)
	})
}

func TestKeeper_InitGenesis(t *testing.T) {
	t.Run("params with initial fields", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper
		require.NotNil(t, marginKeeper)

		params := types.Params{
			LeverageMax:                              sdk.NewDec(2),
			HealthGainFactor:                         sdk.NewDec(1),
			InterestRateMin:                          sdk.NewDecWithPrec(5, 3),
			InterestRateMax:                          sdk.NewDec(3),
			InterestRateDecrease:                     sdk.NewDecWithPrec(1, 1),
			InterestRateIncrease:                     sdk.NewDecWithPrec(1, 1),
			ForceCloseFundPercentage:                 sdk.NewDecWithPrec(1, 1),
			ForceCloseFundAddress:                    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			IncrementalInterestPaymentFundPercentage: sdk.NewDecWithPrec(1, 1),
			IncrementalInterestPaymentFundAddress:    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			PoolOpenThreshold:                        sdk.NewDecWithPrec(1, 1),
			RemovalQueueThreshold:                    sdk.NewDecWithPrec(1, 1),
			EpochLength:                              1,
			MaxOpenPositions:                         10000,
			Pools:                                    []string{},
			SqModifier:                               sdk.MustNewDecFromStr("10000000000000000000000000"),
			SafetyFactor:                             sdk.MustNewDecFromStr("1.05"),
		}
		want := types.GenesisState{Params: &params}
		validatorUpdate := marginKeeper.InitGenesis(ctx, want)

		require.Equal(t, validatorUpdate, []abci.ValidatorUpdate{})

		got := marginKeeper.ExportGenesis(ctx)

		require.Equal(t, got.Params.LeverageMax, marginKeeper.GetMaxLeverageParam(ctx))
		require.Equal(t, got.Params.InterestRateMax, marginKeeper.GetInterestRateMax(ctx))
		require.Equal(t, got.Params.InterestRateMin, marginKeeper.GetInterestRateMin(ctx))
		require.Equal(t, got.Params.InterestRateIncrease, marginKeeper.GetInterestRateIncrease(ctx))
		require.Equal(t, got.Params.InterestRateDecrease, marginKeeper.GetInterestRateDecrease(ctx))
		require.Equal(t, got.Params.HealthGainFactor, marginKeeper.GetHealthGainFactor(ctx))
		require.Equal(t, got.Params.EpochLength, marginKeeper.GetEpochLength(ctx))
	})

	t.Run("params fields set", func(t *testing.T) {
		ctx, app := test.CreateTestAppMargin(false)
		marginKeeper := app.MarginKeeper
		require.NotNil(t, marginKeeper)

		params := types.Params{
			LeverageMax:                              sdk.NewDec(2),
			HealthGainFactor:                         sdk.NewDec(1),
			InterestRateMin:                          sdk.NewDecWithPrec(5, 3),
			InterestRateMax:                          sdk.NewDec(3),
			InterestRateDecrease:                     sdk.NewDecWithPrec(1, 1),
			InterestRateIncrease:                     sdk.NewDecWithPrec(1, 1),
			ForceCloseFundPercentage:                 sdk.NewDecWithPrec(1, 1),
			ForceCloseFundAddress:                    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			IncrementalInterestPaymentFundPercentage: sdk.NewDecWithPrec(1, 1),
			IncrementalInterestPaymentFundAddress:    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			PoolOpenThreshold:                        sdk.NewDecWithPrec(1, 1),
			RemovalQueueThreshold:                    sdk.NewDecWithPrec(1, 1),
			EpochLength:                              1,
			MaxOpenPositions:                         10000,
			SqModifier:                               sdk.MustNewDecFromStr("10000000000000000000000000"),
			SafetyFactor:                             sdk.MustNewDecFromStr("1.05"),
		}
		want := types.GenesisState{Params: &params}
		validatorUpdate := marginKeeper.InitGenesis(ctx, want)

		require.Equal(t, validatorUpdate, []abci.ValidatorUpdate{})

		got := marginKeeper.ExportGenesis(ctx)

		require.Equal(t, *got, want)
	})
}
