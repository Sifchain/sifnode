package keeper

import (
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetEthGasPrice(ctx sdk.Context, ethGasPrice sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := types.EthGasPricePrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(ethGasPrice))
}

func (k Keeper) IsEthGasPriceSet(ctx sdk.Context) bool {
	ethGasPrice := k.GetEthGasPrice(ctx)
	return ethGasPrice != nil
}

func (k Keeper) GetEthGasPrice(ctx sdk.Context) (ethGasPrice *sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := types.EthGasPricePrefix
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &ethGasPrice)
	return
}
