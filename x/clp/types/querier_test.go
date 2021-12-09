package types

import ( // "fmt"
	// "encoding/hex"
	"testing"
)

func Test_NewQueryReqGetPool(t *testing.T) {
	_ = PoolReq()

}
func Test_NewQueryReqLiquidityProvider(t *testing.T) {
	_ = NewQueryReqLiquidityProvider()
}

func Test_NewQueryReqGetAssetList(t *testing.T) {
	_ = NewQueryReqGetAssetList()
}

func Test_Equal(t *testing.T) {
	_ = NewQueryReqLiquidityProviderData()
}
