package network

import (
	"fmt"
	"net"
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/yelinaung/go-haikunator"
)

const (
	NodeHomeDir = ".sifnoded"
	CLIHomeDir  = ".sifnodecli"
	ConfigDir   = "config"
	GentxsDir   = "gentxs"
	NodesDir    = "nodes"
)

type Nodes []Node
type Node struct {
	ChainID                   string `yaml:"chain_id"`
	NodeID                    string `yaml:"node_id"`
	IPv4Address               string `yaml:"ipv4_address"`
	HomeDir                   string `yaml:"-"`
	NodeHomeDir               string `yaml:"-"`
	CLIHomeDir                string `yaml:"-"`
	CLIConfigDir              string `yaml:"-"`
	Moniker                   string `yaml:"moniker"`
	Password                  string `yaml:"password"`
	Address                   string `yaml:"address"`
	PubKey                    string `yaml:"pub_key"`
	ValidatorAddress          string `yaml:"validator_address"`
	ValidatorConsensusAddress string `yaml:"validator_consensus_address"`
	Seed                      bool   `yaml:"is_seed"`
}

func NewNode(rootDir, chainID string, seed bool, lastIPv4Addr string) *Node {
	moniker := haikunator.New(time.Now().UTC().UnixNano()).Haikunate()
	homeDir := fmt.Sprintf("%s/%s/%s/%s", rootDir, NodesDir, chainID, moniker)

	return &Node{
		IPv4Address:  nextIP(lastIPv4Addr),
		ChainID:      chainID,
		HomeDir:      homeDir,
		NodeHomeDir:  fmt.Sprintf("%s/%s", homeDir, NodeHomeDir),
		CLIHomeDir:   fmt.Sprintf("%s/%s", homeDir, CLIHomeDir),
		CLIConfigDir: fmt.Sprintf("%s/%s/%s", homeDir, CLIHomeDir, ConfigDir),
		Moniker:      moniker,
		Password:     generatePassword(),
		Seed:         seed,
	}
}

func generatePassword() string {
	nodePassword, _ := password.Generate(32, 5, 0, false, false)
	return nodePassword
}

func nextIP(lastIPv4Addr string) string {
	ip := net.ParseIP(lastIPv4Addr)
	ip = ip.To4()
	ip[3]++

	return ip.String()
}
