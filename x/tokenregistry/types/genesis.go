package types

import (
	"encoding/json"
	"fmt"
	"strings"

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

func InitialRegistry() Registry {
	entries := Registry{
		Entries: []*RegistryEntry{
			{IsWhitelisted: true, Denom: "ibc/287EE075B7AADDEB240AFE74FA2108CDACA50A7CCD013FA4C1FCD142AFA9CA9A", BaseDenom: "uphoton", IbcChannelId: "channel-0", IbcCounterpartyChannelId: "channel-86", Path: "transfer/channel-0", DisplayName: "uPHOTON", ExternalSymbol: "uPHOTON", Decimals: 6},
			{IsWhitelisted: true, Denom: "ibc/48E40290A494F271890BCFC867EB0940D8A6205DD94750C8EA71750480D65BA9", BaseDenom: "akt", IbcChannelId: "channel-1", IbcCounterpartyChannelId: "channel-12", Path: "transfer/channel-1", DisplayName: "AKT", ExternalSymbol: "AKT", Decimals: 6},
			{IsWhitelisted: true, Denom: "ibc/0F3C9D893A0ADE5738E473BB1A15C44D9715568E0C005D33A02495B444E15225", BaseDenom: "ncat", IbcChannelId: "channel-2", IbcCounterpartyChannelId: "channel-12", Path: "transfer/channel-2", DisplayName: "NCAT", ExternalSymbol: "NCAT", Decimals: 6},
			{IsWhitelisted: true, Denom: "ibc/E0B9629F3DF557C3412F12F6EFE3DACB28B4A30627A27697B6CFAD03A3DE0C96", BaseDenom: "dvpn", IbcChannelId: "channel-3", IbcCounterpartyChannelId: "channel-16", Path: "transfer/channel-3", DisplayName: "dVPN", ExternalSymbol: "dVPN", Decimals: 6},
			{IsWhitelisted: true, Denom: "rowan", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ccel", Decimals: 4, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "causc", Decimals: 6, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cusdt", Decimals: 6, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cusdc", Decimals: 6, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ccro", Decimals: 8, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ccdai", Decimals: 8, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cwbtc", Decimals: 8, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cceth", Decimals: 8, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "crenbtc", Decimals: 8, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ccusdc", Decimals: 8, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "chusd", Decimals: 8, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "campl", Decimals: 9, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ceth", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cdai", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cyfi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "czrx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cwscrt", Decimals: 6, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cwfil", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cuni", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cuma", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ctusd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "csxp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "csushi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "csusd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "csrm", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "csnx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "csand", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "crune", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "creef", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cogn", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cocean", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cmana", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "clrc", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "clon", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "clink", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ciotx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cgrt", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cftm", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cesd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cenj", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ccream", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ccomp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ccocos", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cbond", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cbnt", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cbat", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cband", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cbal", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cant", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "caave", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "c1inch", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cleash", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cshib", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ctidal", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cpaid", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "crndr", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cconv", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "crally", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "crfuel", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cakro", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cb20", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ctshp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "clina", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "cdaofi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
			{IsWhitelisted: true, Denom: "ckeep", Decimals: 18, Permissions: []Permission{Permission_CLP}},
		},
	}

	for i := range entries.Entries {
		if !strings.HasPrefix(entries.Entries[i].Denom, "ibc/") {
			entries.Entries[i].BaseDenom = entries.Entries[i].Denom
		}
	}

	return entries
}
