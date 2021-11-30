package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	bankKeeper types.BankKeeper
	clpKeeper  types.CLPKeeper
	paramStore paramtypes.Subspace
}

func (k Keeper) SetMTP(ctx sdk.Context, mtp *types.MTP) error {
	if err := mtp.Validate(); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetMTPKey(mtp.Asset, mtp.Address)
	store.Set(key, k.cdc.MustMarshal(mtp))
	return nil
}

func (k Keeper) GetMTP(ctx sdk.Context, symbol string, mtpAddress string) (types.MTP, error) {
	var mtp types.MTP
	key := types.GetMTPKey(symbol, mtpAddress)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return mtp, types.ErrMTPDoesNotExist
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &mtp)
	return mtp, nil
}

func (k Keeper) GetMTPIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.MTPPrefix)
}

func (k Keeper) GetMTPs(ctx sdk.Context) []*types.MTP {
	var mtpList []*types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		mtpList = append(mtpList, &mtp)
	}
	return mtpList
}

func (k Keeper) GetMTPsForAsset(ctx sdk.Context, asset string) []*types.MTP {
	var mtpList []*types.MTP
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		if mtp.Asset == asset {
			mtpList = append(mtpList, &mtp)
		}
	}
	return mtpList
}

func (k Keeper) GetAssetsForMTP(ctx sdk.Context, mtpAddress sdk.Address) []string {
	var assetList []string
	iterator := k.GetMTPIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var mtp types.MTP
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &mtp)
		if mtpAddress.String() == mtp.Address {
			assetList = append(assetList, mtp.Asset)
		}
	}
	return assetList
}

func (k Keeper) DestroyMTP(ctx sdk.Context, symbol string, mtpAddress string) error {
	key := types.GetMTPKey(symbol, mtpAddress)
	store := ctx.KVStore(k.storeKey)
	if !store.Has(key) {
		return types.ErrMTPDoesNotExist
	}
	store.Delete(key)
	return nil
}
