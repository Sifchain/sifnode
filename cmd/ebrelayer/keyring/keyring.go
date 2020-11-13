package keyring

import (
	"fmt"

	"github.com/Sifchain/sifnode/app"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
// AccountAddressPrefix = "sif"
)

var (
	// AccountPubKeyPrefix = AccountAddressPrefix + "pub"
	hdpath = *hd.NewFundraiserParams(0, sdk.CoinType, 0)
)

// KeyRing
type KeyRing struct {
	mnemonic string
	moniker  string
	password string
	kb       keys.Keybase
}

// KeyRing
func NewKeyRing(mnemonic string, moniker string, password string) *KeyRing {
	return &KeyRing{
		mnemonic: mnemonic,
		moniker:  moniker,
		password: password,
	}
}

func (r *KeyRing) SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountPubKeyPrefix)
	config.Seal()
}

func (r *KeyRing) GenerateKeyStore() {
	r.kb = keys.NewInMemory(keys.WithSupportedAlgosLedger([]keys.SigningAlgo{keys.Secp256k1, keys.Ed25519}))
}

func (r *KeyRing) GetAccountFromMnemonic() keys.Info {
	account, err := r.kb.CreateAccount(r.moniker, r.mnemonic, "", r.password, hdpath.String(), keys.Secp256k1)
	if err != nil {
		panic(err)
	}

	return account
}

func (r *KeyRing) Address() {
	account := r.GetAccountFromMnemonic()
	fmt.Println(account.GetAddress().String())
}

func (r *KeyRing) Sign(msg []byte) ([]byte, crypto.PubKey, error) {
	pl, pk, err := r.kb.Sign(r.moniker, r.password, msg)
	if err != nil {
		// panic(err)
		fmt.Println(err)
		return nil, nil, err
	}
	return pl, pk, err

	// currentPubKey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, r.GetAccountFromMnemonic().GetPubKey())
	// fmt.Printf("Account PubKey: %s\n", currentPubKey)

	// signedPubKey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pk)
	// fmt.Printf("Signature PubKey: %s\n", signedPubKey)
	// fmt.Println(pl)
}
