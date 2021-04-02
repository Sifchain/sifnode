package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateEthBridgeClaim{}, "ethbridge/MsgCreateEthBridgeClaim", nil)
	cdc.RegisterConcrete(MsgBurn{}, "ethbridge/MsgBurn", nil)
	cdc.RegisterConcrete(MsgLock{}, "ethbridge/MsgLock", nil)
	cdc.RegisterConcrete(MsgUpdateWhiteListValidator{}, "ethbridge/MsgUpdateWhiteListValidator", nil)
	cdc.RegisterConcrete(MsgUpdateCethReceiverAccount{}, "ethbridge/MsgUpdateCethReceiverAccount", nil)
	cdc.RegisterConcrete(MsgRescueCeth{}, "ethbridge/MsgRescueCeth", nil)
	cdc.RegisterConcrete(MsgUpdateGasPrice{}, "ethbridge/MsgUpdateGasPrice", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
