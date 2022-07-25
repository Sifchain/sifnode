package ethbridge_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge"
	"github.com/Sifchain/sifnode/x/ethbridge/test"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// const
const (
	testAddress      = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	lockBurnSequence = uint64(1)
	globalNonce      = uint64(1)
	blockNumber      = uint64(1)
)

var valAddress, _ = sdk.ValAddressFromBech32(testAddress)

func TestExportGenesis(t *testing.T) {
	ctx, keeper := test.CreateTestAppEthBridge(false)
	// Generate State
	receiver := CreateState(ctx, keeper, t)
	state := ethbridge.ExportGenesis(ctx, keeper)
	assert.Equal(t, receiver, state.CrosschainFeeReceiveAccount)

	key := keeper.GetEthereumLockBurnSequencePrefix(test.NetworkDescriptor, valAddress)
	assert.Equal(t, lockBurnSequence, state.EthereumLockBurnSequence[string(key)])

	// After set, global nonce will be +1
	assert.Equal(t, globalNonce+1, state.GlobalNonce[uint32(test.NetworkDescriptor)])

	key = keeper.GetGlobalSequenceToBlockNumberPrefix(ctx, test.NetworkDescriptor, globalNonce)
	assert.Equal(t, blockNumber, state.GlobalNonceBlockNumber[string(key)])
}

func TestInitGenesis(t *testing.T) {
	ctx1, keeper1 := test.CreateTestAppEthBridge(false)
	ctx2, keeper2 := test.CreateTestAppEthBridge(false)
	// Generate State
	receiver := CreateState(ctx1, keeper1, t)
	state := ethbridge.ExportGenesis(ctx1, keeper1)
	assert.Equal(t, state.CrosschainFeeReceiveAccount, receiver)
	state2 := ethbridge.ExportGenesis(ctx2, keeper2)
	assert.Equal(t, state2.CrosschainFeeReceiveAccount, "")

	valUpdates := ethbridge.InitGenesis(ctx2, keeper2, *state)
	assert.Equal(t, len(valUpdates), 0)

	actualReceiver := keeper2.GetCrossChainFeeReceiverAccount(ctx2)
	assert.Equal(t, receiver, actualReceiver.String())
}

func CreateState(ctx sdk.Context, keeper ethbridge.Keeper, t *testing.T) string {
	//Setting Receiver account
	receiver := test.GenerateAddress("")
	keeper.SetCrossChainFeeReceiverAccount(ctx, receiver)
	set := keeper.IsCrossChainFeeReceiverAccount(ctx, receiver)
	assert.True(t, set)

	ethereumLockBurnSequences := keeper.GetEthereumLockBurnSequences(ctx)
	assert.Equal(t, len(ethereumLockBurnSequences), 0)

	globalNonces := keeper.GetGlobalSequences(ctx)
	assert.Equal(t, len(globalNonces), 0)

	keeper.SetEthereumLockBurnSequence(ctx, test.NetworkDescriptor, valAddress, lockBurnSequence)
	keeper.UpdateGlobalSequence(ctx, test.NetworkDescriptor, blockNumber)
	keeper.SetGlobalSequenceToBlockNumber(ctx, test.NetworkDescriptor, globalNonce, blockNumber)

	return receiver.String()
}
