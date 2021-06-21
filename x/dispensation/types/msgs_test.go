package types_test

/*
import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestMsgCreateDistribution_ValidateBasic(t *testing.T) {
	distributor := ed25519.GenPrivKey()
	msg := types.MsgDistribution{
		Distributor:      sdk.AccAddress(distributor.PubKey().Address()),
		DistributionType: types.Airdrop,
		Output:           test.CreatOutputList(2000, "1"),
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_WrongAddress(t *testing.T) {
	distributor := ed25519.GenPrivKey()
	outputList := test.CreatOutputList(1, "1")
	validAddress := ed25519.GenPrivKey().PubKey().Address()
	inValidAddress := validAddress[1:]
	outputList = append(outputList, bank.NewOutput(sdk.AccAddress(inValidAddress),
		sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(1000000)))))
	msg := types.MsgDistribution{
		Distributor:      sdk.AccAddress(distributor.PubKey().Address()),
		DistributionType: types.Airdrop,
		Output:           outputList,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_NonRowan(t *testing.T) {
	distributor := ed25519.GenPrivKey()
	outputlist := test.CreatOutputList(2000, "1")
	outputlist = append(outputlist, bank.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.NewCoins(sdk.NewCoin("dash", sdk.NewInt(10)))))
	msg := types.MsgDistribution{
		Distributor:      sdk.AccAddress(distributor.PubKey().Address()),
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
		DistributionType: types.Airdrop,
		Output:           outputlist,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_ZeroCoins(t *testing.T) {
	distributor := ed25519.GenPrivKey()
	var outputlist []bank.Output
	msg := types.MsgDistribution{
		Distributor:      sdk.AccAddress(distributor.PubKey().Address()),
		DistributionType: types.Airdrop,
		Output:           outputlist,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic_WrongType(t *testing.T) {
	claimer := ed25519.GenPrivKey()
	msg := types.MsgCreateClaim{
		UserClaimAddress: sdk.AccAddress(claimer.PubKey().Address()),
		UserClaimType:    types.Airdrop,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic(t *testing.T) {
	claimer := ed25519.GenPrivKey()
	msg := types.MsgCreateClaim{
		UserClaimAddress: sdk.AccAddress(claimer.PubKey().Address()),
		UserClaimType:    types.ValidatorSubsidy,
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

*/
