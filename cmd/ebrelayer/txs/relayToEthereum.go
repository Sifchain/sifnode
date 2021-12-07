package txs

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
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
	MaxGasPrice = int64(500 * 1000000000)
)

func sleepThread(seconds time.Duration) {
	time.Sleep(time.Second * seconds)
}

// InitRelayConfig set up Ethereum client, validator's transaction auth, and the target contract's address
func InitRelayConfig(
	provider string,
	registry common.Address,
	key *ecdsa.PrivateKey,
	maxFeePerGas *big.Int,
	maxPriorityFeePerGas *big.Int,
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
	transactOptsAuth := bind.NewKeyedTransactor(key)

	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit

	// TODO now, the transaction only works with the gasPrice set.
	// need to investigate if it is a feature not supported by hardhat.
	// revert to gas price set temporarily.
	transactOptsAuth.GasPrice = big.NewInt(150000)
	// GasFeeCap is maxFeePerGas; GasTipCap is maxPriorityFeePerGas
	// transactOptsAuth.GasFeeCap = maxFeePerGas
	// transactOptsAuth.GasTipCap = maxPriorityFeePerGas
	transactOptsAuth.Context = context.Background()

	targetContract := CosmosBridge

	// Get the specific contract's address
	target, err := GetAddressFromBridgeRegistry(client, registry, targetContract, sugaredLogger)
	if err != nil {
		sugaredLogger.Errorw("failed to get cosmos bridger contract address from registry.",
			errorMessageKey, err.Error())
		return nil, nil, common.Address{}, err

	}
	return client, transactOptsAuth, target, nil
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

	batchLen := len(batchProphecyInfo)
	batchClaimData := make([]cosmosbridge.CosmosBridgeClaimData, batchLen)
	batchSignatureData := make([][]cosmosbridge.CosmosBridgeSignatureData, batchLen)
	batchID := make([][32]byte, batchLen)

	// reset the gas limit according to length of batchProphecyInfo
	auth.GasLimit = auth.GasLimit * uint64(batchLen)

	for index, prophecyInfo := range batchProphecyInfo {

		claimData := cosmosbridge.CosmosBridgeClaimData{
			CosmosSender:         []byte(prophecyInfo.CosmosSender),
			CosmosSenderSequence: big.NewInt(int64(prophecyInfo.CosmosSenderSequence)),
			EthereumReceiver:     common.HexToAddress(prophecyInfo.EthereumReceiver),
			TokenAddress:         common.HexToAddress(prophecyInfo.TokenContractAddress),
			Amount:               big.NewInt(prophecyInfo.TokenAmount.Int64()),
			DoublePeg:            prophecyInfo.DoublePeg,
			Nonce:                big.NewInt(int64(prophecyInfo.GlobalSequence)),
			NetworkDescriptor:    int32(prophecyInfo.NetworkDescriptor),
			TokenName:            prophecyInfo.TokenName,
			TokenSymbol:          prophecyInfo.TokenSymbol,
			TokenDecimals:        uint8(prophecyInfo.Decimail),
			CosmosDenom:          prophecyInfo.TokenDenomHash,
		}
		batchClaimData[index] = claimData

		var signatureData = make([]cosmosbridge.CosmosBridgeSignatureData, len(prophecyInfo.EthereumAddress))

		for index, address := range prophecyInfo.EthereumAddress {
			signature := []byte(prophecyInfo.Signatures[index])
			var r [32]byte
			var s [32]byte
			copy(r[:], signature[0:32])
			copy(s[:], signature[32:64])

			tmpSignature := cosmosbridge.CosmosBridgeSignatureData{
				Signer: common.HexToAddress(address),
				V:      signature[64] + 27,
				R:      r,
				S:      s,
			}

			signatureData[index] = tmpSignature
		}

		batchSignatureData[index] = signatureData
		var id [32]byte
		copy(id[:], prophecyInfo.ProphecyId)
		batchID[index] = id
	}

	tx, err := cosmosBridgeInstance.BatchSubmitProphecyClaimAggregatedSigs(
		auth,
		batchID,
		batchClaimData,
		batchSignatureData,
	)

	// sleep 2 seconds to wait for tx to go through before querying.
	sleepThread(2)

	if err != nil {
		sugaredLogger.Errorw(
			"cosmosBridgeInstance.BatchSubmitProphecyClaimAggregatedSigs",
			"batchID", batchID,
			"batchClaimData", batchClaimData,
			"batchSignatureData", batchSignatureData,
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
		// Get the transaction receipt
		receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())

		if err != nil {
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
