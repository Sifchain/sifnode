//go:build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build !FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package app

import (
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
)

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getAppModuleBasics() []module.AppModuleBasic {
	return []module.AppModuleBasic{}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getMaccPerms(maccPerms map[string][]string) map[string][]string {
	return maccPerms
}

type FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SifchainApp struct {
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getStoreKeys() []string {
	return []string{}
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
		app.GetSubspace(clptypes.ModuleName),
	)
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getAppModules(app *SifchainApp, appCodec *codec.Codec) []module.AppModule {
	return []module.AppModule{}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderBeginBlockers() []string {
	return []string{}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderEndBlockers() []string {
	return []string{}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderInitGenesis() []string {
	return []string{}
}

func FEATURE_TOGGLE_MARGIN_CLI_ALPHA_setParamsKeeper(paramsKeeper *paramskeeper.Keeper) {
}
