package network

import (
	"testing"

	. "gopkg.in/check.v1"
)

type networkSuite struct{}

var _ = Suite(&networkSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *networkSuite) SetUpSuite(c *C)              {}
func (s *networkSuite) SetUpTest(c *C)               {}
func (s *networkSuite) TestInitNodes(c *C)           {}
func (s *networkSuite) TestCreateDirs(c *C)          {}
func (s *networkSuite) TestSetDefaultConfig(c *C)    {}
func (s *networkSuite) TestGenerateKey(c *C)         {}
func (s *networkSuite) TestInitChain(c *C)           {}
func (s *networkSuite) TestSetValidatorAddress(c *C) {}
func (s *networkSuite) TestSetNodeID(c *C)           {}
func (s *networkSuite) TestGetSeedNode(c *C)         {}
func (s *networkSuite) TestAddGenesis(c *C)          {}
func (s *networkSuite) TestGenerateTx(c *C)          {}
func (s *networkSuite) TestCollectGenTxs(c *C)       {}
func (s *networkSuite) TestBuild(c *C)               {}
func (s *networkSuite) TearDownTest(c *C)            {}
func (s *networkSuite) TearDownSuite(c *C)           {}
