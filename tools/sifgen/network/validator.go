package network

import (
	"fmt"
	"net"
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/tyler-smith/go-bip39"
	"github.com/yelinaung/go-haikunator"
)

const (
	NodeHomeDir   = ".sifnoded"
	ConfigDir     = "config"
	GentxsDir     = "gentxs"
	ValidatorsDir = "validators"
)

type Validators []Validator
type Validator struct {
	ChainID                   string `yaml:"chain_id"`
	NodeID                    string `yaml:"node_id"`
	IPv4Address               string `yaml:"ipv4_address"`
	HomeDir                   string `yaml:"-"`
	NodeHomeDir               string `yaml:"-"`
	Moniker                   string `yaml:"moniker"`
	Password                  string `yaml:"password"`
	Address                   string `yaml:"address"`
	PubKey                    string `yaml:"pub_key"`
	Mnemonic                  string `yaml:"mnemonic"`
	ValidatorAddress          string `yaml:"validator_address"`
	ValidatorConsensusAddress string `yaml:"validator_consensus_address"`
	Seed                      bool   `yaml:"is_seed"`
}

func NewValidator(rootDir, chainID string, seed bool, lastIPv4Addr string) *Validator {
	moniker := haikunator.New(time.Now().UTC().UnixNano()).Haikunate()
	homeDir := fmt.Sprintf("%s/%s/%s/%s", rootDir, ValidatorsDir, chainID, moniker)

	return &Validator{
		IPv4Address: nextIP(lastIPv4Addr),
		ChainID:     chainID,
		HomeDir:     homeDir,
		NodeHomeDir: fmt.Sprintf("%s/%s", homeDir, NodeHomeDir),
		Moniker:     moniker,
		Password:    generatePassword(),
		Mnemonic:    generateMnemonic(),
		Seed:        seed,
	}
}

func generatePassword() string {
	nodePassword, _ := password.Generate(32, 5, 0, false, false)
	return nodePassword
}

func generateMnemonic() string {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	return mnemonic
}

func nextIP(lastIPv4Addr string) string {
	ip := net.ParseIP(lastIPv4Addr)
	ip = ip.To4()
	ip[3]++

	return ip.String()
}
