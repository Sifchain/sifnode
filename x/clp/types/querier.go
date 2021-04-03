package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

type QueryReqLiquidityProvider struct {
	Symbol    string         `json:"symbol"`
	LpAddress sdk.AccAddress `json:"lp_address"`
}

func NewQueryReqLiquidityProvider(symbol string, lpAddress sdk.AccAddress) QueryReqLiquidityProvider {
	return QueryReqLiquidityProvider{Symbol: symbol, LpAddress: lpAddress}
}

type QueryReqGetAssetList struct {
	LpAddress sdk.AccAddress `json:"lp_address"`
}

type QueryReqGetLiquidityProviderList struct {
	Symbol string `json:"symbol"`
}

func NewQueryReqGetLiquidityProviderList(symbol string) QueryReqGetLiquidityProviderList {
	return QueryReqGetLiquidityProviderList{Symbol: symbol}
}

func NewQueryReqGetAssetList(lpAddress sdk.AccAddress) QueryReqGetAssetList {
	return QueryReqGetAssetList{LpAddress: lpAddress}
}
