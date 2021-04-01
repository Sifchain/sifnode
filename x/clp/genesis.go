package clp

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	keeper.SetParams(ctx, data.Params)
	if data.AddressWhitelist == nil || len(data.AddressWhiteList) == 0 {
	    panic(fmt.Sprintf("AddressWhiteList must be set."))
	}
	keeper.SetClpWhiteList(ctx, data.AddressWhitelist)
	for _, pool := range data.PoolList {
		err := keeper.SetPool(ctx, pool)
		if err != nil {
			panic(fmt.Sprintf("Pool could not be set : %s", pool.String()))
		}
	}
	for _, lp := range data.LiquidityProviderList {
		keeper.SetLiquidityProvider(ctx, lp)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	poolList := keeper.GetPools(ctx)
	liquidityProviders := keeper.GetLiquidityProviders(ctx)
	whiteList := keeper.GetClpWhiteList(ctx)
	return GenesisState{
		Params:                params,
		AddressWhitelist:      whiteList,
		PoolList:              poolList,
		LiquidityProviderList: liquidityProviders,
	}
}

// ValidateGenesis validates the clp genesis parameters
func ValidateGenesis(data GenesisState) error {
	if !data.Params.Validate() {
		return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Params are invalid : %s", data.Params.String()))
	}
	for _, pool := range data.PoolList {
		if !pool.Validate() {
			return errors.Wrap(types.ErrInvalid, fmt.Sprintf("Pool is invalid : %s", pool.String()))
		}
	}
	for _, lp := range data.LiquidityProviderList {
		if !lp.Validate() {
			return errors.Wrap(types.ErrInvalid, fmt.Sprintf("LiquidityProvider is invalid : %s", lp.String()))
		}
	}
	return nil
}
