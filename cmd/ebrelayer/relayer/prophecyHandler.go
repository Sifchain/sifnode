package relayer

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ethereum/go-ethereum/ethclient"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	tmTypes "github.com/tendermint/tendermint/types"
)

// StartProphecyHandler start Cosmos chain subscription and process prophecy completed message
func (sub CosmosSub) StartProphecyHandler(txFactory tx.Factory, completionEvent *sync.WaitGroup) {
	defer completionEvent.Done()
	time.Sleep(time.Second)
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent)
		return
	}

	if err := client.Start(); err != nil {
		sub.SugaredLogger.Errorw("failed to start a sifchain client.",
			errorMessageKey, err.Error())
		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent)
		return
	}

	defer client.Stop() //nolint:errcheck

	// Subscribe to all new blocks
	query := "tm.event = 'NewBlock'"
	results, err := client.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		sub.SugaredLogger.Errorw("sifchain client failed to subscribe to query.",
			errorMessageKey, err.Error(),
			"query", query)
		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent)
		return
	}

	defer func() {
		if err := client.Unsubscribe(context.Background(), "test", query); err != nil {
			sub.SugaredLogger.Errorw("sifchain client failed to unsubscribe query.",
				errorMessageKey, err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	for {
		select {
		case <-quit:
			log.Println("we receive the quit signal and exit")
			return

		case e := <-results:
			data, ok := e.Data.(tmTypes.EventDataNewBlock)
			if !ok {
				sub.SugaredLogger.Errorw("sifchain client failed to extract event data from new block.",
					"EventDataNewBlock", fmt.Sprintf("%v", e.Data))
			}
			blockHeight := data.Block.Height

			ctx := context.Background()
			block, err := client.BlockResults(ctx, &blockHeight)

			if err != nil {
				sub.SugaredLogger.Errorw("sifchain client failed to get a block.",
					errorMessageKey, err.Error())
				continue
			}

			for _, txLog := range block.TxsResults {
				for _, event := range txLog.Events {

					claimType := getOracleClaimType(event.GetType())

					if claimType == types.ProphecyCompleted {

						prophecyInfo, err := txs.ProphecyCompletedEventToProphecyInfo(event.GetAttributes(), sub.SugaredLogger)
						if err != nil {
							sub.SugaredLogger.Errorw("sifchain client failed in get prophecy completed message from event.",
								errorMessageKey, err.Error())
							continue
						}
						if prophecyInfo.NetworkDescriptor == sub.NetworkDescriptor {
							sub.handleNewProphecyCompleted(client, uint64(blockHeight), prophecyInfo)
						}

					}
				}
			}

		}
	}
}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleNewProphecyCompleted(client *tmClient.HTTP, currentBlock uint64,
	prophecyInfo types.ProphecyInfo,
) {
	// Start Ethereum client
	ethClient, err := ethclient.Dial(sub.EthProvider)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to connect to Ethereum node.",
			errorMessageKey, err.Error())
		return
	}

	cosmosBridgeAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, sub.RegistryContractAddress, txs.CosmosBridge, sub.SugaredLogger)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the cosmos bridge address.",
			errorMessageKey, err.Error())
		return
	}

	// get the last global nonce from smart contract
	lastSubmittedNonce, err := GetLastNonceSubmitted(ethClient, cosmosBridgeAddress, sub.SugaredLogger)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the last submitted nonce.",
			errorMessageKey, err.Error())
		return
	}

	// compare the nonce in info with contract
	if prophecyInfo.GlobalNonce <= lastSubmittedNonce.Uint64() {
		sub.SugaredLogger.Errorw("incoming new global nonce less than last submitted nonce in smart contract.",
			"global nonce from cosmos event", prophecyInfo.GlobalNonce,
			"last committed nonce in smart contract", lastSubmittedNonce.Uint64())
		return
	}

	// get prophecy info list from cosmos events
	if prophecyInfo.GlobalNonce == lastSubmittedNonce.Uint64()+1 {
		sub.handleProphecyCompleted(prophecyInfo)
		return
	}

	prophecyInfoMap := sub.getProphecies(client, currentBlock, prophecyInfo.GlobalNonce, lastSubmittedNonce.Uint64())
	// put the latest one into map
	prophecyInfoMap[prophecyInfo.GlobalNonce] = prophecyInfo

	nonce := lastSubmittedNonce.Uint64() + 1
	for nonce <= prophecyInfo.GlobalNonce {
		prophecy, ok := prophecyInfoMap[nonce]
		// must deal prophecy in order, if any one missed in the map, just discontinue
		if !ok {
			sub.SugaredLogger.Errorw("can't get global nonce via scanning the blocks.",
				"get global is ", nonce)
			return
		}
		if !sub.handleProphecyCompleted(prophecy) {
			sub.SugaredLogger.Errorw("fail to process prophecy.")
			return
		}
		nonce++
	}

}

func (sub CosmosSub) getProphecies(client *tmClient.HTTP, currentBlock uint64, currentGlobalNonce uint64, lastSubmittedNonce uint64) map[uint64]types.ProphecyInfo {
	// return variable
	prophecyInfoMap := map[uint64]types.ProphecyInfo{}

	// determine the last possible block we will scan according to prophecy life time
	var endBlock uint64
	if currentBlock > ProphecyLifeTime {
		endBlock = currentBlock - ProphecyLifeTime

	} else {
		endBlock = 0
	}

	ctx := context.Background()
	// if found the next one
	found := false
	for currentBlock > endBlock {
		// not go to priviouos block if found.
		if found {
			break
		}
		currentBlock = currentBlock - 1
		tmpBlock := int64(currentBlock)
		block, err := client.BlockResults(ctx, &tmpBlock)

		if err != nil {
			sub.SugaredLogger.Errorw("sifchain client failed to get a block.",
				errorMessageKey, err.Error())
			continue
		}

		for _, txLog := range block.TxsResults {
			for _, event := range txLog.Events {

				claimType := getOracleClaimType(event.GetType())

				if claimType == types.ProphecyCompleted {

					prophecyInfo, err := txs.ProphecyCompletedEventToProphecyInfo(event.GetAttributes(), sub.SugaredLogger)
					if err != nil {
						sub.SugaredLogger.Errorw("sifchain client failed in get prophecy completed message from event.",
							errorMessageKey, err.Error())
						continue
					}
					if prophecyInfo.NetworkDescriptor == sub.NetworkDescriptor {
						// put not processed prophecy into map
						if prophecyInfo.GlobalNonce > lastSubmittedNonce {
							prophecyInfoMap[prophecyInfo.GlobalNonce] = prophecyInfo

						}
						// already found the minimum global nonce not processeed
						if prophecyInfo.GlobalNonce <= lastSubmittedNonce+1 {
							found = true
						}
					}

				}
			}
		}
	}

	return prophecyInfoMap

}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleProphecyCompleted(
	prophecyInfo types.ProphecyInfo,
) bool {
	sub.SugaredLogger.Infow(
		"get the prophecy completed message.",
		"cosmosMsg", prophecyInfo,
	)

	client, auth, target, err := tryInitRelayConfig(sub)
	if err != nil {
		sub.SugaredLogger.Errorw("failed in init relay config.",
			errorMessageKey, err.Error())
		return false
	}

	// Initialize CosmosBridge instance
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get cosmosBridge instance.",
			errorMessageKey, err.Error())
		return false
	}

	maxRetries := 5
	i := 0

	for i < maxRetries {
		err = txs.RelayProphecyCompletedToEthereum(
			prophecyInfo,
			sub.SugaredLogger,
			client,
			auth,
			cosmosBridgeInstance,
		)

		if err != nil {
			sub.SugaredLogger.Errorw(
				"failed to send new prophecy completed to ethereum",
				errorMessageKey, err.Error(),
			)
		} else {
			break
		}
		i++
	}

	if i == maxRetries {
		sub.SugaredLogger.Errorw(
			"failed to broadcast transaction after 5 attempts",
			errorMessageKey, err.Error(),
		)
		return false
	}

	return true

}
