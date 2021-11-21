package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type Wrapper struct {
	Entries []RegistryEntryParser `json:"entries"`
}

type RegistryEntryParser struct {
	Decimals                 string   `json:"decimals"`
	Denom                    string   `json:"denom"`
	BaseDenom                string   `json:"base_denom"`
	Path                     string   `json:"path"`
	IbcChannelId             string   `json:"ibc_channel_id"`
	IbcCounterpartyChannelId string   `json:"ibc_counterparty_channel_id"`
	DisplayName              string   `json:"display_name"`
	DisplaySymbol            string   `json:"display_symbol"`
	Network                  string   `json:"network"`
	Address                  string   `json:"address"`
	ExternalSymbol           string   `json:"external_symbol"`
	TransferLimit            string   `json:"transfer_limit"`
	Permissions              []string `json:"permissions"`
	UnitDenom                string   `json:"unit_denom"`
	IbcCounterpartyDenom     string   `json:"ibc_counterparty_denom"`
	IbcCounterpartyChainId   string   `json:"ibc_counterparty_chain_id"`
}
type File struct {
	chain     string
	filenames []string
}

func main() {
	var inputs Wrapper
	files := []File{
		{
			chain:     "sifchain-1",
			filenames: []string{"tokenregistry", "registry"},
		},
		{
			chain:     "sifchain-devnet-1",
			filenames: []string{"tokenregistry", "registry"},
		},
		{
			chain:     "sifchain-testnet-1",
			filenames: []string{"tokenregistry", "registry"},
		},
	}

	basepath := "scripts/ibc/tokenregistration/"
	extension := ".json"
	for _, chain := range files {
		for _, filename := range chain.filenames {
			path := filepath.Join(basepath, chain.chain, filename)
			file, err := filepath.Abs(path + extension)
			if err != nil {
				panic(err)
			}
			input, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(input, &inputs)
			if err != nil {
				panic(err)
			}
			newReg := Migration(inputs.Entries)
			for _, e := range newReg {
				fmt.Println(e.Decimals, "  ", e.Permissions, " ", e.Denom)
			}
			updatedList := Wrapper{Entries: newReg}
			f, _ := json.MarshalIndent(updatedList, "", " ")
			outpath := filepath.Join(basepath, chain.chain, filename, "_updated")
			_ = ioutil.WriteFile(outpath+extension, f, 0600)
		}
	}
}

func Migration(entries []RegistryEntryParser) []RegistryEntryParser {
	newreg := make([]RegistryEntryParser, len(entries))
	for i, entry := range entries {
		dec, err := strconv.Atoi(entry.Decimals)
		if err != nil {
			panic(err)
		}
		if dec > 9 && CheckEntryPermissions(entry) {
			entry.Permissions = append(entry.Permissions, "IBCIMPORT")
			entry.IbcCounterpartyDenom = ""
			fmt.Println(entry)
		}
		newreg[i] = entry
	}
	//for _,e := range newreg{
	//	fmt.Println(e.Decimals, "  ",e.Permissions , " " ,e.Denom)
	//}
	return newreg
}

func CheckEntryPermissions(entry RegistryEntryParser) bool {
	requiredPermission := 2
	if len(entry.Permissions) != 2 {
		return false
	}
	for _, permission := range entry.Permissions {
		if permission == "IBCEXPORT" || permission == "CLP" {
			requiredPermission--
		}
	}
	return requiredPermission == 0
}
