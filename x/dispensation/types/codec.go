package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers concrete types on codec
<<<<<<< HEAD
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDistribution{}, "dispensation/MsgCreateDistribution", nil)
	cdc.RegisterConcrete(&Distribution{}, "dispensation/Distribution", nil)
	cdc.RegisterConcrete(&MsgCreateUserClaim{}, "dispensation/claim", nil)
=======
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgDistribution{}, "dispensation/create", nil)
	cdc.RegisterConcrete(MsgCreateClaim{}, "dispensation/claim", nil)
	cdc.RegisterConcrete(MsgRunDistribution{}, "dispensation/run", nil)
>>>>>>> develop
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateDistribution{},
		&MsgCreateUserClaim{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
