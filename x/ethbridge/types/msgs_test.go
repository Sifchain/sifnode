package types_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/stretchr/testify/assert"
)

var (
	cosmosReceivers, _   = test.CreateTestAddrs(1)
	amount               = sdk.NewInt(1000000000000000000)
	symbol               = "eth"
	cosmosSymbol         = "ceth"
	tokenContractAddress = types.NewEthereumAddress("0xbbbbca6a901c926f240b89eacb641d8aec7aeafd")
	ethBridgeAddress     = types.NewEthereumAddress(strings.ToLower("0x30753E4A8aad7F8597332E813735Def5dD395028"))
	ethereumSender       = types.NewEthereumAddress("0x627306090abaB3A6e1400e9345bC60c78a8BEf57")
)

func TestNewMsgLock(t *testing.T) {
	msg := types.NewMsgLock(1, cosmosReceivers[0], ethereumSender, amount, symbol, amount)
	err := msg.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, msg.GetSigners()[0], cosmosReceivers[0])
	msg = types.NewMsgLock(1, cosmosReceivers[0], ethereumSender, amount, "", amount)
	err = msg.ValidateBasic()
	assert.Error(t, err)
}

func TestNewMsgBurn(t *testing.T) {
	msg := types.NewMsgBurn(1, cosmosReceivers[0], ethereumSender, amount, cosmosSymbol, amount)
	err := msg.ValidateBasic()
	assert.NoError(t, err)
	assert.Equal(t, msg.GetSigners()[0], cosmosReceivers[0])
	msg = types.NewMsgBurn(1, cosmosReceivers[0], ethereumSender, amount, "", amount)
	err = msg.ValidateBasic()
	assert.Error(t, err)
}

func TestNewMsgCreateEthBridgeClaim(t *testing.T) {
	valAddress, _ := sdk.ValAddressFromBech32("")
	ethClaim := types.CreateTestEthClaim(t, ethBridgeAddress, tokenContractAddress, valAddress,
		ethereumSender, amount, symbol, types.ClaimType_CLAIM_TYPE_LOCK)
	msg := types.NewMsgCreateEthBridgeClaim(ethClaim)
	err := msg.ValidateBasic()
	assert.Error(t, err)
}
