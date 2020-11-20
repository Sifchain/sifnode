package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryPool              = "pool"
	QueryPools             = "allpools"
	QueryAssetList         = "assetList"
	QueryLiquidityProvider = "liquidityProvider"
)

type QueryReqGetPool struct {
	Ticker string `json:"ticker"`
}

func NewQueryReqGetPool(ticker string) QueryReqGetPool {
	return QueryReqGetPool{Ticker: ticker}
}

type QueryReqLiquidityProvider struct {
	Ticker    string         `json:"ticker"`
	LpAddress sdk.AccAddress `json:"lp_address"`
}

func NewQueryReqLiquidityProvider(ticker string, lpAddress sdk.AccAddress) QueryReqLiquidityProvider {
	return QueryReqLiquidityProvider{Ticker: ticker, LpAddress: lpAddress}
}

type QueryReqGetAssetList struct {
	LpAddress sdk.AccAddress `json:"lp_address"`
}

func NewQueryReqGetAssetList(lpAddress sdk.AccAddress) QueryReqGetAssetList {
	return QueryReqGetAssetList{LpAddress: lpAddress}

}
