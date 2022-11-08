package cli

import (
	"errors"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	whitelistutils "github.com/Sifchain/sifnode/x/tokenregistry/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Token registry transactions sub-commands",
	}
	cmd.AddCommand(
		GetCmdRegister(),
		GetCmdDeregister(),
		GetCmdRegisterAll(),
		GetCmdDeregisterAll(),
		GetCmdSetRegistry(),
	)
	return cmd
}

func GetCmdRegisterAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-all [registry.json]",
		Short: "Add / update tokens on the registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(1)(cmd, args)
			if err != nil {
				return err
			}
			registry, err := whitelistutils.ParseDenoms(clientCtx.Codec, args[0])
			if err != nil {
				return err
			} else if len(registry.Entries) < 1 {
				return errors.New("at least one token entry must be specified in input file")
			}
			for _, entry := range registry.Entries {
				msg := types.MsgRegister{
					From:  clientCtx.GetFromAddress().String(),
					Entry: entry,
				}
				if err := msg.ValidateBasic(); err != nil {
					return err
				}
				err := tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdRegister() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [registry.json]",
		Short: "Add / update token on the registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(1)(cmd, args)
			if err != nil {
				return err
			}
			registry, err := whitelistutils.ParseDenoms(clientCtx.Codec, args[0])
			if err != nil {
				return err
			} else if len(registry.Entries) != 1 {
				return errors.New("exactly one token entry must be specified in input file")
			}
			msg := types.MsgRegister{
				From:  clientCtx.GetFromAddress().String(),
				Entry: registry.Entries[0],
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

func GetCmdSetRegistry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-registry [registry.json]",
		Short: "Reset the registry to the registry.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(1)(cmd, args)
			if err != nil {
				return err
			}
			registry, err := whitelistutils.ParseDenoms(clientCtx.Codec, args[0])
			if err != nil {
				return err
			} else if len(registry.Entries) != 1 {
				return errors.New("exactly one token entry must be specified in input file")
			}
			msg := types.MsgSetRegistry{
				From:     clientCtx.GetFromAddress().String(),
				Registry: &registry,
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

func GetCmdDeregister() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister [denom]",
		Short: "Remove token from the registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(1)(cmd, args)
			if err != nil {
				return err
			}
			msg := types.MsgDeregister{
				From:  clientCtx.GetFromAddress().String(),
				Denom: args[0],
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

func GetCmdDeregisterAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister-all [registry.json]",
		Short: "Remove all tokens listed in registry.json from the registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			err = cobra.ExactArgs(1)(cmd, args)
			if err != nil {
				return err
			}
			registry, err := whitelistutils.ParseDenoms(clientCtx.Codec, args[0])
			if err != nil {
				return err
			} else if len(registry.Entries) < 1 {
				return errors.New("at least one token entry must be specified in input file")
			}
			for _, entry := range registry.Entries {
				msg := types.MsgDeregister{
					From:  clientCtx.GetFromAddress().String(),
					Denom: entry.Denom,
				}
				if err := msg.ValidateBasic(); err != nil {
					return err
				}
				err := tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)

				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
