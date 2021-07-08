package main

import (
	"bufio"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
)

// RunReplayEthereumCmd executes replayEthereumCmd
func RunReplayEthereumCmd(cmd *cobra.Command, args []string) error {
	// Load the validator's Ethereum private key from environment variables
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

	// Parse flag --chain-id
	chainID := viper.GetString(flags.FlagChainID)
	if strings.TrimSpace(chainID) == "" {
		return errors.Errorf("Must specify a 'chain-id'")
	}

	// Parse flag --rpc-url
	rpcURL := viper.GetString(FlagRPCURL)
	if rpcURL != "" {
		_, err := url.Parse(rpcURL)
		if rpcURL != "" && err != nil {
			return errors.Wrapf(err, "invalid RPC URL: %v", rpcURL)
		}
	}

	tendermintNode := args[0]
	web3Provider := args[1]

	if !common.IsHexAddress(args[2]) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", args[1])
	}
	contractAddress := common.HexToAddress(args[2])

	if len(strings.Trim(args[3], "")) == 0 {
		return errors.Errorf("invalid [validator-moniker]: %s", args[2])
	}
	validatorMoniker := args[3]
	mnemonic := args[4]

	fromBlock, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [from-block]: %s", args[5])
	}

	toBlock, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [to-block]: %s", args[6])
	}

	cosmosFromBlock, err := strconv.ParseInt(args[7], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [from-block]: %s", args[7])
	}

	cosmosToBlock, err := strconv.ParseInt(args[8], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [to-block]: %s", args[8])
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}
	sugaredLogger := logger.Sugar()

	// Initialize new Ethereum event listener
	inBuf := bufio.NewReader(cmd.InOrStdin())

	ethSub, err := relayer.NewEthereumSub(inBuf, tendermintNode, cdc, validatorMoniker, chainID, web3Provider,
		contractAddress, privateKey, mnemonic, nil, sugaredLogger)
	if err != nil {
		return err
	}

	ethSub.Replay(fromBlock, toBlock, cosmosFromBlock, cosmosToBlock)

	return nil
}

// RunReplayCosmosCmd executes initRelayerCmd
func RunReplayCosmosCmd(cmd *cobra.Command, args []string) error {
	// Load the validator's Ethereum private key from environment variables
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

	// Parse flag --rpc-url
	rpcURL := viper.GetString(FlagRPCURL)
	if rpcURL != "" {
		_, err := url.Parse(rpcURL)
		if rpcURL != "" && err != nil {
			return errors.Wrapf(err, "invalid RPC URL: %v", rpcURL)
		}
	}

	// Validate and parse arguments
	if len(strings.Trim(args[0], "")) == 0 {
		return errors.Errorf("invalid [tendermint-node]: %s", args[0])
	}
	tendermintNode := args[0]

	if !relayer.IsWebsocketURL(args[1]) {
		return errors.Errorf("invalid [web3-provider]: %s", args[1])
	}
	web3Provider := args[1]

	if !common.IsHexAddress(args[2]) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", args[2])
	}
	contractAddress := common.HexToAddress(args[2])

	fromBlock, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [from-block]: %s", args[3])
	}

	toBlock, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [to-block]: %s", args[4])
	}

	ethFromBlock, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [eth-from-block]: %s", args[3])
	}

	ethToBlock, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [eth-to-block]: %s", args[4])
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}
	sugaredLogger := logger.Sugar()

	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(tendermintNode, web3Provider, contractAddress, privateKey, nil, sugaredLogger)

	cosmosSub.Replay(fromBlock, toBlock, ethFromBlock, ethToBlock)

	return nil
}

// RunListMissedCosmosEventCmd executes initRelayerCmd
func RunListMissedCosmosEventCmd(cmd *cobra.Command, args []string) error {
	// Parse flag --rpc-url
	rpcURL := viper.GetString(FlagRPCURL)
	if rpcURL != "" {
		_, err := url.Parse(rpcURL)
		if rpcURL != "" && err != nil {
			return errors.Wrapf(err, "invalid RPC URL: %v", rpcURL)
		}
	}

	// Validate and parse arguments
	if len(strings.Trim(args[0], "")) == 0 {
		return errors.Errorf("invalid [tendermint-node]: %s", args[0])
	}
	tendermintNode := args[0]

	if !relayer.IsWebsocketURL(args[1]) {
		return errors.Errorf("invalid [web3-provider]: %s", args[1])
	}
	web3Provider := args[1]

	if !common.IsHexAddress(args[2]) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", args[2])
	}
	contractAddress := common.HexToAddress(args[2])

	if !common.IsHexAddress(args[3]) {
		return errors.Errorf("invalid [relayer-ethereum-address]: %s", args[3])
	}
	relayerEthereumAddress := common.HexToAddress(args[3])

	days, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [days]: %s", args[3])
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}
	sugaredLogger := logger.Sugar()

	// Initialize new Cosmos event listener
	listMissedCosmosEvent := relayer.NewListMissedCosmosEvent(tendermintNode, web3Provider, contractAddress, relayerEthereumAddress, days, sugaredLogger)

	listMissedCosmosEvent.ListMissedCosmosEvent()

	return nil
}

// RunStartGasOracleCmd executes startGasOracleCmd
func RunStartGasOracleCmd(cmd *cobra.Command, args []string) error {
	// Load the validator's Ethereum private key from environment variables
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

	// Parse flag --chain-id
	chainID := viper.GetString(flags.FlagChainID)
	if strings.TrimSpace(chainID) == "" {
		return errors.Errorf("Must specify a 'chain-id'")
	}

	// Parse flag --rpc-url
	rpcURL := viper.GetString(FlagRPCURL)
	if rpcURL != "" {
		_, err := url.Parse(rpcURL)
		if rpcURL != "" && err != nil {
			return errors.Wrapf(err, "invalid RPC URL: %v", rpcURL)
		}
	}

	tendermintNode := args[0]
	web3Provider := args[1]

	if !common.IsHexAddress(args[2]) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", args[1])
	}
	contractAddress := common.HexToAddress(args[2])

	if len(strings.Trim(args[3], "")) == 0 {
		return errors.Errorf("invalid [validator-moniker]: %s", args[2])
	}
	validatorMoniker := args[3]
	mnemonic := args[4]

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}
	sugaredLogger := logger.Sugar()

	// Initialize new Ethereum event listener
	inBuf := bufio.NewReader(cmd.InOrStdin())

	ethSub, err := relayer.NewEthereumSub(inBuf, tendermintNode, cdc, validatorMoniker, chainID, web3Provider,
		contractAddress, privateKey, mnemonic, nil, sugaredLogger)
	if err != nil {
		return err
	}

	ethSub.StartGasOracle()

	return nil
}
