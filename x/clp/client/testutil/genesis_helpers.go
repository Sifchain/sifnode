package testutil

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// margintypes "github.com/cosmos/cosmos-sdk/x/margin/types"
	// pmtptypes "github.com/cosmos/cosmos-sdk/x/pmtp/types"
)

func GetBankGenesisState(cfg network.Config, address string, nativeAmount sdk.Int, externalAmount sdk.Int) ([]byte, error) {
	balances := []banktypes.Balance{
		{
			Address: address,
			Coins: sdk.Coins{
				sdk.NewCoin("catk", externalAmount),
				sdk.NewCoin("cbtk", externalAmount),
				sdk.NewCoin("cdash", externalAmount),
				sdk.NewCoin("ceth", externalAmount),
				sdk.NewCoin("clink", externalAmount),
				sdk.NewCoin("rowan", nativeAmount),
			},
		},
	}
	gs := banktypes.DefaultGenesisState()
	gs.Balances = append(gs.Balances, balances...)
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

func GetTokenRegistryGenesisState(cfg network.Config, address string) ([]byte, error) {
	gs := &tokenregistrytypes.GenesisState{
		AdminAccounts: &tokenregistrytypes.AdminAccounts{
			AdminAccounts: []*tokenregistrytypes.AdminAccount{
				&tokenregistrytypes.AdminAccount{
					AdminType:    tokenregistrytypes.AdminType_CLPDEX,
					AdminAddress: address,
				},
				&tokenregistrytypes.AdminAccount{
					AdminType:    tokenregistrytypes.AdminType_ETHBRIDGE,
					AdminAddress: address,
				},
				&tokenregistrytypes.AdminAccount{
					AdminType:    tokenregistrytypes.AdminType_PMTPREWARDS,
					AdminAddress: address,
				},
				&tokenregistrytypes.AdminAccount{
					AdminType:    tokenregistrytypes.AdminType_TOKENREGISTRY,
					AdminAddress: address,
				},
			},
		},
		Registry: &tokenregistrytypes.Registry{
			Entries: []*tokenregistrytypes.RegistryEntry{
				{Denom: "node0token", BaseDenom: "node0token", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "catk", BaseDenom: "catk", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "cbtk", BaseDenom: "cbtk", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "cdash", BaseDenom: "cdash", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "ceth", BaseDenom: "ceth", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "clink", BaseDenom: "clink", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
				{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
			},
		},
	}
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

func GetClpGenesisState(cfg network.Config, pool1Amount sdk.Uint, pool2Amount sdk.Uint) ([]byte, error) {
	SwapPriceNative := sdk.ZeroDec()
	SwapPriceExternal := sdk.ZeroDec()

	pools := []*clptypes.Pool{
		{
			ExternalAsset:                 &clptypes.Asset{Symbol: "cdash"},
			NativeAssetBalance:            pool1Amount,
			ExternalAssetBalance:          pool1Amount,
			PoolUnits:                     pool1Amount,
			SwapPriceNative:               &SwapPriceNative,
			SwapPriceExternal:             &SwapPriceExternal,
			RewardPeriodNativeDistributed: types.ZeroUint(),
		},
		{
			ExternalAsset:                 &clptypes.Asset{Symbol: "ceth"},
			NativeAssetBalance:            pool2Amount,
			ExternalAssetBalance:          pool2Amount,
			PoolUnits:                     pool2Amount,
			SwapPriceNative:               &SwapPriceNative,
			SwapPriceExternal:             &SwapPriceExternal,
			RewardPeriodNativeDistributed: types.ZeroUint(),
		},
	}
	gs := clptypes.DefaultGenesisState()
	gs.PoolList = append(gs.PoolList, pools...)
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}

// func GetMarginGenesisState(cfg network.Config) ([]byte, error) {
// 	gs := margintypes.DefaultGenesis()
// 	bz, err := cfg.Codec.MarshalJSON(gs)
// 	return bz, err
// }

// func GetPmtpGenesisState(cfg network.Config) ([]byte, error) {
// 	gs := pmtptypes.DefaultGenesis()
// 	bz, err := cfg.Codec.MarshalJSON(gs)
// 	return bz, err
// }
