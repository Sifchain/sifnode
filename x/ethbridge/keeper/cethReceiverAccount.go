package keeper

import (
	"bytes"
	protobuftypes "github.com/gogo/protobuf/types"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetCethReceiverAccount(ctx sdk.Context, cethReceiverAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.CethReceiverAccountPrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(&protobuftypes.StringValue{Value: cethReceiverAccount.String()}))
}

func (k Keeper) IsCethReceiverAccount(ctx sdk.Context, cethReceiverAccount sdk.AccAddress) bool {
	account := k.GetCethReceiverAccount(ctx)
	return bytes.Equal(account, cethReceiverAccount)
}

func (k Keeper) IsCethReceiverAccountSet(ctx sdk.Context) bool {
	account := k.GetCethReceiverAccount(ctx)
	return account != nil
}

func (k Keeper) GetCethReceiverAccount(ctx sdk.Context) (cethReceiverAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.CethReceiverAccountPrefix
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &protobuftypes.StringValue{Value: cethReceiverAccount.String()})
	return
}
