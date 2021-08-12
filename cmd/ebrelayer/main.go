package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

const (
	// EnvPrefix defines the environment prefix for the root cmd
	levelDbFile = "relayerdb"
)

func buildRootCmd() *cobra.Command {
	// see cmd/sifnoded/cmd/root.go:37 ; we need to do the
	// same thing in ebrelayer
	encodingConfig := sifapp.MakeTestEncodingConfig()
	authclient.Codec = encodingConfig.Marshaler
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(sifapp.DefaultNodeHome)

	// Read in the configuration file for the sdk
	// config := sdk.GetConfig()
	// config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	// config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	// config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	// config.Seal()

	rootCmd := &cobra.Command{
		Use:   "ebrelayer",
		Short: "Streams live events from Ethereum and Cosmos and relays event information to the opposite chain",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmd.Flags().Set(flags.FlagSkipConfirmation, "true"); err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd)
		},
	}

	log.SetFlags(log.Lshortfile)

	sifapp.SetConfig(true)

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentFlags().String(flags.FlagGas, "gas", fmt.Sprintf(
		"gas limit to set per-transaction; set to %q to calculate required gas automatically (default %d)",
		flags.GasFlagAuto, flags.DefaultGasLimit,
	))
	rootCmd.PersistentFlags().String(flags.FlagGasPrices, "", "Gas prices to determine the transaction fee (e.g. 10uatom)")
	rootCmd.PersistentFlags().Float64(flags.FlagGasAdjustment, flags.DefaultGasAdjustment, "gas adjustment")

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initRelayerCmd(),
		initWitnessCmd(),
		generateBindingsCmd(),
		replayEthereumCmd(),
		replayCosmosBurnLockCmd(),
		replayCosmosSignatureAggregationCmd(),
		listMissedCosmosEventCmd(),
	)
	return rootCmd
}

//	initRelayerCmd
func initRelayerCmd() *cobra.Command {
	//nolint:lll
	initRelayerCmd := &cobra.Command{
		Use:     "init-relayer [networkDescriptor] [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMnemonic]",
		Short:   "Validate credentials and initialize subscriptions to both chains",
		Args:    cobra.ExactArgs(5),
		Example: "ebrelayer init-relayer 1 tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 mnemonic --chain-id=peggy",
		RunE:    RunInitRelayerCmd,
	}
	flags.AddTxFlagsToCmd(initRelayerCmd)

	return initRelayerCmd
}

//	initWitnessCmd
func initWitnessCmd() *cobra.Command {
	//nolint:lll
	initWitnessCmd := &cobra.Command{
		Use:     "init-witness [networkDescriptor] [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMnemonic]",
		Short:   "Validate credentials and initialize subscriptions to both chains",
		Args:    cobra.ExactArgs(5),
		Example: "ebrelayer init-witness 1 tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 mnemonic --chain-id=peggy",
		RunE:    RunInitWitnessCmd,
	}
	flags.AddTxFlagsToCmd(initWitnessCmd)

	return initWitnessCmd
}

//	generateBindingsCmd : Generates ABIs and bindings for Bridge smart contracts which facilitate contract interaction
func generateBindingsCmd() *cobra.Command {
	generateBindingsCmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generates Bridge smart contracts ABIs and bindings",
		Args:    cobra.ExactArgs(0),
		Example: "generate",
		RunE:    RunGenerateBindingsCmd,
	}

	return generateBindingsCmd
}

// RunInitRelayerCmd executes initRelayerCmd
func RunInitRelayerCmd(cmd *cobra.Command, args []string) error {
	// First initialize the Cosmos features we need for the context
	cliContext, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}
	log.Printf("got result from GetClientQueryContext: %v", cliContext)

	// Load the validator's Ethereum private key from environment variables
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

	// Open the level db
	db, err := leveldb.OpenFile(levelDbFile, nil)
	if err != nil {
		log.Fatal("Error opening leveldb: ", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("db.Close filed: ", err.Error())
		}
	}()

	nodeURL, err := cmd.Flags().GetString(flags.FlagNode)
	if err != nil {
		return err
	}
	if nodeURL != "" {
		_, err := url.Parse(nodeURL)
		if nodeURL != "" && err != nil {
			return errors.Wrapf(err, "invalid RPC URL: %v", nodeURL)
		}
	}

	// Validate and parse arguments
	networkDescriptor, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Errorf("%s is invalid network id", args[0])
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

	if len(strings.Trim(args[4], "")) == 0 {
		return errors.Errorf("invalid [validator-moniker]: %s", args[4])
	}
	validatorMoniker := args[4]

	logConfig := zap.NewDevelopmentConfig()
	logConfig.Sampling = nil
	logger, err := logConfig.Build()

	if err != nil {
		log.Fatalln("failed to init zap logging")
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println("failed to sync zap logging")
		}
	}()

	sugaredLogger := logger.Sugar()
	zap.RedirectStdLog(sugaredLogger.Desugar())

	// Initialize new Ethereum event listener
	ethSub := relayer.NewEthereumSub(
		cliContext,
		nodeURL,
		validatorMoniker,
		web3Provider,
		contractAddress,
		nil,
		db,
		sugaredLogger,
	)

	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(oracletypes.NetworkDescriptor(networkDescriptor),
		privateKey,
		tendermintNode,
		web3Provider,
		contractAddress,
		db,
		cliContext,
		validatorMoniker,
		true,
		sugaredLogger)

	waitForAll := sync.WaitGroup{}
	waitForAll.Add(2)
	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())
	go ethSub.Start(txFactory, &waitForAll)
	go cosmosSub.Start(txFactory, &waitForAll)
	waitForAll.Wait()

	return nil
}

// RunInitWitnessCmd executes initWitnessCmd
func RunInitWitnessCmd(cmd *cobra.Command, args []string) error {
	// First initialize the Cosmos features we need for the context
	cliContext, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}
	log.Printf("got result from GetClientQueryContext: %v", cliContext)

	// Load the validator's Ethereum private key from environment variables
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

	// Open the level db
	db, err := leveldb.OpenFile(levelDbFile, nil)
	if err != nil {
		log.Fatal("Error opening leveldb: ", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("db.Close filed: ", err.Error())
		}
	}()

	nodeURL, err := cmd.Flags().GetString(flags.FlagNode)
	if err != nil {
		return err
	}
	if nodeURL != "" {
		_, err := url.Parse(nodeURL)
		if nodeURL != "" && err != nil {
			return errors.Wrapf(err, "invalid RPC URL: %v", nodeURL)
		}
	}

	// Validate and parse arguments
	networkDescriptor, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Errorf("%s is invalid network id", args[0])
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

	if len(strings.Trim(args[4], "")) == 0 {
		return errors.Errorf("invalid [validator-moniker]: %s", args[4])
	}
	validatorMoniker := args[4]

	logConfig := zap.NewDevelopmentConfig()
	logConfig.Sampling = nil
	logger, err := logConfig.Build()

	if err != nil {
		log.Fatalln("failed to init zap logging")
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println("failed to sync zap logging")
		}
	}()

	sugaredLogger := logger.Sugar()
	zap.RedirectStdLog(sugaredLogger.Desugar())

	// Initialize new Ethereum event listener
	ethSub := relayer.NewEthereumSub(
		cliContext,
		nodeURL,
		validatorMoniker,
		web3Provider,
		contractAddress,
		nil,
		db,
		sugaredLogger,
	)

	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(oracletypes.NetworkDescriptor(networkDescriptor),
		privateKey,
		tendermintNode,
		web3Provider,
		contractAddress,
		db,
		cliContext,
		validatorMoniker,
		false,
		sugaredLogger)

	waitForAll := sync.WaitGroup{}
	waitForAll.Add(2)
	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())
	go ethSub.Start(txFactory, &waitForAll)
	go cosmosSub.Start(txFactory, &waitForAll)
	waitForAll.Wait()

	return nil
}

// RunGenerateBindingsCmd : executes the generateBindingsCmd
func RunGenerateBindingsCmd(_ *cobra.Command, _ []string) error {
	contracts := contract.LoadBridgeContracts()

	// Compile contracts, generating contract bins and abis
	err := contract.CompileContracts(contracts)
	if err != nil {
		log.Println(err)
		return err
	}

	// Generate contract bindings from bins and abis
	return contract.GenerateBindings(contracts)
}

func replayEthereumCmd() *cobra.Command {
	//nolint:lll
	replayEthereumCmd := &cobra.Command{
		Use:     "replayEthereum [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMoniker] [fromBlock] [toBlock] [sifFromBlock] [sifEndBlock]",
		Short:   "replay missed ethereum events",
		Args:    cobra.ExactArgs(8),
		Example: "replayEthereum tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 validator 100 200 100 200 --chain-id=peggy",
		RunE:    RunReplayEthereumCmd,
	}

	flags.AddTxFlagsToCmd(replayEthereumCmd)

	return replayEthereumCmd
}

func replayCosmosBurnLockCmd() *cobra.Command {
	//nolint:lll
	replayCosmosBurnLockCmd := &cobra.Command{
		Use:     "replayCosmosBurnLock [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMoniker]",
		Short:   "replay missed cosmos events",
		Args:    cobra.ExactArgs(4),
		Example: "replayCosmos tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 validator",
		RunE:    RunReplayCosmosBurnLockCmd,
	}

	return replayCosmosBurnLockCmd
}

func replayCosmosSignatureAggregationCmd() *cobra.Command {
	//nolint:lll
	replayCosmosSignatureAggregationCmd := &cobra.Command{
		Use:     "replayCosmosSignatureAggregation [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMoniker]",
		Short:   "replay missed cosmos events",
		Args:    cobra.ExactArgs(4),
		Example: "replayCosmos tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 validator",
		RunE:    RunReplayCosmosSignatureAggregationCmd,
	}

	return replayCosmosSignatureAggregationCmd
}

func listMissedCosmosEventCmd() *cobra.Command {
	//nolint:lll
	listMissedCosmosEventCmd := &cobra.Command{
		Use:     "listMissedCosmosEventCmd [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [ebrelayerEthereumAddress]",
		Short:   "replay missed cosmos events",
		Args:    cobra.ExactArgs(4),
		Example: "listMissedCosmosEventCmd tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 0x627306090abaB3A6e1400e9345bC60c78a8BEf57",
		RunE:    RunListMissedCosmosEventCmd,
	}

	return listMissedCosmosEventCmd
}

func main() {
	if err := svrcmd.Execute(buildRootCmd(), sifapp.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
