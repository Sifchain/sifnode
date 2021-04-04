package types

import (
	"encoding/json"
)

// NewGenesisState creates a new GenesisState instance
func NewGenesisState(params Params) GenesisState {

	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	admin := GetDefaultCLPAdmin()
	return &GenesisState{
		Params:           DefaultParams(),
		AddressWhitelist: []string{admin.String()},
	}
}

func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		json.Unmarshal(appState[ModuleName], &genesisState) // todo when module is migrated we need to use codec.JSONMarshler
	}
	return genesisState
}
