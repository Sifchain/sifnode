package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
)

type GenesisState struct {
	DistributionRecords DistributionRecords `json:"distribution_records"`
	Distributions       Distributions       `json:"distributions"`
	Claims              UserClaims          `json:"claims"`
}

// NewGenesisState creates a new GenesisState instance
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
