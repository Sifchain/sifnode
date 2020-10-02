package networks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/Sifchain/sifnode/tools/sifgen/networks/types"
)

const (
	GenesisFile = "genesis.json"
	ConfigFile  = "config.toml"
)

type Utils struct {
	sifDaemon  string
	sifCLI     string
	configPath string
}

func NewUtils(defaultNodeHome string) Utils {
	return Utils{
		sifDaemon:  os.Getenv("SIF_DAEMON"),
		sifCLI:     os.Getenv("SIF_CLI"),
		configPath: fmt.Sprintf("%s/config", defaultNodeHome),
	}
}

func (u Utils) Reset(paths []string) error {
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

func (u Utils) InitChain(chainID, moniker string) (*string, error) {
	return u.shellExec(u.sifDaemon, "init", moniker, "--chain-id", chainID)
}

func (u Utils) SetKeyRingStorage() (*string, error) {
	return u.shellExec(u.sifCLI, "config", "keyring-backend", "file")
}

func (u Utils) SetConfigChainID(chainID string) (*string, error) {
	return u.shellExec(u.sifCLI, "config", "chain-id", chainID)
}

func (u Utils) SetConfigIndent(indent bool) (*string, error) {
	return u.shellExec(u.sifCLI, "config", "indent", fmt.Sprintf("%v", indent))
}

func (u Utils) SetConfigTrustNode(indent bool) (*string, error) {
	return u.shellExec(u.sifCLI, "config", "trust-node", fmt.Sprintf("%v", indent))
}

func (u Utils) AddKey(name, keyPassword string) (*string, error) {
	return u.shellExecInput(u.sifCLI, [][]byte{[]byte(keyPassword + "\n"), []byte(keyPassword + "\n")}, "keys", "add", name)
}

func (u Utils) AddGenesisAccount(address string, coins []string) (*string, error) {
	return u.shellExec(u.sifDaemon, "add-genesis-account", address, strings.Join(coins[:], ","))
}

func (u Utils) GenerateGenesisTxn(name, keyPassword string) (*string, error) {
	return u.shellExecInput(u.sifDaemon,
		[][]byte{[]byte(keyPassword + "\n"), []byte(keyPassword + "\n"), []byte(keyPassword + "\n")},
		"gentx", "--name", name, "--keyring-backend", "file",
	)
}

func (u Utils) CollectGenesisTxns() (*string, error) {
	return u.shellExec(u.sifDaemon, "collect-gentxs")
}

func (u Utils) ExportGenesis() (*string, error) {
	return u.shellExec(u.sifDaemon, "export")
}

func (u Utils) GenesisFilePath() string {
	return fmt.Sprintf("%s/%s", u.configPath, GenesisFile)
}

func (u Utils) ConfigFilePath() string {
	return fmt.Sprintf("%s/%s", u.configPath, ConfigFile)
}

func (u Utils) ScrapePeerGenesis(url string) types.Genesis {
	response, err := http.Get(fmt.Sprintf("%s", url))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var genesis types.Genesis
	if err := json.Unmarshal(body, &genesis); err != nil {
		log.Fatal(err)
	}

	return genesis
}

func (u Utils) SaveGenesis(genesis types.Genesis) error {
	err := ioutil.WriteFile(u.GenesisFilePath(), *genesis.Result.Genesis, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (u Utils) ReplacePeerConfig(peerAddresses []string) error {
	contents, err := ioutil.ReadFile(u.ConfigFilePath())
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

	err = ioutil.WriteFile(u.ConfigFilePath(), []byte(output), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (u Utils) shellExec(cmd string, args ...string) (*string, error) {
	c := exec.Command(cmd, args...)
	var out bytes.Buffer
	c.Stdout = &out

	err := c.Run()
	if err != nil {
		return nil, err
	}

	result := out.String()
	return &result, nil
}

func (u Utils) shellExecInput(cmd string, inputs [][]byte, args ...string) (*string, error) {
	c := exec.Command(cmd, args...)
	stdin, err := c.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	go func() {
		io.Copy(buf, stdout)
	}()

	if err := c.Start(); err != nil {
		panic(err)
	}

	for _, i := range inputs {
		stdin.Write(i)
	}

	if err := c.Wait(); err != nil {
		return nil, err
	}

	result := buf.String()
	return &result, nil
}
