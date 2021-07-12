package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateEthBridgeClaim{}, "ethbridge/MsgCreateEthBridgeClaim", nil)
	cdc.RegisterConcrete(MsgBurn{}, "ethbridge/MsgBurn", nil)
	cdc.RegisterConcrete(MsgLock{}, "ethbridge/MsgLock", nil)
	cdc.RegisterConcrete(MsgUpdateWhiteListValidator{}, "ethbridge/MsgUpdateWhiteListValidator", nil)
	cdc.RegisterConcrete(MsgUpdateNativeTokenReceiverAccount{}, "ethbridge/MsgUpdateNativeTokenReceiverAccount", nil)
	cdc.RegisterConcrete(MsgRescueNativeToken{}, "ethbridge/MsgRescueNativeToken", nil)

}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}
