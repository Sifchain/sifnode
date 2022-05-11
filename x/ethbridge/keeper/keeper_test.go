package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

var (
	cosmosReceivers, _   = test.CreateTestAddrs(1)
	amount               = sdk.NewInt(10)
	doubleAmount         = sdk.NewInt(20)
	decimals             = uint8(18)
	symbol               = "stake"
	name                 = "STAKE"
	tokenContractAddress = types.NewEthereumAddress("0xbbbbca6a901c926f240b89eacb641d8aec7aeafd")
	ethBridgeAddress     = types.NewEthereumAddress(strings.ToLower("0x30753E4A8aad7F8597332E813735Def5dD395028"))
	ethereumSender       = types.NewEthereumAddress("0x627306090abaB3A6e1400e9345bC60c78a8BEf57")
	networkDescriptor    = oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM
)

var testMetadataStake = tokenregistrytypes.TokenMetadata{
	Decimals:          int64(18),
	Name:              "stake",
	NetworkDescriptor: oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM,
	Symbol:            "stake",
	TokenAddress:      "0x0123456789ABCDEF",
}

var testMetadataCeth = tokenregistrytypes.TokenMetadata{
	Decimals:          int64(18),
	Name:              "ceth",
	NetworkDescriptor: oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM,
	Symbol:            "ceth",
	TokenAddress:      "Ox00000000000000000000",
}

var testMetadataRowan = tokenregistrytypes.TokenMetadata{
	Decimals:          int64(18),
	Name:              "rowan",
	NetworkDescriptor: oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_UNSPECIFIED,
	Symbol:            "rowan",
	TokenAddress:      "0x0123456789ABCDEF",
}

var testRopstenToken = tokenregistrytypes.TokenMetadata{
	Decimals:          int64(18),
	Name:              "ropsten",
	NetworkDescriptor: oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM_TESTNET_ROPSTEN,
	Symbol:            "ropsten",
	TokenAddress:      "0x0123456789ABCDEF",
}

func TestProcessClaimLock(t *testing.T) {
	ctx, keeper, _, _, _, _, _, validatorAddresses := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]

	nonce := uint64(1)

	// TODO(timlind): This default does not seem to be in any version history.
	//invalid claim defaults to lock
	//claimType, err := types.StringToClaimType("lkfjdsk")
	//require.Equal(t, claimType.String(), "lock")
	//require.Error(t, err)

	claimType := types.ClaimType_CLAIM_TYPE_LOCK
	require.Equal(t, claimType, types.ClaimType_CLAIM_TYPE_LOCK)
	ethBridgeClaim := types.NewEthBridgeClaim(
		1,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,
		cosmosReceivers[0],
		validator1Pow3,
		amount,
		claimType,
		name,
		decimals,
	)

	status, err := keeper.ProcessClaim(ctx, ethBridgeClaim)

	require.NoError(t, err)
	require.Equal(t, status, oracletypes.StatusText_STATUS_TEXT_PENDING)
	// duplicate execution
	_, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))

	// other validator execute
	ethBridgeClaim = types.NewEthBridgeClaim(
		1,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,       // accounts[0]
		cosmosReceivers[0],
		validator2Pow3,
		amount,
		claimType,
		name,
		decimals,
	)
	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.NoError(t, err)
	require.Equal(t, status, oracletypes.StatusText_STATUS_TEXT_SUCCESS)

}

func TestProcessClaimBurn(t *testing.T) {
	ctx, keeper, _, _, _, _, _, validatorAddresses := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]

	nonce := uint64(1)

	claimType := types.ClaimType_CLAIM_TYPE_BURN

	ethBridgeClaim := types.NewEthBridgeClaim(
		1,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,
		cosmosReceivers[0],
		validator1Pow3,
		amount,
		claimType,
		name,
		decimals,
	)

	status, err := keeper.ProcessClaim(ctx, ethBridgeClaim)

	require.NoError(t, err)
	require.Equal(t, status, oracletypes.StatusText_STATUS_TEXT_PENDING)

	_, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))

	// other validator execute

	ethBridgeClaim = types.NewEthBridgeClaim(
		1,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,       // accounts[0]
		cosmosReceivers[0],
		validator2Pow3,
		amount,
		claimType,
		name,
		decimals,
	)
	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.NoError(t, err)
	require.Equal(t, status.String(), oracletypes.StatusText_STATUS_TEXT_SUCCESS.String())

}

func TestProcessSuccessfulClaimLock(t *testing.T) {
	// ctx, keeper, bankKeeper, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	// receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	// require.Equal(t, receiverCoins, sdk.NewCoins())

	// claimType := types.ClaimType_CLAIM_TYPE_LOCK
	// claimContent := types.NewEthBridgeClaim(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

	// claimBytes, err := json.Marshal(claimContent)
	// require.NoError(t, err)
	// claimString := string(claimBytes)
	// err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	// require.NoError(t, err)

	// receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])

	// require.Equal(t, receiverCoins.String(), "10cstake")

	// // duplicate processSuccessClaim
	// err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	// require.NoError(t, err)

	// receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	// require.Equal(t, "20cstake", receiverCoins.String())
}

func TestProcessSuccessfulClaimBurn(t *testing.T) {
	// ctx, keeper, bankKeeper, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	// receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	// require.Equal(t, receiverCoins, sdk.NewCoins())

	// claimType := types.ClaimType_CLAIM_TYPE_BURN
	// claimContent := types.NewMsgBurn(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

	// claimBytes, err := json.Marshal(claimContent)
	// require.NoError(t, err)
	// claimString := string(claimBytes)
	// err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	// require.NoError(t, err)

	// receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])

	// require.Equal(t, receiverCoins.String(), "10stake")

	// // duplicate processSuccessClaim
	// err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	// require.NoError(t, err)

}

// func TestProcessSuccessfulClaimLock(t *testing.T) {
// 	ctx, keeper, bankKeeper, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

// 	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
// 	require.Equal(t, receiverCoins, sdk.Coins{})

// 	claimType := types.ClaimType_CLAIM_TYPE_LOCK
// 	claimContent := types.NewOracleClaimContent(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

// 	claimBytes, err := json.Marshal(claimContent)
// 	require.NoError(t, err)
// 	claimString := string(claimBytes)
// 	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
// 	require.NoError(t, err)

// 	receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])

// 	require.Equal(t, receiverCoins.String(), "10cstake")

// 	// duplicate processSuccessClaim
// 	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
// 	require.NoError(t, err)

// 	receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
// 	require.Equal(t, "20cstake", receiverCoins.String())
// }

// func TestProcessSuccessfulClaimBurn(t *testing.T) {
// 	ctx, keeper, bankKeeper, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

// 	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
// 	require.Equal(t, receiverCoins, sdk.Coins{})

// 	claimType := types.ClaimType_CLAIM_TYPE_BURN
// 	claimContent := types.NewOracleClaimContent(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

// 	claimBytes, err := json.Marshal(claimContent)
// 	require.NoError(t, err)
// 	claimString := string(claimBytes)
// 	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
// 	require.NoError(t, err)

// 	receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])

// 	require.Equal(t, receiverCoins.String(), "10stake")

// 	// duplicate processSuccessClaim
// 	err = keeper.ProcessSuccessfulClaim(ctx, claimString)
// 	require.NoError(t, err)

// 	receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
// 	require.Equal(t, "20stake", receiverCoins.String())
// }

func TestProcessBurn(t *testing.T) {
	ctx, keeper, bankKeeper, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crossChainFeeConfig, _ := oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)
	crossChainFee := crossChainFeeConfig.FeeCurrency
	denomHash := keeper.AddTokenMetadata(ctx, testMetadataStake)

	msg := types.NewMsgBurn(1, cosmosReceivers[0], ethereumSender, amount, denomHash, amount)
	coins := sdk.NewCoins(sdk.NewCoin(denomHash, amount), sdk.NewCoin(crossChainFee, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	account := keeper.GetAccountKeeper().GetAccount(ctx, cosmosReceivers[0])

	_, err := keeper.ProcessBurn(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataStake, false, false)
	require.NoError(t, err)

	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))
}

func TestProcessBurnCrossChainFee(t *testing.T) {
	ctx, keeper, bankKeeper, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crossChainFeeConfig, _ := oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)
	crossChainFee := crossChainFeeConfig.FeeCurrency

	denomHash := keeper.AddTokenMetadata(ctx, testMetadataCeth)

	msg := types.NewMsgBurn(networkDescriptor, cosmosReceivers[0], ethereumSender, amount, denomHash, amount)
	coins := sdk.NewCoins(sdk.NewCoin(crossChainFee, doubleAmount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	account := keeper.GetAccountKeeper().GetAccount(ctx, cosmosReceivers[0])
	_, err := keeper.ProcessBurn(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataCeth, false, false)
	require.NoError(t, err)

	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))
}

func TestProcessLock(t *testing.T) {
	ctx, keeper, bankKeeper, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	keeper.AddTokenMetadata(ctx, testMetadataStake)

	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crossChainFeeConfig, _ := oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)
	crossChainFee := crossChainFeeConfig.FeeCurrency
	coins := sdk.NewCoins(sdk.NewCoin(crossChainFee, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)
	denomHash := keeper.AddTokenMetadata(ctx, testMetadataRowan)

	msg := types.NewMsgLock(1, cosmosReceivers[0], ethereumSender, amount, denomHash, amount)
	account := keeper.GetAccountKeeper().GetAccount(ctx, cosmosReceivers[0])

	_, err := keeper.ProcessLock(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataRowan, true)
	require.ErrorIs(t, err, sdkerrors.ErrInsufficientFunds)

	coins = sdk.NewCoins(sdk.NewCoin(denomHash, amount), sdk.NewCoin(crossChainFee, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	_, err = keeper.ProcessLock(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataRowan, true)
	require.NoError(t, err)

	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), "")
}

func TestProcessBurnWithReceiver(t *testing.T) {
	ctx, keeper, bankKeeper, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	keeper.AddTokenMetadata(ctx, testMetadataStake)

	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crossChainFeeConfig, _ := oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)
	crossChainFee := crossChainFeeConfig.FeeCurrency
	denomHash := keeper.AddTokenMetadata(ctx, testMetadataStake)

	msg := types.NewMsgBurn(1, cosmosReceivers[0], ethereumSender, amount, denomHash, amount)
	coins := sdk.NewCoins(sdk.NewCoin(denomHash, amount), sdk.NewCoin(crossChainFee, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	account := keeper.GetAccountKeeper().GetAccount(ctx, cosmosReceivers[0])
	_, err = keeper.ProcessBurn(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataStake, false, false)
	require.NoError(t, err)

	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))
}

func TestProcessBurnCrossChainFeeWithReceiver(t *testing.T) {
	ctx, keeper, bankKeeper, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	denomHash := keeper.AddTokenMetadata(ctx, testMetadataCeth)

	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crossChainFeeConfig, _ := oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)
	crossChainFee := crossChainFeeConfig.FeeCurrency

	msg := types.NewMsgBurn(1, cosmosReceivers[0], ethereumSender, amount, denomHash, amount)
	coins := sdk.NewCoins(sdk.NewCoin(crossChainFee, doubleAmount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	account := keeper.GetAccountKeeper().GetAccount(ctx, cosmosReceivers[0])
	_, err = keeper.ProcessBurn(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataCeth, false, false)
	require.NoError(t, err)

	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))
}

func TestProcessLockWithReceiver(t *testing.T) {
	ctx, keeper, bankKeeper, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	denomHash := keeper.AddTokenMetadata(ctx, testMetadataStake)

	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crossChainFeeConfig, _ := oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)
	crossChainFee := crossChainFeeConfig.FeeCurrency
	coins := sdk.NewCoins(sdk.NewCoin(crossChainFee, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins, coins)

	msg := types.NewMsgLock(1, cosmosReceivers[0], ethereumSender, amount, denomHash, amount)

	account := keeper.GetAccountKeeper().GetAccount(ctx, cosmosReceivers[0])
	// require.Nil(t, account)
	_, err = keeper.ProcessLock(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataRowan, true)
	require.ErrorIs(t, err, sdkerrors.ErrInsufficientFunds)

	coins = sdk.NewCoins(sdk.NewCoin(denomHash, amount), sdk.NewCoin(crossChainFee, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	_, err = keeper.ProcessLock(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataRowan, true)
	require.NoError(t, err)

	receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))

}

func TestProcessBurnFirstDoublePeg(t *testing.T) {
	ctx, keeper, bankKeeper, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)
	oracleKeeper.SetAdminAccount(ctx, cosmosSender)
	// add ropsten token into metadata
	denomHash := keeper.AddTokenMetadata(ctx, testRopstenToken)

	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crossChainFeeConfig, _ := oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)
	crossChainFee := crossChainFeeConfig.FeeCurrency
	coins := sdk.NewCoins(sdk.NewCoin(denomHash, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	receiverCoins := bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins, coins)

	// double peg from ropstent to ethereum
	msg := types.NewMsgBurn(networkDescriptor, cosmosReceivers[0], ethereumSender, amount, denomHash, amount)

	account := keeper.GetAccountKeeper().GetAccount(ctx, cosmosReceivers[0])
	_, err = keeper.ProcessBurn(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataCeth, true, false)
	require.ErrorIs(t, err, sdkerrors.ErrInsufficientFunds)

	coins = sdk.NewCoins(sdk.NewCoin(crossChainFee, amount))
	_ = bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	_ = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, cosmosReceivers[0], coins)

	_, err = keeper.ProcessBurn(ctx, cosmosReceivers[0], account.GetSequence(), &msg, testMetadataCeth, true, false)
	require.NoError(t, err)

	receiverCoins = bankKeeper.GetAllBalances(ctx, cosmosReceivers[0])
	require.Equal(t, receiverCoins.String(), string(""))
}

func TestProcessUpdateCrossChainFeeReceiverAccount(t *testing.T) {
	ctx, keeper, _, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)

	err = keeper.ProcessUpdateCrossChainFeeReceiverAccount(ctx, cosmosSender, cosmosSender)
	require.Equal(t, err.Error(), "only admin account can update CrossChainFee receiver account")

	oracleKeeper.SetAdminAccount(ctx, cosmosSender)

	err = keeper.ProcessUpdateCrossChainFeeReceiverAccount(ctx, cosmosSender, cosmosSender)
	require.NoError(t, err)
}

func TestRescueCrossChainFees(t *testing.T) {
	ctx, keeper, bankKeeper, _, oracleKeeper, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	cosmosSender, err := sdk.AccAddressFromBech32(types.TestAddress)
	require.NoError(t, err)

	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crossChainFeeConfig, _ := oracleKeeper.GetCrossChainFeeConfig(ctx, networkIdentity)
	crossChainFee := crossChainFeeConfig.FeeCurrency

	crosschainFee := sdk.NewInt(100)
	err = bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(crossChainFee, crosschainFee)))
	require.NoError(t, err)

	msg := types.NewMsgRescueCrossChainFee(cosmosSender, cosmosSender, crossChainFee, crosschainFee)

	err = keeper.RescueCrossChainFees(ctx, &msg)
	require.Equal(t, err.Error(), "only admin account can call rescue CrossChainFee")

	oracleKeeper.SetAdminAccount(ctx, cosmosSender)

	err = keeper.RescueCrossChainFees(ctx, &msg)
	require.NoError(t, err)
}
