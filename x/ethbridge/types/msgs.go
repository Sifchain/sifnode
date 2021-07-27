package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	gethCommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	burnGasCost = 60000000000 * 393000 // assuming 60gigawei gas prices
	lockGasCost = 60000000000 * 393000
)

// NewMsgLock is a constructor function for MsgLock
func NewMsgLock(
	ethereumChainID int64, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, symbol string, cethAmount sdk.Int) MsgLock {
	return MsgLock{
		EthereumChainId:  ethereumChainID,
		CosmosSender:     cosmosSender.String(),
		EthereumReceiver: ethereumReceiver.String(),
		Amount:           amount,
		Symbol:           symbol,
		CethAmount:       cethAmount,
	}
}

// Route should return the name of the module
func (msg MsgLock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgLock) Type() string { return "lock" }

// ValidateBasic runs stateless checks on the message
func (msg MsgLock) ValidateBasic() error {
	if strconv.FormatInt(msg.EthereumChainId, 10) == "" {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", msg.EthereumChainId)
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

	// if you don't pay enough gas, this tx won't go through
	if msg.CethAmount.LT(sdk.NewInt(lockGasCost)) {
		return ErrCethAmount
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
	ethereumChainID int64, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, symbol string, cethAmount sdk.Int) MsgBurn {
	return MsgBurn{
		EthereumChainId:  ethereumChainID,
		CosmosSender:     cosmosSender.String(),
		EthereumReceiver: ethereumReceiver.String(),
		Amount:           amount,
		Symbol:           symbol,
		CethAmount:       cethAmount,
	}
}

// Route should return the name of the module
func (msg MsgBurn) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBurn) Type() string { return "burn" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBurn) ValidateBasic() error {
	if msg.EthereumChainId == 0 {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", msg.EthereumChainId)
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

	// check that enough ceth is sent to cover the gas cost.
	if msg.CethAmount.LT(sdk.NewInt(burnGasCost)) {
		return ErrCethAmount
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

// NewMsgUpdateCethReceiverAccount is a constructor function for MsgUpdateCethReceiverAccount
func NewMsgUpdateCethReceiverAccount(cosmosSender sdk.AccAddress,
	cethReceiverAccount sdk.AccAddress) MsgUpdateCethReceiverAccount {
	return MsgUpdateCethReceiverAccount{
		CosmosSender:        cosmosSender.String(),
		CethReceiverAccount: cethReceiverAccount.String(),
	}
}

// Route should return the name of the module
func (msg MsgUpdateCethReceiverAccount) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateCethReceiverAccount) Type() string { return "update_ceth_receiver_account" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateCethReceiverAccount) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.CethReceiverAccount == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CethReceiverAccount)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateCethReceiverAccount) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgUpdateCethReceiverAccount) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgRescueCeth is a constructor function for NewMsgRescueCeth
func NewMsgRescueCeth(cosmosSender sdk.AccAddress, cosmosReceiver sdk.AccAddress, cethAmount sdk.Int) MsgRescueCeth {
	return MsgRescueCeth{
		CosmosSender:   cosmosSender.String(),
		CosmosReceiver: cosmosReceiver.String(),
		CethAmount:     cethAmount,
	}
}

// Route should return the name of the module
func (msg MsgRescueCeth) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRescueCeth) Type() string { return "rescue_ceth" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRescueCeth) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.CosmosReceiver == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosReceiver)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRescueCeth) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgRescueCeth) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgUpdateWhiteListValidator is a constructor function for MsgUpdateWhiteListValidator
func NewMsgUpdateWhiteListValidator(cosmosSender sdk.AccAddress,
	validator sdk.ValAddress, operationType string) MsgUpdateWhiteListValidator {
	return MsgUpdateWhiteListValidator{
		CosmosSender:  cosmosSender.String(),
		Validator:     validator.String(),
		OperationType: operationType,
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

// MapOracleClaimsToEthBridgeClaims maps a set of generic oracle claim data into EthBridgeClaim objects
func MapOracleClaimsToEthBridgeClaims(
	ethereumChainID int64,
	bridgeContract EthereumAddress,
	nonce int64,
	symbol string,
	tokenContract EthereumAddress,
	ethereumSender EthereumAddress,
	oracleValidatorClaims map[string]string,
	f func(int64, EthereumAddress, int64, EthereumAddress, sdk.ValAddress, string) (*EthBridgeClaim, error),
) ([]*EthBridgeClaim, error) {

	mappedClaims := make([]*EthBridgeClaim, len(oracleValidatorClaims))
	i := 0
	for validatorBech32, validatorClaim := range oracleValidatorClaims {
		validatorAddress, parseErr := sdk.ValAddressFromBech32(validatorBech32)
		if parseErr != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("failed to parse claim: %s", parseErr))
		}

		mappedClaim, err := f(
			ethereumChainID, bridgeContract, nonce, ethereumSender, validatorAddress, validatorClaim)
		if err != nil {
			return nil, err
		}
		mappedClaims[i] = mappedClaim
		i++
	}

	return mappedClaims, nil
}
