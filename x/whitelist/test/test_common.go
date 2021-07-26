package test

import (
	sifapp "github.com/Sifchain/sifnode/app"
	whitelisttypes "github.com/Sifchain/sifnode/x/whitelist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func CreateTestApp(isCheckTx bool) (*sifapp.SifchainApp, sdk.Context, string) {
	app := sifapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	initTokens := sdk.TokensFromConsensusPower(1000)
	app.BankKeeper.SetSupply(ctx, types.NewSupply(sdk.Coins{}))
	_ = sifapp.AddTestAddrs(app, ctx, 6, initTokens)
	admin := sdk.AccAddress("addr1_______________")
	state := whitelisttypes.GenesisState{
		AdminAccount: admin.String(),
		Whitelist:    nil,
	}
	app.WhitelistKeeper.InitGenesis(ctx, state)
	return app, ctx, admin.String()
}
