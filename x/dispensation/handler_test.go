package dispensation_test

import (
	"fmt"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"

	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
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
	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
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
	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)
	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
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

func TestNewHandler_CreateDistribution_PayRewardsInAnyToken_HappyCase(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 3
	outputList := test.CreatOutputList(recipients, "10")
	runner := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	distributor := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	totalCoins, err := dispensationUtils.TotalOutput(outputList)
	assert.NoError(t, err)
	totalCoins = totalCoins.Add(totalCoins...)
	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)
	// Airdrop distribution type
	msgAirdrop := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, outputList, runner.String())
	res, err := handler(ctx, &msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgAirdrop)
	require.Error(t, err)
	require.Nil(t, res)
	msgAirdropOutput := msgAirdrop.Output
	assert.Equal(t, recipients, len(msgAirdropOutput))
	for i := 0; i < len(msgAirdropOutput); i++ {
		//(testing) So users should get random catk or ceth coins here.
		assert.True(t, msgAirdropOutput[i].Coins.AmountOf("catk").Equal(sdk.NewInt(10)) ||
			msgAirdropOutput[i].Coins.AmountOf("ceth").Equal(sdk.NewInt(10)) || msgAirdropOutput[i].Coins.AmountOf("rowan").Equal(sdk.NewInt(10)))
	}
	msgLM := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, outputList, runner.String())
	res, err = handler(ctx, &msgLM)
	require.NoError(t, err)
	require.NotNil(t, res)
	distributionName := fmt.Sprintf("%d_%s", ctx.BlockHeight(), msgAirdrop.Distributor)
	msgRun := types.NewMsgRunDistribution(runner.String(), distributionName, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, 10)
	res, err = handler(ctx, &msgRun)
	for i := 0; i < len(outputList); i++ {
		lpAddress, _ := sdk.AccAddressFromBech32(outputList[i].Address)
		assert.True(t, keeper.GetBankKeeper().GetBalance(ctx, lpAddress, "ceth").Amount.Equal(sdk.NewInt(10)) ||
			keeper.GetBankKeeper().GetBalance(ctx, lpAddress, "catk").Amount.Equal(sdk.NewInt(10)) || keeper.GetBankKeeper().GetBalance(ctx, lpAddress, "rowan").Amount.Equal(sdk.NewInt(10)))
	}

	require.NoError(t, err)
	require.NotNil(t, res)
	records := keeper.GetRecordsForNameAndStatus(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	assert.Len(t, records.DistributionRecords, 3)
	records = keeper.GetRecordsForNameAndStatus(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	assert.Len(t, records.DistributionRecords, (recipients*2)-3)

	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)
	// Liquidity mining distribution type
	msgLiquidityMining := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, outputList, "")
	res, err = handler(ctx, &msgLiquidityMining)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgLiquidityMining)
	require.Error(t, err)
	require.Nil(t, res)
	msgLiquidityMiningOutput := msgLiquidityMining.Output
	assert.Equal(t, recipients, len(msgLiquidityMiningOutput))
	for i := 0; i < len(msgLiquidityMiningOutput); i++ {
		//(testing) So users should get random catk or ceth coins here.
		assert.True(t, msgLiquidityMiningOutput[i].Coins.AmountOf("catk").Equal(sdk.NewInt(10)) ||
			msgLiquidityMiningOutput[i].Coins.AmountOf("ceth").Equal(sdk.NewInt(10)) || msgLiquidityMiningOutput[i].Coins.AmountOf("rowan").Equal(sdk.NewInt(10)))
	}

	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)
	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)

	// Unspecified distribution type
	msgUnspecified := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, outputList, "")
	res, err = handler(ctx, &msgUnspecified)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgUnspecified)
	require.Error(t, err)
	require.Nil(t, res)
	msgUnspecifiedOutput := msgUnspecified.Output
	assert.Equal(t, recipients, len(msgUnspecifiedOutput))
	for i := 0; i < len(msgUnspecifiedOutput); i++ {
		//(testing) So users should get random catk or ceth coins here.
		assert.True(t, msgUnspecifiedOutput[i].Coins.AmountOf("catk").Equal(sdk.NewInt(10)) ||
			msgUnspecifiedOutput[i].Coins.AmountOf("ceth").Equal(sdk.NewInt(10)) || msgUnspecifiedOutput[i].Coins.AmountOf("rowan").Equal(sdk.NewInt(10)))
	}
	// Validator Subsidy distribution type
	msgValidatorSubsidy := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, outputList, "")
	res, err = handler(ctx, &msgValidatorSubsidy)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgValidatorSubsidy)
	require.Error(t, err)
	require.Nil(t, res)
	msgValidatorSubsidyOutput := msgValidatorSubsidy.Output
	assert.Equal(t, recipients, len(msgValidatorSubsidyOutput))
	for i := 0; i < len(msgValidatorSubsidyOutput); i++ {
		//(testing) So users should get random catk or ceth coins here.
		assert.True(t, msgValidatorSubsidyOutput[i].Coins.AmountOf("catk").Equal(sdk.NewInt(10)) ||
			msgValidatorSubsidyOutput[i].Coins.AmountOf("ceth").Equal(sdk.NewInt(10)) || msgValidatorSubsidyOutput[i].Coins.AmountOf("rowan").Equal(sdk.NewInt(10)))
	}
}

func TestNewHandler_CreateDistribution_PayRewardsInAnyToken_Error(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 2
	outputList := test.CreatOutputList(recipients, "10")
	distributor := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	_, err := dispensationUtils.TotalOutput(nil)
	assert.Error(t, err, "Outputlist is empty")
	totalCoins, _ := dispensationUtils.TotalOutput(outputList)
	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)
	msgAirdrop := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, outputList, "")
	res, err := handler(ctx, &msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgAirdrop)
	require.Error(t, err)
	require.Nil(t, res)
	msgAirdropOutput := msgAirdrop.Output
	assert.Equal(t, recipients, len(msgAirdropOutput))
	for i := 0; i < len(msgAirdropOutput); i++ {
		assert.True(t, msgAirdropOutput[i].Coins.AmountOf("catk").Equal(sdk.NewInt(10)) ||
			msgAirdropOutput[i].Coins.AmountOf("ceth").Equal(sdk.NewInt(10)) || msgAirdropOutput[i].Coins.AmountOf("rowan").Equal(sdk.NewInt(10)))
	}
	msgAirdrop = types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, outputList, "")
	_, err = handler(ctx, &msgAirdrop)
	assert.Error(t, err, "Failed in collecting funds")

	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)

	msgLiquidityMining := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, outputList, "")
	res, err = handler(ctx, &msgLiquidityMining)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgLiquidityMining)
	require.Error(t, err)
	require.Nil(t, res)
	msgLiquidityMiningOutput := msgLiquidityMining.Output
	assert.Equal(t, recipients, len(msgLiquidityMiningOutput))
	for i := 0; i < len(msgLiquidityMiningOutput); i++ {
		assert.True(t, msgLiquidityMiningOutput[i].Coins.AmountOf("catk").Equal(sdk.NewInt(10)) ||
			msgLiquidityMiningOutput[i].Coins.AmountOf("ceth").Equal(sdk.NewInt(10)) || msgLiquidityMiningOutput[i].Coins.AmountOf("rowan").Equal(sdk.NewInt(10)))
	}

	msgLiquidityMining = types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_LIQUIDITY_MINING, outputList, "")
	_, err = handler(ctx, &msgLiquidityMining)
	assert.Error(t, err, "Failed in collecting funds")

	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)

	msgUnspecified := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, outputList, "")
	res, err = handler(ctx, &msgUnspecified)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgUnspecified)
	require.Error(t, err)
	require.Nil(t, res)
	msgUnspecifiedOutput := msgUnspecified.Output
	assert.Equal(t, recipients, len(msgUnspecifiedOutput))
	for i := 0; i < len(msgUnspecifiedOutput); i++ {
		assert.True(t, msgUnspecifiedOutput[i].Coins.AmountOf("catk").Equal(sdk.NewInt(10)) ||
			msgUnspecifiedOutput[i].Coins.AmountOf("ceth").Equal(sdk.NewInt(10)) || msgUnspecifiedOutput[i].Coins.AmountOf("rowan").Equal(sdk.NewInt(10)))
	}

	msgUnspecified = types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED, outputList, "")
	_, err = handler(ctx, &msgUnspecified)
	assert.Error(t, err, "Failed in collecting funds")

	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
	assert.NoError(t, err)

	msgValidatorSubsidy := types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, outputList, "")
	res, err = handler(ctx, &msgValidatorSubsidy)
	require.NoError(t, err)
	require.NotNil(t, res)
	res, err = handler(ctx, &msgValidatorSubsidy)
	require.Error(t, err)
	require.Nil(t, res)
	msgValidatorSubsidyOutput := msgValidatorSubsidy.Output
	assert.Equal(t, recipients, len(msgValidatorSubsidyOutput))
	for i := 0; i < len(msgValidatorSubsidyOutput); i++ {
		assert.True(t, msgValidatorSubsidyOutput[i].Coins.AmountOf("catk").Equal(sdk.NewInt(10)) ||
			msgValidatorSubsidyOutput[i].Coins.AmountOf("ceth").Equal(sdk.NewInt(10)) || msgValidatorSubsidyOutput[i].Coins.AmountOf("rowan").Equal(sdk.NewInt(10)))
	}

	msgValidatorSubsidy = types.NewMsgCreateDistribution(distributor, types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY, outputList, "")
	_, err = handler(ctx, &msgValidatorSubsidy)
	assert.Error(t, err, "Failed in collecting funds")
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
	err = sifapp.AddCoinsToAccount(types.ModuleName, app.BankKeeper, ctx, distributor, totalCoins)
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
	msgRun := types.NewMsgRunDistribution(runner.String(), distributionName, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, 10)
	res, err = handler(ctx, &msgRun)
	require.NoError(t, err)
	require.NotNil(t, res)
	records := keeper.GetRecordsForNameAndStatus(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_COMPLETED)
	assert.Len(t, records.DistributionRecords, 10)
	records = keeper.GetRecordsForNameAndStatus(ctx, distributionName, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING)
	assert.Len(t, records.DistributionRecords, (recipients*2)-10)
	msgRunFalse := types.NewMsgRunDistribution(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(), distributionName, types.DistributionType_DISTRIBUTION_TYPE_AIRDROP, 10)
	res, err = handler(ctx, &msgRunFalse)
	require.NoError(t, err)
	require.NotNil(t, res)
}
