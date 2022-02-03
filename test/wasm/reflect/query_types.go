package reflect

type Text struct {
	Text string `json:"text"`
}

type reflectCustomQuery struct {
	Ping        *struct{} `json:"ping,omitempty"`
	Capitalized *Text     `json:"capitalized,omitempty"`
}

// this is from the go code back to the contract (capitalized or ping)
type reflectCustomQueryResponse struct {
	Msg string `json:"msg"`
}
