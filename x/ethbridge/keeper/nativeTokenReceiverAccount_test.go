package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
)

const (
	TestAddress   = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	SecondAddress = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
)

func TestSetNativeTokenReceiverAccount(t *testing.T) {
	var ctx, keeper, _, _, _, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)

	keeper.SetNativeTokenReceiverAccount(ctx, testCosmosAddress)
	NativeTokenReceiverAccount := keeper.GetNativeTokenReceiverAccount(ctx)
	assert.Equal(t, NativeTokenReceiverAccount, testCosmosAddress)
}

func TestIsNativeTokenReceiverAccount(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)

	keeper.SetNativeTokenReceiverAccount(ctx, testCosmosAddress)
	assert.True(t, keeper.IsNativeTokenReceiverAccount(ctx, testCosmosAddress))
	testCosmosAddress, err = sdk.AccAddressFromBech32(SecondAddress)
	require.NoError(t, err)

	assert.False(t, keeper.IsNativeTokenReceiverAccount(ctx, testCosmosAddress))
}

func TestIsNativeTokenReceiverAccountSet(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	accountSet := keeper.IsNativeTokenReceiverAccountSet(ctx)
	require.Equal(t, accountSet, true)
}
