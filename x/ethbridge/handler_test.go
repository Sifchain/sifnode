package ethbridge_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge"
	ethbridgekeeper "github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oraclekeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
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
	ctx, _, _, _, handler, validatorAddresses, _ := CreateTestHandler(t, 0.7, []int64{3, 7})
	valAddress := validatorAddresses[0]
	//Unrecognized type
	res, err := handler(ctx, testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized ethbridge message type: "))
	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.ClaimType_CLAIM_TYPE_LOCK)
	res, err = handler(ctx, &normalCreateMsg)
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
				require.Equal(t, value, oracletypes.StatusText_STATUS_TEXT_PENDING.String())
			case "claim_type":
				require.Equal(t, value, types.ClaimType_CLAIM_TYPE_LOCK.String())
			case "cosmos_sender":
				require.Equal(t, value, valAddress.String())
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	//Bad Creation
	badCreateMsg := types.CreateTestEthMsg(t, valAddress, types.ClaimType_CLAIM_TYPE_LOCK)
	badCreateMsg.EthBridgeClaim.Nonce = -1
	err = badCreateMsg.ValidateBasic()
	require.Error(t, err)
}

func TestDuplicateMsgs(t *testing.T) {
	ctx, _, _, _, handler, validatorAddresses, _ := CreateTestHandler(t, 0.7, []int64{3, 7})
	valAddress := validatorAddresses[0]
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.ClaimType_CLAIM_TYPE_LOCK)
	res, err := handler(ctx, &normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracletypes.StatusText_STATUS_TEXT_PENDING.String())
			}
		}
	}
	//Duplicate message from same validator
	res, err = handler(ctx, &normalCreateMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))
}

func TestMintSuccess(t *testing.T) {
	//Setup
	ctx, _, bankKeeper, _, handler, validatorAddresses, _ := CreateTestHandler(t, 0.7, []int64{2, 7, 1})
	valAddressVal1Pow2 := validatorAddresses[0]
	valAddressVal2Pow7 := validatorAddresses[1]
	valAddressVal3Pow1 := validatorAddresses[2]
	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, valAddressVal1Pow2, types.ClaimType_CLAIM_TYPE_LOCK)
	res, err := handler(ctx, &normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	//Message from second validator succeeds and mints new tokens
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal2Pow7, types.ClaimType_CLAIM_TYPE_LOCK)
	res, err = handler(ctx, &normalCreateMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetAllBalances(ctx, receiverAddress)
	expectedCoins := sdk.NewCoins(sdk.NewInt64Coin(types.TestCoinsLockedSymbol, types.TestCoinIntAmount))
	require.True(t, receiverCoins.IsEqual(expectedCoins))
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracletypes.StatusText_STATUS_TEXT_SUCCESS.String())
			}
		}
	}
	//Additional message from third validator fails and does not mint
	normalCreateMsg = types.CreateTestEthMsg(t, valAddressVal3Pow1, types.ClaimType_CLAIM_TYPE_LOCK)
	res, err = handler(ctx, &normalCreateMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "prophecy already finalized"))
	receiverCoins = bankKeeper.GetAllBalances(ctx, receiverAddress)
	expectedCoins = sdk.NewCoins(sdk.NewInt64Coin(types.TestCoinsLockedSymbol, types.TestCoinIntAmount))
	require.True(t, receiverCoins.IsEqual(expectedCoins))
}

func TestNoMintFail(t *testing.T) {
	//Setup
	ctx, _, bankKeeper, _, handler, validatorAddresses, _ := CreateTestHandler(t, 0.71, []int64{3, 4, 3})
	valAddressVal1Pow3 := validatorAddresses[0]
	valAddressVal2Pow4 := validatorAddresses[1]
	valAddressVal3Pow3 := validatorAddresses[2]
	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)
	ethClaim1 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow3, testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK)
	ethMsg1 := types.NewMsgCreateEthBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal2Pow4, testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK)
	ethMsg2 := types.NewMsgCreateEthBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal3Pow3, testEthereumAddress, types.AltTestCoinsAmountSDKInt, types.AltTestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK)
	ethMsg3 := types.NewMsgCreateEthBridgeClaim(ethClaim3)
	//Initial message
	res, err := handler(ctx, &ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracletypes.StatusText_STATUS_TEXT_PENDING.String())
			}
		}
	}
	//Different message from second validator succeeds
	res, err = handler(ctx, &ethMsg2)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracletypes.StatusText_STATUS_TEXT_PENDING.String())
			}
		}
	}
	//Different message from third validator succeeds but results in failed prophecy with no minting
	res, err = handler(ctx, &ethMsg3)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			if string(attribute.Key) == statusString {
				require.Equal(t, value, oracletypes.StatusText_STATUS_TEXT_FAILED.String())
			}
		}
	}
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiver1Coins := bankKeeper.GetAllBalances(ctx, receiverAddress)
	require.True(t, receiver1Coins.IsZero())
}

func TestLockFail(t *testing.T) {
	//Setup
	ctx, _, _, _, handler, _, _ := CreateTestHandler(t, 0.7, []int64{2, 7, 1})
	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, UnregisteredValidatorAddress, types.ClaimType_CLAIM_TYPE_LOCK)
	res, err := handler(ctx, &normalCreateMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.Equal(t, err.Error(), "validator must be in whitelist")
}

func TestBurnFail(t *testing.T) {
	//Setup
	ctx, _, _, _, handler, _, _ := CreateTestHandler(t, 0.7, []int64{2, 7, 1})
	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, UnregisteredValidatorAddress, types.ClaimType_CLAIM_TYPE_BURN)
	res, err := handler(ctx, &normalCreateMsg)
	require.Error(t, err)
	require.Nil(t, res)
	require.Equal(t, err.Error(), "validator must be in whitelist")
}

func TestBurnEthFail(t *testing.T) {

}

func TestBurnEthSuccess(t *testing.T) {
	ctx, _, bankKeeper, _, handler, validatorAddresses, _ := CreateTestHandler(t, 0.5, []int64{5})
	valAddressVal1Pow5 := validatorAddresses[0]
	senderSequence := "0"
	coinsToMintAmount := sdk.NewInt(7)
	coinsToMintSymbol := "ether"
	coinsToMintSymbolLocked := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, coinsToMintSymbol)
	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)
	ethereumReceiver := types.NewEthereumAddress(types.AltTestEthereumAddress)
	// Initial message to mint some eth
	ethClaim1 := types.CreateTestEthClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal1Pow5,
		testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.ClaimType_CLAIM_TYPE_LOCK)
	ethMsg1 := types.NewMsgCreateEthBridgeClaim(ethClaim1)
	res, err := handler(ctx, &ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetAllBalances(ctx, receiverAddress)
	mintedCoins := sdk.NewCoins(sdk.NewCoin(coinsToMintSymbolLocked, coinsToMintAmount))
	require.True(t, receiverCoins.IsEqual(mintedCoins))

	// Initial message succeeds and mints eth
	coinsToMintAmount = sdk.NewInt(65000000000 * 300000)
	coinsToMintSymbol = "eth"
	testEthereumAddress = types.NewEthereumAddress(types.AltTestEthereumAddress)
	ethClaim1 = types.CreateTestEthClaim(t, testEthereumAddress, testTokenContractAddress, valAddressVal1Pow5,
		testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.ClaimType_CLAIM_TYPE_LOCK)
	ethMsg1 = types.NewMsgCreateEthBridgeClaim(ethClaim1)
	res, err = handler(ctx, &ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	// Second message succeeds, burns eth and fires correct event
	coinsToBurnAmount := sdk.NewInt(3)
	coinsToBurnSymbol := "ether"
	coinsToBurnSymbolPrefixed := fmt.Sprintf("%v%v", types.PeggedCoinPrefix, coinsToBurnSymbol)
	burnMsg := types.CreateTestBurnMsg(t, types.TestAddress, ethereumReceiver, coinsToBurnAmount, coinsToBurnSymbolPrefixed)
	res, err = handler(ctx, &burnMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
	mintedCoins = sdk.NewCoins(mintedCoins[0], sdk.NewCoin(fmt.Sprintf("%v%v", types.PeggedCoinPrefix, coinsToMintSymbol), coinsToMintAmount))
	burnedCoins := sdk.NewCoins(sdk.NewCoin(coinsToBurnSymbolPrefixed, coinsToBurnAmount))
	remainingCoins := mintedCoins.Sub(burnedCoins)
	senderAddress := receiverAddress
	senderCoins := bankKeeper.GetAllBalances(ctx, senderAddress)
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
				// multiple recipients, skip the comparison
				// require.Equal(t, value, senderAddress.String())
				_, err = sdk.AccAddressFromBech32(value)
				require.NoError(t, err)
			case "recipient":
				// multiple recipients, skip the comparison
				// require.Equal(t, value, TestAddress)
				_, err = sdk.AccAddressFromBech32(value)
				require.NoError(t, err)
			case "spender":
				// multiple recipients, skip the comparison
				// require.Equal(t, value, types.TestAddress)
				_, err = sdk.AccAddressFromBech32(value)
				require.NoError(t, err)
			case "receiver":
				// multiple recipients, skip the comparison
				// require.Equal(t, value, types.TestBridgeModuleSif)
				_, err = sdk.AccAddressFromBech32(value)
				require.NoError(t, err)
			case "burner":
				// multiple recipients, skip the comparison
				// require.Equal(t, value, types.TestBridgeModuleSif)
				_, err = sdk.AccAddressFromBech32(value)
				require.NoError(t, err)
			case moduleString:
				require.Equal(t, value, types.ModuleName)
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
	ethClaim1 = types.CreateTestEthClaim(t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow5, testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.ClaimType_CLAIM_TYPE_LOCK)
	ethMsg1 = types.NewMsgCreateEthBridgeClaim(ethClaim1)
	// Initial message succeeds and mints eth
	res, err = handler(ctx, &ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	// Third message failed since pegged token can be lock.
	lockMsg := types.CreateTestLockMsg(t, types.TestAddress, ethereumReceiver, coinsToBurnAmount, coinsToBurnSymbolPrefixed)
	_, err = handler(ctx, &lockMsg)
	require.NotNil(t, err)
	require.Equal(t, "pegged token cether can't be locked", err.Error())
	// Fourth message OK
	_, err = handler(ctx, &burnMsg)
	require.Nil(t, err)
	// Fifth message fails, not enough eth
	res, err = handler(ctx, &burnMsg)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestUpdateCethReceiverAccountMsg(t *testing.T) {
	ctx, _, bankKeeper, accountKeeper, handler, _, oracleKeeper := CreateTestHandler(t, 0.5, []int64{5})
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000)))
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(cosmosSender))
	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	err = sifapp.AddCoinsToAccount(types.ModuleName, bankKeeper, ctx, cosmosSender, coins)
	require.NoError(t, err)
	testUpdateCethReceiverAccountMsg := types.CreateTestUpdateCethReceiverAccountMsg(t, types.TestAddress, types.TestAddress)
	res, err := handler(ctx, &testUpdateCethReceiverAccountMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestRescueCethMsg(t *testing.T) {
	ctx, _, bankKeeper, accountKeeper, handler, _, oracleKeeper := CreateTestHandler(t, 0.5, []int64{5})
	coins := sdk.NewCoins(sdk.NewCoin(types.CethSymbol, sdk.NewInt(10000)))
	err := bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)
	testRescueCethMsg := types.CreateTestRescueCethMsg(t, types.TestAddress, types.TestAddress, sdk.NewInt(10000))
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(cosmosSender))
	_, err = handler(ctx, &testRescueCethMsg)
	require.Error(t, err)
	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	err = sifapp.AddCoinsToAccount(types.ModuleName, bankKeeper, ctx, cosmosSender, coins)
	require.NoError(t, err)
	res, err := handler(ctx, &testRescueCethMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestUpdateWhiteListValidator(t *testing.T) {
	addrs, validatorAddresses := test.CreateTestAddrs(3)
	testCases := []struct {
		name     string
		sender   sdk.AccAddress
		expected []sdk.ValAddress
		msgs     []types.MsgUpdateWhiteListValidator
	}{
		{
			name:     "Add one",
			sender:   addrs[0],
			expected: []sdk.ValAddress{validatorAddresses[0]},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "add"),
			},
		},
		{
			name:     "Add two",
			sender:   addrs[0],
			expected: []sdk.ValAddress{validatorAddresses[0], validatorAddresses[1]},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "add"),
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[1].String(), "add"),
			},
		},
		{
			name:     "Add two remove last",
			sender:   addrs[0],
			expected: []sdk.ValAddress{validatorAddresses[0]},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "add"),
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[1].String(), "add"),
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[1].String(), "remove"),
			},
		},
		{
			name:     "Add two remove first",
			sender:   addrs[0],
			expected: []sdk.ValAddress{validatorAddresses[1]},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "add"),
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[1].String(), "add"),
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "remove"),
			},
		},
		{
			name:     "Remove when none",
			sender:   addrs[0],
			expected: []sdk.ValAddress{},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "remove"),
			},
		},
		{
			name:     "Add one and remove all",
			sender:   addrs[0],
			expected: []sdk.ValAddress{},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "add"),
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "remove"),
			},
		},
		{
			name:     "Add duplicate, remove all instances",
			sender:   addrs[0],
			expected: []sdk.ValAddress{},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "add"),
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "add"),
				types.CreateTestUpdateWhiteListValidatorMsg(t, addrs[0].String(), validatorAddresses[0].String(), "remove"),
			},
		},
	}
	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			ctx, _, _, accountKeeper, handler, _, oracleKeeper := CreateTestHandler(t, 0.5, []int64{5})
			sender := testCase.sender
			accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(sender))
			oracleKeeper.SetAdminAccount(ctx, sender)
			oracleKeeper.SetOracleWhiteList(ctx, []sdk.ValAddress{})
			for i := range testCase.msgs {
				msg := testCase.msgs[i]
				_, err := handler(ctx, &msg)
				require.NoError(t, err)
			}
			wl := oracleKeeper.GetOracleWhiteList(ctx)
			require.Equal(t, testCase.expected, wl)
		})
	}
}

func CreateTestHandler(t *testing.T, consensusNeeded float64, validatorAmounts []int64) (sdk.Context, ethbridgekeeper.Keeper,
	bankkeeper.Keeper, authkeeper.AccountKeeper, sdk.Handler, []sdk.ValAddress, oraclekeeper.Keeper) {
	ctx, keeper, bankKeeper, accountKeeper, oracleKeeper, _, validatorAddresses := test.CreateTestKeepers(t, consensusNeeded, validatorAmounts, "")
	cethReceiverAccount, _ := sdk.AccAddressFromBech32(types.TestAddress)
	keeper.SetCethReceiverAccount(ctx, cethReceiverAccount)
	handler := ethbridge.NewHandler(keeper)
	return ctx, keeper, bankKeeper, accountKeeper, handler, validatorAddresses, oracleKeeper
}
