package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Pool struct {
	ExternalAsset        Asset          `json:"external_asset"`
	NativeAssetBalance   uint           `json:"native_asset_balance"`
	ExternalAssetBalance uint           `json:"external_asset_balance"`
	PoolUnits            uint           `json:"pool_units"`
	PoolAddress          sdk.AccAddress `json:"pool_address"`
}

func (p Pool) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ExternalAsset: %s
	ExternalAssetBalance: %d
	NativeAssetBalance: %d
	PoolUnits : %d
	PoolAddress :%s`, p.ExternalAsset, p.ExternalAssetBalance, p.NativeAssetBalance, p.PoolUnits, p.PoolAddress))
}

func (p Pool) Validate() bool {
	if p.PoolAddress.Empty() {
		return false
	}
	if !p.ExternalAsset.Validate() {
		return false
	}
	return true
}

// NewPool returns a new Pool
func NewPool(externalAsset Asset, nativeAssetBalance uint, externalAssetBalance uint, poolUnits uint) (Pool, error) {
	pool := Pool{ExternalAsset: externalAsset,
		NativeAssetBalance:   nativeAssetBalance,
		ExternalAssetBalance: externalAssetBalance,
		PoolUnits:            poolUnits}
	nativeAsset := GetSettlementAsset()
	pooladdr, err := GetAddress(pool.ExternalAsset.Ticker, nativeAsset.Ticker)

	if err != nil {
		return Pool{}, err
	}
	pool.PoolAddress = pooladdr
	return pool, nil
}

type Pools []Pool

type LiquidityProvider struct {
	Asset                    Asset          `json:"asset"`
	LiquidityProviderUnits   uint           `json:"liquidity_provider_units"`
	LiquidityProviderAddress sdk.AccAddress `json:"liquidity_provider_address"`
}

func (l LiquidityProvider) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ExternalAsset: %s
	LiquidityProviderUnits: %d
	liquidityOroviderAddress: %s`, l.Asset, l.LiquidityProviderUnits, l.LiquidityProviderAddress))
}

func (l LiquidityProvider) Validate() bool {

	if !l.Asset.Validate() {
		return false
	}
	return true
}

// NewLiquidityProvider returns a new LiquidityProvider
func NewLiquidityProvider(asset Asset, liquidityProviderUnits uint, liquidityProviderAddress sdk.AccAddress) LiquidityProvider {
	return LiquidityProvider{Asset: asset, LiquidityProviderUnits: liquidityProviderUnits, LiquidityProviderAddress: liquidityProviderAddress}
}

// ----------------------------------------------------------------------------
// Client Types

type PoolResponse struct {
	Pool
	Height int64 `json:"height"`
}

func NewPoolResponse(pool Pool, height int64) PoolResponse {
	return PoolResponse{Pool: pool, Height: height}
}

type PoolsResponse struct {
	Pools
	Height int64 `json:"height"`
}

func NewPoolsResponse(pools Pools, height int64) PoolsResponse {
	return PoolsResponse{Pools: pools, Height: height}
}

type LiquidityProviderResponse struct {
	LiquidityProvider
	Height int64 `json:"height"`
}

func NewLiquidityProviderResponse(liquidityProvider LiquidityProvider, height int64) LiquidityProviderResponse {
	return LiquidityProviderResponse{LiquidityProvider: liquidityProvider, Height: height}
}
