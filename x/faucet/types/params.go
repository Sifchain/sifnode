package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
	// TODO: WE need a param for max withdraw amount
)

// Parameter store keys
var (
// KeyParamName          = []byte("ParamName")
)

// ParamKeyTable for faucet module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for faucet at genesis
type Params struct {

	// KeyParamName string `json:"key_param_name"`
}

// NewParams creates a new Params object
func NewParams() Params {
	return Params{}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`
	// TODO: Return all the params as a string
	`)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams()
}
