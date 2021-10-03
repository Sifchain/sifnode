package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/stretchr/testify/assert"
)

func TestGetAndSetFirstLockDoublePeg(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	denom := "denom"

	// test the init value
	FirstLockDoublePeg := keeper.GetFirstLockDoublePeg(ctx, denom)
	assert.Equal(t, FirstLockDoublePeg, false)

	// test the value after set
	keeper.SetFirstLockDoublePeg(ctx, denom)

	FirstLockDoublePeg = keeper.GetFirstLockDoublePeg(ctx, denom)
	assert.Equal(t, FirstLockDoublePeg, true)
}
