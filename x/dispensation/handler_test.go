package dispensation_test

import (
<<<<<<< HEAD
	"testing"

	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
=======
	"fmt"
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
>>>>>>> develop
)

func TestNewHandler_CreateDistribution(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 3000
<<<<<<< HEAD
	inputList := test.CreateInputList(2, "15000000000000000000000")
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
=======
	outputList := test.CreatOutputList(recipients, "10000000000000000000")
	distibutor := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	_, err = keeper.GetBankKeeper().AddCoins(ctx, distibutor, totalCoins)
	assert.NoError(t, err)
	msgAirdrop := types.NewMsgDistribution(distibutor, types.Airdrop, outputList, sdk.AccAddress{})
	res, err := handler(ctx, msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)
	distributionName := fmt.Sprintf("%d_%s", ctx.BlockHeight(), msgAirdrop.Distributor.String())
	for _, e := range res.Events {
		if e.Type == "distribution_started" {
			assert.Len(t, e.Attributes, 3)
			assert.Contains(t, e.Attributes[1].String(), "distribution_name")
			assert.Contains(t, e.Attributes[1].String(), distributionName)
			assert.Contains(t, e.Attributes[2].String(), "distribution_type")
			assert.Contains(t, e.Attributes[2].String(), types.Airdrop.String())
		}
	}
	dr := keeper.GetRecordsForNameAll(ctx, distributionName)
	assert.Len(t, dr, recipients)
}

func TestNewHandler_CreateDistribution_MultipleTypes(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 3000
	outputList := test.CreatOutputList(recipients, "10000000000000000000")
	distibutor := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	_, err = keeper.GetBankKeeper().AddCoins(ctx, distibutor, totalCoins)
	assert.NoError(t, err)
	_, err = keeper.GetBankKeeper().AddCoins(ctx, distibutor, totalCoins)
	assert.NoError(t, err)
	msgAirdrop := types.NewMsgDistribution(distibutor, types.Airdrop, outputList, sdk.AccAddress{})
	res, err := handler(ctx, msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, msgAirdrop)
	require.Error(t, err)
	require.Nil(t, res)
	msgLm := types.NewMsgDistribution(distibutor, types.LiquidityMining, outputList, sdk.AccAddress{})
	res, err = handler(ctx, msgLm)
	require.NoError(t, err)
	require.NotNil(t, res)

>>>>>>> develop
}

func TestNewHandler_CreateClaim(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	address := sdk.AccAddress(crypto.AddressHash([]byte("User1")))
<<<<<<< HEAD
	msgClaim := types.NewMsgCreateUserClaim(address, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	res, err := handler(ctx, &msgClaim)
	require.NoError(t, err)
	require.NotNil(t, res)

	cl, err := keeper.GetClaim(ctx, address.String(), types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY)
	require.NoError(t, err)
	assert.False(t, cl.Locked)
	assert.Equal(t, cl.UserAddress, address.String())
=======
	msgClaim := types.NewMsgCreateClaim(address, types.ValidatorSubsidy)
	res, err := handler(ctx, msgClaim)
	require.NoError(t, err)
	require.NotNil(t, res)

	cl, err := keeper.GetClaim(ctx, address.String(), types.ValidatorSubsidy)
	require.NoError(t, err)
	assert.Equal(t, cl.UserAddress.String(), address.String())
}

func TestNewHandler_RunDistribution(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 3000
	outputList := test.CreatOutputList(recipients, "10000000000000000000")
	distributor := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	runner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	totalCoins = totalCoins.Add(totalCoins...)
	_, err = keeper.GetBankKeeper().AddCoins(ctx, distributor, totalCoins)
	assert.NoError(t, err)
	msgAirdrop := types.NewMsgDistribution(distributor, types.Airdrop, outputList, runner)
	res, err := handler(ctx, msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgLM := types.NewMsgDistribution(distributor, types.LiquidityMining, outputList, runner)
	res, err = handler(ctx, msgLM)
	require.NoError(t, err)
	require.NotNil(t, res)
	distributionName := fmt.Sprintf("%d_%s", ctx.BlockHeight(), msgAirdrop.Distributor.String())
	msgRun := types.NewMsgRunDistribution(runner, distributionName, types.Airdrop)
	res, err = handler(ctx, msgRun)
	require.NoError(t, err)
	require.NotNil(t, res)
	records := keeper.GetRecordsForNameCompleted(ctx, distributionName)
	assert.Len(t, records, 10)
	records = keeper.GetRecordsForNamePending(ctx, distributionName)
	assert.Len(t, records, (recipients*2)-10)
	msgRunFalse := types.NewMsgRunDistribution(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), distributionName, types.Airdrop)
	res, err = handler(ctx, msgRunFalse)
	require.NoError(t, err)
	require.NotNil(t, res)
>>>>>>> develop
}
