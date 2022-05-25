package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestProphecy(t *testing.T) {
	var ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	prophecyID := []byte{1, 2, 3, 4, 5, 6}
	networkDescriptor := types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	cosmosSender := "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	cosmosSenderSequence := uint64(1)
	ethereumReceiver := "0x0000000000000000000000000000000000000000"
	tokenDenomHash := "rowan"
	tokenContractAddress := "0x0000000000000000000000000000000000000000"
	tokenAmount := sdk.NewInt(1)
	crosschainFee := sdk.NewInt(1)
	bridgeToken := true
	globalSequence := uint64(1)
	tokenDecimal := uint8(1)
	tokenName := "name"
	tokenSymbol := "symbol"

	_, ok := keeper.GetProphecy(ctx, prophecyID)
	// ok should false since prophecy not stored yet
	assert.Equal(t, ok, false)

	_, ok = keeper.GetProphecyInfo(ctx, prophecyID)
	// ok should false since prophecy info not stored yet
	assert.Equal(t, ok, false)

	prophecy := types.Prophecy{
		Id:              prophecyID,
		Status:          types.StatusText_STATUS_TEXT_PENDING,
		ClaimValidators: []string{},
	}
	keeper.SetProphecy(ctx, prophecy)

	err := keeper.SetProphecyInfo(
		ctx,
		prophecyID,
		networkDescriptor,
		cosmosSender,
		cosmosSenderSequence,
		ethereumReceiver,
		tokenDenomHash,
		tokenContractAddress,
		tokenAmount,
		crosschainFee,
		bridgeToken,
		globalSequence,
		tokenDecimal,
		tokenName,
		tokenSymbol,
	)
	assert.Equal(t, err, nil)

	_, ok = keeper.GetProphecy(ctx, prophecyID)
	// ok should true since prophecy stored in keeper
	assert.Equal(t, ok, true)

	_, ok = keeper.GetProphecyInfo(ctx, prophecyID)
	// ok should true since prophecy info stored in keeper
	assert.Equal(t, ok, true)

	// test clean outdated prophecy

	// try to clean before ProphecyLiftTime
	keeper.UpdateCurrentHeight(260000)
	keeper.CleanUpProphecy(ctx)

	// still valid, not removed from keeper
	_, ok = keeper.GetProphecyInfo(ctx, prophecyID)
	// ok should true since prophecy info not removed yet
	assert.Equal(t, ok, true)

	// remove prophecy with correct block number
	keeper.UpdateCurrentHeight(521000)

	keeper.CleanUpProphecy(ctx)

	// become invalid, removed from keeper
	_, ok = keeper.GetProphecyInfo(ctx, prophecyID)
	// ok should false since prophecy info removed successful
	assert.Equal(t, ok, false)

}

func TestNetworkDescriptorGlobalNonce(t *testing.T) {
	var ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	prophecyID := []byte{1, 2, 3, 4, 5, 6}
	networkDescriptor := types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	globalSequence := uint64(1)

	// should be false since global sequnce not set yet
	_, ok := keeper.GetProphecyIDByNetworkDescriptorGlobalNonce(ctx, networkDescriptor, globalSequence)
	assert.Equal(t, ok, false)

	keeper.SetGlobalNonceProphecyID(ctx, networkDescriptor, globalSequence, prophecyID)

	id, ok := keeper.GetProphecyIDByNetworkDescriptorGlobalNonce(ctx, networkDescriptor, globalSequence)
	// check the id in keeper is the same with prophecy id set before
	assert.Equal(t, ok, true)
	assert.Equal(t, id, prophecyID)

}

func TestDeleteProphecyInfo(t *testing.T) {
	var ctx, _, _, _, keeper, _, _, _ = test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	prophecyID := []byte{1, 2, 3, 4, 5, 6}
	networkDescriptor := types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
	cosmosSender := "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	cosmosSenderSequence := uint64(1)
	ethereumReceiver := "0x0000000000000000000000000000000000000000"
	tokenDenomHash := "rowan"
	tokenContractAddress := "0x0000000000000000000000000000000000000000"
	tokenAmount := sdk.NewInt(1)
	crosschainFee := sdk.NewInt(1)
	bridgeToken := true
	globalSequence := uint64(1)
	tokenDecimal := uint8(1)
	tokenName := "name"
	tokenSymbol := "symbol"

	_, ok := keeper.GetProphecy(ctx, prophecyID)
	// ok should false since prophecy not stored yet
	assert.Equal(t, ok, false)

	_, ok = keeper.GetProphecyInfo(ctx, prophecyID)
	// ok should false since prophecy info not stored yet
	assert.Equal(t, ok, false)

	_, ok = keeper.GetProphecyIDByNetworkDescriptorGlobalNonce(ctx, networkDescriptor, 1)
	// ok should false since prophecy info not stored yet
	assert.Equal(t, ok, false)

	prophecy := types.Prophecy{
		Id:              prophecyID,
		Status:          types.StatusText_STATUS_TEXT_PENDING,
		ClaimValidators: []string{},
	}
	keeper.SetProphecy(ctx, prophecy)

	err := keeper.SetProphecyInfo(
		ctx,
		prophecyID,
		networkDescriptor,
		cosmosSender,
		cosmosSenderSequence,
		ethereumReceiver,
		tokenDenomHash,
		tokenContractAddress,
		tokenAmount,
		crosschainFee,
		bridgeToken,
		globalSequence,
		tokenDecimal,
		tokenName,
		tokenSymbol,
	)
	assert.Equal(t, err, nil)

	prophecyInfo, ok := keeper.GetProphecyInfo(ctx, prophecyID)
	// ok should true after set prophecy info
	assert.Equal(t, ok, true)

	_, ok = keeper.GetProphecyIDByNetworkDescriptorGlobalNonce(ctx, networkDescriptor, 1)
	// ok should true after set prophecy info
	assert.Equal(t, ok, true)

	keeper.DeleteProphecyInfo(ctx, prophecyInfo)

	_, ok = keeper.GetProphecy(ctx, prophecyID)
	// ok should true because the basic info always stored
	assert.Equal(t, ok, true)

	_, ok = keeper.GetProphecyInfo(ctx, prophecyID)
	// ok should false after prophecy info removed
	assert.Equal(t, ok, false)

	_, ok = keeper.GetProphecyIDByNetworkDescriptorGlobalNonce(ctx, networkDescriptor, 1)
	// ok should false since prophecy info removed
	assert.Equal(t, ok, false)

}
