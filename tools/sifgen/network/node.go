package network

import (
	"fmt"
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
	ChainID                   string
	NodeID                    string
	IPv4Address               string
	HomeDir                   string
	NodeHomeDir               string
	CLIHomeDir                string
	CLIConfigDir              string
	Moniker                   string
	Password                  string
	Address                   string
	PubKey                    string
	ValidatorAddress          string
	ValidatorConsensusAddress string
	Seed                      bool
}

func NewNode(rootDir, chainID string, seed bool) *Node {
	moniker := haikunator.New(time.Now().UTC().UnixNano()).Haikunate()
	homeDir := fmt.Sprintf("%s/%s/%s/%s", rootDir, NodesDir, chainID, moniker)

	return &Node{
		ChainID:      chainID,
		HomeDir:      homeDir,
		NodeHomeDir:  fmt.Sprintf("%s/%s", homeDir, NodeHomeDir),
		CLIHomeDir:   fmt.Sprintf("%s/%s", homeDir, CLIHomeDir),
		CLIConfigDir: fmt.Sprintf("%s/%s/%s", homeDir, CLIHomeDir, ConfigDir),
		Moniker:      moniker,
		Seed:         seed,
	}
}

func (nd *Node) GeneratePassword() {
	nodePassword, _ := password.Generate(32, 5, 0, false, false)
	nd.Password = nodePassword
}
