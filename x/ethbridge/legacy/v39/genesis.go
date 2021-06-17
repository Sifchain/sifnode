package _39

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "ethbridge"

// GenesisState - all ethbridge state that must be provided at genesis
type GenesisState struct {
	PeggyTokens         []string       `json:"peggy_tokens"`
	CethReceiverAccount sdk.AccAddress `json:"ceth_receiver_account"`
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cryptocodec.RegisterCrypto(cdc)
}