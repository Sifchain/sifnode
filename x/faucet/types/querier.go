package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryBalance = "queryBalance"
)

type QueryReqGetFaucetBalance struct {
	Coins sdk.Coins
}

func NewQueryReqGetFaucetBalance(coins sdk.Coins) QueryReqGetFaucetBalance {
	return QueryReqGetFaucetBalance{Coins: coins}
}
