package app

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState(cdc codec.Codec) GenesisState {
	return ModuleBasics.DefaultGenesis(cdc)
}
