package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
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
)

const (
	FlagSymbolTranslatorFile = "symbol-translator-file"
)

// String returns the event type as a string
func (d Event) String() string {
	return [...]string{"unsupported", "burn", "lock", "LogLock", "LogBurn", "LogNewProphecyClaim", "newProphecyClaim", "create_claim"}[d]
}

// EthereumEvent struct is used by LogLock and LogBurn
type EthereumEvent struct {
	To                    []byte
	Symbol                string
	EthereumChainID       *big.Int
	Value                 *big.Int
	Nonce                 *big.Int
	ClaimType             ethbridge.ClaimType
	ID                    [32]byte
	BridgeContractAddress common.Address
	From                  common.Address
	Token                 common.Address
}

// Equal two events
func (e EthereumEvent) Equal(other EthereumEvent) bool {
	return e.EthereumChainID == other.EthereumChainID &&
		e.BridgeContractAddress == other.BridgeContractAddress &&
		bytes.Equal(e.ID[:], other.ID[:]) &&
		e.From == other.From &&
		bytes.Equal(e.To, other.To) &&
		e.Symbol == other.Symbol &&
		e.Value.Cmp(other.Value) == 0 &&
		e.Nonce.Cmp(other.Nonce) == 0 &&
		e.ClaimType == other.ClaimType
}

// String implements fmt.Stringer
func (e EthereumEvent) String() string {
	return fmt.Sprintf("\nChain ID: %v\nBridge contract address: %v\nToken symbol: %v\nToken "+
		"contract address: %v\nSender: %v\nRecipient: %v\nValue: %v\nNonce: %v\nClaim type: %v",
		e.EthereumChainID, e.BridgeContractAddress.Hex(), e.Symbol, e.Token.Hex(), e.From.Hex(),
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
	CosmosSender         []byte
	CosmosSenderSequence *big.Int
	Symbol               string
	Amount               sdk.Int
	EthereumReceiver     common.Address
	ClaimType            Event
}

// NewCosmosMsg creates a new CosmosMsg
func NewCosmosMsg(claimType Event, cosmosSender []byte, cosmosSenderSequence *big.Int, ethereumReceiver common.Address, symbol string,
	amount sdk.Int) CosmosMsg {
	return CosmosMsg{
		ClaimType:            claimType,
		CosmosSender:         cosmosSender,
		CosmosSenderSequence: cosmosSenderSequence,
		EthereumReceiver:     ethereumReceiver,
		Symbol:               symbol,
		Amount:               amount,
	}
}

// String implements fmt.Stringer
func (c CosmosMsg) String() string {
	if c.ClaimType == MsgLock {
		return fmt.Sprintf("\nClaim Type: %v\nCosmos Sender: %v\nCosmos Sender Sequence: %v\nEthereum Recipient: %v"+
			"\nSymbol: %v\nAmount: %v\n",
			c.ClaimType.String(), string(c.CosmosSender), c.CosmosSenderSequence, c.EthereumReceiver.Hex(), c.Symbol, c.Amount)
	}
	return fmt.Sprintf("\nClaim Type: %v\nCosmos Sender: %v\nCosmos Sender Sequence: %v\nEthereum Recipient: %v"+
		"\nSymbol: %v\nAmount: %v\n",
		c.ClaimType.String(), string(c.CosmosSender), c.CosmosSenderSequence, c.EthereumReceiver.Hex(), c.Symbol, c.Amount)
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
	// EthereumSenderNonce is ethereum sender nonce
	EthereumSenderNonce
)

// String returns the event type as a string
func (d CosmosMsgAttributeKey) String() string {
	return [...]string{"unsupported", "cosmos_sender", "cosmos_sender_sequence", "ethereum_receiver", "amount", "symbol", "ethereum_sender", "ethereum_sender_nonce"}[d]
}

// EthereumBridgeClaim for store the EventTypeCreateClaim from cosmos
type EthereumBridgeClaim struct {
	EthereumSender common.Address
	CosmosSender   sdk.ValAddress
	Nonce          sdk.Int
}

// ProphecyClaimUnique for data part of ProphecyClaim transaction in Ethereum
type ProphecyClaimUnique struct {
	CosmosSenderSequence *big.Int
	CosmosSender         []byte
}
