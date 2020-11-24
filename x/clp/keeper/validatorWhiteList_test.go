package keeper_test

import (
	"github.com/Sifchain/sifnode/x/clp/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	V1 := test.GenerateValidatorAddress("A58856F0FD53BF058B4909A21AEC019107BA6")
	V2 := test.GenerateValidatorAddress("A58856F0FD53BF058B4909A21AEC019107BA7")
	ctx, keeper := test.CreateTestAppClp(false)
	keeper.SetValidatorWhiteList(ctx, []sdk.ValAddress{V1, V2})
	vList := keeper.GetValidatorWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, keeper.ExistsValidatorWhiteList(ctx))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	signer := test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA6")
	signer2 := test.GenerateAddress("A58856F0FD53BF058B4909A21AEC019107BA9")
	V1 := test.GenerateValidatorAddress("A58856F0FD53BF058B4909A21AEC019107BA6")
	V2 := test.GenerateValidatorAddress("A58856F0FD53BF058B4909A21AEC019107BA7")
	ctx, keeper := test.CreateTestAppClp(false)
	keeper.SetValidatorWhiteList(ctx, []sdk.ValAddress{V1, V2})
	vList := keeper.GetValidatorWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, keeper.ExistsValidatorWhiteList(ctx))
	assert.True(t, keeper.ValidateAddress(ctx, signer))
	assert.False(t, keeper.ValidateAddress(ctx, signer2))
}
