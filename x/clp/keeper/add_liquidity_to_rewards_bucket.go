package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddLiquidityToRewardsBucket(ctx sdk.Context, signer string, amounts sdk.Coins) (sdk.Coins, error) {
	addr, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return nil, err
	}

	// check that the sender has all the coins in the wallet
	for _, coin := range amounts {
		if !k.bankKeeper.HasBalance(ctx, addr, coin) {
			return nil, types.ErrBalanceNotAvailable
		}
	}

	// send from user to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, amounts); err != nil {
		return nil, err
	}

	// add multiple coins to rewards buckets
	addedCoins, err := k.AddMultipleCoinsToRewardsBuckets(ctx, amounts)
	if err != nil {
		return nil, err
	}

	return addedCoins, nil
}
