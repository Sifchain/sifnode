package relayer

// DONTCOVER

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/syndtr/goleveldb/leveldb"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	tmTypes "github.com/tendermint/tendermint/types"
	"go.uber.org/zap"
)

const (
	cosmosLevelDBKey = "cosmosLastProcessedBlock"
)

// TODO: Move relay functionality out of CosmosSub into a new Relayer parent struct
const errorMessageKey = "errorMessage"

// CosmosSub defines a Cosmos listener that relays events to Ethereum and Cosmos
type CosmosSub struct {
	TmProvider              string
	EthProvider             string
	PrivateKey              *ecdsa.PrivateKey
	DB                      *leveldb.DB
	SugaredLogger           *zap.SugaredLogger
	NetworkDescriptor       oracletypes.NetworkDescriptor
	RegistryContractAddress common.Address
}

// NewCosmosSub initializes a new CosmosSub
func NewCosmosSub(networkDescriptor oracletypes.NetworkDescriptor, privateKey *ecdsa.PrivateKey, tmProvider, ethProvider string, registryContractAddress common.Address,
	db *leveldb.DB, sugaredLogger *zap.SugaredLogger) CosmosSub {

	return CosmosSub{
		NetworkDescriptor:       networkDescriptor,
		TmProvider:              tmProvider,
		PrivateKey:              privateKey,
		EthProvider:             ethProvider,
		RegistryContractAddress: registryContractAddress,
		DB:                      db,
		SugaredLogger:           sugaredLogger,
	}
}

// Start a Cosmos chain subscription
func (sub CosmosSub) Start(completionEvent *sync.WaitGroup) {
	defer completionEvent.Done()
	time.Sleep(time.Second)
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		completionEvent.Add(1)
		go sub.Start(completionEvent)
		return
	}

	if err := client.Start(); err != nil {
		sub.SugaredLogger.Errorw("failed to start a sifchain client.",
			errorMessageKey, err.Error())
		completionEvent.Add(1)
		go sub.Start(completionEvent)
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
		go sub.Start(completionEvent)
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

	var lastProcessedBlock int64

	data, err := sub.DB.Get([]byte(cosmosLevelDBKey), nil)
	if err != nil {
		log.Println("Error getting the last cosmos block from level db", err)
		lastProcessedBlock = 0
	} else {
		lastProcessedBlock = new(big.Int).SetBytes(data).Int64()
	}

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

			// Just start from current block number if never process any block before
			if lastProcessedBlock == 0 {
				lastProcessedBlock = blockHeight
			}
			sub.SugaredLogger.Infow("new transaction witnessed in sifchain client.")

			startBlockHeight := lastProcessedBlock + 1
			sub.SugaredLogger.Infow("cosmos process events for blocks.",
				"startingBlockHeight", startBlockHeight, "currentBlockHeight", blockHeight)

			for blockNumber := startBlockHeight; blockNumber <= blockHeight; {
				tmpBlockNumber := blockNumber

				ctx := context.Background()
				block, err := client.BlockResults(ctx, &tmpBlockNumber)

				if err != nil {
					sub.SugaredLogger.Errorw("sifchain client failed to get a block.",
						errorMessageKey, err.Error())
					continue
				}

				for _, txLog := range block.TxsResults {
					for _, event := range txLog.Events {

						claimType := getOracleClaimType(event.GetType())

						switch claimType {
						case types.MsgBurn, types.MsgLock:
							cosmosMsg, err := txs.BurnLockEventToCosmosMsg(claimType, event.GetAttributes(), sub.SugaredLogger)
							if err != nil {
								sub.SugaredLogger.Errorw("sifchain client failed in get message from event.",
									errorMessageKey, err.Error())
								continue
							}
							if cosmosMsg.NetworkDescriptor == sub.NetworkDescriptor {
								sub.handleBurnLockMsg(cosmosMsg, claimType)
							}
						case types.ProphecyCompleted:
							prophecyInfo, err := txs.ProphecyCompletedEventToProphecyInfo(claimType, event.GetAttributes(), sub.SugaredLogger)
							if err != nil {
								sub.SugaredLogger.Errorw("sifchain client failed in get message from event.",
									errorMessageKey, err.Error())
								continue
							}
							if prophecyInfo.NetworkDescriptor == sub.NetworkDescriptor {
								sub.handleProphecyCompleted(prophecyInfo, claimType)
							}

						}

					}
				}

				lastProcessedBlock = blockNumber
				err = sub.DB.Put([]byte(cosmosLevelDBKey), big.NewInt(lastProcessedBlock).Bytes(), nil)
				if err != nil {
					// if you can't write to leveldb, then error out as something is seriously amiss
					log.Fatalf("Error saving lastProcessedBlock to leveldb: %v", err)
				}
				blockNumber++
			}
		}
	}
}

// GetAllProphecyClaim get all prophecy claims
func GetAllProphecyClaim(client *ethclient.Client, ethereumAddress common.Address, ethFromBlock int64, ethToBlock int64) []types.ProphecyClaimUnique {
	log.Printf("getAllProphecyClaim from %d block to %d block\n", ethFromBlock, ethToBlock)

	var prophecyClaimArray []types.ProphecyClaimUnique

	// Used to recover address from transaction, the clientChainID doesn't work in ganache, hardcoded to 1
	eIP155Signer := ethTypes.NewEIP155Signer(big.NewInt(1))

	CosmosBridgeContractABI := contract.LoadABI(txs.CosmosBridge)
	methodID := CosmosBridgeContractABI.Methods[types.NewProphecyClaim.String()].ID()

	for blockNumber := ethFromBlock; blockNumber < ethToBlock; {
		log.Printf("getAllProphecyClaim current blockNumber is %d\n", blockNumber)

		block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
		if err != nil {
			log.Printf("failed to get block from ethereum, block number is %d\n", blockNumber)
			blockNumber++
			continue
		}

		for _, tx := range block.Transactions() {
			// recover sender from tx
			sender, err := eIP155Signer.Sender(tx)
			if err != nil {
				log.Println("failed to recover sender from tx")
				continue
			}

			// compare tx sender with my ethereum account
			if sender != ethereumAddress {
				// the prophecy claim not sent by me
				continue
			}

			if len(tx.Data()) < 4 {
				log.Println("the tx is not a smart contract call")
				continue
			}

			// compare method id to check if it is NewProphecyClaim method
			if bytes.Compare(tx.Data()[0:4], methodID) != 0 {
				continue
			}

			// decode data via a hardcode method since the abi unpack failed
			prophecyClaim, err := MyDecode(tx.Data()[4:])
			if err != nil {
				log.Printf("decode prophecy claim failed with %s \n", err.Error())
				continue
			}

			// put matched prophecyClaim into result
			prophecyClaimArray = append(prophecyClaimArray, prophecyClaim)
		}
		blockNumber++
	}
	return prophecyClaimArray
}

// MyDecode decode data in ProphecyClaim transaction
func MyDecode(data []byte) (types.ProphecyClaimUnique, error) {
	if len(data) < 32*7+42 {
		return types.ProphecyClaimUnique{}, errors.New("tx data length not enough")
	}

	src := data[64:96]
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	sequence, err := strconv.ParseUint(string(dst), 16, 32)
	if err != nil {
		return types.ProphecyClaimUnique{}, err
	}

	// the length of sifnode acc account is 42

	return types.ProphecyClaimUnique{
		CosmosSenderSequence: big.NewInt(int64(sequence)),
		CosmosSender:         data[32*7 : 32*7+42],
	}, nil
}

// MessageProcessed check if cosmogs message already processed
func MessageProcessed(message types.CosmosMsg, prophecyClaims []types.ProphecyClaimUnique) bool {
	for _, prophecyClaim := range prophecyClaims {
		if bytes.Compare(message.CosmosSender, prophecyClaim.CosmosSender) == 0 &&
			message.CosmosSenderSequence.Cmp(prophecyClaim.CosmosSenderSequence) == 0 {
			return true
		}
	}
	return false
}

// Replay the missed events
func (sub CosmosSub) Replay(fromBlock int64, toBlock int64, ethFromBlock int64, ethToBlock int64) {
	// Start Ethereum client
	ethClient, err := ethclient.Dial(sub.EthProvider)
	if err != nil {
		log.Printf("%s \n", err.Error())
		return
	}

	clientChainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		log.Printf("%s \n", err.Error())
		return
	}
	log.Printf("clientChainID is %d \n", clientChainID)

	// Load the validator's ethereum address
	mySender, err := txs.LoadSender()
	if err != nil {
		log.Println(err)
		return
	}

	ProphecyClaims := GetAllProphecyClaim(ethClient, mySender, ethFromBlock, ethToBlock)

	log.Printf("found out %d prophecy claims I sent from %d to %d block\n", len(ProphecyClaims), ethFromBlock, ethToBlock)

	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client, error as %s\n", err)
		return
	}

	if err := client.Start(); err != nil {
		log.Printf("failed to start a client, error as %s\n", err)
		return
	}

	defer client.Stop() //nolint:errcheck

	for blockNumber := fromBlock; blockNumber < toBlock; {
		tmpBlockNumber := blockNumber

		ctx := context.Background()
		block, err := client.BlockResults(ctx, &tmpBlockNumber)

		blockNumber++
		log.Printf("Replay start to process block %d\n", blockNumber)

		if err != nil {
			log.Printf("failed to start a client %s\n", err.Error())
			continue
		}

		for _, ethLog := range block.TxsResults {
			for _, event := range ethLog.Events {

				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case types.MsgBurn, types.MsgLock:
					log.Println("found out a lock burn message")

					cosmosMsg, err := txs.BurnLockEventToCosmosMsg(claimType, event.GetAttributes(), sub.SugaredLogger)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Printf("found out a lock burn message%s\n", cosmosMsg.String())
					if cosmosMsg.NetworkDescriptor == sub.NetworkDescriptor {
						if !MessageProcessed(cosmosMsg, ProphecyClaims) {
							sub.handleBurnLockMsg(cosmosMsg, claimType)
						} else {
							log.Println("lock burn message already processed by me")
						}
					}
				}
			}
		}
	}
}

// getOracleClaimType sets the OracleClaim's claim type based upon the witnessed event type
func getOracleClaimType(eventType string) types.Event {
	var claimType types.Event
	switch eventType {
	case types.MsgBurn.String():
		claimType = types.MsgBurn
	case types.MsgLock.String():
		claimType = types.MsgLock
	case types.ProphecyCompleted.String():
		claimType = types.ProphecyCompleted
	default:
		claimType = types.Unsupported
	}
	return claimType
}

func tryInitRelayConfig(sub CosmosSub, claimType types.Event) (*ethclient.Client, *bind.TransactOpts, common.Address, error) {

	for i := 0; i < 5; i++ {
		client, auth, target, err := txs.InitRelayConfig(
			sub.EthProvider,
			sub.RegistryContractAddress,
			claimType,
			sub.PrivateKey,
			sub.SugaredLogger,
		)

		if err != nil {
			sub.SugaredLogger.Errorw("failed in init relay config.",
				errorMessageKey, err.Error())
			continue
		}
		return client, auth, target, err
	}

	return nil, nil, common.Address{}, errors.New("hit max initRelayConfig retries")
}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleBurnLockMsg(
	cosmosMsg types.CosmosMsg,
	claimType types.Event,
) {
	sub.SugaredLogger.Infow("handle burn lock message.",
		"cosmosMessage", cosmosMsg.String())

	sub.SugaredLogger.Infow(
		"get the prophecy claim.",
		"cosmosMsg", cosmosMsg,
	)

	client, auth, target, err := tryInitRelayConfig(sub, claimType)
	if err != nil {
		sub.SugaredLogger.Errorw("failed in init relay config.",
			errorMessageKey, err.Error())
		return
	}

	// Initialize CosmosBridge instance
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get cosmosBridge instance.",
			errorMessageKey, err.Error())
		return
	}

	maxRetries := 5
	i := 0

	for i < maxRetries {
		err = txs.RelayProphecyClaimToEthereum(
			cosmosMsg,
			sub.SugaredLogger,
			client,
			auth,
			cosmosBridgeInstance,
		)

		if err != nil {
			sub.SugaredLogger.Errorw(
				"failed to send new prophecyclaim to ethereum",
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
	}
}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleProphecyCompleted(
	prophecyInfo types.ProphecyInfo,
	claimType types.Event,
) {
	sub.SugaredLogger.Infow("handle burn lock message.",
		"cosmosMessage", prophecyInfo)

	sub.SugaredLogger.Infow(
		"get the prophecy claim.",
		"cosmosMsg", prophecyInfo,
	)

	client, auth, target, err := tryInitRelayConfig(sub, claimType)
	if err != nil {
		sub.SugaredLogger.Errorw("failed in init relay config.",
			errorMessageKey, err.Error())
		return
	}

	// Initialize CosmosBridge instance
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get cosmosBridge instance.",
			errorMessageKey, err.Error())
		return
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
				"failed to send new prophecyclaim to ethereum",
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
	}
}
