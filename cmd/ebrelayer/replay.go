package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// RunReplayEthereumCmd executes replayEthereumCmd
func RunReplayEthereumCmd(cmd *cobra.Command, args []string) error {
	cliContext, err := client.GetClientTxContext(cmd)

	if err != nil {
		return err
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
	//mnemonic := args[4]

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

	ethSub := relayer.NewEthereumSub(cliContext, tendermintNode, validatorMoniker, web3Provider,
		contractAddress, nil, nil, sugaredLogger)

	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())
	ethSub.Replay(txFactory, fromBlock, toBlock, cosmosFromBlock, cosmosToBlock)

	return nil
}

// RunReplayCosmosCmd executes initRelayerCmd
func RunReplayCosmosCmd(cmd *cobra.Command, args []string) error {
	cliContext, err := client.GetClientTxContext(cmd)

	if err != nil {
		return err
	}

	// Validate and parse arguments
	networkDescriptor, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Errorf("%s is invalid network descriptor", args[0])
	}

	// Load the validator's Ethereum private key from environment variables
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

	if len(strings.Trim(args[1], "")) == 0 {
		return errors.Errorf("invalid [tendermint-node]: %s", args[1])
	}
	tendermintNode := args[1]

	if !relayer.IsWebsocketURL(args[2]) {
		return errors.Errorf("invalid [web3-provider]: %s", args[2])
	}
	web3Provider := args[2]

	if !common.IsHexAddress(args[3]) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", args[3])
	}
	contractAddress := common.HexToAddress(args[3])

	validatorMoniker := args[4]

	fromBlock, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [from-block]: %s", args[4])
	}

	toBlock, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [to-block]: %s", args[5])
	}

	ethFromBlock, err := strconv.ParseInt(args[7], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [eth-from-block]: %s", args[6])
	}

	ethToBlock, err := strconv.ParseInt(args[8], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [eth-to-block]: %s", args[7])
	}

	// check if the networkDescriptor is valid
	if !oracletypes.NetworkDescriptor(networkDescriptor).IsValid() {
		return errors.Errorf("network id: %d is invalid", networkDescriptor)
	}

	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}

	sugaredLogger := logger.Sugar()

	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(oracletypes.NetworkDescriptor(networkDescriptor),
		privateKey, tendermintNode, web3Provider, contractAddress, nil, cliContext,
		validatorMoniker, true, sugaredLogger)

	cosmosSub.Replay(txFactory, fromBlock, toBlock, ethFromBlock, ethToBlock)

	return nil
}

// RunListMissedCosmosEventCmd executes initRelayerCmd
func RunListMissedCosmosEventCmd(_ *cobra.Command, args []string) error {
	// Validate and parse arguments
	networkDescriptor, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Errorf("%s is invalid network descriptor", args[0])
	}

	// check if the networkDescriptor is valid
	if !oracletypes.NetworkDescriptor(networkDescriptor).IsValid() {
		return errors.Errorf("network id: %d is invalid", networkDescriptor)
	}

	if len(strings.Trim(args[1], "")) == 0 {
		return errors.Errorf("invalid [tendermint-node]: %s", args[1])
	}
	tendermintNode := args[1]

	if !relayer.IsWebsocketURL(args[2]) {
		return errors.Errorf("invalid [web3-provider]: %s", args[2])
	}
	web3Provider := args[2]

	if !common.IsHexAddress(args[3]) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", args[3])
	}
	contractAddress := common.HexToAddress(args[3])

	if !common.IsHexAddress(args[4]) {
		return errors.Errorf("invalid [relayer-ethereum-address]: %s", args[4])
	}
	relayerEthereumAddress := common.HexToAddress(args[4])

	days, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return errors.Errorf("invalid [days]: %s", args[5])
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to init zap logging")
	}
	sugaredLogger := logger.Sugar()

	// Initialize new Cosmos event listener
	listMissedCosmosEvent := relayer.NewListMissedCosmosEvent(oracletypes.NetworkDescriptor(networkDescriptor), tendermintNode, web3Provider, contractAddress, relayerEthereumAddress, days, sugaredLogger)

	listMissedCosmosEvent.ListMissedCosmosEvent()

	return nil
}
