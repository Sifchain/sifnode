package wasm

type SifchainMsg struct {
	Swap         *Swap         `json:"swap,omitempty"`
	AddLiquidity *AddLiquidity `json:"add_liquidity,omitempty"`
}

type Swap struct {
	SentAsset         string `json:"sent_asset,omitempty"`
	ReceivedAssed     string `json:"received_asset,omitempty"`
	SentAmount        string `json:"sent_amount,omitempty"`
	MinReceivedAmount string `json:"min_received_amount,omitempty"`
}

type AddLiquidity struct {
	ExternalAsset       string `json:"external_asset,omitempty"`
	NativeAssetAmount   string `json:"native_asset_amount,omitempty"`
	ExternalAssetAmount string `json:"external_asset_amount,omitempty"`
}
