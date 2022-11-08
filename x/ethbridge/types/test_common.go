package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TestEthereumChainID       = 3
	TestBridgeContractAddress = "0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB"
	TestAddress               = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestValidator             = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestNonce                 = 0
	TestTokenContractAddress  = "0x0000000000000000000000000000000000000000"
	TestEthereumAddress       = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	AltTestEthereumAddress    = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207344"
	Alt2TestEthereumAddress   = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207333"
	TestCoinsSymbol           = "eth"
	TestCoinsLockedSymbol     = "ceth"
	AltTestCoinsAmount        = 12
	AltTestCoinsSymbol        = "eth"
	TestCoinIntAmount         = 10
)

var testCethAmount = sdk.NewInt(65000000000 * 300000)
var TestCoinsAmount = sdk.NewInt(10)
var AltTestCoinsAmountSDKInt = sdk.NewInt(12)

//Ethereum-bridge specific stuff
func CreateTestEthMsg(t *testing.T, validatorAddress sdk.ValAddress, claimType ClaimType) MsgCreateEthBridgeClaim {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	testContractAddress := NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := NewEthereumAddress(TestTokenContractAddress)
	ethClaim := CreateTestEthClaim(t, testContractAddress, testTokenAddress, validatorAddress,
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
	ethClaim := NewEthBridgeClaim(TestEthereumChainID, testContractAddress, TestNonce, symbol,
		testTokenAddress, testEthereumAddress, testCosmosAddress, validatorAddress, amount, claimType)
	return ethClaim
}

func CreateTestBurnMsg(t *testing.T, testCosmosSender string, ethereumReceiver EthereumAddress,
	coinsAmount sdk.Int, coinsSymbol string) MsgBurn {
	testCosmosAddress, err := sdk.AccAddressFromBech32(testCosmosSender)
	require.NoError(t, err)
	burnEth := NewMsgBurn(TestEthereumChainID, testCosmosAddress, ethereumReceiver, coinsAmount, coinsSymbol, testCethAmount)
	return burnEth
}

func CreateTestLockMsg(t *testing.T, testCosmosSender string, ethereumReceiver EthereumAddress,
	coinsAmount sdk.Int, coinsSymbol string) MsgLock {
	testCosmosAddress, err := sdk.AccAddressFromBech32(testCosmosSender)
	require.NoError(t, err)
	lockEth := NewMsgLock(TestEthereumChainID, testCosmosAddress, ethereumReceiver, coinsAmount, coinsSymbol, testCethAmount)
	return lockEth
}

func CreateTestQueryEthProphecyResponse(t *testing.T, validatorAddress sdk.ValAddress, claimType ClaimType,
) QueryEthProphecyResponse {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	testContractAddress := NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := NewEthereumAddress(TestTokenContractAddress)
	ethBridgeClaim := CreateTestEthClaim(t, testContractAddress, testTokenAddress, validatorAddress,
		testEthereumAddress, TestCoinsAmount, TestCoinsSymbol, claimType)
	oracleClaim, _ := CreateOracleClaimFromEthClaim(ethBridgeClaim)
	ethBridgeClaims := []*EthBridgeClaim{ethBridgeClaim}
	return NewQueryEthProphecyResponse(
		oracleClaim.Id,
		oracletypes.NewStatus(oracletypes.StatusText_STATUS_TEXT_PENDING, ""),
		ethBridgeClaims,
	)
}

func CreateTestUpdateCethReceiverAccountMsg(t *testing.T, testCosmosSender string, testCethReceiverAccount string) MsgUpdateCethReceiverAccount {
	accAddress1, err := sdk.AccAddressFromBech32(testCosmosSender)
	require.NoError(t, err)
	accAddress2, err := sdk.AccAddressFromBech32(testCethReceiverAccount)
	require.NoError(t, err)
	msgUpdateCethReceiverAccount := NewMsgUpdateCethReceiverAccount(accAddress1, accAddress2)
	return msgUpdateCethReceiverAccount
}

func CreateTestRescueCethMsg(t *testing.T, testCosmosSender string, testCethReceiverAccount string, cethAmount sdk.Int) MsgRescueCeth {
	accAddress1, err := sdk.AccAddressFromBech32(testCosmosSender)
	require.NoError(t, err)
	accAddress2, err := sdk.AccAddressFromBech32(testCethReceiverAccount)
	require.NoError(t, err)
	MsgRescueCeth := NewMsgRescueCeth(accAddress1, accAddress2, cethAmount)
	return MsgRescueCeth
}

func CreateTestUpdateWhiteListValidatorMsg(_ *testing.T, sender string, validator string, operation string) MsgUpdateWhiteListValidator {
	return MsgUpdateWhiteListValidator{
		CosmosSender:  sender,
		Validator:     validator,
		OperationType: operation,
	}
}
