//go:build TEST_INTEGRATION
// +build TEST_INTEGRATION

package testutil

import (
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func GetBankGenesisState(cfg network.Config, address string) ([]byte, error) {
	externalAmount, _ := sdk.NewIntFromString("100000000000")
	nativeAmount, _ := sdk.NewIntFromString("100000000000000000000000")
	balances := []banktypes.Balance{
		{
			Address: address,
			Coins: sdk.Coins{
				sdk.NewCoin("cusdc", externalAmount),
				sdk.NewCoin("rowan", nativeAmount),
			},
		},
	}
	gs := banktypes.DefaultGenesisState()
	gs.Balances = append(gs.Balances, balances...)
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

func GetAdminGenesisState(cfg network.Config, address string) ([]byte, error) {
	gs := &admintypes.GenesisState{
		AdminAccounts: []*admintypes.AdminAccount{
			{
				AdminType:    admintypes.AdminType_ADMIN,
				AdminAddress: address,
			},
			{
				AdminType:    admintypes.AdminType_CLPDEX,
				AdminAddress: address,
			},
			{
				AdminType:    admintypes.AdminType_MARGIN,
				AdminAddress: address,
			},
			{
				AdminType:    admintypes.AdminType_ETHBRIDGE,
				AdminAddress: address,
			},
			{
				AdminType:    admintypes.AdminType_PMTPREWARDS,
				AdminAddress: address,
			},
			{
				AdminType:    admintypes.AdminType_TOKENREGISTRY,
				AdminAddress: address,
			},
		},
	}
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

func GetTokenRegistryGenesisState(cfg network.Config, address string) ([]byte, error) {
	gs := &tokenregistrytypes.GenesisState{
		Registry: &tokenregistrytypes.Registry{
			Entries: []*tokenregistrytypes.RegistryEntry{
				{Denom: "node0token", BaseDenom: "node0token", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
				{Denom: "cusdc", BaseDenom: "cusdc", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
				{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
			},
		},
	}
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

func GetClpGenesisState(cfg network.Config) ([]byte, error) {
	pools := []*clptypes.Pool{
		{
			ExternalAsset:                 &clptypes.Asset{Symbol: "cusdc"},
			NativeAssetBalance:            sdk.NewUintFromString("100000000000000000000000000"),
			ExternalAssetBalance:          sdk.NewUintFromString("1000000000000"),
			PoolUnits:                     sdk.NewUintFromString("100000000000000000000000000"),
			SwapPriceNative:               nil,
			SwapPriceExternal:             nil,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
	}
	gs := clptypes.DefaultGenesisState()
	gs.PoolList = append(gs.PoolList, pools...)
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

func GetMarginGenesisState(cfg network.Config) ([]byte, error) {
	gs := &margintypes.GenesisState{
		Params: &margintypes.Params{
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
		},
	}
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}
