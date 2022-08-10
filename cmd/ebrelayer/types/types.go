package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracle "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event enum containing supported chain events
type Event byte

const (
	// Unsupported is an invalid Cosmos or Ethereum event
	Unsupported Event = iota
	// MsgBurn is a Cosmos msg of type MsgBurn
	MsgBurn
	// MsgLock is a Cosmos msg of type MsgLock
	MsgLock
	// LogLock is for Ethereum event LogLock
	LogLock
	// LogBurn is for Ethereum event LogBurn
	LogBurn
	// LogNewProphecyClaim is an Ethereum event named 'LogNewProphecyClaim'
	LogNewProphecyClaim
	// NewProphecyClaim for newProphecyClaim method in smart contract
	NewProphecyClaim
	// CreateEthBridgeClaim is a Cosmos msg of type MsgCreateEthBridgeClaim
	CreateEthBridgeClaim
	// ProphecyCompleted is Cosmos event EventTypeProphecyCompleted
	ProphecyCompleted
	// SubmitProphecyClaimAggregatedSigs is Ethereum method name
	SubmitProphecyClaimAggregatedSigs
	FlagSymbolTranslatorFile = "symbol-translator-file"
)

// String returns the event type as a string
func (d Event) String() string {
	return [...]string{
		"unsupported",
		ethbridge.EventTypeBurn,
		ethbridge.EventTypeLock,
		ethbridge.EventTypeProphecyStatus,
		ethbridge.EventTypeCreateClaim,
		ethbridge.EventTypeUpdateWhiteListValidator,
		ethbridge.EventTypeSetCrossChainFee,
		ethbridge.EventTypeUpdateCrossChainFeeReceiverAccount,
		ethbridge.EventTypeSignProphecy,

		// TODO: Suspect these can be removed
		LogLock.String(),                    // TODO: Where is this used?
		LogBurn.String(),                    // TODO: Where is this used?
		LogNewProphecyClaim.String(),        // TODO: Where is this used?
		NewProphecyClaim.String(),           // TODO: Where is this used?
		"newProphecyClaim",                  // TODO: Where is this used?
		"prophecy_completed",                // TODO: Where is this used?
		"submitProphecyClaimAggregatedSigs", // TODO: Where is this used?
	}[d]
}

// EthereumEvent struct is used by LogLock and LogBurn
type EthereumEvent struct {
	To                    []byte
	Symbol                string
	Name                  string
	Decimals              uint8
	NetworkDescriptor     int32
	Value                 *big.Int
	Nonce                 *big.Int
	ClaimType             ethbridge.ClaimType
	BridgeContractAddress common.Address
	From                  common.Address
	Token                 common.Address
	CosmosDenom           string
}

// String implements fmt.Stringer
func (e EthereumEvent) String() string {
	return fmt.Sprintf("\nNetwork ID: %v\nBridge contract address: %v\nToken symbol: %v\nToken "+
		"contract address: %v\nSender: %v\nRecipient: %v\nValue: %v\nNonce: %v\nClaim type: %v",
		e.NetworkDescriptor, e.BridgeContractAddress.Hex(), e.Symbol, e.Token.Hex(), e.From.Hex(),
		string(e.To), e.Value, e.Nonce, e.ClaimType.String())
}

// ProphecyClaimEvent struct which represents a LogNewProphecyClaim event
type ProphecyClaimEvent struct {
	CosmosSender     []byte
	Symbol           string
	ProphecyID       *big.Int
	Amount           sdk.Int
	EthereumReceiver common.Address
	ValidatorAddress common.Address
	TokenAddress     common.Address
	ClaimType        uint8
}

// NewProphecyClaimEvent creates a new ProphecyClaimEvent
func NewProphecyClaimEvent(cosmosSender []byte, symbol string, prophecyID *big.Int, amount sdk.Int, ethereumReceiver,
	validatorAddress, tokenAddress common.Address, claimType uint8) ProphecyClaimEvent {
	return ProphecyClaimEvent{
		CosmosSender:     cosmosSender,
		Symbol:           symbol,
		ProphecyID:       prophecyID,
		Amount:           amount,
		EthereumReceiver: ethereumReceiver,
		ValidatorAddress: validatorAddress,
		TokenAddress:     tokenAddress,
		ClaimType:        claimType,
	}
}

// String implements fmt.Stringer
func (p ProphecyClaimEvent) String() string {
	return fmt.Sprintf("\nProphecy ID: %v\nClaim Type: %v\nSender: %v\n"+
		"Recipient: %v\nSymbol: %v\nToken: %v\nAmount: %v\nValidator: %v\n\n",
		p.ProphecyID, p.ClaimType, string(p.CosmosSender), p.EthereumReceiver.Hex(),
		p.Symbol, p.TokenAddress.Hex(), p.Amount, p.ValidatorAddress.Hex())
}

// CosmosMsg contains data from MsgBurn and MsgLock events
type CosmosMsg struct {
	NetworkDescriptor oracle.NetworkDescriptor
	ProphecyID        []byte
	GlobalSequence    uint64
}

// NewCosmosMsg creates a new CosmosMsg
func NewCosmosMsg(networkDescriptor oracle.NetworkDescriptor, prophecyID []byte, globalSequence uint64) CosmosMsg {
	return CosmosMsg{
		NetworkDescriptor: networkDescriptor,
		ProphecyID:        prophecyID,
		GlobalSequence:    globalSequence,
	}
}

// String implements fmt.Stringer
func (c CosmosMsg) String() string {
	return fmt.Sprintf("\nNetwork id: %v\nProphecy ID: %v\nGlobal sequence: %v\n", c.NetworkDescriptor.String(), c.ProphecyID, c.GlobalSequence)
}

// CosmosMsgAttributeKey enum containing supported attribute keys
type CosmosMsgAttributeKey int

const (
	// UnsupportedAttributeKey unsupported attribute key
	UnsupportedAttributeKey CosmosMsgAttributeKey = iota
	// CosmosSender sender's address on Cosmos network
	CosmosSender
	// CosmosSenderSequence sender's sequence on Cosmos network
	CosmosSenderSequence
	// EthereumReceiver receiver's address on Ethereum network
	EthereumReceiver
	// Amount is coin's value
	Amount
	// Symbol is the coin type
	Symbol
	// EthereumSender is ethereum sender address
	EthereumSender
	// EthereumSenderSequence is ethereum lock burn sequence
	EthereumSenderSequence
	// NetworkDescriptor is different blockchain identity
	NetworkDescriptor
	// EthereumAddresses all relayer's ethereum addresses
	EthereumAddresses
	// Signatures all relayer's signature against one prophecy
	Signatures
	// ProphecyID id of prophecy
	ProphecyID
	// DoublePeg indicates if the token is double pegged
	DoublePeg
	// GlobalSequence an increment value for lock/burn
	GlobalSequence
)

// String returns the event type as a string
func (d CosmosMsgAttributeKey) String() string {
	return [...]string{"unsupported",
		"cosmos_sender",
		"cosmos_sender_sequence",
		"ethereum_receiver", "amount",
		"symbol", "ethereum_sender",
		"ethereum_sender_sequence",
		"network_id",
		"ethereum_addresses",
		"signatures",
		"prophecy_id",
		"double_peg",
		"global_sequence"}[d]
}

// EthereumBridgeClaim for store the EventTypeCreateClaim from cosmos
type EthereumBridgeClaim struct {
	EthereumSender common.Address
	CosmosSender   sdk.ValAddress
	Sequence       sdk.Int
}

// ProphecyClaimUnique for data part of ProphecyClaim transaction in Ethereum
type ProphecyClaimUnique struct {
	ProphecyID []byte
}

type CosmosSignProphecyClaim struct {
	CosmosSender      sdk.ValAddress
	NetworkDescriptor oracle.NetworkDescriptor
	ProphecyID        []byte
}

// ProphecyInfo store all data needed for smart contract call
type ProphecyInfo struct {
	TokenAmount             big.Int
	ProphecyID              []byte
	EthereumSignerAddresses []string
	Signatures              []string
	CosmosSender            string
	EthereumReceiver        string
	TokenSymbol             string
	CosmosSenderSequence    uint64
	GlobalSequence          uint64
	NetworkDescriptor       oracle.NetworkDescriptor
	DoublePeg               bool
}

func (info *ProphecyInfo) String() string {
	return fmt.Sprintf("%#v\n", info)
}
