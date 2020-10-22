package sifgen

import (
	"github.com/Sifchain/sifnode/tools/sifgen/network"
)

type Sifgen struct {
	chainID string
}

func NewSifgen(chainID string) Sifgen {
	return Sifgen{
		chainID: chainID,
	}
}

func (s Sifgen) NetworkCreate(count int, outputDir, startingIPAddress string) {
	net := network.NewNetwork(s.chainID)
	_ = net.Build(count, outputDir, startingIPAddress)
}
