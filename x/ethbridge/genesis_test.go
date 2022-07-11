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
	receiver := CreateState(ctx, keeper, t)
	state := ethbridge.ExportGenesis(ctx, keeper)
	assert.Equal(t, receiver, state.CrosschainFeeReceiveAccount)
}

func TestInitGenesis(t *testing.T) {
	ctx1, keeper1 := test.CreateTestAppEthBridge(false)
	ctx2, keeper2 := test.CreateTestAppEthBridge(false)
	// Generate State
	receiver := CreateState(ctx1, keeper1, t)
	state := ethbridge.ExportGenesis(ctx1, keeper1)
	assert.Equal(t, state.CrosschainFeeReceiveAccount, receiver)

	// not set for state2, so receiver account is empty
	state2 := ethbridge.ExportGenesis(ctx2, keeper2)
	assert.Equal(t, state2.CrosschainFeeReceiveAccount, "")

	// no validator update from the module
	valUpdates := ethbridge.InitGenesis(ctx2, keeper2, *state)
	assert.Equal(t, len(valUpdates), 0)

	// after init the genesis from state, receive account is set
	actualReceiver := keeper2.GetCrossChainFeeReceiverAccount(ctx2)
	assert.Equal(t, receiver, actualReceiver.String())

}

func TestValidateGenesis(t *testing.T) {
	ctx, keeper := test.CreateTestAppEthBridge(false)
	// export genesis but receiver not set
	state := ethbridge.ExportGenesis(ctx, keeper)
	// return err since receiver is empty, not valid cosmos address
	err := ethbridge.ValidateGenesis(*state)
	assert.NotEqual(t, err, nil)
}

func CreateState(ctx sdk.Context, keeper ethbridge.Keeper, t *testing.T) string {
	//Setting Receiver account
	receiver := test.GenerateAddress("")
	keeper.SetCrossChainFeeReceiverAccount(ctx, receiver)
	set := keeper.IsCrossChainFeeReceiverAccount(ctx, receiver)
	assert.True(t, set)

	return receiver.String()
}
