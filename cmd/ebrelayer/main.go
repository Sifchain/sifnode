package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	ebrelayertypes "github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	flag "github.com/spf13/pflag"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
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
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
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
	rootCmd.PersistentFlags().String(
		ebrelayertypes.FlagRelayerDbPath,
		"./relayerdb",
		"Path to the relayerdb directory",
	)
	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initRelayerCmd(),
		generateBindingsCmd(),
		replayEthereumCmd(),
		replayCosmosCmd(),
		listMissedCosmosEventCmd(),
		sendBridgeClaimCmd(),
	)

	return rootCmd
}

//	initRelayerCmd
func initRelayerCmd() *cobra.Command {
	//nolint:lll
	initRelayerCmd := &cobra.Command{
		Use:     "init [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMoniker] [validatorMnemonic]",
		Short:   "Validate credentials and initialize subscriptions to both chains",
		Args:    cobra.ExactArgs(5),
		Example: "ebrelayer init tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 validator mnemonic --chain-id=peggy",
		RunE:    RunInitRelayerCmd,
	}
	//flags.AddQueryFlagsToCmd(initRelayerCmd)
	flags.AddTxFlagsToCmd(initRelayerCmd)

	return initRelayerCmd
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

	levelDbFile, err := cmd.Flags().GetString(ebrelayertypes.FlagRelayerDbPath)
	if err != nil {
		return err
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

	if len(strings.Trim(args[3], "")) == 0 {
		return errors.Errorf("invalid [validator-moniker]: %s", args[3])
	}
	validatorMoniker := args[3]

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

	symbolTranslator, err := buildSymbolTranslator(cmd.Flags())
	if err != nil {
		return err
	}

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

	key, err := txs.LoadPrivateKey()
	if err != nil {
		log.Fatalf("failed to load ETHEREUM_PRIVATE_KEY")
	}

	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(tendermintNode, web3Provider, contractAddress, key, db, sugaredLogger)

	waitForAll := sync.WaitGroup{}
	waitForAll.Add(2)
	txFactory := tx.NewFactoryCLI(cliContext, cmd.Flags())
	go ethSub.Start(txFactory, &waitForAll, symbolTranslator)
	go cosmosSub.Start(&waitForAll, symbolTranslator)
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
		Use:     "replayEthereum [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMoniker] [validatorMnemonic] [fromBlock] [toBlock] [sifFromBlock] [sifEndBlock]",
		Short:   "replay missed ethereum events",
		Args:    cobra.ExactArgs(9),
		Example: "replayEthereum tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 validator mnemonic 100 200 100 200 --chain-id=peggy",
		RunE:    RunReplayEthereumCmd,
	}

	flags.AddTxFlagsToCmd(replayEthereumCmd)

	return replayEthereumCmd
}

func replayCosmosCmd() *cobra.Command {
	//nolint:lll
	replayCosmosCmd := &cobra.Command{
		Use:     "replayCosmos [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [fromBlock] [toBlock] [ethFromBlock] [ethToBlock]",
		Short:   "replay missed cosmos events",
		Args:    cobra.ExactArgs(7),
		Example: "replayCosmos tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 100 200 100 200",
		RunE:    RunReplayCosmosCmd,
	}

	flags.AddTxFlagsToCmd(replayCosmosCmd)

	return replayCosmosCmd
}

func listMissedCosmosEventCmd() *cobra.Command {
	//nolint:lll
	listMissedCosmosEventCmd := &cobra.Command{
		Use:     "listMissedCosmosEventCmd [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [ebrelayerEthereumAddress] [days]",
		Short:   "replay missed cosmos events",
		Args:    cobra.ExactArgs(5),
		Example: "listMissedCosmosEventCmd tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 1",
		RunE:    RunListMissedCosmosEventCmd,
	}

	return listMissedCosmosEventCmd
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

func sendBridgeClaimCmd() *cobra.Command {
	//nolint:lll
	sendBridgeClaimCmd := &cobra.Command{
		// add amount as argument then we can use small amount for testing.
		Use:     "sendBridgeClaimCmd [validatorMoniker] [nonce] [amount] ",
		Short:   "send bridge claim to sifchain",
		Args:    cobra.ExactArgs(3),
		Example: "replayEthereum lisa 0 100 --chain-id=peggy",
		RunE:    RunSendBridgeClaimCmd,
	}

	flags.AddTxFlagsToCmd(sendBridgeClaimCmd)

	return sendBridgeClaimCmd
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
