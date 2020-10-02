package sifgen

import (
	"fmt"

	"testing"

	"github.com/Sifchain/sifnode/tools/sifgen/networks"

	"github.com/stretchr/testify/assert"
	. "gopkg.in/check.v1"
	"syreclabs.com/go/faker"
)

type sifgenSuite struct{}

var (
	_ = Suite(&sifgenSuite{})

	chainID     = faker.Internet().Slug()
	peerAddress = fmt.Sprintf("%s@%s:%v",
		faker.RandomString(40),
		faker.Internet().IpV4Address(),
		faker.Number().NumberInt(5),
	)
	genesisURL = fmt.Sprintf("%s/%s", faker.Internet().Url(), faker.Internet().Slug())
)

func Test(t *testing.T) { TestingT(t) }

func (s *sifgenSuite) SetUpSuite(c *C) {}

func (s *sifgenSuite) TestSifgen(c *C) {
	sifValidator := NewSifgen(validator, localnet, chainID, nil, nil)
	c.Assert(sifValidator.network, Equals, localnet)
	c.Assert(sifValidator.nodeType, Equals, validator)
	c.Assert(sifValidator.chainID, Equals, chainID)

	sfWitness := NewSifgen(witness, localnet, chainID, &peerAddress, &genesisURL)
	c.Assert(sfWitness.network, Equals, localnet)
	c.Assert(sfWitness.nodeType, Equals, witness)
	c.Assert(sfWitness.chainID, Equals, chainID)
}

func (s *sifgenSuite) TestNetworkUtils(c *C) {
	// Ensure we get back a valid network utils struct.
	utils := NetworkUtils()
	_, ok := utils.(networks.NetworkUtils)
	c.Assert(ok, Equals, true)
}

func (s *sifgenSuite) TestNetworkNode(c *C) {
	// Test network node instantiation for a validator.
	sifValidator := NewSifgen(validator, localnet, chainID, nil, nil)

	utils := NetworkUtils()
	_, ok := utils.(networks.NetworkUtils)
	c.Assert(ok, Equals, true)

	node, err := NewNetworkNode(sifValidator, utils)
	c.Assert(err, IsNil)

	_, ok = (*node).(networks.NetworkNode)
	c.Assert(ok, Equals, true)
	c.Assert(assert.NotEmpty(c, (*node).Name()), Equals, true)

	// Test network node instantiation for a witness.
	sifWitness := NewSifgen(witness, localnet, chainID, &peerAddress, &genesisURL)

	utils = NetworkUtils()
	_, ok = utils.(networks.NetworkUtils)
	c.Assert(ok, Equals, true)

	node, err = NewNetworkNode(sifWitness, utils)
	c.Assert(err, IsNil)

	_, ok = (*node).(networks.NetworkNode)
	c.Assert(ok, Equals, true)
	c.Assert(assert.NotEmpty(c, (*node).Name()), Equals, true)

	sifRandom := NewSifgen(faker.Internet().Slug(), faker.Internet().Slug(), chainID, nil, nil)

	utils = NetworkUtils()
	_, ok = utils.(networks.NetworkUtils)
	c.Assert(ok, Equals, true)

	_, err = NewNetworkNode(sifRandom, utils)
	c.Assert(err, NotNil)
}

func (s *sifgenSuite) TestNetwork(c *C) {
	// Test network instantiation, as a validator.
	sifValidator := NewSifgen(validator, localnet, chainID, nil, nil)

	utils := NetworkUtils()
	_, ok := utils.(networks.NetworkUtils)
	c.Assert(ok, Equals, true)

	node, err := NewNetworkNode(sifValidator, utils)
	c.Assert(err, IsNil)

	_, ok = (*node).(networks.NetworkNode)
	c.Assert(ok, Equals, true)
	c.Assert(assert.NotEmpty(c, (*node).Name()), Equals, true)

	network, err := NewNetwork(sifValidator, utils, *node)
	c.Assert(err, IsNil)

	_, ok = (*network).(networks.Network)
	c.Assert(ok, Equals, true)

	// Test network instantiation, as a witness.
	sifWitness := NewSifgen(witness, localnet, chainID, &peerAddress, &genesisURL)

	utils = NetworkUtils()
	_, ok = utils.(networks.NetworkUtils)
	c.Assert(ok, Equals, true)

	node, err = NewNetworkNode(sifWitness, utils)
	c.Assert(err, IsNil)

	_, ok = (*node).(networks.NetworkNode)
	c.Assert(ok, Equals, true)
	c.Assert(assert.NotEmpty(c, (*node).Name()), Equals, true)

	network, err = NewNetwork(sifWitness, utils, *node)
	c.Assert(err, IsNil)

	_, ok = (*network).(networks.Network)
	c.Assert(ok, Equals, true)
}

func (s *sifgenSuite) TearDownSuite(c *C) {}
