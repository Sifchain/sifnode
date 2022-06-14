package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
)

func UnmarshalGenesis(marshaler codec.JSONCodec, state json.RawMessage) GenesisState {
	var genesisState GenesisState
	if state != nil {
		err := marshaler.UnmarshalJSON(state, &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	return genesisState
}
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Registry: InitialRegistry(),
	}
}

func GetGenesisStateFromAppState(marshaler codec.JSONCodec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := marshaler.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	return genesisState
}

func InitialRegistry() *Registry {
	entries := Registry{
		Entries: []*RegistryEntry{
			{Denom: "rowan", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "causc", Decimals: 6, Permissions: []Permission{Permission_CLP}},
			{Denom: "cusdt", Decimals: 6, Permissions: []Permission{Permission_CLP}},
			{Denom: "cusdc", Decimals: 6, Permissions: []Permission{Permission_CLP}},
			{Denom: "cwscrt", Decimals: 6, Permissions: []Permission{Permission_CLP}},
			{Denom: "cwbtc", Decimals: 8, Permissions: []Permission{Permission_CLP}},
			{Denom: "ceth", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cdai", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cyfi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "czrx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cwfil", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cuni", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cuma", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "ctusd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "csxp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "csushi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "csusd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "csrm", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "csnx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "csand", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "crune", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "creef", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cogn", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cocean", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cmana", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "clrc", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "clon", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "clink", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "ciotx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cgrt", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cftm", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cesd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cenj", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "ccream", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "ccomp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "ccocos", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cbond", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cbnt", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cbat", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cband", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cbal", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cant", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "caave", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "c1inch", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cleash", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cshib", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "ctidal", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cpaid", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "crndr", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cconv", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "crally", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "crfuel", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cakro", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cb20", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "ctshp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "clina", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "cdaofi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{Denom: "ckeep", Decimals: 18, Permissions: []Permission{Permission_CLP}},
		},
	}
	for i := range entries.Entries {
		if !strings.HasPrefix(entries.Entries[i].Denom, "ibc/") {
			entries.Entries[i].BaseDenom = entries.Entries[i].Denom
		}
	}
	return &entries
}
