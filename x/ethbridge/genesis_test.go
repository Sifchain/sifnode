package ethbridge_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge"
	"github.com/Sifchain/sifnode/x/ethbridge/test"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestExportGenesis(t *testing.T) {
	ctx, keeper := test.CreateTestAppEthBridge(false)
	// Generate State
	tokenscount, receiver := CreateState(ctx, keeper, t)
	state := ethbridge.ExportGenesis(ctx, keeper)
	assert.Equal(t, len(state.PeggyTokens), tokenscount)
	assert.Equal(t, receiver, state.CrosschainFeeReceiveAccount)
}

func TestInitGenesis(t *testing.T) {
	ctx1, keeper1 := test.CreateTestAppEthBridge(false)
	ctx2, keeper2 := test.CreateTestAppEthBridge(false)
	// Generate State
	tokenscount, receiver := CreateState(ctx1, keeper1, t)
	state := ethbridge.ExportGenesis(ctx1, keeper1)
	assert.Equal(t, len(state.PeggyTokens), tokenscount)
	assert.Equal(t, state.CrosschainFeeReceiveAccount, receiver)
	state2 := ethbridge.ExportGenesis(ctx2, keeper2)
	assert.Equal(t, len(state2.PeggyTokens), 0)
	assert.Equal(t, state2.CrosschainFeeReceiveAccount, "")

	valUpdates := ethbridge.InitGenesis(ctx2, keeper2, *state)
	assert.Equal(t, len(valUpdates), 0)

	tokenslist := keeper2.GetPeggyToken(ctx2)
	assert.Equal(t, len(tokenslist.Tokens), tokenscount)
	actualReceiver := keeper2.GetCrossChainFeeReceiverAccount(ctx2)
	assert.Equal(t, receiver, actualReceiver.String())

}

func CreateState(ctx sdk.Context, keeper ethbridge.Keeper, t *testing.T) (int, string) {
	// Setting Tokens
	tokens := test.GenerateRandomTokens(10)
	for _, token := range tokens {
		keeper.AddPeggyToken(ctx, token)
	}
	gettokens := keeper.GetPeggyToken(ctx)
	assert.Greater(t, len(gettokens.Tokens), 0, "More than one token added")
	assert.LessOrEqual(t, len(gettokens.Tokens), len(tokens), "Add token will ignore duplicates")

	tokenscount := len(gettokens.Tokens)

	//Setting Receiver account
	receiver := test.GenerateAddress("")
	keeper.SetCrossChainFeeReceiverAccount(ctx, receiver)
	set := keeper.IsCrossChainFeeReceiverAccount(ctx, receiver)
	assert.True(t, set)

	return tokenscount, receiver.String()
}
