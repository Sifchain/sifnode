package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

// NewGenesisState creates a new GenesisState instance
func NewGenesisState() GenesisState {
	return GenesisState{
		AddressWhitelist: []string{},
		AdminAddress:     "",
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		AddressWhitelist: []string{},
		AdminAddress:     "",
	}
}

// GetGenesisStateFromAppState gets the GenesisState from raw message
func GetGenesisStateFromAppState(cdc codec.Marshaler, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := cdc.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic("Failed to get genesis state from app state")
		}
	}
	return genesisState
}
