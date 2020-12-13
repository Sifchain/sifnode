package keeper_test

import (
	"log"
	"testing"

	oracleKeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	ctx, _, _, _, _, _ := oracleKeeper.CreateTestKeepers(t, 0.7, []int64{3, 7}, "")
	log.Printf("%v", ctx)

	// addresses := keeper.CreateTestAddrs(2)
	// keeper.SetOracleWhiteList(ctx, addresses[1])
	// vList := keeper.GetOracleWhiteList(ctx)
	// assert.Equal(t, len(vList), 2)
	// assert.True(t, keeper.ExistsOracleWhiteList(ctx))
}

// func TestKeeper_ValidateAddress(t *testing.T) {
// 	signer := test.GenerateAddress("")
// 	signer2 := test.GenerateAddress(test.AddressKey3)
// 	v1 := test.GenerateWhitelistAddress("")
// 	v2 := test.GenerateWhitelistAddress(test.AddressKey2)
// 	ctx, keeper := test.CreateTestAppOracle(false)
// 	keeper.SetOracleWhiteList(ctx, []sdk.AccAddress{v1, v2})
// 	vList := keeper.GetOracleWhiteList(ctx)
// 	assert.Equal(t, len(vList), 2)
// 	assert.True(t, keeper.ExistsOracleWhiteList(ctx))
// 	assert.True(t, keeper.ValidateAddress(ctx, signer))
// 	assert.False(t, keeper.ValidateAddress(ctx, signer2))
// }
