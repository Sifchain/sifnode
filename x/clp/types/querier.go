package types

const (
	QueryPool                  = "pool"
	QueryPools                 = "allpools"
	QueryAssetList             = "assetList"
	QueryLiquidityProvider     = "liquidityProvider"
	QueryLiquidityProviderData = "liquidityProviderData"
	QueryLPList                = "lpList"
	QueryAllLP                 = "allLp"
)

func NewQueryReqGetPool(symbol string) PoolReq {
	return PoolReq{Symbol: symbol}
}
