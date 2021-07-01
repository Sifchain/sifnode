package keeper

import (
	"bytes"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetCethReceiverAccount(ctx sdk.Context, nativeTokenReceiverAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.CethReceiverAccountPrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(&protobuftypes.StringValue{Value: nativeTokenReceiverAccount.String()}))
}

func (k Keeper) IsCethReceiverAccount(ctx sdk.Context, nativeTokenReceiverAccount sdk.AccAddress) bool {
	account := k.GetCethReceiverAccount(ctx)
	return bytes.Equal(account, nativeTokenReceiverAccount)
}

func (k Keeper) IsCethReceiverAccountSet(ctx sdk.Context) bool {
	account := k.GetCethReceiverAccount(ctx)
	return account != nil
}

func (k Keeper) GetCethReceiverAccount(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	key := types.CethReceiverAccountPrefix
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}

	strProto := &protobuftypes.StringValue{}
	k.cdc.MustUnmarshalBinaryBare(bz, strProto)

	if strProto.Value == "" {
		return nil
	}

	accAddress, err := sdk.AccAddressFromBech32(strProto.Value)
	if err != nil {
		ctx.Logger().Error(err.Error(), "error decoding native_tokenreceiveaccount")
		return nil
	}

	return accAddress
}
