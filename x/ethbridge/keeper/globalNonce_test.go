package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

func TestGetAndUpdateGlobalNonce(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM

	// test the init value
	zeroGlobalNonce := uint64(1)
	globalNonce := keeper.GetAndUpdateGlobalNonce(ctx, networkDescriptor)
	assert.Equal(t, globalNonce, zeroGlobalNonce)

	// test the second value
	oneGlobalNonce := uint64(2)
	globalNonce = keeper.GetAndUpdateGlobalNonce(ctx, networkDescriptor)
	assert.Equal(t, globalNonce, oneGlobalNonce)
}
