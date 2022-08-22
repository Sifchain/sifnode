package app

import (
	"github.com/cosmos/cosmos-sdk/x/authz"
	"io"
	"math/big"
	"net/http"
	"os"

	sifchainAnte "github.com/Sifchain/sifnode/app/ante"
	"github.com/Sifchain/sifnode/x/admin"
	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	"github.com/Sifchain/sifnode/x/clp"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/dispensation"
	dispkeeper "github.com/Sifchain/sifnode/x/dispensation/keeper"
	disptypes "github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/Sifchain/sifnode/x/ethbridge"
	ethbridgekeeper "github.com/Sifchain/sifnode/x/ethbridge/keeper"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	ibctransferoverride "github.com/Sifchain/sifnode/x/ibctransfer"
	sctransfertypes "github.com/Sifchain/sifnode/x/ibctransfer/types"
	"github.com/Sifchain/sifnode/x/oracle"
	oraclekeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/Sifchain/sifnode/x/tokenregistry"
	tokenregistrykeeper "github.com/Sifchain/sifnode/x/tokenregistry/keeper"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	feegrant "github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v2/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v2/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v2/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v2/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v2/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v2/modules/core/keeper"
	ibcmock "github.com/cosmos/ibc-go/v2/testing/mock"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"
)

const appName = "sifnode"

var (
	DefaultNodeHome = os.ExpandEnv("$HOME/.sifnoded")

	ModuleBasics = module.NewBasicManager(
		append(
			[]module.AppModuleBasic{
				genutil.AppModuleBasic{},
				auth.AppModuleBasic{},
				bank.AppModuleBasic{},
				capability.AppModuleBasic{},
				crisis.AppModuleBasic{},
				feegrantmodule.AppModuleBasic{},
				authzmodule.AppModuleBasic{},
				staking.AppModuleBasic{},
				mint.AppModuleBasic{},
				distr.AppModuleBasic{},
				gov.NewAppModuleBasic(
					paramsclient.ProposalHandler,
					distrclient.ProposalHandler,
					upgradeclient.ProposalHandler,
					ibcclientclient.UpdateClientProposalHandler,
					ibcclientclient.UpgradeProposalHandler,
				),
				params.AppModuleBasic{},
				upgrade.AppModuleBasic{},
				slashing.AppModuleBasic{},
				evidence.AppModuleBasic{},
				ibc.AppModuleBasic{},
				ibctransferoverride.AppModuleBasic{},

				clp.AppModuleBasic{},
				oracle.AppModuleBasic{},
				ethbridge.AppModuleBasic{},
				dispensation.AppModuleBasic{},
				tokenregistry.AppModuleBasic{},
				admin.AppModuleBasic{},
				vesting.AppModuleBasic{},
			},
			FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getAppModuleBasics()...,
		)...,
	)

	maccPerms = FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getMaccPerms(
		map[string][]string{
			authtypes.FeeCollectorName:     nil,
			distrtypes.ModuleName:          nil,
			minttypes.ModuleName:           {authtypes.Minter},
			stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
			stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
			govtypes.ModuleName:            {authtypes.Burner, authtypes.Staking},
			ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
			sctransfertypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
			ethbridgetypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
			clptypes.ModuleName:            {authtypes.Burner, authtypes.Minter},
			dispensation.ModuleName:        {authtypes.Burner, authtypes.Minter},
		},
	)
)

func init() {
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

	_, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
}

var (
	_ simapp.App              = (*SifchainApp)(nil)
	_ servertypes.Application = (*SifchainApp)(nil)
)

type SifchainApp struct {
	*baseapp.BaseApp
	FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SifchainApp
	legacyAmino       *codec.LegacyAmino //nolint
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	TransferKeeper   ibctransferkeeper.Keeper
	PortKeeper       ibcmock.PortKeeper
	FeegrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCMockKeeper  capabilitykeeper.ScopedKeeper

	ClpKeeper           clpkeeper.Keeper
	OracleKeeper        oraclekeeper.Keeper
	EthbridgeKeeper     ethbridgekeeper.Keeper
	DispensationKeeper  dispkeeper.Keeper
	TokenRegistryKeeper tokenregistrytypes.Keeper
	AdminKeeper         adminkeeper.Keeper

	mm           *module.Manager
	sm           *module.SimulationManager
	configurator module.Configurator
}

func NewSifApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *SifchainApp {
	return NewSifAppWithBlacklist(logger, db, traceStore, loadLatest, skipUpgradeHeights, homePath, invCheckPeriod, encodingConfig, appOpts, []sdk.AccAddress{}, baseAppOptions...)
}

func NewSifAppWithBlacklist(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions, blackListedAddresses []sdk.AccAddress, baseAppOptions ...func(*baseapp.BaseApp),
) *SifchainApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	storeKeys := []string{
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		govtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		evidencetypes.StoreKey,
		ibchost.StoreKey,
		ibctransfertypes.StoreKey,
		feegrant.StoreKey,
		capabilitytypes.StoreKey,
		disptypes.StoreKey,
		ethbridgetypes.StoreKey,
		clptypes.StoreKey,
		oracletypes.StoreKey,
		tokenregistrytypes.StoreKey,
		admintypes.StoreKey,
		authzkeeper.StoreKey,
	}
	storeKeys = append(storeKeys, FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getStoreKeys()...)
	keys := sdk.NewKVStoreKeys(storeKeys...)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	var app = &SifchainApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino, //nolint
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	if homePath == "" {
		homePath = DefaultNodeHome
	}
	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey]) //nolint
	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	scopedIBCMockKeeper := app.CapabilityKeeper.ScopeToModule(ibcmock.ModuleName)

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)

	addrs := app.ModuleAccountAddrs()
	for _, addr := range blackListedAddresses {
		addrs[addr.String()] = true
	}

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		addrs,
	)
	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
	)

	app.FeegrantKeeper = feegrantkeeper.NewKeeper(
		appCodec, keys[feegrant.StoreKey], app.AccountKeeper,
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName), &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	/*
	   Invariants are certain conditions which are checked , which should hold true to certify the network is not byzantine .
	   The invariants are asserted in two places

	   End of every block
	   https://github.com/cosmos/cosmos-sdk/blob/d3a8e1e9535e4d980723fbabdc274cd55cf965ae/x/crisis/abci.go#L20
	   Controlled by `invCheckPeriod`

	   InitGenesis of the Crisis module.
	   https://github.com/cosmos/cosmos-sdk/blob/d3a8e1e9535e4d980723fbabdc274cd55cf965ae/x/crisis/module.go#L147
	   Controlled by `skipGenesisInvariants`

	   We remove the invariant check from crisis module and add to be checked once every 600 blocks instead.
	*/
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)
	skipUpgradeHeights[0] = true
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp)
	app.AdminKeeper = adminkeeper.NewKeeper(appCodec, keys[admintypes.StoreKey])
	app.TokenRegistryKeeper = tokenregistrykeeper.NewKeeper(appCodec, keys[tokenregistrytypes.StoreKey], app.AdminKeeper)

	FEATURE_TOGGLE_MARGIN_CLI_ALPHA_setKeepers(app, keys, &appCodec)
	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		govRouter,
	)
	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register governance hooks
		),
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	transferModule := ibctransferoverride.NewAppModule(app.TransferKeeper, app.TokenRegistryKeeper, app.BankKeeper, appCodec)
	app.OracleKeeper = oraclekeeper.NewKeeper(
		appCodec,
		keys[oracletypes.StoreKey],
		app.StakingKeeper,
		oracletypes.DefaultConsensusNeeded,
	)

	app.EthbridgeKeeper = ethbridgekeeper.NewKeeper(
		appCodec,
		app.BankKeeper,
		app.OracleKeeper,
		app.AccountKeeper,
		app.AdminKeeper,
		keys[ethbridgetypes.StoreKey],
	)

	app.DispensationKeeper = dispkeeper.NewKeeper(
		appCodec,
		keys[disptypes.StoreKey],
		app.BankKeeper,
		app.AccountKeeper,
		app.GetSubspace(disptypes.ModuleName),
	)
	mockModule := ibcmock.NewAppModule(scopedIBCMockKeeper, &app.IBCKeeper.PortKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	ibcRouter.AddRoute(ibcmock.ModuleName, mockModule)
	app.IBCKeeper.SetRouter(ibcRouter)
	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper
	cfg := module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.configurator = cfg

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	appModules := []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper,
			app.StakingKeeper,
			app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeegrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		clp.NewAppModule(app.ClpKeeper, app.BankKeeper),
		oracle.NewAppModule(app.OracleKeeper),
		ethbridge.NewAppModule(app.OracleKeeper, app.BankKeeper, app.AccountKeeper, app.EthbridgeKeeper, &appCodec),
		dispensation.NewAppModule(app.DispensationKeeper, app.BankKeeper, app.AccountKeeper),
		tokenregistry.NewAppModule(app.TokenRegistryKeeper, &appCodec),
		admin.NewAppModule(app.AdminKeeper, &appCodec),
	}
	appModules = append(appModules, FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getAppModules(app, &appCodec)...)
	app.mm = module.NewManager(appModules...)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	orderBeginBlockers := []string{
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		ibchost.ModuleName,
		disptypes.ModuleName,
		transferModule.Name(),
		clptypes.ModuleName,
		ethbridgetypes.ModuleName,
		tokenregistrytypes.ModuleName,
		oracletypes.ModuleName,
		dispensation.ModuleName,
		admintypes.ModuleName,
	}
	orderBeginBlockers = append(orderBeginBlockers, FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderBeginBlockers()...)
	app.mm.SetOrderBeginBlockers(orderBeginBlockers...)

	orderEndBlockers := []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		ibchost.ModuleName,
		disptypes.ModuleName,
		transferModule.Name(),
		clptypes.ModuleName,
		ethbridgetypes.ModuleName,
		tokenregistrytypes.ModuleName,
		oracletypes.ModuleName,
		admintypes.ModuleName,
	}
	orderEndBlockers = append(orderEndBlockers, FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderEndBlockers()...)
	app.mm.SetOrderEndBlockers(orderEndBlockers...)
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	orderInitGenesis := []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		ibchost.ModuleName,
		disptypes.ModuleName,
		transferModule.Name(),
		tokenregistrytypes.ModuleName,
		admintypes.ModuleName,
		clptypes.ModuleName,
		ethbridgetypes.ModuleName,
		oracletypes.ModuleName,
		ethbridge.ModuleName,
		dispensation.ModuleName,
	}
	orderInitGenesis = append(orderInitGenesis, FEATURE_TOGGLE_MARGIN_CLI_ALPHA_getOrderInitGenesis()...)
	app.mm.SetOrderInitGenesis(orderInitGenesis...)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(app.configurator)
	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)
	anteHandler, err := sifchainAnte.NewAnteHandler(
		sifchainAnte.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			StakingKeeper:   app.StakingKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  app.FeegrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	SetupHandlers(app)
	if loadLatest {
		err := app.LoadLatestVersion()
		if err != nil {
			tmos.Exit(err.Error())
		}
		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		app.CapabilityKeeper.Seal()
	}
	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	app.ScopedIBCMockKeeper = scopedIBCMockKeeper
	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	return app
}

func (app *SifchainApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	// Init chain is called when the chain starts up
	// Setting the module version here finalizes the version map which makes sure we can keep track of migrations and not migrate from 1 to 2 every time.
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *SifchainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *SifchainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *SifchainApp) LegacyAmino() *codec.LegacyAmino { //nolint
	return app.legacyAmino
}

// AppCodec returns SimApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SifchainApp) AppCodec() codec.Codec {
	return app.appCodec
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SifchainApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *SifchainApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SifchainApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *SifchainApp) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// LoadHeight loads a particular height
func (app *SifchainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SifchainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
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

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *SifchainApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *SifchainApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *SifchainApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper { //nolint
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(clptypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	FEATURE_TOGGLE_MARGIN_CLI_ALPHA_setParamsKeeper(&paramsKeeper)
	return paramsKeeper
}
