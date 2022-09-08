package ethbridge_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

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
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

const (
	moduleString = "module"
	statusString = "status"
	senderString = "sender"
	power        = 100
	zeroPower    = 0
)

var (
	UnregisteredValidatorAddress = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestAccAddress               = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestAddress                  = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
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
			case "ethereum_sender_sequence":
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
	require.Equal(t, err.Error(), oracletypes.ErrDuplicateMessage.Error())
}

func TestMintSuccess(t *testing.T) {
	//Setup
	ctx, keeper, bankKeeper, _, handler, validatorAddresses, _ := CreateTestHandler(t, 0.7, []int64{2, 7, 1})

	valAddressVal1Pow2 := validatorAddresses[0]
	valAddressVal2Pow7 := validatorAddresses[1]
	valAddressVal3Pow1 := validatorAddresses[2]

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, valAddressVal1Pow2, types.ClaimType_CLAIM_TYPE_LOCK)

	entry := tokenregistrytypes.RegistryEntry{
		Denom:         normalCreateMsg.EthBridgeClaim.Denom,
		DisplaySymbol: normalCreateMsg.EthBridgeClaim.Symbol,
		Decimals:      18,
		Permissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
	}
	keeper.GetTokenRegistryKeeper().SetToken(ctx, &entry)

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

	expectedCoins := sdk.NewCoins(sdk.NewInt64Coin(normalCreateMsg.EthBridgeClaim.Denom, types.TestCoinIntAmount))
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
	_, err = handler(ctx, &normalCreateMsg)
	require.Nil(t, err)
	receiverCoins = bankKeeper.GetAllBalances(ctx, receiverAddress)

	expectedCoins = sdk.NewCoins(sdk.NewInt64Coin(normalCreateMsg.EthBridgeClaim.Denom, types.TestCoinIntAmount))
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
		valAddressVal1Pow3, testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK, types.TestDecimals, types.TestName)
	ethMsg1 := types.NewMsgCreateEthBridgeClaim(ethClaim1)
	ethClaim2 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal2Pow4, testEthereumAddress, types.TestCoinsAmount, types.TestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK, types.TestDecimals, types.TestName)
	ethMsg2 := types.NewMsgCreateEthBridgeClaim(ethClaim2)
	ethClaim3 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal3Pow3, testEthereumAddress, types.AltTestCoinsAmountSDKInt, types.AltTestCoinsSymbol, types.ClaimType_CLAIM_TYPE_LOCK, types.TestDecimals, types.TestName)
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
				require.Equal(t, value, oracletypes.StatusText_STATUS_TEXT_PENDING.String())
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

	addAddress, err := sdk.AccAddressFromBech32(UnregisteredValidatorAddress)
	require.NoError(t, err)
	valAddress := sdk.ValAddress(addAddress)

	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.ClaimType_CLAIM_TYPE_LOCK)
	res, err := handler(ctx, &normalCreateMsg)

	require.Error(t, err)
	require.Nil(t, res)
	require.Equal(t, err.Error(), "validator not in white list")
}

func TestBurnFail(t *testing.T) {
	//Setup
	ctx, _, _, _, handler, _, _ := CreateTestHandler(t, 0.7, []int64{2, 7, 1})
	addAddress, err := sdk.AccAddressFromBech32(UnregisteredValidatorAddress)
	require.NoError(t, err)

	valAddress := sdk.ValAddress(addAddress)
	//Initial message
	normalCreateMsg := types.CreateTestEthMsg(t, valAddress, types.ClaimType_CLAIM_TYPE_BURN)
	res, err := handler(ctx, &normalCreateMsg)

	require.Error(t, err)
	require.Nil(t, res)
	require.Equal(t, err.Error(), "validator not in white list")
}

func TestBurnEthFail(t *testing.T) {

}

func TestBurnEthSuccess(t *testing.T) {
	ctx, keeper, bankKeeper, _, handler, validatorAddresses, _ := CreateTestHandler(t, 0.5, []int64{5})
	valAddressVal1Pow5 := validatorAddresses[0]

	// Initial message to mint some eth
	coinsToMintAmount := sdk.NewInt(7)
	coinsToMintSymbol := "eth"

	testTokenContractAddress := types.NewEthereumAddress(types.TestTokenContractAddress)
	testEthereumAddress := types.NewEthereumAddress(types.TestEthereumAddress)

	ethClaim1 := types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow5, testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.ClaimType_CLAIM_TYPE_LOCK, types.TestDecimals, types.TestName)
	ethMsg1 := types.NewMsgCreateEthBridgeClaim(ethClaim1)

	denomHash := ethClaim1.Denom

	entry := tokenregistrytypes.RegistryEntry{
		Denom:         denomHash,
		DisplaySymbol: ethMsg1.EthBridgeClaim.Symbol,
		Decimals:      18,
		Network:       oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM,
		Permissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
	}
	keeper.GetTokenRegistryKeeper().SetToken(ctx, &entry)

	// Initial message succeeds and mints eth
	res, err := handler(ctx, &ethMsg1)
	require.NoError(t, err)
	require.NotNil(t, res)
	receiverAddress, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	receiverCoins := bankKeeper.GetAllBalances(ctx, receiverAddress)
	mintedCoins := sdk.NewCoins(sdk.NewCoin(ethMsg1.EthBridgeClaim.Denom, coinsToMintAmount))

	require.True(t, receiverCoins.IsEqual(mintedCoins))

	coinsToBurnAmount := sdk.NewInt(3)
	ethereumReceiver := types.NewEthereumAddress(types.AltTestEthereumAddress)

	// Second message succeeds, burns eth and fires correct event
	burnMsg := types.CreateTestBurnMsg(t, types.TestAddress, ethereumReceiver, coinsToBurnAmount, denomHash)

	res, err = handler(ctx, &burnMsg)
	require.NoError(t, err)
	require.NotNil(t, res)

	networkDescriptor := ""
	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case "prophecy_id":
			case "global_sequence":
			case senderString:
			case "recipient":

			case moduleString:
				require.Equal(t, value, types.ModuleName)
			case "network_id":
				networkDescriptor = value
			case "amount":
			case "spender":
			case "receiver":
			case "burner":
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	TestNetworkDescriptorStr := strconv.Itoa(int(types.TestNetworkDescriptor))
	require.Equal(t, networkDescriptor, TestNetworkDescriptorStr)

	coinsToMintAmount = sdk.NewInt(65000000000 * 300000)
	coinsToMintSymbol = "eth"
	testEthereumAddress = types.NewEthereumAddress(types.Alt2TestEthereumAddress)

	ethClaim1 = types.CreateTestEthClaim(
		t, testEthereumAddress, testTokenContractAddress,
		valAddressVal1Pow5, testEthereumAddress, coinsToMintAmount, coinsToMintSymbol, types.ClaimType_CLAIM_TYPE_LOCK, types.TestDecimals, types.TestName)
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

func TestUpdateCrossChainFeeReceiverAccountMsg(t *testing.T) {
	ctx, _, bankKeeper, accountKeeper, handler, _, oracleKeeper := CreateTestHandler(t, 0.5, []int64{5})
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000)))

	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(cosmosSender))
	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	err = bankKeeper.MintCoins(ctx, ethbridge.ModuleName, coins)
	require.NoError(t, err)
	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, ethbridge.ModuleName, cosmosSender, coins)
	require.NoError(t, err)

	testUpdateCrossChainFeeReceiverAccountMsg := types.CreateTestUpdateCrossChainFeeReceiverAccountMsg(
		t, types.TestAddress, types.TestAddress)

	res, err := handler(ctx, &testUpdateCrossChainFeeReceiverAccountMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestRescueCrossChainFeeMsg(t *testing.T) {
	ctx, _, bankKeeper, accountKeeper, handler, _, oracleKeeper := CreateTestHandler(t, 0.5, []int64{5})
	coins := sdk.NewCoins(sdk.NewCoin("ceth", sdk.NewInt(10000)))
	err := bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	require.NoError(t, err)

	testRescueCrossChainFeeMsg := types.CreateTestRescueCrossChainFeeMsg(
		t, types.TestAddress, types.TestAddress, types.TestCrossChainFeeSymbol, sdk.NewInt(10000))

	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)

	accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(cosmosSender))

	_, err = handler(ctx, &testRescueCrossChainFeeMsg)
	require.Error(t, err)

	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	err = bankKeeper.MintCoins(ctx, ethbridge.ModuleName, coins)
	require.NoError(t, err)

	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, ethbridge.ModuleName, cosmosSender, coins)
	require.NoError(t, err)

	res, err := handler(ctx, &testRescueCrossChainFeeMsg)
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
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), power),
			},
		},
		{
			name:     "Add two",
			sender:   addrs[0],
			expected: []sdk.ValAddress{validatorAddresses[0], validatorAddresses[1]},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), power),
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[1].String(), power),
			},
		},
		{
			name:     "Add two remove last",
			sender:   addrs[0],
			expected: []sdk.ValAddress{validatorAddresses[0]},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), power),
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[1].String(), power),
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[1].String(), zeroPower),
			},
		},
		{
			name:     "Add two remove first",
			sender:   addrs[0],
			expected: []sdk.ValAddress{validatorAddresses[1]},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), power),
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[1].String(), power),
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), zeroPower),
			},
		},
		{
			name:     "Remove when none",
			sender:   addrs[0],
			expected: []sdk.ValAddress{},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), zeroPower),
			},
		},
		{
			name:     "Add one and remove all",
			sender:   addrs[0],
			expected: []sdk.ValAddress{},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), power),
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), zeroPower),
			},
		},
		{
			name:     "Add duplicate, remove all instances",
			sender:   addrs[0],
			expected: []sdk.ValAddress{},
			msgs: []types.MsgUpdateWhiteListValidator{
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), power),
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), power),
				types.CreateTestUpdateWhiteListValidatorMsg(t, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, addrs[0].String(), validatorAddresses[0].String(), zeroPower),
			},
		},
	}

	networkDescriptor := oracletypes.NewNetworkIdentity(oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			ctx, _, _, accountKeeper, handler, _, oracleKeeper := CreateTestHandler(t, 0.5, []int64{5})
			sender := testCase.sender

			accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(sender))
			oracleKeeper.SetAdminAccount(ctx, sender)
			oracleKeeper.RemoveOracleWhiteList(ctx, networkDescriptor)

			for i := range testCase.msgs {
				msg := testCase.msgs[i]
				_, err := handler(ctx, &msg)
				require.NoError(t, err)
			}

			wl := oracleKeeper.GetAllValidators(ctx, networkDescriptor)
			for _, address := range wl {
				found := false
				for _, expected := range testCase.expected {
					if address.Equals(expected) {
						found = true
					}
				}
				require.Equal(t, found, true)
			}
		})
	}
}

func TestSetCrossChainFeeMsg(t *testing.T) {
	ctx, _, _, accountKeeper, handler, _, oracleKeeper := CreateTestHandler(t, 0.5, []int64{5})
	feeCurrencyGas := sdk.NewInt(1)
	minimumLockCost := sdk.NewInt(1)
	minimumBurnCost := sdk.NewInt(1)

	testSetAtiveTokenMsg := types.CreateTestSetCrossChainFeeMsg(
		t, types.TestAddress, oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM, "ceth",
		feeCurrencyGas, minimumLockCost, minimumBurnCost)

	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)

	accountKeeper.SetAccount(ctx, authtypes.NewBaseAccountWithAddress(cosmosSender))

	_, err = handler(ctx, &testSetAtiveTokenMsg)
	require.Error(t, err)

	oracleKeeper.SetAdminAccount(ctx, cosmosSender)

	res, err := handler(ctx, &testSetAtiveTokenMsg)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func CreateTestHandler(t *testing.T, consensusNeeded float64, validatorAmounts []int64) (sdk.Context,
	ethbridgekeeper.Keeper, bankkeeper.Keeper, authkeeper.AccountKeeper,
	sdk.Handler, []sdk.ValAddress, oraclekeeper.Keeper) {

	ctx, keeper, bankKeeper, accountKeeper, oracleKeeper, _, _, validators := test.CreateTestKeepers(t, consensusNeeded, validatorAmounts, "")

	CrossChainFeeReceiverAccount, _ := sdk.AccAddressFromBech32(TestAddress)
	keeper.SetCrossChainFeeReceiverAccount(ctx, CrossChainFeeReceiverAccount)
	handler := ethbridge.NewHandler(keeper)

	return ctx, keeper, bankKeeper, accountKeeper, handler, validators, oracleKeeper
}
