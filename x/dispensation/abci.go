package dispensation

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ecoPool = "sif1ct2s3t8u2kffjpaekhtngzv6yc4vm97xajqyl3"
const mintAmountPerBlock = "225000000000000000000"

func BeginBlocker(ctx sdk.Context, k Keeper) {
	mintAmount, ok := sdk.NewIntFromString(mintAmountPerBlock)
	if !ok {
		ctx.Logger().Error("Unable to get mint amount")
		return
	}
	ecoPoolAddress, err := sdk.AccAddressFromBech32(ecoPool)
	if err != nil {
		ctx.Logger().Error("Unable to get address")
		return
	}
	mintCoins := sdk.NewCoins(sdk.NewCoin("rowan", mintAmount))
	err = k.GetBankKeeper().MintCoins(ctx, ModuleName, mintCoins)
	if err != nil {
		ctx.Logger().Error("Unable to mint coins")
		return
	}
	err = k.GetBankKeeper().SendCoinsFromModuleToAccount(ctx, ModuleName, ecoPoolAddress, mintCoins)
	if err != nil {
		panic(fmt.Sprintf("Unable to send %s coins to address %s", mintCoins.String(), ecoPoolAddress.String()))
	}
}
