package dispensation_test

import (
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewHandler(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 3000
	inputList := test.CreatInputList(2, "15000000000000000000000")
	outputList := test.CreatOutputList(recipients, "10000000000000000000")
	err := bank.ValidateInputsOutputs(inputList, outputList)
	assert.NoError(t, err)
	for _, in := range inputList {
		err := keeper.GetBankKeeper().AddCoins(ctx, sdk.AccAddress(in.Address), in.Coins)
		assert.NoError(t, err)
	}
	msgAirdrop := types.NewMsgDistribution(sdk.AccAddress{}, "AR1", types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, inputList, outputList)
	res, err := handler(ctx, &msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)

	dr := keeper.GetRecordsForNameAll(ctx, "AR1")
	assert.Len(t, dr, recipients)
}
