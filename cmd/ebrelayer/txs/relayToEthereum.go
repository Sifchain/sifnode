package txs

// DONTCOVER

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sort"
	"time"

	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/artifacts/contracts/CosmosBridge.sol"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethereumtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const (
	// GasLimit the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit = uint64(2000000)
	// MaxGasPrice for max gas price 500 gwei
)

func sleepThread(seconds time.Duration) {
	time.Sleep(time.Second * seconds)
}

// InitRelayConfig set up Ethereum client, validator's transaction auth, and the target contract's address
func InitRelayConfig(
	provider string,
	registry common.Address,
	key *ecdsa.PrivateKey,
	maxFeePerGas,
	maxPriorityFeePerGas,
	ethereumChainID *big.Int,
	sugaredLogger *zap.SugaredLogger,
) (
	*ethclient.Client,
	*bind.TransactOpts,
	common.Address,
	error,
) {
	// Start Ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		sugaredLogger.Errorw("failed to connect ethereum node.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err
	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth, err := bind.NewKeyedTransactorWithChainID(key, ethereumChainID)
	if err != nil {
		return nil, nil, common.Address{}, err
	}
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit

	// TODO now, the transaction only works with the gasPrice set.
	//GasFeeCap is maxFeePerGas; GasTipCap is maxPriorityFeePerGas
	transactOptsAuth.GasFeeCap = maxFeePerGas
	transactOptsAuth.GasTipCap = maxPriorityFeePerGas
	transactOptsAuth.Context = context.Background()

	targetContract := CosmosBridge

	// Get the specific contract's address
	target, err := GetAddressFromBridgeRegistry(client, registry, targetContract, sugaredLogger)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmos bridger contract address from registry.",
			errorMessageKey, err.Error())
		client.Close()
		return nil, nil, common.Address{}, err

	}
	return client, transactOptsAuth, target, nil
}

type BatchUnit struct {
	prophecyId         [][32]byte
	batchClaimData     []cosmosbridge.CosmosBridgeClaimData
	batchSignatureData [][]cosmosbridge.CosmosBridgeSignatureData
}

type SignatureUnit struct {
	claimData     cosmosbridge.CosmosBridgeClaimData
	signatureData []cosmosbridge.CosmosBridgeSignatureData
	prophecyId    [32]byte
}

func prophecyInfoToSignatureUnit(prophecyInfo *oracletypes.ProphecyInfo) SignatureUnit {
	claimData := cosmosbridge.CosmosBridgeClaimData{
		CosmosSender:         []byte(prophecyInfo.CosmosSender),
		CosmosSenderSequence: big.NewInt(int64(prophecyInfo.CosmosSenderSequence)),
		EthereumReceiver:     common.HexToAddress(prophecyInfo.EthereumReceiver),
		TokenAddress:         common.HexToAddress(prophecyInfo.TokenContractAddress),
		Amount:               big.NewInt(prophecyInfo.TokenAmount.Int64()),
		BridgeToken:          prophecyInfo.BridgeToken,
		Nonce:                big.NewInt(int64(prophecyInfo.GlobalSequence)),
		NetworkDescriptor:    int32(prophecyInfo.NetworkDescriptor),
		TokenName:            prophecyInfo.TokenName,
		TokenSymbol:          prophecyInfo.TokenSymbol,
		TokenDecimals:        uint8(prophecyInfo.Decimail),
		CosmosDenom:          prophecyInfo.TokenDenomHash,
	}

	var signatureData = make([]cosmosbridge.CosmosBridgeSignatureData, len(prophecyInfo.EthereumAddress))

	for index, address := range prophecyInfo.EthereumAddress {
		signature := []byte(prophecyInfo.Signatures[index])
		var r [32]byte
		var s [32]byte
		copy(r[:], signature[0:32])
		copy(s[:], signature[32:64])

		signatureData[index] = cosmosbridge.CosmosBridgeSignatureData{
			Signer: common.HexToAddress(address),
			V:      signature[64] + 27,
			R:      r,
			S:      s,
		}
	}

	sort.Slice(signatureData, func(i, j int) bool {
		return bytes.Compare(signatureData[i].Signer[:], signatureData[j].Signer[:]) < 0
	})

	var id [32]byte
	copy(id[:], prophecyInfo.ProphecyId)

	return SignatureUnit{
		claimData:     claimData,
		signatureData: signatureData,
		prophecyId:    id,
	}
}

func buildBatchClaim(batchProphecyInfo []*oracletypes.ProphecyInfo) BatchUnit {
	batchLen := len(batchProphecyInfo)
	batchClaimData := make([]cosmosbridge.CosmosBridgeClaimData, batchLen)
	batchSignatureData := make([][]cosmosbridge.CosmosBridgeSignatureData, batchLen)
	prophecyID := make([][32]byte, batchLen)

	for index, prophecyInfo := range batchProphecyInfo {
		sUnit := prophecyInfoToSignatureUnit(prophecyInfo)
		batchClaimData[index] = sUnit.claimData
		batchSignatureData[index] = sUnit.signatureData
		prophecyID[index] = sUnit.prophecyId
	}

	return BatchUnit{
		batchClaimData:     batchClaimData,
		batchSignatureData: batchSignatureData,
		prophecyId:         prophecyID,
	}
}

// RelayBatchProphecyCompletedToEthereum send the prophecy aggregation to CosmosBridge contract on the Ethereum network
func RelayBatchProphecyCompletedToEthereum(
	batchProphecyInfo []*oracletypes.ProphecyInfo,
	sugaredLogger *zap.SugaredLogger,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	cosmosBridgeInstance *cosmosbridge.CosmosBridge,
) error {
	if len(batchProphecyInfo) == 0 {
		return nil
	}

	// reset the gas limit according to length of batchProphecyInfo
	auth.GasLimit = auth.GasLimit * uint64(len(batchProphecyInfo))

	batch := buildBatchClaim(batchProphecyInfo)

	sugaredLogger.Errorw(
		"buildBatchClaim",
		"batch", batch,
	)

	tx, err := cosmosBridgeInstance.BatchSubmitProphecyClaimAggregatedSigs(
		auth,
		batch.prophecyId,
		batch.batchClaimData,
		batch.batchSignatureData,
	)

	if err != nil {
		sugaredLogger.Errorw(
			"cosmosBridgeInstance.BatchSubmitProphecyClaimAggregatedSigs",
			"prophecyId", batch.prophecyId,
			"batchClaimData", batch.batchClaimData,
			"batchSignatureData", batch.batchSignatureData,
			errorMessageKey, err,
		)
		return err
	}

	sugaredLogger.Infow("get SubmitProphecyClaimAggregatedSigs tx hash:", "TransactionHash", tx.Hash().Hex())

	var receipt *ethereumtypes.Receipt
	maxRetries := 60
	i := 0

	// if there is an error getting the tx, or if the tx fails, retry 60 times
	for i < maxRetries {
		// sleep 2 seconds to wait for tx to go through before querying.
		sleepThread(2)

		// Get the transaction receipt
		receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())

		sugaredLogger.Debugw("Transaction receipt", "receipt", receipt.Logs)

		if err != nil {
			sugaredLogger.Errorw("Failed to submit to ethereum client", "error", err)
			sleepThread(1)
		} else {
			break
		}
		i++
	}

	if i == maxRetries {
		return errors.New("hit max tx receipt query retries")
	}

	sugaredLogger.Infow(
		"Successfully received transaction receipt after retry",
		"txReceipt", receipt,
	)

	return nil
}
