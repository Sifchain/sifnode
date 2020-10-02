package networks

import (
	"io/ioutil"
	"strings"

	"github.com/Sifchain/sifnode/tools/sifgen/networks/types"

	"gopkg.in/yaml.v3"
)

type Localnet struct {
	defaultNodeHome string
	defaultCLIHome  string
	chainID         string
	utils           NetworkUtils
	node            *NetworkNode
}

func NewLocalnet(defaultNodeHome, defaultCLIHome, chainID string, node NetworkNode, utils NetworkUtils) Localnet {
	return Localnet{
		defaultNodeHome: defaultNodeHome,
		defaultCLIHome:  defaultCLIHome,
		chainID:         chainID,
		utils:           utils,
		node:            &node,
	}
}

// Setup the network; clear any existing config, initialize the blockchain, set the config and perform genesis.
func (l Localnet) Setup() error {
	err := l.utils.Reset([]string{l.defaultNodeHome, l.defaultCLIHome})
	if err != nil {
		return err
	}

	_, err = l.utils.InitChain(l.chainID, (*l.node).Moniker())
	if err != nil {
		return err
	}

	_, err = l.utils.SetKeyRingStorage()
	if err != nil {
		return err
	}

	_, err = l.utils.SetConfigChainID(l.chainID)
	if err != nil {
		return err
	}

	_, err = l.utils.SetConfigIndent(true)
	if err != nil {
		return err
	}

	_, err = l.utils.SetConfigTrustNode(true)
	if err != nil {
		return err
	}

	err = l.generateNodeKey()
	if err != nil {
		return err
	}

	return nil
}

func (l Localnet) Genesis() error {
	address := (*l.node).Address(nil)

	_, isValidator := (*l.node).(*Validator)
	if isValidator {
		if err := l.validatorGenesis(*address); err != nil {
			return err
		}
	} else {
		if err := l.witnessGenesis(*address); err != nil {
			return err
		}
	}

	return nil
}

// Generate a new key for a node.
func (l Localnet) generateNodeKey() error {
	output, err := l.utils.AddKey((*l.node).Moniker(), (*l.node).KeyPassword())
	if err != nil {
		return err
	}

	yml, err := ioutil.ReadAll(strings.NewReader(*output))
	if err != nil {
		return err
	}

	var keys types.Keys

	err = yaml.Unmarshal(yml, &keys)
	if err != nil {
		return err
	}

	(*l.node).Address(&keys[0].Address)

	return nil
}

// Generates the initial transaction(s) for genesis, for a validator.
func (l Localnet) validatorGenesis(address string) error {
	_, err := l.utils.AddGenesisAccount(address, Coins)
	if err != nil {
		return err
	}

	_, err = l.utils.GenerateGenesisTxn((*l.node).Moniker(), (*l.node).KeyPassword())
	if err != nil {
		return err
	}

	_, err = l.utils.CollectGenesisTxns()
	if err != nil {
		return err
	}

	err = (*l.node).CollectPeerAddress()
	if err != nil {
		return err
	}

	return nil
}

// Download the peer's genesis file and update the witness' peer config with the validator's peer address.
func (l Localnet) witnessGenesis(address string) error {
	genesis, err := l.utils.ScrapePeerGenesis((*l.node).GenesisURL())
	if err != nil {
		return err
	}

	err = l.utils.SaveGenesis(genesis)
	if err != nil {
		return err
	}

	err = l.utils.ReplacePeerConfig([]string{(*l.node).PeerAddress()})
	if err != nil {
		return err
	}

	return nil
}
