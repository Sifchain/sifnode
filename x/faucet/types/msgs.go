package faucet

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgRequestCoins struct {
	Coins     sdk.Coins
	Requester sdk.AccAddress
}

func (msg MsgRequestCoins) Name() string { return "request_coins" }
func (msg MsgRequestCoins) Type() string { return "faucet" }
func (msg MsgRequestCoins) ValidateBasic() sdk.Error {
	if msg.Requester.Empty() {
		return sdk.ErrInvalidAddress(msg.Requester.String())
	}
	if !msg.Coins.IsPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
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
