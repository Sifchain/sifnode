//go:build !FEATURE_TOGGLE_SDK_045
// +build !FEATURE_TOGGLE_SDK_045

package cmd

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"

	"github.com/Sifchain/sifnode/app"
)

func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := app.MakeTestEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")
	rootCmd := &cobra.Command{
		Use:   "sifnoded",
		Short: "app Daemon (server)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx = client.ReadHomeFlag(initClientCtx, cmd)
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}
			return server.InterceptConfigsPreRunHandler(cmd, "", nil)
		},
	}
	initRootCmd(rootCmd, encodingConfig)
	return rootCmd, encodingConfig
}
