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

func TestGetGlobalSequenceToBlockNumbers(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM

	// test the init value
	globalNonceOne := uint64(1)
	blockNumber := uint64(100)
	globalNonce := keeper.GetGlobalSequence(ctx, networkDescriptor)
	assert.Equal(t, globalNonce, globalNonceOne)

	// set block number = sequence * 100
	for index := 0; index < 10; index++ {
		keeper.UpdateGlobalSequence(ctx, networkDescriptor, blockNumber+uint64(index)*100)
	}

	// get all GlobalSequenceToBlockNumbers and verify the value
	sequenceToBlockNumbers := keeper.GetGlobalSequenceToBlockNumbers(ctx)
	for _, value := range sequenceToBlockNumbers {
		assert.Equal(t, value.GlobalSequenceKey.NetworkDescriptor, networkDescriptor)
		assert.Equal(t, value.GlobalSequenceKey.GlobalSequence*100, value.BlockNumber.BlockNumber)
	}

}

func TestGetGlobalSequences(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	networkDescriptors := []oracletypes.NetworkDescriptor{
		oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM,
		oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BITCOIN,
		oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM_TESTNET_ROPSTEN,
		oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN,
		oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN_TESTNET,
		oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_GANACHE,
		oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_HARDHAT,
	}

	// set the sequence the same value as NetworkDescriptor enum
	for index := 0; index < len(networkDescriptors); index++ {
		sequence := oracletypes.GlobalSequence{
			GlobalSequence: uint64(networkDescriptors[index]),
		}
		keeper.SetGlobalSequence(ctx, networkDescriptors[index], sequence)
	}

	// check the sequences from keeper
	sequences := keeper.GetGlobalSequences(ctx)
	for _, value := range sequences {
		assert.Equal(t, uint64(value.NetworkDescriptor), value.GlobalSequence.GlobalSequence)

	}
}
