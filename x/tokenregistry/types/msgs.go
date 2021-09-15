package types

import (
	"encoding/json"

	oracletypes "github.com/Sifchain/sifnode/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgRegister{}
var _ sdk.Msg = &MsgDeregister{}

// MsgRegister

func (m *MsgRegister) Route() string {
	return RouterKey
}

func (m *MsgRegister) Type() string {
	return "register"
}

func (m *MsgRegister) ValidateBasic() error {
	if m.Entry == nil {
		return errors.New("no token entry specified")
	}

	if m.Entry.Denom == "" {
		return errors.New("no denom specified")
	}

	coin := sdk.Coin{
		Denom:  m.Entry.Denom,
		Amount: sdk.OneInt(),
	}
	if !coin.IsValid() {
		return errors.New("Denom is not valid")
	}

	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid from address")
	}

	if m.Entry.Decimals < 0 {
		return errors.New("Decimals cannot be negative")
	}

	return nil
}

func (m *MsgRegister) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRegister) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

// MsgDeregister

func (m *MsgDeregister) Route() string {
	return RouterKey
}

func (m *MsgDeregister) Type() string {
	return "deregister"
}

func (m *MsgDeregister) ValidateBasic() error {

	if m.Denom == "" {
		return errors.New("no denom specified")
	}

	coin := sdk.Coin{
		Denom:  m.Denom,
		Amount: sdk.OneInt(),
	}
	if !coin.IsValid() {
		return errors.New("Denom is not valid")
	}

	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid from address")
	}

	return nil
}

func (m *MsgDeregister) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDeregister) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

// NewTokenMetadata is a constructor function for MsgTokenMetadataAdd
func NewTokenMetadata(cosmosSender sdk.AccAddress,
	name string,
	symbol string,
	decimals int64,
	tokenAddress gethcommon.Address,
	network string) TokenMetadataAddRequest {
	return TokenMetadataAddRequest{
		CosmosSender: cosmosSender.String(),
		Metadata: &TokenMetadata{
			Name:         name,
			Symbol:       symbol,
			Decimals:     decimals,
			TokenAddress: tokenAddress.String(),
			Network:      network,
		},
	}
}

// Validate Basic runs stateless checks on the message
func (msg TokenMetadataAddRequest) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.Metadata.Name == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Metadata.Name)
	}

	if msg.Metadata.Symbol == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Metadata.Symbol)
	}

	if msg.Metadata.TokenAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Metadata.TokenAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg TokenMetadataAddRequest) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg TokenMetadataAddRequest) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// Route should return the name of the module
func (msg TokenMetadataAddRequest) Route() string { return RouterKey }

// Type should return the action
func (msg TokenMetadataAddRequest) Type() string { return "add_token_metadata" }

// NewTokenMetadata is a constructor function for MsgTokenMetadataAdd
func DeleteTokenMetadata(cosmosSender sdk.AccAddress, denomHash string) TokenMetadataDeleteRequest {
	return TokenMetadataDeleteRequest{
		CosmosSender: cosmosSender.String(),
		Denom:        denomHash,
	}
}

// Validate Basic runs stateless checks on the message
func (msg TokenMetadataDeleteRequest) ValidateBasic() error {
	if msg.CosmosSender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if msg.Denom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Denom)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg TokenMetadataDeleteRequest) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg TokenMetadataDeleteRequest) GetSigners() []sdk.AccAddress {
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{cosmosSender}
}

// Route should return the name of the module
func (msg TokenMetadataDeleteRequest) Route() string { return RouterKey }

// Type should return the action
func (msg TokenMetadataDeleteRequest) Type() string { return "delete_token_metadata" }

// Check if the token from Sifchain
func (m TokenMetadata) IsSifchain() bool {
	networkID := oracletypes.NetworkDescriptor_value[m.Network]
	networkDescriptor := oracletypes.NetworkDescriptor(networkID)
	return networkDescriptor.IsSifchain()
}
