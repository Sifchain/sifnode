package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/sifnode/types"
    "github.com/cosmos/cosmos-sdk/codec"
)

// CreateUser creates a user
func (k Keeper) CreateUser(ctx sdk.Context, user types.User) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.UserPrefix + user.ID)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(user)
	store.Set(key, value)
}

// GetUser returns the user information
func (k Keeper) GetUser(ctx sdk.Context, key string) (types.User, error) {
	store := ctx.KVStore(k.storeKey)
	var user types.User
	byteKey := []byte(types.UserPrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// SetUser sets a user
func (k Keeper) SetUser(ctx sdk.Context, user types.User) {
	userKey := user.ID
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(user)
	key := []byte(types.UserPrefix + userKey)
	store.Set(key, bz)
}

// DeleteUser deletes a user
func (k Keeper) DeleteUser(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(types.UserPrefix + key))
}

//
// Functions used by querier
//

func listUser(ctx sdk.Context, k Keeper) ([]byte, error) {
	var userList []types.User
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.UserPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var user types.User
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &user)
		userList = append(userList, user)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, userList)
	return res, nil
}

func getUser(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	user, err := k.GetUser(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, user)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// Get creator of the item
func (k Keeper) GetUserOwner(ctx sdk.Context, key string) sdk.AccAddress {
	user, err := k.GetUser(ctx, key)
	if err != nil {
		return nil
	}
	return user.Creator
}


// Check if the key exists in the store
func (k Keeper) UserExists(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.UserPrefix + key))
}
