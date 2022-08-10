//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package cli

import (
	"encoding/json"
	"errors"
	"io/ioutil"

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
		GetUpdatePoolsCmd(),
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

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
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

			leverage, err := cmd.Flags().GetString("leverage")
			if err != nil {
				return err
			}
			leverageDec := sdk.MustNewDecFromStr(leverage)

			msg := types.MsgOpen{
				Signer:           signer.String(),
				CollateralAsset:  collateralAsset,
				CollateralAmount: sdk.NewUintFromString(collateralAmount),
				BorrowAsset:      borrowAsset,
				Position:         positionEnum,
				Leverage:         leverageDec,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().String("collateral_amount", "0", "amount of collateral asset")
	cmd.Flags().String("collateral_asset", "", "symbol of asset")
	cmd.Flags().String("borrow_asset", "", "symbol of asset")
	cmd.Flags().String("position", "", "type of position")
	cmd.Flags().String("leverage", "", "leverage of position")
	_ = cmd.MarkFlagRequired("collateral_amount")
	_ = cmd.MarkFlagRequired("collateral_asset")
	_ = cmd.MarkFlagRequired("borrow_asset")
	_ = cmd.MarkFlagRequired("position")
	_ = cmd.MarkFlagRequired("leverage")
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

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			id, err := cmd.Flags().GetUint64("id")
			if err != nil {
				return err
			}

			msg := types.MsgClose{
				Signer: signer.String(),
				Id:     id,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().Uint64("id", 0, "id of the position")
	_ = cmd.MarkFlagRequired("id")
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

			msg := types.MsgForceClose{
				Signer:     signer.String(),
				MtpAddress: MtpAddress,
				Id:         id,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().String("mtp_address", "", "mtp address")
	cmd.Flags().Uint64("id", 0, "id of the position")
	_ = cmd.MarkFlagRequired("mtp_address")
	_ = cmd.MarkFlagRequired("id")
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

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			msg := types.MsgUpdateParams{
				Signer: signer.String(),
				Params: &types.Params{
					LeverageMax:              sdk.MustNewDecFromStr(viper.GetString("leverage-max")),
					InterestRateMax:          sdk.MustNewDecFromStr(viper.GetString("interest-rate-max")),
					InterestRateMin:          sdk.MustNewDecFromStr(viper.GetString("interest-rate-min")),
					InterestRateIncrease:     sdk.MustNewDecFromStr(viper.GetString("interest-rate-increase")),
					InterestRateDecrease:     sdk.MustNewDecFromStr(viper.GetString("interest-rate-decrease")),
					HealthGainFactor:         sdk.MustNewDecFromStr(viper.GetString("health-gain-factor")),
					ForceCloseThreshold:      sdk.MustNewDecFromStr(viper.GetString("force-close-threshold")),
					PoolOpenThreshold:        sdk.MustNewDecFromStr(viper.GetString("pool-open-threshold")),
					EpochLength:              viper.GetInt64("epoch-length"),
					MaxOpenPositions:         viper.GetUint64("max-open-positions"),
					RemovalQueueThreshold:    sdk.MustNewDecFromStr(viper.GetString("removal-queue-threshold")),
					ForceCloseFundPercentage: sdk.MustNewDecFromStr(viper.GetString("force-close-fund-percentage")),
					InsuranceFundAddress:     viper.GetString("insurance-fund-address"),
					SqModifier:               sdk.MustNewDecFromStr(viper.GetString("sq-modifier")),
					SafetyFactor:             sdk.MustNewDecFromStr(viper.GetString("safety-factor")),
				},
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String("leverage-max", "", "max leverage (integer)")
	cmd.Flags().String("interest-rate-max", "", "max interest rate (decimal)")
	cmd.Flags().String("interest-rate-min", "", "min interest rate (decimal)")
	cmd.Flags().String("interest-rate-increase", "", "interest rate increase (decimal)")
	cmd.Flags().String("interest-rate-decrease", "", "interest rate decrease (decimal)")
	cmd.Flags().String("health-gain-factor", "", "health gain factor (decimal)")
	cmd.Flags().String("force-close-threshold", "", "force close threshold (decimal range 0-1)")
	cmd.Flags().Int64("epoch-length", 1, "epoch length in blocks (integer)")
	cmd.Flags().Uint64("max-open-positions", 10000, "max open positions")
	cmd.Flags().String("removal-queue-threshold", "", "removal queue threshold (decimal range 0-1)")
	cmd.Flags().String("pool-open-threshold", "", "threshold to prevent new positions (decimal range 0-1)")
	cmd.Flags().String("force-close-fund-percentage", "", "percentage of force close proceeds for insurance fund (decimal range 0-1)")
	cmd.Flags().String("insurance-fund-address", "", "address of insurance fund wallet")
	cmd.Flags().String("sq-modifier", "", "the modifier value for the removal queue's sq formula")
	cmd.Flags().String("safety-factor", "", "the safety factor used in liquidation ratio")
	_ = cmd.MarkFlagRequired("leverage-max")
	_ = cmd.MarkFlagRequired("interest-rate-max")
	_ = cmd.MarkFlagRequired("interest-rate-min")
	_ = cmd.MarkFlagRequired("interest-rate-increase")
	_ = cmd.MarkFlagRequired("interest-rate-decrease")
	_ = cmd.MarkFlagRequired("health-gain-factor")
	//_ = cmd.MarkFlagRequired("force-close-threshold")
	_ = cmd.MarkFlagRequired("removal-queue-threshold")
	_ = cmd.MarkFlagRequired("max-open-positions")
	_ = cmd.MarkFlagRequired("pool-open-threshold")
	_ = cmd.MarkFlagRequired("insurance-fund-address")
	_ = cmd.MarkFlagRequired("force-close-fund-percentage")
	_ = cmd.MarkFlagRequired("sq-modifier")
	_ = cmd.MarkFlagRequired("safety-factor")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetUpdatePoolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pools [pools.json]",
		Short: "Update margin enabled pools, and closed pools",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			if signer == nil {
				return errors.New("signer address is missing")
			}

			pools, err := readPoolsJSON(args[0])
			if err != nil {
				return err
			}

			closedPools, err := readPoolsJSON(viper.GetString("closed-pools"))
			if err != nil {
				return err
			}

			msg := types.MsgUpdatePools{
				Signer:      signer.String(),
				Pools:       pools,
				ClosedPools: closedPools,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().String("closed-pools", "", "pools that new positions cannot be opened on")
	_ = cmd.MarkFlagRequired("closed-pools")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func readPoolsJSON(filename string) ([]string, error) {
	var pools []string
	bz, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{}, err
	}
	err = json.Unmarshal(bz, &pools)
	if err != nil {
		return []string{}, err
	}

	return pools, nil
}
