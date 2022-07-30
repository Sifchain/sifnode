package tokenregistry

import (
	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestingONLY_SetRegistry(ctx sdk.Context, keeper types.Keeper) {
	keeper.SetRegistry(ctx, *types.InitialRegistry())
}
