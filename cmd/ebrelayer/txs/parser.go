package txs

import (
	"errors"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"go.uber.org/zap"
	"log"
	"math/big"
	"strings"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
)

const (
	nullAddress           = "0x0000000000000000000000000000000000000000"
	defaultSifchainPrefix = "c"
)

// EthereumEventToEthBridgeClaim parses and packages an Ethereum event struct with a validator address in an EthBridgeClaim msg
func EthereumEventToEthBridgeClaim(valAddr sdk.ValAddress, event types.EthereumEvent, symbolTranslator *symbol_translator.SymbolTranslator, sugaredLogger *zap.SugaredLogger) (ethbridge.EthBridgeClaim, error) {
	witnessClaim := ethbridge.EthBridgeClaim{}

	sugaredLogger.Debug("event", event)

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

	symbol := event.Symbol

	switch event.ClaimType {
	case ethbridge.ClaimType_CLAIM_TYPE_LOCK:
		symbol = strings.ToLower(event.Symbol)
		if symbol == "eth" && !isZeroAddress(event.Token) {
			return witnessClaim, errors.New("symbol \"eth\" must have null address set as token address")
		}
	case ethbridge.ClaimType_CLAIM_TYPE_BURN:
		symbol = symbolTranslator.EthereumToSifchain(symbol)
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

	sugaredLogger.Debug("witnessClaim", witnessClaim)

	return witnessClaim, nil
}

// BurnLockEventToCosmosMsg parses data from a Burn/Lock event witnessed on Cosmos into a CosmosMsg struct
func BurnLockEventToCosmosMsg(claimType types.Event, attributes []abci.EventAttribute, symbolTranslator *symbol_translator.SymbolTranslator, sugaredLogger *zap.SugaredLogger) (types.CosmosMsg, error) {
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
			switch claimType {
			case types.MsgLock:
				// Rowan is special and shouldn't run through the symbol translation system.
				// We normally have an erowan <=> rowan mapping, but we want to transition
				// away from erowan and use the new Rowan token.
				if val == "rowan" {
					symbol = val
				} else {
					symbol = symbolTranslator.SifchainToEthereum(val)
				}
			case types.MsgBurn:
				if !strings.Contains(val, defaultSifchainPrefix) {
					// log.Printf("Can only relay burns of '%v' prefixed coins", defaultSifchainPrefix)
					sugaredLogger.Errorw("only relay burns prefixed coins", "coin symbol", val)
					return types.CosmosMsg{}, errors.New("can only relay burns of '%v' prefixed coins" + defaultSifchainPrefix)
				}
				res := strings.SplitAfter(val, defaultSifchainPrefix)
				symbol = strings.Join(res[1:], "")
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
		sugaredLogger.Infow(
			"current variable values",
			"cosmosSender", cosmosSender,
			"cosmosSenderSequence", cosmosSenderSequence,
			"ethereumReceiver", ethereumReceiver,
			"symbol", symbol,
			"amount", amount,
		)

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
