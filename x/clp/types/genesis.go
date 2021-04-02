package types

import (
	"encoding/json"

<<<<<<< HEAD
	"github.com/cosmos/cosmos-sdk/codec"
=======
	sdk "github.com/cosmos/cosmos-sdk/types"
>>>>>>> marko/0.42
)

// NewGenesisState creates a new GenesisState instance
func NewGenesisState(params Params) GenesisState {

	return GenesisState{
		Params: &params,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
<<<<<<< HEAD
	params := DefaultParams()

	return GenesisState{
		Params: &params,
	}
}

func GetGenesisStateFromAppState(cdc codec.Marshaler, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {

		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
=======
	admin := GetDefaultCLPAdmin()
	return GenesisState{
		Params:           DefaultParams(),
		AddressWhitelist: []sdk.AccAddress{admin},
	}
}

func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		json.Unmarshal(appState[ModuleName], &genesisState) // todo when module is migrated we need to use codec.JSONMarshler
>>>>>>> marko/0.42
	}
	return genesisState
}
