package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_GetDistributions(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	for i := 0; i < 10; i++ {
		name := uuid.New().String()
		distribution := types.NewDistribution(types.Airdrop, name, sdk.AccAddress{})
		err := keeper.SetDistribution(ctx, distribution)
		assert.NoError(t, err)
		res, err := keeper.GetDistribution(ctx, name, types.Airdrop)
		assert.NoError(t, err)
		assert.Equal(t, res.String(), distribution.String())
	}
	list := keeper.GetDistributions(ctx)
	assert.Len(t, list, 10)
}
