package types_test

import (
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"testing"
)

func TestMsgCreateClaim_ValidateBasic_LessAddressesInInputList(t *testing.T) {
	privKey1 := ed25519.GenPrivKey()
	privKey2 := ed25519.GenPrivKey()
	privKey3 := ed25519.GenPrivKey()
	mkey := multisig.PubKeyMultisigThreshold{
		K:       3,
		PubKeys: []crypto.PubKey{privKey1.PubKey(), privKey2.PubKey(), privKey3.PubKey()},
	}
	input1 := bank.Input{
		Address: sdk.AccAddress(privKey1.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	input2 := bank.Input{
		Address: sdk.AccAddress(privKey2.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	inputList := []bank.Input{input1, input2}
	msg := types.MsgDistribution{
		Signer:           mkey,
		DistributionName: "testName",
		DistributionType: types.Airdrop,
		Input:            inputList,
		Output:           test.CreatOutputList(2000, "1"),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic_MoreAddressesInInputList(t *testing.T) {
	privKey1 := ed25519.GenPrivKey()
	privKey2 := ed25519.GenPrivKey()
	privKey3 := ed25519.GenPrivKey()
	mkey := multisig.PubKeyMultisigThreshold{
		K:       2,
		PubKeys: []crypto.PubKey{privKey1.PubKey(), privKey2.PubKey()},
	}
	input1 := bank.Input{
		Address: sdk.AccAddress(privKey1.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	input2 := bank.Input{
		Address: sdk.AccAddress(privKey2.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	input3 := bank.Input{
		Address: sdk.AccAddress(privKey3.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	inputList := []bank.Input{input1, input2, input3}
	msg := types.MsgDistribution{
		Signer:           mkey,
		DistributionName: "testName",
		DistributionType: types.Airdrop,
		Input:            inputList,
		Output:           test.CreatOutputList(2000, "1"),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic_AddressMismatch(t *testing.T) {
	privKey1 := ed25519.GenPrivKey()
	privKey2 := ed25519.GenPrivKey()
	privKey3 := ed25519.GenPrivKey()
	privKey4 := ed25519.GenPrivKey()
	mkey := multisig.PubKeyMultisigThreshold{
		K:       2,
		PubKeys: []crypto.PubKey{privKey1.PubKey(), privKey2.PubKey()},
	}
	input1 := bank.Input{
		Address: sdk.AccAddress(privKey3.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	input2 := bank.Input{
		Address: sdk.AccAddress(privKey4.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	inputList := []bank.Input{input1, input2}
	msg := types.MsgDistribution{
		Signer:           mkey,
		DistributionName: "testName",
		DistributionType: types.Airdrop,
		Input:            inputList,
		Output:           test.CreatOutputList(2000, "1"),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic(t *testing.T) {
	privKey1 := ed25519.GenPrivKey()
	privKey2 := ed25519.GenPrivKey()
	mkey := multisig.PubKeyMultisigThreshold{
		K:       2,
		PubKeys: []crypto.PubKey{privKey1.PubKey(), privKey2.PubKey()},
	}
	input1 := bank.Input{
		Address: sdk.AccAddress(privKey1.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	input2 := bank.Input{
		Address: sdk.AccAddress(privKey2.PubKey().Address()),
		Coins:   sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))},
	}
	inputList := []bank.Input{input1, input2}
	msg := types.MsgDistribution{
		Signer:           mkey,
		DistributionName: "testName",
		DistributionType: types.Airdrop,
		Input:            inputList,
		Output:           test.CreatOutputList(2000, "1"),
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgDistribution_TestMultiSig(t *testing.T) {
	n := 2
	inputCoin := sdk.Coins{sdk.NewCoin("rowan", sdk.NewInt(1000))}
	privkeys := make([]crypto.PrivKey, n)
	pubkeys := make([]crypto.PubKey, n)
	inputList := make([]bank.Input, n)
	signatures := make([][]byte, n)
	for i := 0; i < n; i++ {
		privkeys[i] = ed25519.GenPrivKey()
		pubkeys[i] = privkeys[i].PubKey()
		inputList[i] = bank.Input{
			Address: sdk.AccAddress(pubkeys[i].Address()),
			Coins:   inputCoin,
		}
	}
	mkey := multisig.PubKeyMultisigThreshold{
		K:       uint(n),
		PubKeys: pubkeys,
	}
	msg := types.MsgDistribution{
		Signer:           mkey,
		DistributionName: "testName",
		DistributionType: types.Airdrop,
		Input:            inputList,
		Output:           test.CreatOutputList(2000, "1"),
	}
	multiSignature := multisig.NewMultisig(len(pubkeys))
	// mkey.Pubkeys is the same as pubkeys list
	for i, key := range privkeys {
		signatures[i], _ = key.Sign(msg.GetSignBytes())
		multiSignature.AddSignature(signatures[i], i)
		break
	}
	assert.False(t, mkey.VerifyBytes(msg.GetSignBytes(), multiSignature.Marshal()), "Not enough signatures")
	for i, key := range privkeys {
		signatures[i], _ = key.Sign(msg.GetSignBytes())
		multiSignature.AddSignature(signatures[i], i)
	}
	assert.True(t, mkey.VerifyBytes(msg.GetSignBytes(), multiSignature.Marshal()), "Enough signatures")
}
