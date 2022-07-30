package clp_test

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp"
	"testing"
)

func Test_ReadPoolData(t *testing.T) {
	pools := clp.ReadPoolData()
	fmt.Println(pools)
}
