package relayer

// DONTCOVER

import (
	"context"
	"log"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"go.uber.org/zap"
)

// ListMissedCosmosEvent defines structre that get all not processed signature aggregation completed messages
type ListMissedCosmosEvent struct {
	NetworkDescriptor       oracletypes.NetworkDescriptor
	TmProvider              string
	EthProvider             string
	RegistryContractAddress common.Address
	EthereumAddress         common.Address
	SugaredLogger           *zap.SugaredLogger
}

// NewListMissedCosmosEvent initializes a new CosmosSub
func NewListMissedCosmosEvent(networkDescriptor oracletypes.NetworkDescriptor, tmProvider, ethProvider string, registryContractAddress common.Address,
	ethereumAddress common.Address, sugaredLogger *zap.SugaredLogger) ListMissedCosmosEvent {
	return ListMissedCosmosEvent{
		NetworkDescriptor:       networkDescriptor,
		TmProvider:              tmProvider,
		EthProvider:             ethProvider,
		RegistryContractAddress: registryContractAddress,
		EthereumAddress:         ethereumAddress,
		SugaredLogger:           sugaredLogger,
	}
}

// ListMissedCosmosEvent print all missed cosmos events by this ebrelayer in days
func (list ListMissedCosmosEvent) ListMissedCosmosEvent() {
	log.Println("ListMissedCosmosEvent started")
	// Start Ethereum client
	ethClient, err := ethclient.Dial(list.EthProvider)
	if err != nil {
		log.Printf("%s \n", err.Error())
		return
	}

	cosmosBridgeAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, list.RegistryContractAddress, txs.CosmosBridge, list.SugaredLogger)
	if err != nil {
		log.Printf("failed to get the cosmos bridge address, error as %s\n", err)
		return
	}

	lastSubmittedNonce, err := GetLastNonceSubmitted(ethClient, cosmosBridgeAddress, list.SugaredLogger)
	if err != nil {
		log.Printf("failed to get the last submitted nonce, error as %s\n", err)
		return
	}

	tmClient, err := tmClient.New(list.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client %s\n", err.Error())
		return
	}

	ctx := context.Background()
	block, err := tmClient.Block(ctx, nil)
	if err != nil {
		log.Printf("%s \n", err.Error())
		return
	}

	currentCosmosHeight := block.Block.Header.Height
	var toBlock int64
	if currentCosmosHeight > ProphecyLiftTime {
		toBlock = currentCosmosHeight - ProphecyLiftTime
	} else {
		toBlock = 0
	}

	if err := tmClient.Start(); err != nil {
		log.Printf("failed to start a client %s\n", err.Error())
		return
	}

	defer tmClient.Stop() //nolint:errcheck

	for blockNumber := currentCosmosHeight; blockNumber > toBlock; {
		endLoop := false
		tmpBlockNumber := blockNumber

		block, err := tmClient.BlockResults(ctx, &tmpBlockNumber)
		blockNumber--

		if err != nil {
			continue
		}

		for _, result := range block.TxsResults {
			for _, event := range result.Events {

				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case types.ProphecyCompleted:

					cosmosMsg, err := txs.ProphecyCompletedEventToProphecyInfo(event.GetAttributes(), list.SugaredLogger)
					if err != nil {
						log.Println(err.Error())
						continue
					}

					if cosmosMsg.GlobalNonce <= lastSubmittedNonce.Uint64() {
						endLoop = true
					} else if cosmosMsg.NetworkDescriptor == list.NetworkDescriptor {
						log.Printf("missed cosmos event: %s\n", cosmosMsg.String())
					}
				}
			}
		}

		// exit from loop, not check previous block anymore
		if endLoop {
			break
		}

	}
}
