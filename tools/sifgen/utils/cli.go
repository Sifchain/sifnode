package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Sifchain/sifnode/app"
)

const (
	GenesisFile = "genesis.json"
	ConfigFile  = "config.toml"
)

var (
	DefaultNodeHome = app.DefaultNodeHome
	DefaultCLIHome  = app.DefaultCLIHome
)

type CLIUtils interface {
	Reset() error
	CreateDir(string) error
	CurrentChainID() (*string, error)
	NodeID(nodeDir string) (*string, error)
	ValidatorAddress(nodeDir string) (*string, error)
	ValidatorConsensusAddress(nodeDir string) (*string, error)
	InitChain(string, string, string) (*string, error)
	SetKeyRingStorage() (*string, error)
	SetConfigChainID(string) (*string, error)
	SetConfigIndent(bool) (*string, error)
	SetConfigTrustNode(bool) (*string, error)
	AddKey(string, string, string) (*string, error)
	AddGenesisAccount(string, string, []string) (*string, error)
	GenerateGenesisTxn(string, string, string, string, string, string, string, string, string) (*string, error)
	CollectGenesisTxns(string, string) (*string, error)
	ExportGenesis() (*string, error)
	GenesisFilePath() string
	ConfigFilePath() string
	TransferFunds(string, string, string, string) (*string, error)
	ValidatorPublicKeyAddress() (*string, error)
	CreateValidator(string, string, string, string) (*string, error)
}

type CLI struct {
	chainID    string
	configPath string
}

func NewCLI(chainID string) CLI {
	return CLI{
		chainID:    chainID,
		configPath: fmt.Sprintf("%s/config", app.DefaultNodeHome),
	}
}

func (c CLI) Reset(paths []string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			err = os.RemoveAll(path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c CLI) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (c CLI) CurrentChainID() (*string, error) {
	return c.shellExec("sifnodecli", "config", "chain-id", "--get")
}

func (c CLI) NodeID(nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "tendermint", "show-node-id", "--home", nodeDir)
}

func (c CLI) ValidatorAddress(nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "tendermint", "show-validator", "--home", nodeDir)
}

func (c CLI) ValidatorConsensusAddress(nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "tendermint", "show-address", "--home", nodeDir)
}

func (c CLI) InitChain(chainID, moniker, nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "init", moniker, "--chain-id", chainID, "--home", nodeDir)
}

func (c CLI) SetKeyRingStorage() (*string, error) {
	return c.shellExec("sifnodecli", "config", "keyring-backend", "file")
}

func (c CLI) SetConfigChainID(chainID string) (*string, error) {
	return c.shellExec("sifnodecli", "config", "chain-id", chainID)
}

func (c CLI) SetConfigIndent(indent bool) (*string, error) {
	return c.shellExec("sifnodecli", "config", "indent", fmt.Sprintf("%v", indent))
}

func (c CLI) SetConfigTrustNode(indent bool) (*string, error) {
	return c.shellExec("sifnodecli", "config", "trust-node", fmt.Sprintf("%v", indent))
}

func (c CLI) AddKey(name, keyPassword, cliDir string) (*string, error) {
	return c.shellExecInput("sifnodecli",
		[][]byte{
			[]byte(keyPassword + "\n"),
			[]byte(keyPassword + "\n"),
		}, "keys", "add", name, "--home", cliDir, "--keyring-backend", "file")
}

func (c CLI) AddGenesisAccount(address, nodeDir string, coins []string) (*string, error) {
	return c.shellExec("sifnoded", "add-genesis-account", address, strings.Join(coins[:], ","), "--home", nodeDir)
}

func (c CLI) GenerateGenesisTxn(name, keyPassword, bondAmount, nodeDir, cliDir, outputFile, nodeID, pubKey, ipV4Addr string) (*string, error) {
	return c.shellExecInput("sifnoded",
		[][]byte{[]byte(keyPassword + "\n"), []byte(keyPassword + "\n"), []byte(keyPassword + "\n")},
		"gentx",
		"--name", name,
		"--details", name,
		"--amount", bondAmount,
		"--keyring-backend", "file",
		"--home", nodeDir,
		"--home-client", cliDir,
		"--output-document", outputFile,
		"--node-id", nodeID,
		"--pubkey", pubKey,
		"--ip", ipV4Addr,
	)
}

func (c CLI) CollectGenesisTxns(gentxDir, nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "collect-gentxs", "--gentx-dir", gentxDir, "--home", nodeDir)
}

func (c CLI) ExportGenesis() (*string, error) {
	return c.shellExec("sifnoded", "export")
}

func (c CLI) GenesisFilePath() string {
	return fmt.Sprintf("%s/%s", c.configPath, GenesisFile)
}

func (c CLI) ConfigFilePath() string {
	return fmt.Sprintf("%s/%s", c.configPath, ConfigFile)
}

func (c CLI) TransferFunds(keyPassword, fromAddress, toAddress, coins string) (*string, error) {
	return c.shellExecInput("sifnodecli",
		[][]byte{
			[]byte(keyPassword + "\n"),
			[]byte(keyPassword + "\n"),
		}, "tx", "send", fromAddress, toAddress, coins, "-y")
}

func (c CLI) ValidatorPublicKeyAddress() (*string, error) {
	return c.shellExec("sifnoded", "tendermint", "show-validator")
}

func (c CLI) CreateValidator(moniker, validatorPublicKey, keyPassword, bondAmount string) (*string, error) {
	return c.shellExecInput("sifnodecli",
		[][]byte{
			[]byte(keyPassword + "\n"),
			[]byte(keyPassword + "\n"),
		},
		"tx", "staking", "create-validator",
		"--commission-max-change-rate", "0.1",
		"--commission-max-rate", "0.1",
		"--commission-rate", "0.1",
		"--amount", bondAmount,
		"--pubkey", validatorPublicKey,
		"--moniker", moniker,
		"--chain-id", c.chainID,
		"--min-self-delegation", "1",
		"--gas", "auto",
		"--from", moniker,
		"--keyring-backend", "file",
		"-y")
}

func (c CLI) shellExec(cmd string, args ...string) (*string, error) {
	cm := exec.Command(cmd, args...)
	var out bytes.Buffer
	cm.Stdout = &out

	err := cm.Run()
	if err != nil {
		return nil, err
	}

	result := out.String()
	return &result, nil
}

func (c CLI) shellExecInput(cmd string, inputs [][]byte, args ...string) (*string, error) {
	cm := exec.Command(cmd, args...)
	var stderr bytes.Buffer
	cm.Stderr = &stderr

	stdin, err := cm.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cm.StdoutPipe()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	go func() {
		io.Copy(buf, stdout)
	}()

	if err := cm.Start(); err != nil {
		panic(err)
	}

	for _, i := range inputs {
		stdin.Write(i)
	}

	if err := cm.Wait(); err != nil {
		fmt.Println(stderr.String())
		return nil, err
	}

	result := buf.String()
	return &result, nil
}
