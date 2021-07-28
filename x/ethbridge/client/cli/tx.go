package cli

import (
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

// GetCmdCreateEthBridgeClaim is the CLI command for creating a claim on an ethereum prophecy
//nolint:lll
func GetCmdCreateEthBridgeClaim() *cobra.Command {
	return &cobra.Command{
		Use:   "create-claim [bridge-registry-contract] [nonce] [symbol] [name] [decimals] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount] [claim-type] --network-descriptor [network-descriptor] --token-contract-address [token-contract-address]",
		Short: "create a claim on an ethereum prophecy",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			flags := cmd.Flags()

			ethereumChainIDStr, err := flags.GetString(types.FlagEthereumChainID)
			if err != nil {
				return err
			}

			ethereumChainID, err := strconv.Atoi(ethereumChainIDStr)
			if err != nil {
				return err
			}

			tokenContractString, err := flags.GetString(types.FlagTokenContractAddr)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(tokenContractString) {
				return errors.Errorf("invalid [token-contract-address]: %s", tokenContractString)
			}

			tokenContract := types.NewEthereumAddress(tokenContractString)

			if !common.IsHexAddress(args[0]) {
				return errors.Errorf("invalid [bridge-registry-contract]: %s", args[0])
			}
			bridgeContract := types.NewEthereumAddress(args[0])

			nonce, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			symbol := args[2]

			name := args[3]

			decimals, err := strconv.ParseInt(args[4], 10, 32)
			if err != nil {
				return err
			}

			ethereumSender := types.NewEthereumAddress(args[5])
			if !common.IsHexAddress(args[5]) {
				return errors.Errorf("invalid [ethereum-sender-address]: %s", args[0])
			}

			cosmosReceiver, err := sdk.AccAddressFromBech32(args[6])
			if err != nil {
				return err
			}

			validator, err := sdk.ValAddressFromBech32(args[7])
			if err != nil {
				return err
			}

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[8]) {
				return types.ErrInvalidAmount
			}

			bigIntAmount, ok := sdk.NewIntFromString(args[8])
			if !ok {
				return types.ErrInvalidAmount
			}

			if bigIntAmount.LTE(sdk.NewInt(0)) {
				return types.ErrInvalidAmount
			}

			claimType, exist := types.ClaimType_value[args[9]]
			if !exist {
				return err
			}
			ct := types.ClaimType(claimType)

			networkDescriptor := oracletypes.NetworkDescriptor(ethereumChainID)

			denomHash := types.GetDenomHash(networkDescriptor, tokenContractString, int32(decimals), name, symbol)

			ethBridgeClaim := types.NewEthBridgeClaim(networkDescriptor, bridgeContract, nonce, symbol, tokenContract,
				ethereumSender, cosmosReceiver, validator, bigIntAmount, ct, name, int32(decimals), denomHash)

			msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, flags, &msg)
		},
	}
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

			ethereumChainIDStr, err := flags.GetString(types.FlagEthereumChainID)
			if err != nil {
				return err
			}

			ethereumChainID, err := strconv.Atoi(ethereumChainIDStr)
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

			msg := types.NewMsgBurn(oracletypes.NetworkDescriptor(ethereumChainID), cosmosSender, ethereumReceiver, amount, symbol, crossChainFee)
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

			ethereumChainIDStr, err := flags.GetString(types.FlagEthereumChainID)
			if err != nil {
				return err
			}

			ethereumChainID, err := strconv.Atoi(ethereumChainIDStr)
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

			msg := types.NewMsgLock(oracletypes.NetworkDescriptor(ethereumChainID), cosmosSender, ethereumReceiver, amount, symbol, crossChainFee)
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
		Use:   "set-cross-chain-fee [cosmos-sender-address] [network-id] [cross-chain-fee]",
		Short: "This should be used to set crosschain fee for a network.",
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

			crossChainFee := args[2]

			msg := types.NewMsgSetFeeInfo(cosmosSender, oracletypes.NetworkDescriptor(networkDescriptor), crossChainFee)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
