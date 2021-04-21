package types

import (
	"fmt"
)

const (
	QueryPool              = "pool"
	QueryPools             = "allpools"
	QueryAssetList         = "assetList"
	QueryLiquidityProvider = "liquidityProvider"
	QueryLPList            = "lpList"
	QueryAllLP             = "allLp"
)

func NewQueryReqGetPool(symbol string) PoolReq {
	return PoolReq{Symbol: symbol}
}

func NewQueryReqLiquidityProvider(symbol string, lpAddress fmt.Stringer) LiquidityProviderReq {
	return LiquidityProviderReq{Symbol: symbol, LpAddress: lpAddress.String()}
}

func NewQueryReqGetAssetList(lpAddress fmt.Stringer) AssetListReq {
	return AssetListReq{LpAddress: lpAddress.String()}
}
