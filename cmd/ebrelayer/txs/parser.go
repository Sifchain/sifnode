package txs

import (
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

const (
	nullAddress           = "0x0000000000000000000000000000000000000000"
	defaultSifchainPrefix = "c"
	defaultEthereumPrefix = "e"
)

// EthereumEventToEthBridgeClaim parses and packages an Ethereum event struct with a validator address in an EthBridgeClaim msg
func EthereumEventToEthBridgeClaim(valAddr sdk.AccAddress, event types.EthereumEvent) (ethbridge.EthBridgeClaim, error) {
	witnessClaim := ethbridge.EthBridgeClaim{}

	// chainID type casting (*big.Int -> int)
	chainID := int(event.EthereumChainID.Int64())

	bridgeContractAddress := ethbridge.NewEthereumAddress(event.BridgeContractAddress.Hex())

	// Sender type casting (address.common -> string)
	sender := ethbridge.NewEthereumAddress(event.From.Hex())

	// Recipient type casting ([]bytes -> sdk.AccAddress)
	recipient, err := sdk.AccAddressFromBech32(string(event.To))
	if err != nil {
		return witnessClaim, err
	}
	if recipient.Empty() {
		return witnessClaim, errors.New("empty recipient address")
	}

	// Sender type casting (address.common -> string)
	tokenContractAddress := ethbridge.NewEthereumAddress(event.Token.Hex())

	// Symbol formatted to lowercase
	symbol := strings.ToLower(event.Symbol)
	switch event.ClaimType {
	case ethbridge.ClaimType_CLAIM_TYPE_LOCK:
		if symbol == "eth" && !isZeroAddress(event.Token) {
			return witnessClaim, errors.New("symbol \"eth\" must have null address set as token address")
		}
	case ethbridge.ClaimType_CLAIM_TYPE_BURN:
		if !strings.Contains(symbol, defaultEthereumPrefix) {
			log.Printf("Can only relay burns of '%v' prefixed tokens", defaultEthereumPrefix)
			return witnessClaim, errors.New("symbol of burn token must start with prefix")
		}
		res := strings.SplitAfter(symbol, defaultEthereumPrefix)
		symbol = strings.Join(res[1:], "")
	}

	amount := sdk.NewIntFromBigInt(event.Value)

	// Nonce type casting (*big.Int -> int)
	nonce := int(event.Nonce.Int64())

	// Package the information in a unique EthBridgeClaim
	witnessClaim.EthereumChainId = int64(chainID)
	witnessClaim.BridgeContractAddress = bridgeContractAddress.String()
	witnessClaim.Nonce = int64(nonce)
	witnessClaim.TokenContractAddress = tokenContractAddress.String()
	witnessClaim.Symbol = symbol
	witnessClaim.EthereumSender = sender.String()
	witnessClaim.ValidatorAddress = valAddr.String()
	witnessClaim.CosmosReceiver = recipient.String()
	witnessClaim.Amount = amount
	witnessClaim.ClaimType = event.ClaimType

	return witnessClaim, nil
}

// ProphecyClaimToSignedOracleClaim packages and signs a prophecy claim's data, returning a new oracle claim
func ProphecyClaimToSignedOracleClaim(event types.ProphecyClaimEvent, key *ecdsa.PrivateKey) (OracleClaim, error) {
	oracleClaim := OracleClaim{}

	// Generate a hashed claim message which contains ProphecyClaim's data
	message := GenerateClaimMessage(event)

	// Sign the message using the validator's private key
	signature, err := SignClaim(PrefixMsg(message), key)
	if err != nil {
		return oracleClaim, err
	}

	oracleClaim.ProphecyID = event.ProphecyID
	var message32 [32]byte
	copy(message32[:], message)
	oracleClaim.Message = message32
	oracleClaim.Signature = signature
	return oracleClaim, nil
}

// CosmosMsgToProphecyClaim parses event data from a CosmosMsg, packaging it as a ProphecyClaim
func CosmosMsgToProphecyClaim(event types.CosmosMsg) ProphecyClaim {
	claimType := event.ClaimType
	cosmosSender := event.CosmosSender
	cosmosSenderSequence := event.CosmosSenderSequence
	ethereumReceiver := event.EthereumReceiver
	symbol := event.Symbol
	amount := event.Amount

	prophecyClaim := ProphecyClaim{
		ClaimType:            claimType,
		CosmosSender:         cosmosSender,
		CosmosSenderSequence: cosmosSenderSequence,
		EthereumReceiver:     ethereumReceiver,
		Symbol:               symbol,
		Amount:               amount,
	}
	return prophecyClaim
}

// BurnLockEventToCosmosMsg parses data from a Burn/Lock event witnessed on Cosmos into a CosmosMsg struct
func BurnLockEventToCosmosMsg(claimType types.Event, attributes []abci.EventAttribute, sugaredLogger *zap.SugaredLogger) (types.CosmosMsg, error) {
	var cosmosSender []byte
	var cosmosSenderSequence *big.Int
	var ethereumReceiver common.Address
	var symbol string
	var amount sdk.Int

	attributeNumber := 0

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())

		// Set variable based on the attribute's key
		switch key {
		case types.CosmosSender.String():
			cosmosSender = []byte(val)
			attributeNumber++
		case types.CosmosSenderSequence.String():
			attributeNumber++
			tempSequence := new(big.Int)
			tempSequence, ok := tempSequence.SetString(val, 10)
			if !ok {
				// log.Println("Invalid account sequence:", val)
				sugaredLogger.Errorw("Invalid account sequence", "account sequence", val)
				return types.CosmosMsg{}, errors.New("invalid account sequence: " + val)
			}
			cosmosSenderSequence = tempSequence
		case types.EthereumReceiver.String():
			attributeNumber++
			if !common.IsHexAddress(val) {
				// log.Printf("Invalid recipient address: %v", val)
				sugaredLogger.Errorw("Invalid recipient address", "recipient address", val)

				return types.CosmosMsg{}, errors.New("invalid recipient address: " + val)
			}
			ethereumReceiver = common.HexToAddress(val)
		case types.Symbol.String():
			attributeNumber++
			if claimType == types.MsgBurn {
				if !strings.Contains(val, defaultSifchainPrefix) {
					// log.Printf("Can only relay burns of '%v' prefixed coins", defaultSifchainPrefix)
					sugaredLogger.Errorw("only relay burns prefixed coins", "coin symbol", val)
					return types.CosmosMsg{}, errors.New("can only relay burns of '%v' prefixed coins" + defaultSifchainPrefix)
				}
				res := strings.SplitAfter(val, defaultSifchainPrefix)
				symbol = strings.Join(res[1:], "")
			} else {
				symbol = val
			}
		case types.Amount.String():
			attributeNumber++
			tempAmount, ok := sdk.NewIntFromString(val)
			if !ok {
				// log.Println("Invalid amount:", val)
				sugaredLogger.Errorw("Invalid amount", "amount", val)

				return types.CosmosMsg{}, errors.New("invalid amount:" + val)
			}
			amount = tempAmount
		}
	}

	if attributeNumber < 5 {
		sugaredLogger.Errorw("message not complete", "attributeNumber", attributeNumber)
		return types.CosmosMsg{}, errors.New("message not complete")
	}
	return types.NewCosmosMsg(claimType, cosmosSender, cosmosSenderSequence, ethereumReceiver, symbol, amount), nil
}

// AttributesToEthereumBridgeClaim parses data from event to EthereumBridgeClaim
func AttributesToEthereumBridgeClaim(attributes []abci.EventAttribute) (types.EthereumBridgeClaim, error) {
	var cosmosSender sdk.ValAddress
	var ethereumSenderNonce sdk.Int
	var ethereumSender common.Address
	var err error

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())

		// Set variable based on the attribute's key
		switch key {
		case types.CosmosSender.String():
			cosmosSender, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				return types.EthereumBridgeClaim{}, err
			}

		case types.EthereumSender.String():
			if !common.IsHexAddress(val) {
				log.Printf("Invalid recipient address: %v", val)
				return types.EthereumBridgeClaim{}, errors.New("invalid recipient address: " + val)
			}
			ethereumSender = common.HexToAddress(val)

		case types.EthereumSenderNonce.String():
			tempNonce, ok := sdk.NewIntFromString(val)
			if !ok {
				log.Println("Invalid nonce:", val)
				return types.EthereumBridgeClaim{}, errors.New("invalid nonce:" + val)
			}
			ethereumSenderNonce = tempNonce
		}
	}

	return types.EthereumBridgeClaim{
		EthereumSender: ethereumSender,
		CosmosSender:   cosmosSender,
		Nonce:          ethereumSenderNonce,
	}, nil
}

// isZeroAddress checks an Ethereum address and returns a bool which indicates if it is the null address
func isZeroAddress(address common.Address) bool {
	return address == common.HexToAddress(nullAddress)
}
