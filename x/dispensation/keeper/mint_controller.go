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

func (k Keeper) GetMintController(ctx sdk.Context) types.MintController {
	controller := types.MintController{}
	store := ctx.KVStore(k.storeKey)
	if !k.Exists(ctx, types.MintControllerPrefix) {
		return controller
	}
	bz := store.Get(types.MintControllerPrefix)
	k.cdc.MustUnmarshal(bz, &controller)
	return controller
}

func (k Keeper) AddMintAmount(ctx sdk.Context, mintedCoin sdk.Coin) {
	controller := k.GetMintController(ctx)
	controller.TotalCounter = controller.TotalCounter.Add(mintedCoin)
	k.SetMintController(ctx, controller)
}

func (k Keeper) TokensCanBeMinted(ctx sdk.Context) bool {
	controller := k.GetMintController(ctx)
	maxMintAmount, ok := sdk.NewIntFromString(types.MaxMintAmount)
	if !ok {
		return ok
	}
	return controller.TotalCounter.IsLT(sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, maxMintAmount))
}
