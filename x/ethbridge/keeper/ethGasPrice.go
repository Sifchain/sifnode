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

func (k Keeper) SetGasMultiplier(ctx sdk.Context, GasMultiplier sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := types.GasMultiplierPrefix
	store.Set(key, k.cdc.MustMarshalBinaryBare(GasMultiplier))
}

func (k Keeper) IsGasMultiplierSet(ctx sdk.Context) bool {
	GasMultiplier := k.GetGasMultiplier(ctx)
	return GasMultiplier != nil
}

func (k Keeper) GetGasMultiplier(ctx sdk.Context) (GasMultiplier *sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := types.GasMultiplierPrefix
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &GasMultiplier)
	return
}
