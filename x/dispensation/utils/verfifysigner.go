package utils

import (
	"bytes"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/tendermint/tendermint/crypto"
)

func VerifyInputList(inputList []bank.Input, pubKeys []crypto.PubKey) error {
	for _, i := range inputList {
		addressFound := false
		for _, signPubKeys := range pubKeys {
			if bytes.Equal(signPubKeys.Address().Bytes(), i.Address.Bytes()) {
				addressFound = true
				continue
			}
		}
		if !addressFound {
			return errors.Wrap(types.ErrKeyInvalid, i.Address.String())
		}
	}
	return nil
}
