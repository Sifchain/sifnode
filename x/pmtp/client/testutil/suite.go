//go:build integration
// +build integration

package testutil

import (
	sifapp "github.com/Sifchain/sifnode/app"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	mnemonic       string
	address        string
	nativeAmount   types.Int
	externalAmount types.Int
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	sifapp.SetConfig(false)

	s.mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
	s.address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"
	s.nativeAmount, _ = types.NewIntFromString("999999000000000000000000000")
	s.externalAmount, _ = types.NewIntFromString("500000000000000000000000")

	s.cfg.Mnemonics = []string{s.mnemonic}
	s.cfg.StakingTokens = s.nativeAmount

	bz, err := GetBankGenesisState(s.cfg, s.address, s.nativeAmount, s.externalAmount)
	s.Require().NoError(err)
	s.cfg.GenesisState["bank"] = bz

	bz, err = GetTokenRegistryGenesisState(s.cfg, s.address)
	s.Require().NoError(err)
	s.cfg.GenesisState["tokenregistry"] = bz

	bz, err = GetClpGenesisState(s.cfg, types.NewUint(3000000000000000000), types.NewUint(2000000000000000000))
	s.Require().NoError(err)
	s.cfg.GenesisState["clp"] = bz

	// bz, err = GetMarginGenesisState(s.cfg)
	// s.Require().NoError(err)
	// s.cfg.GenesisState["margin"] = bz

	// bz, err = GetPmtpGenesisState(s.cfg)
	// s.Require().NoError(err)
	// s.cfg.GenesisState["pmtp"] = bz

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")

	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestRowanBalanceExists() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	var genesisState banktypes.GenesisState
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(s.cfg.GenesisState["bank"], &genesisState))
	s.Require().Equal(genesisState.Balances[0].Address, s.address)
	s.Require().Equal(genesisState.Balances[0].Coins[5], types.NewCoin("rowan", s.nativeAmount))

	out, err := QueryBalancesExec(clientCtx, val.Address)
	s.Require().NoError(err)

	var balancesRes banktypes.QueryAllBalancesResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &balancesRes), out.String())

	s.Require().Contains(balancesRes.Balances, types.NewCoin("rowan", s.nativeAmount))
}

func (s *IntegrationTestSuite) TestCLPExists() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	out, err := QueryClpPoolsExec(clientCtx)
	s.Require().NoError(err)

	var poolsRes clptypes.PoolsRes
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolsRes), out.String())

	s.Require().Contains(
		poolsRes.Pools,
		clptypes.NewPool(
			&clptypes.Asset{Symbol: "cdash"},
			types.NewUint(3000000000000000000),
			types.NewUint(3000000000000000000),
			types.NewUint(3000000000000000000),
		),
	)
}
