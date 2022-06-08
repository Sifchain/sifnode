package testutil

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// margintypes "github.com/cosmos/cosmos-sdk/x/margin/types"
)

func GetBankGenesisState(cfg network.Config, address string) ([]byte, error) {
	amount, _ := sdk.NewIntFromString("999000000000000000000000000000000")
	balances := []banktypes.Balance{
		{
			Address: address,
			Coins: sdk.Coins{
				sdk.NewCoin("ceth", amount),
				sdk.NewCoin("cusdc", amount),
				sdk.NewCoin("cusdt", amount),
				sdk.NewCoin("ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", amount), // atom
				sdk.NewCoin("ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D", amount), // luna
				sdk.NewCoin("ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA", amount), // juno
				sdk.NewCoin("rowan", amount),
				// sdk.NewCoin("stake", amount),
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
				{
					AdminType:    tokenregistrytypes.AdminType_CLPDEX,
					AdminAddress: address,
				},
				{
					AdminType:    tokenregistrytypes.AdminType_ETHBRIDGE,
					AdminAddress: address,
				},
				{
					AdminType:    tokenregistrytypes.AdminType_PMTPREWARDS,
					AdminAddress: address,
				},
				{
					AdminType:    tokenregistrytypes.AdminType_TOKENREGISTRY,
					AdminAddress: address,
				},
			},
		},
		Registry: &tokenregistrytypes.Registry{
			Entries: []*tokenregistrytypes.RegistryEntry{
				{Denom: "node0token", BaseDenom: "node0token", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
				{Denom: "ceth", BaseDenom: "ceth", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
				{Denom: "cusdc", BaseDenom: "cusdc", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
				{Denom: "cusdt", BaseDenom: "cusdt", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
				{Denom: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", BaseDenom: "uatom", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
				{Denom: "ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D", BaseDenom: "uluna", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
				{Denom: "ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA", BaseDenom: "ujuno", Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP, tokenregistrytypes.Permission_IBCEXPORT, tokenregistrytypes.Permission_IBCIMPORT}},
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
			ExternalAsset:                 &clptypes.Asset{Symbol: "ceth"},
			NativeAssetBalance:            sdk.NewUintFromString("49352380611368792060339203"),
			ExternalAssetBalance:          sdk.NewUintFromString("1576369012576526264262"),
			PoolUnits:                     sdk.NewUintFromString("49352380611368792060339203"),
			SwapPriceNative:               nil,
			SwapPriceExternal:             nil,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
		{
			ExternalAsset:                 &clptypes.Asset{Symbol: "cusdc"},
			NativeAssetBalance:            sdk.NewUintFromString("52798591956187184978275830"),
			ExternalAssetBalance:          sdk.NewUintFromString("5940239555604"),
			PoolUnits:                     sdk.NewUintFromString("52798591956187184978275830"),
			SwapPriceNative:               nil,
			SwapPriceExternal:             nil,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
		{
			ExternalAsset:                 &clptypes.Asset{Symbol: "cusdt"},
			NativeAssetBalance:            sdk.NewUintFromString("1550459183129248235861408"),
			ExternalAssetBalance:          sdk.NewUintFromString("174248776094"),
			PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
			SwapPriceNative:               nil,
			SwapPriceExternal:             nil,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
		{
			ExternalAsset:                 &clptypes.Asset{Symbol: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"},
			NativeAssetBalance:            sdk.NewUintFromString("200501596725333601567765449"),
			ExternalAssetBalance:          sdk.NewUintFromString("708998027178"),
			PoolUnits:                     sdk.NewUintFromString("200501596725333601567765449"),
			SwapPriceNative:               nil,
			SwapPriceExternal:             nil,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
		{
			ExternalAsset:                 &clptypes.Asset{Symbol: "ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D"},
			NativeAssetBalance:            sdk.NewUintFromString("29315228314524379224549414"),
			ExternalAssetBalance:          sdk.NewUintFromString("29441954962"),
			PoolUnits:                     sdk.NewUintFromString("29315228314524379224549414"),
			SwapPriceNative:               nil,
			SwapPriceExternal:             nil,
			RewardPeriodNativeDistributed: sdk.ZeroUint(),
		},
		{
			ExternalAsset:                 &clptypes.Asset{Symbol: "ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA"},
			NativeAssetBalance:            sdk.NewUintFromString("32788415426458039601937058"),
			ExternalAssetBalance:          sdk.NewUintFromString("139140831718"),
			PoolUnits:                     sdk.NewUintFromString("32788415426458039601937058"),
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
			LeverageMax:          sdk.NewUintFromString("1"),
			HealthGainFactor:     sdk.MustNewDecFromStr("1.0"),
			InterestRateMin:      sdk.MustNewDecFromStr("0.005"),
			InterestRateMax:      sdk.MustNewDecFromStr("3.0"),
			InterestRateDecrease: sdk.MustNewDecFromStr("0.10"),
			InterestRateIncrease: sdk.MustNewDecFromStr("0.10"),
			ForceCloseThreshold:  sdk.MustNewDecFromStr("0.10"),
			EpochLength:          1,
			Pools:                []string{"cusdt", "cusdc"},
		},
	}
	bz, err := cfg.Codec.MarshalJSON(gs)
	return bz, err
}
