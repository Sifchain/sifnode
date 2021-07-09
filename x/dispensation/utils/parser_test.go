package utils_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	config.Seal()
}

func createInput(t *testing.T, filename string) {
	in, err := sdk.AccAddressFromBech32("sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd")
	assert.NoError(t, err)
	out, err := sdk.AccAddressFromBech32("sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5")
	assert.NoError(t, err)
	coin := sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(10))}
	inputList := []types.Input{types.NewInput(in, coin), types.NewInput(out, coin)}
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

// TODO Add the following utils as its own separate cmd

func TestAddressFilter(t *testing.T) {
	var addresStrings []string
	file, err := filepath.Abs("addrs.json")
	if err != nil {
		panic("Err getting filepath :" + err.Error())
	}
	o, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Err Reading file :" + err.Error())
	}
	err = json.Unmarshal(o, &addresStrings)
	if err != nil {
		panic("Err Unmarshall :" + err.Error())
	}
	for _, add := range addresStrings {
		_, err := sdk.AccAddressFromBech32(add)
		if err != nil {
			fmt.Println("Invalid :", add)
		}
	}
}

func TestSplitBetweenReciepients(t *testing.T) {
	t.Skip()
	type funders struct {
		Address           string
		PercentageFunding float64
		CalculatedAmount  sdk.Int
	}
	var investors []funders
	investors = append(investors, funders{
		Address:           "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
		PercentageFunding: 50.000,
	})
	investors = append(investors, funders{
		Address:           "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
		PercentageFunding: 50.000,
	})
	file := "expected_output.json"
	outputs, err := utils.ParseOutput(file)
	assert.NoError(t, err)
	total := sdk.ZeroDec()
	for _, out := range outputs {
		total = total.Add(out.Coins.AmountOf("rowan").ToDec())
	}
	inputList := make([]types.Input, len(investors))
	var totalPercentage float64
	for _, investor := range investors {
		totalPercentage = totalPercentage + investor.PercentageFunding
		percentage := sdk.NewDec(int64(investor.PercentageFunding))
		denom := sdk.NewDec(100)
		investor.CalculatedAmount = percentage.Quo(denom).Mul(total).TruncateInt()
		add, err := sdk.AccAddressFromBech32(investor.Address)
		assert.NoError(t, err)
		in := types.NewInput(add, sdk.Coins{sdk.NewCoin("rowan", investor.CalculatedAmount)})
		inputList = append(inputList, in)
	}
	assert.True(t, totalPercentage == 100.00, "Total Percentage is not 100%")
	tempInput := utils.TempInput{In: inputList}
	f, _ := json.MarshalIndent(tempInput, "", " ")
	_ = ioutil.WriteFile("input.json", f, 0600)
}
