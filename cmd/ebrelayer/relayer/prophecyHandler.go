package relayer

import (
	"context"
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
)

const wakeupTimer = 60

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

	t := time.NewTicker(time.Second * wakeupTimer)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	for {
		select {
		case <-quit:
			sub.SugaredLogger.Warn("we receive the quit signal and exit")
			return

		case <-t.C:
			sub.SugaredLogger.Info("timer triggered, start to check cosmos message")
			sub.handleNewProphecyCompleted(client)
		}
	}
}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleNewProphecyCompleted(client *tmClient.HTTP) {
	ctx := context.Background()
	currentBlock, err := client.Block(ctx, nil)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the latest block from cosmos.",
			errorMessageKey, err.Error())
		return
	}

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

	// put the latest one into map
	maxGlobalNonce, prophecyInfoMap := sub.getUnhandledProphecies(client, uint64(currentBlock.Block.Height), lastSubmittedNonce.Uint64())

	nonce := lastSubmittedNonce.Uint64() + 1
	batchStartNonce := nonce
	batchEndNonce := lastSubmittedNonce.Uint64()

	for nonce <= maxGlobalNonce {
		prophecy, ok := prophecyInfoMap[nonce]
		// must deal prophecy in order, if any one missed in the map, just discontinue
		if !ok {
			sub.SugaredLogger.Errorw("can't get global nonce via scanning the blocks.",
				"expected global nonce is ", nonce)

			if !sub.handleBatchProphecyCompleted(prophecyInfoMap, batchStartNonce, batchEndNonce) {
				sub.SugaredLogger.Errorw("fail to process prophecy.")
				return
			}
			return
		}
		batchEndNonce = prophecy.GlobalNonce
		if batchEndNonce > batchStartNonce+5 {
			if !sub.handleBatchProphecyCompleted(prophecyInfoMap, batchStartNonce, batchEndNonce) {
				sub.SugaredLogger.Errorw("fail to process prophecy.")
				return
			}
			batchStartNonce = batchEndNonce + 1
		}
		nonce++
	}

	sub.handleBatchProphecyCompleted(prophecyInfoMap, batchStartNonce, batchEndNonce)
}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleBatchProphecyCompleted(
	prophecyInfoMap map[uint64]types.ProphecyInfo,
	batchStartNonce uint64,
	batchEndNonce uint64) bool {

	sub.SugaredLogger.Infow(
		"handle batch prophecy completed.",
		"cosmosMsg", prophecyInfoMap,
	)

	var batchProphecyInfo []types.ProphecyInfo
	for batchStartNonce <= batchEndNonce {
		batchProphecyInfo = append(batchProphecyInfo, prophecyInfoMap[batchStartNonce])
		batchStartNonce++
	}

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
		err = txs.RelayBatchProphecyCompletedToEthereum(
			batchProphecyInfo,
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

	// process successfully complete
	return true

}

func (sub CosmosSub) getUnhandledProphecies(client *tmClient.HTTP, currentBlock uint64, lastSubmittedNonce uint64) (uint64, map[uint64]types.ProphecyInfo) {
	// return variable
	prophecyInfoMap := map[uint64]types.ProphecyInfo{}
	maxGlobalNonce := uint64(0)

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
							if prophecyInfo.GlobalNonce > maxGlobalNonce {
								// record the max global nonce
								maxGlobalNonce = prophecyInfo.GlobalNonce

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
		
	}
	return maxGlobalNonce, prophecyInfoMap
}
