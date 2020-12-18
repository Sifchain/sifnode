package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Query endpoints supported by the faucet querier
const (
	QueryBalance = "queryBalance"
)

type QueryReqGetFaucetBalance struct {
	FaucetAddress sdk.AccAddress `json:"faucet_address"`
}

func NewQueryReqGetFaucetBalance(faucetAddress sdk.AccAddress) QueryReqGetFaucetBalance {
	return QueryReqGetFaucetBalance{FaucetAddress: faucetAddress}

}
