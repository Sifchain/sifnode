package keeper

import (
	"strconv"

	"github.com/Sifchain/sifnode/x/trees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateTree(ctx sdk.Context, msg types.MsgCreateTree) (string, error) {
	ctx.Logger().Error("Keeper got it")
	count := k.GetTreeCount(ctx)
	var tree = types.Tree{
		Name:     msg.Name,
		Seller:   msg.Seller,
		Price:    msg.Price,
		Category: msg.Category,
		Id:       strconv.FormatInt(count, 10),
		Status:   false,
	}
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefix(types.TreeKey + tree.Id)
	store.Set(key, k.cdc.MustMarshalBinaryBare(tree))
	k.SetTreeCount(ctx, count+1)
	return tree.Id, nil
}

// GetTreeCount get the total number of tree
func (k Keeper) GetTreeCount(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(types.TreeCountKey)
	bz := store.Get(byteKey)
	ctx.Logger().Error(string(bz))

	ctx.Logger().Error("Keeper got it")

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseInt(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to int64
		panic("cannot decode count")
	}

	return count
}

// SetTreeCount set the total number of tree
func (k Keeper) SetTreeCount(ctx sdk.Context, count int64) {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.KeyPrefix(types.TreeCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

func (k Keeper) GetTree(ctx sdk.Context, key string) types.Tree {
	store := ctx.KVStore(k.storeKey)
	var tree types.Tree
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.TreeKey+key)), &tree)
	return tree
}

func (k Keeper) HasTree(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.KeyPrefix(types.TreeKey + id))
}

func (k Keeper) GetTreeOwner(ctx sdk.Context, key string) string {
	return k.GetTree(ctx, key).Seller.String()
}

func (k Keeper) GetAllTrees(ctx sdk.Context) (msgs []types.Tree) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.TreeKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Tree
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}
