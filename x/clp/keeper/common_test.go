package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
)

// createTestInput Returns a simapp with custom StakingKeeper
// to avoid messing with the hooks.
func createTestInput() (*codec.LegacyAmino, *sifapp.SifchainApp, sdk.Context) {
	app := sifapp.Setup(false)
	sifapp.SetConfig(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.ClpKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.BankKeeper,
		app.AccountKeeper,
		app.GetSubspace(types.ModuleName),
	)
	return app.LegacyAmino(), app, ctx
}
