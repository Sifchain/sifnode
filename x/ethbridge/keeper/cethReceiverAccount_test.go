package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/internal"
)

const (
	TestAddress   = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	SecondAddress = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
)

func TestSetCethReceiverAccount(t *testing.T) {
	var ctx, keeper, _, _, _, _, _ = internal.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)

	keeper.SetCethReceiverAccount(ctx, testCosmosAddress)
	CethReceiverAccount := keeper.GetCethReceiverAccount(ctx)
	assert.Equal(t, CethReceiverAccount, testCosmosAddress)
}

func TestIsCethReceiverAccount(t *testing.T) {
	ctx, keeper, _, _, _, _, _ := internal.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)

	keeper.SetCethReceiverAccount(ctx, testCosmosAddress)
	assert.True(t, keeper.IsCethReceiverAccount(ctx, testCosmosAddress))
	testCosmosAddress, err = sdk.AccAddressFromBech32(SecondAddress)
	require.NoError(t, err)

	assert.False(t, keeper.IsCethReceiverAccount(ctx, testCosmosAddress))
}

func TestIsCethReceiverAccountSet(t *testing.T) {
	ctx, keeper, _, _, _, _, _ := internal.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	accountSet := keeper.IsCethReceiverAccountSet(ctx)
	require.Equal(t, accountSet, true)
}
