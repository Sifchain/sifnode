package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/whitelist/types"
)

func (k keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	if state.AdminAccount != "" {
		addr, err := sdk.AccAddressFromBech32(state.AdminAccount)
		if err != nil {
			panic(err)
		}
		k.SetAdminAccount(ctx, addr)
	}

	if state.Registry != nil {
		k.SetDenomWhitelist(ctx, *state.Registry)
	}

	return []abci.ValidatorUpdate{}
}

func (k keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	wl := k.GetDenomWhitelist(ctx)

	return &types.GenesisState{
		AdminAccount: k.GetAdminAccount(ctx).String(),
		Registry:     &wl,
	}
}
