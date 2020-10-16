package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreatePool{}, "clp/CreatePool", nil)
	cdc.RegisterConcrete(MsgAddLiquidity{}, "clp/AddLiquidity", nil)
	cdc.RegisterConcrete(MsgRemoveLiquidity{}, "clp/RemoveLiquidity", nil)
	cdc.RegisterConcrete(MsgSwap{}, "clp/Swap", nil)
	cdc.RegisterConcrete(MsgDecommissionPool{}, "clp/DecommissionPool", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
