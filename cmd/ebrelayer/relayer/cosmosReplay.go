package relayer

import (
	"bytes"
	"context"
	"log"
	"math/big"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"go.uber.org/zap"
)

// ProphecyLiftTime signature info life time on chain
const ProphecyLiftTime = 520000

// ReplayBurnLock the missed burn lock events
func (sub CosmosSub) ReplayBurnLock(txFactory tx.Factory) {
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client, error as %s\n", err)
		return
	}

	if err := client.Start(); err != nil {
		log.Printf("failed to start a client, error as %s\n", err)
		return
	}

	fromBlock, toBlock, err := GetScanBlockScope(client)
	if err != nil {
		log.Printf("failed to get the scaned block scope, error as %s\n", err)
		return
	}

	accAddr, err := GetAccAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		log.Printf("failed to get the account address, error as %s\n", err)
		return
	}

	ProphecyClaims := sub.getAllSignSigature(client, accAddr, fromBlock, toBlock)

	sub.ReplayBurnLockWithBlocks(txFactory, client, accAddr, ProphecyClaims, fromBlock, toBlock)
}

// ReplayBurnLockWithBlocks replay the missed burn lock events
func (sub CosmosSub) ReplayBurnLockWithBlocks(
	txFactory tx.Factory,
	client *tmClient.HTTP,
	accAddr sdk.AccAddress,
	ProphecyClaims []types.ProphecyClaimUnique,
	fromBlock int64,
	toBlock int64) {

	log.Printf("ReplayBurnLockWithBlocks from %d to %d block\n", fromBlock, toBlock)

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

// ReplaySignatureAggregation replay the missed signature aggregation events
func (sub CosmosSub) ReplaySignatureAggregation(txFactory tx.Factory) {
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client, error as %s\n", err)
		return
	}

	if err := client.Start(); err != nil {
		log.Printf("failed to start a client, error as %s\n", err)
		return
	}

	fromBlock, toBlock, err := GetScanBlockScope(client)
	if err != nil {
		log.Printf("failed to get the scaned block scope, error as %s\n", err)
		return
	}

	// Start Ethereum client
	ethClient, err := ethclient.Dial(sub.EthProvider)
	if err != nil {
		log.Printf("%s \n", err.Error())
		return
	}

	cosmosBridgeAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, sub.RegistryContractAddress, txs.CosmosBridge, sub.SugaredLogger)
	if err != nil {
		log.Printf("failed to get the cosmos bridge address, error as %s\n", err)
		return
	}

	lastSubmittedNonce, err := GetLastNonceSubmitted(ethClient, cosmosBridgeAddress, sub.SugaredLogger)
	if err != nil {
		log.Printf("failed to get the last submitted nonce, error as %s\n", err)
		return
	}

	if err != nil {
		log.Printf("failed to get the last submitted nonce, error as %s\n", err)
		return
	}

	sub.ReplaySignatureAggregationWithScope(txFactory, client, fromBlock, toBlock, lastSubmittedNonce.Uint64())

}

// ReplaySignatureAggregationWithScope to check missed ProphecyCompleted events
func (sub CosmosSub) ReplaySignatureAggregationWithScope(txFactory tx.Factory, client *tmClient.HTTP, fromBlock int64, toBlock int64, lastSubmittedNonce uint64) {

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
					if prophecyInfo.NetworkDescriptor == sub.NetworkDescriptor &&
						prophecyInfo.GlobalNonce > lastSubmittedNonce {
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
func (sub CosmosSub) getAllSignSigature(client *tmClient.HTTP, accAddress sdk.AccAddress, fromBlock int64, toBlock int64) []types.ProphecyClaimUnique {
	log.Printf("Replay get all ethereum bridge claim from block %d to block %d\n", fromBlock, toBlock)
	claims := []types.ProphecyClaimUnique{}

	for blockNumber := fromBlock; blockNumber < toBlock; {
		tmpBlockNumber := blockNumber

		ctx := context.Background()
		block, err := client.BlockResults(ctx, &tmpBlockNumber)

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
						claims = append(claims, types.ProphecyClaimUnique{
							ProphecyID: claim.ProphecyID,
						})
					}
				}
			}
		}
	}

	return claims
}

// GetLastNonceSubmitted get last nonce submitted in cosmos bridge contract
func GetLastNonceSubmitted(client *ethclient.Client, cosmosBridgeAddress common.Address, sugaredLogger *zap.SugaredLogger) (*big.Int, error) {

	// Initialize CosmosBridge instance
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(cosmosBridgeAddress, client)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmosBridge instance.",
			errorMessageKey, err.Error())
		return nil, err
	}
	return cosmosBridgeInstance.LastNonceSubmitted(nil)
}

// GetScanBlockScope get the block scope for scan
func GetScanBlockScope(client *tmClient.HTTP) (int64, int64, error) {
	currentBlock, err := client.BlockResults(context.Background(), nil)
	if err != nil {
		return 0, 0, err
	}
	toBlock := currentBlock.Height
	var fromBlock int64
	if toBlock > ProphecyLiftTime {
		fromBlock = toBlock - ProphecyLiftTime
	} else {
		fromBlock = 0
	}
	return fromBlock, toBlock, nil
}

// GetAccAddressFromKeyring get the address from key ring and keyname
func GetAccAddressFromKeyring(k keyring.Keyring, keyname string) (sdk.AccAddress, error) {
	keyInfo, err := k.Key(keyname)
	if err != nil {
		return nil, err
	}
	return keyInfo.GetAddress(), nil
}
