package keeper_test

import (
	"github.com/Sifchain/sifnode/x/clp/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	v1 := test.GenerateValidatorAddress("")
	v2 := test.GenerateValidatorAddress(test.AddressKey2)
	ctx, keeper := test.CreateTestAppClp(false)
	keeper.SetValidatorWhiteList(ctx, []sdk.ValAddress{v1, v2})
	vList := keeper.GetValidatorWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, keeper.ExistsValidatorWhiteList(ctx))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	signer := test.GenerateAddress("")
	signer2 := test.GenerateAddress(test.AddressKey3)
	v1 := test.GenerateValidatorAddress("")
	v2 := test.GenerateValidatorAddress(test.AddressKey2)
	ctx, keeper := test.CreateTestAppClp(false)
	keeper.SetValidatorWhiteList(ctx, []sdk.ValAddress{v1, v2})
	vList := keeper.GetValidatorWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, keeper.ExistsValidatorWhiteList(ctx))
	assert.True(t, keeper.ValidateAddress(ctx, signer))
	assert.False(t, keeper.ValidateAddress(ctx, signer2))
}
