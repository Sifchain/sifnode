package main

import (
	"flag"
	"os"

	"github.com/Sifchain/sifnode/tools/sifgen"
)

func main() {
	nodeType := flag.String("t", os.Getenv("NODE_TYPE"), "The node type [validator|witness].")
	network := flag.String("n", os.Getenv("NETWORK"), "The network [localnet|testnet|mainnet].")
	chainID := flag.String("c", os.Getenv("CHAIN_ID"), "The ID of the chain.")
	peerAddress := flag.String("p", "", "The address of the peer to sync with (<hash>@<ip>:<port>).")
	genesisURL := flag.String("u", "", "The URL to download the Genesis file from.")
	flag.Parse()

	sif := sifgen.NewSifgen(*nodeType, *network, *chainID, *peerAddress, *genesisURL)
	sif.Run()
}
