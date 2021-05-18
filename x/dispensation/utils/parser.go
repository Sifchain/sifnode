package utils

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"io/ioutil"
	"path/filepath"
)

type TempInput struct {
	In []bank.Input `json:"Input"`
}
type TempOutput struct {
	Out []bank.Output `json:"Output"`
}

func ParseInput(fp string) ([]bank.Input, error) {
	var inputs TempInput
	file, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(input, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs.In, nil
}

func ParseOutput(fp string) ([]bank.Output, error) {
	var outputs TempOutput
	file, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	o, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(o, &outputs)
	if err != nil {
		return nil, err
	}
	return outputs.Out, nil
}

func TotalOutput(output []bank.Output) sdk.Coins {
	total := output[0].Coins
	for i := 1; i < len(output); i++ {
		total = total.Add(output[i].Coins...)
	}
	return total
}
