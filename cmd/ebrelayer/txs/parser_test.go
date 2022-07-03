package txs

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

var (
	sugaredLogger = NewZapSugaredLogger()
)

func NewZapSugaredLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

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
		TestNetworkDescriptor, testBridgeContractAddress, TestNonce, strings.ToLower(TestSymbol), testTokenContractAddress,
		testEthereumAddress, testCosmosAddress, testCosmosValidatorBech32Address, testSDKAmount, ethbridge.ClaimType_CLAIM_TYPE_LOCK,
		TestName, TestDecimals)

	// Create test ethereum event
	ethereumEvent := CreateTestLogEthereumEvent(t)

	symbolTranslator := symbol_translator.NewSymbolTranslator()
	ethBridgeClaim, err := EthereumEventToEthBridgeClaim(testCosmosValidatorBech32Address, ethereumEvent, symbolTranslator, sugaredLogger)
	require.NoError(t, err)

	require.Equal(t, expectedEthBridgeClaim, &ethBridgeClaim)
}

func TestEthereumEventToEthBridgeClaimAcceptsEthSymbolNonNullAddress(t *testing.T) {
	testBridgeContractAddress := ethbridge.NewEthereumAddress(TestBridgeContractAddress)
	testTokenContractAddress := ethbridge.NewEthereumAddress(TestEthereumAddress2)
	testEthereumAddress := ethbridge.NewEthereumAddress(TestEthereumAddress1)
	testCosmosAddress, _ := sdk.AccAddressFromBech32(TestCosmosAddress1)

	symbolTranslator := symbol_translator.NewSymbolTranslator()

	testRawCosmosValidatorAddress, _ := sdk.AccAddressFromBech32(TestCosmosAddress2)
	testCosmosValidatorBech32Address := sdk.ValAddress(testRawCosmosValidatorAddress)

	ethereumEvent := CreateTestLogEthereumEvent(t)
	ethereumEvent.Symbol = "eth"
	ethereumEvent.Token = common.Address(testTokenContractAddress)

	expectedEthBridgeClaim := ethbridge.NewEthBridgeClaim(
		TestNetworkDescriptor, testBridgeContractAddress, TestNonce,
		"eth", testTokenContractAddress,
		testEthereumAddress, testCosmosAddress,
		testCosmosValidatorBech32Address, testSDKAmount, ethbridge.ClaimType_CLAIM_TYPE_LOCK,
		TestName, TestDecimals)

	ethBridgeClaim, err := EthereumEventToEthBridgeClaim(testCosmosValidatorBech32Address,
		ethereumEvent,
		symbolTranslator,
		sugaredLogger)

	require.NoError(t, err)

	require.Equal(t, expectedEthBridgeClaim, &ethBridgeClaim)
}

func TestDenomCalculated(t *testing.T) {
	address := ethbridge.NewEthereumAddress("0X0000000000000000000000000000000000000000")
	// Test strings with uppercase values
	expectedDenom := "sifBridge00010x0000000000000000000000000000000000000000"
	actualDenom := ethbridge.GetDenom(1, address)
	require.Equal(t, expectedDenom, actualDenom)
	// Test strings with lowercase values
	actualDenom = ethbridge.GetDenom(1, address)
	require.Equal(t, expectedDenom, actualDenom)
}

func TestBurnEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgBurn
	expectedMsgBurn := CreateTestCosmosMsg(t, types.MsgBurn)

	// Create MsgBurn attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t, types.MsgBurn)
	msgBurn, err := BurnLockEventToCosmosMsg(cosmosMsgAttributes, sugaredLogger)

	require.Nil(t, err)
	require.Equal(t, expectedMsgBurn, msgBurn)
}

func TestLockEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgLock
	expectedMsgLock := CreateTestCosmosMsg(t, types.MsgLock)

	// Create MsgLock attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t, types.MsgLock)
	msgLock, err := BurnLockEventToCosmosMsg(cosmosMsgAttributes, sugaredLogger)

	require.Nil(t, err)
	require.Equal(t, expectedMsgLock, msgLock)
}

func TestFailedBurnEventToCosmosMsg(t *testing.T) {
	// Create MsgBurn attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgIncompleteAttributes(t, types.MsgBurn)
	_, err := BurnLockEventToCosmosMsg(cosmosMsgAttributes, sugaredLogger)

	require.Error(t, err)
}

func TestFailedLockEventToCosmosMsg(t *testing.T) {
	// Create MsgLock attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgIncompleteAttributes(t, types.MsgLock)
	_, err := BurnLockEventToCosmosMsg(cosmosMsgAttributes, sugaredLogger)

	require.Error(t, err)
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
