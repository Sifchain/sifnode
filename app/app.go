package app

import (
	"encoding/json"

	//"github.com/Sifchain/sifnode/x/clp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/tendermint/tendermint/libs/log"
	"net/http"

	tmos "github.com/tendermint/tendermint/libs/os"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"

	//"github.com/Sifchain/sifnode/x/faucet"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	sifappparams "github.com/Sifchain/sifnode/app/params"
	baseapp "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
)

const appName = "sifnode"

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.sifnodecli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.sifnoded")
	ModuleBasics    = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			upgradeclient.ProposalHandler,
		),
		params.AppModuleBasic{},
		//clp.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		//oracle.AppModuleBasic{},
		//ethbridge.AppModuleBasic{},
		//faucet.AppModuleBasic{},
		slashing.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner, authtypes.Staking},
		//ethbridge.ModuleName:      {authtypes.Burner, authtypes.Minter},
		//clp.ModuleName: {authtypes.Burner, authtypes.Minter},
		//faucet.ModuleName:         {authtypes.Minter},
	}
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}
)

//func MakeCodec() *codec.Codec {
//	var cdc = codec.New()
//
//	ModuleBasics.RegisterCodec(cdc)
//	vesting.RegisterCodec(cdc) // Need to verify if we need this
//	sdk.RegisterCodec(cdc)
//	codec.RegisterCrypto(cdc)
//	codec.RegisterEvidences(cdc)
//
//	return cdc.Seal()
//}

type SifchainApp struct {
	*bam.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Marshaler
	interfaceRegistry types.InterfaceRegistry
	invCheckPeriod    uint

	keys    map[string]*sdk.KVStoreKey
	tKeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	AccountKeeper  authkeeper.AccountKeeper
	paramsKeeper   paramskeeper.Keeper
	UpgradeKeeper  upgradekeeper.Keeper
	govKeeper      govkeeper.Keeper
	bankKeeper     bankkeeper.Keeper
	stakingKeeper  stakingkeeper.Keeper
	slashingKeeper slashingkeeper.Keeper
	distrKeeper    distrkeeper.Keeper

	// Sifchain keepers
	//EthBridgeKeeper ethbridge.Keeper
	//OracleKeeper    oracle.Keeper
	//clpKeeper       clp.Keeper
	mm *module.Manager
	//faucetKeeper    faucet.Keeper
	sm *module.SimulationManager
}

func NewSifChainApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool, homePath string,
	invCheckPeriod uint, encodingConfig sifappparams.EncodingConfig, baseAppOptions ...func(*bam.BaseApp),
) *SifchainApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	file, _ := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	l := log.NewTMFmtLogger(file)
	_ = l.Log("base app options : ", baseAppOptions)
	bApp := bam.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		//	bam.MainStoreKey,
		authtypes.StoreKey,
		stakingtypes.StoreKey,
		banktypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		//oracle.StoreKey,
		//ethbridge.StoreKey,
		//clp.StoreKey,
		govtypes.StoreKey,
		//faucet.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
	)

	tKeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	var app = &SifchainApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tKeys:             tKeys,
	}

	app.paramsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])
	bApp.SetParamStore(app.paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)

	app.bankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.ModuleAccountAddrs(),
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.bankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	app.distrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.AccountKeeper,
		app.bankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)

	app.slashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
	app.stakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	//app.OracleKeeper = oracle.NewKeeper(
	//	app.cdc,
	//	keys[oracle.StoreKey],
	//	app.stakingKeeper,
	//	oracle.DefaultConsensusNeeded,
	//)

	//app.EthBridgeKeeper = ethbridge.NewKeeper(
	//	app.cdc,
	//	app.SupplyKeeper,
	//	app.OracleKeeper,
	//	keys[ethbridge.StoreKey],
	//)

	//app.clpKeeper = clp.NewKeeper(
	//	app.cdc,
	//	keys[clp.StoreKey],
	//	app.bankKeeper,
	//	app.SupplyKeeper,
	//	app.subspaces[clp.ModuleName])

	//app.faucetKeeper = faucet.NewKeeper(
	//	app.SupplyKeeper,
	//	app.cdc,
	//	keys[faucet.StoreKey],
	//	app.bankKeeper)

	// This map defines heights to skip for updates
	// The mapping represents height to bool. if the value is true for a height that height
	// will be skipped even if we have a update proposal for it

	//skipUpgradeHeights := make(map[int64]bool)
	//skipUpgradeHeights[0] = true
	//app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath)

	//app.UpgradeKeeper.SetUpgradeHandler("sifUpdate1", func(ctx sdk.Context, plan upgrade.Plan) {
	//	clp.InitGenesis(ctx, app.clpKeeper, clp.DefaultGenesisState())
	//})
	//app.SetStoreLoader(bam.StoreLoaderWithUpgrade(&types.StoreUpgrades{
	//	Added: []string{clp.ModuleName},
	//}))

	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath)

	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	app.govKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.bankKeeper,
		&stakingKeeper,
		govRouter,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		//vesting.NewAppModule(app.AccountKeeper, app.bankKeeper),
		bank.NewAppModule(appCodec, app.bankKeeper, app.AccountKeeper),
		gov.NewAppModule(appCodec, app.govKeeper, app.AccountKeeper, app.bankKeeper),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.AccountKeeper, app.bankKeeper, app.stakingKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.AccountKeeper, app.bankKeeper, app.stakingKeeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.AccountKeeper, app.bankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		params.NewAppModule(app.paramsKeeper),
	)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(distrtypes.ModuleName,
		stakingtypes.ModuleName,
		//faucet.ModuleName,
		upgradetypes.ModuleName)

	app.mm.SetOrderEndBlockers(
		stakingtypes.ModuleName,
		govtypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		slashingtypes.ModuleName,

		genutiltypes.ModuleName,
		//oracle.ModuleName,
		//ethbridge.ModuleName,
		//clp.ModuleName,
		govtypes.ModuleName,
		//faucet.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(module.NewConfigurator(app.MsgServiceRouter(), app.GRPCQueryRouter()))

	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.bankKeeper, app.AccountKeeper),
		gov.NewAppModule(appCodec, app.govKeeper, app.AccountKeeper, app.bankKeeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.AccountKeeper, app.bankKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.AccountKeeper, app.bankKeeper, app.stakingKeeper),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.AccountKeeper, app.bankKeeper, app.stakingKeeper),
		params.NewAppModule(app.paramsKeeper),
	)
	app.sm.RegisterStoreDecoders()

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.SetAnteHandler(
		ante.NewAnteHandler(
			app.AccountKeeper, app.bankKeeper, ante.DefaultSigVerificationGasConsumer,
			encodingConfig.TxConfig.SignModeHandler(),
		),
	)

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)
	//app.MountMemoryStores(memKeys)

	if loadLatest {
		err := app.LoadLatestVersion()
		if err != nil {
			tmos.Exit(err.Error())
		}

	}
	return app
}

type GenesisState map[string]json.RawMessage

//func NewDefaultGenesisState() GenesisState {
//	return ModuleBasics.DefaultGenesis()
//}

func MakeCodecs() (codec.Marshaler, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

func (app *SifchainApp) Name() string { return app.BaseApp.Name() }

func (app *SifchainApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	app.legacyAmino.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *SifchainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}
func (app *SifchainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *SifchainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *SifchainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *SifchainApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

func (app *SifchainApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *SifchainApp) AppCodec() codec.Marshaler {
	return app.appCodec
}

func (app *SifchainApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

func (app *SifchainApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key
func (app *SifchainApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tKeys[storeKey]
}

func (app *SifchainApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *SifchainApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}

func initParamsKeeper(appCodec codec.BinaryMarshaler, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	//app.subspaces[clp.ModuleName] = app.paramsKeeper.Subspace(clp.DefaultParamspace)
	return paramsKeeper
}

func (app *SifchainApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)

	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(apiSvr.ClientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

func (app *SifchainApp) RegisterTxService(clientCtx client.Context) {
	tx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *SifchainApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}
