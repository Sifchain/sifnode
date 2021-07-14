package keeper

import (
	"bytes"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetCrossChainFeeReceiverAccount(ctx sdk.Context, crossChainFeeReceiverAccount sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.CrossChainFeeReceiverAccountPrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(&protobuftypes.StringValue{Value: crossChainFeeReceiverAccount.String()}))
}

func (k Keeper) IsCrossChainFeeReceiverAccount(ctx sdk.Context, crossChainFeeReceiverAccount sdk.AccAddress) bool {
	account := k.GetCrossChainFeeReceiverAccount(ctx)
	return bytes.Equal(account, crossChainFeeReceiverAccount)
}

func (k Keeper) IsCrossChainFeeReceiverAccountSet(ctx sdk.Context) bool {
	account := k.GetCrossChainFeeReceiverAccount(ctx)
	return account != nil
}

func (k Keeper) GetCrossChainFeeReceiverAccount(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	key := types.CrossChainFeeReceiverAccountPrefix
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
		ctx.Logger().Error(err.Error(), "error decoding crosschain fee receive account")
		return nil
	}

	return accAddress
}
