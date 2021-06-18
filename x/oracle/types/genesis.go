package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	AddressWhitelist []sdk.ValAddress `json:"address_whitelist"`
	AdminAddress     sdk.AccAddress   `json:"admin_address"`
	Prophecies       []DBProphecy     `json:"prophecies"`
}

// NewGenesisState creates a default GenesisState instance
func NewGenesisState() GenesisState {
	return DefaultGenesisState()
}

// DefaultGenesisState creates default genesis for new chains
func DefaultGenesisState() GenesisState {
	return GenesisState{
		AdminAddress:     sdk.AccAddress{},
		AddressWhitelist: []sdk.ValAddress{},
		Prophecies:       []DBProphecy{},
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
