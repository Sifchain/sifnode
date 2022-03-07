package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) Validate() bool {
	if !p.ExternalAsset.Validate() {
		return false
	}
	return true
}

// NewPool returns a new Pool
func NewPool(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint) Pool {
	pool := Pool{ExternalAsset: externalAsset,
		NativeAssetBalance:   nativeAssetBalance,
		ExternalAssetBalance: externalAssetBalance,
		PoolUnits:            poolUnits}

	return pool
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
func NewLiquidityProvider(asset *Asset, liquidityProviderUnits sdk.Uint, liquidityProviderAddress sdk.AccAddress) LiquidityProvider {
	return LiquidityProvider{Asset: asset, LiquidityProviderUnits: liquidityProviderUnits, LiquidityProviderAddress: liquidityProviderAddress.String()}
}

// ----------------------------------------------------------------------------
// Client Types

func NewLiquidityProviderResponse(liquidityProvider LiquidityProvider, height int64, nativeBalance string, externalBalance string) LiquidityProviderRes {
	return LiquidityProviderRes{LiquidityProvider: &liquidityProvider, Height: height, NativeAssetBalance: nativeBalance, ExternalAssetBalance: externalBalance}
}

func NewLiquidityProviderDataResponse(liquidityProviderData []*LiquidityProviderData, height int64) LiquidityProviderDataRes {
	return LiquidityProviderDataRes{LiquidityProviderData: liquidityProviderData, Height: height}
}

func NewLiquidityProviderData(liquidityProvider LiquidityProvider, nativeBalance string, externalBalance string) LiquidityProviderData {
	return LiquidityProviderData{LiquidityProvider: &liquidityProvider, NativeAssetBalance: nativeBalance, ExternalAssetBalance: externalBalance}
}

// ----------------------------------------------------------------------------
// Pmtp types
func CreatePmtpParams(data Params) PmtpRateParams {
	return PmtpRateParams{
		PmtpPeriodBlockRate:    sdk.Dec{},
		PmtpCurrentRunningRate: sdk.Dec{},
	}
}
