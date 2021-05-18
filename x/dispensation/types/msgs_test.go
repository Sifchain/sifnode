package types_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestMsgCreateClaim_ValidateBasic(t *testing.T) {
	distributor := ed25519.GenPrivKey()
	msg := types.MsgDistribution{
		Distributor:      sdk.AccAddress(distributor.PubKey().Address()),
		DistributionName: "testName",
		DistributionType: types.Airdrop,
		Output:           test.CreatOutputList(2000, "1"),
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_NonRowan(t *testing.T) {
	distributor := ed25519.GenPrivKey()
	outputlist := test.CreatOutputList(2000, "1")
	outputlist = append(outputlist, bank.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.NewCoins(sdk.NewCoin("dash", sdk.NewInt(10)))))
	msg := types.MsgDistribution{
		Distributor:      sdk.AccAddress(distributor.PubKey().Address()),
		DistributionName: "testName",
		DistributionType: types.Airdrop,
		Output:           outputlist,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_MultipleCoins(t *testing.T) {
	distributor := ed25519.GenPrivKey()
	outputlist := test.CreatOutputList(2000, "1")
	outputlist = append(outputlist, bank.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(10)), sdk.NewCoin("dash", sdk.NewInt(100)))))
	msg := types.MsgDistribution{
		Distributor:      sdk.AccAddress(distributor.PubKey().Address()),
		DistributionName: "testName",
		DistributionType: types.Airdrop,
		Output:           outputlist,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_NoName(t *testing.T) {
	distributor := ed25519.GenPrivKey()
	msg := types.MsgDistribution{
		Distributor:      sdk.AccAddress(distributor.PubKey().Address()),
		DistributionName: "",
		DistributionType: types.Airdrop,
		Output:           test.CreatOutputList(2000, "1"),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}
