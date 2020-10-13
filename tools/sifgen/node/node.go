package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Sifchain/sifnode/tools/sifgen/node/types"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"

	"github.com/BurntSushi/toml"
	"github.com/sethvargo/go-password/password"
	"gopkg.in/yaml.v3"
)

type Node struct {
	chainID                       string
	moniker                       string
	nodeKeyAddress                string
	nodeKeyPassword               string
	nodePeerAddress               string
	nodeValidatorPublicKeyAddress string
	seedAddress                   *string
	genesisURL                    *string
	CLI                           utils.CLIUtils
}

func NewNode(chainID string, moniker, seedAddress, genesisURL *string) *Node {
	return &Node{
		chainID:     chainID,
		moniker:     *moniker,
		seedAddress: seedAddress,
		genesisURL:  genesisURL,
		CLI:         utils.NewCLI(chainID),
	}
}

// Validate config.
func (n *Node) Validate() error {
	if err := n.validateChainID(); err != nil {
		return err
	}

	if err := n.validateMoniker(); err != nil {
		return err
	}

	return nil
}

// Pre-flight setup.
func (n *Node) Setup() error {
	err := n.CLI.Reset()
	if err != nil {
		return err
	}

	_, err = n.CLI.InitChain(n.chainID, n.moniker)
	if err != nil {
		return err
	}

	_, err = n.CLI.SetKeyRingStorage()
	if err != nil {
		return err
	}

	_, err = n.CLI.SetConfigChainID(n.chainID)
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

	err = n.generateNodeKeyAddress()
	if err != nil {
		return err
	}

	return nil
}

// Genesis init.
func (n *Node) Genesis(deposit []string) error {
	if n.seedAddress == nil {
		return n.seedGenesis(deposit)
	}

	return n.validatorGenesis()
}

// Promote to a full validator.
func (n *Node) Promote(validatorPublicKey, keyPassword, bondAmount string) error {
	_, err := n.CLI.CreateValidator(n.moniker, validatorPublicKey, keyPassword, bondAmount)
	if err != nil {
		return err
	}

	return nil
}

// Get Chain ID.
func (n *Node) ChainID() string {
	return n.chainID
}

// Get node moniker.
func (n *Node) Moniker() string {
	return n.moniker
}

// Set/Get the node key address.
func (n *Node) NodeKeyAddress(address *string) *string {
	if address != nil {
		n.nodeKeyAddress = *address
	} else {
		return &n.nodeKeyAddress
	}

	return nil
}

// Get the node key password.
func (n *Node) NodeKeyPassword() string {
	return n.nodeKeyPassword
}

// Get the node peer address.
func (n *Node) NodePeerAddress() string {
	return n.nodePeerAddress
}

// Get the node validator public key address.
func (n *Node) NodeValidatorPublicKeyAddress() string {
	_ = n.collectNodeValidatorPublicKeyAddress()
	return n.nodeValidatorPublicKeyAddress
}

// Get the node seed address.
func (n *Node) SeedAddress() *string {
	return n.seedAddress
}

// Update the peer list.
func (n *Node) UpdatePeerList(peerList []string) error {
	err := n.replacePeerConfig(peerList)
	if err != nil {
		return err
	}

	return nil
}

// Generate a new key for a node.
func (n *Node) generateNodeKeyAddress() error {
	if err := n.generateNodeKeyPassword(); err != nil {
		return err
	}

	output, err := n.CLI.AddKey(n.moniker, n.nodeKeyPassword)
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

	n.nodeKeyAddress = keys[0].Address

	return nil
}

// Generate a password for the new node key.
func (n *Node) generateNodeKeyPassword() error {
	keyPassword, err := password.Generate(32, 5, 0, false, false)
	if err != nil {
		return err
	}

	n.nodeKeyPassword = keyPassword
	return nil
}

// Generates the initial transaction(s) for genesis, for a seed.
func (n *Node) seedGenesis(deposit []string) error {
	_, err := n.CLI.AddGenesisAccount(n.nodeKeyAddress, deposit)
	if err != nil {
		return err
	}

	_, err = n.CLI.GenerateGenesisTxn(n.moniker, n.nodeKeyPassword)
	if err != nil {
		return err
	}

	_, err = n.CLI.CollectGenesisTxns()
	if err != nil {
		return err
	}

	err = n.collectNodePeerAddress()
	if err != nil {
		return err
	}

	return nil
}

// Validate Chain ID.
func (n *Node) validateChainID() error {
	chainID, err := n.CLI.CurrentChainID()
	if err != nil {
		return err
	}

	if strings.TrimSpace(*chainID) != n.chainID {
		return errors.New("chain ID does not match")
	}

	return nil
}

// Validate the moniker.
func (n *Node) validateMoniker() error {
	config, err := n.parseConfig()
	if err != nil {
		return err
	}

	if config.Moniker != n.moniker {
		return errors.New("moniker does not match")
	}

	return nil
}

// Download the genesis file and update the peer config with the seeder's address.
func (n *Node) validatorGenesis() error {
	genesis, err := n.downloadGenesis()
	if err != nil {
		return err
	}

	if err = n.saveGenesis(genesis); err != nil {
		return err
	}

	err = n.replacePeerConfig([]string{*n.seedAddress})
	if err != nil {
		return err
	}

	return nil
}

// Download the genesis.
func (n *Node) downloadGenesis() (types.Genesis, error) {
	var genesis types.Genesis

	response, err := http.Get(fmt.Sprintf("%s", *n.genesisURL))
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

// Replace peer config.
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

// Parse config.toml.
func (n *Node) parseConfig() (types.Config, error) {
	var config types.Config

	content, err := ioutil.ReadFile(n.CLI.ConfigFilePath())
	if err != nil {
		return config, err
	}

	if _, err := toml.Decode(string(content), &config); err != nil {
		return config, err
	}

	return config, nil
}

// Save the genesis.
func (n *Node) saveGenesis(genesis types.Genesis) error {
	err := ioutil.WriteFile(n.CLI.GenesisFilePath(), *genesis.Result.Genesis, 0600)
	if err != nil {
		return err
	}

	return nil
}

// Collect our peer address from genesis.
func (n *Node) collectNodePeerAddress() error {
	output, err := n.CLI.ExportGenesis()
	if err != nil {
		return err
	}

	var genesisAppState types.GenesisAppState
	if err := json.Unmarshal([]byte(*output), &genesisAppState); err != nil {
		return err
	}

	n.nodePeerAddress = genesisAppState.AppState.Genutil.Gentxs[0].Value.Memo

	return nil
}

// Collect the validator public key address for the node from the key.
func (n *Node) collectNodeValidatorPublicKeyAddress() error {
	output, err := n.CLI.ValidatorPublicKeyAddress()
	if err != nil {
		return err
	}

	n.nodeValidatorPublicKeyAddress = *output

	return nil
}
