package types

import (
	"testing"
	"github.com/Sifchain/sifnode/simapp"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/oracle"
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
	AddressKey1 							= "A58856F0FD53BF058B4909A21AEC019107BA6"

)


var testCethAmount = sdk.NewInt(65000000000 * 300000)
var TestCoinsAmount = sdk.NewInt(10)
var AltTestCoinsAmountSDKInt = sdk.NewInt(12)

//// returns context and app with params set on account keeper
func CreateTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))
	_ = simapp.AddTestAddrs(app, ctx, 6, initTokens)
	return app, ctx
}

func CreateTestAppEthBridge(isCheckTx bool) (sdk.Context, keeper.Keeper) {
	ctx, app := GetSimApp(isCheckTx)
	return ctx, app.EthBridgeKeeper
}

func GetSimApp(isCheckTx bool) (sdk.Context, *simapp.SimApp) {
	app, ctx := CreateTestApp(isCheckTx)
	return ctx, app
}

func GenerateRandomTokens(numberOfTokens int) []string {
	var tokenList []string
	tokens := []string{"ceth", "cbtc", "ceos", "cbch", "cbnb", "cusdt", "cada", "ctrx", "cacoin", "cbcoin", "ccoin", "cdcoin"}
	rand.Seed(time.Now().Unix())
	for i := 0; i < numberOfTokens; i++ {
		// initialize global pseudo random generator
		randToken := tokens[rand.Intn(len(tokens))]

		tokenList = append(tokenList, randToken)
	}
	return tokenList
}

func GenerateAddress(key string) sdk.AccAddress {
	if key == "" {
		key = AddressKey1
	}
	var buffer bytes.Buffer
	buffer.WriteString(key)
	buffer.WriteString(strconv.Itoa(100))
	res, _ := sdk.AccAddressFromHex(buffer.String())
	bech := res.String()
	addr := buffer.String()
	res, err := sdk.AccAddressFromHex(addr)

	if err != nil {
		panic(err)
	}

	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}
	return res
}

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
) EthBridgeClaim {
	testCosmosAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err1)
	ethClaim := NewEthBridgeClaim(
		TestEthereumChainID, testContractAddress, TestNonce, symbol,
		testTokenAddress, testEthereumAddress, testCosmosAddress, validatorAddress, amount, claimType)
	return ethClaim
}

func CreateTestBurnMsg(t *testing.T, testCosmosSender string, ethereumReceiver EthereumAddress,
	coinsAmount sdk.Int, coinsSymbol string) MsgBurn {
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)
	burnEth := NewMsgBurn(TestEthereumChainID, testCosmosAddress, ethereumReceiver, coinsAmount, coinsSymbol, testCethAmount)
	return burnEth
}

func CreateTestLockMsg(t *testing.T, testCosmosSender string, ethereumReceiver EthereumAddress,
	coinsAmount sdk.Int, coinsSymbol string) MsgLock {
	testCosmosAddress, err := sdk.AccAddressFromBech32(TestAddress)
	require.NoError(t, err)
	lockEth := NewMsgLock(TestEthereumChainID, testCosmosAddress, ethereumReceiver, coinsAmount, coinsSymbol, testCethAmount)
	return lockEth
}

func CreateTestQueryEthProphecyResponse(
	cdc *codec.Codec, t *testing.T, validatorAddress sdk.ValAddress, claimType ClaimType,
) QueryEthProphecyResponse {
	testEthereumAddress := NewEthereumAddress(TestEthereumAddress)
	testContractAddress := NewEthereumAddress(TestBridgeContractAddress)
	testTokenAddress := NewEthereumAddress(TestTokenContractAddress)
	ethBridgeClaim := CreateTestEthClaim(t, testContractAddress, testTokenAddress, validatorAddress,
		testEthereumAddress, TestCoinsAmount, TestCoinsSymbol, claimType)
	oracleClaim, _ := CreateOracleClaimFromEthClaim(cdc, ethBridgeClaim)
	ethBridgeClaims := []EthBridgeClaim{ethBridgeClaim}

	return NewQueryEthProphecyResponse(
		oracleClaim.ID,
		oracle.NewStatus(oracle.PendingStatusText, ""),
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
