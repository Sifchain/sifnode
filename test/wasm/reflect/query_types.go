package reflect

type Text struct {
	Text string `json:"text"`
}

type ReflectCustomQuery struct {
	Ping        *struct{} `json:"ping,omitempty"`
	Capitalized *Text     `json:"capitalized,omitempty"`
}

// this is from the go code back to the contract (capitalized or ping)
type ReflectCustomQueryResponse struct {
	Msg string `json:"msg"`
}
