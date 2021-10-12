package utils_test

import (
	"encoding/json"

	"github.com/Sifchain/sifnode/x/tokenregistry/test"
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	"github.com/Sifchain/sifnode/x/tokenregistry/utils"

	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createInput(t *testing.T, filename string) {
	denomEntry := types.RegistryEntry{
		Denom:    "ceth",
		Decimals: 18,
	}
	denomEntryList := []*types.RegistryEntry{&denomEntry}
	list := types.Registry{Entries: denomEntryList}
	file, err := json.MarshalIndent(list, "", " ")
	assert.NoError(t, err)
	_ = ioutil.WriteFile(filename, file, 0600)
}

func TestParseDenoms(t *testing.T) {
	app, _, _ := test.CreateTestApp(false)
	filepath := "denoms.json"
	defer os.Remove(filepath)
	createInput(t, filepath)
	whitelist, err := utils.ParseDenoms(app.AppCodec(), filepath)
	assert.NoError(t, err)
	assert.Len(t, whitelist.Entries, 1)
}
