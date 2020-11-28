package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	tmos "github.com/tendermint/tendermint/libs/os"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/x/gov"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/Sifchain/sifnode/x/ethbridge"
	"github.com/Sifchain/sifnode/x/oracle"

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
		gov.NewAppModuleBasic(
			upgradeclient.ProposalHandler,
		),
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		//	clp.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		oracle.AppModuleBasic{},
		ethbridge.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner, supply.Staking},
		ethbridge.ModuleName:      {supply.Burner, supply.Minter},
		//	clp.ModuleName:            {supply.Burner, supply.Minter},
	}
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc) // Need to verify if we need this
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc.Seal()
}

type NewApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	subspaces map[string]params.Subspace

	AccountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	StakingKeeper staking.Keeper
	SupplyKeeper  supply.Keeper
	paramsKeeper  params.Keeper
	UpgradeKeeper upgrade.Keeper
	govKeeper     gov.Keeper

	// Peggy keepers
	EthBridgeKeeper ethbridge.Keeper
	OracleKeeper    oracle.Keeper
	//	clpKeeper clp.Keeper
	mm *module.Manager

	sm *module.SimulationManager
}

var _ simapp.App = (*NewApp)(nil)

func NewInitApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *NewApp {

	f, _ := os.Create("testlog.log")
	defer f.Close()
	loger := log.NewTMLogger(f)
	loger.Info("Starting to setup app ")

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		auth.StoreKey,
		staking.StoreKey,
		supply.StoreKey,
		params.StoreKey,
		upgrade.StoreKey,
		oracle.StoreKey,
		ethbridge.StoreKey,
		//		clp.StoreKey,
		gov.StoreKey,
	)

	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	var app = &NewApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	//	app.subspaces[clp.ModuleName] = app.paramsKeeper.Subspace(clp.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())

	app.AccountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(
		app.AccountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	app.SupplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.AccountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	app.StakingKeeper = staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.SupplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	app.StakingKeeper = *app.StakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(),
	)

	app.OracleKeeper = oracle.NewKeeper(
		app.cdc,
		keys[oracle.StoreKey],
		app.StakingKeeper,
		oracle.DefaultConsensusNeeded,
	)

	app.EthBridgeKeeper = ethbridge.NewKeeper(
		app.cdc,
		app.SupplyKeeper,
		app.OracleKeeper,
	)

	//app.clpKeeper = clp.NewKeeper(
	//	app.cdc,
	//	keys[clp.StoreKey],
	//	app.bankKeeper,
	//	app.SupplyKeeper,
	//	app.subspaces[clp.ModuleName])

	skipUpgradeHeights := make(map[int64]bool)
	skipUpgradeHeights[0] = true
	app.UpgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], app.cdc)
	loger.Info("Trying to upgrade ")
	//app.UpgradeKeeper.SetUpgradeHandler("testupgrade", func(ctx sdk.Context, plan upgrade.Plan) {
	//
	//	f, err := os.Create("testlog.log")
	//	defer f.Close()
	//	loger := log.NewTMLogger(f)
	//	loger.Info("Starting to execute upgrade plan")
	//	ethAsset := clp.NewAsset("Ethereum","ETH","ceth")
	//	loger.Info("Asset Created")
	//	pool,err := clp.NewPool(ethAsset,sdk.NewUint(100),sdk.NewUint(100),sdk.NewUint(10))
	//	if err!= nil{
	//		loger.Info("Pool Not Created" ,err.Error())
	//		return
	//	}
	//	loger.Info("Pool Created")
	//	err = app.clpKeeper.SetPool(ctx,pool)
	//	if err!= nil{
	//		loger.Info("Pool Not Set" ,err.Error())
	//		return
	//	}
	//	loger.Info("Pool Set")
	//})

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))

	app.govKeeper = gov.NewKeeper(
		app.cdc,
		keys[gov.StoreKey],
		app.subspaces[gov.ModuleName],
		app.SupplyKeeper,
		app.StakingKeeper,
		govRouter,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.bankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		oracle.NewAppModule(app.OracleKeeper),
		ethbridge.NewAppModule(app.OracleKeeper, app.SupplyKeeper, app.AccountKeeper, app.EthBridgeKeeper, app.cdc),
		//		clp.NewAppModule(app.clpKeeper, app.bankKeeper, app.SupplyKeeper),
		gov.NewAppModule(app.govKeeper, app.AccountKeeper, app.SupplyKeeper),
	)

	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		upgrade.ModuleName,
		staking.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		staking.ModuleName,
		gov.ModuleName,
	)

	app.mm.SetOrderInitGenesis(
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
		oracle.ModuleName,
		ethbridge.ModuleName,
		//		clp.ModuleName,
		gov.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.AccountKeeper,
			app.SupplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *NewApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState

	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *NewApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *NewApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *NewApp) Codec() *codec.Codec {
	return app.cdc
}

func (app *NewApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key
func (app *NewApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tKeys[storeKey]
}

func (app *NewApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

func (app *NewApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *NewApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}
