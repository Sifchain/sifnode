package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all clp state that must be provided at genesis
//TODO: Add parameters to Genesis state ,such as minimum liquidity required to create a pool
type GenesisState struct {
	Params                Params             `json:"params" yaml:"params"`
	AddressWhitelist      []sdk.AccAddress   `json:"address_whitelist"`
	PoolList              Pools              `json:"pool_list"`
	LiquidityProviderList LiquidityProviders `json:"liquidity_provider_list"`
}

// NewGenesisState creates a new GenesisState instance
func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		json.Unmarshal(appState[ModuleName], &genesisState) // todo when module is migrated we need to use codec.JSONMarshler
	}
	return genesisState
}
