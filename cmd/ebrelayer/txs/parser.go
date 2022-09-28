package txs

import (
	"errors"
	"log"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
)

// EthereumEventToEthBridgeClaim parses and packages an Ethereum event struct with a validator address in an EthBridgeClaim msg
func EthereumEventToEthBridgeClaim(valAddr sdk.ValAddress, event types.EthereumEvent, symbolTranslator *symbol_translator.SymbolTranslator, sugaredLogger *zap.SugaredLogger) (ethbridge.EthBridgeClaim, error) {
	ethBridgeClaim := ethbridge.EthBridgeClaim{}

	networkDescriptor := oracletypes.NetworkDescriptor(event.NetworkDescriptor)

	if err := ethbridge.ValidateNetworkDescriptor(networkDescriptor); err != nil {
		return ethBridgeClaim, err
	}

	bridgeContractAddress := ethbridge.NewEthereumAddress(event.BridgeContractAddress.Hex())

	// Sender type casting (address.common -> string)
	sender := ethbridge.NewEthereumAddress(event.From.Hex())

	// Recipient type casting ([]bytes -> sdk.AccAddress)
	recipient, err := sdk.AccAddressFromBech32(string(event.To))
	if err != nil {
		return ethBridgeClaim, err
	}
	if recipient.Empty() {
		return ethBridgeClaim, errors.New("empty recipient address")
	}

	// Sender type casting (address.common -> string)
	tokenContractAddress := ethbridge.NewEthereumAddress(event.Token.Hex())

	// Symbol formatted to lowercase
	symbol := strings.ToLower(event.Symbol)
	if event.ClaimType == ethbridge.ClaimType_CLAIM_TYPE_BURN {
		symbol = symbolTranslator.EthereumToSifchain(symbol)
	}

	amount := sdk.NewIntFromBigInt(event.Value)

	// Package the information in a unique EthBridgeClaim
	ethBridgeClaim.NetworkDescriptor = networkDescriptor
	ethBridgeClaim.BridgeContractAddress = bridgeContractAddress.String()
	ethBridgeClaim.EthereumLockBurnSequence = event.Nonce.Uint64()
	ethBridgeClaim.TokenContractAddress = tokenContractAddress.String()
	ethBridgeClaim.Symbol = symbol
	ethBridgeClaim.EthereumSender = sender.String()
	ethBridgeClaim.ValidatorAddress = valAddr.String()
	ethBridgeClaim.CosmosReceiver = recipient.String()
	ethBridgeClaim.Amount = amount
	ethBridgeClaim.ClaimType = event.ClaimType
	ethBridgeClaim.Decimals = int64(event.Decimals)
	ethBridgeClaim.TokenName = event.Name
	// the nonce from ethereum event is lock burn nonce, not transaction nonce
	ethBridgeClaim.Denom = ethbridge.GetDenom(networkDescriptor, tokenContractAddress)
	ethBridgeClaim.CosmosDenom = event.CosmosDenom
	return ethBridgeClaim, nil
}

// BurnLockEventToCosmosMsg parses data from a Burn/Lock event witnessed on Cosmos into a CosmosMsg struct
func BurnLockEventToCosmosMsg(attributes []abci.EventAttribute, sugaredLogger *zap.SugaredLogger) (types.CosmosMsg, error) {
	var prophecyID []byte
	var networkDescriptor oracletypes.NetworkDescriptor
	var globalSequence uint64

	attributeNumber := 0

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())

		// Set variable based on the attribute's key
		switch key {
		case types.ProphecyID.String():
			prophecyID = []byte(val)
			attributeNumber++

		case types.NetworkDescriptor.String():
			attributeNumber++
			tmpNetworkDescriptor, err := oracletypes.ParseNetworkDescriptor(val)

			if err != nil {
				sugaredLogger.Errorw("network id can't parse", "networkDescriptor", val)
				return types.CosmosMsg{}, errors.New("network id is invalid")
			}
			networkDescriptor = tmpNetworkDescriptor

		case types.GlobalSequence.String():
			attributeNumber++
			tempGlobalSequence, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				sugaredLogger.Errorw("globalSequence can't parse", "globalSequence", val)
				return types.CosmosMsg{}, errors.New("globalSequence can't parse")
			}
			globalSequence = tempGlobalSequence
		}
	}

	if attributeNumber < 3 {
		sugaredLogger.Errorw("message not complete", "attributeNumber", attributeNumber)
		return types.CosmosMsg{}, errors.New("message not complete")
	}
	return types.NewCosmosMsg(networkDescriptor, prophecyID, globalSequence), nil
}

// AttributesToEthereumBridgeClaim parses data from event to EthereumBridgeClaim
func AttributesToEthereumBridgeClaim(attributes []abci.EventAttribute) (types.EthereumBridgeClaim, error) {
	var cosmosSender sdk.ValAddress
	var ethereumSenderSequence sdk.Int
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

		case types.EthereumSenderSequence.String():
			tempSequence, ok := sdk.NewIntFromString(val)
			if !ok {
				log.Println("Invalid nonce:", val)
				return types.EthereumBridgeClaim{}, errors.New("invalid nonce:" + val)
			}
			ethereumSenderSequence = tempSequence
		}
	}

	return types.EthereumBridgeClaim{
		EthereumSender: ethereumSender,
		CosmosSender:   cosmosSender,
		Sequence:       ethereumSenderSequence,
	}, nil
}

// AttributesToCosmosSignProphecyClaim parses data from event to EthereumBridgeClaim
func AttributesToCosmosSignProphecyClaim(attributes []abci.EventAttribute) (types.CosmosSignProphecyClaim, error) {
	var cosmosSender sdk.ValAddress
	var networkDescriptor oracletypes.NetworkDescriptor
	var prophecyID []byte
	var err error
	attributeNumber := 0

	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		val := string(attribute.GetValue())

		// Set variable based on the attribute's key
		switch key {
		case types.CosmosSender.String():
			cosmosSender, err = sdk.ValAddressFromBech32(val)
			if err != nil {
				return types.CosmosSignProphecyClaim{}, err
			}

		case types.NetworkDescriptor.String():
			attributeNumber++
			tempNetworkDescriptor, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				log.Printf("network id can't parse networkDescriptor is %s\n", val)
				return types.CosmosSignProphecyClaim{}, errors.New("network id can't parse")
			}
			networkDescriptor = oracletypes.NetworkDescriptor(uint32(tempNetworkDescriptor))

			// check if the networkDescriptor is valid
			if !networkDescriptor.IsValid() {
				return types.CosmosSignProphecyClaim{}, errors.New("network id is invalid")
			}

		case types.ProphecyID.String():
			prophecyID = []byte(val)
		}
	}

	return types.CosmosSignProphecyClaim{
		CosmosSender:      cosmosSender,
		NetworkDescriptor: networkDescriptor,
		ProphecyID:        prophecyID,
	}, nil
}

// isZeroAddress checks an Ethereum address and returns a bool which indicates if it is the null address
func isZeroAddress(address common.Address) bool {
	return address == common.HexToAddress(nullAddress)
}
