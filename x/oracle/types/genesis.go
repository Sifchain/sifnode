package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all clp state that must be provided at genesis
//TODO: Add parameters to Genesis state ,such as minimum liquidity required to create a pool
type GenesisState struct {
	AddressWhitelist []sdk.ValAddress `json:"address_whitelist"`
}

// NewGenesisState creates a new GenesisState instance
func NewGenesisState() GenesisState {
	return GenesisState{
		AddressWhitelist: []sdk.ValAddress{},
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		AddressWhitelist: []sdk.ValAddress{},
	}
}

func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
