package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	if m.Entry.Decimals < 0 {
		return errors.New("Decimals cannot be less than zero")
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
		if entry.Decimals < 0 {
			return errors.New("Decimals cannot be less than zero")
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
