package key

//
//import (
//	"github.com/Sifchain/sifnode/app"
//	"github.com/cosmos/cosmos-sdk/client/keys"
//	"github.com/cosmos/cosmos-sdk/crypto/hd"
//	"github.com/cosmos/cosmos-sdk/crypto/keyring"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/tyler-smith/go-bip39"
//)
//
//var (
//	hdpath = *hd.NewFundraiserParams(0, sdk.CoinType, 0)
//)
///*    * [\#5904](https://github.com/cosmos/cosmos-sdk/pull/5904) `Keybase` -> `Keyring` interfaces migration. `LegacyKeybase` interface is added in order
//  to guarantee limited backward compatibility with the old Keybase interface for the sole purpose of migrating keys across the new keyring backends. `NewLegacy`
//  constructor is provided [\#5889](https://github.com/cosmos/cosmos-sdk/pull/5889) to allow for smooth migration of keys from the legacy LevelDB based implementation
//  to new keyring backends. Plus, the package and the new keyring no longer depends on the sdk.Config singleton. Please consult the [package documentation](https://github.com/cosmos/cosmos-sdk/tree/master/crypto/keyring/doc.go) for more
//  information on how to implement the new `Keyring` interface.
//    * [\#5858](https://github.com/cosmos/cosmos-sdk/pull/5858) Make Keyring store keys by name and address's hexbytes representation.*/
//type Key struct {
//	Name             *string
//	Mnemonic         string
//	Password         *string
//	Address          string
//	ValidatorAddress string
//	ConsensusAddress string
//	Keybase          keyring.Keyring
//}
//
//func NewKey(name, password *string) *Key {
//	return &Key{
//		Name:     name,
//		Password: password,
//		Keybase:  keys.NewInMemory(keys.WithSupportedAlgosLedger([]keys.SigningAlgo{keys.Secp256k1, keys.Ed25519})),
//	}
//}
//
//func (k *Key) GenerateMnemonic() {
//	entropy, _ := bip39.NewEntropy(256)
//	mnemonic, _ := bip39.NewMnemonic(entropy)
//	k.Mnemonic = mnemonic
//}
//
//func (k *Key) RecoverFromMnemonic(mnemonic string) error {
//	k.setConfig()
//	k.Mnemonic = mnemonic
//
//	account, err := k.Keybase.CreateAccount(*k.Name, k.Mnemonic, "", *k.Password, hdpath.String(), keys.Secp256k1)
//	if err != nil {
//		return err
//	}
//
//	consensusAddress := sdk.ConsAddress(account.GetPubKey().Address()).String()
//	validatorAddress := sdk.ValAddress(account.GetPubKey().Address()).String()
//
//	k.Address = account.GetAddress().String()
//	k.ValidatorAddress = validatorAddress
//	k.ConsensusAddress = consensusAddress
//
//	return nil
//}
//
//func (k *Key) setConfig() {
//	config := sdk.GetConfig()
//	config.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountPubKeyPrefix)
//	config.SetBech32PrefixForValidator(app.ValidatorAddressPrefix, app.ValidatorPubKeyPrefix)
//	config.SetBech32PrefixForConsensusNode(app.ConsNodeAddressPrefix, app.ConsNodePubKeyPrefix)
//	config.Seal()
//}
