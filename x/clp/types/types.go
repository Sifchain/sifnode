package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) Validate() bool {
	return p.ExternalAsset.Validate()
}

// NewPool returns a new Pool
func NewPool(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint) Pool {
	pool := Pool{
		ExternalAsset:        externalAsset,
		NativeAssetBalance:   nativeAssetBalance,
		ExternalAssetBalance: externalAssetBalance,
		PoolUnits:            poolUnits,
	}

	return pool
}

func (p *Pool) ExtractValues(to Asset) (sdk.Uint, sdk.Uint, bool) {
	var X, Y sdk.Uint
	var toRowan bool

	if to.IsSettlementAsset() {
		Y = p.NativeAssetBalance
		X = p.ExternalAssetBalance
		toRowan = true
	} else {
		X = p.NativeAssetBalance
		Y = p.ExternalAssetBalance
		toRowan = false
	}

	return X, Y, toRowan
}

// UpdateBalances updates the balances of the pool
// TODO: It is unclear what the inputs are here.  We should document this function.
// TODO: DOCUMENT THIS
// TODO: CAPITAL LETTERS ARE GLOBAL VARIABLES.
func (p *Pool) UpdateBalances(toRowan bool, X, x, Y, swapResult sdk.Uint) { //nolint:gocritic
	if toRowan {
		p.ExternalAssetBalance = X.Add(x)
		p.NativeAssetBalance = Y.Sub(swapResult)
	} else {
		p.NativeAssetBalance = X.Add(x)
		p.ExternalAssetBalance = Y.Sub(swapResult)
	}
}

type (
	Pools              []Pool
	LiquidityProviders []LiquidityProvider
)

func (l LiquidityProvider) Validate() bool {
	return l.Asset.Validate()
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

func NewPmtpParamsResponse(params *PmtpParams, pmtpRateParams PmtpRateParams, pmtpEpoch PmtpEpoch, height int64) PmtpParamsRes {
	return PmtpParamsRes{Params: params, PmtpRateParams: &pmtpRateParams, PmtpEpoch: &pmtpEpoch, Height: height}
}

func NewLiquidityProtectionParamsResponse(params *LiquidityProtectionParams, rateParams LiquidityProtectionRateParams, height int64) LiquidityProtectionParamsRes {
	return LiquidityProtectionParamsRes{Params: params, RateParams: &rateParams, Height: height}
}

// todo: document
func (p *Pool) ExtractDebt(X, Y sdk.Uint, toRowan bool) (sdk.Uint, sdk.Uint) { //nolint:gocritic
	if toRowan {
		Y = Y.Add(p.NativeLiabilities)
		X = X.Add(p.ExternalLiabilities)
	} else {
		X = X.Add(p.NativeLiabilities)
		Y = Y.Add(p.ExternalLiabilities)
	}

	return X, Y
}

func StringCompare(a, b string) bool {
	return a == b
}
