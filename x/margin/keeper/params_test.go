//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper_test

import (
	"fmt"
	"testing"

	"github.com/Sifchain/sifnode/x/margin/test"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_ParamGetter(t *testing.T) {
	ctx, app := test.CreateTestAppMargin(false)
	marginKeeper := app.MarginKeeper

	data := types.GenesisState{Params: &types.Params{
		LeverageMax:          sdk.NewUint(10),
		InterestRateMax:      sdk.NewDec(5),
		InterestRateMin:      sdk.NewDec(1),
		InterestRateIncrease: sdk.NewDec(1),
		InterestRateDecrease: sdk.NewDec(1),
		HealthGainFactor:     sdk.NewDec(1),
		EpochLength:          1,
		ForceCloseThreshold:  sdk.NewDec(1), //TODO get real default
	}}
	marginKeeper.InitGenesis(ctx, data)

	paramGetterTests := []struct {
		name   string
		want   string
		method func(sdk.Context) string
	}{
		{
			name:   "LeverageMax",
			want:   data.Params.LeverageMax.String(),
			method: func(ctx sdk.Context) string { return marginKeeper.GetLeverageParam(ctx).String() },
		},
		{
			name:   "InterestRateMax",
			want:   data.Params.InterestRateMax.String(),
			method: func(ctx sdk.Context) string { return marginKeeper.GetInterestRateMax(ctx).String() },
		},
		{
			name:   "InterestRateMin",
			want:   data.Params.InterestRateMin.String(),
			method: func(ctx sdk.Context) string { return marginKeeper.GetInterestRateMin(ctx).String() },
		},
		{
			name:   "InterestRateIncrease",
			want:   data.Params.InterestRateIncrease.String(),
			method: func(ctx sdk.Context) string { return marginKeeper.GetInterestRateIncrease(ctx).String() },
		},
		{
			name:   "InterestRateDecrease",
			want:   data.Params.InterestRateDecrease.String(),
			method: func(ctx sdk.Context) string { return marginKeeper.GetInterestRateDecrease(ctx).String() },
		},
		{
			name:   "HealthGainFactor",
			want:   data.Params.HealthGainFactor.String(),
			method: func(ctx sdk.Context) string { return marginKeeper.GetHealthGainFactor(ctx).String() },
		},
		{
			name:   "EpochLength",
			want:   fmt.Sprint(data.Params.EpochLength),
			method: func(ctx sdk.Context) string { return fmt.Sprint(marginKeeper.GetEpochLength(ctx)) },
		},
		{
			name:   "ForceCloseThreshold",
			want:   fmt.Sprint(data.Params.ForceCloseThreshold),
			method: func(ctx sdk.Context) string { return fmt.Sprint(marginKeeper.GetForceCloseThreshold(ctx)) },
		},
	}

	for _, tt := range paramGetterTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method(ctx)

			require.Equal(t, got, tt.want)
		})
	}
}
