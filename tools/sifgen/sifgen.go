package sifgen

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Sifchain/sifnode/tools/sifgen/network"
	"github.com/Sifchain/sifnode/tools/sifgen/node"
)

type Sifgen struct {
	chainID string
}

func NewSifgen(chainID string) Sifgen {
	return Sifgen{
		chainID: chainID,
	}
}

func (s Sifgen) NetworkCreate(count int, outputDir, startingIPAddress string, outputFile string) {
	net := network.NewNetwork(s.chainID)
	summary, err := net.Build(count, outputDir, startingIPAddress)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = ioutil.WriteFile(outputFile, []byte(*summary), 0600); err != nil {
		log.Fatal(err)
		return
	}
}

func (s Sifgen) NodeCreate(peerAddress, genesisURL *string) {
	witness := node.NewNode(s.chainID, peerAddress, genesisURL)
	summary, err := witness.Build()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(*summary)
}
