package types_test

import (
	"reflect"
	"testing"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
)

func TestTypes_ParamSetPairs(t *testing.T) {
	want := sdk.NewUint(10)

	p := types.Params{
		LeverageMax: want,
	}

	got := p.ParamSetPairs()[0]

	require.Equal(t, got.Key, []byte("LeverageMax"))
	require.Equal(t, got.Value, &want)
}

func TestTypes_ParamKeyTable(t *testing.T) {
	want := paramtypes.NewKeyTable().RegisterParamSet(&types.Params{})
	got := types.ParamKeyTable()

	reflect.DeepEqual(got, want)
}
