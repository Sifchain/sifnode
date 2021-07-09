// TODO Fix me. This is broke with the new update and is having go mod issues
package main

func main() {}

//
// import (
// 	"flag"
// 	"fmt"
// 	"os"
//
// 	"github.com/tendermint/cosmos-rosetta-gateway/cosmos/launchpad"
// 	"github.com/tendermint/cosmos-rosetta-gateway/service"
// )
//
// var (
// 	flagAppRPC        = flag.String("app-rpc", "localhost:1317", "Application's RPC endpoint.")
// 	flagTendermintRPC = flag.String("tendermint-rpc", "localhost:26657", "Tendermint's RPC endpoint.")
// 	flagBlockchain    = flag.String("blockchain", "sifchain", "Application's name (e.g. Cosmos Hub)")
// 	flagNetworkDescriptor     = flag.String("network", "localnet", "Network's identifier (e.g. cosmos-hub-3, testnet-1, etc)")
// 	flagOfflineMode   = flag.Bool("offline", false, "Flag that forces the rosetta service to run in offline mode, some endpoints won't work.")
// 	flagAddrPrefix    = flag.String("prefix", "sif", "Bech32 prefix of address (e.g. cosmos, iaa, xrn:)")
// 	flagPort          = flag.Uint("port", 8080, "The port where the service is exposed.")
// )
//
// func main() {
// 	flag.Parse()
//
// 	h, err := service.New(
// 		service.Options{Port: uint32(*flagPort)},
// 		launchpad.NewLaunchpadNetwork(launchpad.Options{
// 			CosmosEndpoint:     *flagAppRPC,
// 			TendermintEndpoint: *flagTendermintRPC,
// 			Blockchain:         *flagBlockchain,
// 			Network:            *flagNetworkDescriptor,
// 			AddrPrefix:         *flagAddrPrefix,
// 			OfflineMode:        *flagOfflineMode,
// 		}),
// 	)
// 	if err != nil {
// 		fmt.Fprintln(flag.CommandLine.Output(), err)
// 		os.Exit(2)
// 	}
//
// 	fmt.Fprintf(flag.CommandLine.Output(), "Listening at http://localhost:%d\n", *flagPort)
// 	err = h.Start()
// 	if err != nil {
// 		fmt.Fprintln(flag.CommandLine.Output(), err)
// 		os.Exit(2)
// 	}
// }
