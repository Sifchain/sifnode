package keeper_test

import (
	"github.com/Sifchain/sifnode/x/clp/test"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	v1 := test.GenerateWhitelistAddress("")
	v2 := test.GenerateWhitelistAddress(test.AddressKey2)
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	clpKeeper.SetClpWhiteList(ctx, []sdk.AccAddress{v1, v2})
	vList := clpKeeper.GetClpWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, clpKeeper.ExistsClpWhiteList(ctx))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	signer := test.GenerateAddress("")
	signer2 := test.GenerateAddress(test.AddressKey3)
	v1 := test.GenerateWhitelistAddress("")
	v2 := test.GenerateWhitelistAddress(test.AddressKey2)
	ctx, app := test.CreateTestAppClp(false)
	clpKeeper := app.ClpKeeper
	clpKeeper.SetClpWhiteList(ctx, []sdk.AccAddress{v1, v2})
	vList := clpKeeper.GetClpWhiteList(ctx)
	assert.Equal(t, len(vList), 2)
	assert.True(t, clpKeeper.ExistsClpWhiteList(ctx))
	assert.True(t, clpKeeper.ValidateAddress(ctx, signer))
	assert.False(t, clpKeeper.ValidateAddress(ctx, signer2))
}
