package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

// NewGenesisState creates a new GenesisState instance
func NewGenesisState() GenesisState {
	return *DefaultGenesisState()
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		ValidatorWhitelist:      []*GenesisValidatorWhiteList{},
		AdminAddress:            "",
		Prophecies:              []*Prophecy{},
		ConsensusNeeded:         []*GenesisConsensusNeeded{},
		CrossChainFee:           []*GenesisCrossChainFeeConfig{},
		WitnessLockBurnSequence: []*GenesisWitnessLockBurnSequence{},
		ProphecyInfo:            []*GenesisProphecyInfo{},
	}
}

// GetGenesisStateFromAppState gets the GenesisState from raw message
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := cdc.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic("Failed to get genesis state from app state")
		}
	}
	return genesisState
}
