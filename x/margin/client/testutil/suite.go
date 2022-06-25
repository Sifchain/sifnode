//go:build FEATURE_TOGGLE_SDK_045 && FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_SDK_045,FEATURE_TOGGLE_MARGIN_CLI_ALPHA

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
		LeverageMax:          sdk.NewUintFromString("2"),
		HealthGainFactor:     sdk.MustNewDecFromStr("1.0"),
		InterestRateMin:      sdk.MustNewDecFromStr("0.005"),
		InterestRateMax:      sdk.MustNewDecFromStr("3.0"),
		InterestRateDecrease: sdk.MustNewDecFromStr("0.10"),
		InterestRateIncrease: sdk.MustNewDecFromStr("0.10"),
		ForceCloseThreshold:  sdk.MustNewDecFromStr("0.10"),
		EpochLength:          1,
		Pools:                []string{"cusdt"},
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

	out, err := MsgMarginOpenExec(clientCtx, from, collateralAsset, collateralAmount, borrowAsset, position)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(uint32(0), txResp.Code)

	height, _ := s.network.LatestHeight()

	testCases := []struct {
		height                    int64
		expectedPool              clptypes.Pool
		expectedSwapPriceNative   sdk.Dec
		expectedSwapPriceExternal sdk.Dec
		forcedClosed              bool
		expectedMtp               margintypes.MTP
	}{
		{
			height: 1,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"), // 1560459 rowan
				ExternalAssetBalance:          sdk.NewUintFromString("174248776094"),              // 169838 cusdt
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("4409900942"), // 4409 cusdt
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("0.900000000000000000"),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.113114828359000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.840573906129206560"),
			expectedMtp: margintypes.MTP{
				Id:               uint64(1),
				Address:          from.String(),
				CollateralAsset:  collateralAsset,
				CollateralAmount: collateralAmount,
				LiabilitiesP:     collateralAmount,
				LiabilitiesI:     sdk.NewUintFromString("4656613983300"),
				CustodyAsset:     borrowAsset,
				CustodyAmount:    sdk.NewUintFromString("4409900942"),
				Leverage:         sdk.NewUintFromString("2"),
				MtpHealth:        sdk.MustNewDecFromStr("0.102579668460296506"),
				Position:         margintypes.Position_LONG,
			},
		},
		{
			height: 2,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("169838875152"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("4409900942"),
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("1.000000000000000000"),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.110252109898000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("9.070121205951520641"),
			expectedMtp: margintypes.MTP{
				Id:               uint64(1),
				Address:          from.String(),
				CollateralAsset:  collateralAsset,
				CollateralAmount: collateralAmount,
				LiabilitiesP:     collateralAmount,
				LiabilitiesI:     sdk.NewUintFromString("30000000018626455933200"),
				CustodyAsset:     borrowAsset,
				CustodyAmount:    sdk.NewUintFromString("4409900942"),
				Leverage:         sdk.NewUintFromString("2"),
				MtpHealth:        sdk.MustNewDecFromStr("0.102579668455396542"),
				Position:         margintypes.Position_LONG,
			},
		},
		{
			height: 3,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("169838875152"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("4409900942"),
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("1.000000000000000000"),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.110252109898000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("9.070121205951520641"),
			forcedClosed:              true,
		},
		{
			height: 4,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("169838875152"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("4409900942"),
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("1.000000000000000000"),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.110252109898000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("9.070121205951520641"),
			forcedClosed:              true,
		},
		{
			height: 5,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: borrowAsset},
				NativeAssetBalance:            sdk.NewUintFromString("1540459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("169838875152"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				ExternalLiabilities:           sdk.NewUintFromString("0"),
				ExternalCustody:               sdk.NewUintFromString("4409900942"),
				NativeLiabilities:             collateralAmount,
				NativeCustody:                 sdk.NewUintFromString("0"),
				Health:                        sdk.MustNewDecFromStr("0.993550297802862968"),
				InterestRate:                  sdk.MustNewDecFromStr("1.000000000000000000"),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.110252109898000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("9.070121205951520641"),
			forcedClosed:              true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		_, err := s.network.WaitForHeight(height + tc.height)
		s.Require().NoError(err)

		s.Run(fmt.Sprintf("height: %d (%d)", height+tc.height, tc.height), func() {
			out, err := QueryClpPoolExec(clientCtx, borrowAsset)
			s.Require().NoError(err)

			var poolRes clptypes.PoolRes
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolRes), out.String())

			tc.expectedPool.SwapPriceNative = &tc.expectedSwapPriceNative
			tc.expectedPool.SwapPriceExternal = &tc.expectedSwapPriceExternal
			s.T().Logf("pool: %v", poolRes.Pool)
			s.Require().Equal(poolRes.Pool, &tc.expectedPool)

			out, err = QueryMarginPositionsForAddressExec(clientCtx, val.Address)
			s.Require().NoError(err)

			var positionsRes margintypes.PositionsForAddressResponse
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

			if tc.forcedClosed {
				s.Require().Empty(positionsRes.Mtps)
			} else {
				s.Require().NotEmpty(positionsRes.Mtps)
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

	out, err = QueryMarginPositionsForAddressExec(clientCtx, val.Address)
	s.Require().NoError(err)

	var positionsRes margintypes.PositionsForAddressResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

	s.Require().NotEmpty(positionsRes.Mtps)

	out, err = MsgMarginCloseExec(clientCtx, from, positionsRes.Mtps[0].Id)
	s.Require().NoError(err)

	respType = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(uint32(0), txResp.Code)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	out, err = QueryMarginPositionsForAddressExec(clientCtx, val.Address)
	s.Require().NoError(err)

	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &positionsRes), out.String())

	s.Require().Empty(positionsRes.Mtps)
}
