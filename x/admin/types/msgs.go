package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

var _ sdk.Msg = &MsgAddAccount{}
var _ sdk.Msg = &MsgRemoveAccount{}
var _ sdk.Msg = &MsgSetParams{}
var _ legacytx.LegacyMsg = &MsgAddAccount{}
var _ legacytx.LegacyMsg = &MsgRemoveAccount{}
var _ legacytx.LegacyMsg = &MsgSetParams{}

func (m *MsgAddAccount) Route() string {
	return RouterKey
}

func (m *MsgAddAccount) Type() string {
	return "add_account"
}

func (m *MsgAddAccount) ValidateBasic() error {
	return nil
}

func (m *MsgAddAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgAddAccount) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgRemoveAccount) Route() string {
	return RouterKey
}

func (m *MsgRemoveAccount) Type() string {
	return "remove_account"
}

func (m *MsgRemoveAccount) ValidateBasic() error {
	return nil
}

func (m *MsgRemoveAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRemoveAccount) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgSetParams) Route() string {
	return RouterKey
}

func (m *MsgSetParams) Type() string {
	return "set_params"
}

func (m *MsgSetParams) ValidateBasic() error {
	return nil
}

func (m *MsgSetParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgSetParams) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
