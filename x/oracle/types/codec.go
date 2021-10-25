package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc defines the module codec
var ModuleCdc *codec.LegacyAmino //nolint

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	ModuleCdc.Seal()
}
