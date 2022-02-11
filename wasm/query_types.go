package wasm

type SifchainQuery struct {
	Pool *PoolQuery `json:"pool,omitempty"`
}

type PoolQuery struct {
	ExternalAsset string `json:"external_asset,omitempty"`
}

type PoolResponse struct {
	ExternalAsset        string `json:"external_asset,omitempty"`
	ExternalAssetBalance string `json:"external_asset_balance,omitempty"`
	NativeAssetBalance   string `json:"native_asset_balance,omitempty"`
	PoolUnits            string `json:"pool_units,omitempty"`
}
