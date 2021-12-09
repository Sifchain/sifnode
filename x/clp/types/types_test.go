package types

// "fmt"

func Test_NewPool() {
	_ := NewPool(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_Validate() {
	_ := Validate(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewPool() {
	_ := NewPool(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewPoolsResponse() {
	_ := NewPoolsResponse(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewLiquidityProvider() {
	_ := NewLiquidityProvider(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewLiquidityProviderResponse() {
	_ := NewLiquidityProviderResponse(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewLiquidityProviderDataResponse() {
	_ := NewLiquidityProviderDataResponse(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}

func Test_NewLiquidityProviderData() {
	_ := NewLiquidityProviderData(externalAsset *Asset, nativeAssetBalance, externalAssetBalance, poolUnits sdk.Uint)
}
