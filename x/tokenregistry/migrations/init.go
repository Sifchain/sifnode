package migrations

import (
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Init(ctx sdk.Context, keeper tokenregistrytypes.Keeper) {
	addr, err := sdk.AccAddressFromBech32("sif1tpypxpppcf5lea47vcvgy09675nllmcucxydvu")
	if err != nil {
		panic(err)
	}

	keeper.SetAdminAccount(ctx, addr)

	registry := tokenregistrytypes.DefaultRegistry()

	for _, t := range registry.Entries {
		keeper.SetToken(ctx, t)
	}
}
