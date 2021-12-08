package cli

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Margin transactions sub-commands",
	}
	cmd.AddCommand(
		GetOpenLongCmd(),
	)
	return cmd
}

func GetOpenLongCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-long",
		Short: "Open long position",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateralAsset, err := cmd.Flags().GetString("collateral_asset")
			if err != nil {
				return err
			}

			collateralAmount, err := cmd.Flags().GetUint64("collateral_amount")
			if err != nil {
				return err
			}

			borrowAsset, err := cmd.Flags().GetString("borrow_asset")
			if err != nil {
				return err
			}

			msg := types.MsgOpenLong{
				Signer:           clientCtx.GetFromAddress().String(),
				CollateralAsset:  collateralAsset,
				CollateralAmount: sdk.NewUint(collateralAmount),
				BorrowAsset:      borrowAsset,
			}

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().Uint64("collateral_amount", 0, "amount of collateral asset < max_uint64")
	cmd.Flags().String("collateral_asset", "", "symbol of asset")
	cmd.Flags().String("borrow_asset", "", "symbol of asset")
	flags.AddTxFlagsToCmd(cmd)
	return cmd

}
