package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
)

func TestGetAndSetFirstLockDoublePeg(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	denom := "denom"
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM

	// test the init value
	FirstLockDoublePeg := keeper.GetFirstLockDoublePeg(ctx, denom, networkDescriptor)
	assert.Equal(t, FirstLockDoublePeg, false)

	// test the value after set
	keeper.SetFirstLockDoublePeg(ctx, denom, networkDescriptor)

	FirstLockDoublePeg = keeper.GetFirstLockDoublePeg(ctx, denom, networkDescriptor)
	assert.Equal(t, FirstLockDoublePeg, true)
}
