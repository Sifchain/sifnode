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
	burnGasCost = 160000000000 * 393000 // assuming 160gigawei gas prices
	lockGasCost = 160000000000 * 393000
)

// MsgLock defines a message for locking coins and triggering a related event
type MsgLock struct {
	CosmosSender     sdk.AccAddress  `json:"cosmos_sender" yaml:"cosmos_sender"`
	Amount           sdk.Int         `json:"amount" yaml:"amount"`
	Symbol           string          `json:"symbol" yaml:"symbol"`
	EthereumChainID  int             `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	EthereumReceiver EthereumAddress `json:"ethereum_receiver" yaml:"ethereum_receiver"`
	CethAmount       sdk.Int         `json:"ceth_amount" yaml:"ceth_amount"`
}

// NewMsgLock is a constructor function for MsgLock
func NewMsgLock(
	ethereumChainID int, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, symbol string, cethAmount sdk.Int) MsgLock {
	return MsgLock{
		EthereumChainID:  ethereumChainID,
		CosmosSender:     cosmosSender,
		EthereumReceiver: ethereumReceiver,
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
	if strconv.Itoa(msg.EthereumChainID) == "" {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", msg.EthereumChainID)
	}

	if msg.CosmosSender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	if msg.EthereumReceiver.String() == "" {
		return ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.EthereumReceiver.String()) {
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgLock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.CosmosSender}
}

// MsgBurn defines a message for burning coins and triggering a related event
type MsgBurn struct {
	CosmosSender     sdk.AccAddress  `json:"cosmos_sender" yaml:"cosmos_sender"`
	Amount           sdk.Int         `json:"amount" yaml:"amount"`
	Symbol           string          `json:"symbol" yaml:"symbol"`
	EthereumChainID  int             `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	EthereumReceiver EthereumAddress `json:"ethereum_receiver" yaml:"ethereum_receiver"`
	CethAmount       sdk.Int         `json:"ceth_amount" yaml:"ceth_amount"`
}

// NewMsgBurn is a constructor function for MsgBurn
func NewMsgBurn(
	ethereumChainID int, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, symbol string, cethAmount sdk.Int) MsgBurn {
	return MsgBurn{
		EthereumChainID:  ethereumChainID,
		CosmosSender:     cosmosSender,
		EthereumReceiver: ethereumReceiver,
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
	if strconv.Itoa(msg.EthereumChainID) == "" {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", msg.EthereumChainID)
	}
	if msg.CosmosSender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}
	if msg.EthereumReceiver.String() == "" {
		return ErrInvalidEthAddress
	}
	if !gethCommon.IsHexAddress(msg.EthereumReceiver.String()) {
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
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.CosmosSender}
}

// MsgCreateEthBridgeClaim defines a message for creating claims on the ethereum bridge
type MsgCreateEthBridgeClaim EthBridgeClaim

// NewMsgCreateEthBridgeClaim is a constructor function for MsgCreateBridgeClaim
func NewMsgCreateEthBridgeClaim(ethBridgeClaim EthBridgeClaim) MsgCreateEthBridgeClaim {
	return MsgCreateEthBridgeClaim(ethBridgeClaim)
}

// Route should return the name of the module
func (msg MsgCreateEthBridgeClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateEthBridgeClaim) Type() string { return "create_bridge_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateEthBridgeClaim) ValidateBasic() error {
	if msg.CosmosReceiver.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosReceiver.String())
	}

	if msg.ValidatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.ValidatorAddress.String())
	}

	if msg.Nonce < 0 {
		return ErrInvalidEthNonce
	}

	if !gethCommon.IsHexAddress(msg.EthereumSender.String()) {
		return ErrInvalidEthAddress
	}
	if !gethCommon.IsHexAddress(msg.BridgeContractAddress.String()) {
		return ErrInvalidEthAddress
	}
	if !gethCommon.IsHexAddress(msg.TokenContractAddress.String()) {
		return ErrInvalidEthAddress
	}
	if strings.ToLower(msg.Symbol) == "eth" &&
		msg.TokenContractAddress != NewEthereumAddress("0x0000000000000000000000000000000000000000") {
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
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

// MsgUpdateCethReceiverAccount add or remove validator from whitelist
type MsgUpdateCethReceiverAccount struct {
	CosmosSender        sdk.AccAddress `json:"cosmos_sender" yaml:"cosmos_sender"`
	CethReceiverAccount sdk.AccAddress `json:"ceth_receiver_account" yaml:"ceth_receiver_account"`
}

// NewMsgUpdateCethReceiverAccount is a constructor function for MsgUpdateCethReceiverAccount
func NewMsgUpdateCethReceiverAccount(cosmosSender sdk.AccAddress,
	cethReceiverAccount sdk.AccAddress) MsgUpdateCethReceiverAccount {
	return MsgUpdateCethReceiverAccount{
		CosmosSender:        cosmosSender,
		CethReceiverAccount: cethReceiverAccount,
	}
}

// Route should return the name of the module
func (msg MsgUpdateCethReceiverAccount) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateCethReceiverAccount) Type() string { return "update_ceth_receiver_account" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateCethReceiverAccount) ValidateBasic() error {
	if msg.CosmosSender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	if msg.CethReceiverAccount.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CethReceiverAccount.String())
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
	return []sdk.AccAddress{msg.CosmosSender}
}

// MsgRescueCeth transfer the ceth from ethbridge module to an account
type MsgRescueCeth struct {
	CosmosSender   sdk.AccAddress `json:"cosmos_sender" yaml:"cosmos_sender"`
	CosmosReceiver sdk.AccAddress `json:"cosmos_receiver" yaml:"cosmos_receiver"`
	CethAmount     sdk.Int        `json:"ceth_amount" yaml:"ceth_amount"`
}

// NewMsgRescueCeth is a constructor function for NewMsgRescueCeth
func NewMsgRescueCeth(cosmosSender sdk.AccAddress, cosmosReceiver sdk.AccAddress, cethAmount sdk.Int) MsgRescueCeth {
	return MsgRescueCeth{
		CosmosSender:   cosmosSender,
		CosmosReceiver: cosmosReceiver,
		CethAmount:     cethAmount,
	}
}

// Route should return the name of the module
func (msg MsgRescueCeth) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRescueCeth) Type() string { return "rescue_ceth" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRescueCeth) ValidateBasic() error {
	if msg.CosmosSender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}
	if msg.CosmosReceiver.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosReceiver.String())
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
	return []sdk.AccAddress{msg.CosmosSender}
}

// MsgUpdateWhiteListValidator add or remove validator from whitelist
type MsgUpdateWhiteListValidator struct {
	NetworkDescriptor uint32         `json:"network_descriptor" yaml:"network_descriptor"`
	CosmosSender      sdk.AccAddress `json:"cosmos_sender" yaml:"cosmos_sender"`
	Validator         sdk.ValAddress `json:"validator" yaml:"validator"`
	Power             uint32         `json:"power" yaml:"power"`
}

// NewMsgUpdateWhiteListValidator is a constructor function for MsgUpdateWhiteListValidator
func NewMsgUpdateWhiteListValidator(networkDescriptor uint32, cosmosSender sdk.AccAddress,
	validator sdk.ValAddress, power uint32) MsgUpdateWhiteListValidator {
	return MsgUpdateWhiteListValidator{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      cosmosSender,
		Validator:         validator,
		Power:             power,
	}
}

// Route should return the name of the module
func (msg MsgUpdateWhiteListValidator) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateWhiteListValidator) Type() string { return "update_whitelist_validator" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateWhiteListValidator) ValidateBasic() error {
	if msg.CosmosSender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	if msg.Validator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator.String())
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
	return []sdk.AccAddress{msg.CosmosSender}
}

// MapOracleClaimsToEthBridgeClaims maps a set of generic oracle claim data into EthBridgeClaim objects
func MapOracleClaimsToEthBridgeClaims(
	ethereumChainID int, bridgeContract EthereumAddress, nonce int, symbol string,
	tokenContract EthereumAddress, ethereumSender EthereumAddress,
	oracleValidatorClaims map[string]string,
	f func(int, EthereumAddress, int, EthereumAddress, sdk.ValAddress, string,
	) (EthBridgeClaim, error),
) ([]EthBridgeClaim, error) {
	mappedClaims := make([]EthBridgeClaim, len(oracleValidatorClaims))
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
