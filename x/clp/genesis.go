package clp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	k.SetParams(ctx, data.Params)

	if data.AddressWhitelist == nil || len(data.AddressWhitelist) == 0 {
		panic(fmt.Sprintf("AddressWhiteList must be set."))
	}

	wl := make([]sdk.AccAddress, len(data.AddressWhitelist))
	if data.AddressWhitelist != nil {
		for i, entry := range data.AddressWhitelist {
			wlAddress, err := sdk.AccAddressFromBech32(entry)
			if err != nil {
				panic(err)
			}
			wl[i] = wlAddress
		}

		k.SetClpWhiteList(ctx, wl)
	}

	k.SetClpWhiteList(ctx, wl)
	for _, pool := range data.PoolList {
		err := k.SetPool(ctx, pool)
		if err != nil {
			panic(fmt.Sprintf("Pool could not be set : %s", pool.String()))
		}
	}
	for _, lp := range data.LiquidityProviders {
		k.SetLiquidityProvider(ctx, lp)
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)

	poolList := keeper.GetPools(ctx)

	liquidityProviders := keeper.GetLiquidityProviders(ctx)

	whiteList := keeper.GetClpWhiteList(ctx)

	wl := make([]string, len(whiteList))
	for i, entry := range whiteList {
		wl[i] = entry.String()
	}

	return types.GenesisState{
		Params:             params,
		AddressWhitelist:   wl,
		PoolList:           poolList,
		LiquidityProviders: liquidityProviders,
	}
}

// ValidateGenesis validates the clp genesis parameters
func ValidateGenesis(data types.GenesisState) error {
	if !data.Params.Validate() {
		return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("Params are invalid : %s", data.Params.String()))
	}
	for _, pool := range data.PoolList {
		if !pool.Validate() {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("Pool is invalid : %s", pool.String()))
		}
	}
	for _, lp := range data.LiquidityProviders {
		if !lp.Validate() {
			return sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("LiquidityProvider is invalid : %s", lp.String()))
		}
	}
	return nil
}
