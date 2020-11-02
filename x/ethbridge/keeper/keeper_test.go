package keeper

import (
	"encoding/json"
	"fmt"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/oracle"
	_ "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

var (
	cosmosReceivers, _   = CreateTestAddrs(1)
	amount, _            = strconv.ParseInt("10", 10, 64)
	symbol               = "stake"
	tokenContractAddress = types.NewEthereumAddress("0xbbbbca6a901c926f240b89eacb641d8aec7aeafd") // loopring
)

func TestProcessClaim(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	validator1Pow3 := validatorAddresses[0]
	nonce, err := strconv.Atoi("1")
	require.NoError(t, err)
	//!!!!!!!!1
	ethBridgeAddress := strings.ToLower("0x30753E4A8aad7F8597332E813735Def5dD395028")
	//!!!!!!!!!1
	claimType, err := types.StringToClaimType("lock")
	require.NoError(t, err)
	ethBridgeClaim := types.NewEthBridgeClaim(
		5777,
		types.NewEthereumAddress(ethBridgeAddress), // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		types.NewEthereumAddress("0x627306090abaB3A6e1400e9345bC60c78a8BEf57"), // accounts[0]
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

}

func TestProcessSuccessfulClaim(t *testing.T) {
	ctx, keeper, _, _, _, _ := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	claimType, err := types.StringToClaimType("lock")
	require.NoError(t, err)
	fmt.Println(amount)
	claimContent := types.NewOracleClaimContent(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

	claimBytes, err := json.Marshal(claimContent)
	claimString := string(claimBytes)
	fmt.Println(claimString)
	require.Panics(t, func() { keeper.ProcessSuccessfulClaim(ctx, claimString) }, "the code did not panic")

}
