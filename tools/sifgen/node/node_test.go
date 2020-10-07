package node

import (
	"fmt"
	"testing"

	"github.com/Sifchain/sifnode/tools/sifgen/faucet"
	"github.com/Sifchain/sifnode/tools/sifgen/node/types"

	"github.com/MakeNowJust/heredoc/v2"
	. "gopkg.in/check.v1"
	"syreclabs.com/go/faker"
)

type nodeSuite struct{}

var (
	_               = Suite(&nodeSuite{})
	nodeKeyAddress  = fmt.Sprintf("sif%s", faker.RandomString(39))
	nodePeerAddress = fmt.Sprintf("%s@%s:%v",
		faker.RandomString(40),
		faker.Internet().IpV4Address(),
		faker.Number().NumberInt(5),
	)
	nodeSeedAddress = faker.Internet().IpV4Address()
	genesisURL      = fmt.Sprintf("%s/%s", faker.Internet().Url(), faker.Internet().Slug())
	nodeValidatorPublicKeyAddress = fmt.Sprintf("sifvalconspub1zcjduepq%s", faker.RandomString(58))
)

type mockCLIUtils struct{}

func (c mockCLIUtils) Reset() error                                       { return nil }
func (c mockCLIUtils) InitChain(chainID, moniker string) (*string, error) { return nil, nil }
func (c mockCLIUtils) SetKeyRingStorage() (*string, error)                { return nil, nil }
func (c mockCLIUtils) SetConfigChainID(chainID string) (*string, error)   { return nil, nil }
func (c mockCLIUtils) SetConfigIndent(indent bool) (*string, error)       { return nil, nil }
func (c mockCLIUtils) SetConfigTrustNode(trust bool) (*string, error)     { return nil, nil }

func (c mockCLIUtils) AddKey(name, keyPassword string) (*string, error) {
	key := heredoc.Doc(`
- name: foobar
  type: local
  address: ` + nodeKeyAddress + `
  pubkey: sifpub1addwnpepqgxzxu0ftntzarv4equjj5g7df3ftjwyy7ylymm9083xwmrtlzen2wn5k7g
  mnemonic: ""
  threshold: 0
  pubkeys: []
`)
	return &key, nil
}

func (c mockCLIUtils) AddGenesisAccount(name string, coins []string) (*string, error) {
	return nil, nil
}
func (c mockCLIUtils) GenerateGenesisTxn(name, keyPassword string) (*string, error) {
	return nil, nil
}
func (c mockCLIUtils) CollectGenesisTxns() (*string, error) { return nil, nil }

func (c mockCLIUtils) ExportGenesis() (*string, error) {
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
            "memo": "` + nodePeerAddress + `"
          }
        }
      ]
    }
  }
}
`)
	return &genesis, nil
}

func (c mockCLIUtils) GenesisFilePath() string {
	return fmt.Sprintf("%s/%s", faker.Internet().Slug(), fmt.Sprintf("%s.json", faker.Internet().Slug()))
}
func (c mockCLIUtils) ConfigFilePath() string {
	return fmt.Sprintf("%s/%s", faker.Internet().Slug(), fmt.Sprintf("%s.toml", faker.Internet().Slug()))
}

func (c mockCLIUtils) ScrapePeerGenesis(url string) (types.Genesis, error) {
	return types.Genesis{}, nil
}

func (c mockCLIUtils) SaveGenesis(genesis types.Genesis) error        { return nil }
func (c mockCLIUtils) ReplacePeerConfig(peerAddresses []string) error { return nil }
func (c mockCLIUtils) TransferFunds(keyPassword, fromAddress, toAddress, coins string) (*string, error) {
	return nil, nil
}
func (c mockCLIUtils) ValidatorPublicKeyAddress() (*string, error) { return &nodeValidatorPublicKeyAddress, nil }
func (c mockCLIUtils) CreateValidator(string, string, string, string) (*string, error) {
	return nil, nil
}

func Test(t *testing.T) { TestingT(t) }

func (s *nodeSuite) SetUpSuite(c *C) {}

func (s *nodeSuite) TestNode(c *C) {
	n := &Node{
		chainID:     faker.Lorem().Word(),
		moniker:     faker.Internet().Slug(),
		seedAddress: &nodeSeedAddress,
		genesisURL:  &genesisURL,
		CLI:         mockCLIUtils{},
	}

	err := n.Setup()
	c.Assert(err, IsNil)

	err = n.Genesis(faucet.NewFaucet(n.chainID).DefaultDeposit())
	c.Assert(err, IsNil)

	key := n.NodeKeyAddress(&nodeKeyAddress)
	c.Assert(key, IsNil)

	_nodeKeyAddress := n.NodeKeyAddress(nil)
	c.Assert(*_nodeKeyAddress, Equals, nodeKeyAddress)

	c.Assert(n.NodeKeyPassword(), NotNil)
	c.Assert(n.NodePeerAddress(), NotNil)

	c.Assert(n.collectNodePeerAddress(), IsNil)
	c.Assert(n.NodePeerAddress(), Equals, nodePeerAddress)

	c.Assert(n.collectNodeValidatorPublicKeyAddress(), IsNil)
	c.Assert(n.NodeValidatorPublicKeyAddress(), Equals, nodeValidatorPublicKeyAddress)

	_nodeSeedAddress := n.SeedAddress()
	c.Assert(*_nodeSeedAddress, Equals, nodeSeedAddress)

	c.Assert(n.generateNodeKeyAddress(), IsNil)
	c.Assert(n.generateNodeKeyPassword(), IsNil)
	c.Assert(n.seedGenesis(faucet.NewFaucet(n.chainID).DefaultDeposit()), IsNil)
	c.Assert(n.validatorGenesis(), IsNil)
}

func (s *nodeSuite) TearDownSuite(c *C) {}
