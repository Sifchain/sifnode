package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultMinCreatePoolThreshold uint64 = 100
	DefaultPmtpStartBlock         int64  = 11
	DefaultPmtpEndBlock           int64  = 72010
)

// Parameter store keys
var (
	KeyMinCreatePoolThreshold   = []byte("MinCreatePoolThreshold")
	KeyPmtpPeriodGovernanceRate = []byte("PmtpPeriodGovernanceRate")
	KeyPmtpEpochLength          = []byte("PmtpEpochLength")
	KeyPmtpStartBlock           = []byte("PmtpStartBlock")
	KeyPmtpEndBlock             = []byte("PmtpEndBlock")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable for clp module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(minThreshold uint64, pmtpGovernanaceRate sdk.Dec, pmtpEpochLength, pmtpStartBlock, pmtpEndBlock int64) Params {
	return Params{
		MinCreatePoolThreshold:   minThreshold,
		PmtpPeriodGovernanceRate: pmtpGovernanaceRate,
		PmtpPeriodEpochLength:    pmtpEpochLength,
		PmtpPeriodStartBlock:     pmtpStartBlock,
		PmtpPeriodEndBlock:       pmtpEndBlock,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinCreatePoolThreshold, &p.MinCreatePoolThreshold, validateMinCreatePoolThreshold),
		paramtypes.NewParamSetPair(KeyPmtpPeriodGovernanceRate, &p.PmtpPeriodGovernanceRate, validatePmtpPeriodGovernanceRate),
		paramtypes.NewParamSetPair(KeyPmtpEpochLength, &p.PmtpPeriodEpochLength, validatePmtpPeriodEpochLength),
		paramtypes.NewParamSetPair(KeyPmtpStartBlock, &p.PmtpPeriodStartBlock, validatePmtpStartBlock),
		paramtypes.NewParamSetPair(KeyPmtpEndBlock, &p.PmtpPeriodEndBlock, validatePmtpEndBlock),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return Params{
		MinCreatePoolThreshold:   DefaultMinCreatePoolThreshold,
		PmtpPeriodGovernanceRate: sdk.MustNewDecFromStr("0.10"),
		PmtpPeriodEpochLength:    14400,
		PmtpPeriodStartBlock:     DefaultPmtpStartBlock,
		PmtpPeriodEndBlock:       DefaultPmtpEndBlock,
	}
}

func (p Params) Validate() error {
	if err := validateMinCreatePoolThreshold(p.MinCreatePoolThreshold); err != nil {
		return err
	}
	if err := validatePmtpPeriodGovernanceRate(p.PmtpPeriodGovernanceRate); err != nil {
		return err
	}
	if err := validatePmtpPeriodEpochLength(p.PmtpPeriodEpochLength); err != nil {
		return err
	}
	if err := validatePmtpStartBlock(p.PmtpPeriodStartBlock); err != nil {
		return err
	}
	if err := validatePmtpEndBlock(p.PmtpPeriodEndBlock); err != nil {
		return err
	}
	if p.PmtpPeriodEndBlock <= p.PmtpPeriodStartBlock {
		return fmt.Errorf(
			"end block (%d) must be after begin block (%d)",
			p.PmtpPeriodEndBlock, p.PmtpPeriodStartBlock,
		)
	}
	if (p.PmtpPeriodEndBlock-p.PmtpPeriodStartBlock+1)%p.PmtpPeriodEpochLength != 0 {
		return fmt.Errorf("all epochs must have equal number of blocks : %d", p.PmtpPeriodEpochLength)
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

func validatePmtpPeriodGovernanceRate(i interface{}) error { // TODO determine all checks
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("pmtp governanace rate must be positive: %d", v)
	}
	return nil
}

func validatePmtpPeriodEpochLength(i interface{}) error { // TODO determine all checks
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("pmtp epoch length must be positive: %d", v)
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
