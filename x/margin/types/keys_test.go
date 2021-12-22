package types

import (
	"reflect"
	"testing"
)

func TestTypes_GetMTPKey(t *testing.T) {
	got := GetMTPKey("ceth", "xxx", "xxx")
	want := []byte{1, 99, 101, 116, 104, 120, 120, 120, 120, 120, 120}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
