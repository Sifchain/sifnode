package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
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
)

type CLIUtils interface {
	Reset([]string) error
	DaemonPath() (*string, error)
	ResetState(string) (*string, error)
	CreateDir(string) error
	MoveFile(string, string) (*string, error)
	NodeID(string) (*string, error)
	ValidatorAddress(string) (*string, error)
	ValidatorConsensusAddress(string) (*string, error)
	InitChain(string, string, string) (*string, error)
	AddKey(string, string, string, string) (*string, error)
	AddGenesisAccount(string, string, []string) (*string, error)
	AddGenesisCLPAdmin(string, string) (*string, error)
	SetGenesisOracleAdmin(string, string) (*string, error)
	GenerateGenesisTxn(string, string, string, string, string, string, string, string) (*string, error)
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
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c CLI) DaemonPath() (*string, error) {
	return c.shellExec("which", "sifnoded")
}

func (c CLI) ResetState(nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "unsafe-reset-all", "--home", nodeDir)
}

func (c CLI) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (c CLI) MoveFile(src, dest string) (*string, error) {
	return c.shellExec("mv", src, dest)
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

func (c CLI) AddKey(name, mnemonic, keyPassword, cliDir string) (*string, error) {
	return c.shellExecInput("sifnoded",
		[][]byte{
			[]byte(mnemonic + "\n"),
			[]byte("\n"),
			[]byte(keyPassword + "\n"),
			[]byte(keyPassword + "\n"),
		}, "keys", "add", name, "--home", cliDir, "-i", "--keyring-backend", "file")
}

func (c CLI) AddGenesisAccount(address, nodeDir string, coins []string) (*string, error) {
	return c.shellExec("sifnoded", "add-genesis-account", address, strings.Join(coins[:], ","), "--home", nodeDir)
}

func (c CLI) AddGenesisCLPAdmin(address, nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "add-genesis-clp-admin", address, "--home", nodeDir)
}

func (c CLI) SetGenesisOracleAdmin(address, nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "set-genesis-oracle-admin", address, "--home", nodeDir)
}

func (c CLI) GenerateGenesisTxn(name, keyPassword, bondAmount, nodeDir, outputFile, nodeID, pubKey, ipV4Addr string) (*string, error) {
	return c.shellExecInput("sifnoded",
		[][]byte{[]byte(keyPassword + "\n"), []byte(keyPassword + "\n"), []byte(keyPassword + "\n")},
		"gentx",
		"--name", name,
		"--details", name,
		"--amount", bondAmount,
		"--keyring-backend", "file",
		"--home", nodeDir,
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
	return c.shellExecInput("sifnoded",
		[][]byte{
			[]byte(keyPassword + "\n"),
			[]byte(keyPassword + "\n"),
		}, "tx", "send", fromAddress, toAddress, coins, "-y")
}

func (c CLI) ValidatorPublicKeyAddress() (*string, error) {
	return c.shellExec("sifnoded", "tendermint", "show-validator")
}

func (c CLI) CreateValidator(moniker, validatorPublicKey, keyPassword, bondAmount string) (*string, error) {
	return c.shellExecInput("sifnoded",
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
	var errOut bytes.Buffer
	cm.Stdout = &out
	cm.Stderr = &errOut

	err := cm.Run()
	if err != nil {
		return nil, fmt.Errorf("error executing %s %s: %s \n %s", cmd, strings.Join(args, " "), err.Error(), errOut.String())
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
		_, err := io.Copy(buf, stdout)
		if err != nil {
			log.Println("io.Copy failed: ", err.Error())
		}
	}()

	if err := cm.Start(); err != nil {
		panic(err)
	}

	for _, i := range inputs {
		_, err := stdin.Write(i)
		if err != nil {
			log.Println("Write failed: ", err.Error())
		}
	}

	if err := cm.Wait(); err != nil {
		fmt.Println(stderr.String())
		return nil, err
	}

	result := buf.String()
	return &result, nil
}
