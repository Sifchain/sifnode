package types

import (
	"encoding/json"
	"fmt"
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

var Assets = []string{"rowan",
	"cusdt",
	"cusdc",
	"ccro",
	"cwbtc",
	"ceth",
	"cdai",
	"cyfi",
	"czrx",
	"cwscrt",
	"cwfil",
	"cuni",
	"cuma",
	"ctusd",
	"csxp",
	"csushi",
	"csusd",
	"csrm",
	"csnx",
	"csand",
	"crune",
	"creef",
	"cogn",
	"cocean",
	"cmana",
	"clrc",
	"clon",
	"clink",
	"ciotx",
	"cgrt",
	"cftm",
	"cesd",
	"cenj",
	"ccream",
	"ccomp",
	"ccocos",
	"cbond",
	"cbnt",
	"cbat",
	"cband",
	"cbal",
	"cant",
	"caave",
	"c1inch",
	"cleash",
	"cshib",
	"ctidal",
	"cpaid",
	"crndr",
	"cconv",
	"crfuel",
	"cakro",
	"cb20",
	"ctshp",
	"clina",
	"cdaofi",
	"ckeep",
	"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
	"ibc/6D717BFF5537D129035BAB39F593D638BA258A9F8D86FB7ECCEAB05B6950CC3E",
	"ibc/21CB41565FCA19AB6613EE06B0D56E588E0DC3E53FF94BA499BB9635794A1A35",
	"crly",
	"ibc/D87BC708A791246AA683D514C273736F07579CBD56C9CA79B7823F9A01C16270",
	"ibc/11DFDFADE34DCE439BA732EBA5CD8AA804A544BA1ECC0882856289FAF01FE53F",
	"ibc/B21954812E6E642ADC0B5ACB233E02A634BF137C572575BF80F7C0CC3DB2E74D",
	"ibc/2CC6F10253D563A7C238096BA63D060F7F356E37D5176E517034B8F730DB4AB6",
	"caxs",
	"cdfyn",
	"cdnxc",
	"cdon",
	"cern",
	"cfrax",
	"cfxs",
	"ckft",
	"cmatic",
	"cmetis",
	"cpols",
	"csaito",
	"ctoke",
	"czcn",
	"czcx",
	"cust",
	"cbtsg",
	"cquick",
	"cldo",
	"crail",
	"cpond",
	"cdino",
	"cufo",
	"ibc/F279AB967042CAC10BFF70FAECB179DCE37AAAE4CD4C1BC4565C2BBC383BC0FA",
	"ibc/C5C8682EB9AA1313EF1B12C991ADCDA465B80C05733BFB2972E2005E01BCE459",
	"ibc/B4314D0E670CB43C88A5DCA09F76E5E812BD831CC2FEC6E434C9E5A9D1F57953",
	"cratom",
	"cfis",
	"ibc/17F5C77854734CFE1301E6067AA42CDF62DAF836E4467C635E6DB407853C6082",
	"ibc/F141935FF02B74BDC6B8A0BD6FE86A23EE25D10E89AA0CD9158B3D92B63FDF4D",
	"ibc/ACA7D0100794F39DF3FF0C5E31638B24737321C24F32C2C486A24C78DD8F2029",
	"ibc/7B8A3357032F3DB000ACFF3B2C9F8E77B932F21004FC93B5A8F77DE24161A573",
	"coh",
	"ibc/7876FB1D317D993F1F54185DF6E405C7FE070B71E3A53AE0CEA5A86AC878EB7A",
	"ccsms",
	"clgcy",
	"ibc/3313DFB885C0C0EBE85E307A529985AFF7CA82239D404329BDF294E357FBC73A",
	"cmc",
	"cinj",
	"cpush",
	"cgala",
	"cosqth",
	"cnewo",
	"cuos",
	"cxft",
	"ibc/F20C4E30E4202C11FE009D6D58B2FF212C99084CB6F767287A51A93EFD960086",
	"ibc/57BB0CFF9782730595988FD330AA41605B0628E11507BABC1207B830A23493B9",
	"ibc/345D30E8ED06B47FC538ED131D99D16126F07CD6F8B35DE96AAF4C1E445AF466",
	"ibc/E46B030074825C99488BC57FD2DA711B0650FEF2BD24B61C228BBE3BCD73E69E",
	"ibc/7B1E1EFA6808065DA759354B6F21433156F4BF5DF2CF96DCBBC91738683748AF",
	"ibc/84506C652F91EA3742B9E00C4240BB039466DBAC48BD12872D2C1BA3FCFCA31E",
	"ibc/B650115F83DF4CA83E406A0ABDCE0BC284DC0B382DEFF634321D256FA8AFE2B9",
	"ibc/C8D8DAB01D770335E61A09D5468FBD6AEA080794AB4B866CCAFB7AD85DD270FB",
	"ibc/41139CF1224ADAF97D1E1466815F50E9BBF19A8C311B8331A333035DA938A5CF",
	"ccudos",
	"ibc/902CFB7D533886C25315A4602EB1938968565215A770E6E2EBA0842FC14A62C9",
	"ibc/37AE3DD9177BAD68DC2E39BD43FB1E70C3BB719FAEBAB42F0D03132A2E23A7BF",
}

func InitialRegistry() *Registry {
	var entries []*RegistryEntry
	for _, asset := range Assets {
		entries = append(entries, &RegistryEntry{
			Decimals:    18,
			Denom:       asset,
			Permissions: []Permission{Permission_CLP},
		})
	}
	return &Registry{Entries: entries}
}

//func InitialRegistry() *Registry {
//	entries := Registry{
//		Entries: []*RegistryEntry{
//			{Denom: "rowan", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "causc", Decimals: 6, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cusdt", Decimals: 6, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cusdc", Decimals: 6, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cwscrt", Decimals: 6, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cwbtc", Decimals: 8, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ceth", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cdai", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cyfi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "czrx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cwfil", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cuni", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cuma", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ctusd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "csxp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "csushi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "csusd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "csrm", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "csnx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "csand", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "crune", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "creef", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cogn", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cocean", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cmana", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "clrc", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "clon", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "clink", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ciotx", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cgrt", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cftm", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cesd", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cenj", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ccream", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ccomp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ccocos", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cbond", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cbnt", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cbat", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cband", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cbal", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cant", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "caave", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "c1inch", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cleash", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cshib", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ctidal", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cpaid", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "crndr", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cconv", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "crally", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "crfuel", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cakro", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cb20", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ctshp", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "clina", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "cdaofi", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//			{Denom: "ckeep", Decimals: 18, Permissions: []Permission{Permission_CLP}},
//		},
//	}
//	for i := range entries.Entries {
//		if !strings.HasPrefix(entries.Entries[i].Denom, "ibc/") {
//			entries.Entries[i].BaseDenom = entries.Entries[i].Denom
//		}
//	}
//	return &entries
//}
