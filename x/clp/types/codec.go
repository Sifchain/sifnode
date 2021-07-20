package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePool{}, "clp/CreatePool", nil)
	cdc.RegisterConcrete(&MsgAddLiquidity{}, "clp/AddLiquidity", nil)
	cdc.RegisterConcrete(&MsgRemoveLiquidity{}, "clp/RemoveLiquidity", nil)
	cdc.RegisterConcrete(&MsgSwap{}, "clp/Swap", nil)
	cdc.RegisterConcrete(&MsgDecommissionPool{}, "clp/DecommissionPool", nil)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRemoveLiquidity{},
		&MsgCreatePool{},
		&MsgAddLiquidity{},
		&MsgSwap{},
		&MsgDecommissionPool{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
