package relayer

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/artifacts/contracts/CosmosBridge.sol"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ethereum/go-ethereum/ethclient"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"

	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/instrumentation"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const wakeupTimer = 60

// StartProphecyHandler start Cosmos chain subscription and process prophecy completed message
func (sub CosmosSub) StartProphecyHandler(txFactory tx.Factory, completionEvent *sync.WaitGroup, symbolTranslator *symbol_translator.SymbolTranslator) {
	defer completionEvent.Done()
	time.Sleep(time.Second)
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		// completionEvent.Add(1)
		// go sub.Start(txFactory, completionEvent, symbolTranslator)
		return
	}

	if err := client.Start(); err != nil {
		sub.SugaredLogger.Errorw("failed to start a sifchain client.",
			errorMessageKey, err.Error())
		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent, symbolTranslator)
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

// Get all not processed Prophecy via rpc and handle them in batch
func (sub CosmosSub) handleNewProphecyCompleted(client *tmClient.HTTP) {
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

	prophecyInfoArray := GetAllProphciesCompleted(sub.TmProvider, sub.NetworkDescriptor, lastSubmittedNonce.Uint64()+1)

	batches := (len(prophecyInfoArray) + 4) / 5

	for batch := 0; batch < batches; batch++ {
		// for batch < batches {
		end := (batch + 1) * 5
		if end > len(prophecyInfoArray) {
			end = len(prophecyInfoArray)
		}

		batchProphecyInfo := prophecyInfoArray[batch*5 : end]

		if !sub.handleBatchProphecyCompleted(batchProphecyInfo) {
			sub.SugaredLogger.Errorw("fail to process prophecy.")
			return
		}
	}
}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleBatchProphecyCompleted(
	batchProphecyInfo []*oracletypes.ProphecyInfo) bool {

	sub.SugaredLogger.Infow(
		"handle batch prophecy completed.",
		"prophecyInfoArray", batchProphecyInfo,
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
		return false
	}

	// process successfully complete
	instrumentation.PeggyCheckpointZap(sub.SugaredLogger, instrumentation.ProphecyClaimSubmitted)
	return true

}

// GetAllProphciesCompleted usage
// 1. Call ethereum and get lastNonceSubmitted
// 2. Call this function with the lastNonceSubmitted on ethereum side
// 3. This function returns all of the prophecies that need to be relayed from sifchain to that EVM chain
// TODO add a limit of maximum of n prophecies to query for
func GetAllProphciesCompleted(rpcServer string, networkDescriptor oracletypes.NetworkDescriptor, startGlobalSequence uint64) []*oracletypes.ProphecyInfo {
	conn, err := grpc.Dial("0.0.0.0:9090", grpc.WithInsecure())
	if err != nil {
		return []*oracletypes.ProphecyInfo{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := ethbridgetypes.NewQueryClient(conn)
	request := ethbridgetypes.QueryProphciesCompletedRequest{
		NetworkDescriptor: networkDescriptor,
		GlobalSequence:    startGlobalSequence,
	}
	response, err := client.ProphciesCompleted(ctx, &request)
	if err != nil {
		return []*oracletypes.ProphecyInfo{}
	}
	return response.ProphecyInfo
}
