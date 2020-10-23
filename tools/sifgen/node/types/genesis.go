package types

import (
	"encoding/json"
)

type Value struct {
	Memo string `json:"memo"`
}

type Gentxs []Gentx
type Gentx struct {
	Type  string `json:"type"`
	Value Value  `json:"value"`
}

type Genutil struct {
	Gentxs Gentxs `json:"gentxs"`
}

type AppState struct {
	Genutil Genutil `json:"genutil"`
}

type GenesisAppState struct {
	AppState AppState `json:"app_state"`
}

type Result struct {
	Genesis *json.RawMessage `json:"genesis"`
}

type Genesis struct {
	Result Result `json:"result"`
}
