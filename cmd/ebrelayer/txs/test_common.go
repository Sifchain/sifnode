package txs

import (
	"encoding/binary"
	"math/big"
	"strconv"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	// EthereumPrivateKey config field which holds the user's private key
	EthereumPrivateKey        = "ETHEREUM_PRIVATE_KEY"
	TestNetworkID             = 1
	TestBridgeContractAddress = "0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2"
	TestLockClaimType         = 1
	TestBurnClaimType         = 2
	TestProphecyID            = 20
	TestNonce                 = 19
	TestEthTokenAddress       = "0x0000000000000000000000000000000000000000"
	TestSymbol                = "CETH"
	TestEthereumAddress1      = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359"
	TestEthereumAddress2      = "0xc230f38FF05860753840e0d7cbC66128ad308B67"
	TestCosmosAddress1        = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestCosmosAddress2        = "cosmos1l5h2x255pvdy9l4z0hf9tr8zw7k657s97wyz7y"
	TestCosmosValAddress      = "cosmosvaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pn7fqfk"
	TestExpectedMessage       = "8d46d2f689aa50a0dde8563f4ab1c90f4f74a80817ad18052ef1aa8bd5a0fd96"
	TestCosmosAddressSequence = 1
	TestExpectedSignature     = "f3b43b87b8b3729d6b380a640954d14e425acd603bc98f5da8437cba9e492e7805c437f977900cf9cfbeb9e0e2e6fc5b189723b0979efff1fc2796a2daf4fd3a01" //nolint:lll
	TestAddrHex               = "970e8128ab834e8eac17ab8e3812f010678cf791"
	TestPrivHex               = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	TestNullAddress           = "0x0000000000000000000000000000000000000000"
	TestOtherAddress          = "0x1000000000000000000000000000000000000000"
)

var testAmount = big.NewInt(5)
var testSDKAmount = sdk.NewIntFromBigInt(testAmount)

// CreateTestLogEthereumEvent creates a sample EthereumEvent event for testing purposes
func CreateTestLogEthereumEvent(t *testing.T) types.EthereumEvent {
	networkID := oracletypes.NetworkID(TestNetworkID)
	testBridgeContractAddress := common.HexToAddress(TestBridgeContractAddress)
	// Convert int to [32]byte
	var testProphecyID []byte
	var testProphecyID32 [32]byte
	testProphecyID = make([]byte, 32)
	binary.LittleEndian.PutUint64(testProphecyID, uint64(TestProphecyID))
	copy(testProphecyID32[:], testProphecyID)
	testEthereumSender := common.HexToAddress(TestEthereumAddress1)
	testCosmosRecipient := []byte(TestCosmosAddress1)
	testTokenAddress := common.HexToAddress(TestEthTokenAddress)
	testAmount := testAmount
	testNonce := big.NewInt(int64(TestNonce))

	return types.EthereumEvent{NetworkID: networkID,
		BridgeContractAddress: testBridgeContractAddress,
		ID:                    testProphecyID32,
		From:                  testEthereumSender,
		To:                    testCosmosRecipient, Token: testTokenAddress,
		Symbol: TestSymbol, Value: testAmount, Nonce: testNonce, ClaimType: ethbridge.ClaimType_CLAIM_TYPE_LOCK}
}

// CreateTestProphecyClaimEvent creates a sample ProphecyClaimEvent for testing purposes
func CreateTestProphecyClaimEvent(t *testing.T) types.ProphecyClaimEvent {
	testProphecyID := big.NewInt(int64(TestProphecyID))
	testEthereumReceiver := common.HexToAddress(TestEthereumAddress1)
	testValidatorAddress := common.HexToAddress(TestEthereumAddress2)
	testTokenAddress := common.HexToAddress(TestEthTokenAddress)
	testAmount := testSDKAmount

	return types.NewProphecyClaimEvent([]byte(TestCosmosAddress1), TestSymbol,
		testProphecyID, testAmount, testEthereumReceiver, testValidatorAddress,
		testTokenAddress, TestBurnClaimType)
}

// CreateTestCosmosMsg creates a sample Cosmos Msg for testing purposes
func CreateTestCosmosMsg(t *testing.T, claimType types.Event) types.CosmosMsg {
	testCosmosSender := []byte(TestCosmosAddress1)
	testEthereumReceiver := common.HexToAddress(TestEthereumAddress1)
	testAmount := testSDKAmount

	var symbol string
	if claimType == types.MsgBurn {
		res := strings.SplitAfter(strings.ToLower(TestSymbol), defaultSifchainPrefix)
		symbol = strings.Join(res[1:], "")
	} else {
		symbol = TestSymbol
	}

	// Create new Cosmos Msg
	cosmosMsg := types.NewCosmosMsg(oracletypes.NetworkID(TestNetworkID), claimType, testCosmosSender, big.NewInt(TestCosmosAddressSequence),
		testEthereumReceiver, symbol, testAmount)

	return cosmosMsg
}

// CreateCosmosMsgAttributes creates expected attributes for a MsgBurn/MsgLock for testing purposes
func CreateCosmosMsgAttributes(t *testing.T, claimType types.Event) []abci.EventAttribute {
	attributes := [7]abci.EventAttribute{}

	// (key, value) pairing for "cosmos_sender" key
	pairCosmosSender := abci.EventAttribute{
		Key:   []byte("cosmos_sender"),
		Value: []byte(TestCosmosAddress1),
	}

	// (key, value) pairing for "cosmos_sender_sequence" key
	pairCosmosSenderSequence := abci.EventAttribute{
		Key:   []byte("cosmos_sender_sequence"),
		Value: []byte(strconv.Itoa(TestCosmosAddressSequence)),
	}

	// (key, value) pairing for "ethereum_receiver" key
	pairEthereumReceiver := abci.EventAttribute{
		Key:   []byte("ethereum_receiver"),
		Value: []byte(common.HexToAddress(TestEthereumAddress1).Hex()), // .Bytes() doesn't seem to work here
	}

	// (key, value) pairing for "symbol" key
	var symbol string
	if claimType == types.MsgBurn {
		symbol = strings.ToLower(TestSymbol)
	} else {
		symbol = TestSymbol
	}
	pairSymbol := abci.EventAttribute{
		Key:   []byte("symbol"),
		Value: []byte(symbol),
	}

	// (key, value) pairing for "amount" key
	pairAmount := abci.EventAttribute{
		Key:   []byte("amount"),
		Value: []byte(testAmount.String()),
	}

	// (key, value) pairing for "token_contract_address" key
	pairTokenContract := abci.EventAttribute{
		Key:   []byte("token_contract_address"),
		Value: []byte(common.HexToAddress(TestEthTokenAddress).Hex()),
	}

	// (key, value) pairing for "ethereum_chain_id" key
	pairEthereumChainID := abci.EventAttribute{
		Key:   []byte("ethereum_chain_id"),
		Value: []byte(strconv.Itoa(TestNetworkID)),
	}

	// Assign pairs to attributes array
	attributes[0] = pairCosmosSender
	attributes[1] = pairCosmosSenderSequence
	attributes[2] = pairEthereumReceiver
	attributes[3] = pairTokenContract
	attributes[4] = pairSymbol
	attributes[5] = pairAmount
	attributes[6] = pairEthereumChainID
	return attributes[:]
}

// CreateCosmosMsgIncompleteAttributes creates a MsgBurn/MsgLock for testing purposes missing some attributes
func CreateCosmosMsgIncompleteAttributes(t *testing.T, claimType types.Event) []abci.EventAttribute {
	attributes := [3]abci.EventAttribute{}
	// (key, value) pairing for "cosmos_sender" key
	pairCosmosSender := abci.EventAttribute{
		Key:   []byte("cosmos_sender"),
		Value: []byte(TestCosmosAddress1),
	}

	// (key, value) pairing for "cosmos_sender_sequence" key
	pairCosmosSenderSequence := abci.EventAttribute{
		Key:   []byte("cosmos_sender_sequence"),
		Value: []byte(strconv.Itoa(TestCosmosAddressSequence)),
	}

	// (key, value) pairing for "ethereum_receiver" key
	pairEthereumReceiver := abci.EventAttribute{
		Key:   []byte("ethereum_receiver"),
		Value: []byte(common.HexToAddress(TestEthereumAddress1).Hex()), // .Bytes() doesn't seem to work here
	}

	// Assign pairs to attributes array
	attributes[0] = pairCosmosSender
	attributes[1] = pairCosmosSenderSequence
	attributes[2] = pairEthereumReceiver

	return attributes[:]
}

// CreateEthereumBridgeClaimAttributes creates some attributes for ethereum bridge claim
func CreateEthereumBridgeClaimAttributes(t *testing.T) []abci.EventAttribute {
	attributes := [3]abci.EventAttribute{}

	// (key, value) pairing for "cosmos_sender" key
	pairCosmosSender := abci.EventAttribute{
		Key:   []byte("cosmos_sender"),
		Value: []byte(TestCosmosValAddress),
	}

	// (key, value) pairing for "cosmos_sender_sequence" key
	pairCosmosSenderSequence := abci.EventAttribute{
		Key:   []byte("cosmos_sender_sequence"),
		Value: []byte(strconv.Itoa(TestCosmosAddressSequence)),
	}

	// (key, value) pairing for "ethereum_receiver" key
	pairEthereumReceiver := abci.EventAttribute{
		Key:   []byte("ethereum_sender"),
		Value: []byte(common.HexToAddress(TestEthereumAddress1).Hex()), // .Bytes() doesn't seem to work here
	}

	// Assign pairs to attributes array
	attributes[0] = pairCosmosSender
	attributes[1] = pairCosmosSenderSequence
	attributes[2] = pairEthereumReceiver

	return attributes[:]
}

// CreateInvalidCosmosSenderEthereumBridgeClaimAttributes creates some invalide attributes for ethereum bridge claim
func CreateInvalidCosmosSenderEthereumBridgeClaimAttributes(t *testing.T) []abci.EventAttribute {
	attributes := [3]abci.EventAttribute{}

	// (key, value) pairing for "cosmos_sender" key
	pairCosmosSender := abci.EventAttribute{
		Key:   []byte("cosmos_sender"),
		Value: []byte(TestEthereumAddress1),
	}

	// (key, value) pairing for "cosmos_sender_sequence" key
	pairCosmosSenderSequence := abci.EventAttribute{
		Key:   []byte("cosmos_sender_sequence"),
		Value: []byte(strconv.Itoa(TestCosmosAddressSequence)),
	}

	// (key, value) pairing for "ethereum_receiver" key
	pairEthereumReceiver := abci.EventAttribute{
		Key:   []byte("ethereum_sender"),
		Value: []byte(common.HexToAddress(TestEthereumAddress1).Hex()), // .Bytes() doesn't seem to work here
	}

	// Assign pairs to attributes array
	attributes[0] = pairCosmosSender
	attributes[1] = pairCosmosSenderSequence
	attributes[2] = pairEthereumReceiver

	return attributes[:]
}

// CreateInvalidEthereumSenderEthereumBridgeClaimAttributes creates some attributes for ethereum bridge claim
func CreateInvalidEthereumSenderEthereumBridgeClaimAttributes(t *testing.T) []abci.EventAttribute {
	attributes := [3]abci.EventAttribute{}

	// (key, value) pairing for "cosmos_sender" key
	pairCosmosSender := abci.EventAttribute{
		Key:   []byte("cosmos_sender"),
		Value: []byte(TestCosmosValAddress),
	}

	// (key, value) pairing for "cosmos_sender_sequence" key
	pairCosmosSenderSequence := abci.EventAttribute{
		Key:   []byte("cosmos_sender_sequence"),
		Value: []byte(strconv.Itoa(TestCosmosAddressSequence)),
	}

	// (key, value) pairing for "ethereum_receiver" key
	pairEthereumReceiver := abci.EventAttribute{
		Key:   []byte("ethereum_sender"),
		Value: []byte(TestCosmosValAddress), // .Bytes() doesn't seem to work here
	}

	// Assign pairs to attributes array
	attributes[0] = pairCosmosSender
	attributes[1] = pairCosmosSenderSequence
	attributes[2] = pairEthereumReceiver

	return attributes[:]
}

// CreateInvalidSequenceEthereumBridgeClaimAttributes creates some attributes for ethereum bridge claim
func CreateInvalidSequenceEthereumBridgeClaimAttributes(t *testing.T) []abci.EventAttribute {
	attributes := [3]abci.EventAttribute{}

	// (key, value) pairing for "cosmos_sender" key
	pairCosmosSender := abci.EventAttribute{
		Key:   []byte("cosmos_sender"),
		Value: []byte(TestCosmosValAddress),
	}

	// (key, value) pairing for "cosmos_sender_sequence" key
	pairCosmosSenderSequence := abci.EventAttribute{
		Key:   []byte("cosmos_sender_sequence"),
		Value: []byte(TestCosmosValAddress),
	}

	// (key, value) pairing for "ethereum_receiver" key
	pairEthereumReceiver := abci.EventAttribute{
		Key:   []byte("ethereum_sender"),
		Value: []byte("wrong sequence"), // .Bytes() doesn't seem to work here
	}

	// Assign pairs to attributes array
	attributes[0] = pairCosmosSender
	attributes[1] = pairCosmosSenderSequence
	attributes[2] = pairEthereumReceiver

	return attributes[:]
}
