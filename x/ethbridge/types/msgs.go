package types

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	gethCommon "github.com/ethereum/go-ethereum/common"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgLock is a constructor function for MsgLock
func NewMsgLock(
	networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, symbol string, nativeTokenAmount sdk.Int) MsgLock {
	return MsgLock{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      cosmosSender.String(),
		EthereumReceiver:  ethereumReceiver.String(),
		Amount:            amount,
		Symbol:            symbol,
		NativeTokenAmount: nativeTokenAmount,
	}
}

// Route should return the name of the module
func (msg MsgLock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgLock) Type() string { return "lock" }

// ValidateBasic runs stateless checks on the message
func (msg MsgLock) ValidateBasic() error {
	if strconv.FormatInt(int64(msg.NetworkDescriptor), 10) == "" {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", msg.NetworkDescriptor)
	}

	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.EthereumReceiver == "" {
		return ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.EthereumReceiver) {
		return ErrInvalidEthAddress
	}

	if msg.Amount.LTE(sdk.NewInt(0)) {
		return ErrInvalidAmount
	}

	if len(msg.Symbol) == 0 {
		return ErrInvalidSymbol
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgLock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgLock) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgBurn is a constructor function for MsgBurn
func NewMsgBurn(
	networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, symbol string, nativeTokenAmount sdk.Int) MsgBurn {
	return MsgBurn{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      cosmosSender.String(),
		EthereumReceiver:  ethereumReceiver.String(),
		Amount:            amount,
		Symbol:            symbol,
		NativeTokenAmount: nativeTokenAmount,
	}
}

// Route should return the name of the module
func (msg MsgBurn) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBurn) Type() string { return "burn" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBurn) ValidateBasic() error {
	if msg.NetworkDescriptor == 0 {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", msg.NetworkDescriptor)
	}

	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.EthereumReceiver == "" {
		return ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.EthereumReceiver) {
		return ErrInvalidEthAddress
	}

	if msg.Amount.LTE(sdk.NewInt(0)) {
		return ErrInvalidAmount
	}

	prefixLength := len(PeggedCoinPrefix)
	if len(msg.Symbol) <= prefixLength+1 {
		return ErrInvalidBurnSymbol
	}

	symbolPrefix := msg.Symbol[:prefixLength]
	if symbolPrefix != PeggedCoinPrefix {
		return ErrInvalidBurnSymbol
	}

	symbolSuffix := msg.Symbol[prefixLength:]
	if len(symbolSuffix) == 0 {
		return ErrInvalidBurnSymbol
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgCreateEthBridgeClaim is a constructor function for MsgCreateBridgeClaim
func NewMsgCreateEthBridgeClaim(ethBridgeClaim *EthBridgeClaim) MsgCreateEthBridgeClaim {
	return MsgCreateEthBridgeClaim{
		EthBridgeClaim: ethBridgeClaim,
	}
}

// Route should return the name of the module
func (msg MsgCreateEthBridgeClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateEthBridgeClaim) Type() string { return "create_bridge_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateEthBridgeClaim) ValidateBasic() error {
	if msg.EthBridgeClaim.CosmosReceiver == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.EthBridgeClaim.CosmosReceiver)
	}

	if msg.EthBridgeClaim.ValidatorAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.EthBridgeClaim.ValidatorAddress)
	}

	if msg.EthBridgeClaim.Nonce < 0 {
		return ErrInvalidEthNonce
	}

	if !gethCommon.IsHexAddress(msg.EthBridgeClaim.EthereumSender) {
		return ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.EthBridgeClaim.BridgeContractAddress) {
		return ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.EthBridgeClaim.TokenContractAddress) {
		return ErrInvalidEthAddress
	}

	if strings.ToLower(msg.EthBridgeClaim.Symbol) == "eth" &&
		NewEthereumAddress(msg.EthBridgeClaim.TokenContractAddress) != NewEthereumAddress("0x0000000000000000000000000000000000000000") {
		return ErrInvalidEthSymbol
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateEthBridgeClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgCreateEthBridgeClaim) GetSigners() []sdk.AccAddress {
	validatorAddress, err := sdk.ValAddressFromBech32(msg.EthBridgeClaim.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sdk.AccAddress(validatorAddress)}
}

// NewMsgUpdateNativeTokenReceiverAccount is a constructor function for MsgUpdateNativeTokenReceiverAccount
func NewMsgUpdateNativeTokenReceiverAccount(cosmosSender sdk.AccAddress,
	nativeTokenReceiverAccount sdk.AccAddress) MsgUpdateNativeTokenReceiverAccount {
	return MsgUpdateNativeTokenReceiverAccount{
		CosmosSender:               cosmosSender.String(),
		NativeTokenReceiverAccount: nativeTokenReceiverAccount.String(),
	}
}

// Route should return the name of the module
func (msg MsgUpdateNativeTokenReceiverAccount) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateNativeTokenReceiverAccount) Type() string {
	return "update_native_token_receiver_account"
}

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateNativeTokenReceiverAccount) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.NativeTokenReceiverAccount == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.NativeTokenReceiverAccount)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateNativeTokenReceiverAccount) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgUpdateNativeTokenReceiverAccount) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgRescueNativeToken is a constructor function for NewMsgRescueNativeToken
func NewMsgRescueNativeToken(cosmosSender sdk.AccAddress, cosmosReceiver sdk.AccAddress, nativeToken string, nativeTokenAmount sdk.Int) MsgRescueNativeToken {
	return MsgRescueNativeToken{
		CosmosSender:      cosmosSender.String(),
		CosmosReceiver:    cosmosReceiver.String(),
		NativeTokenSymbol: nativeToken,
		NativeTokenAmount: nativeTokenAmount,
	}
}

// Route should return the name of the module
func (msg MsgRescueNativeToken) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRescueNativeToken) Type() string { return "rescue_native_token" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRescueNativeToken) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.CosmosReceiver == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosReceiver)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRescueNativeToken) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgRescueNativeToken) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgUpdateWhiteListValidator is a constructor function for MsgUpdateWhiteListValidator
func NewMsgUpdateWhiteListValidator(networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress,
	validator sdk.ValAddress, power uint32) MsgUpdateWhiteListValidator {
	return MsgUpdateWhiteListValidator{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      cosmosSender.String(),
		Validator:         validator.String(),
		Power:             power,
	}
}

// Route should return the name of the module
func (msg MsgUpdateWhiteListValidator) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateWhiteListValidator) Type() string { return "update_whitelist_validator" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateWhiteListValidator) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.Validator == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateWhiteListValidator) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgUpdateWhiteListValidator) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgSetNativeToken is a constructor function for MsgSetNativeToken
func NewMsgSetNativeToken(cosmosSender sdk.AccAddress, networkDescriptor oracletypes.NetworkDescriptor, nativeToken string) MsgSetNativeToken {
	return MsgSetNativeToken{
		CosmosSender:      cosmosSender.String(),
		NetworkDescriptor: networkDescriptor,
		NativeToken:       nativeToken,
	}
}

// Route should return the name of the module
func (msg MsgSetNativeToken) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetNativeToken) Type() string { return "set_native_token" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetNativeToken) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if !msg.NetworkDescriptor.IsValid() {
		return errors.New("network descriptor is invalid")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetNativeToken) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSetNativeToken) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}
