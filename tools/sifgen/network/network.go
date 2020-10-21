package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Sifchain/sifnode/tools/sifgen/network/types"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"
	"github.com/pelletier/go-toml"
	"github.com/sethvargo/go-password/password"
	"github.com/yelinaung/go-haikunator"
	"gopkg.in/yaml.v3"
)

type Network struct {
	chainID string
}

func NewNetwork(chainID string) *Network {
	return &Network{
		chainID: chainID,
	}
}

func (n *Network) Build(count int, outputDir, startingIPAddress string) error {
	cli := utils.NewCLI(n.chainID)

	var keylist types.NetworkKeys

	nodesDir := fmt.Sprintf("%s/nodes", outputDir)
	if err := cli.CreateDir(nodesDir); err != nil {
		return err
	}

	gentxsDir := fmt.Sprintf("%s/gentxs", outputDir)
	if err := cli.CreateDir(gentxsDir); err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		moniker := haikunator.New(time.Now().UTC().UnixNano()).Haikunate()

		rootCliDir := fmt.Sprintf("%s/%s/%s/%s", nodesDir, n.chainID, moniker, ".sifnodecli")
		if err := cli.CreateDir(rootCliDir); err != nil {
			return err
		}

		rootNodeDir := fmt.Sprintf("%s/%s/%s/%s", nodesDir, n.chainID, moniker, ".sifnoded")
		if err := cli.CreateDir(rootNodeDir); err != nil {
			return err
		}

		configDir := fmt.Sprintf("%s/config", rootCliDir)
		if err := cli.CreateDir(configDir); err != nil {
			return err
		}

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

		if err = ioutil.WriteFile(fmt.Sprintf("%s/config.toml", configDir), data, 0600); err != nil {
			return err
		}

		keyPassword, err := password.Generate(32, 5, 0, false, false)
		if err != nil {
			return err
		}

		output, err := cli.AddKey(moniker, keyPassword, rootCliDir)
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

		keys[0].Moniker = moniker
		keys[0].Password = keyPassword

		// Initialize Genesis.
		_, err = cli.InitChain(n.chainID, moniker, rootNodeDir)
		if err != nil {
			return err
		}

		// Replace the staking denom.
		var genesis types.Genesis

		body, err := ioutil.ReadFile(cli.GenesisFilePath())
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &genesis); err != nil {
			return err
		}

		genesis.AppState.Staking.Params.BondDenom = types.BondDenom
		content, err := json.Marshal(genesis)
		if err != nil {
			return err
		}

		if err = ioutil.WriteFile(cli.GenesisFilePath(), content, 0600); err != nil {
			return err
		}

		// Get the Node ID.
		nodeID, err := cli.NodeID(rootNodeDir)
		if err != nil {
			return err
		}

		keys[0].NodeID = strings.TrimSuffix(*nodeID, "\n")

		// Add genesis accounts (to the primary).
		if i != 0 {
			dir := fmt.Sprintf("%s/%s/%s/%s", nodesDir, n.chainID, keylist[0][0].Moniker, ".sifnoded")
			_, _ = cli.AddGenesisAccount(keys[0].Address, dir, []string{"1000000000000000trowan"})
		} else {
			_, _ = cli.AddGenesisAccount(keys[0].Address, rootNodeDir, []string{"1000000000000000trowan"})
		}

		// GenerateTXs.
		if i != 0 {
			dir := fmt.Sprintf("%s/%s/%s/%s", nodesDir, n.chainID, keylist[0][0].Moniker, ".sifnoded")
			_, _ = cli.GenerateGenesisTxn(
				keys[0].Moniker,
				keys[0].Password,
				"1000000000trowan",
				dir,
				rootCliDir,
				fmt.Sprintf("%s/%s.json", gentxsDir, keys[0].Moniker),
				keys[0].NodeID,
			)
		} else {
			_, _ = cli.GenerateGenesisTxn(
				keys[0].Moniker,
				keys[0].Password,
				"1000000000trowan",
				rootNodeDir,
				rootCliDir,
				fmt.Sprintf("%s/%s.json", gentxsDir, keys[0].Moniker),
				keys[0].NodeID,
			)
		}

		keylist = append(keylist, keys)
	}

	// CollectTXs.
	rootNodeDir := fmt.Sprintf("%s/%s/%s/%s", nodesDir, n.chainID, keylist[0][0].Moniker, ".sifnoded")
	_, _ = cli.CollectGenesisTxns(gentxsDir, rootNodeDir)

	fmt.Printf("%+v\n", keylist)

	return nil
}
