package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRequestCoins{}
	_ sdk.Msg = &MsgAddCoins{}
)

type MsgRequestCoins struct {
	Coins     sdk.Coins
	Requester sdk.AccAddress
}

func NewMsgRequestCoins(requester sdk.AccAddress, coins sdk.Coins) MsgRequestCoins {
	return MsgRequestCoins{Requester: requester, Coins: coins}
}

func (msg MsgRequestCoins) Route() string {
	return RouterKey
}

func (msg MsgRequestCoins) Type() string {
	return "request_coins"
}

func (msg MsgRequestCoins) ValidateBasic() error {
	if msg.Requester.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Requester.String())
	}
	if !msg.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "Bids must be positive")
	}
	return nil
}

func (msg MsgRequestCoins) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}
func (msg MsgRequestCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Requester}
}

type MsgAddCoins struct {
	Signer sdk.AccAddress
	Coins  sdk.Coins
}

func NewMsgAddCoins(signer sdk.AccAddress, coins sdk.Coins) MsgAddCoins {
	return MsgAddCoins{Signer: signer, Coins: coins}
}

func (msg MsgAddCoins) Route() string {
	return RouterKey
}

func (msg MsgAddCoins) Type() string {
	return "add_coins"
}

func (msg MsgAddCoins) ValidateBasic() error {
	if msg.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer.String())
	}
	return nil
}

func (msg MsgAddCoins) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}
func (msg MsgAddCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
