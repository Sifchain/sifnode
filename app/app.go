package app

import (
	"encoding/json"
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"

	tmos "github.com/tendermint/tendermint/libs/os"

	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/x/gov"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/Sifchain/sifnode/x/ethbridge"
	"github.com/Sifchain/sifnode/x/faucet"
	"github.com/Sifchain/sifnode/x/oracle"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
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
			paramsclient.ProposalHandler,
		),
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		clp.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		oracle.AppModuleBasic{},
		ethbridge.AppModuleBasic{},
		faucet.AppModuleBasic{},
		slashing.AppModuleBasic{},
		dispensation.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner, supply.Staking},
		ethbridge.ModuleName:      {supply.Burner, supply.Minter},
		clp.ModuleName:            {supply.Burner, supply.Minter},
		dispensation.ModuleName:   {supply.Burner, supply.Minter},
		faucet.ModuleName:         {supply.Minter},
	}
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc) // Need to verify if we need this
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc.Seal()
}

func init() {
	sdk.PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

type SifchainApp struct {
	*bam.BaseApp
	Cdc *codec.Codec

	invCheckPeriod uint

	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	subspaces map[string]params.Subspace

	AccountKeeper      auth.AccountKeeper
	paramsKeeper       params.Keeper
	UpgradeKeeper      upgrade.Keeper
	GovKeeper          gov.Keeper
	BankKeeper         bank.Keeper
	StakingKeeper      staking.Keeper
	SlashingKeeper     slashing.Keeper
	DistributionKeeper distr.Keeper
	SupplyKeeper       supply.Keeper

	// Peggy keepers
	EthBridgeKeeper    ethbridge.Keeper
	OracleKeeper       oracle.Keeper
	ClpKeeper          clp.Keeper
	DispensationKeeper dispensation.Keeper
	mm                 *module.Manager
	FaucetKeeper       faucet.Keeper
	sm                 *module.SimulationManager
}

var _ simapp.App = (*SifchainApp)(nil)

func NewInitApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *SifchainApp {

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
		clp.StoreKey,
		gov.StoreKey,
		faucet.StoreKey,
		distr.StoreKey,
		slashing.StoreKey,
		dispensation.StoreKey,
	)

	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	var app = &SifchainApp{
		BaseApp:        bApp,
		Cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	app.paramsKeeper = params.NewKeeper(app.Cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[clp.ModuleName] = app.paramsKeeper.Subspace(clp.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[dispensation.ModuleName] = app.paramsKeeper.Subspace(dispensation.DefaultParamspace)

	app.AccountKeeper = auth.NewAccountKeeper(
		app.Cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	app.BankKeeper = bank.NewBaseKeeper(
		app.AccountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	app.SupplyKeeper = supply.NewKeeper(
		app.Cdc,
		keys[supply.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		maccPerms,
	)

	stakingKeeper := staking.NewKeeper(
		app.Cdc,
		keys[staking.StoreKey],
		app.SupplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	app.DistributionKeeper = distr.NewKeeper(app.Cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName], &stakingKeeper,
		app.SupplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())

	app.SlashingKeeper = slashing.NewKeeper(
		app.Cdc, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)

	app.StakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.DistributionKeeper.Hooks(), app.SlashingKeeper.Hooks()))

	app.OracleKeeper = oracle.NewKeeper(
		app.Cdc,
		keys[oracle.StoreKey],
		app.StakingKeeper,
		oracle.DefaultConsensusNeeded,
	)

	app.EthBridgeKeeper = ethbridge.NewKeeper(
		app.Cdc,
		app.SupplyKeeper,
		app.OracleKeeper,
		keys[ethbridge.StoreKey],
	)

	app.ClpKeeper = clp.NewKeeper(
		app.Cdc,
		keys[clp.StoreKey],
		app.BankKeeper,
		app.SupplyKeeper,
		app.subspaces[clp.ModuleName])

	app.DispensationKeeper = dispensation.NewKeeper(
		app.Cdc,
		keys[dispensation.StoreKey],
		app.BankKeeper,
		app.SupplyKeeper,
	)

	app.FaucetKeeper = faucet.NewKeeper(
		app.SupplyKeeper,
		app.Cdc,
		keys[faucet.StoreKey],
		app.BankKeeper)

	// This map defines heights to skip for updates
	// The mapping represents height to bool. if the value is true for a height that height
	// will be skipped even if we have a update proposal for it

	skipUpgradeHeights := make(map[int64]bool)
	skipUpgradeHeights[0] = true
	app.UpgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], app.Cdc)
	app.paramsKeeper = params.NewKeeper(app.Cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	SetupHandlers(app)

	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper))
	app.GovKeeper = gov.NewKeeper(
		app.Cdc,
		keys[gov.StoreKey],
		app.subspaces[gov.ModuleName],
		app.SupplyKeeper,
		app.StakingKeeper,
		govRouter,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		distr.NewAppModule(app.DistributionKeeper, app.AccountKeeper, app.SupplyKeeper, app.StakingKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.StakingKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		oracle.NewAppModule(app.OracleKeeper),
		ethbridge.NewAppModule(app.OracleKeeper, app.SupplyKeeper, app.AccountKeeper, app.EthBridgeKeeper, app.Cdc),
		clp.NewAppModule(app.ClpKeeper, app.BankKeeper, app.SupplyKeeper),
		faucet.NewAppModule(app.FaucetKeeper, app.BankKeeper, app.SupplyKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.SupplyKeeper),
		dispensation.NewAppModule(app.DispensationKeeper, app.BankKeeper, app.SupplyKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(distr.ModuleName,
		slashing.ModuleName,
		faucet.ModuleName,
		dispensation.ModuleName,
		upgrade.ModuleName)

	app.mm.SetOrderEndBlockers(
		staking.ModuleName,
		gov.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
		oracle.ModuleName,
		ethbridge.ModuleName,
		clp.ModuleName,
		gov.ModuleName,
		faucet.ModuleName,
		dispensation.ModuleName,
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

func (app *SifchainApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	app.Cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *SifchainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *SifchainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *SifchainApp) Codec() *codec.Codec {
	return app.Cdc
}

func (app *SifchainApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key
func (app *SifchainApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tKeys[storeKey]
}

func (app *SifchainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

func (app *SifchainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
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
