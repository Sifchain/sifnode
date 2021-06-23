package keeper

import (
	"strconv"

	"github.com/Sifchain/sifnode/x/trees/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) FillLimitOrders(ctx sdk.Context) {
	// limitOrders := k.GetAllOrdersByTreeId(ctx, strconv.FormatInt(0, 10))
	store := ctx.KVStore(k.storeKey)
	iterator := k.GetLimitOrdersIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var order types.LimitOrder
		var tree types.Tree
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &order)
		k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.TreeKey+order.TreeId)), &tree)
		ctx.Logger().Error(tree.Id)
		ctx.Logger().Error(tree.Price.String())
		ctx.Logger().Error(order.MaxPrice.String())
		if !tree.Status && !order.Executed {
			if tree.Price[0].Amount.Int64() < order.MaxPrice[0].Amount.Int64() {
				tree.Status = true
				store.Set(types.KeyPrefix(types.TreeKey+tree.Id), k.cdc.MustMarshalBinaryBare(tree))
			}
		}
		order.Executed = true
		store.Set(types.GetLimitedOrderKey(tree.Id, order.OrderId), k.cdc.MustMarshalBinaryBare(order))

	}
}

func (k Keeper) CreateLimitOrder(ctx sdk.Context, msg types.MsgBuyTree) (string, error) {
	ctx.Logger().Error("Limit orderKeeper got it")
	count := k.GetOrderCountByTreeId(ctx, msg.Id)
	var order = types.LimitOrder{
		Buyer:    msg.Buyer,
		MaxPrice: msg.Price,
		TreeId:   msg.Id,
		OrderId:  strconv.FormatInt(count, 10),
		Executed: false,
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetLimitedOrderKey(order.TreeId, order.OrderId)
	store.Set(key, k.cdc.MustMarshalBinaryBare(order))
	k.SetOrderCountByTreeId(ctx, order.TreeId, count+1)
	// return tree.Id, nil
	return order.Buyer.String(), nil
}

func (k Keeper) GetOrderCountByTreeId(ctx sdk.Context, id string) int64 {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.GetLimitOrderCountKey(id)
	bz := store.Get(byteKey)
	ctx.Logger().Error(string(bz))
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

func (k Keeper) SetOrderCountByTreeId(ctx sdk.Context, id string, count int64) {
	store := ctx.KVStore(k.storeKey)
	byteKey := types.GetLimitOrderCountKey(id)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

func (k Keeper) GetAllOrdersByTreeId(ctx sdk.Context, treeId string) (msgs []types.LimitOrder) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.OrderKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.LimitOrder
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		if msg.TreeId == treeId {
			msgs = append(msgs, msg)
		}
	}

	return
}

func (k Keeper) GetLimitOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.OrderKey))
}
