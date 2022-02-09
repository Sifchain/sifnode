package wasm

type SifchainMsg struct {
	Swap *Swap `json:"swap,omitempty"`
}

type Swap struct {
	SentAsset         string `json:"sent_asset,omitempty"`
	ReceivedAssed     string `json:"received_asset,omitempty"`
	SentAmount        string `json:"sent_amount,omitempty"`
	MinReceivedAmount string `json:"min_received_amount,omitempty"`
}
