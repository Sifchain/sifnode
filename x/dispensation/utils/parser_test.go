package utils_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
)

const (
	AccountAddressPrefix = "sif"
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
)

func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	// config.Seal()
}

func createInput(t *testing.T, filename string) {
	in, err := sdk.AccAddressFromBech32("sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd")
	assert.NoError(t, err)
	out, err := sdk.AccAddressFromBech32("sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5")
	assert.NoError(t, err)
	coin := sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(10)))
	inputList := []types.Input{types.NewInput(in, coin), types.NewInput(out, coin)}
	tempInput := utils.TempInput{In: inputList}
	file, _ := json.MarshalIndent(tempInput, "", " ")
	_ = ioutil.WriteFile(filename, file, 0o600)
}

func createOutput(filename string, count int) {
	outputList := test.CreatOutputList(count, "10000000000000000000")
	tempInput := utils.TempOutput{Out: outputList}
	file, _ := json.MarshalIndent(tempInput, "", " ")
	_ = ioutil.WriteFile(filename, file, 0o600)
}

func removeFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	assert.NoError(t, err)
}

func init() {
	SetConfig()
}

func TestParseInput(t *testing.T) {
	file := "input.json"
	createInput(t, file)
	defer removeFile(t, file)
	inputs, err := utils.ParseInput(file)
	assert.NoError(t, err)
	assert.Equal(t, len(inputs), 2)
}

func TestParseOutput(t *testing.T) {
	file := "output.json"
	count := 3000
	createOutput(file, count)
	defer removeFile(t, file)
	outputs, err := utils.ParseOutput(file)
	assert.NoError(t, err)
	assert.Equal(t, len(outputs), count)
}
