package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: How do I add the specified account to send from in the MsgRequestCoinsStruct?
type MsgRequestCoins struct {
	Coins     sdk.Coins
	Requester sdk.AccAddress
}

func (msg MsgRequestCoins) Name() string {
	return "request_coins"
}
func (msg MsgRequestCoins) Type() string {
	return "faucet"
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
