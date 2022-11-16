package main

import (
	"fmt"
	"os"

	"github.com/Sifchain/sifnode/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	encodingConfig := app.MakeTestEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")
	app.SetConfig(false)

	rootCmd := &cobra.Command{
		Use: "siftest",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
		//RunE: run,
	}

	rootCmd.AddCommand(GetVerifyCmd(), GetTestCmd())

	err := svrcmd.Execute(rootCmd, app.DefaultNodeHome)
	if err != nil {
		panic(err)
	}
}

func GetTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run tests",
		RunE:  runTest,
	}
	flags.AddTxFlagsToCmd(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")
	return cmd
}

func GetVerifyCmd() *cobra.Command {
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify transaction results",
	}

	verifyCmd.AddCommand(GetVerifyRemove(), GetVerifyAdd(), GetVerifyOpen(), GetVerifyClose())

	return verifyCmd
}

/* VerifySwap verifies amounts sent and received from wallet address.
 */
func VerifySwap(clientCtx client.Context, key keyring.Info) {

}

func GetVerifyOpen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-position --height --from --collateralAmount --leverage --collateral-asset --borrow-asset",
		Short: "Verify a margin long position open",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("verifying open...\n")
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			collateralAmount := sdk.NewUintFromString(viper.GetString("collateralAmount"))
			leverageDec, err := sdk.NewDecFromStr(viper.GetString("leverage"))
			if err != nil {
				panic(err)
			}

			err = VerifyOpenLong(clientCtx,
				viper.GetString("from"),
				int64(viper.GetUint64("height")),
				collateralAmount,
				viper.GetString("collateral-asset"),
				viper.GetString("borrow-asset"),
				leverageDec)
			if err != nil {
				panic(err)
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	//cmd.Flags().Uint64("height", 0, "height of transaction")
	cmd.Flags().String("from", "", "address of transactor")
	cmd.Flags().String("collateralAmount", "0", "collateral provided")
	cmd.Flags().String("leverage", "0", "leverage")
	cmd.Flags().String("collateral-asset", "", "collateral asset")
	cmd.Flags().String("borrow-asset", "", "borrow asset")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("collateralAmount")
	_ = cmd.MarkFlagRequired("leverage")
	_ = cmd.MarkFlagRequired("collateral-asset")
	_ = cmd.MarkFlagRequired("height")
	return cmd
}

func VerifyOpenLong(clientCtx client.Context,
	from string,
	height int64,
	collateralAmount sdk.Uint,
	collateralAsset,
	borrowAsset string,
	leverage sdk.Dec) error {

	return nil
}
