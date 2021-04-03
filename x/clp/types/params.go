package types

import (
	"bytes"
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace                    = ModuleName
	DefaultMinCreatePoolThreshold uint64 = 100
)

// Parameter store keys
var (
	KeyMinCreatePoolThreshold = []byte("MinCreatePoolThreshold")
)
var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable for clp module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(minThreshold uint64) *Params {
	return &Params{
		MinCreatePoolThreshold: minThreshold,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinCreatePoolThreshold, &p.MinCreatePoolThreshold, validateMinCreatePoolThreshold),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() *Params {
	return NewParams(DefaultMinCreatePoolThreshold)
}

func (p Params) Validate() bool {
	if err := validateMinCreatePoolThreshold(p.MinCreatePoolThreshold); err != nil {
		return false
	}
	return true
}

func validateMinCreatePoolThreshold(i interface{}) error {
	v, ok := i.(uint)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("min create pool threshold must be positive: %d", v)
	}
	return nil
}

func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}
