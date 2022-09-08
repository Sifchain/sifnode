package cli

import (
	"errors"

	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func GetAdminCloseAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin-close-all",
		Short: "Force close margin position",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			takeMarginFund, err := cmd.Flags().GetBool("take_margin_fund")
			if err != nil {
				return err
			}

			msg := types.MsgAdminCloseAll{
				Signer:         signer.String(),
				TakeMarginFund: takeMarginFund,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().Bool("take_margin_fund", true, "boolean value to indicate weather margin fund will be deducted on close")
	_ = cmd.MarkFlagRequired("take_margin_fund")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetAdminCloseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin-close",
		Short: "Force close margin position",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			MtpAddress, err := cmd.Flags().GetString("mtp_address")
			if err != nil {
				return err
			}

			id, err := cmd.Flags().GetUint64("id")
			if err != nil {
				return err
			}

			takeMarginFund, err := cmd.Flags().GetBool("take_margin_fund")
			if err != nil {
				return err
			}

			msg := types.MsgAdminClose{
				Signer:         signer.String(),
				MtpAddress:     MtpAddress,
				Id:             id,
				TakeMarginFund: takeMarginFund,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().String("mtp_address", "", "mtp address")
	cmd.Flags().Uint64("id", 0, "id of the position")
	_ = cmd.MarkFlagRequired("mtp_address")
	_ = cmd.MarkFlagRequired("id")
	cmd.Flags().Bool("take_margin_fund", true, "boolean value to indicate weather margin fund will be deducted on close")
	_ = cmd.MarkFlagRequired("take_margin_fund")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
