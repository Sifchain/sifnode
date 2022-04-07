package types

const (
	QueryPool                  = "pool"
	QueryPools                 = "allpools"
	QueryAssetList             = "assetList"
	QueryLiquidityProvider     = "liquidityProvider"
	QueryLiquidityProviderData = "liquidityProviderData"
	QueryLPList                = "lpList"
	QueryAllLP                 = "allLp"
	QueryParams                = "params"
	QueryRewardParams          = "rewardParams"
)

func NewQueryReqGetPool(symbol string) PoolReq {
	return PoolReq{Symbol: symbol}
}
