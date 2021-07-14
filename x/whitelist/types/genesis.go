package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

func UnmarshalGenesis(marshaler codec.JSONMarshaler, state json.RawMessage) GenesisState {
	var genesisState GenesisState
	if state != nil {
		err := marshaler.UnmarshalJSON(state, &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}

	return genesisState
}
