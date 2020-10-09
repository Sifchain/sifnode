package node

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Sifchain/sifnode/tools/sifgen/faucet"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/yelinaung/go-haikunator"
	. "gopkg.in/check.v1"
	"syreclabs.com/go/faker"
)

type nodeSuite struct {
	node   *Node
	tmpDir string
	server *httptest.Server
}

var (
	_               = Suite(&nodeSuite{})
	nodeKeyAddress  = fmt.Sprintf("sif%s", faker.RandomString(39))
	nodePeerAddress = fmt.Sprintf("%s@%s:%v",
		faker.RandomString(40),
		faker.Internet().IpV4Address(),
		faker.Number().NumberInt(5),
	)
	chainID                       = haikunator.New(time.Now().UTC().UnixNano()).Haikunate()
	moniker                       = "moniker"
	nodeSeedAddress               = faker.Internet().IpV4Address()
	genesisURL                    *string
	nodeValidatorPublicKeyAddress = fmt.Sprintf("sifvalconspub1zcjduepq%s", faker.RandomString(58))
	nodeConfigFixtureFile         = "config.toml"
	nodeConfigFixturePath         = fmt.Sprintf("%s/%s", "../../../test/unit/fixtures", nodeConfigFixtureFile)
	nodeGenesisFixtureFile        = "genesis.json"
	nodeGenesisFixturePath        = fmt.Sprintf("%s/%s", "../../../test/unit/fixtures", nodeGenesisFixtureFile)
	tmpDir                        *string
	nodePeerList                  = []string{
		fmt.Sprintf("%s@%s:%v",
			faker.RandomString(40),
			faker.Internet().IpV4Address(),
			faker.Number().NumberInt(5)),
		fmt.Sprintf("%s@%s:%v",
			faker.RandomString(40),
			faker.Internet().IpV4Address(),
			faker.Number().NumberInt(5)),
	}
)

type mockCLIUtils struct{}

func (c mockCLIUtils) Reset() error                                       { return nil }
func (c mockCLIUtils) CurrentChainID() (*string, error)                   { return &chainID, nil }
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
	return fmt.Sprintf("%s/%s", *tmpDir, nodeGenesisFixtureFile)
}

func (c mockCLIUtils) ConfigFilePath() string {
	return fmt.Sprintf("%s/%s", *tmpDir, "config.toml")
}

func (c mockCLIUtils) TransferFunds(keyPassword, fromAddress, toAddress, coins string) (*string, error) {
	return nil, nil
}

func (c mockCLIUtils) ValidatorPublicKeyAddress() (*string, error) {
	return &nodeValidatorPublicKeyAddress, nil
}

func (c mockCLIUtils) CreateValidator(string, string, string, string) (*string, error) {
	return nil, nil
}

func Test(t *testing.T) { TestingT(t) }

func (s *nodeSuite) SetUpSuite(c *C) {
	dir, err := ioutil.TempDir("/tmp", faker.RandomString(32))
	c.Assert(err, IsNil)

	tmpDir = &dir
	s.tmpDir = *tmpDir

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile(nodeGenesisFixturePath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write(content)
	})

	s.server = httptest.NewServer(h)
	genesisURL = &s.server.URL
}

func (s *nodeSuite) SetUpTest(c *C) {
	err := copyFixture(nodeConfigFixturePath, fmt.Sprintf("%s/%s", s.tmpDir, nodeConfigFixtureFile))
	c.Assert(err, IsNil)

	s.node = &Node{
		chainID:     chainID,
		moniker:     moniker,
		seedAddress: &nodeSeedAddress,
		genesisURL:  genesisURL,
		CLI:         mockCLIUtils{},
	}
}

func (s *nodeSuite) TestValidate(c *C) {
	err := s.node.Validate()
	c.Assert(err, IsNil)

	s.node.chainID = faker.Lorem().Word()
	err = s.node.Validate()
	c.Assert(err, NotNil)

	s.node.moniker = faker.Lorem().Word()
	err = s.node.Validate()
	c.Assert(err, NotNil)

	s.node.chainID = chainID
	s.node.moniker = moniker
	err = s.node.Validate()
	c.Assert(err, IsNil)
}

func (s *nodeSuite) TestSetup(c *C) {
	err := s.node.Setup()
	c.Assert(err, IsNil)
}

func (s *nodeSuite) TestGenesis(c *C) {
	err := s.node.Genesis(faucet.NewFaucet(chainID).DefaultDeposit())
	c.Assert(err, IsNil)
}

func (s *nodeSuite) TestChainID(c *C) {
	c.Assert(s.node.ChainID(), Equals, chainID)
}

func (s *nodeSuite) TestMoniker(c *C) {
	c.Assert(s.node.Moniker(), Equals, moniker)
}

func (s *nodeSuite) TestNodeKeyAddress(c *C) {
	key := s.node.NodeKeyAddress(&nodeKeyAddress)
	c.Assert(key, IsNil)

	_nodeKeyAddress := s.node.NodeKeyAddress(nil)
	c.Assert(*_nodeKeyAddress, Equals, nodeKeyAddress)
}

func (s *nodeSuite) TestNodeKeyPassword(c *C) {
	c.Assert(s.node.NodeKeyPassword(), NotNil)
}

func (s *nodeSuite) TestNodePeerAddress(c *C) {
	c.Assert(s.node.NodePeerAddress(), NotNil)
	c.Assert(s.node.collectNodePeerAddress(), IsNil)
	c.Assert(s.node.NodePeerAddress(), Equals, nodePeerAddress)
}

func (s *nodeSuite) TestNodeValidatorPublicKeyAddress(c *C) {
	c.Assert(s.node.collectNodeValidatorPublicKeyAddress(), IsNil)
	c.Assert(s.node.NodeValidatorPublicKeyAddress(), Equals, nodeValidatorPublicKeyAddress)
}

func (s *nodeSuite) TestSeedAddress(c *C) {
	_nodeSeedAddress := s.node.SeedAddress()
	c.Assert(*_nodeSeedAddress, Equals, nodeSeedAddress)
}

func (s *nodeSuite) TestUpdatePeerList(c *C) {
	err := s.node.UpdatePeerList(nodePeerList)
	c.Assert(err, IsNil)

	config, err := s.node.parseConfig()
	c.Assert(err, IsNil)
	c.Assert(config.P2P.PersistentPeers, Equals, strings.Join(nodePeerList, ","))
}

func (s *nodeSuite) TearDownTest(c *C) {}

func (s *nodeSuite) TearDownSuite(c *C) {
	err := os.RemoveAll(s.tmpDir)
	c.Assert(err, IsNil)

	s.server.Close()
}

func copyFixture(sourcePath, destPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
