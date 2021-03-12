package types

const (
	QueryBalance = "queryBalance"
)

func NewQueryReqGetFaucetBalance(denom string) QueryReqGetFaucetBalance {
	return QueryReqGetFaucetBalance{Denom: denom}
}
