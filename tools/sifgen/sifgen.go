package sifgen

import (
	"github.com/Sifchain/sifnode/tools/sifgen/network"
	"io/ioutil"
	"log"
)

type Sifgen struct {
	chainID string
}

func NewSifgen(chainID string) Sifgen {
	return Sifgen{
		chainID: chainID,
	}
}

func (s Sifgen) NetworkCreate(count int, outputDir, startingIPAddress string, outputFile *string) {
	net := network.NewNetwork(s.chainID)
	summary, err := net.Build(count, outputDir, startingIPAddress)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = ioutil.WriteFile(*outputFile, []byte(*summary), 0600); err != nil {
		log.Fatal(err)
		return
	}
}
