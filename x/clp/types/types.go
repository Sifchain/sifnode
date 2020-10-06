package types

import (
	"fmt"
	"strings"
)

type Asset struct {
	SourceChain string `json:"source_chain"`
	Symbol      string `json:"symbol"`
	Ticker      string `json:"ticker"`
}

// NewAsset returns a new Asset
func NewAsset(sourceChain string, symbol string, ticker string) Asset {
	return Asset{
		SourceChain: sourceChain,
		Symbol:      symbol,
		Ticker:      ticker,
	}
}

// implement fmt.Stringer
func (a Asset) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SourceChain: %s
Symbol: %s
Ticker: %s`, a.SourceChain, a.Symbol, a.Ticker))
}

type Pool struct {
	ExternalAsset        Asset  `json:"external_asset"`
	NativeAssetBalance   uint   `json:"native_asset_balance"`
	ExternalAssetBalance uint   `json:"external_asset_balance"`
	PoolUnits            uint   `json:"pool_units"`
	PoolAddress          string `json:"pool_address"`
}

func (p Pool) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ExternalAsset: %s
	NativeAssetBalance: %d
	NativeAssetBalance: %d
	PoolUnits : %s
	PoolAddress :%d`, p.ExternalAsset, p.ExternalAssetBalance, p.NativeAssetBalance, p.PoolAddress, p.PoolUnits))
}

// NewPool returns a new Pool
func NewPool(externalAsset Asset, nativeAssetBalance uint, externalAssetBalance uint, poolUnits uint, poolAddress string) Pool {
	return Pool{ExternalAsset: externalAsset, NativeAssetBalance: nativeAssetBalance, ExternalAssetBalance: externalAssetBalance, PoolUnits: poolUnits, PoolAddress: poolAddress}
}

type LiquidityProvider struct {
	Asset                    Asset `json:"asset"`
	LiquidityProviderUnits   uint  `json:"liquidity_provider_units"`
	LiquidityProviderAddress uint  `json:"liquidity_provider_address"`
}

func (l LiquidityProvider) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ExternalAsset: %s
	NativeAssetBalance: %d
	NativeAssetBalance: %d`, l.Asset, l.LiquidityProviderAddress, l.LiquidityProviderUnits))
}

// NewLiquidityProvider returns a new LiquidityProvider
func NewLiquidityProvider(asset Asset, liquidityProviderUnits uint, liquidityProviderAddress uint) LiquidityProvider {
	return LiquidityProvider{Asset: asset, LiquidityProviderUnits: liquidityProviderUnits, LiquidityProviderAddress: liquidityProviderAddress}
}
