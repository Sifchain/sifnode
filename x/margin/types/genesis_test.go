package types_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/stretchr/testify/require"
)

func TestTypes_DefaultGenesis(t *testing.T) {
	got := types.DefaultGenesis()

	require.NotNil(t, got.Params.LeverageMax)
	require.NotNil(t, got.Params.InterestRateMin)
	require.NotNil(t, got.Params.InterestRateMax)
	require.NotNil(t, got.Params.InterestRateIncrease)
	require.NotNil(t, got.Params.InterestRateDecrease)
	require.NotNil(t, got.Params.HealthGainFactor)
	require.NotNil(t, got.Params.EpochLength)
}
