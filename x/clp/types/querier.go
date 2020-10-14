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

func NewQueryReqGetPool(ticker string, sourceChain string) QueryReqGetPool {
	return QueryReqGetPool{Ticker: ticker, SourceChain: sourceChain}
}

type QueryReqLiquidityProvider struct {
	Ticker    string `json:"ticker"`
	LpAddress string `json:"lp_address"`
}

func NewQueryReqLiquidityProvider(ticker string, lpAddress string) QueryReqLiquidityProvider {
	return QueryReqLiquidityProvider{Ticker: ticker, LpAddress: lpAddress}
}
