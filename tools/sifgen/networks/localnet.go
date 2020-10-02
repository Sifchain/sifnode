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

func (l Localnet) Setup() error {
	err := l.utils.Reset([]string{l.defaultNodeHome, l.defaultCLIHome})
	if err != nil {
		return err
	}

	_, err = l.utils.InitChain(l.chainID, (*l.node).Name())
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

func (l Localnet) generateNodeKey() error {
	output, err := l.utils.AddKey((*l.node).Name(), (*l.node).KeyPassword())
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

func (l Localnet) validatorGenesis(address string) error {
	_, err := l.utils.AddGenesisAccount(address, Coins)
	if err != nil {
		return err
	}

	_, err = l.utils.GenerateGenesisTxn((*l.node).Name(), (*l.node).KeyPassword())
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
