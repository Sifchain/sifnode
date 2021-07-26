package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

func UnmarshalGenesis(marshaler codec.JSONMarshaler, state json.RawMessage) GenesisState {
	var genesisState GenesisState
	if state != nil {
		err := marshaler.UnmarshalJSON(state, &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}

	return genesisState
}

func GetGenesisStateFromAppState(marshaler codec.JSONMarshaler, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := marshaler.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	return genesisState
}

func DefaultWhitelist() DenomWhitelist {
	return DenomWhitelist{
		DenomWhitelistEntries: []*DenomWhitelistEntry{
			{Denom: "ccel", Decimals: 4},
			{Denom: "causc", Decimals: 6},
			{Denom: "cusdt", Decimals: 6},
			{Denom: "cusdc", Decimals: 6},
			{Denom: "ccro", Decimals: 8},
			{Denom: "ccdai", Decimals: 8},
			{Denom: "cwbtc", Decimals: 8},
			{Denom: "cceth", Decimals: 8},
			{Denom: "crenbtc", Decimals: 8},
			{Denom: "ccusdc", Decimals: 8},
			{Denom: "chusd", Decimals: 8},
			{Denom: "campl", Decimals: 9},
			{Denom: "ceth", Decimals: 18},
			{Denom: "cdai", Decimals: 18},
			{Denom: "cyfi", Decimals: 18},
			{Denom: "czrx", Decimals: 18},
			{Denom: "cwscrt", Decimals: 18},
			{Denom: "cwfil", Decimals: 18},
			{Denom: "cwbtc", Decimals: 18},
			{Denom: "cuni", Decimals: 18},
			{Denom: "cuma", Decimals: 18},
			{Denom: "ctusd", Decimals: 18},
			{Denom: "csxp", Decimals: 18},
			{Denom: "csushi", Decimals: 18},
			{Denom: "csusd", Decimals: 18},
			{Denom: "csrm", Decimals: 18},
			{Denom: "csnx", Decimals: 18},
			{Denom: "csand", Decimals: 18},
			{Denom: "crune", Decimals: 18},
			{Denom: "creef", Decimals: 18},
			{Denom: "cogn", Decimals: 18},
			{Denom: "cocean", Decimals: 18},
			{Denom: "cmana", Decimals: 18},
			{Denom: "clrc", Decimals: 18},
			{Denom: "clon", Decimals: 18},
			{Denom: "clink", Decimals: 18},
			{Denom: "ciotx", Decimals: 18},
			{Denom: "cgrt", Decimals: 18},
			{Denom: "cftm", Decimals: 18},
			{Denom: "cesd", Decimals: 18},
			{Denom: "cenj", Decimals: 18},
			{Denom: "ccream", Decimals: 18},
			{Denom: "ccomp", Decimals: 18},
			{Denom: "ccocos", Decimals: 18},
			{Denom: "cbond", Decimals: 18},
			{Denom: "cbnt", Decimals: 18},
			{Denom: "cbat", Decimals: 18},
			{Denom: "cband", Decimals: 18},
			{Denom: "cbal", Decimals: 18},
			{Denom: "cant", Decimals: 18},
			{Denom: "caave", Decimals: 18},
			{Denom: "c1inch", Decimals: 18},
		},
	}
}
