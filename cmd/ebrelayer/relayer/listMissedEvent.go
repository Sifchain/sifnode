package relayer

// DONTCOVER

import (
	"context"
	"log"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"go.uber.org/zap"
)

// ListMissedCosmosEvent defines a Cosmos listener that relays events to Ethereum and Cosmos
type ListMissedCosmosEvent struct {
	NetworkID               uint32
	TmProvider              string
	EthProvider             string
	RegistryContractAddress common.Address
	EthereumAddress         common.Address
	Days                    int64
	SugaredLogger           *zap.SugaredLogger
}

// NewListMissedCosmosEvent initializes a new CosmosSub
func NewListMissedCosmosEvent(networkID uint32, tmProvider, ethProvider string, registryContractAddress common.Address,
	ethereumAddress common.Address, days int64, sugaredLogger *zap.SugaredLogger) ListMissedCosmosEvent {
	return ListMissedCosmosEvent{
		NetworkID:               networkID,
		TmProvider:              tmProvider,
		EthProvider:             ethProvider,
		RegistryContractAddress: registryContractAddress,
		EthereumAddress:         ethereumAddress,
		Days:                    days,
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

	header, err := ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Printf("%s \n", err.Error())
		return
	}

	currentEthHeight := header.Number.Int64()
	// estimate blocks by one block every 15 seconds
	blocks := 4 * 60 * 24 * list.Days
	ethFromHeight := currentEthHeight - blocks
	if ethFromHeight < 0 {
		ethFromHeight = 0
	}

	ProphecyClaims := GetAllProphecyClaim(ethClient, list.EthereumAddress, ethFromHeight, currentEthHeight)

	log.Printf("found out %d prophecy claims I sent from %d to %d block\n", len(ProphecyClaims), ethFromHeight, currentEthHeight)

	client, err := tmClient.New(list.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client %s\n", err.Error())
		return
	}

	ctx := context.Background()
	block, err := client.Block(ctx, nil)
	if err != nil {
		log.Printf("%s \n", err.Error())
		return
	}

	currentCosmosHeight := block.Block.Header.Height
	// estimate blocks by one block every 6 seconds
	blocks = 10 * 60 * 24 * list.Days
	cosmosFromHeight := currentCosmosHeight - blocks
	if cosmosFromHeight < 0 {
		cosmosFromHeight = 0
	}

	if err := client.Start(); err != nil {
		log.Printf("failed to start a client %s\n", err.Error())
		return
	}

	defer client.Stop() //nolint:errcheck

	for blockNumber := cosmosFromHeight; blockNumber < currentCosmosHeight; {
		tmpBlockNumber := blockNumber

		block, err := client.BlockResults(ctx, &tmpBlockNumber)
		blockNumber++

		if err != nil {
			continue
		}

		for _, result := range block.TxsResults {
			for _, event := range result.Events {

				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case types.MsgBurn, types.MsgLock:

					cosmosMsg, networkID, err := txs.BurnLockEventToCosmosMsg(claimType, event.GetAttributes(), list.SugaredLogger)
					if err != nil {
						log.Println(err.Error())
						continue
					}

					if networkID == list.NetworkID && !MessageProcessed(cosmosMsg, ProphecyClaims) {
						log.Printf("missed cosmos event: %s\n", cosmosMsg.String())
					}
				}
			}
		}
	}
}
