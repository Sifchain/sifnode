package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/relayer"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
)

const (
	// FlagRPCURL defines the URL for the tendermint RPC connection
	FlagRPCURL = "rpc-url" // TODO: deprecate in favor of cosmos reading from flags.FlagNode
	// EnvPrefix defines the environment prefix for the root cmd
	levelDbFile = "relayerdb"
)

var rootCmd = &cobra.Command{
	Use:          "ebrelayer",
	Short:        "Streams live events from Ethereum and Cosmos and relays event information to the opposite chain",
	SilenceUsage: true,
}

func init() {

	encodingConfig := sifapp.MakeTestEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(sifapp.DefaultNodeHome)

	log.SetFlags(log.Lshortfile)

	sifapp.SetConfig(true)

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend,
		"Select keyring's backend (os|file|test)")
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	// Set FlagNode as seen in rpc.StatusCommand()
	rootCmd.PersistentFlags().StringP(flags.FlagNode, "n", "tcp://localhost:26657", "Tendermint node to connect to")
	rootCmd.PersistentFlags().String(flags.FlagGas, "gas", fmt.Sprintf(
		"gas limit to set per-transaction; set to %q to calculate required gas automatically (default %d)",
		flags.GasFlagAuto, flags.DefaultGasLimit,
	))
	rootCmd.PersistentFlags().String(flags.FlagGasPrices, "", "Gas prices to determine the transaction fee (e.g. 10uatom)")
	rootCmd.PersistentFlags().Float64(flags.FlagGasAdjustment, flags.DefaultGasAdjustment, "gas adjustment")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Set FlagNode from first arg, so that SetCmdClientContextHandler sets it in context / connection.
		if len(args) > 0 {
			rootCmd.PersistentFlags().Set(flags.FlagNode, args[0])
		}

		viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID))

		if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
			return err
		}

		return nil
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		initRelayerCmd(),
		generateBindingsCmd(),
		replayEthereumCmd(),
		replayCosmosCmd(),
		listMissedCosmosEventCmd(),
	)
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
	mnemonic := args[4]

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
	inBuf := bufio.NewReader(cmd.InOrStdin())
	ethSub, err := relayer.NewEthereumSub(inBuf, client.GetClientContextFromCmd(cmd), validatorMoniker, chainID, web3Provider,
		contractAddress, privateKey, mnemonic, db, sugaredLogger)
	if err != nil {
		return err
	}
	// Initialize new Cosmos event listener
	cosmosSub := relayer.NewCosmosSub(tendermintNode, web3Provider, contractAddress, privateKey, db, sugaredLogger)

	waitForAll := sync.WaitGroup{}
	waitForAll.Add(2)
	go ethSub.Start(&waitForAll)
	go cosmosSub.Start(&waitForAll)
	waitForAll.Wait()

	return nil
}

// RunGenerateBindingsCmd : executes the generateBindingsCmd
func RunGenerateBindingsCmd(cmd *cobra.Command, args []string) error {
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

func main() {
	if err := svrcmd.Execute(rootCmd, sifapp.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
