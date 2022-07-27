package types

import (
	"encoding/json"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// NewMsgLock is a constructor function for MsgLock
func NewMsgLock(
	networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, denomHash string, crossChainFee sdk.Int) MsgLock {
	return MsgLock{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      cosmosSender.String(),
		EthereumReceiver:  ethereumReceiver.String(),
		Amount:            amount,
		DenomHash:         denomHash,
		CrosschainFee:     crossChainFee,
	}
}

// Route should return the name of the module
func (msg MsgLock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgLock) Type() string { return "lock" }

// ValidateNetworkDescriptor returns an error if the network type is out of the
// range we require (four base-10 digits)
func ValidateNetworkDescriptor(networkDescriptor oracletypes.NetworkDescriptor) error {
	if networkDescriptor < 0 || networkDescriptor > 9999 {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", networkDescriptor)
	}
	return nil
}

// ValidateBasic runs stateless checks on the message
func (msg MsgLock) ValidateBasic() error {
	err := ValidateNetworkDescriptor(msg.NetworkDescriptor)
	if err != nil {
		return err
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

	if len(msg.DenomHash) == 0 {
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
func (msg MsgLock) GetProphecyID(
	sequence uint64,
	tokenAddress string,
	tokenName string,
	tokenSymbol string,
	tokenDecimals uint8,
	bridgeToken bool,
	globalNonce uint64,
) []byte {
	return ComputeProphecyID(
		msg.CosmosSender,
		sequence,
		msg.EthereumReceiver,
		tokenAddress,
		msg.Amount,
		tokenName,
		tokenSymbol,
		tokenDecimals,
		msg.NetworkDescriptor,
		bridgeToken,
		globalNonce,
		msg.DenomHash,
	)
}

// NewMsgBurn is a constructor function for MsgBurn
func NewMsgBurn(
	networkDescriptor oracletypes.NetworkDescriptor, cosmosSender sdk.AccAddress,
	ethereumReceiver EthereumAddress, amount sdk.Int, denomHash string, crosschainFee sdk.Int) MsgBurn {
	return MsgBurn{
		NetworkDescriptor: networkDescriptor,
		CosmosSender:      cosmosSender.String(),
		EthereumReceiver:  ethereumReceiver.String(),
		Amount:            amount,
		DenomHash:         denomHash,
		CrosschainFee:     crosschainFee,
	}
}

// Route should return the name of the module
func (msg MsgBurn) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBurn) Type() string { return "burn" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBurn) ValidateBasic() error {
	err := ValidateNetworkDescriptor(msg.NetworkDescriptor)
	if err != nil {
		return err
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

	// prefixLength := len(PeggedCoinPrefix)
	// if len(msg.DenomHash) <= prefixLength+1 {
	// 	return ErrInvalidBurnSymbol
	// }

	// symbolPrefix := msg.DenomHash[:prefixLength]
	// if symbolPrefix != PeggedCoinPrefix {
	// 	return ErrInvalidBurnSymbol
	// }

	// symbolSuffix := msg.DenomHash[prefixLength:]
	// if len(symbolSuffix) == 0 {
	// 	return ErrInvalidBurnSymbol
	// }

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

// GetProphecyID get prophecy ID for burn message
func (msg MsgBurn) GetProphecyID(
	sequence uint64,
	tokenAddress string,
	tokenName string,
	tokenSymbol string,
	tokenDecimals uint8,
	bridgeToken bool,
	globalNonce uint64,
) []byte {
	return ComputeProphecyID(
		msg.CosmosSender,
		sequence,
		msg.EthereumReceiver,
		tokenAddress,
		msg.Amount,
		tokenName,
		tokenSymbol,
		tokenDecimals,
		msg.NetworkDescriptor,
		bridgeToken,
		globalNonce,
		msg.DenomHash,
	)
}

// ComputeProphecyID compute the prophecy id
func ComputeProphecyID(
	cosmosSender string,
	sequence uint64,
	ethereumReceiver string,
	tokenAddress string,
	amount sdk.Int,
	tokenName string,
	tokenSymbol string,
	tokenDecimals uint8,
	networkDescriptor oracletypes.NetworkDescriptor,
	bridgeToken bool,
	globalNonce uint64,
	cosmosDenom string,
) []byte {

	bytesTy, _ := abi.NewType("bytes", "bytes", nil)
	boolTy, _ := abi.NewType("bool", "bool", nil)
	uint8Ty, _ := abi.NewType("uint8", "uint8", nil)
	int32Ty, _ := abi.NewType("int32", "int32", nil)
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	addressTy, _ := abi.NewType("address", "address", nil)
	stringTy, _ := abi.NewType("string", "string", nil)
	uint128Ty, _ := abi.NewType("uint128", "uint128", nil)

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
			Type: stringTy,
		},
		{
			Type: stringTy,
		},
		{
			Type: uint8Ty,
		},
		{
			Type: int32Ty,
		},
		{
			Type: boolTy,
		},
		{
			Type: uint128Ty,
		},
		{
			Type: stringTy,
		},
	}

	bytes, _ := arguments.Pack(
		[]byte(cosmosSender),
		big.NewInt(int64(sequence)),
		gethCommon.HexToAddress(ethereumReceiver),
		gethCommon.HexToAddress(tokenAddress),
		big.NewInt(amount.Int64()),
		tokenName,
		tokenSymbol,
		tokenDecimals,
		networkDescriptor,
		bridgeToken,
		big.NewInt(int64(globalNonce)),
		cosmosDenom,
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
	if err := ValidateNetworkDescriptor(msg.EthBridgeClaim.NetworkDescriptor); err != nil {
		return err
	}

	if msg.EthBridgeClaim.CosmosReceiver == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.EthBridgeClaim.CosmosReceiver)
	}

	if msg.EthBridgeClaim.ValidatorAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.EthBridgeClaim.ValidatorAddress)
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
func NewMsgSetFeeInfo(cosmosSender sdk.AccAddress,
	networkDescriptor oracletypes.NetworkDescriptor,
	feeCurrency string,
	feeCurrencyGas sdk.Int,
	minimumLockCost sdk.Int,
	minimumBurnCost sdk.Int) MsgSetFeeInfo {
	return MsgSetFeeInfo{
		CosmosSender:      cosmosSender.String(),
		NetworkDescriptor: networkDescriptor,
		FeeCurrency:       feeCurrency,
		FeeCurrencyGas:    feeCurrencyGas,
		MinimumLockCost:   minimumLockCost,
		MinimumBurnCost:   minimumBurnCost,
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

	if err := ValidateNetworkDescriptor(msg.NetworkDescriptor); err != nil {
		return err
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
	cosmosSender, err := sdk.ValAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sdk.AccAddress(cosmosSender)}
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

	if err := ValidateNetworkDescriptor(msg.NetworkDescriptor); err != nil {
		return err
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

// Route should return the name of the module
func (msg MsgUpdateConsensusNeeded) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUpdateConsensusNeeded) Type() string { return "update_consensus_needed" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateConsensusNeeded) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if err := ValidateNetworkDescriptor(msg.NetworkDescriptor); err != nil {
		return err
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUpdateConsensusNeeded) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgUpdateConsensusNeeded) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// NewMsgUpdateConsensusNeeded is a constructor function for MsgUpdateConsensusNeeded
func NewMsgUpdateConsensusNeeded(cosmosSender string, networkDescriptor oracletypes.NetworkDescriptor, consensusNeeded uint32) MsgUpdateConsensusNeeded {
	return MsgUpdateConsensusNeeded{
		CosmosSender:      cosmosSender,
		NetworkDescriptor: networkDescriptor,
		ConsensusNeeded:   consensusNeeded,
	}
}

var _ sdk.Msg = &MsgSetBlacklist{}

// Route should return the name of the module
func (msg MsgSetBlacklist) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetBlacklist) Type() string { return "setblacklist" }

func (msg MsgSetBlacklist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg *MsgSetBlacklist) ValidateBasic() error {
	return nil
}

func (msg *MsgSetBlacklist) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
