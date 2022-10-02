//go:build TEST_INTEGRATION
// +build TEST_INTEGRATION

package cli_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/margin/client/cli"
	testutilcli "github.com/cosmos/cosmos-sdk/testutil/cli"
	testnetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     testnetwork.Config
	network *testnetwork.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	sifapp.SetConfig(false)

	cfg := testnetwork.DefaultConfig()
	cfg.NumValidators = 1
	encConfig := sifapp.MakeTestEncodingConfig()
	cfg.InterfaceRegistry = encConfig.InterfaceRegistry
	cfg.Codec = encConfig.Marshaler
	cfg.TxConfig = encConfig.TxConfig
	// cfg.AppConstructor = func(val testnetwork.Validator) servertypes.Application {
	// 	db := dbm.NewMemDB()
	// 	return sifapp.NewSifApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, val.Dir, 0, encConfig, sifapp.EmptyAppOptions{})
	// }

	s.cfg = cfg
	s.network = testnetwork.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestOpenCmd() {
	val := s.network.Validators[0]

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...
	args := []string{
		"--collateral_asset=rowan",
		"--borrow_asset=atom",
		"--collateral_amount=1000",
		"--position=long",
		"--from=" + val.Address.String(),
		"-y",
	}
	_, err := testutilcli.ExecTestCLICmd(val.ClientCtx, cli.GetOpenCmd(), args)
	require.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestCloseCmd() {
	val := s.network.Validators[0]

	args := []string{
		"--id=1",
		"--from=" + val.Address.String(),
		"-y",
	}
	_, err := testutilcli.ExecTestCLICmd(val.ClientCtx, cli.GetCloseCmd(), args)
	require.NoError(s.T(), err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
