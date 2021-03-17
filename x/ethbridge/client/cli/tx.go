package cli

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

// GetCmdCreateEthBridgeClaim is the CLI command for creating a claim on an ethereum prophecy
//nolint:lll
func GetCmdCreateEthBridgeClaim(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-claim [bridge-registry-contract] [nonce] [symbol] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount] [claim-type] --ethereum-chain-id [ethereum-chain-id] --token-contract-address [token-contract-address]",
		Short: "create a claim on an ethereum prophecy",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			ethereumChainIDString := viper.GetString(types.FlagEthereumChainID)
			ethereumChainID, err := strconv.Atoi(ethereumChainIDString)
			if err != nil {
				return err
			}

			tokenContractString := viper.GetString(types.FlagTokenContractAddr)
			if !common.IsHexAddress(tokenContractString) {
				return errors.Errorf("invalid [token-contract-address]: %s", tokenContractString)
			}
			tokenContract := types.NewEthereumAddress(tokenContractString)

			if !common.IsHexAddress(args[0]) {
				return errors.Errorf("invalid [bridge-registry-contract]: %s", args[0])
			}
			bridgeContract := types.NewEthereumAddress(args[0])

			nonce, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			symbol := args[2]
			ethereumSender := types.NewEthereumAddress(args[3])
			if !common.IsHexAddress(args[3]) {
				return errors.Errorf("invalid [ethereum-sender-address]: %s", args[0])
			}
			cosmosReceiver, err := sdk.AccAddressFromBech32(args[4])
			if err != nil {
				return err
			}

			validator, err := sdk.ValAddressFromBech32(args[5])
			if err != nil {
				return err
			}

			var digitCheck = regexp.MustCompile(`^[0-9]+$`)
			if !digitCheck.MatchString(args[6]) {
				return types.ErrInvalidAmount
			}

			bigIntAmount, ok := sdk.NewIntFromString(args[6])
			if !ok {
				fmt.Println("SetString: error")
				return types.ErrInvalidAmount
			}

			if bigIntAmount.LTE(sdk.NewInt(0)) {
				return types.ErrInvalidAmount
			}

			claimType, err := types.StringToClaimType(args[7])
			if err != nil {
				return err
			}

			ethBridgeClaim := types.NewEthBridgeClaim(ethereumChainID, bridgeContract, nonce, symbol, tokenContract,
				ethereumSender, cosmosReceiver, validator, bigIntAmount, claimType)

			msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBurn is the CLI command for burning some of your eth and triggering an event
//nolint:lll
func GetCmdBurn(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burn [cosmos-sender-address] [ethereum-receiver-address] [amount] [symbol] [cethAmount] --ethereum-chain-id [ethereum-chain-id]",
		Short: "burn cETH or cERC20 on the Cosmos chain",
		Long: `This should be used to burn cETH or cERC20. It will burn your coins on the Cosmos Chain, removing them from your account and deducting them from the supply.
		It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the withdrawal of the original ETH/ERC20 to you from the Ethereum contract!`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			ethereumChainIDString := viper.GetString(types.FlagEthereumChainID)
			ethereumChainID, err := strconv.Atoi(ethereumChainIDString)
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

			cethAmount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("Error parsing ceth amount")
			}

			msg := types.NewMsgBurn(ethereumChainID, cosmosSender, ethereumReceiver, amount, symbol, cethAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdLock is the CLI command for locking some of your coins and triggering an event
func GetCmdLock(cdc *codec.Codec) *cobra.Command {
	//nolint:lll
	return &cobra.Command{
		Use:   "lock [cosmos-sender-address] [ethereum-receiver-address] [amount] [symbol] [cethAmount] --ethereum-chain-id [ethereum-chain-id]",
		Short: "This should be used to lock Cosmos-originating coins (eg: ATOM). It will lock up your coins in the supply module, removing them from your account. It will also trigger an event on the Cosmos Chain for relayers to watch so that they can trigger the minting of the pegged token on Etherum to you!",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			ethereumChainIDString := viper.GetString(types.FlagEthereumChainID)
			ethereumChainID, err := strconv.Atoi(ethereumChainIDString)
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

			cethAmount, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("Error parsing ceth amount")
			}

			msg := types.NewMsgLock(ethereumChainID, cosmosSender, ethereumReceiver, amount, symbol, cethAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUpdateWhiteListValidator is the CLI command for update the validator whitelist
func GetCmdUpdateWhiteListValidator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update_whitelist_validator [cosmos-sender-address] [validator-address] [operation-type] --node [node-address]",
		Short: "This should be used to update the validator whitelist.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			validatorAddress, err := sdk.ValAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			operationType := args[2]
			if operationType != "add" && operationType != "remove" {
				return errors.Errorf("invalid [operation-type]: %s", args[2])
			}

			msg := types.NewMsgUpdateWhiteListValidator(cosmosSender, validatorAddress, operationType)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUpdateCethReceiverAccount is the CLI command for update the validator whitelist
func GetCmdUpdateCethReceiverAccount(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update_ceth_receiver_account [cosmos-sender-address] [ceth_receiver_account] --node [node-address]",
		Short: "This should be used to set the ceth receiver account.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			cosmosSender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			cethReceiverAccount, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateCethReceiverAccount(cosmosSender, cethReceiverAccount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
