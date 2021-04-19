package utils

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"io/ioutil"
	"path/filepath"
)

type TempInput struct {
	In []types.Input `json:"Input"`
}
type TempOutput struct {
	Out []types.Output `json:"Output"`
}

func ParseInput(fp string) ([]types.Input, error) {
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

func ParseOutput(fp string) ([]types.Output, error) {
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
