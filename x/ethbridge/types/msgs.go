package types

import (
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"strings"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

// NewMsgLock is a constructor function for MsgLock
func NewMsgLock(
	networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, symbol string, crossChainFee sdk.Int) MsgLock {
	return MsgLock{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      cosmosSender.String(),
		EthereumReceiver:  ethereumReceiver.String(),
		Amount:            amount,
		Symbol:            symbol,
		CrosschainFee:     crossChainFee,
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

// GetProphecyID get prophecy ID for lock message
func (msg MsgLock) GetProphecyID(doublePeggy bool, sequence, globalNonce uint64) []byte {
	return ComputeProphecyID(
		msg.CosmosSender,
		sequence,
		msg.EthereumReceiver,
		// TODO need get the token address from token's symbol
		msg.EthereumReceiver,
		msg.Amount,
		doublePeggy,
		globalNonce,
	)
}

// NewMsgBurn is a constructor function for MsgBurn
func NewMsgBurn(
	networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, symbol string, crosschainFee sdk.Int) MsgBurn {
	return MsgBurn{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      cosmosSender.String(),
		EthereumReceiver:  ethereumReceiver.String(),
		Amount:            amount,
		Symbol:            symbol,
		CrosschainFee:     crosschainFee,
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

// GetProphecyID get prophecy ID for lock message
func (msg MsgBurn) GetProphecyID(doublePeggy bool, sequence, globalNonce uint64) []byte {

	return ComputeProphecyID(
		msg.CosmosSender,
		sequence,
		msg.EthereumReceiver,
		// TODO need get the token address from token's symbol
		msg.EthereumReceiver,
		msg.Amount,
		doublePeggy,
		globalNonce,
	)
}

// ComputeProphecyID compute the prophecy id
func ComputeProphecyID(cosmosSender string, sequence uint64, ethereumReceiver string, tokenAddress string, amount sdk.Int,
	doublePeggy bool, globalNonce uint64) []byte {

	bytesTy, _ := abi.NewType("bytes", nil)
	boolTy, _ := abi.NewType("bool", nil)
	uint128Ty, _ := abi.NewType("uint128", nil)
	uint256Ty, _ := abi.NewType("uint256", nil)
	addressTy, _ := abi.NewType("address", nil)

	arguments := abi.Arguments{
		{
			Type: bytesTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: boolTy,
		},
		{
			Type: uint128Ty,
		},
	}

	bytes, _ := arguments.Pack(
		[]byte(cosmosSender),
		big.NewInt(int64(sequence)),

		gethCommon.HexToAddress(ethereumReceiver),
		// TODO need get the token address from token's symbol
		gethCommon.HexToAddress(ethereumReceiver),
		big.NewInt(amount.Int64()),
		doublePeggy,
		big.NewInt(int64(globalNonce)),
	)

	hashBytes := crypto.Keccak256(bytes)
	return hashBytes
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

// NewMsgUpdateCrossChainFeeReceiverAccount is a constructor function for MsgUpdateCrossChainFeeReceiverAccount
func NewMsgUpdateCrossChainFeeReceiverAccount(cosmosSender sdk.AccAddress,
	crosschainFeeReceiver sdk.AccAddress) MsgUpdateCrossChainFeeReceiverAccount {
	return MsgUpdateCrossChainFeeReceiverAccount{
		CosmosSender:          cosmosSender.String(),
		CrosschainFeeReceiver: crosschainFeeReceiver.String(),
	}
}

// Route should return the name of the module
func (msg MsgUpdateCrossChainFeeReceiverAccount) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateCrossChainFeeReceiverAccount) Type() string {
	return "update_crosschain_fee_receiver_account"
}

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateCrossChainFeeReceiverAccount) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.CrosschainFeeReceiver == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CrosschainFeeReceiver)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateCrossChainFeeReceiverAccount) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgUpdateCrossChainFeeReceiverAccount) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgRescueCrossChainFee is a constructor function for NewMsgRescueCrossChainFee
func NewMsgRescueCrossChainFee(cosmosSender sdk.AccAddress, cosmosReceiver sdk.AccAddress, crosschainFeeSymbol string, crosschainFee sdk.Int) MsgRescueCrossChainFee {
	return MsgRescueCrossChainFee{
		CosmosSender:        cosmosSender.String(),
		CosmosReceiver:      cosmosReceiver.String(),
		CrosschainFeeSymbol: crosschainFeeSymbol,
		CrosschainFee:       crosschainFee,
	}
}

// Route should return the name of the module
func (msg MsgRescueCrossChainFee) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRescueCrossChainFee) Type() string { return "rescue_crosschain_fee" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRescueCrossChainFee) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.CosmosReceiver == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosReceiver)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRescueCrossChainFee) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgRescueCrossChainFee) GetSigners() []sdk.AccAddress {
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

// NewMsgSetFeeInfo is a constructor function for MsgSetFeeInfo
func NewMsgSetFeeInfo(cosmosSender sdk.AccAddress, networkDescriptor oracletypes.NetworkDescriptor, feeCurrency string) MsgSetFeeInfo {
	return MsgSetFeeInfo{
		CosmosSender:      cosmosSender.String(),
		NetworkDescriptor: networkDescriptor,
		FeeCurrency:       feeCurrency,
	}
}

// Route should return the name of the module
func (msg MsgSignProphecy) Route() string { return RouterKey }

// NewMsgSignProphecy is a constructor function for MsgSignProphecy
func NewMsgSignProphecy(cosmosSender string, networkDescriptor oracletypes.NetworkDescriptor, prophecyID []byte, ethereumAddress, signature string) MsgSignProphecy {
	return MsgSignProphecy{
		CosmosSender:      cosmosSender,
		NetworkDescriptor: networkDescriptor,
		ProphecyId:        prophecyID,
		EthereumAddress:   ethereumAddress,
		Signature:         signature,
	}
}

// Type should return the action
func (msg MsgSignProphecy) Type() string { return "sign_prophecy" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSignProphecy) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if !msg.NetworkDescriptor.IsValid() {
		return errors.New("network descriptor is invalid")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSignProphecy) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSignProphecy) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// Route should return the name of the module
func (msg MsgSetFeeInfo) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetFeeInfo) Type() string { return "set_crosschain_fee_info" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetFeeInfo) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if !msg.NetworkDescriptor.IsValid() {
		return errors.New("network descriptor is invalid")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetFeeInfo) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgSetFeeInfo) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}
