package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Pool struct {
	ExternalAsset        Asset    `json:"external_asset"`
	NativeAssetBalance   sdk.Uint `json:"native_asset_balance"`
	ExternalAssetBalance sdk.Uint `json:"external_asset_balance"`
	PoolUnits            sdk.Uint `json:"pool_units"`
}

func (p Pool) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ExternalAsset: %s
	ExternalAssetBalance: %s
	NativeAssetBalance: %s
	PoolUnits : %s`, p.ExternalAsset, p.ExternalAssetBalance, p.NativeAssetBalance, p.PoolUnits))
}

func (p Pool) Validate() bool {
	if !p.ExternalAsset.Validate() {
		return false
	}
	return true
}

// NewPool returns a new Pool
func NewPool(externalAsset Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint) (Pool, error) {
	pool := Pool{ExternalAsset: externalAsset,
		NativeAssetBalance:   nativeAssetBalance,
		ExternalAssetBalance: externalAssetBalance,
		PoolUnits:            poolUnits}

	return pool, nil
}

type Pools []Pool
type LiquidityProviders []LiquidityProvider

type LiquidityProvider struct {
	Asset                    Asset          `json:"asset"`
	LiquidityProviderUnits   sdk.Uint       `json:"liquidity_provider_units"`
	LiquidityProviderAddress sdk.AccAddress `json:"liquidity_provider_address"`
}

func (l LiquidityProvider) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ExternalAsset: %s
	LiquidityProviderUnits: %s
	liquidityOroviderAddress: %s`, l.Asset, l.LiquidityProviderUnits, l.LiquidityProviderAddress))
}

func (l LiquidityProvider) Validate() bool {

	if !l.Asset.Validate() {
		return false
	}
	return true
}

// NewLiquidityProvider returns a new LiquidityProvider
func NewLiquidityProvider(asset Asset, liquidityProviderUnits sdk.Uint, liquidityProviderAddress sdk.AccAddress) LiquidityProvider {
	return LiquidityProvider{Asset: asset, LiquidityProviderUnits: liquidityProviderUnits, LiquidityProviderAddress: liquidityProviderAddress}
}

// ----------------------------------------------------------------------------
// Client Types

type PoolResponse struct {
	Pool
	ClpModuleAddress string `json:"clp_module_address"`
	Height           int64  `json:"height"`
}

func NewPoolResponse(pool Pool, height int64, address string) PoolResponse {
	return PoolResponse{Pool: pool, Height: height, ClpModuleAddress: address}
}

type PoolsResponse struct {
	Pools
	ClpModuleAddress string `json:"clp_module_address"`
	Height           int64  `json:"height"`
}

func NewPoolsResponse(pools Pools, height int64, address string) PoolsResponse {
	return PoolsResponse{Pools: pools, Height: height, ClpModuleAddress: address}
}

type LiquidityProviderResponse struct {
	LiquidityProvider
	NativeAssetBalance   string `json:"native_asset_balance"`
	ExternalAssetBalance string `json:"external_asset_balance"`
	Height               int64  `json:"height"`
}

func NewLiquidityProviderResponse(liquidityProvider LiquidityProvider, height int64, nativeBalance string, externalBalance string) LiquidityProviderResponse {
	return LiquidityProviderResponse{LiquidityProvider: liquidityProvider, Height: height, NativeAssetBalance: nativeBalance, ExternalAssetBalance: externalBalance}
}

type AssetListResponse struct {
	Assets
	Height int64 `json:"height"`
}

func NewAssetListResponse(assets Assets, height int64) AssetListResponse {
	return AssetListResponse{Assets: assets, Height: height}
}

type LpListResponse struct {
	LiquidityProviders
	Height int64 `json:"height"`
}

func NewLpListResponse(liquidityProviders LiquidityProviders, height int64) *LpListResponse {
	return &LpListResponse{LiquidityProviders: liquidityProviders, Height: height}
}
