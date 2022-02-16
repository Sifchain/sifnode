package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultParamspace                    = ModuleName
	DefaultMinCreatePoolThreshold uint64 = 100
	DefaultPmtpStartBlock         int64  = 0
	DefaultPmtpEndBlock           int64  = 0
)

// Parameter store keys
var (
	KeyMinCreatePoolThreshold = []byte("MinCreatePoolThreshold")
	KeyPmtpNativeWeight       = []byte("PmtpNativeWeight")
	KeyPmtpExternalWeight     = []byte("PmtpExternalWeight")
	KeyPmtpStartBlock         = []byte("PmtpStartBlock")
	KeyPmtpEndBlock           = []byte("PmtpEndBlock")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable for clp module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(minThreshold uint64, pmtpNativeWeight, pmtpExternalWeight sdk.Dec, pmtpStartBlock, pmtpEndBlock int64) Params {
	return Params{
		MinCreatePoolThreshold: minThreshold,
		PmtpNativeWeight:       pmtpNativeWeight,
		PmtpExternalWeight:     pmtpExternalWeight,
		PmtpStartBlock:         pmtpStartBlock,
		PmtpEndBlock:           pmtpEndBlock,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinCreatePoolThreshold, &p.MinCreatePoolThreshold, validateMinCreatePoolThreshold),
		paramtypes.NewParamSetPair(KeyPmtpNativeWeight, &p.PmtpNativeWeight, validatePmtpNativeWeight),
		paramtypes.NewParamSetPair(KeyPmtpExternalWeight, &p.PmtpExternalWeight, validatePmtpExternalWeight),
		paramtypes.NewParamSetPair(KeyPmtpStartBlock, &p.PmtpStartBlock, validatePmtpStartBlock),
		paramtypes.NewParamSetPair(KeyPmtpEndBlock, &p.PmtpEndBlock, validatePmtpEndBlock),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return Params{
		MinCreatePoolThreshold: DefaultMinCreatePoolThreshold,
		PmtpNativeWeight:       sdk.NewDecWithPrec(5, 1),
		PmtpExternalWeight:     sdk.NewDecWithPrec(5, 1),
		PmtpStartBlock:         DefaultPmtpStartBlock,
		PmtpEndBlock:           DefaultPmtpEndBlock,
	}
}

func (p Params) Validate() error { // TODO determine all checks
	if err := validateMinCreatePoolThreshold(p.MinCreatePoolThreshold); err != nil {
		return err
	}
	if err := validatePmtpNativeWeight(p.PmtpNativeWeight); err != nil {
		return err
	}
	if err := validatePmtpExternalWeight(p.PmtpExternalWeight); err != nil {
		return err
	}
	if err := validatePmtpStartBlock(p.PmtpStartBlock); err != nil {
		return err
	}
	if err := validatePmtpEndBlock(p.PmtpEndBlock); err != nil {
		return err
	}
	if p.PmtpEndBlock < (p.PmtpStartBlock) {
		return fmt.Errorf(
			"end block (%d) must be after begin block (%d)",
			p.PmtpEndBlock, p.PmtpStartBlock,
		)
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

func validatePmtpNativeWeight(i interface{}) error { // TODO determine all checks
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("pmtp native weight threshold must be positive: %d", v)
	}
	return nil
}

func validatePmtpExternalWeight(i interface{}) error { // TODO determine all checks
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("pmtp external weight threshold must be positive: %d", v)
	}
	return nil
}

func validatePmtpStartBlock(i interface{}) error { // TODO determine all checks
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("pmtp start block cannot be negative: %d", v)
	}
	return nil
}

func validatePmtpEndBlock(i interface{}) error { // TODO determine all checks
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("pmtp end block cannot be negative: %d", v)
	}
	return nil
}

func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}
