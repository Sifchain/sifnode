package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sifapp "github.com/Sifchain/sifnode/app"
)

func CreateTestApp(isCheckTx bool) (*sifapp.SifchainApp, sdk.Context) {
	app := sifapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction)
	_ = sifapp.AddTestAddrs(app, ctx, 6, initTokens)
	return app, ctx
}

func CreateTestAppMargin(isCheckTx bool) (sdk.Context, *sifapp.SifchainApp) {
	sifapp.SetConfig((false))
	app := sifapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	return ctx, app
}

func CreateTestAppMarginFromGenesis(isCheckTx bool, genesisTransformer func(*sifapp.SifchainApp, sifapp.GenesisState) sifapp.GenesisState) (sdk.Context, *sifapp.SifchainApp) {
	sifapp.SetConfig(false)
	app := sifapp.SetupFromGenesis(isCheckTx, genesisTransformer)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	return ctx, app
}
