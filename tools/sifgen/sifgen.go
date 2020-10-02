package sifgen

import (
	"fmt"

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
	peerAddress *string
	genesisURL  *string
}

func NewSifgen(nodeType, network, chainID string, peerAddress, genesisURL *string) Sifgen {
	return Sifgen{
		nodeType:    nodeType,
		network:     network,
		chainID:     chainID,
		peerAddress: peerAddress,
		genesisURL:  genesisURL,
	}
}

func (s Sifgen) Run() {
	utils := NetworkUtils()
	node, err := NewNetworkNode(s, utils)
	if err != nil {
		panic(err)
	}

	network, err := NewNetwork(s, utils, *node)
	if err != nil {
		panic(err)
	}

	err = (*network).Setup()
	if err != nil {
		panic(err)
	}

	err = (*network).Genesis()
	if err != nil {
		panic(err)
	}

	s.summary(*node)
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
		Moniker: ` + node.Moniker() + `
		Address: ` + address + `
		Password: ` + node.KeyPassword() + `
	`))
}

func NetworkUtils() networks.NetworkUtils {
	return networks.NewUtils(app.DefaultNodeHome)
}

func NewNetworkNode(s Sifgen, utils networks.NetworkUtils) (*networks.NetworkNode, error) {
	var node networks.NetworkNode

	switch s.nodeType {
	case validator:
		node = networks.NewValidator(utils)
	case witness:
		node = networks.NewWitness(*s.peerAddress, *s.genesisURL, utils)
	default:
		return nil, notImplemented(s.nodeType)
	}

	return &node, nil
}

func NewNetwork(s Sifgen, utils networks.NetworkUtils, node networks.NetworkNode) (*networks.Network, error) {
	var network networks.Network

	switch s.network {
	case localnet:
		network = networks.NewLocalnet(app.DefaultNodeHome, app.DefaultCLIHome, s.chainID, node, utils)
	case testnet:
		return nil, notImplemented(s.network)
	case mainnet:
		return nil, notImplemented(s.network)
	default:
		return nil, notImplemented(s.network)
	}

	return &network, nil
}

func notImplemented(item string) error {
	return fmt.Errorf("%s not implemented", item)
}
