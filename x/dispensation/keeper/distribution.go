package keeper

import (
	"fmt"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

// This package adds set and get operations for distributions

func (k Keeper) SetDistribution(ctx sdk.Context, ar types.Distribution) error {
	if !ar.Validate() {
		return errors.Wrapf(types.ErrInvalid, "Record Details : %s", ar.String())
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionsKey(ar.DistributionName, ar.DistributionType, ar.Runner)
	store.Set(key, k.cdc.MustMarshal(&ar))
	return nil
}

func (k Keeper) GetDistribution(ctx sdk.Context, name string, distributionType types.DistributionType, authorizedRunner string) (*types.Distribution, error) {
	var ar types.Distribution
	store := ctx.KVStore(k.storeKey)
	key := types.GetDistributionsKey(name, distributionType, authorizedRunner)
	if !k.Exists(ctx, key) {
		return &ar, errors.Wrapf(types.ErrInvalid, "Record Does not Exist : %s", ar.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &ar)
	return &ar, nil
}

func (k Keeper) ExistsDistribution(ctx sdk.Context, name string, distributionType types.DistributionType, authorizedRunner string) bool {
	key := types.GetDistributionsKey(name, distributionType, authorizedRunner)
	return k.Exists(ctx, key)
}

func (k Keeper) GetDistributionIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DistributionsPrefix)
}

func (k Keeper) GetDistributions(ctx sdk.Context) *types.Distributions {
	var res types.Distributions
	iterator := k.GetDistributionIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dl types.Distribution
		bytesValue := iterator.Value()
		err := k.cdc.Unmarshal(bytesValue, &dl)
		if err != nil {
			// Log unmarshal distribution error instead of panic.
			ctx.Logger().Error(fmt.Sprintf("Unmarshal failed for distribution bytes : %s ", bytesValue))
			continue
		}
		res.Distributions = append(res.Distributions, &dl)
	}
	return &res
}
