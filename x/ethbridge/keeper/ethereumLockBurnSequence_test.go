package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
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
