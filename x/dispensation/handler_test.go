package dispensation_test

import (
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"testing"
)

func TestNewHandler_CreateDistribution(t *testing.T) {
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
	msgAirdrop := types.NewMsgCreateDistribution(sdk.AccAddress{}, "AR1", types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, inputList, outputList)
	res, err := handler(ctx, &msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)

	dr := append(keeper.GetRecordsForName(ctx, "AR1", types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords,
		keeper.GetRecordsForName(ctx, "AR1", types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED).DistributionRecords...)
	assert.Len(t, dr, recipients)
}

func TestNewHandler_CreateClaim(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	address := sdk.AccAddress(crypto.AddressHash([]byte("User1")))
	msgClaim := types.NewMsgCreateUserClaim(address, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	res, err := handler(ctx, &msgClaim)
	require.NoError(t, err)
	require.NotNil(t, res)

	cl, err := keeper.GetClaim(ctx, address.String(), types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	require.NoError(t, err)
	assert.False(t, cl.Locked)
	assert.Equal(t, cl.UserAddress, address.String())
}
