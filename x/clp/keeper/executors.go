package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

func CreatePool(ctx sdk.Context, keeper Keeper, poolUints sdk.Uint, msg types.MsgCreatePool) (*types.Pool, error) {
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(msg.ExternalAssetAmount.Uint64()))
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Ticker, sdk.NewIntFromUint64(msg.NativeAssetAmount.Uint64()))
	if !keeper.HasCoins(ctx, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}
	pool, err := types.NewPool(msg.ExternalAsset, msg.NativeAssetAmount, msg.ExternalAssetAmount, poolUints)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}
	// Send coins from suer to pool
	err = keeper.SendCoins(ctx, msg.Signer, pool.PoolAddress, sdk.Coins{externalAssetCoin, nativeAssetCoin})
	if err != nil {
		return nil, err
	}
	// Pool creator becomes the first LP
	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	return &pool, nil
}

func CreateLiquidityProvider(ctx sdk.Context, asset types.Asset, keeper Keeper, lpunits sdk.Uint, lpaddress sdk.AccAddress) types.LiquidityProvider {
	lp := types.NewLiquidityProvider(asset, lpunits, lpaddress)
	keeper.SetLiquidityProvider(ctx, lp)
	return lp
}

func AddLiquidity(ctx sdk.Context, keeper Keeper, msg types.MsgAddLiquidity, pool types.Pool, newPoolUnits sdk.Uint, lpUnits sdk.Uint) (*types.LiquidityProvider, error) {

	// Verify user has coins to add liquidity
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(msg.ExternalAssetAmount.Uint64()))
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Ticker, sdk.NewIntFromUint64(msg.NativeAssetAmount.Uint64()))
	if !keeper.HasCoins(ctx, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}
	// Send from user to pool
	err := keeper.SendCoins(ctx, msg.Signer, pool.PoolAddress, sdk.Coins{externalAssetCoin, nativeAssetCoin})
	if err != nil {
		return nil, err
	}

	pool.PoolUnits = newPoolUnits
	pool.NativeAssetBalance = pool.NativeAssetBalance.Add(msg.NativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(msg.ExternalAssetAmount)

	// Create new Liquidity provider or add liquidity units
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Ticker, msg.Signer.String())
	if err != nil {
		lp = CreateLiquidityProvider(ctx, msg.ExternalAsset, keeper, lpUnits, msg.Signer)
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
	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Save LP
	keeper.SetLiquidityProvider(ctx, lp)
	return &lp, err
}

func RemoveLiquidityProvider(ctx sdk.Context, keeper Keeper, coins sdk.Coins, pool types.Pool, lp types.LiquidityProvider) error {
	err := keeper.SendCoins(ctx, pool.PoolAddress, lp.LiquidityProviderAddress, coins)
	if err != nil {
		return errors.Wrap(types.ErrUnableToAddBalance, err.Error())
	}
	DestroyLiquidityProvider(ctx, keeper, lp)
	return nil
}

func DestroyLiquidityProvider(ctx sdk.Context, keeper Keeper, lp types.LiquidityProvider) {
	keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
}

func DestroyPool(ctx sdk.Context, keeper Keeper, pool types.Pool) error {
	err := keeper.DestroyPool(ctx, pool.ExternalAsset.Ticker)
	if err != nil {
		return errors.Wrap(types.ErrUnableToDestroyPool, err.Error())
	}
	return nil
}

func RemoveLiquidity(ctx sdk.Context, keeper Keeper, pool types.Pool, externalAssetCoin sdk.Coin,
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
	err := keeper.SetPool(ctx, pool)
	if err != nil {
		return errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Send coins from pool to user
	if !sendCoins.Empty() {
		if !keeper.HasCoins(ctx, pool.PoolAddress, sendCoins) {
			return types.ErrNotEnoughLiquidity
		}
		err = keeper.SendCoins(ctx, pool.PoolAddress, lp.LiquidityProviderAddress, sendCoins)
		if err != nil {
			return err
		}
	}

	if lpUnitsLeft.IsZero() {
		DestroyLiquidityProvider(ctx, keeper, lp)
	} else {
		lp.LiquidityProviderUnits = lpUnitsLeft
		keeper.SetLiquidityProvider(ctx, lp)
	}
	return nil
}

func FinalizeSwap(ctx sdk.Context, keeper Keeper, sentAmount sdk.Uint, finalPool types.Pool, outPool types.Pool, msg types.MsgSwap) error {
	err := keeper.SetPool(ctx, finalPool)
	if err != nil {
		return errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Adding balance to users account ,Received Asset is the asset the user wants to receive
	// Case 1 . Adding his ETH and deducting from  RWN:ETH pool
	// Case 2 , Adding his XCT and deducting from  RWN:XCT pool
	sentCoin := sdk.NewCoin(msg.ReceivedAsset.Ticker, sdk.NewIntFromUint64(sentAmount.Uint64()))
	err = keeper.SendCoins(ctx, outPool.PoolAddress, msg.Signer, sdk.Coins{sentCoin})
	if err != nil {
		return err
	}
	return nil
}
