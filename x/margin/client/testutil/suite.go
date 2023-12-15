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
		HealthGainFactor:                         sdk.MustNewDecFromStr("0.000000022"),
		InterestRateDecrease:                     sdk.MustNewDecFromStr("0.000000000333333333"),
		InterestRateIncrease:                     sdk.MustNewDecFromStr("0.000000000333333333"),
		InterestRateMin:                          sdk.MustNewDecFromStr("0.00000021"),
		InterestRateMax:                          sdk.MustNewDecFromStr("0.00000001"),
		LeverageMax:                              sdk.MustNewDecFromStr("10.0"),
		EpochLength:                              1,
		RemovalQueueThreshold:                    sdk.MustNewDecFromStr("0.35"),
		MaxOpenPositions:                         10000,
		ForceCloseFundPercentage:                 sdk.MustNewDecFromStr("1.0"),
		ForceCloseFundAddress:                    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
		IncrementalInterestPaymentEnabled:        true,
		IncrementalInterestPaymentFundPercentage: sdk.MustNewDecFromStr("0.35"),
		IncrementalInterestPaymentFundAddress:    "sif15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqhns3lt",
		PoolOpenThreshold:                        sdk.MustNewDecFromStr("0.65"),
		SqModifier:                               sdk.MustNewDecFromStr("10000000000000000000000000"),
		SafetyFactor:                             sdk.MustNewDecFromStr("1.05"),
		WhitelistingEnabled:                      false,
		Pools:                                    []string{"cusdc"},
		ClosedPools:                              []string{},
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
	collateralAsset := "cusdc"
	collateralAmount := sdk.NewUintFromString("100000000") // 1000 cusdc
	borrowAsset := "rowan"
	position := "long"
	leverage := sdk.MustNewDecFromStr("10.0")

	// before opening position check the pool state
	out, err := QueryClpPoolExec(clientCtx, collateralAsset)
	s.Require().NoError(err)

	var poolRes clptypes.PoolRes
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolRes), out.String())

	height, err := s.network.LatestHeight()
	s.Require().NoError(err)

	spn := sdk.MustNewDecFromStr("0.01")
	spe := sdk.MustNewDecFromStr("100")

	expectedPool := clptypes.Pool{
		ExternalAsset:                  &clptypes.Asset{Symbol: collateralAsset},
		NativeAssetBalance:             sdk.NewUintFromString("100000000000000000000000000"),
		ExternalAssetBalance:           sdk.NewUintFromString("1000000000000"),
		PoolUnits:                      sdk.NewUintFromString("100000000000000000000000000"),
		SwapPriceNative:                &spn,
		SwapPriceExternal:              &spe,
		ExternalLiabilities:            sdk.NewUintFromString("0"),
		ExternalCustody:                sdk.NewUintFromString("0"),
		NativeLiabilities:              sdk.NewUintFromString("0"),
		NativeCustody:                  sdk.NewUintFromString("0"),
		Health:                         sdk.MustNewDecFromStr("1.0"),
		InterestRate:                   sdk.MustNewDecFromStr("0.00000021"),
		RewardPeriodNativeDistributed:  sdk.ZeroUint(),
		LastHeightInterestRateComputed: height,
		UnsettledExternalLiabilities:   sdk.ZeroUint(),
		UnsettledNativeLiabilities:     sdk.ZeroUint(),
		BlockInterestNative:            sdk.ZeroUint(),
		BlockInterestExternal:          sdk.ZeroUint(),
	}

	s.T().Log("pool:", *poolRes.Pool)
	s.T().Log("expected pool:", expectedPool)
	s.Require().Equal(poolRes.Pool, &expectedPool)

	// open position
	out, err = MsgMarginOpenExec(clientCtx, from, collateralAsset, collateralAmount, borrowAsset, position, leverage)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)

	// fmt.Println("txResp:", txResp)

	// s.Require().Contains(txResp.RawLog, "failed")
	// s.Require().Contains(txResp.RawLog, "cusdc: margin not enabled for pool")
	s.Require().Equal(uint32(0), txResp.Code)

	// check the pool again at opening block
	out, err = QueryClpPoolExec(clientCtx, collateralAsset)
	s.Require().NoError(err)

	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolRes), out.String())

	height, err = s.network.LatestHeight()
	s.Require().NoError(err)

	spn = sdk.MustNewDecFromStr("0.01")
	spe = sdk.MustNewDecFromStr("100")

	expectedPool = clptypes.Pool{
		ExternalAsset:                  &clptypes.Asset{Symbol: collateralAsset},
		NativeAssetBalance:             sdk.NewUintFromString("99900399600399600399600400"),
		ExternalAssetBalance:           sdk.NewUintFromString("1000100000000"),
		PoolUnits:                      sdk.NewUintFromString("100000000000000000000000000"),
		SwapPriceNative:                &spn,
		SwapPriceExternal:              &spe,
		ExternalLiabilities:            sdk.NewUintFromString("900000000"),
		ExternalCustody:                sdk.NewUintFromString("0"),
		NativeLiabilities:              sdk.NewUintFromString("0"),
		NativeCustody:                  sdk.NewUintFromString("99600399600399600399600"),
		Health:                         sdk.MustNewDecFromStr("0.999100899100899101"),
		InterestRate:                   sdk.MustNewDecFromStr("0.00000021"),
		RewardPeriodNativeDistributed:  sdk.ZeroUint(),
		LastHeightInterestRateComputed: height,
		UnsettledExternalLiabilities:   sdk.ZeroUint(),
		UnsettledNativeLiabilities:     sdk.ZeroUint(),
		BlockInterestNative:            sdk.ZeroUint(),
		BlockInterestExternal:          sdk.ZeroUint(),
	}

	s.T().Log("pool:", *poolRes.Pool)
	s.T().Log("expected pool:", expectedPool)
	s.Require().Equal(poolRes.Pool, &expectedPool)

	testCases := []struct {
		expectedPool              clptypes.Pool
		expectedSwapPriceNative   sdk.Dec
		expectedSwapPriceExternal sdk.Dec
		forcedClosed              bool
		expectedMtp               margintypes.MTP
	}{
		{
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: collateralAsset},
				NativeAssetBalance:            sdk.NewUintFromString("99900399600399600399600400"),
				ExternalAssetBalance:          sdk.NewUintFromString("1000100000000"),
				PoolUnits:                     sdk.NewUintFromString("100000000000000000000000000"),
				ExternalLiabilities:           sdk.NewUintFromString("900000000"),
				ExternalCustody:               sdk.NewUintFromString("0"),
				NativeLiabilities:             sdk.NewUintFromString("0"),
				NativeCustody:                 sdk.NewUintFromString("99600399600399600399600"),
				Health:                        sdk.MustNewDecFromStr("0.999100899100899101"),
				InterestRate:                  sdk.MustNewDecFromStr("0.00000021"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.010010000000000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("99.900099900099900099"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              sdk.NewUintFromString("900000000"),
				InterestPaidCollateral:   sdk.NewUintFromString("189"),
				InterestPaidCustody:      sdk.NewUintFromString("18824475520921252"),
				InterestUnpaidCollateral: sdk.ZeroUint(),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("99600380775924079478348"),
				Leverage:                 sdk.MustNewDecFromStr("10.0"),
				MtpHealth:                sdk.MustNewDecFromStr("1.103355497777777778"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
		{
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: collateralAsset},
				NativeAssetBalance:            sdk.NewUintFromString("99900399600399600399600400"),
				ExternalAssetBalance:          sdk.NewUintFromString("1000100000000"),
				PoolUnits:                     sdk.NewUintFromString("100000000000000000000000000"),
				ExternalLiabilities:           sdk.NewUintFromString("900000000"),
				ExternalCustody:               sdk.NewUintFromString("0"),
				NativeLiabilities:             sdk.NewUintFromString("0"),
				NativeCustody:                 sdk.NewUintFromString("99600399600399600399600"),
				Health:                        sdk.MustNewDecFromStr("0.999100899100899101"),
				InterestRate:                  sdk.MustNewDecFromStr("0.00000021"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              sdk.NewUintFromString("900000000"),
				InterestPaidCollateral:   sdk.NewUintFromString("189"),
				InterestPaidCustody:      sdk.NewUintFromString("18824475520921252"),
				InterestUnpaidCollateral: sdk.ZeroUint(),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("99600380775924079478348"),
				Leverage:                 sdk.MustNewDecFromStr("10.0"),
				MtpHealth:                sdk.MustNewDecFromStr("1.103355497777777778"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
		{
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: collateralAsset},
				NativeAssetBalance:            sdk.NewUintFromString("99900399600399600399600400"),
				ExternalAssetBalance:          sdk.NewUintFromString("1000100000000"),
				PoolUnits:                     sdk.NewUintFromString("100000000000000000000000000"),
				ExternalLiabilities:           sdk.NewUintFromString("900000000"),
				ExternalCustody:               sdk.NewUintFromString("0"),
				NativeLiabilities:             sdk.NewUintFromString("0"),
				NativeCustody:                 sdk.NewUintFromString("99600399600399600399600"),
				Health:                        sdk.MustNewDecFromStr("0.999100899100899101"),
				InterestRate:                  sdk.MustNewDecFromStr("0.00000021"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              sdk.NewUintFromString("900000000"),
				InterestPaidCollateral:   sdk.NewUintFromString("189"),
				InterestPaidCustody:      sdk.NewUintFromString("18824475520921252"),
				InterestUnpaidCollateral: sdk.ZeroUint(),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("99600380775924079478348"),
				Leverage:                 sdk.MustNewDecFromStr("10.0"),
				MtpHealth:                sdk.MustNewDecFromStr("1.103355497777777778"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
		{
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: collateralAsset},
				NativeAssetBalance:            sdk.NewUintFromString("99900399600399600399600400"),
				ExternalAssetBalance:          sdk.NewUintFromString("1000100000000"),
				PoolUnits:                     sdk.NewUintFromString("100000000000000000000000000"),
				ExternalLiabilities:           sdk.NewUintFromString("900000000"),
				ExternalCustody:               sdk.NewUintFromString("0"),
				NativeLiabilities:             sdk.NewUintFromString("0"),
				NativeCustody:                 sdk.NewUintFromString("99600399600399600399600"),
				Health:                        sdk.MustNewDecFromStr("0.999100899100899101"),
				InterestRate:                  sdk.MustNewDecFromStr("0.00000021"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              sdk.NewUintFromString("900000000"),
				InterestPaidCollateral:   sdk.NewUintFromString("189"),
				InterestPaidCustody:      sdk.NewUintFromString("18824475520921252"),
				InterestUnpaidCollateral: sdk.ZeroUint(),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("99600380775924079478348"),
				Leverage:                 sdk.MustNewDecFromStr("10.0"),
				MtpHealth:                sdk.MustNewDecFromStr("1.103355497777777778"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
		{
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: collateralAsset},
				NativeAssetBalance:            sdk.NewUintFromString("99900399600399600399600400"),
				ExternalAssetBalance:          sdk.NewUintFromString("1000100000000"),
				PoolUnits:                     sdk.NewUintFromString("100000000000000000000000000"),
				ExternalLiabilities:           sdk.NewUintFromString("900000000"),
				ExternalCustody:               sdk.NewUintFromString("0"),
				NativeLiabilities:             sdk.NewUintFromString("0"),
				NativeCustody:                 sdk.NewUintFromString("99600399600399600399600"),
				Health:                        sdk.MustNewDecFromStr("0.999100899100899101"),
				InterestRate:                  sdk.MustNewDecFromStr("0.00000021"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
				UnsettledExternalLiabilities:  sdk.ZeroUint(),
				UnsettledNativeLiabilities:    sdk.ZeroUint(),
				BlockInterestNative:           sdk.ZeroUint(),
				BlockInterestExternal:         sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.111669710297000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.954979800030249027"),
			expectedMtp: margintypes.MTP{
				Address:                  from.String(),
				CollateralAsset:          collateralAsset,
				CollateralAmount:         collateralAmount,
				Liabilities:              sdk.NewUintFromString("900000000"),
				InterestPaidCollateral:   sdk.NewUintFromString("189"),
				InterestPaidCustody:      sdk.NewUintFromString("18824475520921252"),
				InterestUnpaidCollateral: sdk.ZeroUint(),
				CustodyAsset:             borrowAsset,
				CustodyAmount:            sdk.NewUintFromString("99600380775924079478348"),
				Leverage:                 sdk.MustNewDecFromStr("10.0"),
				MtpHealth:                sdk.MustNewDecFromStr("1.103355497777777778"),
				Position:                 margintypes.Position_LONG,
				Id:                       uint64(1),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		err := s.network.WaitForNextBlock()
		s.Require().NoError(err)

		height, err := s.network.LatestHeight()
		s.Require().NoError(err)

		s.Run(fmt.Sprintf("height: %d", height), func() {
			out, err := QueryClpPoolExec(clientCtx, collateralAsset)
			s.Require().NoError(err)

			var poolRes clptypes.PoolRes
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolRes), out.String())

			tc.expectedPool.SwapPriceNative = &tc.expectedSwapPriceNative
			tc.expectedPool.SwapPriceExternal = &tc.expectedSwapPriceExternal
			tc.expectedPool.LastHeightInterestRateComputed = height
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
