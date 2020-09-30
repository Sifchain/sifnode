package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Sifchain/sifnode/app"

	"github.com/sethvargo/go-password/password"
	"github.com/yelinaung/go-haikunator"
	"gopkg.in/yaml.v2"
)

const (
	daemon = "daemon"
	cli    = "cli"
)

// Because the binary names differ between local and what we end up with in the containers.
var (
	defaultNetwork  = os.Getenv("NETWORK")
	defaultChainID  = os.Getenv("CHAIN_ID")
	defaultNodeType = os.Getenv("NODE_TYPE")
	defaultCoins    = []string{"1000rowan", "100000000stake"}

	executables = map[string][]string{
		daemon: {"sifnoded", "sifd"},
		cli:    {"sifnodecli", "sifcli"},
	}
)

func main() {
	network := flag.String("n", defaultNetwork, "The network [localnet|testnet|mainnet].")
	chainID := flag.String("c", defaultChainID, "The ID of the chain.")
	nodeType := flag.String("t", defaultNodeType, "The node type [validator|witness].")
	flag.Parse()

	genesis := NewGenesis(*network, *chainID, *nodeType)
	genesis.build()
}

type Keys []Key
type Key struct {
	Name    string `json:"name" yaml:"name"`
	KeyType string `json:"type" yaml:"type"`
	Address string `json:"address" yaml:"address"`
	PubKey  string `json:"pubkey" yaml:"pubkey"`
}

type Genesis struct {
	network  string
	chainID  string
	nodeType string
	moniker  haikunator.Name
}

func NewGenesis(network, chainID, nodeType string) Genesis {
	return Genesis{
		network:  network,
		chainID:  chainID,
		nodeType: nodeType,
		moniker:  haikunator.New(time.Now().UTC().UnixNano()),
	}
}

func (g Genesis) build() {
	switch g.network {
	case "localnet":
		g.reset()
		g.localnet()
	}
}

func (g Genesis) localnet() {
	keyName := g.moniker.Haikunate()
	keyPwd := g.generatePassword()

	g.initChain()
	g.setKeyringStorage()
	keys := g.addKey(keyName, keyPwd)
	g.addGenesisAccount(keys[0].Address, strings.Join(defaultCoins[:], ","))
	g.setConfig("json", true, true)
	g.generateGenesisTx(keyName, keyPwd)
	g.collectGenesisTxns()

	fmt.Printf("%s initialized.\n\nValidator Details\n-----------------\nName: %s\nAddress: %s\nPassword: %s\n",
		g.network, keyName, keys[0].Address, keyPwd)
}

// Look for the binaries. These differ between local and k8s environments.
func (g Genesis) executable(executableType string) *string {
	if len(executables[executableType]) == 0 {
		panic(fmt.Errorf("unknown type %s\n", executableType))
	}

	for _, exe := range executables[executableType] {
		path, err := exec.LookPath(exe)
		if err == nil {
			return &path
		}
	}

	return nil
}

// Initializes a new chain.
func (g Genesis) initChain() {
	g.cmd(*g.executable(daemon), "init", g.moniker.Haikunate(), "--chain-id", g.chainID)
}

// Sets the key ring storage.
func (g Genesis) setKeyringStorage() {
	g.cmd(*g.executable(cli), "config", "keyring-backend", "file")
}

// Add a new validator key.
func (g Genesis) addKey(name, pwd string) Keys {
	r := g.cmdWithInput(
		*g.executable(cli),
		[][]byte{[]byte(pwd + "\n"), []byte(pwd + "\n")},
		"keys", "add", name,
	)

	yml, err := ioutil.ReadAll(strings.NewReader(r))
	if err != nil {
		panic(err)
	}

	var keys Keys
	err = yaml.Unmarshal(yml, &keys)
	if err != nil {
		panic(err)
	}

	return keys
}

// Generate a new, random password for a key.
func (g Genesis) generatePassword() string {
	pwd, err := password.Generate(32, 5, 0, false, false)
	if err != nil {
		panic(err)
	}

	return pwd
}

// Add genesis account.
func (g Genesis) addGenesisAccount(address, coins string) {
	g.cmd(*g.executable(daemon), "add-genesis-account", address, coins)
}

// Set config.
func (g Genesis) setConfig(output string, indent, trust bool) {
	g.setConfigChainID()
	g.setConfigOutput(output)
	g.setConfigIndent(indent)
	g.setConfigTrustNode(trust)
}

// Set chain-id.
func (g Genesis) setConfigChainID() {
	g.cmd(*g.executable(cli), "config", "chain-id", g.chainID)
}

// Set the output type.
func (g Genesis) setConfigOutput(output string) {
	g.cmd(*g.executable(cli), "config", "output", output)
}

// Set indenting.
func (g Genesis) setConfigIndent(indent bool) {
	g.cmd(*g.executable(cli), "config", "indent", fmt.Sprintf("%v", indent))
}

// Trust the node?
func (g Genesis) setConfigTrustNode(trust bool) {
	g.cmd(*g.executable(cli), "config", "trust-node", fmt.Sprintf("%v", trust))
}

// Generate the genesis transaction.
func (g Genesis) generateGenesisTx(name, pwd string) {
	g.cmdWithInput(
		*g.executable(daemon),
		[][]byte{[]byte(pwd + "\n"), []byte(pwd + "\n"), []byte(pwd + "\n")},
		"gentx", "--name", name, "--keyring-backend", "file",
	)
}

// Collect the genesis transactions.
func (g Genesis) collectGenesisTxns() {
	g.cmd(*g.executable(daemon), "collect-gentxs")
}

// Removes any existing config in $HOME.
func (g Genesis) reset() {
	roots := []string{app.DefaultCLIHome, app.DefaultNodeHome}
	for _, root := range roots {
		if _, err := os.Stat(root); !os.IsNotExist(err) {
			err = os.RemoveAll(root)
			if err != nil {
				panic(err)
			}
		}
	}
}

//  Wrapper for exec.Command.
func (g Genesis) cmd(name string, args ...string) string {
	c := exec.Command(name, args...)
	var out bytes.Buffer
	c.Stdout = &out

	err := c.Run()
	if err != nil {
		panic(err)
	}

	return out.String()
}

// Wrapper for exec.Command, with inputs.
func (g Genesis) cmdWithInput(name string, inputs [][]byte, args ...string) string {
	c := exec.Command(name, args...)
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
