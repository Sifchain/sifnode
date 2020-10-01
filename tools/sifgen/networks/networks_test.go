package networks

import (
	"testing"

	. "gopkg.in/check.v1"
)

type networksSuite struct{}

var _ = Suite(&networksSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *networksSuite) SetUpSuite(c *C)    {}
func (s *networksSuite) TearDownSuite(c *C) {}
