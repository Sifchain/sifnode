package dispensation_test

import (
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewHandler(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	inputList := test.GenerateInputList("15000000000000000000")
	outputList := test.GenerateOutputList("10000000000000000000")

	for _, in := range inputList {
		_, err := keeper.GetBankKeeper().AddCoins(ctx, in.Address, in.Coins)
		assert.NoError(t, err)
	}
	msgAirdrop := types.NewMsgDistribution(sdk.AccAddress{}, "AR1", inputList, outputList)
	res, err := handler(ctx, msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, out := range outputList {
		ok := keeper.GetBankKeeper().HasCoins(ctx, out.Address, out.Coins)
		assert.True(t, ok)
	}
}
