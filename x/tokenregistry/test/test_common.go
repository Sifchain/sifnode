package test

import (
	sifapp "github.com/Sifchain/sifnode/app"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func CreateTestApp(isCheckTx bool) (*sifapp.SifchainApp, sdk.Context, string) {
	sifapp.SetConfig(false)
	app := sifapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)
	_ = sifapp.AddTestAddrs(app, ctx, 6, initTokens)
	admin := sdk.AccAddress("addr1_______________")
	state := tokenregistrytypes.GenesisState{
		AdminAccounts: GetAdmins(admin.String()),
		Registry:      nil,
	}
	app.TokenRegistryKeeper.InitGenesis(ctx, state)
	return app, ctx, admin.String()
}

func GetAdmins(address string) *tokenregistrytypes.AdminAccounts {
	return &tokenregistrytypes.AdminAccounts{AdminAccounts: []*tokenregistrytypes.AdminAccount{
		{
			AdminType:    tokenregistrytypes.AdminType_TOKENREGISTRY,
			AdminAddress: address,
		},
	}}
}
