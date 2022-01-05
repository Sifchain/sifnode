package cli_test

import (
	"os"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/margin/client/cli"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	testutilcli "github.com/cosmos/cosmos-sdk/testutil/cli"
	testnetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     testnetwork.Config
	network *testnetwork.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := testnetwork.DefaultConfig()
	cfg.NumValidators = 1
	encConfig := sifapp.MakeTestEncodingConfig()
	cfg.InterfaceRegistry = encConfig.InterfaceRegistry
	cfg.Codec = encConfig.Marshaler
	cfg.TxConfig = encConfig.TxConfig
	cfg.AppConstructor = func(val testnetwork.Validator) servertypes.Application {
		db := dbm.NewMemDB()
		return sifapp.NewSifApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, val.Dir, 0, encConfig, sifapp.EmptyAppOptions{})
	}

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

func (s *IntegrationTestSuite) TestOpenLongCmd() {
	val := s.network.Validators[0]

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls.
	// ...
	args := []string{
		"--collateral_asset=rowan",
		"--borrow_asset=atom",
		"--collateral_amount=1000",
		"--from=" + val.Address.String(),
		"-y",
	}
	_, err := testutilcli.ExecTestCLICmd(val.ClientCtx, cli.GetOpenLongCmd(), args)
	require.NoError(s.T(), err)
}

func (s *IntegrationTestSuite) TestCloseLongCmd() {
	val := s.network.Validators[0]

	args := []string{
		"--collateral_asset=rowan",
		"--borrow_asset=atom",
		"--from=" + val.Address.String(),
		"-y",
	}
	_, err := testutilcli.ExecTestCLICmd(val.ClientCtx, cli.GetCloseLongCmd(), args)
	require.NoError(s.T(), err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
