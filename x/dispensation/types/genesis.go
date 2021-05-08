package types

import (
	"encoding/json"
)

// NewGenesisState creates a new GenesisState instance
func NewGenesisState() GenesisState {
	return GenesisState{
		DistributionRecords: nil,
		Distributions:       nil,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		DistributionRecords: nil,
		Distributions:       nil,
	}
}

func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		ModuleCdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
