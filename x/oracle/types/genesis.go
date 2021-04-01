package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all clp state that must be provided at genesis
//TODO: Add parameters to Genesis state ,such as minimum liquidity required to create a pool
type GenesisState struct {
	AddressWhitelist []sdk.ValAddress `json:"address_whitelist"`
	AdminAddress     sdk.AccAddress   `json:"admin_address"`
}

// NewGenesisState creates a new GenesisState instance
func NewGenesisState() GenesisState {
	return GenesisState{
		AddressWhitelist: []sdk.ValAddress{},
		AdminAddress:     sdk.AccAddress{},
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		AddressWhitelist: []sdk.ValAddress{},
		AdminAddress:     sdk.AccAddress{},
	}
}

// GetGenesisStateFromAppState gets the GenesisState from raw message
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		json.Unmarshal(appState[ModuleName], &genesisState) // todo when module is migrated we need to use codec.JSONMarshler
	}
	return genesisState
}
