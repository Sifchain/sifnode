package keeper

import (
	"errors"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func (k Keeper) CreatePool(ctx sdk.Context, poolUints sdk.Uint, msg *types.MsgCreatePool) (*types.Pool, error) {
	// Defensive programming
	if msg == nil {
		return nil, errors.New("MsgCreatePool can not be nil")
	}
	extInt, ok := k.ParseToInt(msg.ExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}

	nativeInt, ok := k.ParseToInt(msg.NativeAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}

	addr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, extInt)
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, nativeInt)
	if !k.bankKeeper.HasBalance(ctx, addr, externalAssetCoin) && !k.bankKeeper.HasBalance(ctx, addr, nativeAssetCoin) {
		return nil, types.ErrBalanceNotAvailable
	}

	pool, err := types.NewPool(msg.ExternalAsset, msg.NativeAssetAmount, msg.ExternalAssetAmount, poolUints)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}

	// Send coins from user to pool
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{externalAssetCoin, nativeAssetCoin})
	if err != nil {
		return nil, err
	}

	// Pool creator becomes the first LP
	err = k.SetPool(ctx, &pool)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}

	return &pool, nil
}

func (k Keeper) CreateLiquidityProvider(ctx sdk.Context, asset *types.Asset, lpunits sdk.Uint, lpaddress fmt.Stringer) types.LiquidityProvider {
	lp := types.NewLiquidityProvider(asset, lpunits, lpaddress)
	k.SetLiquidityProvider(ctx, &lp)

	return lp
}

func (k Keeper) AddLiquidity(ctx sdk.Context, msg *types.MsgAddLiquidity, pool types.Pool, newPoolUnits sdk.Uint, lpUnits sdk.Uint) (*types.LiquidityProvider, error) {

	// Verify user has coins to add liquidiy
	extInt, ok := k.ParseToInt(msg.ExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}

	nativeInt, ok := k.ParseToInt(msg.NativeAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}

	var coins sdk.Coins
	if extInt != sdk.ZeroInt() {
		externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, extInt)
		coins = coins.Add(externalAssetCoin)
	}

	if nativeInt != sdk.ZeroInt() {
		nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, nativeInt)
		coins = coins.Add(nativeAssetCoin)
	}

	addr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	if !k.bankKeeper.HasBalance(ctx, addr, coins[0]) && !k.bankKeeper.HasBalance(ctx, addr, coins[1]) {
		return nil, types.ErrBalanceNotAvailable
	}
	// Send from user to pool
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	pool.PoolUnits = newPoolUnits
	pool.NativeAssetBalance = pool.NativeAssetBalance.Add(msg.NativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(msg.ExternalAssetAmount)

	// Create new Liquidity provider or add liquidity units
	lp, err := k.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer)
	if err != nil {
		lp = k.CreateLiquidityProvider(ctx, msg.ExternalAsset, lpUnits, addr)
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
	err = k.SetPool(ctx, &pool)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Save LP
	k.SetLiquidityProvider(ctx, &lp)
	return &lp, err
}

func (k Keeper) RemoveLiquidityProvider(ctx sdk.Context, coins sdk.Coins, lp types.LiquidityProvider) error {
	lpaddr, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lpaddr, coins)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToAddBalance, err.Error())
	}
	k.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	return nil
}

func (k Keeper) DecommissionPool(ctx sdk.Context, pool types.Pool) error {
	err := k.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToDestroyPool, err.Error())
	}
	return nil
}

func (k Keeper) RemoveLiquidity(ctx sdk.Context, pool types.Pool, externalAssetCoin sdk.Coin,
	nativeAssetCoin sdk.Coin, lp types.LiquidityProvider, lpUnitsLeft, poolOriginalEB, poolOriginalNB sdk.Uint) error {
	lpAddr, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	if err != nil {
		return err
	}

	sendCoins := sdk.Coins{}
	if !externalAssetCoin.IsZero() && !externalAssetCoin.IsNegative() {
		sendCoins = sendCoins.Add(externalAssetCoin)
	}

	if !nativeAssetCoin.IsZero() && !nativeAssetCoin.IsNegative() {
		sendCoins = sendCoins.Add(nativeAssetCoin)
	}
	// Verify if Swap makes the pool too shallow in one of the assets
	if externalAssetCoin.Amount.GTE(sdk.Int(poolOriginalEB)) || nativeAssetCoin.Amount.GTE(sdk.Int(poolOriginalNB)) {
		return sdkerrors.Wrap(types.ErrPoolTooShallow, "Pool Balance nil after adjusting asymmetry")
	}
	err = k.SetPool(ctx, &pool)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Send coins from pool to user
	if !sendCoins.Empty() {
		for _, coin := range sendCoins {
			if !k.bankKeeper.HasBalance(ctx, types.GetCLPModuleAddress(), coin) {
				return types.ErrNotEnoughLiquidity
			}
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lpAddr, sendCoins)
		if err != nil {
			return err
		}
	}

	if lpUnitsLeft.IsZero() {
		k.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	} else {
		lp.LiquidityProviderUnits = lpUnitsLeft
		k.SetLiquidityProvider(ctx, &lp)
	}
	return nil
}

func (k Keeper) InitiateSwap(ctx sdk.Context, sentCoin sdk.Coin, swapper sdk.AccAddress) error {
	if !k.bankKeeper.HasBalance(ctx, swapper, sentCoin) {
		return types.ErrBalanceNotAvailable
	}
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, swapper, types.ModuleName, sdk.Coins{sentCoin})
	if err != nil {
		return err
	}
	return nil

}
func (k Keeper) FinalizeSwap(ctx sdk.Context, sentAmount string, finalPool types.Pool, msg types.MsgSwap) error {
	err := k.SetPool(ctx, &finalPool)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	sentAmountInt, ok := k.ParseToInt(sentAmount)
	if !ok {
		return types.ErrUnableToParseInt
	}
	// Adding balance to users account ,Received Asset is the asset the user wants to receive
	// Case 1 . Adding his ETH and deducting from  RWN:ETH pool
	// Case 2 , Adding his XCT and deducting from  RWN:XCT pool

	addr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}
	sentCoin := sdk.NewCoin(msg.ReceivedAsset.Symbol, sentAmountInt)
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{sentCoin})
	if err != nil {
		return err
	}
	return nil
}

// Use strings instead of Unit/Int in between conventions
func (k Keeper) ParseToInt(nu string) (sdk.Int, bool) {
	return sdk.NewIntFromString(nu)
}
