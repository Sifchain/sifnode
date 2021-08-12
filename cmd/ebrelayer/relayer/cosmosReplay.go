package relayer

import (
	"context"
	"log"
	"math/big"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	"go.uber.org/zap"
)

// ProphecyLiftTime signature info life time on chain
const ProphecyLiftTime = 520000

// ReplayCosmosBurnLock the missed burn lock events from cosmos
func (sub CosmosSub) ReplayCosmosBurnLock(txFactory tx.Factory) {
	client, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client, error as %s\n", err)
		return
	}

	if err := client.Start(); err != nil {
		log.Printf("failed to start a client, error as %s\n", err)
		return
	}

	fromBlock, toBlock, err := GetScannedBlockScope(client)
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

	sub.ReplayCosmosBurnLockWithBlocks(txFactory, client, accAddr, ProphecyClaims, fromBlock, toBlock)
}

// ReplayCosmosBurnLockWithBlocks replay the missed burn lock events from cosmos
func (sub CosmosSub) ReplayCosmosBurnLockWithBlocks(
	txFactory tx.Factory,
	client *tmClient.HTTP,
	accAddr sdk.AccAddress,
	ProphecyClaims []types.ProphecyClaimUnique,
	fromBlock int64,
	toBlock int64) {

	log.Printf("ReplayCosmosBurnLockWithBlocks from %d to %d block\n", fromBlock, toBlock)

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

// ReplaySignatureAggregation replay the missed signature aggregation events from cosmos
func (sub CosmosSub) ReplaySignatureAggregation(txFactory tx.Factory) {
	// start the cosmos client
	tmClient, err := tmClient.New(sub.TmProvider, "/websocket")
	if err != nil {
		log.Printf("failed to initialize a client, error as %s\n", err)
		return
	}

	if err := tmClient.Start(); err != nil {
		log.Printf("failed to start a client, error as %s\n", err)
		return
	}

	fromBlock, toBlock, err := GetScannedBlockScope(tmClient)
	if err != nil {
		log.Printf("failed to get the scaned block scope, error as %s\n", err)
		return
	}

	// Start Ethereum client
	ethClient, err := ethclient.Dial(sub.EthProvider)
	if err != nil {
		log.Printf("failed to connect to Ethereum node, error as %s\n", err)
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

	sub.ReplaySignatureAggregationWithScope(txFactory, tmClient, fromBlock, toBlock, lastSubmittedNonce.Uint64())
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

// getAllSignSigature returns all sign signature messages by check events
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

// GetScannedBlockScope get the block scope for scan
func GetScannedBlockScope(client *tmClient.HTTP) (int64, int64, error) {
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
