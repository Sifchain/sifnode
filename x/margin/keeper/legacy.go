package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewLegacyQuerier(k types.Keeper) sdk.Querier {
	return nil
}

func NewLegacyHandler(k types.Keeper) sdk.Handler {
	return nil
}
