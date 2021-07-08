package state

import (
	"github.com/Sifchain/sifnode/x/ibc_sifchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func AddToDenomWhiteList(ctx sdk.Context, cdc codec.BinaryMarshaler, denom string) error {
	whitelist, err := GetWhiteList(ctx, cdc)
	if err != nil {
		return err
	}
	whitelist.DenomWhitelist[denom] = true
	SetWhiteList(ctx, cdc, whitelist)
	return nil
}

func AddToChannelWhiteList(ctx sdk.Context, cdc codec.BinaryMarshaler, channel string, port string) error {
	whitelist, err := GetWhiteList(ctx, cdc)
	if err != nil {
		return err
	}
	whitelist.ChannelPortWhitelist[channel] = port
	SetWhiteList(ctx, cdc, whitelist)
	return nil
}

func SetWhiteList(ctx sdk.Context, cdc codec.BinaryMarshaler, whiteList types.WhiteList) {
	store := ctx.KVStore(sdk.NewKVStoreKey(types.StoreKey))
	key := types.GetWhiteListKey()
	store.Set(key, cdc.MustMarshalBinaryBare(&whiteList))
}

func GetWhiteList(ctx sdk.Context, cdc codec.BinaryMarshaler) (types.WhiteList, error) {
	whiteList := types.WhiteList{}
	store := ctx.KVStore(sdk.NewKVStoreKey(types.StoreKey))
	key := types.GetWhiteListKey()
	if !Exists(ctx, key) {
		return whiteList, errors.New("WhiteList Does not exist")
	}
	bz := store.Get(key)
	cdc.MustUnmarshalBinaryBare(bz, &whiteList)
	return whiteList, nil
}

func Exists(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(sdk.NewKVStoreKey(types.StoreKey))
	return store.Has(key)
}
