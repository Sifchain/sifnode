package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/Sifchain/sifnode/tools/sifgen/genesis"
	"github.com/Sifchain/sifnode/tools/sifgen/key"
	"github.com/Sifchain/sifnode/tools/sifgen/node/types"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"

	"github.com/BurntSushi/toml"
	"github.com/sethvargo/go-password/password"
	"github.com/yelinaung/go-haikunator"
	"gopkg.in/yaml.v3"
)

type Node struct {
	ChainID     string    `yaml:"chain_id"`
	PeerAddress *string   `yaml:"-"`
	GenesisURL  *string   `yaml:"-"`
	Moniker     string    `yaml:"moniker"`
	Address     string    `yaml:"address"`
	Password    string    `yaml:"password"`
	Mnemonic    *string   `yaml:"mnemonic"`
	Key         *key.Key  `yaml:"-"`
	CLI         utils.CLI `yaml:"-"`
}

func Reset(chainID string, nodeDir *string) error {
	var directory string
	if nodeDir == nil {
		directory = common.DefaultNodeHome
	} else {
		directory = *nodeDir
	}

	_, err := utils.NewCLI(chainID).ResetState(directory)
	if err != nil {
		return err
	}

	return nil
}

func NewNode(chainID string, moniker, mnemonic, peerAddress, genesisURL *string) *Node {
	password, _ := password.Generate(32, 5, 0, false, false)
	if moniker == nil {
		name := haikunator.New(time.Now().UTC().UnixNano()).Haikunate()
		moniker = &name
	}

	return &Node{
		ChainID:     chainID,
		PeerAddress: peerAddress,
		GenesisURL:  genesisURL,
		Moniker:     *moniker,
		Password:    password,
		Mnemonic:    mnemonic,
		CLI:         utils.NewCLI(chainID),
		Key:         key.NewKey(moniker, &password),
	}
}

func (n *Node) Build() (*string, error) {
	if err := n.setup(); err != nil {
		return nil, err
	}

	if err := n.genesis(); err != nil {
		return nil, err
	}

	summary := n.summary()
	return &summary, nil
}

func (n *Node) setup() error {
	if err := n.CLI.Reset([]string{common.DefaultNodeHome, common.DefaultCLIHome}); err != nil {
		return err
	}

	_, err := n.CLI.InitChain(n.ChainID, n.Moniker, common.DefaultNodeHome)
	if err != nil {
		return err
	}

	_, err = n.CLI.SetKeyRingStorage()
	if err != nil {
		return err
	}

	_, err = n.CLI.SetConfigChainID(n.ChainID)
	if err != nil {
		return err
	}

	_, err = n.CLI.SetConfigIndent(true)
	if err != nil {
		return err
	}

	_, err = n.CLI.SetConfigTrustNode(true)
	if err != nil {
		return err
	}

	if n.Mnemonic == nil || *n.Mnemonic == "" {
		err = n.generateMnemonic()
		if err != nil {
			return err
		}
	}

	err = n.generateNodeKeyAddress()
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) genesis() error {
	if n.GenesisURL != nil {
		return n.networkGenesis()
	}

	return n.seedGenesis()
}

func (n *Node) networkGenesis() error {
	genesis, err := n.downloadGenesis()
	if err != nil {
		return err
	}

	if err = n.saveGenesis(genesis); err != nil {
		return err
	}

	err = n.replacePeerConfig([]string{*n.PeerAddress})
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) seedGenesis() error {
	_, err := n.CLI.AddGenesisAccount(n.Address, common.DefaultNodeHome, common.ToFund)
	if err != nil {
		return err
	}

	gentxDir, err := ioutil.TempDir("", "gentx")
	if err != nil {
		return err
	}

	outputFile := fmt.Sprintf("%s/%s", gentxDir, "gentx.json")
	nodeID, _ := n.CLI.NodeID(common.DefaultNodeHome)

	pubKey, err := n.CLI.ValidatorAddress(common.DefaultNodeHome)
	if err != nil {
		return err
	}

	_, err = n.CLI.GenerateGenesisTxn(
		n.Moniker,
		n.Password,
		common.ToBond,
		common.DefaultNodeHome,
		common.DefaultCLIHome,
		outputFile,
		strings.TrimSuffix(*nodeID, "\n"),
		strings.TrimSuffix(*pubKey, "\n"),
		"127.0.0.1")
	if err != nil {
		return err
	}

	_, err = n.CLI.CollectGenesisTxns(gentxDir, common.DefaultNodeHome)
	if err != nil {
		return err
	}

	if err = genesis.ReplaceStakingBondDenom(common.DefaultNodeHome); err != nil {
		return err
	}

	return nil
}

func (n *Node) generateNodeKeyAddress() error {
	output, err := n.CLI.AddKey(n.Moniker, *n.Mnemonic, n.Password, common.DefaultCLIHome)
	if err != nil {
		return err
	}

	yml, err := ioutil.ReadAll(strings.NewReader(*output))
	if err != nil {
		return err
	}

	var keys common.Keys

	err = yaml.Unmarshal(yml, &keys)
	if err != nil {
		return err
	}

	n.Address = keys[0].Address

	return nil
}

func (n *Node) generateMnemonic() error {
	n.Key.GenerateMnemonic()
	n.Mnemonic = &n.Key.Mnemonic

	return nil
}

func (n *Node) downloadGenesis() (types.Genesis, error) {
	var genesis types.Genesis

	response, err := http.Get(fmt.Sprintf("%v", *n.GenesisURL))
	if err != nil {
		return genesis, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return genesis, err
	}

	if err := json.Unmarshal(body, &genesis); err != nil {
		return genesis, err
	}

	return genesis, nil
}

func (n *Node) saveGenesis(genesis types.Genesis) error {
	err := ioutil.WriteFile(n.CLI.GenesisFilePath(), *genesis.Result.Genesis, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) replacePeerConfig(peerAddresses []string) error {
	config, err := n.parseConfig()
	if err != nil {
		return err
	}

	file, err := os.Create(n.CLI.ConfigFilePath())
	if err != nil {
		return err
	}

	config.P2P.PersistentPeers = strings.Join(peerAddresses[:], ",")
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func (n *Node) parseConfig() (common.NodeConfig, error) {
	var config common.NodeConfig

	content, err := ioutil.ReadFile(n.CLI.ConfigFilePath())
	if err != nil {
		return config, err
	}

	if _, err := toml.Decode(string(content), &config); err != nil {
		return config, err
	}

	n.updateConfigDefaults(&config)

	return config, nil
}

func (n *Node) updateConfigDefaults(config *common.NodeConfig) {
	config.P2P.MaxNumInboundPeers = common.MaxNumInboundPeers
	config.P2P.MaxNumOutboundPeers = common.MaxNumOutboundPeers
	config.P2P.AllowDuplicateIP = common.AllowDuplicateIP
}

func (n *Node) summary() string {
	yml, _ := yaml.Marshal(n)
	return string(yml)
}
