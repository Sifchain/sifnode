package types

import (
	"bytes"
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultMinCreatePoolThreshold uint64 = 100
	DefaultEnableSwap             bool   = true
)

// Parameter store keys
var (
	KeyMinCreatePoolThreshold = []byte("MinCreatePoolThreshold")
	KeyEnableSwap             = []byte("EnableSwap")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable for clp module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(minThreshold uint64, enableSwap bool) Params {
	return Params{
		MinCreatePoolThreshold: minThreshold,
		EnableSwap:             enableSwap,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinCreatePoolThreshold, &p.MinCreatePoolThreshold, validateMinCreatePoolThreshold),
		paramtypes.NewParamSetPair(KeyEnableSwap, &p.EnableSwap, validateEnableSwap),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return Params{
		MinCreatePoolThreshold: DefaultMinCreatePoolThreshold,
		EnableSwap:             DefaultEnableSwap,
	}
}

func (p Params) Validate() error {
	if err := validateMinCreatePoolThreshold(p.MinCreatePoolThreshold); err != nil {
		return err
	}
	if err := validateEnableSwap(p.EnableSwap); err != nil {
		return err
	}
	return nil
}

func validateMinCreatePoolThreshold(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("min create pool threshold must be positive: %d", v)
	}
	return nil
}

func validateEnableSwap(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}
