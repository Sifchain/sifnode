package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgDeleteUser{}

type MsgDeleteUser struct {
  ID      string         `json:"id" yaml:"id"`
  Creator sdk.AccAddress `json:"creator" yaml:"creator"`
}

func NewMsgDeleteUser(id string, creator sdk.AccAddress) MsgDeleteUser {
  return MsgDeleteUser{
    ID: id,
		Creator: creator,
	}
}

func (msg MsgDeleteUser) Route() string {
  return RouterKey
}

func (msg MsgDeleteUser) Type() string {
  return "DeleteUser"
}

func (msg MsgDeleteUser) GetSigners() []sdk.AccAddress {
  return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

func (msg MsgDeleteUser) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg MsgDeleteUser) ValidateBasic() error {
  if msg.Creator.Empty() {
    return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
  }
  return nil
}