package faucet

import (
	"github.com/Sifchain/sifnode/tools/sifgen/utils"
)

type Reserve interface {
	DefaultDeposit() []string
	Transfer(string, string, string, string) error
}

type Faucet struct {
	chainID string
	CLI     utils.CLIUtils
}

func NewFaucet(chainID string) Faucet {
	return Faucet{
		chainID: chainID,
		CLI:     utils.NewCLI(chainID),
	}
}

func (f Faucet) DefaultDeposit() []string {
	return []string{
		"1000000000000000stake",
		"1000000000000rowan",
	}
}
