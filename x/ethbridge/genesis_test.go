package ethbridge

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestExportGenesis(t *testing.T) {
	ctx, keeper := types.CreateTestAppEthBridge(false)
	// Generate State
	tokenscount, receivercount := CreateState(ctx, keeper, t)
	state := ExportGenesis(ctx, keeper)
	assert.Equal(t, len(state.PeggyTokens), tokenscount)
	assert.Equal(t, len(state.CethReceiverAccount), receivercount)

}

func TestInitGenesis(t *testing.T) {
	ctx1, keeper1 := types.CreateTestAppEthBridge(false)
	ctx2, keeper2 := types.CreateTestAppEthBridge(false)
	// Generate State
	tokenscount, receivercount := CreateState(ctx1, keeper1, t)
	state := ExportGenesis(ctx1, keeper1)
	assert.Equal(t, len(state.PeggyTokens), tokenscount)
	assert.Equal(t, len(state.CethReceiverAccount), receivercount)
	state2 := ExportGenesis(ctx2, keeper2)
	assert.Equal(t, len(state2.PeggyTokens), 0)
	assert.Equal(t, len(state2.CethReceiverAccount), 0)

	valUpdates := InitGenesis(ctx2, keeper2, state)
	assert.Equal(t, len(valUpdates), 0)

	tokenslist := keeper2.GetPeggyToken(ctx2)
	assert.Equal(t, len(tokenslist), tokenscount)
	receiver := keeper2.GetCethReceiverAccount(ctx2)
	assert.Equal(t, len(receiver), receivercount)

}

func CreateState(ctx sdk.Context, keeper Keeper, t *testing.T) (int, int) {
	// Setting Tokens
	tokens := types.GenerateRandomTokens(10)
	for _, token := range tokens {
		err := keeper.AddPeggyToken(ctx, token)
		assert.NoError(t, err)
	}
	gettokens := keeper.GetPeggyToken(ctx)
	assert.Greater(t, len(gettokens), 0, "More than one token added")
	assert.LessOrEqual(t, len(gettokens), len(tokens), "Add token will ignore duplicates")

	tokenscount := len(gettokens)

	//Setting Receiver account
	receiver := types.GenerateAddress("")
	keeper.SetCethReceiverAccount(ctx, []sdk.AccAddress{receiver})
	set := keeper.IsCethReceiverAccount(ctx, []sdk.AccAddress{receiver})
	assert.True(t, set)

	return tokenscount, 1
}
