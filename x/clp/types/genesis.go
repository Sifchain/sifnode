package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
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
	admin := GetDefaultCLPAdmin()
	return GenesisState{
		Params:           DefaultParams(),
		AddressWhitelist: []sdk.AccAddress{admin},
	}
}

func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
