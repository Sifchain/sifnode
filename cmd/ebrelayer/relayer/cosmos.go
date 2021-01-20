package relayer

// DONTCOVER

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmKv "github.com/tendermint/tendermint/libs/kv"
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

	defer client.Unsubscribe(context.Background(), "test", query)

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
					// Parse event data, then package it as a ProphecyClaim and relay to the Ethereum Network
					err := sub.handleBurnLockMsg(event.GetAttributes(), claimType)
					if err != nil {
						sub.Logger.Error(err.Error())
					}
				}
			}
		case <-quit:
			return
		}
	}
}

func (sub CosmosSub) getAll(ethFromBlock int64, ethToBlock int64) error {
	log.Printf("getAll %d %d\n", ethFromBlock, ethToBlock)
	// Start Ethereum client
	client, err := ethclient.Dial(sub.EthProvider)
	if err != nil {
		log.Printf("%s \n", err.Error())
		return nil
	}

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		sub.Logger.Error(err.Error())
		return nil
	}
	log.Printf("clientChainID is %d \n", clientChainID)

	// Used to recover address from transaction
	eIP155Signer := ethTypes.NewEIP155Signer(big.NewInt(1))

	// Load the validator's ethereum address
	_, err = txs.LoadSender()
	if err != nil {
		log.Println(err)
		return nil
	}

	CosmosBridgeContractABI := contract.LoadABI(txs.CosmosBridge)
	methodID := CosmosBridgeContractABI.Methods[types.NewProphecyClaim.String()].ID()

	// method := types.NewProphecyClaim.String()
	fmt.Printf("method ID is %v \n", methodID)

	for blockNumber := ethFromBlock; blockNumber < ethToBlock; {
		fmt.Printf("loop blockNumber is %d \n", blockNumber)

		block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
		if err != nil {
			blockNumber++
			continue
		}
		for _, tx := range block.Transactions() {

			sender, err := eIP155Signer.Sender(tx)

			if err != nil {
				continue
			}
			fmt.Printf("sender is %s \n", sender.String())

			_, err = CosmosBridgeContractABI.MethodById(methodID)

			if err != nil {
				fmt.Println("sender print is  2")
				log.Fatal(err)
			}

			fmt.Printf("Data method part  is  %v \n", tx.Data()[:])
			if bytes.Compare(tx.Data()[0:4], methodID) != 0 {
				continue
			}

			// nolint
			type FunctionInputs struct {
				ClaimType            ethbridge.ClaimType
				CosmosSender         []byte
				CosmosSenderSequence *big.Int

				EthereumReceiver common.Address
				Symbol           string
				Amount           *big.Int
			}

			functionInputs := FunctionInputs{}

			err = CosmosBridgeContractABI.Unpack(&functionInputs, types.NewProphecyClaim.String(), tx.Data())

			if err != nil {
				fmt.Printf(" errir is %s \n", err.Error())
				continue
			}

			fmt.Printf("functionInputs is %v\n", functionInputs)
		}
		blockNumber++
	}
	return nil
}

// Replay the missed events
func (sub CosmosSub) Replay(fromBlock int64, toBlock int64, ethFromBlock int64, ethToBlock int64) {
	err := sub.getAll(ethFromBlock, ethToBlock)
	if err != nil {
		log.Fatal(err)
		return
	}

	if fromBlock > 0 {
		return
	}

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

	// TODO  junius
	// read all txs and transform to message from eth block and end eth block
	// match message to address compare with smart contract address
	// parse the data to smart contract call, the method and the arguments.
	// to check if match with msg burn, msg lock.

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
					// Parse event data, then package it as a ProphecyClaim and relay to the Ethereum Network
					err := sub.handleBurnLockMsg(event.GetAttributes(), claimType)
					if err != nil {
						sub.Logger.Error(err.Error())
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
func (sub CosmosSub) handleBurnLockMsg(attributes []tmKv.Pair, claimType types.Event) error {
	cosmosMsg, err := txs.BurnLockEventToCosmosMsg(claimType, attributes)
	if err != nil {
		fmt.Println(err)
		return err
	}
	sub.Logger.Info(cosmosMsg.String())

	prophecyClaim := txs.CosmosMsgToProphecyClaim(cosmosMsg)
	err = txs.RelayProphecyClaimToEthereum(sub.EthProvider, sub.RegistryContractAddress,
		claimType, prophecyClaim, sub.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
