package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
)

func (k keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	if state.AdminAccounts != nil {
		for _, adminAccount := range state.AdminAccounts.AdminAccounts {
			addr, err := sdk.AccAddressFromBech32(adminAccount.Address)
			if err != nil {
				panic(err)
			}
			k.SetAdminAccount(ctx, adminAccount.ModuleName, addr)
		}
	}
	if state.Registry != nil {
		k.SetRegistry(ctx, *state.Registry)
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
