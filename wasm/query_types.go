package wasm

type SifchainQuery struct {
	Ping *struct{} `json:"ping,omitempty"`
}

// this is from the go code back to the contract (capitalized or ping)
type SifchainQueryResponse struct {
	Msg string `json:"msg"`
}
