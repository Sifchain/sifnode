package relayer

// DONTCOVER

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
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
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	tmTypes "github.com/tendermint/tendermint/types"
)

// TODO: Move relay functionality out of CosmosSub into a new Relayer parent struct

// CosmosSub defines a Cosmos listener that relays events to Ethereum and Cosmos
type CosmosSub struct {
	TmProvider              string
	EthProvider             string
	RegistryContractAddress common.Address
	PrivateKey              *ecdsa.PrivateKey
	Logger                  tmLog.Logger
}

// NewCosmosSub initializes a new CosmosSub
func NewCosmosSub(tmProvider, ethProvider string, registryContractAddress common.Address,
	privateKey *ecdsa.PrivateKey, logger tmLog.Logger) CosmosSub {
	return CosmosSub{
		TmProvider:              tmProvider,
		EthProvider:             ethProvider,
		RegistryContractAddress: registryContractAddress,
		PrivateKey:              privateKey,
		Logger:                  logger,
	}
}

// Start a Cosmos chain subscription
func (sub CosmosSub) Start(completionEvent *sync.WaitGroup) {
	defer completionEvent.Done()
	time.Sleep(time.Second)
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.Logger.Error("failed to initialize a client", "err", err)
		completionEvent.Add(1)
		go sub.Start(completionEvent)
		return
	}
	client.SetLogger(sub.Logger)

	if err := client.Start(); err != nil {
		sub.Logger.Error("failed to start a client", "err", err)
		completionEvent.Add(1)
		go sub.Start(completionEvent)
		return
	}

	defer client.Stop() //nolint:errcheck

	// Subscribe to all tendermint transactions
	query := "tm.event = 'Tx'"
	out, err := client.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		sub.Logger.Error("failed to subscribe to query", "err", err, "query", query)
		completionEvent.Add(1)
		go sub.Start(completionEvent)
		return
	}

	defer func() {
		if err := client.Unsubscribe(context.Background(), "test", query); err != nil {
			sub.Logger.Error("Unsubscribe failed: ", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	for {
		select {
		case result := <-out:
			tx, ok := result.Data.(tmTypes.EventDataTx)
			if !ok {
				sub.Logger.Error("new tx: error while extracting event data from new tx")
			}

			sub.Logger.Info("New transaction witnessed")

			// Iterate over each event in the transaction
			for _, event := range tx.Result.Events {
				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case types.MsgBurn, types.MsgLock:
					cosmosMsg, err := txs.BurnLockEventToCosmosMsg(claimType, event.GetAttributes())
					if err != nil {
						fmt.Println(err)
						continue
					}
					// Parse event data, then package it as a ProphecyClaim and relay to the Ethereum Network
					sub.handleBurnLockMsg(cosmosMsg, claimType)
				}
			}
		case <-quit:
			return
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

			// compare method id to check if it is NewProphecyClaim method
			if bytes.Compare(tx.Data()[0:4], methodID) != 0 {
				continue
			}

			// decode data via a hardcode method since the abi unpack failed
			prophecyClaim, err := MyDecode(tx.Data()[4:])
			if err != nil {
				fmt.Printf("decode prophecy claim failed with %s \n", err.Error())
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

	src := data[64:96]
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	sequence, err := strconv.ParseUint(string(dst), 16, 32)
	if err != nil {
		fmt.Printf("Decode data failed with %s \n", err.Error())
		return types.ProphecyClaimUnique{}, err
	}
	fmt.Printf("CosmosSenderSequence is %d \n", sequence)

	// the length of sifnode acc account is 42
	cosmosSender := string(data[32*7 : 32*7+42])
	fmt.Printf("CosmosSender is %s \n", cosmosSender)

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

	fmt.Printf("found out %d prophecy claims I sent from %d to %d block", len(ProphecyClaims), ethFromBlock, ethToBlock)

	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.Logger.Error("failed to initialize a client", "err", err)
		return
	}
	client.SetLogger(sub.Logger)

	if err := client.Start(); err != nil {
		sub.Logger.Error("failed to start a client", "err", err)
		return
	}

	defer client.Stop() //nolint:errcheck

	for blockNumber := fromBlock; blockNumber < toBlock; {
		tmpBlockNumber := blockNumber
		block, err := client.BlockResults(&tmpBlockNumber)
		blockNumber++
		sub.Logger.Info(fmt.Sprintf("Replay start to process block %d", blockNumber))

		if err != nil {
			sub.Logger.Error(fmt.Sprintf("failed to start a client %s", err))
			continue
		}

		for _, log := range block.TxsResults {
			for _, event := range log.Events {

				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case types.MsgBurn, types.MsgLock:
					sub.Logger.Info(fmt.Sprintf("found out a lock burn message\n"))

					cosmosMsg, err := txs.BurnLockEventToCosmosMsg(claimType, event.GetAttributes())
					if err != nil {
						fmt.Println(err)
						continue
					}
					sub.Logger.Info(fmt.Sprintf("found out a lock burn message%s\n", cosmosMsg.String()))

					if !MessageProcessed(cosmosMsg, ProphecyClaims) {
						sub.handleBurnLockMsg(cosmosMsg, claimType)
					} else {
						sub.Logger.Info(fmt.Sprintf("lock burn message already processed by me\n"))
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
	default:
		claimType = types.Unsupported
	}
	return claimType
}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleBurnLockMsg(cosmosMsg types.CosmosMsg, claimType types.Event) {
	sub.Logger.Info(cosmosMsg.String())

	prophecyClaim := txs.CosmosMsgToProphecyClaim(cosmosMsg)
	err := txs.RelayProphecyClaimToEthereum(sub.EthProvider, sub.RegistryContractAddress,
		claimType, prophecyClaim, sub.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}
}
