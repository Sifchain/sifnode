package sifgen

import (
	"testing"

	. "gopkg.in/check.v1"
)

type sifgenSuite struct{}

var _ = Suite(&sifgenSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *sifgenSuite) SetUpSuite(c *C) {}

func (s *sifgenSuite) TestSifgen(c *C)       {}
func (s *sifgenSuite) TestSummary(c *C)      {}
func (s *sifgenSuite) TestNetworkUtils(c *C) {}
func (s *sifgenSuite) TestNetworkNode(c *C)  {}
func (s *sifgenSuite) TestNetwork(c *C)      {}

func (s *sifgenSuite) TearDownSuite(c *C) {}
