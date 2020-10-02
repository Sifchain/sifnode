package networks

import (
	"fmt"
	"testing"

	"github.com/Sifchain/sifnode/tools/sifgen/networks/types"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/stretchr/testify/assert"
	. "gopkg.in/check.v1"
	"syreclabs.com/go/faker"
)

type networksSuite struct{}

const (
	defaultNodeHome = "/tmp"
	defaultCLIHome  = "/tmp"
	genesisFile     = "genesis.json"
	configFile      = "config.toml"
)

var (
	_           = Suite(&networksSuite{})
	nodeAddress = fmt.Sprintf("sif%s", faker.RandomString(39))
	peerAddress = fmt.Sprintf("%s@%s:%v",
		faker.RandomString(40),
		faker.Internet().IpV4Address(),
		faker.Number().NumberInt(5),
	)
	genesisURL = fmt.Sprintf("%s/%s", faker.Internet().Url(), faker.Internet().Slug())
)

type mockNetworkUtils struct{}

func (u mockNetworkUtils) InitChain(chainID, moniker string) (*string, error) { return nil, nil }

func (u mockNetworkUtils) Reset(paths []string) error {
	return nil
}

func (u mockNetworkUtils) SetKeyRingStorage() (*string, error)              { return nil, nil }
func (u mockNetworkUtils) SetConfigChainID(chainID string) (*string, error) { return nil, nil }
func (u mockNetworkUtils) SetConfigIndent(indent bool) (*string, error)     { return nil, nil }
func (u mockNetworkUtils) SetConfigTrustNode(trust bool) (*string, error)   { return nil, nil }

func (u mockNetworkUtils) AddKey(name, keyPassword string) (*string, error) {
	key := heredoc.Doc(`
- name: foobar
  type: local
  address: ` + nodeAddress + `
  pubkey: sifpub1addwnpepqgxzxu0ftntzarv4equjj5g7df3ftjwyy7ylymm9083xwmrtlzen2wn5k7g
  mnemonic: ""
  threshold: 0
  pubkeys: []
`)
	return &key, nil
}

func (u mockNetworkUtils) AddGenesisAccount(name string, coins []string) (*string, error) {
	return nil, nil
}
func (u mockNetworkUtils) GenerateGenesisTxn(name, keyPassword string) (*string, error) {
	return nil, nil
}
func (u mockNetworkUtils) CollectGenesisTxns() (*string, error) { return nil, nil }

func (u mockNetworkUtils) ExportGenesis() (*string, error) {
	genesis := heredoc.Doc(`
{
  "genesis_time": "2020-10-01T20:29:57.428357Z",
  "chain_id": "sifchain",
  "app_hash": "",
  "app_state": {
    "params": null,
    "supply": {
      "supply": []
    },
    "sifnode": {},
    "genutil": {
      "gentxs": [
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "memo": "` + peerAddress + `"
          }
        }
      ]
    }
  }
}
`)
	return &genesis, nil
}

func (u mockNetworkUtils) GenesisFilePath() string {
	return fmt.Sprintf("%s/%s", defaultNodeHome, genesisFile)
}
func (u mockNetworkUtils) ConfigFilePath() string {
	return fmt.Sprintf("%s/%s", defaultNodeHome, configFile)
}

func (u mockNetworkUtils) ScrapePeerGenesis(url string) (types.Genesis, error) {
	return types.Genesis{}, nil
}

func (u mockNetworkUtils) SaveGenesis(genesis types.Genesis) error        { return nil }
func (u mockNetworkUtils) ReplacePeerConfig(peerAddresses []string) error { return nil }

func Test(t *testing.T) { TestingT(t) }

func (s *networksSuite) SetUpSuite(c *C) {}

func (s *networksSuite) TestLocalnetValidator(c *C) {
	node := NewValidator(mockNetworkUtils{})
	network := NewLocalnet(defaultNodeHome, defaultCLIHome, faker.Lorem().Word(), node, mockNetworkUtils{})

	err := network.Setup()
	c.Assert(err, IsNil)

	err = network.Genesis()
	c.Assert(err, IsNil)

	c.Assert(node.Name(), NotNil)
	c.Assert(*node.Address(nil), Equals, nodeAddress)
	c.Assert(node.PeerAddress(), Equals, peerAddress)
	c.Assert(node.KeyPassword(), NotNil)

	node.peerAddress = ""
	c.Assert(assert.Empty(c, node.PeerAddress()), Equals, true)

	err = node.CollectPeerAddress()
	c.Assert(err, IsNil)
	c.Assert(node.PeerAddress(), Equals, peerAddress)
}

func (s *networksSuite) TestLocalnetWitness(c *C) {
	node := NewWitness(peerAddress, genesisURL, mockNetworkUtils{})
	network := NewLocalnet(defaultNodeHome, defaultCLIHome, faker.Lorem().Word(), node, mockNetworkUtils{})

	err := network.Setup()
	c.Assert(err, IsNil)

	err = network.Genesis()
	c.Assert(err, IsNil)

	c.Assert(node.Name(), NotNil)
	c.Assert(*node.Address(nil), Equals, nodeAddress)
	c.Assert(node.KeyPassword(), NotNil)
	c.Assert(node.GenesisURL(), Equals, genesisURL)
}

func (s *networksSuite) TearDownSuite(c *C) {}
