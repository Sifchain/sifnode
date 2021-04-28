package key

import (
	"testing"

	"github.com/tyler-smith/go-bip39"
)

var (
	name             = "cool-frost"
	address          = "sif1ftf2l7y3d330zxun9dyh4he7rak7uytalu9dlx"
	validatorAddress = "sifvaloper1ftf2l7y3d330zxun9dyh4he7rak7uytak7d8lk"
	consensusAddress = "sifvalcons1ftf2l7y3d330zxun9dyh4he7rak7uytazd7mnh"
	random           = "qY3XtZc4a16jTnQWsJUwDvxfC2giHmSB"
	mnemonic         = "flock toss tip service element interest leisure bright subway critic copy lazy zero limb unveil reveal ecology slab detail wrong smooth fun pond choice"
)

func TestGenerateMnemonic(t *testing.T) {
	k := NewKey(name, random)
	k.GenerateMnemonic()

	if !bip39.IsMnemonicValid(k.Mnemonic) {
		t.Error("mnemonic is invalid")
	}
}

func TestRecoverFromMnemonic(t *testing.T) {
	k := NewKey(name, random)
	if err := k.RecoverFromMnemonic(mnemonic); err != nil {
		t.Error(err)
	}

	if k.Address != address {
		t.Errorf("expected %s, got %s", address, k.Address)
	}

	if k.ValidatorAddress != validatorAddress {
		t.Errorf("expected %s, got %s", validatorAddress, k.ValidatorAddress)
	}

	if k.ConsensusAddress != consensusAddress {
		t.Errorf("expected %s, got %s", consensusAddress, k.ConsensusAddress)
	}
}
