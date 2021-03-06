package app

import (
	"encoding/json"
	"fmt"
	"github.com/Sifchain/sifnode/x/clp"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/tendermint/tendermint/libs/log"
	"io/ioutil"
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
		supply.AppModuleBasic{},
		clp.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		oracle.AppModuleBasic{},
		ethbridge.AppModuleBasic{},
		faucet.AppModuleBasic{},
		slashing.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner, supply.Staking},
		ethbridge.ModuleName:      {supply.Burner, supply.Minter},
		clp.ModuleName:            {supply.Burner, supply.Minter},
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
	cdc *codec.Codec

	invCheckPeriod uint

	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	subspaces map[string]params.Subspace

	AccountKeeper  auth.AccountKeeper
	paramsKeeper   params.Keeper
	UpgradeKeeper  upgrade.Keeper
	govKeeper      gov.Keeper
	bankKeeper     bank.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	distrKeeper    distr.Keeper
	SupplyKeeper   supply.Keeper

	// Peggy keepers
	EthBridgeKeeper ethbridge.Keeper
	OracleKeeper    oracle.Keeper
	clpKeeper       clp.Keeper
	mm              *module.Manager
	faucetKeeper    faucet.Keeper
	sm              *module.SimulationManager
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
	)

	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	var app = &SifchainApp{
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
	app.subspaces[clp.ModuleName] = app.paramsKeeper.Subspace(clp.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)

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

	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.SupplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	app.distrKeeper = distr.NewKeeper(app.cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName], &stakingKeeper,
		app.SupplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)

	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()))

	app.OracleKeeper = oracle.NewKeeper(
		app.cdc,
		keys[oracle.StoreKey],
		app.stakingKeeper,
		oracle.DefaultConsensusNeeded,
	)

	app.EthBridgeKeeper = ethbridge.NewKeeper(
		app.cdc,
		app.SupplyKeeper,
		app.OracleKeeper,
		keys[ethbridge.StoreKey],
	)

	app.clpKeeper = clp.NewKeeper(
		app.cdc,
		keys[clp.StoreKey],
		app.bankKeeper,
		app.SupplyKeeper,
		app.subspaces[clp.ModuleName])

	app.faucetKeeper = faucet.NewKeeper(
		app.SupplyKeeper,
		app.cdc,
		keys[faucet.StoreKey],
		app.bankKeeper)

	// This map defines heights to skip for updates
	// The mapping represents height to bool. if the value is true for a height that height
	// will be skipped even if we have a update proposal for it

	skipUpgradeHeights := make(map[int64]bool)
	skipUpgradeHeights[0] = true
	app.UpgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], app.cdc)

	app.UpgradeKeeper.SetUpgradeHandler("changePoolFormula", func(ctx sdk.Context, plan upgrade.Plan) {
		ctx.Logger().Info("Starting to execute upgrade plan for pool re-balance")
		appState, vallist, err := app.ExportAppStateAndValidators(true, []string{})
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to export app state: %s", err))
			return
		}
		appStateJSON, err := cdc.MarshalJSON(appState)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to marshal application genesis state: %s", err.Error()))
			return
		}
		valList, err := json.MarshalIndent(vallist, "", " ")
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to marshal application genesis state: %s", err.Error()))
		}

		err = ioutil.WriteFile("State-Export.json", appStateJSON, 0600)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to write state to file: %s", err.Error()))
		}
		err = ioutil.WriteFile("Validator-Export.json", valList, 0600)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to write Validator List to file: %s", err.Error()))
		}

		allPools := app.clpKeeper.GetPools(ctx)
		lps := clp.LiquidityProviders{}
		poolList := clp.Pools{}
		hasError := false
		for _, pool := range allPools {
			lpList := app.clpKeeper.GetLiquidityProvidersForAsset(ctx, pool.ExternalAsset)
			temp := sdk.ZeroUint()
			tempExternal := sdk.ZeroUint()
			tempNative := sdk.ZeroUint()
			for _, lp := range lpList {
				withdrawNativeAssetAmount, withdrawExternalAssetAmount, _, _ := clp.CalculateWithdrawal(pool.PoolUnits, pool.NativeAssetBalance.String(),
					pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(), sdk.NewUint(clp.MaxWbasis).String(), sdk.NewInt(0))
				newLpUnits, lpUnits, err := clp.CalculatePoolUnits(pool.ExternalAsset.Symbol, temp, tempNative, tempExternal,
					withdrawNativeAssetAmount, withdrawExternalAssetAmount)
				if err != nil {
					hasError = true
					ctx.Logger().Error(fmt.Sprintf("failed to calculate pool units for | Pool : %s | LP %s ", pool.String(), lp.String()))
					break
				}
				lp.LiquidityProviderUnits = lpUnits
				if !lp.Validate() {
					hasError = true
					ctx.Logger().Error(fmt.Sprintf("Invalid | LP %s ", lp.String()))
					break
				}
				lps = append(lps, lp)
				tempExternal = tempExternal.Add(withdrawExternalAssetAmount)
				tempNative = tempNative.Add(withdrawNativeAssetAmount)
				temp = newLpUnits
			}
			pool.PoolUnits = temp
			if !app.clpKeeper.ValidatePool(pool) {
				hasError = true
				ctx.Logger().Error(fmt.Sprintf("Invalid | Pool %s ", pool.String()))
				break
			}
			poolList = append(poolList, pool)
		}
		// If we have error dont set state
		if hasError {
			ctx.Logger().Error("Failed to execute upgrade plan for pool re-balance")
		}
		// If we have no errors , Set state .
		if !hasError {
			for _, pool := range poolList {
				_ = app.clpKeeper.SetPool(ctx, pool)
			}
			for _, l := range lps {
				app.clpKeeper.SetLiquidityProvider(ctx, l)
			}
		}
	})

	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		keys[gov.StoreKey],
		app.subspaces[gov.ModuleName],
		app.SupplyKeeper,
		app.stakingKeeper,
		govRouter,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.bankKeeper, app.AccountKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		distr.NewAppModule(app.distrKeeper, app.AccountKeeper, app.SupplyKeeper, app.stakingKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.AccountKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.AccountKeeper, app.SupplyKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		oracle.NewAppModule(app.OracleKeeper),
		ethbridge.NewAppModule(app.OracleKeeper, app.SupplyKeeper, app.AccountKeeper, app.EthBridgeKeeper, app.cdc),
		clp.NewAppModule(app.clpKeeper, app.bankKeeper, app.SupplyKeeper),
		faucet.NewAppModule(app.faucetKeeper, app.bankKeeper, app.SupplyKeeper),
		gov.NewAppModule(app.govKeeper, app.AccountKeeper, app.SupplyKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(distr.ModuleName,
		slashing.ModuleName,
		faucet.ModuleName,
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

	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *SifchainApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *SifchainApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *SifchainApp) Codec() *codec.Codec {
	return app.cdc
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
