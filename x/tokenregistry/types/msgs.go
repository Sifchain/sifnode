package types

import (
	"encoding/json"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgRegister{}
var _ sdk.Msg = &MsgDeregister{}
var _ sdk.Msg = &MsgSetRegistry{}

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
	if m.Entry.Decimals <= 0 {
		return errors.New("Decimals cannot be zero")
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

// MsgRegisterAll

func (m *MsgRegisterAll) Route() string {
	return RouterKey
}

func (m *MsgRegisterAll) Type() string {
	return "register-all"
}

func (m *MsgRegisterAll) ValidateBasic() error {
	if m.Entries == nil {
		return errors.New("no token entry specified")
	}

	for _, entry := range m.Entries {

		if entry == nil {
			return errors.New("entry not initialized")
		}
		if entry.Denom == "" {
			return errors.New("no denom specified")
		}
		coin := sdk.Coin{
			Denom:  entry.Denom,
			Amount: sdk.OneInt(),
		}
		if !coin.IsValid() {
			return errors.New("Denom is not valid")
		}
		_, err := sdk.AccAddressFromBech32(m.From)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid from address")
		}
		if entry.Decimals <= 0 {
			return errors.New("Decimals cannot be zero")
		}
	}
	return nil
}

func (m *MsgRegisterAll) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRegisterAll) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// MsgSetRegistry

func (m *MsgSetRegistry) Route() string {
	return RouterKey
}

func (m *MsgSetRegistry) Type() string {
	return "set_registry"
}

func (m *MsgSetRegistry) ValidateBasic() error {
	if m.Registry == nil {
		return errors.New("no token entry specified")
	}

	// verify the entries
	for _, entry := range m.Registry.Entries {
		if entry.Denom == "" {
			return errors.New("no denom specified")
		}
		coin := sdk.Coin{
			Denom:  entry.Denom,
			Amount: sdk.OneInt(),
		}
		if !coin.IsValid() {
			return errors.New("Denom is not valid")
		}
		if entry.Decimals <= 0 {
			return errors.New("Decimals cannot be zero")
		}
	}

	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid from address")
	}
	return nil
}

func (m *MsgSetRegistry) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgSetRegistry) GetSigners() []sdk.AccAddress {
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

// MsgDeregisterAll

func (m *MsgDeregisterAll) Route() string {
	return RouterKey
}

func (m *MsgDeregisterAll) Type() string {
	return "deregister-akk"
}

func (m *MsgDeregisterAll) ValidateBasic() error {

	for _, denom := range m.Denoms {

		if denom == "" {
			return errors.New("no denom specified")
		}
		coin := sdk.Coin{
			Denom:  denom,
			Amount: sdk.OneInt(),
		}
		if !coin.IsValid() {
			return errors.New("Denom is not valid")
		}
		_, err := sdk.AccAddressFromBech32(m.From)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid from address")
		}
	}
	return nil
}

func (m *MsgDeregisterAll) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDeregisterAll) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewTokenMetadataAddRequest is a constructor function for TokenMetadataAddRequest
func NewTokenMetadataAddRequest(cosmosSender sdk.AccAddress,
	name string,
	symbol string,
	decimals int64,
	tokenAddress gethcommon.Address,
	networkDescriptor oracletypes.NetworkDescriptor) TokenMetadataAddRequest {
	return TokenMetadataAddRequest{
		CosmosSender: cosmosSender.String(),
		Metadata: &TokenMetadata{
			Name:              name,
			Symbol:            symbol,
			Decimals:          decimals,
			TokenAddress:      tokenAddress.String(),
			NetworkDescriptor: networkDescriptor,
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
