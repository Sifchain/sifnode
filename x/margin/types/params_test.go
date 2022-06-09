package types_test

import (
	"reflect"
	"testing"

	"github.com/Sifchain/sifnode/x/margin/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func TestTypes_ParamKeyTable(t *testing.T) {
	want := paramtypes.NewKeyTable().RegisterParamSet(&types.Params{})
	got := types.ParamKeyTable()

	reflect.DeepEqual(got, want)
}
