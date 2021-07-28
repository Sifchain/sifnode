package txs

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

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
	// Calculate Denom Hash String
	denomHash := ethbridge.GetDenomHash(TestNetworkDescriptor, TestEthTokenAddress, TestDecimals, TestName, TestSymbol)

	// Set up expected EthBridgeClaim
	expectedEthBridgeClaim := ethbridge.NewEthBridgeClaim(
		TestNetworkDescriptor, testBridgeContractAddress, TestNonce, strings.ToLower(TestSymbol), testTokenContractAddress,
		testEthereumAddress, testCosmosAddress, testCosmosValidatorBech32Address, testSDKAmount, ethbridge.ClaimType_CLAIM_TYPE_LOCK, TestName, TestDecimals, denomHash)

	// Create test ethereum event
	ethereumEvent := CreateTestLogEthereumEvent(t)

	ethBridgeClaim, err := EthereumEventToEthBridgeClaim(testCosmosValidatorBech32Address, ethereumEvent)
	require.NoError(t, err)

	require.Equal(t, expectedEthBridgeClaim, &ethBridgeClaim)
}

func TestDenomHashHandCrafted(t *testing.T) {
	expectedDenom := "BTC"
	actualDenom := ethbridge.GetDenomHash(2, "0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2", 9, "wBTC", "WBTC")

	require.Equal(t, expectedDenom, actualDenom)
}

func TestDenomCalculated(t *testing.T) {
	expectedDenom := "7ab5ab7aa7d978577efb3b68b9115976ec8e0312d04062350f4b40822a3870ce"
	actualDenom := ethbridge.GetDenomHash(1, "0x0000000000000000000000000000000000000000", 18, "Ethereum", "ETH")

	require.Equal(t, expectedDenom, actualDenom)
}

func TestBurnEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgBurn
	expectedMsgBurn := CreateTestCosmosMsg(t, types.MsgBurn)

	// Create MsgBurn attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t, types.MsgBurn)
	msgBurn, err := BurnLockEventToCosmosMsg(types.MsgBurn, cosmosMsgAttributes, sugaredLogger)

	require.Nil(t, err)
	require.Equal(t, expectedMsgBurn, msgBurn)
}

func TestLockEventToCosmosMsg(t *testing.T) {
	// Set up expected MsgLock
	expectedMsgLock := CreateTestCosmosMsg(t, types.MsgLock)

	// Create MsgLock attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgAttributes(t, types.MsgLock)
	msgLock, err := BurnLockEventToCosmosMsg(types.MsgLock, cosmosMsgAttributes, sugaredLogger)

	require.Nil(t, err)
	require.Equal(t, expectedMsgLock, msgLock)
}

func TestFailedBurnEventToCosmosMsg(t *testing.T) {
	// Create MsgBurn attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgIncompleteAttributes(t, types.MsgBurn)
	_, err := BurnLockEventToCosmosMsg(types.MsgBurn, cosmosMsgAttributes, sugaredLogger)

	require.Error(t, err)
}

func TestFailedLockEventToCosmosMsg(t *testing.T) {
	// Create MsgLock attributes as input parameter
	cosmosMsgAttributes := CreateCosmosMsgIncompleteAttributes(t, types.MsgLock)
	_, err := BurnLockEventToCosmosMsg(types.MsgLock, cosmosMsgAttributes, sugaredLogger)

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
