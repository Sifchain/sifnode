package types_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func TestTypes_RegisterLegacyAminoCodec(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)
}

func TestTypes_RegisterInterfaces(t *testing.T) {
	registry := cdctypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
}
