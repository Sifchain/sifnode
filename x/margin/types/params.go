package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyLeverageMaxParam = []byte("LeverageMax")
)

var _ paramtypes.ParamSet = (*Params)(nil)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLeverageMaxParam, &p.LeverageMax, validateLeverageMax),
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func validateLeverageMax(i interface{}) error {
	return nil
}
