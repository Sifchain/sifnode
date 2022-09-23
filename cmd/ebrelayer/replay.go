package main

import (
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"log"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
)

// RunReplayEthereumCmd executes replayEthereumCmd
func RunReplayEthereumCmd(cmd *cobra.Command, args []string) error {
	cliContext, err := client.GetClientTxContext(cmd)

	if err != nil {
		return err
	}

	if cliContext.From == "" {
		log.Println("Received empty clientContext.From, needed for validating cosmos transaction. Check if --from flag is set")
		return errors.New("Missing from flag ")
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

	symbolTranslator, err := buildSymbolTranslator(cmd.Flags())
	if err != nil {
		return err
	}

	ethSub := relayer.NewEthereumSub(cliContext, tendermintNode, validatorMoniker, web3Provider,
		contractAddress, nil, nil, sugaredLogger)

	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())
	ethSub.Replay(txFactory, fromBlock, toBlock, cosmosFromBlock, cosmosToBlock, symbolTranslator)

	return nil
}

// RunReplayCosmosCmd executes initRelayerCmd
func RunReplayCosmosCmd(cmd *cobra.Command, args []string) error {
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

	symbolTranslator, err := buildSymbolTranslator(cmd.Flags())
	if err != nil {
		return err
	}

	key, err := txs.LoadPrivateKey()
	if err != nil {
		log.Fatalf("failed to load ETHEREUM_PRIVATE_KEY")
	}

	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(tendermintNode, web3Provider, contractAddress, key, nil, sugaredLogger)

	cosmosSub.Replay(symbolTranslator, fromBlock, toBlock, ethFromBlock, ethToBlock)

	return nil
}

// RunListMissedCosmosEventCmd executes initRelayerCmd
func RunListMissedCosmosEventCmd(cmd *cobra.Command, args []string) error {
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

	symbolTranslator, err := buildSymbolTranslator(cmd.Flags())
	if err != nil {
		return err
	}

	// Initialize new Cosmos event listener
	listMissedCosmosEvent := relayer.NewListMissedCosmosEvent(tendermintNode, web3Provider, contractAddress, relayerEthereumAddress, days, sugaredLogger)

	listMissedCosmosEvent.ListMissedCosmosEvent(symbolTranslator)

	return nil
}
