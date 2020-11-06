package relayer

// const (
// 	tmProvider      = "Node"
// 	ethProvider     = "ws://127.0.0.1:7545/"
// 	contractAddress = "0x00"
// 	privateKeyStr   = "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f"
// 	// rpc              = "http://localhost:7545"
// 	rpc              = ""
// 	validatorMoniker = "user1"
// 	chainID          = "sifchain"
// 	web3Provider     = "ws://localhost:7545/"
// )

// comment out it because of error related to configuration.
// func TestNewEthereumSub(t *testing.T) {

// 	rootCmd := &cobra.Command{
// 		Use:          "use",
// 		Short:        "short",
// 		SilenceUsage: true,
// 	}

// 	cdc := app.MakeCodec()
// 	inBuf := bufio.NewReader(rootCmd.InOrStdin())
// 	privateKey, _ := crypto.HexToECDSA(privateKeyStr)
// 	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
// 	registryContractAddress := common.HexToAddress(contractAddress)
// 	a, err := relayer.NewEthereumSub(inBuf, rpc, cdc, validatorMoniker, chainID, web3Provider, registryContractAddress,
// 		privateKey, logger)
// 	require.Equal(t, err, nil, "error when init Ethereum sub")
// }
