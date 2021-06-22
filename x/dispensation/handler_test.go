package dispensation_test

import (
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
)

func TestNewHandler_CreateDistribution(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 3000
	outputList := test.CreatOutputList(recipients, "10000000000000000000")
	distributor := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	err = keeper.GetBankKeeper().AddCoins(ctx, distributor, totalCoins)
	assert.NoError(t, err)

	msgAirdrop := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, outputList, "")
	res, err := handler(ctx, &msgAirdrop)
	distributionName := fmt.Sprintf("%d_%s", ctx.BlockHeight(), msgAirdrop.Distributor)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, e := range res.Events {
		if e.Type == "distribution_started" {
			assert.Len(t, e.Attributes, 3)
			assert.Contains(t, e.Attributes[1].String(), "distribution_name")
			assert.Contains(t, e.Attributes[1].String(), distributionName)
			assert.Contains(t, e.Attributes[2].String(), "distribution_type")
			assert.Contains(t, e.Attributes[2].String(), types.DistributionType_DISTRIBUTION_TYPE_AIRDROP.String())
		}
	}
	dr := keeper.GetRecordsForName(ctx, distributionName)
	assert.Len(t, dr.DistributionRecords, recipients)
	dr = keeper.GetRecordsForNameAndStatus(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	assert.Len(t, dr.DistributionRecords, recipients)
}

func TestNewHandler_CreateDistribution_MultipleTypes(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 3000
	outputList := test.CreatOutputList(recipients, "10000000000000000000")
	distributor := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	err = keeper.GetBankKeeper().AddCoins(ctx, distributor, totalCoins)
	assert.NoError(t, err)
	err = keeper.GetBankKeeper().AddCoins(ctx, distributor, totalCoins)
	assert.NoError(t, err)
	msgAirdrop := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, outputList, "")
	res, err := handler(ctx, &msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgAirdrop)
	require.Error(t, err)
	require.Nil(t, res)
	msgLm := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, outputList, "")
	res, err = handler(ctx, &msgLm)
	require.NoError(t, err)
	require.NotNil(t, res)
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
	assert.Equal(t, cl.UserAddress, address.String())
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
	err = keeper.GetBankKeeper().AddCoins(ctx, distributor, totalCoins)
	assert.NoError(t, err)
	msgAirdrop := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, outputList, runner.String())
	res, err := handler(ctx, &msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)
	msgLM := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, outputList, runner.String())
	res, err = handler(ctx, &msgLM)
	require.NoError(t, err)
	require.NotNil(t, res)
	distributionName := fmt.Sprintf("%d_%s", ctx.BlockHeight(), msgAirdrop.Distributor)
	msgRun := types.NewMsgRunDistribution(runner.String(), distributionName, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
	res, err = handler(ctx, &msgRun)
	require.NoError(t, err)
	require.NotNil(t, res)
	records := keeper.GetRecordsForNameAndStatus(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	assert.Len(t, records.DistributionRecords, 10)
	records = keeper.GetRecordsForNameAndStatus(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	assert.Len(t, records.DistributionRecords, (recipients*2)-10)
	msgRunFalse := types.NewMsgRunDistribution(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(), distributionName, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP)
	res, err = handler(ctx, &msgRunFalse)
	require.NoError(t, err)
	require.NotNil(t, res)
}
