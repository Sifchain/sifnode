package keeper

import (
	"encoding/json"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	cosmosReceivers, _                         = CreateTestAddrs(1)
	amount               int64                 = 10
	symbol                                     = "stake"
	tokenContractAddress                       = types.NewEthereumAddress("0xbbbbca6a901c926f240b89eacb641d8aec7aeafd")
	ethBridgeAddress     types.EthereumAddress = types.NewEthereumAddress(strings.ToLower("0x30753E4A8aad7F8597332E813735Def5dD395028"))
	ethereumSender                             = types.NewEthereumAddress("0x627306090abaB3A6e1400e9345bC60c78a8BEf57")
	//BadValidatorAddress                        = sdk.ValAddress(CreateTestPubKeys(1)[0].Address().Bytes())
)

func TestProcessClaimLock(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]

	nonce := 1
	//invalid claim defaults to lock
	claimType, err := types.StringToClaimType("lkfjdsk")
	require.Equal(t, claimType.String(), "lock")
	require.Error(t, err)

	claimType, err = types.StringToClaimType("lock")
	require.NoError(t, err)

	ethBridgeClaim := types.NewEthBridgeClaim(
		5777,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,
		cosmosReceivers[0],
		validator1Pow3,
		amount,
		claimType,
	)

	status, err := keeper.ProcessClaim(ctx, ethBridgeClaim)

	require.NoError(t, err)
	require.Equal(t, status.Text, oracle.PendingStatusText)
	// duplicate execution
	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))

	// other validator execute

	ethBridgeClaim = types.NewEthBridgeClaim(
		5777,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,       // accounts[0]
		cosmosReceivers[0],
		validator2Pow3,
		amount,
		claimType,
	)
	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, oracle.SuccessStatusText)

}

func TestProcessClaimBurn(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]

	nonce := 1

	claimType, err := types.StringToClaimType("burn")
	require.NoError(t, err)

	ethBridgeClaim := types.NewEthBridgeClaim(
		5777,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,
		cosmosReceivers[0],
		validator1Pow3,
		amount,
		claimType,
	)

	status, err := keeper.ProcessClaim(ctx, ethBridgeClaim)

	require.NoError(t, err)
	require.Equal(t, status.Text, oracle.PendingStatusText)

	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))

	// other validator execute

	ethBridgeClaim = types.NewEthBridgeClaim(
		5777,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,       // accounts[0]
		cosmosReceivers[0],
		validator2Pow3,
		amount,
		claimType,
	)
	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.NoError(t, err)
	require.Equal(t, status.Text, oracle.SuccessStatusText)

}
func TestProcessSuccessfulClaimLock(t *testing.T) {
	ctx, keeper, bankKeeper, _, _, _ := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	receiverCoins := bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins, sdk.Coins{})

	claimType, err := types.StringToClaimType("lock")
	require.NoError(t, err)
	claimContent := types.NewOracleClaimContent(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

	claimBytes, err := json.Marshal(claimContent)
	require.NoError(t, err)
	claimString := string(claimBytes)
	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	require.NoError(t, err)

	receiverCoins = bankKeeper.GetCoins(ctx, cosmosReceivers[0])

	require.Equal(t, receiverCoins.String(), "10cstake")

	// duplicate processSuccessClaim
	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	require.NoError(t, err)

	receiverCoins = bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, "20cstake", receiverCoins.String())
}

func TestProcessSuccessfulClaimBurn(t *testing.T) {
	ctx, keeper, bankKeeper, _, _, _ := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	receiverCoins := bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins, sdk.Coins{})

	claimType, err := types.StringToClaimType("burn")
	require.NoError(t, err)
	claimContent := types.NewOracleClaimContent(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

	claimBytes, err := json.Marshal(claimContent)
	require.NoError(t, err)
	claimString := string(claimBytes)
	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	require.NoError(t, err)

	receiverCoins = bankKeeper.GetCoins(ctx, cosmosReceivers[0])

	require.Equal(t, receiverCoins.String(), "10stake")

	// duplicate processSuccessClaim
	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	require.NoError(t, err)

	receiverCoins = bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, "20stake", receiverCoins.String())
}
func TestProcessBurn(t *testing.T) {
	ctx, keeper, bankKeeper, _, _, _ := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	coins := sdk.NewCoins(sdk.NewInt64Coin("stake", amount))

	receiverCoins := bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins, sdk.Coins{})
	err := keeper.ProcessBurn(ctx, cosmosReceivers[0], coins)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "insufficient account funds"))

	//process successful claim to get stake

	claimType, err := types.StringToClaimType("burn")
	require.NoError(t, err)
	claimContent := types.NewOracleClaimContent(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

	claimBytes, err := json.Marshal(claimContent)
	require.NoError(t, err)
	claimString := string(claimBytes)
	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	require.NoError(t, err)

	err = keeper.ProcessBurn(ctx, cosmosReceivers[0], coins)
	require.NoError(t, err)
	// lock stake

	receiverCoins = bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))

}

func TestProcessLock(t *testing.T) {

	ctx, keeper, bankKeeper, _, _, _ := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	receiverCoins := bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins, sdk.Coins{})

	coins := sdk.NewCoins(sdk.NewInt64Coin("stake", amount))
	err := keeper.ProcessLock(ctx, cosmosReceivers[0], coins)
	require.True(t, strings.Contains(err.Error(), "insufficient account funds"))

	//process successful claim to get stake

	claimType, err := types.StringToClaimType("burn")
	require.NoError(t, err)
	claimContent := types.NewOracleClaimContent(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

	claimBytes, err := json.Marshal(claimContent)
	require.NoError(t, err)
	claimString := string(claimBytes)
	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	require.NoError(t, err)

	receiverCoins = bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), "10stake")

	err = keeper.ProcessLock(ctx, cosmosReceivers[0], coins)
	require.NoError(t, err)

	receiverCoins = bankKeeper.GetCoins(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))

}
