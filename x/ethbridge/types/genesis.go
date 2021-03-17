package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all ethbridge state that must be provided at genesis
type GenesisState struct {
	CethAccount sdk.AccAddress `json:"ceth_account"`
}

// NewGenesisState creates a new GenesisState instance
func NewGenesisState() GenesisState {
	return GenesisState{
		CethAccount: sdk.AccAddress{},
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		CethAccount: sdk.AccAddress{},
	}
}

// GetGenesisStateFromAppState gets the GenesisState from raw message
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
