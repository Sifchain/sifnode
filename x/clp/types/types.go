package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) Validate() bool {
	if !p.ExternalAsset.Validate() {
		return false
	}
	return true
}

// NewPool returns a new Pool
func NewPool(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint) (Pool, error) {
	pool := Pool{ExternalAsset: externalAsset,
		NativeAssetBalance:   nativeAssetBalance,
		ExternalAssetBalance: externalAssetBalance,
		PoolUnits:            poolUnits}

	return pool, nil
}

type Pools []Pool
type LiquidityProviders []LiquidityProvider

func (l LiquidityProvider) Validate() bool {

	if !l.Asset.Validate() {
		return false
	}
	return true
}

// NewLiquidityProvider returns a new LiquidityProvider
func NewLiquidityProvider(asset *Asset, liquidityProviderUnits sdk.Uint, liquidityProviderAddress fmt.Stringer) LiquidityProvider {
	return LiquidityProvider{Asset: asset, LiquidityProviderUnits: liquidityProviderUnits, LiquidityProviderAddress: liquidityProviderAddress.String()}
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
	Pools            []*Pool `json:"pools"`
	ClpModuleAddress string  `json:"clp_module_address"`
	Height           int64   `json:"height"`
}

func NewPoolsResponse(pools []*Pool, height int64, address string) PoolsResponse {
	return PoolsResponse{Pools: pools, Height: height, ClpModuleAddress: address}
}

func NewLiquidityProviderResponse(liquidityProvider LiquidityProvider, height int64, nativeBalance string, externalBalance string) LiquidityProviderRes {
	return LiquidityProviderRes{LiquidityProvider: &liquidityProvider, Height: height, NativeAssetBalance: nativeBalance, ExternalAssetBalance: externalBalance}
}

func NewAssetListResponse(assets []*Asset, height int64) AssetListRes {
	return AssetListRes{Assets: assets, Height: height}
}

type LpListResponse struct {
	LiquidityProviders
	Height int64 `json:"height"`
}

func NewLpListResponse(liquidityProviders LiquidityProviders, height int64) *LpListResponse {
	return &LpListResponse{LiquidityProviders: liquidityProviders, Height: height}
}
