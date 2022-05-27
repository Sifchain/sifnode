package testutil

import (
	"fmt"

	sifapp "github.com/Sifchain/sifnode/app"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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

	// bz, err = GetMarginGenesisState(s.cfg)
	// s.Require().NoError(err)
	// s.cfg.GenesisState["margin"] = bz

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")

	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestRowanBalanceExists() {
	s.T().Log("#################################")
	s.T().Log("TestRowanBalanceExists")
	s.T().Log("#################################")

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	amount, _ := sdk.NewIntFromString("999000000000000000000000000000000")

	var genesisState banktypes.GenesisState
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(s.cfg.GenesisState["bank"], &genesisState))
	s.Require().Equal(genesisState.Balances[0].Address, s.address)
	s.Require().Contains(genesisState.Balances[0].Coins, sdk.NewCoin("rowan", amount))

	out, err := QueryBalancesExec(clientCtx, val.Address)
	s.Require().NoError(err)

	var balancesRes banktypes.QueryAllBalancesResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &balancesRes), out.String())

	s.Require().Contains(balancesRes.Balances, sdk.NewCoin("rowan", amount))
}

func (s *IntegrationTestSuite) TestCLPsExists() {
	s.T().Log("#################################")
	s.T().Log("TestCLPsExists")
	s.T().Log("#################################")

	return

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	out, err := QueryClpPoolsExec(clientCtx)
	s.Require().NoError(err)

	var poolsRes clptypes.PoolsRes
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolsRes), out.String())

	SwapPriceNative := sdk.MustNewDecFromStr("0.000031941094000000")
	SwapPriceExternal := sdk.MustNewDecFromStr("31307.631790289925000000")
	s.Require().Contains(
		poolsRes.Pools,
		&clptypes.Pool{
			ExternalAsset:                 &clptypes.Asset{Symbol: "ceth"},
			NativeAssetBalance:            sdk.NewUintFromString("49352380611368792060339203"),
			ExternalAssetBalance:          sdk.NewUintFromString("1576369012576526264262"),
			PoolUnits:                     sdk.NewUintFromString("49352380611368792060339203"),
			SwapPriceNative:               &SwapPriceNative,
			SwapPriceExternal:             &SwapPriceExternal,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
	)
	SwapPriceNative = sdk.MustNewDecFromStr("0.112507537332000000")
	SwapPriceExternal = sdk.MustNewDecFromStr("8.888293386477892504")
	s.Require().Contains(
		poolsRes.Pools,
		&clptypes.Pool{
			ExternalAsset:                 &clptypes.Asset{Symbol: "cusdc"},
			NativeAssetBalance:            sdk.NewUintFromString("52798591956187184978275830"),
			ExternalAssetBalance:          sdk.NewUintFromString("5940239555604"),
			PoolUnits:                     sdk.NewUintFromString("52798591956187184978275830"),
			SwapPriceNative:               &SwapPriceNative,
			SwapPriceExternal:             &SwapPriceExternal,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
	)
	SwapPriceNative = sdk.MustNewDecFromStr("0.112385271402000000")
	SwapPriceExternal = sdk.MustNewDecFromStr("8.897963118404021251")
	s.Require().Contains(
		poolsRes.Pools,
		&clptypes.Pool{
			ExternalAsset:                 &clptypes.Asset{Symbol: "cusdt"},
			NativeAssetBalance:            sdk.NewUintFromString("1550459183129248235861408"),
			ExternalAssetBalance:          sdk.NewUintFromString("174248776094"),
			PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
			SwapPriceNative:               &SwapPriceNative,
			SwapPriceExternal:             &SwapPriceExternal,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
	)
	SwapPriceNative = sdk.MustNewDecFromStr("0.003536121601000000")
	SwapPriceExternal = sdk.MustNewDecFromStr("282.795704697257746702")
	s.Require().Contains(
		poolsRes.Pools,
		&clptypes.Pool{
			ExternalAsset:                 &clptypes.Asset{Symbol: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"},
			NativeAssetBalance:            sdk.NewUintFromString("200501596725333601567765449"),
			ExternalAssetBalance:          sdk.NewUintFromString("708998027178"),
			PoolUnits:                     sdk.NewUintFromString("200501596725333601567765449"),
			SwapPriceNative:               &SwapPriceNative,
			SwapPriceExternal:             &SwapPriceExternal,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
	)
	SwapPriceNative = sdk.MustNewDecFromStr("0.001004322895000000")
	SwapPriceExternal = sdk.MustNewDecFromStr("995.695712134925308505")
	s.Require().Contains(
		poolsRes.Pools,
		&clptypes.Pool{
			ExternalAsset:                 &clptypes.Asset{Symbol: "ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D"},
			NativeAssetBalance:            sdk.NewUintFromString("29315228314524379224549414"),
			ExternalAssetBalance:          sdk.NewUintFromString("29441954962"),
			PoolUnits:                     sdk.NewUintFromString("29315228314524379224549414"),
			SwapPriceNative:               &SwapPriceNative,
			SwapPriceExternal:             &SwapPriceExternal,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
	)
	SwapPriceNative = sdk.MustNewDecFromStr("0.004243597317000000")
	SwapPriceExternal = sdk.MustNewDecFromStr("235.649126292703172595")
	s.Require().Contains(
		poolsRes.Pools,
		&clptypes.Pool{
			ExternalAsset:                 &clptypes.Asset{Symbol: "ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA"},
			NativeAssetBalance:            sdk.NewUintFromString("32788415426458039601937058"),
			ExternalAssetBalance:          sdk.NewUintFromString("139140831718"),
			PoolUnits:                     sdk.NewUintFromString("32788415426458039601937058"),
			SwapPriceNative:               &SwapPriceNative,
			SwapPriceExternal:             &SwapPriceExternal,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
	)
}

func (s *IntegrationTestSuite) TestPMTPDefaultParams() {
	s.T().Log("#################################")
	s.T().Log("TestPMTPDefaultParams")
	s.T().Log("#################################")

	return

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	out, err := QueryClpPmtpParamsExec(clientCtx)
	s.Require().NoError(err)

	var pmtpParamsRes clptypes.PmtpParamsRes
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &pmtpParamsRes), out.String())

	s.Require().Equal(pmtpParamsRes.Params, &clptypes.PmtpParams{
		PmtpPeriodGovernanceRate: sdk.ZeroDec(),
		PmtpPeriodEpochLength:    1,
		PmtpPeriodStartBlock:     0,
		PmtpPeriodEndBlock:       0,
	})
}

func (s *IntegrationTestSuite) TestModifyPMTPRates() {
	s.T().Log("#################################")
	s.T().Log("TestModifyPMTPRates")
	s.T().Log("#################################")

	return

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	from := val.Address

	blockRate := sdk.MustNewDecFromStr("0.000000458623032662")
	runningRate := sdk.MustNewDecFromStr("1.308075140599690284")

	out, err := MsgClpModifyPmtpRatesExec(
		clientCtx,
		from,
		blockRate,
		runningRate,
	)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(uint32(0), txResp.Code)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	out, err = QueryClpPmtpParamsExec(clientCtx)
	s.Require().NoError(err)

	var pmtpParamsRes clptypes.PmtpParamsRes
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &pmtpParamsRes), out.String())

	s.Require().Equal(pmtpParamsRes.Params, &clptypes.PmtpParams{
		PmtpPeriodGovernanceRate: sdk.ZeroDec(),
		PmtpPeriodEpochLength:    1,
		PmtpPeriodStartBlock:     0,
		PmtpPeriodEndBlock:       0,
	})
	s.Require().Equal(pmtpParamsRes.PmtpRateParams, &clptypes.PmtpRateParams{
		PmtpCurrentRunningRate: runningRate,
		PmtpPeriodBlockRate:    blockRate,
		PmtpInterPolicyRate:    runningRate,
	})
}

func (s *IntegrationTestSuite) TestEndPMTPPolicy() {
	s.T().Log("#################################")
	s.T().Log("TestEndPMTPPolicy")
	s.T().Log("#################################")

	return

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	from := val.Address

	out, err := MsgClpEndPolicyExec(
		clientCtx,
		from,
	)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(uint32(0), txResp.Code)

	out, err = QueryClpPmtpParamsExec(clientCtx)
	s.Require().NoError(err)

	var pmtpParamsRes clptypes.PmtpParamsRes
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &pmtpParamsRes), out.String())

	s.Require().Equal(pmtpParamsRes.Params, &clptypes.PmtpParams{
		PmtpPeriodGovernanceRate: sdk.ZeroDec(),
		PmtpPeriodEpochLength:    1,
		PmtpPeriodStartBlock:     0,
		PmtpPeriodEndBlock:       0,
	})
	s.Require().Equal(pmtpParamsRes.PmtpRateParams, &clptypes.PmtpRateParams{
		PmtpCurrentRunningRate: sdk.ZeroDec(),
		PmtpPeriodBlockRate:    sdk.ZeroDec(),
		PmtpInterPolicyRate:    sdk.ZeroDec(),
	})
}

func (s *IntegrationTestSuite) TestSetNewPMTPPolicy() {
	s.T().Log("#################################")
	s.T().Log("TestSetNewPMTPPolicy")
	s.T().Log("#################################")

	return

	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	from := val.Address
	startBlock := sdk.NewInt(11)
	endBlock := sdk.NewInt(100)
	epochLength := sdk.NewInt(10)
	rGov := sdk.MustNewDecFromStr("0.10")

	out, err := MsgClpUpdatePmtpParamsExec(
		clientCtx,
		from,
		startBlock,
		endBlock,
		epochLength,
		rGov,
	)
	s.Require().NoError(err)

	var respType proto.Message = &sdk.TxResponse{}
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), respType), out.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(uint32(0), txResp.Code)

	out, err = QueryClpPmtpParamsExec(clientCtx)
	s.Require().NoError(err)

	var pmtpParamsRes clptypes.PmtpParamsRes
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &pmtpParamsRes), out.String())

	s.Require().Equal(pmtpParamsRes.Params, &clptypes.PmtpParams{
		PmtpPeriodGovernanceRate: sdk.MustNewDecFromStr("0.100000000000000000"),
		PmtpPeriodEpochLength:    epochLength.Int64(),
		PmtpPeriodStartBlock:     startBlock.Int64(),
		PmtpPeriodEndBlock:       endBlock.Int64(),
	})
	s.Require().Equal(pmtpParamsRes.PmtpRateParams, &clptypes.PmtpRateParams{
		PmtpCurrentRunningRate: sdk.ZeroDec(),
		PmtpPeriodBlockRate:    sdk.ZeroDec(),
		PmtpInterPolicyRate:    sdk.ZeroDec(),
	})

	testCases := []struct {
		height                    int64
		expectedPool              clptypes.Pool
		expectedSwapPriceNative   sdk.Dec
		expectedSwapPriceExternal sdk.Dec
	}{
		{
			height: 8,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: "cusdc"},
				NativeAssetBalance:            sdk.NewUintFromString("52798591956187184978275830"),
				ExternalAssetBalance:          sdk.NewUintFromString("5940239555604"),
				PoolUnits:                     sdk.NewUintFromString("52798591956187184978275830"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.112507537332000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.888293386477892504"),
		},
		{
			height: 11,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: "cusdc"},
				NativeAssetBalance:            sdk.NewUintFromString("52798591956187184978275830"),
				ExternalAssetBalance:          sdk.NewUintFromString("5940239555604"),
				PoolUnits:                     sdk.NewUintFromString("52798591956187184978275830"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.113584975076283602"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.803981330500189701"),
		},
		{
			height: 12,
			expectedPool: clptypes.Pool{
				ExternalAsset:                 &clptypes.Asset{Symbol: "cusdc"},
				NativeAssetBalance:            sdk.NewUintFromString("52798591956187184978275830"),
				ExternalAssetBalance:          sdk.NewUintFromString("5940239555604"),
				PoolUnits:                     sdk.NewUintFromString("52798591956187184978275830"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.114672730992312279"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.720469036914894168"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		_, err := s.network.WaitForHeight(tc.height)
		s.Require().NoError(err)

		s.Run(fmt.Sprintf("height: %d", tc.height), func() {
			out, err := QueryClpPoolsExec(clientCtx)
			s.Require().NoError(err)

			var poolsRes clptypes.PoolsRes
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &poolsRes), out.String())

			s.T().Log("######################################")
			s.T().Log(poolsRes.Pools[1])
			s.T().Log("######################################")

			tc.expectedPool.SwapPriceNative = &tc.expectedSwapPriceNative
			tc.expectedPool.SwapPriceExternal = &tc.expectedSwapPriceExternal
			s.Require().Contains(poolsRes.Pools, &tc.expectedPool)
		})
	}
}

func (s *IntegrationTestSuite) TestResetPMTPParams() {
	s.T().Log("#################################")
	s.T().Log("TestResetPMTPParams")
	s.T().Log("#################################")
}

func (s *IntegrationTestSuite) TestEndPolicy() {
	s.T().Log("#################################")
	s.T().Log("TestEndPolicy")
	s.T().Log("#################################")
}
