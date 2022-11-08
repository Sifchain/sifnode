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
	QueryPmtpParams            = "pmtpParams"
)

func NewQueryReqGetPool(symbol string) PoolReq {
	return PoolReq{Symbol: symbol}
}
