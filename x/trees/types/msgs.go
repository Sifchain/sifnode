package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreateTree{}
)

type MsgCreateTree struct {
	Name     string         `json:"name"`
	Seller   sdk.AccAddress `json:"seller"`
	Price    sdk.Coins      `json:"price"`
	Category string         `json:"category"`
}

func NewMsgCreateTree(name string, seller sdk.AccAddress, price sdk.Coins, category string) *MsgCreateTree {
	return &MsgCreateTree{Name: name, Seller: seller, Price: price, Category: category}
}

func (msg MsgCreateTree) Route() string {
	return ModuleName
}

func (msg MsgCreateTree) Type() string {
	return "trees"
}

func (msg MsgCreateTree) ValidateBasic() error {
	if msg.Seller.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Seller.String())
	}
	if msg.Price[0].Amount.Int64() <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Price must be greater than Zero")
	}
	_, err := sdk.AccAddressFromBech32(msg.Seller.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "Invalid Seller Address")
	}
	return nil
}

func (msg MsgCreateTree) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateTree) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Seller}
}

type MsgBuyTree struct {
	Buyer sdk.AccAddress `json:"buyer"`
	Id    string         `json:"id"`
	Price sdk.Coins      `json:"price"`
}

func NewMsgBuyTree(buyer sdk.AccAddress, id string, price sdk.Coins) MsgBuyTree {
	return MsgBuyTree{Buyer: buyer, Price: price, Id: id}
}

func (msg MsgBuyTree) Route() string {
	return ModuleName
}

func (msg MsgBuyTree) Type() string {
	return "trees"
}

func (msg MsgBuyTree) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Buyer.String())
	}
	if msg.Price[0].Amount.Int64() <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Price must be greater than Zero")
	}
	_, err := sdk.AccAddressFromBech32(msg.Buyer.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "Invalid Buyer Address")
	}
	return nil
}

func (msg MsgBuyTree) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

func (msg MsgBuyTree) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
