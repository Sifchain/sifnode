package types_test

import (
	"reflect"
	"testing"

	"github.com/Sifchain/sifnode/x/margin/types"
)

func TestTypes_GetMTPKey(t *testing.T) {
	got := types.GetMTPKey("ceth", "xxx", "xxx", types.Position_LONG)
	want := []byte{1, 99, 101, 116, 104, 120, 120, 120, 120, 120, 120}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
