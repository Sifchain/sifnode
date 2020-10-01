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

func (u Utils) InitChain(chainID, moniker string) {
	u.ShellExec(u.sifDaemon, "init", moniker, "--chain-id", chainID)
}

func (u Utils) SetKeyRingStorage() {
	u.ShellExec(u.sifCLI, "config", "keyring-backend", "file")
}

func (u Utils) SetConfigChainID(chainID string) {
	u.ShellExec(u.sifCLI, "config", "chain-id", chainID)
}

func (u Utils) SetConfigIndent(indent bool) {
	u.ShellExec(u.sifCLI, "config", "indent", fmt.Sprintf("%v", indent))
}

func (u Utils) SetConfigTrustNode(indent bool) {
	u.ShellExec(u.sifCLI, "config", "trust-node", fmt.Sprintf("%v", indent))
}

func (u Utils) AddKey(name, keyPassword string) string {
	return u.ShellExecInput(u.sifCLI, [][]byte{[]byte(keyPassword + "\n"), []byte(keyPassword + "\n")}, "keys", "add", name)
}

func (u Utils) AddGenesisAccount(address string, coins []string) {
	u.ShellExec(u.sifDaemon, "add-genesis-account", address, strings.Join(coins[:], ","))
}

func (u Utils) GenerateGenesisTxn(name, keyPassword string) {
	u.ShellExecInput(u.sifDaemon,
		[][]byte{[]byte(keyPassword + "\n"), []byte(keyPassword + "\n"), []byte(keyPassword + "\n")},
		"gentx", "--name", name, "--keyring-backend", "file",
	)
}

func (u Utils) CollectGenesisTxns() {
	u.ShellExec(u.sifDaemon, "collect-gentxs")
}

func (u Utils) ExportGenesis() string {
	return u.ShellExec(u.sifDaemon, "export")
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

func (u Utils) SaveGenesis(genesis types.Genesis) {
	err := ioutil.WriteFile(u.GenesisFilePath(), *genesis.Result.Genesis, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func (u Utils) ReplacePeerConfig(peerAddresses []string) {
	contents, err := ioutil.ReadFile(u.ConfigFilePath())
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

func (u Utils) ShellExec(cmd string, args ...string) string {
	c := exec.Command(cmd, args...)
	var out bytes.Buffer
	c.Stdout = &out

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}

	return out.String()
}

func (u Utils) ShellExecInput(cmd string, inputs [][]byte, args ...string) string {
	c := exec.Command(cmd, args...)
	stdin, err := c.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := c.StdoutPipe()
	if err != nil {
		panic(err)
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
		panic(err)
	}

	return buf.String()
}
