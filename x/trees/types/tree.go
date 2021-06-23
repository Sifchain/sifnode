package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateTree{}
)

type Tree struct {
	Name     string         `json:"name"`
	Seller   sdk.AccAddress `json:"seller"`
	Price    sdk.Coins      `json:"price"`
	Category string         `json:"category"`
	Id       string         `json:"id"`
	Status   bool           `json:"status"`
}

type LimitOrder struct {
	Buyer    sdk.AccAddress `json:"buyer"`
	MaxPrice sdk.Coins      `json:"max_price"`
	TreeId   string         `json:"tree_id"`
	OrderId  string         `json:"order_id"`
	Executed bool           `json:"executed"`
}
