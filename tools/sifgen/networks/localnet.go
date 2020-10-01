package networks

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Sifchain/sifnode/tools/sifgen/networks/types"

	"gopkg.in/yaml.v3"
)

type Localnet struct {
	defaultNodeHome string
	defaultCLIHome  string
	chainID         string
	utils           Utils
	networkNode     *NetworkNode
}

func NewLocalnet(defaultNodeHome, defaultCLIHome, chainID string, networkNode *NetworkNode) Localnet {
	return Localnet{
		defaultNodeHome: defaultNodeHome,
		defaultCLIHome:  defaultCLIHome,
		chainID:         chainID,
		utils:           NewUtils(defaultNodeHome),
		networkNode:     networkNode,
	}
}

func (l Localnet) Reset() {
	if _, err := os.Stat(l.defaultNodeHome); !os.IsNotExist(err) {
		err = os.RemoveAll(l.defaultNodeHome)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(l.defaultCLIHome); !os.IsNotExist(err) {
		err = os.RemoveAll(l.defaultCLIHome)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (l Localnet) Setup() {
	l.utils.InitChain(l.chainID, (*l.networkNode).Name())
	l.utils.SetKeyRingStorage()
	l.utils.SetConfigChainID(l.chainID)
	l.utils.SetConfigIndent(true)
	l.utils.SetConfigTrustNode(true)
	l.generateNodeKey()
}

func (l Localnet) Genesis() {
	address := (*l.networkNode).Address(nil)

	_, isValidator := (*l.networkNode).(*Validator)
	if isValidator {
		l.validatorGenesis(*address)
	} else {
		l.witnessGenesis(*address)
	}
}

func (l Localnet) generateNodeKey() {
	yml, err := ioutil.ReadAll(strings.NewReader(l.utils.AddKey((*l.networkNode).Name(), (*l.networkNode).KeyPassword())))
	if err != nil {
		log.Fatal(err)
	}

	var keys types.Keys

	err = yaml.Unmarshal(yml, &keys)
	if err != nil {
		log.Fatal(err)
	}

	(*l.networkNode).Address(&keys[0].Address)
}

func (l Localnet) validatorGenesis(address string) {
	l.utils.AddGenesisAccount(address, Coins)
	l.utils.GenerateGenesisTxn((*l.networkNode).Name(), (*l.networkNode).KeyPassword())
	l.utils.CollectGenesisTxns()
	(*l.networkNode).CollectPeerAddress()
}

func (l Localnet) witnessGenesis(address string) {
	l.utils.SaveGenesis(l.utils.ScrapePeerGenesis((*l.networkNode).GenesisURL()))
	l.utils.ReplacePeerConfig([]string{(*l.networkNode).PeerAddress()})
}
