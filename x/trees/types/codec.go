package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateTree{}, "trees/CreateTree", nil)
	cdc.RegisterConcrete(MsgBuyTree{}, "trees/BuyTree", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
