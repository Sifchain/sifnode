package main

import (
	"fmt"
	"log"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/Sifchain/sifnode/x/instrumentation"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	ebrelayertypes "github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	flag "github.com/spf13/pflag"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	networkDescriptorFlag             = "network-descriptor"
	tendermintNodeFlag                = "tendermint-node"
	web3ProviderFlag                  = "web3-provider"
	bridgeRegistryContractAddressFlag = "bridge-registry-contract-address"
	validatorMnemonicFlag             = "validator-mnemonic"
	maxFeePerGasFlag                  = "maxFeePerGasFlag"
	maxPriorityFeePerGasFlag          = "maxPriorityFeePerGasFlag"
	ethereumChainIdFlag               = "ethereum-chain-id"
	sifnodeGrpcEntryPointFlag         = "sifnode-grpc"
	defaultSifnodeGrpc                = "0.0.0.0:9090"
)

func buildRootCmd() *cobra.Command {
	// see cmd/sifnoded/cmd/root.go:37 ; we need to do the
	// same thing in ebrelayer
	encodingConfig := sifapp.MakeTestEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
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

			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
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
	rootCmd.PersistentFlags().String(
		ebrelayertypes.FlagSymbolTranslatorFile,
		"",
		"Path to a json file containing an array of sifchain denom => Ethereum symbol pairs",
	)
	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initRelayerCmd(),
		initWitnessCmd(),
	)
	return rootCmd
}

//	initRelayerCmd
func initRelayerCmd() *cobra.Command {
	//nolint:lll
	initRelayerCmd := &cobra.Command{
		Use:     "init-relayer --network-descriptor 1 --tendermint-node tcp://localhost:26657 --web3-provider ws://localhost:7545/ --bridge-registry-contract-address 0x --validator-mnemonic mnemonic --chain-id=peggy --sifnode-grpc 0.0.0.0:9090",
		Short:   "Validate credentials and initialize subscriptions to both chains",
		Args:    cobra.ExactArgs(0),
		Example: "ebrelayer init-relayer --network-descriptor 1 --tendermint-node tcp://localhost:26657 --web3-provider ws://localhost:7545/ --bridge-registry-contract-address 0x --validator-mnemonic mnemonic  --chain-id=peggy",
		RunE:    RunInitRelayerCmd,
	}
	flags.AddTxFlagsToCmd(initRelayerCmd)
	AddRelayerFlagsToCmd(initRelayerCmd)

	return initRelayerCmd
}

//	initWitnessCmd
func initWitnessCmd() *cobra.Command {
	//nolint:lll
	initWitnessCmd := &cobra.Command{
		Use:     "init-witness --network-descriptor 1 --tendermint-node tcp://localhost:26657 --web3-provider ws://localhost:7545/ --bridge-registry-contract-address 0x --validator-mnemonic mnemonic ",
		Short:   "Validate credentials and initialize subscriptions to both chains",
		Args:    cobra.ExactArgs(0),
		Example: "ebrelayer init-witness --network-descriptor 1 --tendermint-node tcp://localhost:26657 --web3-provider ws://localhost:7545/ --bridge-registry-contract-address 0x --validator-mnemonic mnemonic  --chain-id=peggy --sifnode-grpc 0.0.0.0:9090",
		RunE:    RunInitWitnessCmd,
	}
	flags.AddTxFlagsToCmd(initWitnessCmd)
	AddRelayerFlagsToCmd(initWitnessCmd)

	return initWitnessCmd
}

func parseGasArguments(cmd *cobra.Command, networkDescriptor *big.Int) (maxFeePerGas, maxPriorityFeePerGas, ethereumChainId *big.Int) {
	maxFeePerGasString, err := cmd.Flags().GetString(maxFeePerGasFlag)
	if err != nil {
		log.Fatalln("failed to parse maxFeePerGasFlag")
	}

	maxPriorityFeePerGasString, err := cmd.Flags().GetString(maxPriorityFeePerGasFlag)
	if err != nil {
		log.Fatalln("failed to parse maxPriorityFeePerGasFlag")
	}

	ethereumChainIdString, err := cmd.Flags().GetString(ethereumChainIdFlag)
	if err != nil {
		log.Fatalln("failed to parse ethereumChainIdFlag", ethereumChainIdFlag)
	}

	maxFeePerGas, _ = (&big.Int{}).SetString(maxFeePerGasString, 10)
	maxPriorityFeePerGas, _ = (&big.Int{}).SetString(maxPriorityFeePerGasString, 10)
	ethereumChainId, _ = (&big.Int{}).SetString(ethereumChainIdString, 10)
	if ethereumChainId == nil {
		ethereumChainId = big.NewInt(0).Set(networkDescriptor)
	}

	return maxFeePerGas, maxPriorityFeePerGas, ethereumChainId
}

// RunInitRelayerCmd executes initRelayerCmd
func RunInitRelayerCmd(cmd *cobra.Command, args []string) error {
	// First initialize the Cosmos features we need for the context
	cliContext, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	// Load the validator's Ethereum private key from environment variables
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

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
	networkDescriptorString, err := cmd.Flags().GetString(networkDescriptorFlag)
	if err != nil {
		return errors.Errorf("network descriptor is invalid: %s", err.Error())
	}

	if len(networkDescriptorString) == 0 {
		return errors.New("network descriptor is empty")
	}

	networkDescriptor, err := strconv.Atoi(networkDescriptorString)
	if err != nil {
		return errors.Errorf("network id: %s is invalid", networkDescriptorString)
	}

	// check if the networkDescriptor is valid
	if !oracletypes.NetworkDescriptor(networkDescriptor).IsValid() {
		return errors.Errorf("network id: %d is invalid", networkDescriptor)
	}

	tendermintNode, err := cmd.Flags().GetString(tendermintNodeFlag)
	if err != nil {
		return errors.Errorf("tendermint node is invalid: %s", err.Error())
	}

	if len(tendermintNode) == 0 {
		return errors.New("tendermint node is empty")
	}

	web3Provider, err := cmd.Flags().GetString(web3ProviderFlag)
	if err != nil {
		return errors.Errorf("web3 provider is invalid: %s", err.Error())
	}

	if len(web3Provider) == 0 {
		return errors.New("web3 provider is empty")
	}

	contractAddressString, err := cmd.Flags().GetString(bridgeRegistryContractAddressFlag)
	if err != nil {
		return errors.Errorf("contract address is invalid: %s", err.Error())
	}

	if len(contractAddressString) == 0 {
		return errors.New("contract address is empty")
	}

	if !common.IsHexAddress(contractAddressString) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", contractAddressString)
	}
	contractAddress := common.HexToAddress(contractAddressString)

	validatorMoniker, err := cmd.Flags().GetString(validatorMnemonicFlag)
	if err != nil {
		return errors.Errorf("validator moniker is invalid: %s", err.Error())
	}

	if len(validatorMoniker) == 0 {
		return errors.New("validator moniker is empty")
	}

	sifnodeGrpc, err := cmd.Flags().GetString(sifnodeGrpcEntryPointFlag)
	if err != nil {
		// set default value
		sifnodeGrpc = defaultSifnodeGrpc
	}

	if len(sifnodeGrpc) == 0 {
		sifnodeGrpc = defaultSifnodeGrpc
	}

	logConfig := zap.NewDevelopmentConfig()
	logConfig.Sampling = nil
	logConfig.Encoding = "json"
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

	symbolTranslator, err := buildSymbolTranslator(cmd.Flags())
	if err != nil {
		return err
	}

	instrumentation.PeggyCheckpointZap(
		sugaredLogger,
		instrumentation.Startup,
		"starting", true,
	)

	// Initialize new Ethereum event listener
	ethSub := relayer.NewEthereumSub(
		cliContext,
		nodeURL,
		validatorMoniker,
		oracletypes.NetworkDescriptor(networkDescriptor),
		web3Provider,
		contractAddress,
		sugaredLogger,
		sifnodeGrpc,
	)

	bigNetworkDescriptor := big.NewInt(int64(networkDescriptor))
	maxFeePerGas, maxPriorityFeePerGas, ethereumChainId := parseGasArguments(cmd, bigNetworkDescriptor)

	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(oracletypes.NetworkDescriptor(networkDescriptor),
		privateKey,
		tendermintNode,
		web3Provider,
		contractAddress,
		cliContext,
		validatorMoniker,
		sugaredLogger,
		maxFeePerGas,
		maxPriorityFeePerGas,
		ethereumChainId,
		sifnodeGrpc,
	)

	waitForAll := sync.WaitGroup{}
	waitForAll.Add(2)
	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())
	go ethSub.Start(txFactory, &waitForAll, symbolTranslator)
	go cosmosSub.StartProphecyHandler(txFactory, &waitForAll, symbolTranslator)
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

	// Load the validator's Ethereum private key from environment variables
	privateKey, err := txs.LoadPrivateKey()
	if err != nil {
		return errors.Errorf("invalid [ETHEREUM_PRIVATE_KEY] environment variable")
	}

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

	networkDescriptorString, err := cmd.Flags().GetString(networkDescriptorFlag)
	if err != nil {
		return errors.Errorf("network descriptor is invalid: %s", err.Error())
	}

	if len(networkDescriptorString) == 0 {
		return errors.New("network descriptor is empty")
	}

	networkDescriptor, err := strconv.Atoi(networkDescriptorString)
	if err != nil {
		return errors.Errorf("network id: %s is invalid", networkDescriptorString)
	}

	// check if the networkDescriptor is valid
	if !oracletypes.NetworkDescriptor(networkDescriptor).IsValid() {
		return errors.Errorf("network id: %d is invalid", networkDescriptor)
	}

	tendermintNode, err := cmd.Flags().GetString(tendermintNodeFlag)
	if err != nil {
		return errors.Errorf("tendermint node is invalid: %s", err.Error())
	}

	if len(tendermintNode) == 0 {
		return errors.New("tendermint node is empty")
	}

	web3Provider, err := cmd.Flags().GetString(web3ProviderFlag)
	if err != nil {
		return errors.Errorf("web3 provider is invalid: %s", err.Error())
	}

	if len(web3Provider) == 0 {
		return errors.New("web3 provider is empty")
	}

	contractAddressString, err := cmd.Flags().GetString(bridgeRegistryContractAddressFlag)
	if err != nil {
		return errors.Errorf("contract address is invalid: %s", err.Error())
	}

	if len(contractAddressString) == 0 {
		return errors.New("contract address is empty")
	}

	if !common.IsHexAddress(contractAddressString) {
		return errors.Errorf("invalid [bridge-registry-contract-address]: %s", contractAddressString)
	}
	contractAddress := common.HexToAddress(contractAddressString)

	validatorMoniker, err := cmd.Flags().GetString(validatorMnemonicFlag)
	if err != nil {
		return errors.Errorf("validator moniker is invalid: %s", err.Error())
	}

	if len(validatorMoniker) == 0 {
		return errors.New("validator moniker is empty")
	}

	sifnodeGrpc, err := cmd.Flags().GetString(sifnodeGrpcEntryPointFlag)
	if err != nil {
		// set default value
		sifnodeGrpc = defaultSifnodeGrpc
	}

	if len(sifnodeGrpc) == 0 {
		sifnodeGrpc = defaultSifnodeGrpc
	}

	logConfig := zap.NewDevelopmentConfig()
	logConfig.Encoding = "json"
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

	symbolTranslator, err := buildSymbolTranslator(cmd.Flags())
	if err != nil {
		return err
	}

	sugaredLogger := logger.Sugar()
	zap.RedirectStdLog(sugaredLogger.Desugar())

	// Initialize new Ethereum event listener
	ethSub := relayer.NewEthereumSub(
		cliContext,
		nodeURL,
		validatorMoniker,
		oracletypes.NetworkDescriptor(networkDescriptor),
		web3Provider,
		contractAddress,
		sugaredLogger,
		sifnodeGrpc,
	)

	bigNetworkDescriptor := big.NewInt(int64(networkDescriptor))
	maxFeePerGas, maxPriorityFeePerGas, ethereumChainId := parseGasArguments(cmd, bigNetworkDescriptor)

	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(oracletypes.NetworkDescriptor(networkDescriptor),
		privateKey,
		tendermintNode,
		web3Provider,
		contractAddress,
		cliContext,
		validatorMoniker,
		sugaredLogger,
		maxFeePerGas,
		maxPriorityFeePerGas,
		ethereumChainId,
		sifnodeGrpc,
	)

	waitForAll := sync.WaitGroup{}
	waitForAll.Add(2)
	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())
	go ethSub.Start(txFactory, &waitForAll, symbolTranslator)
	go cosmosSub.Start(txFactory, &waitForAll, symbolTranslator)
	waitForAll.Wait()

	return nil
}

// AddRelayerFlagsToCmd adds all common flags to relayer commands.
func AddRelayerFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(
		networkDescriptorFlag,
		"",
		"The network descriptor for the chain",
	)
	cmd.Flags().String(
		tendermintNodeFlag,
		"",
		"Sifchain node address",
	)
	cmd.Flags().String(
		web3ProviderFlag,
		"",
		"Ethereum web3 service address",
	)
	cmd.Flags().String(
		bridgeRegistryContractAddressFlag,
		"",
		"Ethereum bridge registry contract address",
	)
	cmd.Flags().String(
		validatorMnemonicFlag,
		"",
		"Validator mnemonic",
	)
	cmd.Flags().String(
		maxFeePerGasFlag,
		"30000000000000",
		"maxFeePerGasFlag for ethereum in wei",
	)
	cmd.Flags().String(
		maxPriorityFeePerGasFlag,
		"1000000000",
		"maxFeePerGasFlag for ethereum in wei",
	)
	cmd.Flags().String(
		ethereumChainIdFlag,
		"",
		"the Ethereum chain id (defaults to --network-descriptor)",
	)
	cmd.Flags().String(
		sifnodeGrpcEntryPointFlag,
		"",
		"The sifnode grpc url",
	)
}

func buildSymbolTranslator(flags *flag.FlagSet) (*symbol_translator.SymbolTranslator, error) {
	filename, err := flags.GetString(ebrelayertypes.FlagSymbolTranslatorFile)
	// If FlagSymbolTranslatorFile isn't specified, just use an empty SymbolTranslator
	if err != nil || filename == "" {
		return symbol_translator.NewSymbolTranslator(), nil
	}

	symbolTranslator, err := symbol_translator.NewSymbolTranslatorFromJSONFile(filename)
	if err != nil {
		return nil, err
	}

	return symbolTranslator, nil
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
