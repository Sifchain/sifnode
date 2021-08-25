package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
)

const (
	testNetwork       = 1
	testAddress       = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	testLockBurnNonce = uint64(10)
	testInitNonce     = uint64(0)
)

func TestSetEthereumLockBurnNonce(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	keeper.SetEthereumLockBurnNonce(ctx, testNetwork, testCosmosAddress, testLockBurnNonce)

	lockBurnNonce := keeper.GetEthereumLockBurnNonce(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, lockBurnNonce, testLockBurnNonce)
}

func TestGetEthereumLockBurnNonce(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.ValAddressFromBech32(testAddress)
	require.NoError(t, err)

	lockBurnNonce := keeper.GetEthereumLockBurnNonce(ctx, testNetwork, testCosmosAddress)
	assert.Equal(t, lockBurnNonce, testInitNonce)
}
