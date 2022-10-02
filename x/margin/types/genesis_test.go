package types_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/stretchr/testify/require"
)

func TestTypes_DefaultGenesis(t *testing.T) {
	got := types.DefaultGenesis()

	require.NotNil(t, got.Params.LeverageMax)
	require.Equal(t, "2.000000000000000000", got.Params.LeverageMax.String())
	require.NotNil(t, got.Params.InterestRateMin)
	require.Equal(t, "0.005000000000000000", got.Params.InterestRateMin.String())
	require.NotNil(t, got.Params.InterestRateMax)
	require.Equal(t, "3.000000000000000000", got.Params.InterestRateMax.String())
	require.NotNil(t, got.Params.InterestRateIncrease)
	require.Equal(t, "0.100000000000000000", got.Params.InterestRateIncrease.String())
	require.NotNil(t, got.Params.InterestRateDecrease)
	require.Equal(t, "0.100000000000000000", got.Params.InterestRateDecrease.String())
	require.NotNil(t, got.Params.HealthGainFactor)
	require.Equal(t, "1.000000000000000000", got.Params.HealthGainFactor.String())
	require.NotNil(t, got.Params.EpochLength)
	require.Equal(t, int64(1), got.Params.EpochLength)
}
