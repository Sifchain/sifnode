package dispensation

import (
	"fmt"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	// Verify mintTokens
	mintAmount, ok := sdk.NewIntFromString(types.MintAmountPerBlock)
	if !ok {
		ctx.Logger().Error("Unable to get mint amount")
		return
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
		panic(fmt.Sprintf("Unable to send %s coins to address %s", mintCoins.String(), ecoPoolAddress.String()))
	}
	k.AddMintAmount(ctx, mintCoins[0])
}
