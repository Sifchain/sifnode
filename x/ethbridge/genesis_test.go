package ethbridge_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge"
	"github.com/Sifchain/sifnode/x/ethbridge/test"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	"math/rand"
)

// const
const (
	testAddress      = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	lockBurnSequence = uint64(1)
	globalNonce      = uint64(1)
	blockNumber      = uint64(1)
)

var valAddress, _ = sdk.ValAddressFromBech32(testAddress)

func TestExportGenesisExportsCorrectValue(t *testing.T) {
	ctx, keeper := test.CreateTestAppEthBridge(false)
	// Generate State
	receiver := createState(ctx, keeper, t)
	state := ethbridge.ExportGenesis(ctx, keeper)

	assert.NotNil(t, state, "ExportGenesis should have non-nil output")

	// Verify CrosschainFeeReceiveAccount
	assert.Equal(t, receiver, state.CrosschainFeeReceiveAccount)
	// Verify EtehreumLockBurnSequence
	assert.Equal(t, keeper.GetEthereumLockBurnSequences(ctx), state.EthereumLockBurnSequence)
	// Verify GlobalNonce
	assert.Equal(t, keeper.GetGlobalSequences(ctx), state.GlobalNonce)
	// Verify GlobalSequenceBlockNumber
	assert.Equal(t, keeper.GetGlobalSequenceToBlockNumbers(ctx), state.GlobalNonceBlockNumber)
}

// InitGenesis and ExportGenesis should be inverse of each other.
// aka InitGenesis(ExportGenesis(keeper)) === keeper
func TestInitGenesisWithExportGenesisDataSuccessful(t *testing.T) {
	ctx1, oldKeeper := test.CreateTestAppEthBridge(false)
	ctx2, newKeeper := test.CreateTestAppEthBridge(false)
	// Generate State
	createState(ctx1, oldKeeper, t)

	exportedState := ethbridge.ExportGenesis(ctx1, oldKeeper)
	// no validator update from the module
	valUpdates := ethbridge.InitGenesis(ctx2, newKeeper, *exportedState)
	assert.Equal(t, len(valUpdates), 0)

	// after init the genesis from state, receive account is set
	assert.Equal(t, oldKeeper.GetCrossChainFeeReceiverAccount(ctx1), newKeeper.GetCrossChainFeeReceiverAccount(ctx2))

	assert.Equal(t, oldKeeper.GetEthereumLockBurnSequences(ctx1), newKeeper.GetEthereumLockBurnSequences(ctx2))

	// check block number for network
	actualBlockNumber := newKeeper.GetGlobalSequence(ctx2, test.NetworkDescriptor)
	assert.Equal(t, blockNumber, actualBlockNumber)

	// TODO: Need to make the states more complex, Import actually fails, it is comparing default values.
	assert.Equal(t, oldKeeper.GetGlobalSequences(ctx1), newKeeper.GetGlobalSequences(ctx2))

	// check block number for network and global nonce
	// actualBlockNumber = keeper2.GetGlobalSequenceToBlockNumber(ctx2, test.NetworkDescriptor, globalNonce)
	// assert.Equal(t, blockNumber, actualBlockNumber)

	assert.Equal(t, oldKeeper.GetGlobalSequenceToBlockNumbers(ctx1), newKeeper.GetGlobalSequenceToBlockNumbers(ctx2))
}

func TestInitGenesisWithExportGenesisNonEmptyEthereumLockBurnSequence(t *testing.T) {
	ctx1, oldKeeper, _, _, _, _, _, validators := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	ctx2, newKeeper := test.CreateTestAppEthBridge(false)

	testValidatorAddress := validators[0]

	expectedEvmNetworkToLockBurnSequence := make(map[oracletypes.NetworkDescriptor]uint64)
	expectedEvmNetworkToLockBurnSequence[oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN_TESTNET] = rand.Uint64()
	expectedEvmNetworkToLockBurnSequence[oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_GANACHE] = rand.Uint64()
	expectedEvmNetworkToLockBurnSequence[oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_HARDHAT] = rand.Uint64()

	for networkDescriptor, sequence := range expectedEvmNetworkToLockBurnSequence {
		oldKeeper.SetEthereumLockBurnSequence(ctx1, networkDescriptor, testValidatorAddress, sequence)
	}

	exportedState := ethbridge.ExportGenesis(ctx1, oldKeeper)
	ethbridge.InitGenesis(ctx2, newKeeper, *exportedState)

	assert.Equal(t, oldKeeper.GetEthereumLockBurnSequences(ctx1), newKeeper.GetEthereumLockBurnSequences(ctx2))
	assert.Equal(t, len(expectedEvmNetworkToLockBurnSequence), len(newKeeper.GetEthereumLockBurnSequences(ctx2)))
}

func TestInitGenesisWithExportGenesisEthereumLockBurnSequenceMultipleValidators(t *testing.T) {
	ctx1, oldKeeper, _, _, _, _, _, validators := test.CreateTestKeepers(t, 0.7, []int64{3, 3, 2}, "")
	ctx2, newKeeper := test.CreateTestAppEthBridge(false)

	for _, validator := range validators {
		oldKeeper.SetEthereumLockBurnSequence(ctx1, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_GANACHE, validator, rand.Uint64())
		oldKeeper.SetEthereumLockBurnSequence(ctx1, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_HARDHAT, validator, rand.Uint64())
	}

	exportedState := ethbridge.ExportGenesis(ctx1, oldKeeper)
	ethbridge.InitGenesis(ctx2, newKeeper, *exportedState)

	assert.Equal(t, oldKeeper.GetEthereumLockBurnSequences(ctx1), newKeeper.GetEthereumLockBurnSequences(ctx2))
}

func TestInitGenesisWithExportGenesisDataGlobalSequenceSuccess(t *testing.T) {
	ctx1, oldKeeper := test.CreateTestAppEthBridge(false)
	ctx2, newKeeper := test.CreateTestAppEthBridge(false)

	// TODO: Cleanup
	oldKeeper.UpdateGlobalSequence(ctx1, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN_TESTNET, 5)
	oldKeeper.UpdateGlobalSequence(ctx1, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN_TESTNET, 34)
	oldKeeper.UpdateGlobalSequence(ctx1, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN_TESTNET, 51)

	oldKeeper.UpdateGlobalSequence(ctx1, test.NetworkDescriptor, 13)
	oldKeeper.UpdateGlobalSequence(ctx1, test.NetworkDescriptor, 53)

	exportedState := ethbridge.ExportGenesis(ctx1, oldKeeper)
	ethbridge.InitGenesis(ctx2, newKeeper, *exportedState)

	assert.Equal(t, oldKeeper.GetGlobalSequences(ctx1), newKeeper.GetGlobalSequences(ctx2))

}

func createState(ctx sdk.Context, keeper ethbridge.Keeper, t *testing.T) string {
	//Setting Receiver account
	receiver := test.GenerateAddress("")
	keeper.SetCrossChainFeeReceiverAccount(ctx, receiver)
	set := keeper.IsCrossChainFeeReceiverAccount(ctx, receiver)
	assert.True(t, set)

	ethereumLockBurnSequences := keeper.GetEthereumLockBurnSequences(ctx)
	assert.Equal(t, len(ethereumLockBurnSequences), 0, "New instances should have 0")

	globalNonces := keeper.GetGlobalSequences(ctx)
	assert.Equal(t, len(globalNonces), 0)

	keeper.SetEthereumLockBurnSequence(ctx, test.NetworkDescriptor, valAddress, lockBurnSequence)
	keeper.UpdateGlobalSequence(ctx, test.NetworkDescriptor, blockNumber)
	keeper.SetGlobalSequenceToBlockNumber(ctx, test.NetworkDescriptor, globalNonce, blockNumber)

	return receiver.String()
}
