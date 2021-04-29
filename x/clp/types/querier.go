package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func NewQueryReqLiquidityProvider(symbol string, lpAddress sdk.AccAddress) LiquidityProviderReq {
	return LiquidityProviderReq{Symbol: symbol, LpAddress: lpAddress.String()}
}

func NewQueryReqGetAssetList(lpAddress sdk.AccAddress) AssetListReq {
	return AssetListReq{LpAddress: lpAddress.String()}
}
