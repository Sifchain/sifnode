package cli

import (
	"errors"

	"github.com/Sifchain/sifnode/x/admin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Admin key management transactions sub-commands",
	}
	cmd.AddCommand(
		GetCmdAdd(),
		GetCmdRemove(),
		GetCmdSetParams(),
	)
	return cmd
}

func GetCmdAdd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-account [address] [type]",
		Short: "Add an account",
		Long: `Add an account. Valid types are:
	CLPDEX
	PMTPREWARDS
	TOKENREGISTRY
	ETHBRIDGE
	ADMIN
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(2)(cmd, args)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			adminType, ok := types.AdminType_value[args[1]]
			if !ok {
				return errors.New("invalid admin type")
			}

			msg := types.MsgAddAccount{
				Signer: clientCtx.GetFromAddress().String(),
				Account: &types.AdminAccount{
					AdminType:    types.AdminType(adminType),
					AdminAddress: addr.String(),
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRemove() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-account [address] [type]",
		Short: "Remove an account",
		Long: `Remove an account. Valid types are:
	CLPDEX
	PMTPREWARDS
	TOKENREGISTRY
	ETHBRIDGE
	ADMIN
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(2)(cmd, args)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			adminType, ok := types.AdminType_value[args[1]]
			if !ok {
				return errors.New("invalid admin type")
			}

			msg := types.MsgRemoveAccount{
				Signer: clientCtx.GetFromAddress().String(),
				Account: &types.AdminAccount{
					AdminType:    types.AdminType(adminType),
					AdminAddress: addr.String(),
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdSetParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-params",
		Short: "Set parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fee := sdk.NewUintFromString(viper.GetString("submit-proposal-fee"))

			msg := types.MsgSetParams{
				Signer: clientCtx.GetFromAddress().String(),
				Params: &types.Params{
					SubmitProposalFee: fee,
				},
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().String("submit-proposal-fee", "5000000000000000000000", "fee to submit proposals")
	_ = cmd.MarkFlagRequired("submit-proposal-fee")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
