//go:build integration
// +build integration

package testutil

import (
	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	mnemonic string
	address  string
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	sifapp.SetConfig(false)

	s.mnemonic = "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
	s.address = "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"

	s.cfg.Mnemonics = []string{s.mnemonic}

	genesisState := banktypes.DefaultGenesisState()
	genesisState.Balances = append(genesisState.Balances, banktypes.Balance{Address: s.address, Coins: types.Coins{types.Coin{Denom: "rowan", Amount: types.NewInt(1000000000000000000)}}})
	bz, err := s.cfg.Codec.MarshalJSON(genesisState)
	s.Require().NoError(err)
	s.cfg.GenesisState["bank"] = bz

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
	s.Require().Contains(genesisState.Balances[0].Coins, types.Coin{Denom: "rowan", Amount: types.NewInt(1000000000000000000)})

	out, err := QueryBalancesExec(clientCtx, val.Address)
	s.Require().NoError(err)

	var balancesRes banktypes.QueryAllBalancesResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &balancesRes), out.String())

	s.Require().Contains(balancesRes.Balances, types.Coin{Denom: "rowan", Amount: types.NewInt(1000000000000000000)})
}
