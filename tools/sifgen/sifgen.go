package sifgen

import (
	"fmt"
	"log"

	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/tools/sifgen/networks"

	"github.com/MakeNowJust/heredoc/v2"
)

const (
	validator = "validator"
	witness   = "witness"

	localnet = "localnet"
	testnet  = "testnet"
	mainnet  = "mainnet"
)

type Sifgen struct {
	nodeType    string
	network     string
	chainID     string
	peerAddress string
	genesisURL  string
}

func NewSifgen(nodeType, network, chainID, peerAddress, genesisURL string) Sifgen {
	return Sifgen{
		nodeType:    nodeType,
		network:     network,
		chainID:     chainID,
		peerAddress: peerAddress,
		genesisURL:  genesisURL,
	}
}

func (s Sifgen) Run() {
	node := *s.newNode(s.nodeType)
	network := s.newNetwork(s.network, s.chainID, node)

	err := (*network).Setup()
	if err != nil {
		panic(err)
	}

	err = (*network).Genesis()
	if err != nil {
		panic(err)
	}

	s.summary(node)
}

func (s Sifgen) newNode(nodeType string) *networks.NetworkNode {
	var node networks.NetworkNode

	switch nodeType {
	case validator:
		node = networks.NewValidator(s.networkUtils())
	case witness:
		node = networks.NewWitness(s.peerAddress, s.genesisURL, s.networkUtils())
	default:
		s.notImplemented(nodeType)
	}

	return &node
}

func (s Sifgen) newNetwork(networkType, chainID string, node networks.NetworkNode) *networks.Network {
	var network networks.Network

	switch networkType {
	case localnet:
		network = networks.NewLocalnet(app.DefaultNodeHome, app.DefaultCLIHome, chainID, node, s.networkUtils())
	case testnet:
		s.notImplemented(networkType)
	case mainnet:
		s.notImplemented(networkType)
	default:
		s.notImplemented(networkType)
	}

	return &network
}

func (s Sifgen) networkUtils() networks.NetworkUtils {
	return networks.NewUtils(app.DefaultNodeHome)
}

func (s Sifgen) notImplemented(item string) {
	log.Fatal(fmt.Sprintf("%s not implemented", item))
}

func (s Sifgen) summary(node networks.NetworkNode) {
	var address string

	_, isValidator := node.(*networks.Validator)
	if isValidator {
		address = fmt.Sprintf("%s (%s)", *node.Address(nil), node.PeerAddress())
	} else {
		address = fmt.Sprintf("%s", *node.Address(nil))
	}

	fmt.Println(heredoc.Doc(`
		Node Details
		============
		Name: ` + node.Name() + `
		Address: ` + address + `
		Password: ` + node.KeyPassword() + `
	`))
}
