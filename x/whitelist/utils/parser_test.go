package utils

import (
	"encoding/json"
	"github.com/Sifchain/sifnode/x/whitelist/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func createInput(t *testing.T, filename string) {
	denomEntry := types.DenomWhitelistEntry{
		IsWhitelisted: true,
		Denom:         "ceth",
		Decimals:      18,
	}
	denomEntryList := []*types.DenomWhitelistEntry{&denomEntry}
	list := types.DenomWhitelist{DenomWhitelistEntries: denomEntryList}
	file, err := json.MarshalIndent(list, "", " ")
	assert.NoError(t, err)
	_ = ioutil.WriteFile(filename, file, 0600)
}

func TestParseDenoms(t *testing.T) {
	filepath := "denoms.json"
	defer os.Remove(filepath)
	createInput(t, filepath)
	whitelist, err := ParseDenoms(filepath)
	assert.NoError(t, err)
	assert.Len(t, whitelist.DenomWhitelistEntries, 1)
}
