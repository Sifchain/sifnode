//go:build TEST_INTEGRATION
// +build TEST_INTEGRATION

package testutil

import (
	"fmt"

	sifapp "github.com/Sifchain/sifnode/app"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
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

	bz, err = GetAdminGenesisState(s.cfg, s.address)
	s.Require().NoError(err)
	s.cfg.GenesisState["admin"] = bz

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

func (s *IntegrationTestSuite) TestA1_MarginParams() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	out, err := QueryMarginParamsExec(clientCtx)
	s.Require().NoError(err)

	var res margintypes.ParamsResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())

	s.Require().Equal(res.Params, &margintypes.Params{
		LeverageMax:           sdk.MustNewDecFromStr("2.0"),
		HealthGainFactor:      sdk.MustNewDecFromStr("1.0"),
		InterestRateMin:       sdk.MustNewDecFromStr("0.005"),
		InterestRateMax:       sdk.MustNewDecFromStr("3.0"),
		InterestRateDecrease:  sdk.MustNewDecFromStr("0.001"),
		InterestRateIncrease:  sdk.MustNewDecFromStr("0.001"),
		RemovalQueueThreshold: sdk.MustNewDecFromStr("0.1"),
		EpochLength:           1,
		MaxOpenPositions:      10000,
		Pools:                 []string{"cusdt"},
	})
}

func (s *IntegrationTestSuite) TestA2_MarginPositionsForAddress() {
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
	collateralAmount := sdk.NewUintFromString("10000000000000000000000") // 10000 rowan
	borrowAsset := "cusdt"
	position := "long"
	leverage := sdk.MustNewDecFromStr("2.0")

	out, err := MsgMarginOpenExec(clientCtx, from, collateralAsset, collateralAmount, borrowAsset, position, leverage)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)

	// fmt.Println("txResp:", txResp)

	s.Require().Equal(uint32(0), txResp.Code)

	testCases := []struct {
		height                    int64
		expectedPool              clptypes.Pool
		expectedSwapPriceNative   sdk.Dec
		expectedSwapPriceExternal sdk.Dec
		forcedClosed              bool
		expectedMtp               margintypes.MTP
	}{
		{
			height: 9,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"), // 1560459 rowan
				ExternalAssetBalance:          sdk.NewUintFromString("172022630705"),              // 169838 cusdt
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("2226145389"), // 4409 cusdt
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("0.013000000000000000"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              collateralAmount,
				InterestUnpaidCollateral: sdk.NewUintFromString("72759593488"),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("2226145389"),
				Leverage:                 sdk.MustNewDecFromStr("2.0"),
				MtpHealth:                sdk.MustNewDecFromStr("0.168454370237483891"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
		{
			height: 10,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("172022630705"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("2226145389"),
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("0.014000000000000000"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              collateralAmount,
				InterestUnpaidCollateral: sdk.NewUintFromString("72760009820"),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("2226145389"),
				Leverage:                 sdk.MustNewDecFromStr("2.0"),
				MtpHealth:                sdk.MustNewDecFromStr("0.168454370237277422"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
		{
			height: 11,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("172022630705"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("2226145389"),
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("0.015000000000000000"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              collateralAmount,
				InterestUnpaidCollateral: sdk.NewUintFromString("72760703708"),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("2226145389"),
				Leverage:                 sdk.MustNewDecFromStr("2.0"),
				MtpHealth:                sdk.MustNewDecFromStr("0.168454370237277420"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
		{
			height: 12,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("172022630705"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("2226145389"),
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("0.016000000000000000"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              collateralAmount,
				InterestUnpaidCollateral: sdk.NewUintFromString("72761397596"),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("2226145389"),
				Leverage:                 sdk.MustNewDecFromStr("2.0"),
				MtpHealth:                sdk.MustNewDecFromStr("0.168454370237277418"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
		{
			height: 13,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("172022630705"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("2226145389"),
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("0.017000000000000000"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              collateralAmount,
				InterestUnpaidCollateral: sdk.NewUintFromString("72761987401"),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("2226145389"),
				Leverage:                 sdk.MustNewDecFromStr("2.0"),
				MtpHealth:                sdk.MustNewDecFromStr("0.168454370237277416"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		_, err := s.network.WaitForHeight(tc.height)
		s.Require().NoError(err)

		s.Run(fmt.Sprintf("height: %d", tc.height), func() {
			out, err := QueryClpPoolExec(clientCtx, borrowAsset)
			s.Require().NoError(err)

			var poolRes clptypes.PoolRes
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolRes), out.String())

			tc.expectedPool.SwapPriceNative = &tc.expectedSwapPriceNative
			tc.expectedPool.SwapPriceExternal = &tc.expectedSwapPriceExternal
			s.T().Log("pool:", *poolRes.Pool)
			s.T().Log("expected pool:", tc.expectedPool)
			s.Require().Equal(poolRes.Pool, &tc.expectedPool)

			out, err = QueryMarginPositionsForAddressExec(clientCtx, val.Address)
			s.Require().NoError(err)

			var positionsRes margintypes.PositionsForAddressResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

			if tc.forcedClosed {
				s.Require().Empty(positionsRes.Mtps)
			} else {
				s.Require().NotEmpty(positionsRes.Mtps)
				s.T().Log("mtp:", *positionsRes.Mtps[0])
				s.T().Log("expected mtp:", tc.expectedMtp)
				s.Require().Equal(positionsRes.Mtps[0], &tc.expectedMtp)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestC_CloseLongMTP() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	from := val.Address

	out, err := QueryMarginPositionsForAddressExec(clientCtx, val.Address)
	s.Require().NoError(err)

	var positionsRes margintypes.PositionsForAddressResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

	s.Require().NotEmpty(positionsRes.Mtps)

	out, err = MsgMarginCloseExec(clientCtx, from, positionsRes.Mtps[0].Id)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(uint32(0), txResp.Code)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	out, err = QueryMarginPositionsForAddressExec(clientCtx, val.Address)
	s.Require().NoError(err)

	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

	s.Require().Empty(positionsRes.Mtps)
}
