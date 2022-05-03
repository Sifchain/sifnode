package utils

import (
	"io/ioutil"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func ParseDenoms(codec codec.JSONCodec, dir string) (tokenregistrytypes.Registry, error) {
	var denoms tokenregistrytypes.Registry
	file, err := filepath.Abs(dir)
	if err != nil {
		return denoms, err
	}
	o, err := ioutil.ReadFile(file)
	if err != nil {
		return denoms, err
	}
	err = codec.UnmarshalJSON(o, &denoms)
	if err != nil {
		return denoms, err
	}
	return denoms, err
}
