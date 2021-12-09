package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// RegisterLegacyAminoCodec registers concrete types on the Amino codec
//lint:ignore SA1019 Legacy handler has to use legacy/deprecated features
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateEthBridgeClaim{}, "ethbridge/MsgCreateEthBridgeClaim", nil)
	cdc.RegisterConcrete(MsgBurn{}, "ethbridge/MsgBurn", nil)
	cdc.RegisterConcrete(MsgLock{}, "ethbridge/MsgLock", nil)
	cdc.RegisterConcrete(MsgUpdateWhiteListValidator{}, "ethbridge/MsgUpdateWhiteListValidator", nil)
	cdc.RegisterConcrete(MsgUpdateCrossChainFeeReceiverAccount{}, "ethbridge/MsgUpdateCrossChainFeeReceiverAccount", nil)
	cdc.RegisterConcrete(MsgRescueCrossChainFee{}, "ethbridge/MsgRescueCrossChainFee", nil)
	cdc.RegisterConcrete(MsgSetFeeInfo{}, "ethbridge/MsgSetFeeInfo", nil)
	cdc.RegisterConcrete(MsgSignProphecy{}, "ethbridge/MsgSignProphecy", nil)
	cdc.RegisterConcrete(MsgUpdateConsensusNeeded{}, "ethbridge/MsgUpdateConsensusNeeded", nil)
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
