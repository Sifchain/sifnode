package utils_test

import (
	"encoding/json"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
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
	config.Seal()
}

func createInput(t *testing.T, filename string) {
	in, err := sdk.AccAddressFromBech32("sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd")
	assert.NoError(t, err)
	out, err := sdk.AccAddressFromBech32("sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5")
	assert.NoError(t, err)
	coin := sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(10))}
	inputList := []bank.Input{bank.NewInput(in, coin), bank.NewInput(out, coin)}
	tempInput := utils.TempInput{In: inputList}
	file, _ := json.MarshalIndent(tempInput, "", " ")
	_ = ioutil.WriteFile(filename, file, 0600)
}

func createOutput(filename string, count int) {
	outputList := test.CreatOutputList(count, "10000000000000000000")
	tempInput := utils.TempOutput{Out: outputList}
	file, _ := json.MarshalIndent(tempInput, "", " ")
	_ = ioutil.WriteFile(filename, file, 0600)
}

func removeFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	assert.NoError(t, err)
}
func init() {
	SetConfig()
}
func TestParseInput(t *testing.T) {
	filename := "input.json"
	createInput(t, filename)
	defer removeFile(t, filename)
	inputs, err := utils.ParseInput(filename)
	assert.NoError(t, err)
	assert.Equal(t, len(inputs), 2)
}

func TestParseOutput(t *testing.T) {
	filename := "output.json"
	count := 3000
	createOutput(filename, count)
	//defer removeFile(t, filename)
	outputs, err := utils.ParseOutput(filename)
	assert.NoError(t, err)
	assert.Equal(t, len(outputs), count)
}
