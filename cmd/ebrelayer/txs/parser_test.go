package txs

import (
	"math/big"
	"os"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

func TestLogLockToEthBridgeClaim(t *testing.T) {
	// Set up testing variables
	testBridgeContractAddress := ethbridge.NewEthereumAddress(TestBridgeContractAddress)
	testTokenContractAddress := ethbridge.NewEthereumAddress(TestEthTokenAddress)
	testEthereumAddress := ethbridge.NewEthereumAddress(TestEthereumAddress1)
	// Cosmos account address
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestCosmosAddress1)
	require.NoError(t, err)
	// Cosmos validator address
	testRawCosmosValidatorAddress, err := sdk.AccAddressFromBech32(TestCosmosAddress2)
	require.NoError(t, err)
	testCosmosValidatorBech32Address := sdk.ValAddress(testRawCosmosValidatorAddress)

	// Set up expected EthBridgeClaim
	expectedEthBridgeClaim := ethbridge.NewEthBridgeClaim(
		TestEthereumChainID, testBridgeContractAddress, TestNonce, strings.ToLower(TestSymbol), testTokenContractAddress,
		testEthereumAddress, testCosmosAddress, testCosmosValidatorBech32Address, testSDKAmount, TestLockClaimType)

	// Create test ethereum event
	ethereumEvent := CreateTestLogEthereumEvent(t)

	ethBridgeClaim, err := EthereumEventToEthBridgeClaim(testCosmosValidatorBech32Address, ethereumEvent)
	require.NoError(t, err)

	require.Equal(t, expectedEthBridgeClaim, ethBridgeClaim)
}

func TestProphecyClaimToSignedOracleClaim(t *testing.T) {
	// Set ETHEREUM_PRIVATE_KEY env variable
	os.Setenv(EthereumPrivateKey, TestPrivHex)
	// Get and load private key from env variables
	rawKey := os.Getenv(EthereumPrivateKey)
	privateKey, _ := crypto.HexToECDSA(rawKey)

	// Create new test ProphecyClaimEvent
	prophecyClaimEvent := CreateTestProphecyClaimEvent(t)
	// Generate claim message from ProphecyClaim
	message := GenerateClaimMessage(prophecyClaimEvent)

	// Prepare the message (required for signature verification on contract)
	prefixedHashedMsg := PrefixMsg(message)

	// Sign the message using the validator's private key
	signature, err := SignClaim(prefixedHashedMsg, privateKey)
	require.NoError(t, err)

	var message32 [32]byte
	copy(message32[:], message)

	// Set up expected OracleClaim
	expectedOracleClaim := OracleClaim{
		ProphecyID: big.NewInt(int64(TestProphecyID)),
		Message:    message32,
		Signature:  signature,
	}

	// Map the test ProphecyClaim to a signed OracleClaim
	oracleClaim, err := ProphecyClaimToSignedOracleClaim(prophecyClaimEvent, privateKey)
	require.NoError(t, err)

	require.Equal(t, expectedOracleClaim, oracleClaim)
}

func TestBurnEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgBurn
	expectedMsgBurn := CreateTestCosmosMsg(t, types.MsgBurn)

	// Create MsgBurn attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t, types.MsgBurn)
	msgBurn, err := BurnLockEventToCosmosMsg(types.MsgBurn, cosmosMsgAttributes)

	require.Nil(t, err)
	require.Equal(t, expectedMsgBurn, msgBurn)
}

func TestLockEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgLock
	expectedMsgLock := CreateTestCosmosMsg(t, types.MsgLock)

	// Create MsgLock attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t, types.MsgLock)
	msgLock, err := BurnLockEventToCosmosMsg(types.MsgLock, cosmosMsgAttributes)

	require.Nil(t, err)
	require.Equal(t, expectedMsgLock, msgLock)
}

func TestFailedBurnEventToCosmosMsg(t *testing.T) {
	// Create MsgBurn attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgIncompleteAttributes(t, types.MsgBurn)
	_, err := BurnLockEventToCosmosMsg(types.MsgBurn, cosmosMsgAttributes)

	require.Error(t, err)
}

func TestFailedLockEventToCosmosMsg(t *testing.T) {
	// Create MsgLock attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgIncompleteAttributes(t, types.MsgLock)
	_, err := BurnLockEventToCosmosMsg(types.MsgLock, cosmosMsgAttributes)

	require.Error(t, err)
}

func TestMsgBurnToProphecyClaim(t *testing.T) {
	// Parse expected symbol
	res := strings.SplitAfter(strings.ToLower(TestSymbol), defaultSifchainPrefix)
	symbol := strings.Join(res[1:], "")

	// Set up expected ProphecyClaim
	expectedProphecyClaim := ProphecyClaim{
		ClaimType:            types.MsgBurn,
		CosmosSender:         []byte(TestCosmosAddress1),
		CosmosSenderSequence: big.NewInt(1),
		EthereumReceiver:     common.HexToAddress(TestEthereumAddress1),
		Symbol:               symbol,
		Amount:               testSDKAmount,
	}

	// Create a MsgBurn as input parameter
	testCosmosMsgBurn := CreateTestCosmosMsg(t, types.MsgBurn)
	prophecyClaim := CosmosMsgToProphecyClaim(testCosmosMsgBurn)

	require.Equal(t, expectedProphecyClaim, prophecyClaim)
}

func TestMsgLockToProphecyClaim(t *testing.T) {
	// Set up expected ProphecyClaim
	expectedProphecyClaim := ProphecyClaim{
		ClaimType:            types.MsgLock,
		CosmosSender:         []byte(TestCosmosAddress1),
		CosmosSenderSequence: big.NewInt(1),
		EthereumReceiver:     common.HexToAddress(TestEthereumAddress1),
		Symbol:               TestSymbol,
		Amount:               testSDKAmount,
	}

	// Create a MsgLock as input parameter
	testCosmosMsgLock := CreateTestCosmosMsg(t, types.MsgLock)
	prophecyClaim := CosmosMsgToProphecyClaim(testCosmosMsgLock)

	require.Equal(t, expectedProphecyClaim, prophecyClaim)
}

func TestIsZeroAddress(t *testing.T) {
	falseRes := isZeroAddress(common.HexToAddress(TestOtherAddress))
	require.False(t, falseRes)

	trueRes := isZeroAddress(common.HexToAddress(TestNullAddress))
	require.True(t, trueRes)
}

func TestAttributesToEthereumBridgeClaim(t *testing.T) {
	attributes := CreateEthereumBridgeClaimAttributes(t)
	claim, err := AttributesToEthereumBridgeClaim(attributes)
	require.NotEqual(t, claim, nil)
	require.Equal(t, err, nil)
}

func TestInvalidCosmosSenderAttributesToEthereumBridgeClaim(t *testing.T) {
	attributes := CreateInvalidCosmosSenderEthereumBridgeClaimAttributes(t)
	_, err := AttributesToEthereumBridgeClaim(attributes)
	require.Error(t, err)
}

func TestInvalidEthereumSenderAttributesToEthereumBridgeClaim(t *testing.T) {
	attributes := CreateInvalidEthereumSenderEthereumBridgeClaimAttributes(t)
	_, err := AttributesToEthereumBridgeClaim(attributes)
	require.Error(t, err)
}

func TestInvalidSequenceAttributesToEthereumBridgeClaim(t *testing.T) {
	attributes := CreateInvalidSequenceEthereumBridgeClaimAttributes(t)
	_, err := AttributesToEthereumBridgeClaim(attributes)
	require.Error(t, err)
}
