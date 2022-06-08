package testutil

import (
	"fmt"

	sifapp "github.com/Sifchain/sifnode/app"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
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
	amount, _ := sdk.NewIntFromString("999000000000000000000000000000000")
	s.cfg.StakingTokens = amount

	bz, err := GetBankGenesisState(s.cfg, s.address)
	s.Require().NoError(err)
	s.cfg.GenesisState["bank"] = bz

	bz, err = GetTokenRegistryGenesisState(s.cfg, s.address)
	s.Require().NoError(err)
	s.cfg.GenesisState["tokenregistry"] = bz

	bz, err = GetClpGenesisState(s.cfg)
	s.Require().NoError(err)
	s.cfg.GenesisState["clp"] = bz

	bz, err = GetMarginGenesisState(s.cfg)
	s.Require().NoError(err)
	s.cfg.GenesisState["margin"] = bz

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")

	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestA_MarginPositionsForAddress() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	out, err := QueryMarginPositionsForAddressExec(clientCtx, val.Address)
	s.Require().NoError(err)

	var positionsRes margintypes.PositionsForAddressResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

	s.Require().Empty(positionsRes.Mtps)
}

func (s *IntegrationTestSuite) TestB_OpenLongMTP() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	from := val.Address
	collateralAsset := "rowan"
	collateralAmount := sdk.NewUintFromString("10000000000000000000000")
	borrowAsset := "cusdt"
	position := "long"

	out, err := MsgMarginOpenExec(clientCtx, from, collateralAsset, collateralAmount, borrowAsset, position)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(uint32(0), txResp.Code)

	height, _ := s.network.LatestHeight()

	testCases := []struct {
		height       int64
		forcedClosed bool
		expectedMtp  margintypes.MTP
	}{
		{
			height: height + 1,
			expectedMtp: margintypes.MTP{
				Id:               uint64(1),
				Address:          from.String(),
				CollateralAsset:  collateralAsset,
				CollateralAmount: collateralAmount,
				LiabilitiesP:     collateralAmount,
				LiabilitiesI:     sdk.NewUintFromString("4656613983300"),
				CustodyAsset:     borrowAsset,
				CustodyAmount:    sdk.NewUintFromString("4409900942"),
				Leverage:         sdk.NewUintFromString("1"),
				MtpHealth:        sdk.MustNewDecFromStr("0.101538596164020080"),
				Position:         margintypes.Position_LONG,
			},
		},
		{
			height: height + 5,
			expectedMtp: margintypes.MTP{
				Id:               uint64(1),
				Address:          from.String(),
				CollateralAsset:  collateralAsset,
				CollateralAmount: collateralAmount,
				LiabilitiesP:     collateralAmount,
				LiabilitiesI:     sdk.NewUintFromString("30000000018626846731702"),
				CustodyAsset:     borrowAsset,
				CustodyAmount:    sdk.NewUintFromString("4409900942"),
				Leverage:         sdk.NewUintFromString("1"),
				MtpHealth:        sdk.MustNewDecFromStr("0.077830267765395576"),
				Position:         margintypes.Position_LONG,
			},
		},
		{
			height:       height + 10,
			forcedClosed: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		_, err := s.network.WaitForHeight(tc.height)
		s.Require().NoError(err)

		s.Run(fmt.Sprintf("height: %d", tc.height), func() {
			out, err = QueryMarginPositionsForAddressExec(clientCtx, val.Address)
			s.Require().NoError(err)

			var positionsRes margintypes.PositionsForAddressResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

			if tc.forcedClosed {
				s.Require().Empty(positionsRes.Mtps)
			} else {
				s.T().Logf("mtp: %v", positionsRes.Mtps[0])
				s.Require().Equal(positionsRes.Mtps[0], &tc.expectedMtp)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestC_CloseLongMTP() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	from := val.Address
	id := uint64(1)

	out, err := MsgMarginCloseExec(clientCtx, from, id)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(uint32(0), txResp.Code)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	out, err = QueryMarginPositionsForAddressExec(clientCtx, val.Address)
	s.Require().NoError(err)

	var positionsRes margintypes.PositionsForAddressResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

	s.Require().Empty(positionsRes.Mtps)
}
