package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Sifchain/sifnode/tools/sifgen/network/types"
	"github.com/Sifchain/sifnode/tools/sifgen/network/utils"
	"gopkg.in/yaml.v3"
)

type Network struct {
	chainID string
	CLI     utils.CLI
}

func NewNetwork(chainID string) *Network {
	return &Network{
		chainID: chainID,
		CLI:     utils.NewCLI(chainID),
	}
}

func (n *Network) Build(count int, outputDir, seedIPv4Addr string) error {
	if err := n.CLI.Reset([]string{outputDir}); err != nil {
		return err
	}

	initDirs := []string{
		fmt.Sprintf("%s/%s", outputDir, NodesDir),
		fmt.Sprintf("%s/%s", outputDir, GentxsDir),
	}

	if err := n.createDirs(initDirs); err != nil {
		return err
	}

	gentxDir := fmt.Sprintf("%s/%s", outputDir, GentxsDir)
	nodes := n.initNodes(count, outputDir, seedIPv4Addr)

	for _, node := range nodes {
		appDirs := []string{node.NodeHomeDir, node.CLIHomeDir, node.CLIConfigDir}
		if err := n.createDirs(appDirs); err != nil {
			return err
		}

		if err := n.setDefaultConfig(fmt.Sprintf("%s/%s/%s/%s", node.HomeDir, CLIHomeDir, ConfigDir, utils.ConfigFile)); err != nil {
			return err
		}

		if err := n.generateKey(node); err != nil {
			return err
		}

		if err := n.initChain(node); err != nil {
			return err
		}

		if err := n.setValidatorAddress(node); err != nil {
			return err
		}

		if err := n.setValidatorConsensusAddress(node); err != nil {
			return err
		}

		if err := n.replaceStakingBondDenom(node); err != nil {
			return err
		}

		if err := n.setNodeID(node); err != nil {
			return err
		}

		if !node.Seed {
			seedNode := n.getSeedNode(nodes)
			if err := n.addGenesis(node.Address, seedNode.NodeHomeDir); err != nil {
				return err
			}

			if err := n.generateTx(node, seedNode.NodeHomeDir, gentxDir); err != nil {
				return err
			}
		} else {
			if err := n.addGenesis(node.Address, node.NodeHomeDir); err != nil {
				return err
			}

			if err := n.generateTx(node, node.NodeHomeDir, gentxDir); err != nil {
				return err
			}
		}
	}

	seedNode := n.getSeedNode(nodes)
	if err := n.collectGenTxs(gentxDir, seedNode.NodeHomeDir); err != nil {
		return err
	}

	if err := n.setPeers(nodes); err != nil {
		return err
	}

	n.summary(nodes)

	return nil
}

func (n *Network) summary(nodes []*Node) {
	for _, node := range nodes {
		yml, _ := yaml.Marshal(node)
		fmt.Println(string(yml))
	}
}

func (n *Network) initNodes(count int, outputDir, seedIPv4Addr string) []*Node {
	var nodes []*Node
	var lastIPv4Addr string

	for i := 0; i < count; i++ {
		seed := false
		if i == 0 {
			seed = true
		}

		if seed {
			lastIPv4Addr = seedIPv4Addr
		}

		node := NewNode(outputDir, n.chainID, seed, lastIPv4Addr)
		nodes = append(nodes, node)

		lastIPv4Addr = node.IPv4Address
	}

	return nodes
}

func (n *Network) createDirs(toCreate []string) error {
	for _, dir := range toCreate {
		if err := n.CLI.CreateDir(dir); err != nil {
			return err
		}
	}

	return nil
}

func (n *Network) setDefaultConfig(configPath string) error {
	config := types.CLIConfig{
		ChainID:        n.chainID,
		Indent:         true,
		KeyringBackend: "file",
		TrustNode:      true,
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func (n *Network) generateKey(node *Node) error {
	output, err := n.CLI.AddKey(node.Moniker, node.Password, fmt.Sprintf("%s/%s", node.HomeDir, CLIHomeDir))
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

	node.Address = keys[0].Address
	node.PubKey = keys[0].PubKey

	return nil
}

func (n *Network) initChain(node *Node) error {
	_, err := n.CLI.InitChain(node.ChainID, node.Moniker, node.NodeHomeDir)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) setValidatorAddress(node *Node) error {
	output, err := n.CLI.ValidatorAddress(node.NodeHomeDir)
	if err != nil {
		return err
	}

	node.ValidatorAddress = strings.TrimSuffix(*output, "\n")

	return nil
}

func (n *Network) setValidatorConsensusAddress(node *Node) error {
	output, err := n.CLI.ValidatorConsensusAddress(node.NodeHomeDir)
	if err != nil {
		return err
	}

	node.ValidatorConsensusAddress = strings.TrimSuffix(*output, "\n")

	return nil
}

func (n *Network) replaceStakingBondDenom(node *Node) error {
	var genesis types.Genesis

	genesisPath := fmt.Sprintf("%s/config/%s", node.NodeHomeDir, utils.GenesisFile)

	body, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &genesis); err != nil {
		return err
	}

	genesis.AppState.Staking.Params.BondDenom = types.TokenDenom
	content, err := json.Marshal(genesis)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(genesisPath, content, 0600); err != nil {
		return err
	}

	return nil
}

func (n *Network) setNodeID(node *Node) error {
	output, err := n.CLI.NodeID(node.NodeHomeDir)
	if err != nil {
		return err
	}

	node.NodeID = strings.TrimSuffix(*output, "\n")

	return nil
}

func (n *Network) getSeedNode(nodes []*Node) *Node {
	for _, node := range nodes {
		if node.Seed {
			return node
		}
	}

	return &Node{}
}

func (n *Network) addGenesis(address, nodeHomeDir string) error {
	_, err := n.CLI.AddGenesisAccount(address, nodeHomeDir, types.ToFund)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) generateTx(node *Node, nodeDir, outputDir string) error {
	_, err := n.CLI.GenerateGenesisTxn(
		node.Moniker,
		node.Password,
		types.ToBond,
		nodeDir,
		node.CLIHomeDir,
		fmt.Sprintf("%s/%s.json", outputDir, node.Moniker),
		node.NodeID,
		node.ValidatorAddress,
		node.IPv4Address,
	)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) collectGenTxs(gentxDir, nodeDir string) error {
	_, err := n.CLI.CollectGenesisTxns(gentxDir, nodeDir)
	if err != nil {
		return err
	}

	return nil
}

func (n *Network) generatePeerList(nodes []*Node, idx int) []string {
	var peers []string
	for i, node := range nodes {
		if i != idx {
			peers = append(peers, fmt.Sprintf("%s@%s:26657", node.NodeID, node.IPv4Address))
		}
	}

	return peers
}

func (n *Network) setPeers(nodes []*Node) error {
	for i, node := range nodes {
		var config types.NodeConfig

		configFile := fmt.Sprintf("%s/%s/%s", node.NodeHomeDir, ConfigDir, utils.ConfigFile)

		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}

		if _, err := toml.Decode(string(content), &config); err != nil {
			return err
		}

		file, err := os.Create(configFile)
		if err != nil {
			return err
		}

		config.P2P.PersistentPeers = strings.Join(n.generatePeerList(nodes, i)[:], ",")
		if err := toml.NewEncoder(file).Encode(config); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}
