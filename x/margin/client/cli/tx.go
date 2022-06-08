package cli

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/spf13/viper"

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
		GetOpenCmd(),
		GetCloseCmd(),
		GetForceCloseCmd(),
		GetUpdateParamsCmd(),
	)
	return cmd
}

func GetOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open margin position",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateralAsset, err := cmd.Flags().GetString("collateral_asset")
			if err != nil {
				return err
			}

			collateralAmount, err := cmd.Flags().GetString("collateral_amount")
			if err != nil {
				return err
			}

			borrowAsset, err := cmd.Flags().GetString("borrow_asset")
			if err != nil {
				return err
			}

			position, err := cmd.Flags().GetString("position")
			if err != nil {
				return err
			}
			positionEnum := types.GetPositionFromString(position)

			msg := types.MsgOpen{
				Signer:           clientCtx.GetFromAddress().String(),
				CollateralAsset:  collateralAsset,
				CollateralAmount: sdk.NewUintFromString(collateralAmount),
				BorrowAsset:      borrowAsset,
				Position:         positionEnum,
			}

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().String("collateral_amount", "0", "amount of collateral asset")
	cmd.Flags().String("collateral_asset", "", "symbol of asset")
	cmd.Flags().String("borrow_asset", "", "symbol of asset")
	cmd.Flags().String("position", "", "type of position")
	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func GetCloseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close margin position",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := cmd.Flags().GetUint64("id")
			if err != nil {
				return err
			}

			msg := types.MsgClose{
				Signer: clientCtx.GetFromAddress().String(),
				Id:     id,
			}

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().Uint64("id", 0, "id of the position")
	flags.AddTxFlagsToCmd(cmd)
	return cmd

}

func GetForceCloseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "force-close",
		Short: "Force close margin position",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			MtpAddress, err := cmd.Flags().GetString("mtp_address")
			if err != nil {
				return err
			}

			id, err := cmd.Flags().GetUint64("id")
			if err != nil {
				return err
			}

			msg := types.MsgForceClose{
				Signer:     clientCtx.GetFromAddress().String(),
				MtpAddress: MtpAddress,
				Id:         id,
			}

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().Uint64("id", 0, "id of the position")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetUpdateParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params",
		Short: "Update margin params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgUpdateParams{
				Signer: clientCtx.GetFromAddress().String(),
				Params: &types.Params{
					LeverageMax:          sdk.NewUintFromString(viper.GetString("leverage-max")),
					InterestRateMax:      sdk.MustNewDecFromStr(viper.GetString("interest-rate-max")),
					InterestRateMin:      sdk.MustNewDecFromStr(viper.GetString("interest-rate-min")),
					InterestRateIncrease: sdk.MustNewDecFromStr(viper.GetString("interest-rate-increase")),
					InterestRateDecrease: sdk.MustNewDecFromStr(viper.GetString("interest-rate-decrease")),
					HealthGainFactor:     sdk.MustNewDecFromStr(viper.GetString("health-gain-factor")),
					ForceCloseThreshold:  sdk.MustNewDecFromStr(viper.GetString("force-close-threshold")),
					EpochLength:          viper.GetInt64("epoch-length"),
				},
			}

			err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("leverage-max", "", "max leverage")
	cmd.Flags().String("interest-rate-max", "", "max interest rate")
	cmd.Flags().String("interest-rate-min", "", "min interest rate")
	cmd.Flags().String("interest-rate-increase", "", "interest rate increase")
	cmd.Flags().String("interest-rate-decrease", "", "interest rate decrease")
	cmd.Flags().String("health-gain-factor", "", "health gain factor")
	cmd.Flags().String("force-close-threshold", "", "force close threshold")
	cmd.Flags().Int64("epoch-length", 1, "epoch length in blocks")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}