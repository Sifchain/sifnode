package keeper

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetMintController(ctx sdk.Context, mintController types.MintController) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.MintControllerPrefix, k.cdc.MustMarshal(&mintController))
}

func (k Keeper) GetMintController(ctx sdk.Context) (types.MintController, bool) {
	controller := types.MintController{}
	store := ctx.KVStore(k.storeKey)
	if !k.Exists(ctx, types.MintControllerPrefix) {
		return controller, false
	}
	bz := store.Get(types.MintControllerPrefix)
	k.cdc.MustUnmarshal(bz, &controller)
	return controller, true
}

func (k Keeper) AddMintAmount(ctx sdk.Context, mintedCoin sdk.Coin) error {
	controller, found := k.GetMintController(ctx)
	if !found {
		return types.ErrNotFoundMintController
	}
	controller.TotalCounter = controller.TotalCounter.Add(mintedCoin)
	k.SetMintController(ctx, controller)
	return nil
}

func (k Keeper) TokensCanBeMinted(ctx sdk.Context) bool {
	controller, found := k.GetMintController(ctx)
	if !found {
		ctx.Logger().Error(types.ErrNotFoundMintController.Error())
		return false
	}
	totalCounter := controller.TotalCounter
	maxMintAmount, ok := sdk.NewIntFromString(types.MaxMintAmount)
	if !ok {
		return ok
	}
	maxMintAmountCoin := sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, maxMintAmount)
	return totalCounter.IsLT(maxMintAmountCoin)
}

func (k Keeper) IsLastBlock(ctx sdk.Context) bool {
	controller, found := k.GetMintController(ctx)
	if !found {
		ctx.Logger().Error(types.ErrNotFoundMintController.Error())
		return false
	}
	totalCounter := controller.TotalCounter.Amount

	maxMintAmount, ok := sdk.NewIntFromString(types.MaxMintAmount)
	if !ok {
		return ok
	}
	blockMintAmount, ok := sdk.NewIntFromString(types.MintAmountPerBlock)
	if !ok {
		return ok
	}
	return maxMintAmount.Sub(totalCounter).LTE(blockMintAmount)
}
