package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddLiquidityToRewardsBucket(ctx sdk.Context, signer string, amounts sdk.Coins) (sdk.Coins, error) {
	// check that the sender has all the coins in the wallet
	for _, coin := range amounts {
		if !k.bankKeeper.HasBalance(ctx, sdk.AccAddress(signer), coin) {
			return nil, types.ErrBalanceNotAvailable
		}
	}

	// send from user to module
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(signer), types.ModuleName, amounts)
	if err != nil {
		return nil, err
	}

	// add multiple coins to rewards buckets
	addedCoins, err := k.AddMultipleCoinsToRewardsBuckets(ctx, amounts)
	if err != nil {
		return nil, err
	}

	return addedCoins, nil
}
