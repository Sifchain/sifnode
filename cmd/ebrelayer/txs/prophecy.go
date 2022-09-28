package txs

import (
	"errors"
	"math/big"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// ProphecyCompletedEventToProphecyInfo parses data from a prophecy completed event witnessed on Cosmos
func ProphecyCompletedEventToProphecyInfo(attributes []abci.EventAttribute, sugaredLogger *zap.SugaredLogger) (types.ProphecyInfo, error) {
	var prophecyID []byte
	var cosmosSender []byte
	var cosmosSenderSequence *big.Int
	var ethereumReceiver common.Address
	var symbol string
	var amount big.Int
	var networkDescriptor oracletypes.NetworkDescriptor
	var doublePeg bool
	var globalSequence uint64

	var ethereumSignerAddresses []string
	var signatures []string

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
				return types.ProphecyInfo{}, errors.New("network id is invalid")
			}
			networkDescriptor = tmpNetworkDescriptor

		case types.CosmosSender.String():
			cosmosSender = []byte(val)
			attributeNumber++

		case types.CosmosSenderSequence.String():
			attributeNumber++
			tempSequence := new(big.Int)
			tempSequence, ok := tempSequence.SetString(val, 10)
			if !ok {
				sugaredLogger.Errorw("Invalid account sequence", "account sequence", val)
				return types.ProphecyInfo{}, errors.New("invalid account sequence: " + val)
			}
			cosmosSenderSequence = tempSequence

		case types.EthereumReceiver.String():
			attributeNumber++
			if !common.IsHexAddress(val) {
				sugaredLogger.Errorw("Invalid recipient address", "recipient address", val)

				return types.ProphecyInfo{}, errors.New("invalid recipient address: " + val)
			}
			ethereumReceiver = common.HexToAddress(val)

		case types.Symbol.String():
			attributeNumber++
			symbol = val

		case types.Amount.String():
			attributeNumber++
			tempAmount, ok := sdk.NewIntFromString(val)
			if !ok {
				sugaredLogger.Errorw("Invalid amount", "amount", val)

				return types.ProphecyInfo{}, errors.New("invalid amount:" + val)
			}
			amount = *big.NewInt(tempAmount.Int64())

		case types.DoublePeg.String():
			attributeNumber++
			tmpDoublePeg, err := strconv.ParseBool(val)
			if err != nil {
				sugaredLogger.Errorw("double peg can't parse", "doublePeg", val)
				return types.ProphecyInfo{}, errors.New("network id can't parse")
			}
			doublePeg = tmpDoublePeg

		case types.GlobalSequence.String():
			attributeNumber++
			tempGlobalSequence, ok := sdk.NewIntFromString(val)
			if !ok {
				sugaredLogger.Errorw("Invalid global nonce", "global nonce", val)

				return types.ProphecyInfo{}, errors.New("invalid amount:" + val)
			}
			globalSequence = tempGlobalSequence.Uint64()

		case types.EthereumAddresses.String():
			attributeNumber++
			ethereumSignerAddresses = strings.Split(val, ",")

		case types.Signatures.String():
			attributeNumber++
			signatures = strings.Split(val, ",")
		}
	}

	if attributeNumber < 11 {
		sugaredLogger.Errorw("message not complete", "attributeNumber", attributeNumber)
		return types.ProphecyInfo{}, errors.New("message not complete")
	}

	return types.ProphecyInfo{
		ProphecyID:              prophecyID,
		NetworkDescriptor:       oracletypes.NetworkDescriptor(networkDescriptor),
		CosmosSender:            string(cosmosSender),
		CosmosSenderSequence:    cosmosSenderSequence.Uint64(),
		EthereumReceiver:        ethereumReceiver.String(),
		TokenSymbol:             symbol,
		TokenAmount:             amount,
		DoublePeg:               doublePeg,
		GlobalSequence:          globalSequence,
		EthereumSignerAddresses: ethereumSignerAddresses,
		Signatures:              signatures,
	}, nil
}
