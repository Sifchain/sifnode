package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func (k keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	a := *types.InitialAdminAccounts()
	for _, admin := range a.AdminAccounts {
		k.SetAdminAccount(ctx, admin)
	}
	return []abci.ValidatorUpdate{}
}

func (k keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	wl := k.GetRegistry(ctx)
	return &types.GenesisState{
		AdminAccounts: k.GetAdminAccounts(ctx),
		Registry:      &wl,
	}
}
