package bank

import (
	"github.com/Sifchain/sifnode/tools/sifgen/utils"
)

type Bank struct {
	chainID string
	CLI     utils.CLIUtils
}

func NewBank(chainID string) Bank {
	return Bank{
		chainID: chainID,
		CLI:     utils.NewCLI(chainID),
	}
}

func (b Bank) DefaultDeposit() []string {
	return []string{
		"1000000000000000stake",
		"1000000000000rowan",
	}
}
