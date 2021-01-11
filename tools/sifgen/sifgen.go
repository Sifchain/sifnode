package sifgen

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/MakeNowJust/heredoc"
	"github.com/Sifchain/sifnode/tools/sifgen/key"
	"github.com/Sifchain/sifnode/tools/sifgen/network"
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

func (s Sifgen) NetworkCreate(count int, outputDir, startingIPAddress string, outputFile string) {
	net := network.NewNetwork(*s.chainID)
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

func (s Sifgen) NetworkReset(networkDir string) {
	if err := network.Reset(*s.chainID, networkDir); err != nil {
		log.Fatal(err)
	}
}

func (s Sifgen) NodeCreate(moniker, mnemonic string, adminCLPAddresses []string, adminOracleAddress, ipAddr string, peerAddress, genesisURL *string, printDetails, withCosmovisor *bool) {
	validator := node.NewNode(*s.chainID,
		moniker,
		mnemonic,
		adminCLPAddresses,
		adminOracleAddress,
		ipAddr,
		peerAddress,
		genesisURL,
		withCosmovisor)

	summary, err := validator.Build()
	if err != nil {
		log.Fatal(err)
		return
	}

	if *printDetails {
		fmt.Println(*summary)
	}
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
