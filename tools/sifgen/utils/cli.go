package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/tools/sifgen/node/types"
)

const (
	GenesisFile = "genesis.json"
	ConfigFile  = "config.toml"
)

type CLIUtils interface {
	Reset() error
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
	ScrapePeerGenesis(string) (types.Genesis, error)
	SaveGenesis(types.Genesis) error
	ReplacePeerConfig([]string) error
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

func (c CLI) ScrapePeerGenesis(url string) (types.Genesis, error) {
	var genesis types.Genesis

	response, err := http.Get(fmt.Sprintf("%s", url))
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

func (c CLI) SaveGenesis(genesis types.Genesis) error {
	err := ioutil.WriteFile(c.GenesisFilePath(), *genesis.Result.Genesis, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (c CLI) ReplacePeerConfig(peerAddresses []string) error {
	contents, err := ioutil.ReadFile(c.ConfigFilePath())
	if err != nil {
		return err
	}

	lines := strings.Split(string(contents), "\n")
	for i, line := range lines {
		if strings.Contains(line, "persistent_peers = \"\"") {
			lines[i] = fmt.Sprintf("persistent_peers = \"%s\"", strings.Join(peerAddresses[:], ","))
		}
	}

	output := strings.Join(lines, "\n")

	err = ioutil.WriteFile(c.ConfigFilePath(), []byte(output), 0600)
	if err != nil {
		return err
	}

	return nil
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
