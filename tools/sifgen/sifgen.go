package sifgen

import (
	"fmt"
	"log"

	"github.com/MakeNowJust/heredoc"
	"github.com/Sifchain/sifnode/tools/sifgen/key"
	"github.com/Sifchain/sifnode/tools/sifgen/node"
)

type Sifgen struct {
	chainID *string
}

func NewSifgen(chainID *string) Sifgen {
	return Sifgen{
		chainID: chainID,
	}
}

func (s Sifgen) NewNode() *node.Node {
	return &node.Node{}
}

func (s Sifgen) NodeReset(nodeHomeDir *string) {
	if err := node.Reset(*s.chainID, nodeHomeDir); err != nil {
		log.Fatal(err)
	}
}

func (s Sifgen) KeyGenerateMnemonic(name, password *string) {
	key := key.NewKey(name, password)
	key.GenerateMnemonic()
	fmt.Println(key.Mnemonic)
}

func (s Sifgen) KeyRecoverFromMnemonic(mnemonic string) {
	key := key.NewKey(nil, nil)
	if err := key.RecoverFromMnemonic(mnemonic); err != nil {
		log.Fatal(err)
	}

	fmt.Println(heredoc.Doc(`
		Address: ` + key.Address + `
		Validator Address: ` + key.ValidatorAddress + `
		Consensus Address: ` + key.ConsensusAddress + `
	`))
}
