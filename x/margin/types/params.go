package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyLeverageMaxParam          = []byte("LeverageMax")
	KeyInterestRateMaxParam      = []byte("InterestRateMax")
	KeyInterestRateMinParam      = []byte("InterestRateMin")
	KeyInterestRateIncreaseParam = []byte("InterestRateIncrease")
	KeyInterestRateDecreaseParam = []byte("InterestRateDecrease")
	KeyHealthGainFactorParam     = []byte("HealthGainFactor")
	KeyEpochLengthParam          = []byte("EpochLength")
	KeyForceCloseThresholdParam  = []byte("ForceCloseThreshold")
	KeyPoolsParam                = []byte("Pools")
)

var _ paramtypes.ParamSet = (*Params)(nil)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLeverageMaxParam, &p.LeverageMax, validate),
		paramtypes.NewParamSetPair(KeyInterestRateMaxParam, &p.InterestRateMax, validate),
		paramtypes.NewParamSetPair(KeyInterestRateMinParam, &p.InterestRateMin, validate),
		paramtypes.NewParamSetPair(KeyInterestRateIncreaseParam, &p.InterestRateIncrease, validate),
		paramtypes.NewParamSetPair(KeyInterestRateDecreaseParam, &p.InterestRateDecrease, validate),
		paramtypes.NewParamSetPair(KeyHealthGainFactorParam, &p.HealthGainFactor, validate),
		paramtypes.NewParamSetPair(KeyEpochLengthParam, &p.EpochLength, validate),
		paramtypes.NewParamSetPair(KeyForceCloseThresholdParam, &p.ForceCloseThreshold, validate),
		paramtypes.NewParamSetPair(KeyPoolsParam, &p.Pools, validate),
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func validate(i interface{}) error {
	return nil
}
