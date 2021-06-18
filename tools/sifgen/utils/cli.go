package utils

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Sifchain/sifnode/app"
)

const (
	GenesisFile = "genesis.json"
	ConfigTOML  = "config.toml"
	AppTOML = "app.toml"
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
	GenerateGenesisTxn(string, string, string, string, string, string, string, string, string) (*string, error)
	CollectGenesisTxns(string, string) (*string, error)
	ExportGenesis() (*string, error)
	GenesisFilePath() string
	ConfigFilePath() string
	AppFilePath() string
	TransferFunds(string, string, string, string) (*string, error)
	ValidatorPublicKeyAddress() (*string, error)
	CreateValidator(string, string, string, string) (*string, error)
}

type CLI struct {
	chainID        string
	configPath     string
	keyringBackend string
}

func NewCLI(chainID, keyringBackend string) CLI {
	return CLI{
		chainID:        chainID,
		configPath:     fmt.Sprintf("%s/config", app.DefaultNodeHome),
		keyringBackend: keyringBackend,
	}
}

func (c CLI) Reset(paths []string) error {
	for _, _path := range paths {
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
	switch c.keyringBackend {
	case keyring.BackendFile:
		return c.AddKeyToFileBackend(name, mnemonic, keyPassword, cliDir)
	default:
		var input [][]byte
		input = c.formatInputs([]string{mnemonic, ""})
		return c.shellExecInput("sifnoded", input, "keys", "add", name, "--home", cliDir, "-i", "--keyring-backend", c.keyringBackend)
	}
}

// AddKeyToFileBackend
//
// Adding a key to the file backend is different enough from the other backends that it's
// worth splitting it out.  This is usually only called by AddKey.  (It needs a few things
// from an interactive session - the mnemonic and the password repeated twice)
func (c CLI) AddKeyToFileBackend(name, mnemonic, keyPassword, cliDir string) (*string, error) {
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
	return c.shellExec("sifnoded", "add-genesis-clp-admin", address, "--home", nodeDir, "--keyring-backend", c.keyringBackend)
}

func (c CLI) SetGenesisOracleAdmin(address, nodeDir string) (*string, error) {
	return c.shellExec("sifnoded", "set-genesis-oracle-admin", address, "--home", nodeDir, "--keyring-backend", c.keyringBackend)
}

func (c CLI) GenerateGenesisTxn(name, keyPassword, bondAmount, nodeDir, outputFile, nodeID, pubKey, ipV4Addr, chainID string) (*string, error) {
	var input [][]byte
	if c.keyringBackend == keyring.BackendFile {
		input = c.formatInputs([]string{keyPassword, keyPassword, keyPassword})
	}

	return c.shellExecInput("sifnoded", input,
		"gentx", name, bondAmount,
		"--details", name,
		"--keyring-backend", c.keyringBackend,
		"--home", nodeDir,
		"--output-document", outputFile,
		"--node-id", nodeID,
		"--pubkey", pubKey,
		"--ip", ipV4Addr,
		"--chain-id", chainID,
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
	return fmt.Sprintf("%s/%s", c.configPath, ConfigTOML)
}

func (c CLI) AppFilePath() string {
	return fmt.Sprintf("%s/%s", c.configPath, AppTOML)
}

func (c CLI) TransferFunds(keyPassword, fromAddress, toAddress, coins string) (*string, error) {
	var input [][]byte
	if c.keyringBackend == keyring.BackendFile {
		input = c.formatInputs([]string{keyPassword, keyPassword})
	}

	return c.shellExecInput("sifnoded", input, "tx", "send", fromAddress, toAddress, coins, "-y")
}

func (c CLI) ValidatorPublicKeyAddress() (*string, error) {
	return c.shellExec("sifnoded", "tendermint", "show-validator")
}

func (c CLI) CreateValidator(moniker, validatorPublicKey, keyPassword, bondAmount string) (*string, error) {
	var input [][]byte
	if c.keyringBackend == keyring.BackendFile {
		input = c.formatInputs([]string{keyPassword, keyPassword})
	}

	return c.shellExecInput("sifnoded", input,
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
		"--keyring-backend", c.keyringBackend,
		"-y")
}

func (c CLI) formatInputs(inputs []string) [][]byte {
	formatted := make([][]byte, 0)
	for _, input := range inputs {
		formatted = append(formatted, []byte(input+"\n"))
	}

	return formatted
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
