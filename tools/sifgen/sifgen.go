package sifgen

import (
	"fmt"
	"log"
	"time"

	"github.com/Sifchain/sifnode/tools/sifgen/faucet"
	"github.com/Sifchain/sifnode/tools/sifgen/node"

	"github.com/yelinaung/go-haikunator"
	"gopkg.in/yaml.v3"
)

type Output struct {
	ChainID                   string `yaml:"chain_id"`
	Moniker                   string `yaml:"moniker"`
	KeyAddress                string `yaml:"key_address"`
	KeyPassword               string `yaml:"key_password"`
	PeerAddress               string `yaml:"peer_address"`
	ValidatorPublicKeyAddress string `yaml:"validator_public_key_address"`
}

type Sifgen struct {
	chainID string
}

func NewSifgen(chainID string) Sifgen {
	return Sifgen{
		chainID: chainID,
	}
}

func (s Sifgen) NodeCreate(seedAddress, genesisURL *string) {
	moniker := haikunator.New(time.Now().UTC().UnixNano()).Haikunate()
	nd := node.NewNode(s.chainID, &moniker, seedAddress, genesisURL)

	if err := nd.Setup(); err != nil {
		log.Fatal(err)
	}

	if err := nd.Genesis(faucet.NewFaucet(s.chainID).DefaultDeposit()); err != nil {
		log.Fatal(err)
	}

	s.summary(nd)
}

func (s Sifgen) NodePromote(moniker, validatorPublicKey, keyPassword, bondAmount string) {
	nd := node.NewNode(s.chainID, &moniker, nil, nil)
	if err := nd.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := nd.Promote(validatorPublicKey, keyPassword, bondAmount); err != nil {
		log.Fatal(err)
	}
}

func (s Sifgen) NodePeers(moniker string, peerList []string) {
	nd := node.NewNode(s.chainID, &moniker, nil, nil)
	if err := nd.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := nd.UpdatePeerList(peerList); err != nil {
		log.Fatal(err)
	}
}

func (s Sifgen) Transfer(fromKeyPassword, fromKeyAddress, toKeyAddress, amount string) {
	if err := faucet.NewFaucet(s.chainID).Transfer(fromKeyPassword, fromKeyAddress, toKeyAddress, amount); err != nil {
		log.Fatal(err)
	}
}

func (s Sifgen) summary(node *node.Node) {
	output := Output{
		ChainID:                   node.ChainID(),
		Moniker:                   node.Moniker(),
		KeyAddress:                *node.NodeKeyAddress(nil),
		KeyPassword:               node.NodeKeyPassword(),
		PeerAddress:               node.NodePeerAddress(),
		ValidatorPublicKeyAddress: node.NodeValidatorPublicKeyAddress(),
	}

	yml, _ := yaml.Marshal(output)
	fmt.Println(string(yml))
}
