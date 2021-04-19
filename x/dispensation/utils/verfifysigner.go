package utils

import (
	"bytes"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	"github.com/tendermint/tendermint/crypto"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

)

func VerifyInputList(inputList []banktypes.Input, pubKeys []crypto.PubKey) error {
	for _, i := range inputList {
		addressFound := false
		for _, signPubKeys := range pubKeys {
			if bytes.Equal(signPubKeys.Address().Bytes(), []byte(i.Address)) {
				addressFound = true
				continue
			}
		}
		if !addressFound {
			return errors.Wrap(types.ErrKeyInvalid, i.Address)
		}
	}
	return nil
}
