package ethbridge

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/oracle"
)

const (
	moduleString = "module"
	statusString = "status"
	senderString = "sender"
)

var (
	UnregisteredValidatorAddress = sdk.ValAddress("cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq")
)

func TestBasicMsgs(t *testing.T) {
	//Setup
	ctx, _, _, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{3, 7})

	valAddress := validatorAddresses[0]

	//Unrecognized type
	res, err := handler(ctx, sdk.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized ethbridge message type: "))

	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.LockText)
	res, err = handler(ctx, normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case "module":
				require.Equal(t, value, types.ModuleName)
			case senderString:
				require.Equal(t, value, valAddress.String())
			case "ethereum_sender":
				require.Equal(t, value, types.TestEthereumAddress)
			case "ethereum_sender_nonce":
				require.Equal(t, value, strconv.Itoa(types.TestNonce))
			case "cosmos_receiver":
				require.Equal(t, value, types.TestAddress)
			case "amount":
				require.Equal(t, value, strconv.FormatInt(10, 10))
			case "symbol":
				require.Equal(t, value, types.TestCoinsSymbol)
			case "token_contract_address":
				require.Equal(t, value, types.TestTokenContractAddress)
			case statusString:
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			case "claim_type":
				require.Equal(t, value, types.ClaimTypeToString[types.LockText])
			case "cosmos_sender":
				require.Equal(t, value, valAddress.String())
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	//Bad Creation
	badCreateMsg := types.CreateTestEthMsg(t, valAddress, types.LockText)
	badCreateMsg.Nonce = -1
	err = badCreateMsg.ValidateBasic()
	require.Error(t, err)
}

func TestDuplicateMsgs(t *testing.T) {
	ctx, _, _, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{3, 7})

	valAddress := validatorAddresses[0]

	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.LockText)
	res, err := handler(ctx, normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Duplicate message from same validator
	res, err = handler(ctx, normalCreateMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))
}

func TestMintSuccess(t *testing.T) {
	//Setup
	ctx, _, bankKeeper, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.7, []int64{2, 7, 1})

	valAddressVal1Pow2 := validatorAddresses[0]
	valAddressVal2Pow7 := validatorAddresses[1]
	valAddressVal3Pow1 := validatorAddresses[2]

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, valAddressVal1Pow2, types.LockText)
	res, err := handler(ctx, normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	//Message from second validator succeeds and mints new tokens
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal2Pow7, types.LockText)
	res, err = handler(ctx, normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetCoins(ctx, receiverAddress)
	expectedCoins := sdk.Coins{sdk.NewInt64Coin(types.TestCoinsLockedSymbol, types.TestCoinIntAmount)}
	require.True(t, receiverCoins.IsEqual(expectedCoins))
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.SuccessStatusText])
			}
		}
	}

	//Additional message from third validator fails and does not mint
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal3Pow1, types.LockText)
	res, err = handler(ctx, normalCreateMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "prophecy already finalized"))
	receiverCoins = bankKeeper.GetCoins(ctx, receiverAddress)
	expectedCoins = sdk.Coins{sdk.NewInt64Coin(types.TestCoinsLockedSymbol, types.TestCoinIntAmount)}
	require.True(t, receiverCoins.IsEqual(expectedCoins))
}

func TestNoMintFail(t *testing.T) {
	//Setup
	ctx, _, bankKeeper, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.71, []int64{3, 4, 3})

	valAddressVal1Pow3 := validatorAddresses[0]
	valAddressVal2Pow4 := validatorAddresses[1]
	valAddressVal3Pow3 := validatorAddresses[2]

	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)

	ethClaim1 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow3, testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.LockText)
	ethMsg1 := NewMsgCreateEthBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal2Pow4, testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.LockText)
	ethMsg2 := NewMsgCreateEthBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal3Pow3, testEthereumAddress, types.AltTestCoinsAmountSDKInt, types.AltTestCoinsSymbol, types.LockText)
	ethMsg3 := NewMsgCreateEthBridgeClaim(ethClaim3)

	//Initial message
	res, err := handler(ctx, ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Different message from second validator succeeds
	res, err = handler(ctx, ethMsg2)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.PendingStatusText])
			}
		}
	}

	//Different message from third validator succeeds but results in failed prophecy with no minting
	res, err = handler(ctx, ethMsg3)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracle.StatusTextToString[oracle.FailedStatusText])
			}
		}
	}
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiver1Coins := bankKeeper.GetCoins(ctx, receiverAddress)
	require.True(t, receiver1Coins.IsZero())
}

func TestLockFail(t *testing.T) {
	//Setup
	ctx, _, _, _, _, _, handler := CreateTestHandler(t, 0.7, []int64{2, 7, 1})

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, UnregisteredValidatorAddress, types.LockText)
	res, err := handler(ctx, normalCreateMsg)

	require.Error(t, err)
	require.Nil(t, res)
	require.Equal(t, err.Error(), "validator must be in whitelist")
}

func TestBurnFail(t *testing.T) {
	//Setup
	ctx, _, _, _, _, _, handler := CreateTestHandler(t, 0.7, []int64{2, 7, 1})

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, UnregisteredValidatorAddress, types.BurnText)
	res, err := handler(ctx, normalCreateMsg)

	require.Error(t, err)
	require.Nil(t, res)
	require.Equal(t, err.Error(), "validator must be in whitelist")
}

func TestBurnEthFail(t *testing.T) {

}

func TestBurnEthSuccess(t *testing.T) {
	ctx, _, bankKeeper, _, _, validatorAddresses, handler := CreateTestHandler(t, 0.5, []int64{5})
	valAddressVal1Pow5 := validatorAddresses[0]

	// Initial message to mint some eth
	coinsToMintAmount := sdk.NewInt(7)
	coinsToMintSymbol := "ether"
	coinsToMintSymbolLocked := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, coinsToMintSymbol)

	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)

	ethClaim1 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow5, testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.LockText)
	ethMsg1 := NewMsgCreateEthBridgeClaim(ethClaim1)

	// Initial message succeeds and mints eth
	res, err := handler(ctx, ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetCoins(ctx, receiverAddress)
	mintedCoins := sdk.Coins{sdk.NewCoin(coinsToMintSymbolLocked, coinsToMintAmount)}
	require.True(t, receiverCoins.IsEqual(mintedCoins))

	coinsToMintAmount = sdk.NewInt(65000000000 * 300000)
	coinsToMintSymbol = "eth"
	testEthereumAddress = types.NewEthereumAddress(types.AltTestEthereumAddress)

	ethClaim1 = types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow5, testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.LockText)
	ethMsg1 = NewMsgCreateEthBridgeClaim(ethClaim1)

	// Initial message succeeds and mints eth
	res, err = handler(ctx, ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	coinsToBurnAmount := sdk.NewInt(3)
	coinsToBurnSymbol := "ether"
	coinsToBurnSymbolPrefixed := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, coinsToBurnSymbol)

	ethereumReceiver := types.NewEthereumAddress(types.AltTestEthereumAddress)

	// Second message succeeds, burns eth and fires correct event
	burnMsg := types.CreateTestBurnMsg(t, types.TestAddress, ethereumReceiver, coinsToBurnAmount,
		coinsToBurnSymbolPrefixed)
	res, err = handler(ctx, burnMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	senderAddress := receiverAddress
	burnedCoins := sdk.Coins{sdk.NewCoin(coinsToBurnSymbolPrefixed, coinsToBurnAmount)}
	senderSequence := "0"
	remainingCoins := mintedCoins.Sub(burnedCoins)
	senderCoins := bankKeeper.GetCoins(ctx, senderAddress)
	require.True(t, senderCoins.IsEqual(remainingCoins))
	eventEthereumChainID := ""
	eventCosmosSender := ""
	eventCosmosSenderSequence := ""
	eventEthereumReceiver := ""
	eventAmount := ""
	eventSymbol := ""
	eventCethAmount := sdk.NewInt(0)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case senderString:
				// multiple recipient in burn, skip the comparison
				// require.Equal(t, value, senderAddress.String())
			case "recipient":
				// multiple recipient in burn, skip the comparison
				// require.Equal(t, value, TestAddress)
			case moduleString:
				require.Equal(t, value, ModuleName)
			case "ethereum_chain_id":
				eventEthereumChainID = value
			case "cosmos_sender":
				eventCosmosSender = value
			case "cosmos_sender_sequence":
				eventCosmosSenderSequence = value
			case "ethereum_receiver":
				eventEthereumReceiver = value
			case "amount":
				eventAmount = value
			case "symbol":
				eventSymbol = value
			case "ceth_amount":
				var ok bool
				eventCethAmount, ok = sdk.NewIntFromString(value)
				require.Equal(t, ok, true)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	require.Equal(t, eventEthereumChainID, strconv.Itoa(types.TestEthereumChainID))
	require.Equal(t, eventCosmosSender, senderAddress.String())
	require.Equal(t, eventCosmosSenderSequence, senderSequence)
	require.Equal(t, eventEthereumReceiver, ethereumReceiver.String())
	require.Equal(t, eventAmount, coinsToBurnAmount.String())
	require.Equal(t, eventSymbol, coinsToBurnSymbolPrefixed)
	require.Equal(t, eventCethAmount, sdk.NewInt(65000000000*300000))

	coinsToMintAmount = sdk.NewInt(65000000000 * 300000)
	coinsToMintSymbol = "eth"
	testEthereumAddress = types.NewEthereumAddress(types.Alt2TestEthereumAddress)

	ethClaim1 = types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow5, testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.LockText)
	ethMsg1 = NewMsgCreateEthBridgeClaim(ethClaim1)

	// Initial message succeeds and mints eth
	res, err = handler(ctx, ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Third message failed since pegged token can be lock.
	lockMsg := types.CreateTestLockMsg(t, types.TestAddress, ethereumReceiver, coinsToBurnAmount,
		coinsToBurnSymbolPrefixed)
	_, err = handler(ctx, lockMsg)
	require.NotNil(t, err)
	require.Equal(t, "Pegged token cether can't be lock.", err.Error())

	// Fourth message OK
	_, err = handler(ctx, burnMsg)
	require.Nil(t, err)

	// Fifth message fails, not enough eth
	res, err = handler(ctx, burnMsg)
	require.Error(t, err)
	require.Nil(t, res)
}
