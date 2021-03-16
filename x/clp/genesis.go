package clp

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	keeper.SetParams(ctx, data.Params)
	if data.AddressWhitelist != nil {
		keeper.SetClpWhiteList(ctx, data.AddressWhitelist)
	}
	for _, pool := range data.PoolList {
		err := keeper.SetPool(ctx, pool)
		if err != nil {
			panic(fmt.Sprintf("Pool could not be set : %s", pool.String()))
		}
	}
	for _, lp := range data.LiquidityProviders {
		keeper.SetLiquidityProvider(ctx, lp)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	poolList := keeper.GetPools(ctx)
	liquidityProviders := keeper.GetLiquidityProviders(ctx)
	whiteList := keeper.GetClpWhiteList(ctx)
	wl := make([]string, 0, len(whiteList))
	for i, entry := range whiteList {
		wl[i] = entry.String()
	}

	return types.GenesisState{
		Params:             &params,
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
