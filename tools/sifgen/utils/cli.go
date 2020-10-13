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

type CLIUtils interface {
	Reset() error
	CurrentChainID() (*string, error)
	InitChain(string, string) (*string, error)
	SetKeyRingStorage() (*string, error)
	SetConfigChainID(string) (*string, error)
	SetConfigIndent(bool) (*string, error)
	SetConfigTrustNode(bool) (*string, error)
	AddKey(string, string) (*string, error)
	AddGenesisAccount(string, []string) (*string, error)
	GenerateGenesisTxn(string, string) (*string, error)
	CollectGenesisTxns() (*string, error)
	ExportGenesis() (*string, error)
	GenesisFilePath() string
	ConfigFilePath() string
	TransferFunds(string, string, string, string) (*string, error)
	ValidatorPublicKeyAddress() (*string, error)
	CreateValidator(string, string, string, string) (*string, error)
}

type CLI struct {
	chainID    string
	sifDaemon  string
	sifCLI     string
	configPath string
}

func NewCLI(chainID string) CLI {
	return CLI{
		chainID:    chainID,
		sifDaemon:  os.Getenv("SIF_DAEMON"),
		sifCLI:     os.Getenv("SIF_CLI"),
		configPath: fmt.Sprintf("%s/config", app.DefaultNodeHome),
	}
}

func (c CLI) Reset() error {
	paths := []string{app.DefaultNodeHome, app.DefaultCLIHome}
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

func (c CLI) CurrentChainID() (*string, error) {
	return c.shellExec(c.sifCLI, "config", "chain-id", "--get")
}

func (c CLI) InitChain(chainID, moniker string) (*string, error) {
	return c.shellExec(c.sifDaemon, "init", moniker, "--chain-id", chainID)
}

func (c CLI) SetKeyRingStorage() (*string, error) {
	return c.shellExec(c.sifCLI, "config", "keyring-backend", "file")
}

func (c CLI) SetConfigChainID(chainID string) (*string, error) {
	return c.shellExec(c.sifCLI, "config", "chain-id", chainID)
}

func (c CLI) SetConfigIndent(indent bool) (*string, error) {
	return c.shellExec(c.sifCLI, "config", "indent", fmt.Sprintf("%v", indent))
}

func (c CLI) SetConfigTrustNode(indent bool) (*string, error) {
	return c.shellExec(c.sifCLI, "config", "trust-node", fmt.Sprintf("%v", indent))
}

func (c CLI) AddKey(name, keyPassword string) (*string, error) {
	return c.shellExecInput(c.sifCLI, [][]byte{[]byte(keyPassword + "\n"), []byte(keyPassword + "\n")}, "keys", "add", name)
}

func (c CLI) AddGenesisAccount(address string, coins []string) (*string, error) {
	return c.shellExec(c.sifDaemon, "add-genesis-account", address, strings.Join(coins[:], ","))
}

func (c CLI) GenerateGenesisTxn(name, keyPassword string) (*string, error) {
	return c.shellExecInput(c.sifDaemon,
		[][]byte{[]byte(keyPassword + "\n"), []byte(keyPassword + "\n"), []byte(keyPassword + "\n")},
		"gentx", "--name", name, "--keyring-backend", "file",
	)
}

func (c CLI) CollectGenesisTxns() (*string, error) {
	return c.shellExec(c.sifDaemon, "collect-gentxs")
}

func (c CLI) ExportGenesis() (*string, error) {
	return c.shellExec(c.sifDaemon, "export")
}

func (c CLI) GenesisFilePath() string {
	return fmt.Sprintf("%s/%s", c.configPath, GenesisFile)
}

func (c CLI) ConfigFilePath() string {
	return fmt.Sprintf("%s/%s", c.configPath, ConfigFile)
}

func (c CLI) TransferFunds(keyPassword, fromAddress, toAddress, coins string) (*string, error) {
	return c.shellExecInput(c.sifCLI,
		[][]byte{
			[]byte(keyPassword + "\n"),
			[]byte(keyPassword + "\n"),
		}, "tx", "send", fromAddress, toAddress, coins, "-y")
}

func (c CLI) ValidatorPublicKeyAddress() (*string, error) {
	return c.shellExec(c.sifDaemon, "tendermint", "show-validator")
}

func (c CLI) CreateValidator(moniker, validatorPublicKey, keyPassword, bondAmount string) (*string, error) {
	return c.shellExecInput(c.sifCLI,
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
