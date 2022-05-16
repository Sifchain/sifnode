package cli

import (
	"encoding/json"
	"io/ioutil"
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

func parseNetworkDescriptor(networkDescriptorStr string) (oracletypes.NetworkDescriptor, error) {
	networkDescriptor, err := strconv.Atoi(networkDescriptorStr)
	if err != nil {
		return -1, err
	} else if networkDescriptor < 0 || networkDescriptor > 9999 {
		return -1, errors.Errorf("Invalid %s. Valid range: [0-9999], received %d", types.FlagEthereumChainID, networkDescriptor)
	} else if !oracletypes.NetworkDescriptor(networkDescriptor).IsValid() {
		return -1, errors.Errorf("Invalid %s. Invalid value, received %d", types.FlagEthereumChainID, networkDescriptor)
	}
	return oracletypes.NetworkDescriptor(networkDescriptor), nil
}

// GetCmdBurn is the CLI command for burning some of your eth and triggering an event
//nolint:lll
func GetCmdBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [cosmos-sender-address] [ethereum-receiver-address] [amount] [symbol] [crossChainFee] --network-descriptor [network-descriptor]",
		Short: "burn CrossChainFee or cERC20 on the Cosmos chain",
		Long: `This should be used to burn CrossChainFee or cERC20. It will burn your coins on the Cosmos Chain, removing them from your account and deducting them from the supply.
		It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the withdrawal of the original ETH/ERC20 to you from the Ethereum contract!`,
		Args: cobra.ExactArgs(5),
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

			networkDescriptor, err := parseNetworkDescriptor(networkDescriptorStr)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[1]) {
				return errors.Errorf("invalid [ethereum-receiver-address]: %s", args[1])
			}
			ethereumReceiver := types.NewEthereumAddress(args[1])

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[2]) {
				return types.ErrInvalidAmount
			}
			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return err
			}

			if amount.LTE(sdk.NewInt(0)) {
				return types.ErrInvalidAmount
			}

			symbol := args[3]

			crossChainFee, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("Error parsing cross-chain-fee amount")
			}

			msg := types.NewMsgBurn(networkDescriptor, cosmosSender, ethereumReceiver, amount, symbol, crossChainFee)
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
		Use:   "lock [cosmos-sender-address] [ethereum-receiver-address] [amount] [symbol] [crossChainFee] --network-descriptor [network-descriptor]",
		Short: "This should be used to lock Cosmos-originating coins (eg: ATOM). It will lock up your coins in the supply module, removing them from your account. It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the minting of the pegged token on Etherum to you!",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// TODO: Rather use standard --from arg in integration test, and remove first redudant param in this command.
			clientCtx = clientCtx.WithFromAddress(cosmosSender)

			flags := cmd.Flags()

			networkDescriptorStr, err := flags.GetString(types.FlagEthereumChainID)
			if err != nil {
				return err
			}

			networkDescriptor, err := parseNetworkDescriptor(networkDescriptorStr)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[1]) {
				return errors.Errorf("invalid [ethereum-receiver-address]: %s", args[1])
			}
			ethereumReceiver := types.NewEthereumAddress(args[1])

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[2]) {
				return types.ErrInvalidAmount
			}

			amount, ok := sdk.NewIntFromString(args[2])

			if !ok {
				return err
			}
			if amount.LTE(sdk.NewInt(0)) {
				return types.ErrInvalidAmount
			}

			symbol := args[3]

			crossChainFee, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("Error parsing cross-chain-fee amount")
			}

			msg := types.NewMsgLock(networkDescriptor, cosmosSender, ethereumReceiver, amount, symbol, crossChainFee)
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
		Use:   "update-whitelist-validator [cosmos-sender-address] [network-id]  [validator-address] [power] --node [node-address]",
		Short: "This should be used to update the validator whitelist.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			networkDescriptor, err := strconv.Atoi(args[1])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}

			validatorAddress, err := sdk.ValAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			power, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateWhiteListValidator(oracletypes.NetworkDescriptor(networkDescriptor), cosmosSender, validatorAddress, uint32(power))
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
		Use:   "update-cross-chain-fee-receiver-account [cosmos-sender-address] [cross-chain-fee-receiver-account]",
		Short: "This should be used to set the crosschain fee receiver account.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			crossChainFeeReceiverAccount, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateCrossChainFeeReceiverAccount(cosmosSender, crossChainFeeReceiverAccount)
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
		Use:   "rescue-cross-chain-fee [cosmos-sender-address] [cross-chain-fee-receiver-account] [cross-chain-fee] [cross-chain-fee-amount]",
		Short: "This should be used to send cross-chain-fee from ethbridge to an account.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			crossChainFeeReceiverAccount, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			crossChainFee, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return errors.New("Error parsing cross-chain-fee amount")
			}

			crosschainFeeSymbol := args[3]

			msg := types.NewMsgRescueCrossChainFee(cosmosSender, crossChainFeeReceiverAccount, crosschainFeeSymbol, crossChainFee)
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
		// TODO: Rename variable network-id to network descriptor
		Use:   "set-cross-chain-fee [cosmos-sender-address] [network-id] [fee-currency] [fee-currency-gas] [minimum-lock-cost] [minimum-burn-cost]",
		Short: "This should be used to set crosschain fee for a network.",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return err
			}

			networkDescriptor, err := strconv.Atoi(args[1])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}

			feeCurrency := args[2]

			feeCurrencyGas, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("Error parsing feeCurrencyGas")
			}

			minimumLockCost, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("Error parsing minimumLockCost")
			}

			minimumBurnCost, ok := sdk.NewIntFromString(args[5])
			if !ok {
				return errors.New("Error parsing minimumBurnCost")
			}

			msg := types.NewMsgSetFeeInfo(cosmosSender,
				oracletypes.NetworkDescriptor(networkDescriptor),
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
		Use:   "sign signature [cosmos-sender-address] [network-id] [prophecy-id] [ethereum-address] [signature]",
		Short: "This should be used to sign a prophecy.",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			networkDescriptor, err := strconv.Atoi(args[1])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}

			prophecyID := args[2]
			ethereumAddress := args[3]
			signature := args[4]

			msg := types.NewMsgSignProphecy(args[0], oracletypes.NetworkDescriptor(networkDescriptor), []byte(prophecyID), ethereumAddress, signature)
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
		Use:   "update-consensus-needed [cosmos-sender-address] [network-id] [consensus-needed]",
		Short: "This should be used to update consensus-needed.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.New("Error cosmos sender address")
			}

			networkDescriptor, err := strconv.Atoi(args[1])
			if err != nil {
				return errors.New("Error parsing network descriptor")
			}

			consensusNeeded, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return errors.New("Error parsing consensus needed")
			}

			if consensusNeeded > 100 {
				return errors.New("Error consensus needed value too large")
			}

			msg := types.NewMsgUpdateConsensusNeeded(args[0], oracletypes.NetworkDescriptor(networkDescriptor), uint32(consensusNeeded))
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

			contents, err := ioutil.ReadFile(file)
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
