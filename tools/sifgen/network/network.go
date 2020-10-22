package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Sifchain/sifnode/tools/sifgen/network/types"
	"github.com/Sifchain/sifnode/tools/sifgen/network/utils"
	"github.com/pelletier/go-toml"
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

func (n *Network) initNodes(count int, outputDir string) []*Node {
	var nodes []*Node
	for i := 0; i < count; i++ {
		seed := false
		if i == 0 {
			seed = true
		}

		nodes = append(nodes, NewNode(outputDir, n.chainID, seed))
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
	config := types.Config{
		ChainID:        n.chainID,
		Indent:         true,
		KeyringBackend: "file",
		TrustNode:      true,
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(configPath, data, 0600); err != nil {
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

func (n *Network) Build(count int, outputDir, startingIPAddress string) error {
	initDirs := []string{
		fmt.Sprintf("%s/%s", outputDir, NodesDir),
		fmt.Sprintf("%s/%s", outputDir, GentxsDir),
	}

	if err := n.createDirs(initDirs); err != nil {
		return err
	}

	gentxDir := fmt.Sprintf("%s/%s", outputDir, GentxsDir)
	fmt.Println(gentxDir)

	nodes := n.initNodes(count, outputDir)

	for _, node := range nodes {
		node.GeneratePassword()

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

		fmt.Printf("%+v\n\n", node)
	}

	seedNode := n.getSeedNode(nodes)
	if err := n.collectGenTxs(gentxDir, seedNode.NodeHomeDir); err != nil {
		return err
	}

	// TODO: update the peer lists in the config files.

	return nil
}
