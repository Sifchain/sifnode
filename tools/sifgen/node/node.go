package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/tools/sifgen/common"
	"github.com/Sifchain/sifnode/tools/sifgen/genesis"
	"github.com/Sifchain/sifnode/tools/sifgen/key"
	"github.com/Sifchain/sifnode/tools/sifgen/node/types"
	"github.com/Sifchain/sifnode/tools/sifgen/utils"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/sethvargo/go-password/password"
	"gopkg.in/yaml.v3"
)

type Node struct {
	CLI                       utils.CLI     `yaml:"-"`
	AdminCLPAddresses         []string      `yaml:"admin_clp_addresses"`
	AdminOracleAddress        string        `yaml:"admin_oracle_address"`
	AdminTokenRegistryAddress string        `yaml:"admin_token_registry_address"`
	ChainID                   string        `yaml:"chain_id"`
	Moniker                   string        `yaml:"moniker"`
	Mnemonic                  string        `yaml:"mnemonic"`
	IPAddr                    string        `yaml:"ip_address"`
	Address                   string        `yaml:"address"`
	Password                  string        `yaml:"password"`
	BondAmount                string        `yaml:"-"`
	MintAmount                string        `yaml:"-"`
	MinCLPCreatePoolThreshold uint64        `yaml:"-"`
	GovMaxDepositPeriod       time.Duration `yaml:"-"`
	GovVotingPeriod           time.Duration `yaml:"-"`
	PeerAddress               string        `yaml:"-"`
	GenesisURL                string        `yaml:"-"`
	Standalone                bool          `yaml:"-"`
	WithCosmovisor            bool          `yaml:"-"`
}

func Reset(chainID string, nodeDir *string) error {
	var directory string
	if nodeDir == nil {
		directory = common.DefaultNodeHome
	} else {
		directory = *nodeDir
	}

	_, err := utils.NewCLI(chainID, keyring.BackendTest).ResetState(directory)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) Build() (*string, error) {
	if _, err := os.Stat(fmt.Sprintf("%v/config/genesis.json", common.DefaultNodeHome)); err == nil {
		return nil, nil
	}

	if err := n.setup(); err != nil {
		return nil, err
	}

	if err := n.genesis(); err != nil {
		return nil, err
	}

	if err := n.setupCosmovisor(); err != nil {
		return nil, err
	}

	if err := n.cleanup(); err != nil {
		return nil, err
	}

	summary := n.summary()
	return &summary, nil
}

func (n *Node) setup() error {
	_, err := n.CLI.InitChain(n.ChainID, n.Moniker, common.DefaultNodeHome)
	if err != nil {
		return err
	}

	if err := n.generatePassword(); err != nil {
		return err
	}

	err = n.generateNodeKeyAddress()
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) genesis() error {
	if n.Standalone {
		return n.seedGenesis()
	}

	return n.networkGenesis()
}

func (n *Node) networkGenesis() error {
	genesis, err := n.downloadGenesis()
	if err != nil {
		return err
	}

	if err = n.saveGenesis(genesis); err != nil {
		return err
	}

	err = n.replaceConfigTOML()
	if err != nil {
		return err
	}

	err = n.replaceAppTOML()
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) seedGenesis() error {
	_, err := n.CLI.AddGenesisAccount(n.Address, common.DefaultNodeHome, []string{n.MintAmount})
	if err != nil {
		return err
	}

	if len(n.AdminCLPAddresses) != 0 {
		for _, adminAddress := range n.AdminCLPAddresses {
			_, err := n.CLI.AddGenesisCLPAdmin(adminAddress, common.DefaultNodeHome)
			if err != nil {
				return err
			}
		}
	}

	if n.AdminOracleAddress != "" {
		_, err = n.CLI.SetGenesisOracleAdmin(n.AdminOracleAddress, common.DefaultNodeHome)
		if err != nil {
			return err
		}
	}

	if n.AdminTokenRegistryAddress != "" {
		_, err = n.CLI.SetGenesisWhitelisterAdmin(n.AdminTokenRegistryAddress, common.DefaultNodeHome)
		if err != nil {
			return err
		}
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
		n.BondAmount,
		common.DefaultNodeHome,
		outputFile,
		strings.TrimSuffix(*nodeID, "\n"),
		strings.TrimSuffix(*pubKey, "\n"),
		n.IPAddr,
		n.ChainID)
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

	if err = genesis.ReplaceCLPMinCreatePoolThreshold(common.DefaultNodeHome, n.MinCLPCreatePoolThreshold); err != nil {
		return err
	}

	if err = genesis.ReplaceGovDepositParamsMinDeposit(common.DefaultNodeHome, common.StakeTokenDenom); err != nil {
		return err
	}

	if err = genesis.ReplaceGovDepositParamsMaxDepositPeriod(common.DefaultNodeHome, n.GovMaxDepositPeriod); err != nil {
		return err
	}

	if err = genesis.ReplaceGovVotingParamsVotingPeriod(common.DefaultNodeHome, n.GovVotingPeriod); err != nil {
		return err
	}

	if err = FEATURE_TOGGLE_MARGIN_CLI_ALPHA_seedGenesis(); err != nil {
		return err
	}

	err = n.replaceConfigTOML()
	if err != nil {
		return err
	}

	err = n.replaceAppTOML()
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) generatePassword() error {
	password, err := password.Generate(32, 5, 0, false, false)
	if err != nil {
		return err
	}

	n.Password = password

	return nil
}

func (n *Node) generateNodeKeyAddress() error {
	if n.Standalone {
		output, err := n.CLI.AddKey(n.Moniker, n.Mnemonic, n.Password, common.DefaultNodeHome)
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
	} else {
		tmpKeyring := key.NewKey(n.Moniker, n.Password)
		if err := tmpKeyring.RecoverFromMnemonic(n.Mnemonic); err != nil {
			return err
		}

		n.Address = tmpKeyring.Address
	}

	return nil
}

func (n *Node) downloadGenesis() (types.Genesis, error) {
	var genesis types.Genesis

	response, err := http.Get(fmt.Sprintf("%v", n.GenesisURL))
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

func (n *Node) replaceConfigTOML() error {
	config, err := n.parseConfigTOML()
	if err != nil {
		return err
	}

	file, err := os.Create(n.CLI.ConfigFilePath())
	if err != nil {
		return err
	}

	if !n.Standalone {
		addressList := []string{n.PeerAddress}
		config.P2P.PersistentPeers = strings.Join(addressList[:], ",")
	}

	if n.IPAddr != "" {
		config.P2P.ExternalAddress = fmt.Sprintf("%v:%v", n.IPAddr, common.P2PPort)
	}

	config.RPC.CorsAllowedOrigins = []string{"*"}
	config.RPC.CorsAllowedHeaders = []string{"*"}
	config.P2P.MaxNumInboundPeers = common.MaxNumInboundPeers
	config.P2P.MaxNumOutboundPeers = common.MaxNumOutboundPeers
	config.P2P.AllowDuplicateIP = common.AllowDuplicateIP

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func (n *Node) replaceAppTOML() error {
	config, err := n.parseAppTOML()
	if err != nil {
		return err
	}

	file, err := os.Create(n.CLI.AppFilePath())
	if err != nil {
		return err
	}

	config.API.Enable = true
	config.API.Swagger = true
	config.API.EnabledUnsafeCors = true
	config.Grpc.Enable = true

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func (n *Node) parseConfigTOML() (common.ConfigTOML, error) {
	var config common.ConfigTOML

	content, err := ioutil.ReadFile(n.CLI.ConfigFilePath())
	if err != nil {
		return config, err
	}

	if _, err := toml.Decode(string(content), &config); err != nil {
		return config, err
	}

	return config, nil
}

func (n *Node) parseAppTOML() (common.AppTOML, error) {
	var config common.AppTOML

	content, err := ioutil.ReadFile(n.CLI.AppFilePath())
	if err != nil {
		return config, err
	}

	if _, err := toml.Decode(string(content), &config); err != nil {
		return config, err
	}

	return config, nil
}

func (n *Node) setupCosmovisor() error {
	if !n.WithCosmovisor {
		return nil
	}

	if err := n.CLI.CreateDir(fmt.Sprintf("%v/cosmovisor/genesis/bin", app.DefaultNodeHome)); err != nil {
		return err
	}

	if err := n.CLI.CreateDir(fmt.Sprintf("%v/cosmovisor/upgrades", app.DefaultNodeHome)); err != nil {
		return err
	}

	path := fmt.Sprintf("%v:%v", os.Getenv("PATH"), fmt.Sprintf("%v/cosmovisor/genesis/bin", app.DefaultNodeHome))
	if err := os.Setenv("PATH", path); err != nil {
		return err
	}

	daemon, err := n.CLI.DaemonPath()
	if err != nil {
		return err
	}

	_, err = n.CLI.MoveFile(strings.TrimSuffix(*daemon, "\n"), fmt.Sprintf("%v/cosmovisor/genesis/bin/", app.DefaultNodeHome))
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) summary() string {
	yml, _ := yaml.Marshal(n)
	return string(yml)
}

func (n *Node) cleanup() error {
	if n.Standalone {
		_path := fmt.Sprintf("%v/%v", common.DefaultNodeHome, "keyring-test")
		dir, err := ioutil.ReadDir(_path)
		for _, d := range dir {
			_ = os.RemoveAll(path.Join([]string{_path, d.Name()}...))
		}

		if err != nil {
			return err
		}
	}

	return nil
}
