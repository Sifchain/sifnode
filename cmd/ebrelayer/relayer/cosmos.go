package relayer

// DONTCOVER

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ebrelayertypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	CliContext              client.Context
	ValidatorName           string
	SignatureAggregator     bool
}

// NewCosmosSub initializes a new CosmosSub
func NewCosmosSub(networkDescriptor oracletypes.NetworkDescriptor, privateKey *ecdsa.PrivateKey, tmProvider, ethProvider string, registryContractAddress common.Address,
	db *leveldb.DB, cliContext client.Context, validatorName string, signatureAggregator bool, sugaredLogger *zap.SugaredLogger) CosmosSub {

	return CosmosSub{
		NetworkDescriptor:       networkDescriptor,
		TmProvider:              tmProvider,
		PrivateKey:              privateKey,
		EthProvider:             ethProvider,
		RegistryContractAddress: registryContractAddress,
		DB:                      db,
		CliContext:              cliContext,
		SignatureAggregator:     signatureAggregator,
		ValidatorName:           validatorName,
		SugaredLogger:           sugaredLogger,
	}
}

// Start a Cosmos chain subscription
func (sub CosmosSub) Start(txFactory tx.Factory, completionEvent *sync.WaitGroup) {
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
							// the relayer for signature aggregator not handle burn and lock
							if !sub.SignatureAggregator {
								cosmosMsg, err := txs.BurnLockEventToCosmosMsg(event.GetAttributes(), sub.SugaredLogger)
								if err != nil {
									sub.SugaredLogger.Errorw("sifchain client failed in get burn lock message from event.",
										errorMessageKey, err.Error())
									continue
								}
								if cosmosMsg.NetworkDescriptor == sub.NetworkDescriptor {
									sub.handleBurnLockMsg(txFactory, cosmosMsg)
								}
							}

						case types.ProphecyCompleted:
							// the relayer for signature aggregator just handle the prophecy completed
							if sub.SignatureAggregator {
								prophecyInfo, err := txs.ProphecyCompletedEventToProphecyInfo(event.GetAttributes(), sub.SugaredLogger)
								if err != nil {
									sub.SugaredLogger.Errorw("sifchain client failed in get prophecy completed message from event.",
										errorMessageKey, err.Error())
									continue
								}
								if prophecyInfo.NetworkDescriptor == sub.NetworkDescriptor {
									sub.handleProphecyCompleted(prophecyInfo)
								}
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

// MessageProcessed check if cosmogs message already processed
func MessageProcessed(prophecyID []byte, prophecyClaims []types.ProphecyClaimUnique) bool {
	for _, prophecyClaim := range prophecyClaims {
		if bytes.Compare(prophecyID, prophecyClaim.ProphecyID) == 0 {

			return true
		}
	}
	return false
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

func tryInitRelayConfig(sub CosmosSub) (*ethclient.Client, *bind.TransactOpts, common.Address, error) {

	for i := 0; i < 5; i++ {
		client, auth, target, err := txs.InitRelayConfig(
			sub.EthProvider,
			sub.RegistryContractAddress,
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
	txFactory tx.Factory,
	cosmosMsg types.CosmosMsg,
) {
	sub.SugaredLogger.Infow("handle burn lock message.",
		"cosmosMessage", cosmosMsg.String())

	sub.SugaredLogger.Infow(
		"get the prophecy claim.",
		"cosmosMsg", cosmosMsg,
	)

	valAddr, err := GetValAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		sub.SugaredLogger.Infow(
			"get the prophecy claim.",
			"cosmosMsg", err,
		)
	}

	signProphecy := ebrelayertypes.NewMsgSignProphecy(valAddr.String(), cosmosMsg.NetworkDescriptor,
		cosmosMsg.ProphecyID, "", "")

	txs.SignProphecyToCosmos(txFactory, signProphecy, sub.CliContext, sub.SugaredLogger)

}
