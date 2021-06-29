package types_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMsgCreateDistribution_ValidateBasic(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           test.CreatOutputList(2000, "1"),
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_WrongAddress(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	outputList := test.CreatOutputList(1, "1")
	validAddress := sdk.AccAddress("addr2_______________")
	inValidAddress := validAddress[1:]
	authorizedRunner := sdk.AccAddress("addr3_______________")
	outputList = append(outputList, banktypes.NewOutput(inValidAddress,
		sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(1000000)))))
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputList,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_NonRowan(t *testing.T) {
	distributor := sdk.AccAddress([]byte("addr1_______________"))
	outputlist := test.CreatOutputList(2000, "1")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	outputlist = append(outputlist, banktypes.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.NewCoins(sdk.NewCoin("dash", sdk.NewInt(10)))))
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputlist,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_MultipleCoins(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	outputlist := test.CreatOutputList(2000, "1")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	outputlist = append(outputlist, banktypes.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(10)), sdk.NewCoin("dash", sdk.NewInt(100)))))
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputlist,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_ZeroCoins(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	var outputlist []banktypes.Output
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputlist,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic_WrongType(t *testing.T) {
	claimer := sdk.AccAddress("addr1_______________")
	msg := types.MsgCreateUserClaim{
		UserClaimAddress: claimer.String(),
		UserClaimType:    types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic(t *testing.T) {
	claimer := sdk.AccAddress("addr1_______________")
	msg := types.MsgCreateUserClaim{
		UserClaimAddress: claimer.String(),
		UserClaimType:    types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY,
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}
