package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/Sifchain/sifnode/app"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/siflogger"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"

	cfg "github.com/tendermint/tendermint/config"
)

const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

// If a new config is created, change some of the default tendermint settings
func interceptLoadConfig() (conf *cfg.Config, err error) {
	tmpConf := cfg.DefaultConfig()
	err = viper.Unmarshal(tmpConf)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}
	rootDir := tmpConf.RootDir
	configFilePath := filepath.Join(rootDir, "config/config.toml")
	// Intercept only if the file doesn't already exist

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// the following parse config is needed to create directories
		conf, _ = tcmd.ParseConfig() // NOTE: ParseConfig() creates dir/files as necessary.
		conf.ProfListenAddress = "localhost:6060"
		conf.P2P.RecvRate = 5120000
		conf.P2P.SendRate = 5120000
		conf.TxIndex.IndexAllKeys = true
		conf.Consensus.TimeoutCommit = 5 * time.Second
		cfg.WriteConfigFile(configFilePath, conf)
		// Fall through, just so that its parsed into memory.
	}

	if conf == nil {
		conf, err = tcmd.ParseConfig() // NOTE: ParseConfig() creates dir/files as necessary.
		if err != nil {
			panic(err)
		}
	}

	appConfigFilePath := filepath.Join(rootDir, "config/app.toml")
	if _, err := os.Stat(appConfigFilePath); os.IsNotExist(err) {
		appConf, _ := config.ParseConfig()
		config.WriteConfigFile(appConfigFilePath, appConf)
	}

	viper.SetConfigName("app")
	err = viper.MergeInConfig()

	return conf, err
}

func main() {
	cdc := app.MakeCodec()

	app.SetConfig()

	logger := siflogger.New(siflogger.JSON)
	ctx := server.NewContext(cfg.DefaultConfig(), logger)
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "sifnoded",
		Short:             "app Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.MigrateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
	)
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(AddGenesisCLPAdminCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(flags.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(debug.Cmd(cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "AU", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	var cache sdk.MultiStorePersistentCache

	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}
	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}
	return app.NewInitApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		aApp := app.NewInitApp(logger, db, traceStore, false, uint(1))
		err := aApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	aApp := app.NewInitApp(logger, db, traceStore, true, uint(1))

	return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
