package migrations

import (
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Init(ctx sdk.Context, keeper tokenregistrytypes.Keeper) {
	registry := tokenregistrytypes.DefaultRegistry()

	for _, t := range registry.Entries {
		keeper.SetToken(ctx, t)
	}
}
