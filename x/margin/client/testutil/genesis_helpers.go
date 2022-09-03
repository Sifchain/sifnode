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
	amount, _ := sdk.NewIntFromString("999000000000000000000000000000000")
	balances := []banktypes.Balance{
		{
			Address: address,
			Coins: sdk.Coins{
				sdk.NewCoin("cusdt", amount),
				sdk.NewCoin("rowan", amount),
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
				{Denom: "cusdt", BaseDenom: "cusdt", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
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
			ExternalAsset:                 &clptypes.Asset{Symbol: "cusdt"},
			NativeAssetBalance:            sdk.NewUintFromString("1550459183129248235861408"),
			ExternalAssetBalance:          sdk.NewUintFromString("174248776094"),
			PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
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
		},
	}
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}
