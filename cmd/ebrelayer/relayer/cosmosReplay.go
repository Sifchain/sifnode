package relayer

import (
	"bytes"
	"context"
	"log"
	"math/big"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
)

// ReplayBurnLock the missed burn lock events
func (sub CosmosSub) ReplayBurnLock(txFactory tx.Factory, fromBlock int64, toBlock int64) {
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client, error as %s\n", err)
		return
	}

	if err := client.Start(); err != nil {
		log.Printf("failed to start a client, error as %s\n", err)
		return
	}

	accAddr, err := GetAccAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		log.Printf("failed to get the account address, error as %s\n", err)
		return
	}

	ProphecyClaims := sub.getAllSignSigature(accAddr, fromBlock, toBlock)

	log.Printf("found out %d prophecy claims I sent from %d to %d block\n", len(ProphecyClaims), fromBlock, toBlock)

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

					cosmosMsg, err := txs.BurnLockEventToCosmosMsg(event.GetAttributes(), sub.SugaredLogger)
					if err != nil {
						log.Println(err)
						continue
					}
					log.Printf("found out a lock burn message%s\n", cosmosMsg.String())
					if cosmosMsg.NetworkDescriptor == sub.NetworkDescriptor {
						if !MessageProcessed(cosmosMsg.ProphecyID, ProphecyClaims) {
							sub.handleBurnLockMsg(txFactory, cosmosMsg)
						} else {
							log.Println("lock burn message already processed by me")
						}
					}
				}
			}
		}
	}
}

// ReplaySignatureAggregation to check missed ProphecyCompleted events
func (sub CosmosSub) ReplaySignatureAggregation(txFactory tx.Factory, fromBlock int64, toBlock int64, ethFromBlock int64, ethToBlock int64) {
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

	// scan cosmos blocks
	for blockNumber := fromBlock; blockNumber < toBlock; {
		tmpBlockNumber := blockNumber

		ctx := context.Background()
		block, err := client.BlockResults(ctx, &tmpBlockNumber)

		if err != nil {
			log.Printf("failed to start a client %s\n", err.Error())
			continue
		}

		blockNumber++
		log.Printf("Replay start to process block %d\n", blockNumber)

		for _, ethLog := range block.TxsResults {
			for _, event := range ethLog.Events {

				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case types.ProphecyCompleted:

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
}

// GetAllProphecyClaim get all prophecy claims from smart contract
func GetAllProphecyClaim(client *ethclient.Client, ethereumAddress common.Address, ethFromBlock int64, ethToBlock int64) []types.ProphecyClaimUnique {
	log.Printf("getAllProphecyClaim from %d block to %d block\n", ethFromBlock, ethToBlock)

	var prophecyClaimArray []types.ProphecyClaimUnique

	// Used to recover address from transaction, the clientChainID doesn't work in ganache, hardcoded to 1
	eIP155Signer := ethTypes.NewEIP155Signer(big.NewInt(1))

	CosmosBridgeContractABI := contract.LoadABI(txs.CosmosBridge)
	methodID := CosmosBridgeContractABI.Methods[types.SubmitProphecyClaimAggregatedSigs.String()].ID()

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

// getAllSignSigature
func (sub CosmosSub) getAllSignSigature(accAddress sdk.AccAddress, fromBlock int64, toBlock int64) []types.ProphecyClaimUnique {
	log.Printf("Replay get all ethereum bridge claim from block %d to block %d\n", fromBlock, toBlock)

	var claimArray []types.ProphecyClaimUnique
	tmClient, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a cosmos client, error is %s\n", err.Error())
		return claimArray
	}

	if err := tmClient.Start(); err != nil {
		log.Printf("failed to start a cosmos client, error is %s\n", err.Error())
		return claimArray
	}

	defer tmClient.Stop() //nolint:errcheck

	for blockNumber := fromBlock; blockNumber < toBlock; {
		tmpBlockNumber := blockNumber

		ctx := context.Background()
		block, err := tmClient.BlockResults(ctx, &tmpBlockNumber)

		blockNumber++
		log.Printf("Replay start to process block %d\n", blockNumber)

		if err != nil {
			log.Printf("failed to get the block %s\n", err.Error())
			continue
		}

		for _, result := range block.TxsResults {
			for _, event := range result.Events {
				log.Printf("Replay get an event %s\n", event.GetType())
				if event.GetType() == "sign_prophecy" {
					claim, err := txs.AttributesToCosmosSignProphecyClaim(event.GetAttributes())
					if err != nil {
						continue
					}

					// Check if sender is me
					if claim.CosmosSender.Equals(accAddress) {
						claimArray = append(claimArray, types.ProphecyClaimUnique{
							ProphecyID: claim.ProphecyID,
						})
					}
				}
			}
		}
	}

	return claimArray
}

// GetAccAddressFromKeyring get the address from key ring and keyname
func GetAccAddressFromKeyring(k keyring.Keyring, keyname string) (sdk.AccAddress, error) {
	keyInfo, err := k.Key(keyname)
	if err != nil {
		return nil, err
	}
	return keyInfo.GetAddress(), nil
}
