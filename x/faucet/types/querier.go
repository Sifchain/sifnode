package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Query endpoints supported by the faucet querier
const (
	QueryBalance = "queryBalance"
)

type QueryReqGetFaucetBalance struct {
	FaucetAddress sdk.AccAddress `json:"faucet_address"`
	Coins         sdk.Coins
}

func NewQueryReqGetFaucetBalance(faucetAddress sdk.AccAddress, coins sdk.Coins) QueryReqGetFaucetBalance {
	return QueryReqGetFaucetBalance{FaucetAddress: faucetAddress}

}
