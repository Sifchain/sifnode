package types_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
)

func TestTypes_ParamSetPairs(t *testing.T) {
	p := types.Params{
		LeverageMax:          sdk.NewUint(10),
		InterestRateMax:      sdk.NewDec(5),
		InterestRateMin:      sdk.NewDec(1),
		InterestRateIncrease: sdk.NewDec(2),
		InterestRateDecrease: sdk.NewDec(3),
		HealthGainFactor:     sdk.NewDec(9),
		EpochLength:          45,
		ForceCloseThreshold:  sdk.NewDec(1), //TODO get real default
	}

	got := p.ParamSetPairs()

	paramSetPairsTest := []struct {
		key          string
		param        interface{}
		paramSetPair paramtypes.ParamSetPair
	}{
		{
			key:          "LeverageMax",
			param:        &p.LeverageMax,
			paramSetPair: got[0],
		},
		{
			key:          "InterestRateMax",
			param:        &p.InterestRateMax,
			paramSetPair: got[1],
		},
		{
			key:          "InterestRateMin",
			param:        &p.InterestRateMin,
			paramSetPair: got[2],
		},
		{
			key:          "InterestRateIncrease",
			param:        &p.InterestRateIncrease,
			paramSetPair: got[3],
		},
		{
			key:          "InterestRateDecrease",
			param:        &p.InterestRateDecrease,
			paramSetPair: got[4],
		},
		{
			key:          "HealthGainFactor",
			param:        &p.HealthGainFactor,
			paramSetPair: got[5],
		},
		{
			key:          "EpochLength",
			param:        &p.EpochLength,
			paramSetPair: got[6],
		},
		{
			key:          "ForceCloseThreshold",
			param:        &p.ForceCloseThreshold,
			paramSetPair: got[1],
		},
	}

	for _, tt := range paramSetPairsTest {
		tt := tt
		name := fmt.Sprintf("param %v", tt.key)
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.paramSetPair.Key, []byte(tt.key))
			require.Equal(t, tt.paramSetPair.Value, tt.param)
		})
	}
}

func TestTypes_ParamKeyTable(t *testing.T) {
	want := paramtypes.NewKeyTable().RegisterParamSet(&types.Params{})
	got := types.ParamKeyTable()

	reflect.DeepEqual(got, want)
}
