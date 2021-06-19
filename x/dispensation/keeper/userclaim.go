package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (k Keeper) SetClaim(ctx sdk.Context, ar types.UserClaim) error {
	if !ar.Validate() {
		return errors.Wrapf(types.ErrInvalid, "Claim Details : %s", ar.String())
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserClaimKey(ar.UserAddress, ar.UserClaimType)
	store.Set(key, k.cdc.MustMarshalBinaryBare(&ar))
	return nil
}

func (k Keeper) GetClaim(ctx sdk.Context, recipient string, userClaimType types.DistributionType) (types.UserClaim, error) {
	var ar types.UserClaim
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserClaimKey(recipient, userClaimType)
	if !k.Exists(ctx, key) {
		return ar, errors.Wrapf(types.ErrInvalid, "Claim Does not Exist : %s", ar.String())
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshalBinaryBare(bz, &ar)
	return ar, nil
}

func (k Keeper) ExistsClaim(ctx sdk.Context, recipient string, userClaimType types.DistributionType) bool {
	key := types.GetUserClaimKey(recipient, userClaimType)
	return k.Exists(ctx, key)
}
func (k Keeper) DeleteClaim(ctx sdk.Context, recipient string, userClaimType types.DistributionType) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserClaimKey(recipient, userClaimType)
	store.Delete(key)
}

func (k Keeper) GetUserClaimsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.UserClaimPrefix)
}

func (k Keeper) GetClaims(ctx sdk.Context) []types.UserClaim {
	var res []types.UserClaim
	iterator := k.GetUserClaimsIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic("Failed to close iterator")
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var dl types.UserClaim
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dl)
		res = append(res, dl)
	}
	return res
}

func (k Keeper) GetClaimsByType(ctx sdk.Context, userClaimType types.DistributionType) []types.UserClaim {
	var res []types.UserClaim
	iterator := k.GetUserClaimsIterator(ctx)
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic("Failed to close iterator")
		}
	}(iterator)
	for ; iterator.Valid(); iterator.Next() {
		var dl types.UserClaim
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshalBinaryBare(bytesValue, &dl)
		if dl.UserClaimType == userClaimType {
			res = append(res, dl)
		}
	}
	return res
}
