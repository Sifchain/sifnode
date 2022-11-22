package cli

import (
	"encoding/json"
	"github.com/Sifchain/sifnode/x/ethbridge/utils"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

// GetCmdBurn is the CLI command for burning some of your eth and triggering an event
//
//nolint:lll
func GetCmdBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [ethereum-receiver-address] [amount] [symbol] [crossChainFee] --network-descriptor [network-descriptor]",
		Short: "burn CrossChainFee or cERC20 on the Cosmos chain",
		Long: `This should be used to burn CrossChainFee or cERC20. It will burn your coins on the Cosmos Chain, removing them from your account and deducting them from the supply.
		It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the withdrawal of the original ETH/ERC20 to you from the Ethereum contract!`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			flags := cmd.Flags()

			networkDescriptorStr, err := flags.GetString(types.FlagEthereumChainID)
			if err != nil {
				return err
			}

			networkDescriptor, err := oracletypes.ParseNetworkDescriptor(networkDescriptorStr)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return errors.Errorf("invalid [ethereum-receiver-address]: %s", args[0])
			}
			ethereumReceiver := types.NewEthereumAddress(args[0])

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[1]) {
				return types.ErrInvalidAmount
			}
			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return err
			}

			if amount.LTE(sdk.NewInt(0)) {
				return types.ErrInvalidAmount
			}

			symbol := args[2]

			crossChainFee, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("Error parsing cross-chain-fee amount")
			}

			msg := types.NewMsgBurn(networkDescriptor, clientCtx.FromAddress, ethereumReceiver, amount, symbol, crossChainFee)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdLock is the CLI command for locking some of your coins and triggering an event
func GetCmdLock() *cobra.Command {
	//nolint:lll
	cmd := &cobra.Command{
		Use:   "lock [ethereum-receiver-address] [amount] [symbol] [crossChainFee] --network-descriptor [network-descriptor]",
		Short: "This should be used to lock Cosmos-originating coins (eg: ATOM). It will lock up your coins in the supply module, removing them from your account. It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the minting of the pegged token on Ethereum to you!",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			flags := cmd.Flags()

			networkDescriptorStr, err := flags.GetString(types.FlagEthereumChainID)
			if err != nil {
				return err
			}

			networkDescriptor, err := oracletypes.ParseNetworkDescriptor(networkDescriptorStr)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return errors.Errorf("invalid [ethereum-receiver-address]: %s", args[0])
			}
			ethereumReceiver := types.NewEthereumAddress(args[0])

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[1]) {
				return types.ErrInvalidAmount
			}

			amount, ok := sdk.NewIntFromString(args[1])

			if !ok {
				return err
			}
			if amount.LTE(sdk.NewInt(0)) {
				return types.ErrInvalidAmount
			}

			symbol := args[2]

			crossChainFee, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("Error parsing cross-chain-fee amount")
			}

			msg := types.NewMsgLock(networkDescriptor, clientCtx.FromAddress, ethereumReceiver, amount, symbol, crossChainFee)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, flags, &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateWhiteListValidator is the CLI command for update the validator whitelist
func GetCmdUpdateWhiteListValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-whitelist-validator [network-id]  [validator-address] [power] --node [node-address]",
		Short: "This should be used to update the validator whitelist.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			networkDescriptor, err := oracletypes.ParseNetworkDescriptor(args[0])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}

			validatorAddress, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			power, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateWhiteListValidator(networkDescriptor, clientCtx.FromAddress, validatorAddress, uint32(power))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateCrossChainFeeReceiverAccount is the CLI command to update the sifchain account that receives the cross-chain-fee proceeds
func GetCmdUpdateCrossChainFeeReceiverAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-cross-chain-fee-receiver-account  [cross-chain-fee-receiver-account]",
		Short: "This should be used to set the crosschain fee receiver account.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			crossChainFeeReceiverAccount, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateCrossChainFeeReceiverAccount(clientCtx.FromAddress, crossChainFeeReceiverAccount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRescueCrossChainFee is the CLI command to send the message to transfer cross-chain-fee from ethbridge module to account
func GetCmdRescueCrossChainFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rescue-cross-chain-fee  [cross-chain-fee-receiver-account] [cross-chain-fee] [cross-chain-fee-amount]",
		Short: "This should be used to send cross-chain-fee from ethbridge to an account.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			crossChainFeeReceiverAccount, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			crossChainFee, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("Error parsing cross-chain-fee amount")
			}

			crosschainFeeSymbol := args[2]

			msg := types.NewMsgRescueCrossChainFee(clientCtx.FromAddress, crossChainFeeReceiverAccount, crosschainFeeSymbol, crossChainFee)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSetCrossChainFee is the CLI command to send the message to set crosschain fee for network
func GetCmdSetCrossChainFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-cross-chain-fee  [network-descriptor] [fee-currency] [fee-currency-gas] [minimum-lock-cost] [minimum-burn-cost]",
		Short: "This should be used to set crosschain fee for a network.",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			networkDescriptor, err := oracletypes.ParseNetworkDescriptor(args[0])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}

			feeCurrency := args[1]

			feeCurrencyGas, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return errors.New("Error parsing feeCurrencyGas")
			}

			minimumLockCost, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("Error parsing minimumLockCost")
			}

			minimumBurnCost, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("Error parsing minimumBurnCost")
			}

			msg := types.NewMsgSetFeeInfo(clientCtx.FromAddress,
				networkDescriptor,
				feeCurrency,
				feeCurrencyGas,
				minimumLockCost,
				minimumBurnCost)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdSignProphecy is the CLI command to send the message to sign a prophecy
func GetCmdSignProphecy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign signature [network-id] [prophecy-id] [ethereum-address] [signature]",
		Short: "This should be used to sign a prophecy.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			networkDescriptor, err := oracletypes.ParseNetworkDescriptor(args[0])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}

			prophecyID := args[1]
			ethereumAddress := args[2]
			signature := args[3]

			msg := types.NewMsgSignProphecy(clientCtx.FromAddress.String(), networkDescriptor, []byte(prophecyID), ethereumAddress, signature)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUpdateConsensusNeeded is the CLI command to send the message to update consensusNeeded
func GetCmdUpdateConsensusNeeded() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-consensus-needed [network-id] [consensus-needed]",
		Short: "This should be used to update consensus-needed.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			networkDescriptor, err := oracletypes.ParseNetworkDescriptor(args[0])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}

			consensusNeeded, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return errors.New("Error parsing consensus needed")
			}

			if consensusNeeded > 100 {
				return errors.New("Error consensus needed value too large")
			}

			msg := types.NewMsgUpdateConsensusNeeded(clientCtx.FromAddress.String(), networkDescriptor, uint32(consensusNeeded))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdSetBlacklist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-blacklist [msgsetblacklist.json]",
		Short: "Set the ethereum address blacklist.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgSetBlacklist{}

			file, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			contents, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			err = json.Unmarshal(contents, &msg)
			if err != nil {
				return err
			}

			msg.From = clientCtx.FromAddress.String()

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdPause() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-pause [pause]",
		Short: "pause or unpause Lock and Burn transactions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			isPaused, err := utils.ParseStringToBool(args[0])
			if err != nil {
				return err
			}
			signer := clientCtx.GetFromAddress()
			msg := types.MsgPause{
				Signer:   signer.String(),
				IsPaused: isPaused,
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
