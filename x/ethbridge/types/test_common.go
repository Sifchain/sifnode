package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestNetworkDescriptor     = oracletypes.NetworkDescriptor(1)
	TestBridgeContractAddress = "0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB"
	TestAddress               = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestValidator             = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestNonce                 = 0
	TestTokenContractAddress  = "0x0000000000000000000000000000000000000000"
	TestEthereumAddress       = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	AltTestEthereumAddress    = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207344"
	Alt2TestEthereumAddress   = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207333"
	TestCoinsSymbol           = "eth"
	TestCrossChainFeeSymbol   = "ceth"
	AltTestCoinsAmount        = 12
	AltTestCoinsSymbol        = "eth"
	TestCoinIntAmount         = 10
	TestProphecyID            = "test_prophecy_id"
)

var testcrossChainFee = sdk.NewInt(65000000000 * 300000)
var TestCoinsAmount = sdk.NewInt(10)
var AltTestCoinsAmountSDKInt = sdk.NewInt(12)

//Ethereum-bridge specific stuff
func CreateTestEthMsg(t *testing.T, validatorAddress sdk.ValAddress, claimType ClaimType) MsgCreateEthBridgeClaim {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	testContractAddress := NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := NewEthereumAddress(TestTokenContractAddress)
	ethClaim := CreateTestEthClaim(
		t, testContractAddress, testTokenAddress, validatorAddress,
		testEthereumAddress, TestCoinsAmount, TestCoinsSymbol, claimType)
	ethMsg := NewMsgCreateEthBridgeClaim(ethClaim)
	return ethMsg
}

func CreateTestEthClaim(
	t *testing.T, testContractAddress EthereumAddress, testTokenAddress EthereumAddress,
	validatorAddress sdk.ValAddress, testEthereumAddress EthereumAddress, amount sdk.Int, symbol string, claimType ClaimType,
) *EthBridgeClaim {
	testCosmosAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err1)
	ethClaim := NewEthBridgeClaim(
		TestNetworkDescriptor, testContractAddress, TestNonce, symbol,
		testTokenAddress, testEthereumAddress, testCosmosAddress, validatorAddress, amount, claimType)
	return ethClaim
}

func CreateTestBurnMsg(t *testing.T, testCosmosSender string, ethereumReceiver EthereumAddress,
	coinsAmount sdk.Int, coinsSymbol string) MsgBurn {
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)
	burnEth := NewMsgBurn(TestNetworkDescriptor, testCosmosAddress, ethereumReceiver, coinsAmount, coinsSymbol, testcrossChainFee)
	return burnEth
}

func CreateTestLockMsg(t *testing.T, testCosmosSender string, ethereumReceiver EthereumAddress,
	coinsAmount sdk.Int, coinsSymbol string) MsgLock {
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)
	lockEth := NewMsgLock(TestNetworkDescriptor, testCosmosAddress, ethereumReceiver, coinsAmount, coinsSymbol, testcrossChainFee)
	return lockEth
}

func CreateTestQueryEthProphecyResponse(t *testing.T, validatorAddress sdk.ValAddress, claimType ClaimType,
) QueryEthProphecyResponse {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	testContractAddress := NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := NewEthereumAddress(TestTokenContractAddress)
	ethBridgeClaim := CreateTestEthClaim(t, testContractAddress, testTokenAddress, validatorAddress,
		testEthereumAddress, TestCoinsAmount, TestCoinsSymbol, claimType)
	ethBridgeClaims := []string{ethBridgeClaim.ValidatorAddress}

	return NewQueryEthProphecyResponse(
		ethBridgeClaim.GetProphecyID(),
		oracletypes.StatusText_STATUS_TEXT_PENDING,
		ethBridgeClaims,
	)
}

func CreateTestUpdateCrossChainFeeReceiverAccountMsg(t *testing.T, testCosmosSender string, testCrossChainFeeReceiverAccount string) MsgUpdateCrossChainFeeReceiverAccount {
	accAddress1, err := sdk.AccAddressFromBech32(testCosmosSender)
	require.NoError(t, err)
	accAddress2, err := sdk.AccAddressFromBech32(testCrossChainFeeReceiverAccount)
	require.NoError(t, err)

	msgUpdateCrossChainFeeReceiverAccount := NewMsgUpdateCrossChainFeeReceiverAccount(accAddress1, accAddress2)
	return msgUpdateCrossChainFeeReceiverAccount
}

func CreateTestRescueCrossChainFeeMsg(t *testing.T, testCosmosSender string, testCrossChainFeeReceiverAccount string, crosschainFeeSymbol string, crosschainFee sdk.Int) MsgRescueCrossChainFee {
	accAddress1, err := sdk.AccAddressFromBech32(testCosmosSender)
	require.NoError(t, err)
	accAddress2, err := sdk.AccAddressFromBech32(testCrossChainFeeReceiverAccount)
	require.NoError(t, err)

	MsgRescueCrossChainFee := NewMsgRescueCrossChainFee(accAddress1, accAddress2, crosschainFeeSymbol, crosschainFee)
	return MsgRescueCrossChainFee
}

func CreateTestUpdateWhiteListValidatorMsg(_ *testing.T, networkDescriptor oracletypes.NetworkDescriptor, sender string, validator string, power uint32) MsgUpdateWhiteListValidator {
	return MsgUpdateWhiteListValidator{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      sender,
		Validator:         validator,
		Power:             power,
	}
}

func CreateTestSetCrossChainFeeMsg(t *testing.T, testCosmosSender string, networkDescriptor oracletypes.NetworkDescriptor, crossChainFee string) MsgSetFeeInfo {
	accAddress, err := sdk.AccAddressFromBech32(testCosmosSender)
	require.NoError(t, err)

	msgSetFeeInfo := NewMsgSetFeeInfo(accAddress, networkDescriptor, crossChainFee)
	return msgSetFeeInfo
}
