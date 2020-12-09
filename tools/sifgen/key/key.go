package key

import (
	"github.com/Sifchain/sifnode/app"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tyler-smith/go-bip39"
)

var (
	hdpath = *hd.NewFundraiserParams(0, sdk.CoinType, 0)
)

type Key struct {
	Name             *string
	Mnemonic         string
	Password         *string
	Address          string
	ValidatorAddress string
	ConsensusAddress string
	Keybase          keys.Keybase
}

func NewKey(name, password *string) *Key {
	return &Key{
		Name:     name,
		Password: password,
		Keybase:  keys.NewInMemory(keys.WithSupportedAlgosLedger([]keys.SigningAlgo{keys.Secp256k1, keys.Ed25519})),
	}
}

func (k *Key) GenerateMnemonic() {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	k.Mnemonic = mnemonic
}

func (k *Key) RecoverFromMnemonic(mnemonic string) error {
	k.setConfig()
	k.Mnemonic = mnemonic

	account, err := k.Keybase.CreateAccount(*k.Name, k.Mnemonic, "", *k.Password, hdpath.String(), keys.Secp256k1)
	if err != nil {
		return err
	}

	consensusAddress := sdk.ConsAddress(account.GetPubKey().Address()).String()
	validatorAddress := sdk.ValAddress(account.GetPubKey().Address()).String()

	k.Address = account.GetAddress().String()
	k.ValidatorAddress = validatorAddress
	k.ConsensusAddress = consensusAddress

	return nil
}

func (k *Key) setConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(app.ValidatorAddressPrefix, app.ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(app.ConsNodeAddressPrefix, app.ConsNodePubKeyPrefix)
	config.Seal()
}
