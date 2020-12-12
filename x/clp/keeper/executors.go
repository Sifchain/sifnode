package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

func (k Keeper) CreatePool(ctx sdk.Context, poolUints sdk.Uint, msg types.MsgCreatePool) (*types.Pool, error) {
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, sdk.NewIntFromUint64(msg.ExternalAssetAmount.Uint64()))
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, sdk.NewIntFromUint64(msg.NativeAssetAmount.Uint64()))
	if !k.HasCoins(ctx, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}
	pool, err := types.NewPool(msg.ExternalAsset, msg.NativeAssetAmount, msg.ExternalAssetAmount, poolUints)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}
	// Send coins from user to pool

	err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Signer, types.ModuleName, sdk.Coins{externalAssetCoin, nativeAssetCoin})
	if err != nil {
		return nil, err
	}
	// Pool creator becomes the first LP
	err = k.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	return &pool, nil
}

func (k Keeper) CreateLiquidityProvider(ctx sdk.Context, asset types.Asset, lpunits sdk.Uint, lpaddress sdk.AccAddress) types.LiquidityProvider {
	lp := types.NewLiquidityProvider(asset, lpunits, lpaddress)
	k.SetLiquidityProvider(ctx, lp)
	return lp
}

func (k Keeper) AddLiquidity(ctx sdk.Context, msg types.MsgAddLiquidity, pool types.Pool, newPoolUnits sdk.Uint, lpUnits sdk.Uint) (*types.LiquidityProvider, error) {

	// Verify user has coins to add liquidity
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, sdk.NewIntFromUint64(msg.ExternalAssetAmount.Uint64()))
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, sdk.NewIntFromUint64(msg.NativeAssetAmount.Uint64()))
	if !k.HasCoins(ctx, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}
	// Send from user to pool
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Signer, types.ModuleName, sdk.Coins{externalAssetCoin, nativeAssetCoin})
	if err != nil {
		return nil, err
	}

	pool.PoolUnits = newPoolUnits
	pool.NativeAssetBalance = pool.NativeAssetBalance.Add(msg.NativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(msg.ExternalAssetAmount)

	// Create new Liquidity provider or add liquidity units
	lp, err := k.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer.String())
	if err != nil {
		lp = k.CreateLiquidityProvider(ctx, msg.ExternalAsset, lpUnits, msg.Signer)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeCreateLiquidityProvider,
				sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
		})
		lpUnits = sdk.ZeroUint()
	}
	lp.LiquidityProviderUnits = lp.LiquidityProviderUnits.Add(lpUnits)
	// Save new pool balances
	err = k.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Save LP
	k.SetLiquidityProvider(ctx, lp)
	return &lp, err
}

func (k Keeper) RemoveLiquidityProvider(ctx sdk.Context, coins sdk.Coins, lp types.LiquidityProvider) error {
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lp.LiquidityProviderAddress, coins)
	if err != nil {
		return errors.Wrap(types.ErrUnableToAddBalance, err.Error())
	}
	k.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress.String())
	return nil
}

func (k Keeper) DecommissionPool(ctx sdk.Context, pool types.Pool) error {
	err := k.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	if err != nil {
		return errors.Wrap(types.ErrUnableToDestroyPool, err.Error())
	}
	return nil
}

func (k Keeper) RemoveLiquidity(ctx sdk.Context, pool types.Pool, externalAssetCoin sdk.Coin,
	nativeAssetCoin sdk.Coin, lp types.LiquidityProvider, lpUnitsLeft, poolOriginalEB, poolOriginalNB sdk.Uint) error {

	sendCoins := sdk.Coins{}
	if !externalAssetCoin.IsZero() && !externalAssetCoin.IsNegative() {
		sendCoins = sendCoins.Add(externalAssetCoin)
	}

	if !nativeAssetCoin.IsZero() && !nativeAssetCoin.IsNegative() {
		sendCoins = sendCoins.Add(nativeAssetCoin)
	}
	// Verify if Swap makes the pool too shallow in one of the assets
	if externalAssetCoin.Amount.GTE(sdk.Int(poolOriginalEB)) || nativeAssetCoin.Amount.GTE(sdk.Int(poolOriginalNB)) {
		return errors.Wrap(types.ErrPoolTooShallow, "Pool Balance nil after adjusting asymmetry")
	}
	err := k.SetPool(ctx, pool)
	if err != nil {
		return errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Send coins from pool to user
	if !sendCoins.Empty() {
		if !k.HasCoins(ctx, types.GetCLPModuleAddress(), sendCoins) {
			return types.ErrNotEnoughLiquidity
		}
		err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lp.LiquidityProviderAddress, sendCoins)
		if err != nil {
			return err
		}
	}

	if lpUnitsLeft.IsZero() {
		k.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress.String())
	} else {
		lp.LiquidityProviderUnits = lpUnitsLeft
		k.SetLiquidityProvider(ctx, lp)
	}
	return nil
}

func (k Keeper) InitiateSwap(ctx sdk.Context, sentCoin sdk.Coin, swapper sdk.AccAddress) error {
	if !k.HasCoins(ctx, swapper, sdk.Coins{sentCoin}) {
		return types.ErrBalanceNotAvailable
	}
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, swapper, types.ModuleName, sdk.Coins{sentCoin})
	if err != nil {
		return err
	}
	return nil

}
func (k Keeper) FinalizeSwap(ctx sdk.Context, sentAmount sdk.Uint, finalPool types.Pool, msg types.MsgSwap) error {
	err := k.SetPool(ctx, finalPool)
	if err != nil {
		return errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Adding balance to users account ,Received Asset is the asset the user wants to receive
	// Case 1 . Adding his ETH and deducting from  RWN:ETH pool
	// Case 2 , Adding his XCT and deducting from  RWN:XCT pool
	sentCoin := sdk.NewCoin(msg.ReceivedAsset.Symbol, sdk.NewIntFromUint64(sentAmount.Uint64()))
	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.Signer, sdk.Coins{sentCoin})
	if err != nil {
		return err
	}
	return nil
}
