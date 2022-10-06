package dispensation

import (
	"fmt"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	// Verify mintTokens
	if !k.TokensCanBeMinted(ctx) {
		return
	}
	mintAmount, ok := sdk.NewIntFromString(types.MintAmountPerBlock)
	if !ok {
		ctx.Logger().Error("Unable to get mint amount")
		return
	}

	if k.IsLastBlock(ctx) {
		maxMintAmount, ok := sdk.NewIntFromString(types.MaxMintAmount)
		if !ok {
			ctx.Logger().Error("Unable to get max mint amount")
			return
		}
		controller, found := k.GetMintController(ctx)
		if !found {
			ctx.Logger().Error(types.ErrNotFoundMintController.Error())
			return
		}
		mintAmount = maxMintAmount.Sub(controller.TotalCounter.Amount)
	}

	mintCoins := sdk.NewCoins(sdk.NewCoin(clptypes.GetSettlementAsset().Symbol, mintAmount))
	if !mintCoins.IsValid() || mintCoins.Len() != 1 {
		ctx.Logger().Error(fmt.Sprintf("Trying to mint invalid coins %v", mintCoins))
		return
	}
	// Get Ecosystem Pool Address
	ecoPoolAddress, err := sdk.AccAddressFromBech32(types.EcoPool)
	if err != nil {
		ctx.Logger().Error("Unable to get address")
		return
	}
	// Mint Tokens
	err = k.GetBankKeeper().MintCoins(ctx, ModuleName, mintCoins)
	if err != nil {
		ctx.Logger().Error("Unable to mint coins")
		return
	}
	// Send newly minted tokens to EcosystemPool
	err = k.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, ecoPoolAddress, mintCoins)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Unable to send %s coins to address %s", mintCoins.String(), ecoPoolAddress.String()))
	}
	err = k.AddMintAmount(ctx, mintCoins[0])
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("Unable to set Mint controller | Err : %s", err.Error()))
	}
}
