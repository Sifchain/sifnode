package keeper_test

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	testNetwork          = 1
	testAddress          = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	testLockBurnSequence = uint64(10)
	testInitNonce        = uint64(0)
)

func TestSetEthereumLockBurnSequence(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	keeper.SetEthereumLockBurnSequence(ctx, testNetwork, testCosmosAddress, testLockBurnSequence)

	LockBurnSequence := keeper.GetEthereumLockBurnSequence(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, LockBurnSequence, testLockBurnSequence)
}

func TestGetEthereumLockBurnSequence(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	LockBurnSequence := keeper.GetEthereumLockBurnSequence(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, LockBurnSequence, testInitNonce)
}

func TestGetEthereumLockBurnSequences(t *testing.T) {
	// create some validators
	validatorPowers := []int64{3, 3, 3, 3, 3}
	var ctx, keeper, _, _, _, _, _, valAddresses = test.CreateTestKeepers(t, 0.7, validatorPowers, "")
	networkDescriptor := oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	sequence := uint64(100)

	// set different sequence number
	for index := 0; index < len(validatorPowers); index++ {
		keeper.SetEthereumLockBurnSequence(ctx, networkDescriptor, valAddresses[index], sequence*uint64(index))
	}

	// verify the EthereumLockBurnSequences data from keeper
	for _, value := range keeper.GetEthereumLockBurnSequences(ctx) {
		index := 0
		for ; index < len(validatorPowers); index++ {
			if bytes.Compare(valAddresses[index], value.EthereumLockBurnSequenceKey.ValidatorAddress) == 0 {
				break
			}
		}
		assert.Less(t, index, len(validatorPowers))
		assert.Equal(t, value.EthereumLockBurnSequenceKey.NetworkDescriptor, networkDescriptor)
		assert.Equal(t, value.EthereumLockBurnSequence.EthereumLockBurnSequence, sequence*uint64(index))
	}
}
