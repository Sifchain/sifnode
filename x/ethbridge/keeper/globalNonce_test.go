package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

func TestGetAndUpdateGlobalSequence(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM

	// test the init value
	globalNonceOne := uint64(1)
	blockNumber := uint64(100)
	globalNonce := keeper.GetGlobalSequence(ctx, networkDescriptor)
	assert.Equal(t, globalNonce, globalNonceOne)

	// test the second value
	keeper.UpdateGlobalSequence(ctx, networkDescriptor, blockNumber)

	globalNonceTwo := uint64(2)
	globalNonce = keeper.GetGlobalSequence(ctx, networkDescriptor)
	assert.Equal(t, globalNonce, globalNonceTwo)
}

func TestGetGlobalSequenceToBlockNumber(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM

	// test the init value
	globalNonceOne := uint64(1)
	initNonce := keeper.GetGlobalSequenceToBlockNumber(ctx, networkDescriptor, globalNonceOne)
	assert.Equal(t, initNonce, uint64(0))

	// set the block number
	blockNumber := uint64(100)
	globalNonce := keeper.GetGlobalSequence(ctx, networkDescriptor)
	assert.Equal(t, globalNonce, globalNonceOne)
	keeper.UpdateGlobalSequence(ctx, networkDescriptor, blockNumber)

	testBlockNumber := keeper.GetGlobalSequenceToBlockNumber(ctx, networkDescriptor, globalNonceOne)
	assert.Equal(t, testBlockNumber, blockNumber)
}
