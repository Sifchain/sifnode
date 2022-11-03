package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
//
//lint:ignore SA1019 Legacy handler has to use legacy/deprecated features
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegister{}, "tokenregistry/MsgRegister", nil)
	cdc.RegisterConcrete(&MsgRegisterAll{}, "tokenregistry/MsgRegisterAll", nil)
	cdc.RegisterConcrete(&MsgSetRegistry{}, "tokenregistry/MsgSetRegistry", nil)
	cdc.RegisterConcrete(&MsgDeregister{}, "tokenregistry/MsgDeregister", nil)
	cdc.RegisterConcrete(&MsgDeregisterAll{}, "tokenregistry/MsgDeregisterAll", nil)
	cdc.RegisterConcrete(&TokenMetadataAddRequest{}, "tokenregistry/TokenMetadataAddRequest", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegister{},
		&MsgRegisterAll{},
		&MsgSetRegistry{},
		&MsgDeregister{},
		&MsgDeregisterAll{},
		&TokenMetadataAddRequest{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
