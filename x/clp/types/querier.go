package types

const (
	QueryPool              = "pool"
	QueryPools             = "allpools"
	QueryLiquidityProvider = "liquidityProvider"
)

type QueryReqGetPool struct {
	Ticker      string `json:"ticker"`
	SourceChain string `json:"source_chain"`
}

type QueryReqLiquidityProvider struct {
	Ticker    string `json:"ticker"`
	LpAddress string `json:"lp_address"`
}
