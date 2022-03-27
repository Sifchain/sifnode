package keeper

import (
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
)

func (keeper Keeper) CreatePool(ctx sdk.Context, poolUints sdk.Uint, msg *types.MsgCreatePool) (*types.Pool, error) {
	// Defensive programming
	if msg == nil {
		return nil, errors.New("MsgCreatePool can not be nil")
	}
	extInt, ok := keeper.ParseToInt(msg.ExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}

	nativeInt, ok := keeper.ParseToInt(msg.NativeAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}

	addr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, extInt)
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, nativeInt)
	if !keeper.bankKeeper.HasBalance(ctx, addr, externalAssetCoin) && !keeper.bankKeeper.HasBalance(ctx, addr, nativeAssetCoin) {
		return nil, types.ErrBalanceNotAvailable
	}
	pool := types.NewPool(msg.ExternalAsset, msg.NativeAssetAmount, msg.ExternalAssetAmount, poolUints)
	// Send coins from user to pool
	err = keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(externalAssetCoin, nativeAssetCoin))
	if err != nil {
		return nil, err
	}
	// Pool creator becomes the first LP
	err = keeper.SetPool(ctx, &pool)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	return &pool, nil
}

func (keeper Keeper) CreateLiquidityProvider(ctx sdk.Context, asset *types.Asset, lpunits sdk.Uint, lpaddress sdk.AccAddress) types.LiquidityProvider {
	lp := types.NewLiquidityProvider(asset, lpunits, lpaddress)
	keeper.SetLiquidityProvider(ctx, &lp)

	return lp
}

func (keeper Keeper) AddLiquidity(ctx sdk.Context, msg *types.MsgAddLiquidity, pool types.Pool, newPoolUnits sdk.Uint, lpUnits sdk.Uint) (*types.LiquidityProvider, error) {

	// Verify user has coins to add liquidiy
	extInt, ok := keeper.ParseToInt(msg.ExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}

	nativeInt, ok := keeper.ParseToInt(msg.NativeAssetAmount.String())
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

	if !keeper.bankKeeper.HasBalance(ctx, addr, coins[0]) && !keeper.bankKeeper.HasBalance(ctx, addr, coins[1]) {
		return nil, types.ErrBalanceNotAvailable
	}
	// Send from user to pool
	err = keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	pool.PoolUnits = newPoolUnits
	pool.NativeAssetBalance = pool.NativeAssetBalance.Add(msg.NativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(msg.ExternalAssetAmount)

	// Create new Liquidity provider or add liquidity units
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer)
	if err != nil {
		lp = keeper.CreateLiquidityProvider(ctx, msg.ExternalAsset, lpUnits, addr)
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
	err = keeper.SetPool(ctx, &pool)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Save LP
	keeper.SetLiquidityProvider(ctx, &lp)
	return &lp, err
}

func (keeper Keeper) RemoveLiquidityProvider(ctx sdk.Context, coins sdk.Coins, lp types.LiquidityProvider) error {
	lpaddr, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	if err != nil {
		return err
	}

	err = keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lpaddr, coins)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToAddBalance, err.Error())
	}
	keeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	return nil
}

func (keeper Keeper) DecommissionPool(ctx sdk.Context, pool types.Pool) error {
	err := keeper.DestroyPool(ctx, pool.ExternalAsset.Symbol)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToDestroyPool, err.Error())
	}
	return nil
}

func (keeper Keeper) RemoveLiquidity(ctx sdk.Context, pool types.Pool, externalAssetCoin sdk.Coin,
	nativeAssetCoin sdk.Coin, lp types.LiquidityProvider, lpUnitsLeft, poolOriginalEB, poolOriginalNB sdk.Uint) error {
	lpAddr, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
	if err != nil {
		return err
	}

	sendCoins := sdk.NewCoins()
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

	err = keeper.SetPool(ctx, &pool)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Send coins from pool to user
	if !sendCoins.Empty() {
		for _, coin := range sendCoins {
			if !keeper.bankKeeper.HasBalance(ctx, types.GetCLPModuleAddress(), coin) {
				return types.ErrNotEnoughLiquidity
			}
		}
		err = keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lpAddr, sendCoins)
		if err != nil {
			return err
		}
	}

	if lpUnitsLeft.IsZero() {
		keeper.DestroyLiquidityProvider(ctx, lp.Asset.Symbol, lp.LiquidityProviderAddress)
	} else {
		lp.LiquidityProviderUnits = lpUnitsLeft
		keeper.SetLiquidityProvider(ctx, &lp)
	}
	return nil
}

func (keeper Keeper) InitiateSwap(ctx sdk.Context, sentCoin sdk.Coin, swapper sdk.AccAddress) error {
	if !keeper.bankKeeper.HasBalance(ctx, swapper, sentCoin) {
		return types.ErrBalanceNotAvailable
	}
	err := keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, swapper, types.ModuleName, sdk.NewCoins(sentCoin))
	if err != nil {
		return err
	}
	return nil

}
func (keeper Keeper) FinalizeSwap(ctx sdk.Context, sentAmount string, finalPool types.Pool, msg types.MsgSwap) error {
	err := keeper.SetPool(ctx, &finalPool)
	if err != nil {
		return sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	sentAmountInt, ok := keeper.ParseToInt(sentAmount)
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
	err = keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(sentCoin))
	if err != nil {
		return err
	}
	return nil
}

// Use strings instead of Unit/Int in between conventions
func (keeper Keeper) ParseToInt(nu string) (sdk.Int, bool) {
	return sdk.NewIntFromString(nu)
}
