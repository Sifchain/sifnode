package types

const (
	QueryBalance = "queryBalance"
)

type QueryReqGetFaucetBalance struct {
	Denom string
}

func NewQueryReqGetFaucetBalance(denom string) QueryReqGetFaucetBalance {
	return QueryReqGetFaucetBalance{Denom: denom}
}
