package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/oracle/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	v1 := test.GenerateWhitelistAddress("")
	v2 := test.GenerateWhitelistAddress(test.AddressKey2)
	ctx, keeper := test.CreateTestAppOracle(false)
	keeper.SetOracleWhiteList(ctx, []sdk.AccAddress{v1, v2})
	vList := keeper.GetOracleWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, keeper.ExistsOracleWhiteList(ctx))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	signer := test.GenerateAddress("")
	signer2 := test.GenerateAddress(test.AddressKey3)
	v1 := test.GenerateWhitelistAddress("")
	v2 := test.GenerateWhitelistAddress(test.AddressKey2)
	ctx, keeper := test.CreateTestAppOracle(false)
	keeper.SetOracleWhiteList(ctx, []sdk.AccAddress{v1, v2})
	vList := keeper.GetOracleWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, keeper.ExistsOracleWhiteList(ctx))
	assert.True(t, keeper.ValidateAddress(ctx, signer))
	assert.False(t, keeper.ValidateAddress(ctx, signer2))
}
