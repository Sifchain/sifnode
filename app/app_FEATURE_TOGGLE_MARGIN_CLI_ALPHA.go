//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package app

import (
	"github.com/Sifchain/sifnode/x/admin"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin"
	marginkeeper "github.com/Sifchain/sifnode/x/margin/keeper"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
)

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getAppModuleBasics() []module.AppModuleBasic {
	return []module.AppModuleBasic{
		margin.AppModuleBasic{},
	}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getMaccPerms(maccPerms map[string][]string) map[string][]string {
	maccPerms[margintypes.ModuleName] = []string{authtypes.Burner, authtypes.Minter}
	return maccPerms
}

type FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SifchainApp struct {
	MarginKeeper marginkeeper.Keeper
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getStoreKeys() []string {
	return []string{
		margintypes.StoreKey,
	}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_setKeepers(app *SifchainApp, keys map[string]*types.KVStoreKey, appCodec *codec.Codec) {
	app.ClpKeeper = clpkeeper.NewKeeper(
		*appCodec,
		keys[clptypes.StoreKey],
		app.BankKeeper,
		app.AccountKeeper,
		app.TokenRegistryKeeper,
		app.AdminKeeper,
		app.MintKeeper,
		func() margintypes.Keeper { return app.MarginKeeper },
		app.GetSubspace(clptypes.ModuleName),
	)
	app.MarginKeeper = marginkeeper.NewKeeper(
		keys[margintypes.StoreKey],
		*appCodec,
		app.BankKeeper,
		app.ClpKeeper,
		app.AdminKeeper,
		app.GetSubspace(margintypes.ModuleName),
	)
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getAppModules(app *SifchainApp, appCodec *codec.Codec) []module.AppModule {
	return []module.AppModule{
		margin.NewAppModule(app.MarginKeeper, appCodec),
	}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderBeginBlockers() []string {
	return []string{
		margin.ModuleName,
		admin.ModuleName,
	}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderEndBlockers() []string {
	return []string{
		margin.ModuleName,
		admin.ModuleName,
	}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderInitGenesis() []string {
	return []string{
		margin.ModuleName,
	}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_setParamsKeeper(paramsKeeper *paramskeeper.Keeper) {
	paramsKeeper.Subspace(margintypes.ModuleName).WithKeyTable(margintypes.ParamKeyTable())
}
