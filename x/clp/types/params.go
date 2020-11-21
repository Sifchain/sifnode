package types

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace                  = ModuleName
	DefaultMinCreatePoolThreshold uint = 100
)

// Parameter store keys
var (
	KeyMinCreatePoolThreshold = []byte("MinCreatePoolThreshold")
)
var _ params.ParamSet = (*Params)(nil)

// ParamKeyTable for clp module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for clp at genesis
type Params struct {
	MinCreatePoolThreshold uint `json:"min_create_pool_threshold"`
}

// NewParams creates a new Params object
func NewParams(minThreshold uint) Params {
	return Params{
		MinCreatePoolThreshold: minThreshold,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`
	MinCreatePoolThreshold : %d
	`, p.MinCreatePoolThreshold)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMinCreatePoolThreshold, &p.MinCreatePoolThreshold, validateMinCreatePoolThreshold),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
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
