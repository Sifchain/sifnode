package clp

import (
	"fmt"
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"strconv"
)

// NewHandler creates an sdk.Handler for all the clp type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreatePool:
			return handleMsgCreatePool(ctx, k, msg)
		case MsgDecommissionPool:
			return handleMsgDecommissionPool(ctx, k, msg)
		case MsgAddLiquidity:
			return handleMsgAddLiquidity(ctx, k, msg)
		case MsgRemoveLiquidity:
			return handleMsgRemoveLiquidity(ctx, k, msg)
		case MsgSwap:
			return handleMsgSwap(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgDecommissionPool(ctx sdk.Context, keeper Keeper, msg MsgDecommissionPool) (*sdk.Result, error) {
	// Verify pool
	pool, err := keeper.GetPool(ctx, msg.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	if pool.NativeAssetBalance.GTE(sdk.NewUint(uint64(keeper.GetParams(ctx).MinCreatePoolThreshold))) {
		return nil, types.ErrBalanceTooHigh
	}
	// Get all LP's for the pool
	lpList := keeper.GetLiqudityProvidersForAsset(ctx, pool.ExternalAsset)
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	// iterate over Lp list and refund them there tokens
	// Return both RWN and EXTERNAL ASSET
	for _, lp := range lpList {
		withdrawNativeAsset, withdrawExternalAsset, _, _ := CalculateWithdrawal(poolUnits, nativeAssetBalance.String(), externalAssetBalance.String(),
			lp.LiquidityProviderUnits.String(), sdk.NewInt(10000).String(), sdk.ZeroInt())
		poolUnits = poolUnits.Sub(lp.LiquidityProviderUnits)
		nativeAssetBalance = nativeAssetBalance.Sub(withdrawNativeAsset)
		externalAssetBalance = externalAssetBalance.Sub(withdrawExternalAsset)
		withdrawNativeCoins := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(withdrawNativeAsset.Uint64()))
		withdrawExternalCoins := sdk.NewCoin(msg.Ticker, sdk.NewIntFromUint64(withdrawExternalAsset.Uint64()))
		err = keeper.BankKeeper.SendCoins(ctx, pool.PoolAddress, lp.LiquidityProviderAddress, sdk.Coins{withdrawExternalCoins, withdrawNativeCoins})
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToAddBalance, err.Error())
		}
		keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	}
	// Pool should be empty at this point
	// Destroy the pool
	err = keeper.DestroyPool(ctx, pool.ExternalAsset.Ticker)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToDestroyPool, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDecommissionPool,
			sdk.NewAttribute(types.AttributeKeyPool, pool.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreatePool(ctx sdk.Context, keeper Keeper, msg MsgCreatePool) (*sdk.Result, error) {
	// Verify min threshold
	MinThreshold := sdk.NewUint(uint64(keeper.GetParams(ctx).MinCreatePoolThreshold))
	if msg.NativeAssetAmount.LT(MinThreshold) { // Need to verify
		return nil, types.ErrTotalAmountTooLow
	}
	// Check if pool already exists
	if keeper.ExistsPool(ctx, msg.ExternalAsset.Ticker) {
		return nil, types.ErrUnableToCreatePool
	}

	asset := msg.ExternalAsset
	// Verify user has balance to create a new pool
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(msg.ExternalAssetAmount.Uint64()))
	nativeAssetCoin := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(msg.NativeAssetAmount.Uint64()))
	if !keeper.BankKeeper.HasCoins(ctx, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}

	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	poolUnits, lpunits, err := calculatePoolUnits(sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), nativeBalance, externalBalance)
	if err != nil {
		return nil, err
	}
	pool, err := NewPool(asset, nativeBalance, externalBalance, poolUnits)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}
	// Send coins from suer to pool
	err = keeper.BankKeeper.SendCoins(ctx, msg.Signer, pool.PoolAddress, sdk.Coins{externalAssetCoin, nativeAssetCoin})
	if err != nil {
		return nil, err
	}
	// Pool creator becomes the first LP
	lp := NewLiquidityProvider(asset, lpunits, msg.Signer)
	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	keeper.SetLiquidityProvider(ctx, lp)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyPool, pool.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			types.EventTypeCreateLiquidityProvider,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgAddLiquidity(ctx sdk.Context, keeper Keeper, msg MsgAddLiquidity) (*sdk.Result, error) {
	createNewLP := false
	// Get pool
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	newPoolUnits, lpUnits, err := calculatePoolUnits(pool.PoolUnits, pool.NativeAssetBalance, pool.ExternalAssetBalance, msg.NativeAssetAmount, msg.ExternalAssetAmount)
	if err != nil {
		return nil, err
	}
	// Get lp , if lp doesnt exist create lp
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Ticker, msg.Signer.String())
	if err != nil {
		createNewLP = true
	}
	// Verify user has coins to add liquidity
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(msg.ExternalAssetAmount.Uint64()))
	nativeAssetCoin := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(msg.NativeAssetAmount.Uint64()))
	if !keeper.BankKeeper.HasCoins(ctx, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}
	// Send from user to pool
	err = keeper.BankKeeper.SendCoins(ctx, msg.Signer, pool.PoolAddress, sdk.Coins{externalAssetCoin, nativeAssetCoin})
	if err != nil {
		return nil, err
	}

	pool.PoolUnits = newPoolUnits
	pool.NativeAssetBalance = pool.NativeAssetBalance.Add(msg.NativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Add(msg.ExternalAssetAmount)
	// Create lp if needed
	// Doesn't look like this can occur, as creating a pool creates this as well.  Not sure if this is a valid scenario
	if createNewLP {
		lp := NewLiquidityProvider(msg.ExternalAsset, lpUnits, msg.Signer)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeCreateLiquidityProvider,
				sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
		})
	} else {
		lp.LiquidityProviderUnits = lp.LiquidityProviderUnits.Add(lpUnits)
	}
	// Save new pool balances
	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Save LP
	keeper.SetLiquidityProvider(ctx, lp)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveLiquidity(ctx sdk.Context, keeper Keeper, msg MsgRemoveLiquidity) (*sdk.Result, error) {
	// Get pool
	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Ticker)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	//Get LP
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Ticker, msg.Signer.String())
	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}

	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WBasisPoints.String(), msg.Asymmetry)

	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(withdrawExternalAssetAmount.Uint64()))
	nativeAssetCoin := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(withdrawNativeAssetAmount.Uint64()))

	// Subtract Value from pool
	pool.PoolUnits = pool.PoolUnits.Sub(lp.LiquidityProviderUnits).Add(lpUnitsLeft)
	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	// Check if withdrawal makes pool too shallow , checking only for asymetric withdraw.
	if !msg.Asymmetry.IsZero() && (pool.ExternalAssetBalance.IsZero() || pool.NativeAssetBalance.IsZero()) {
		return nil, errors.Wrap(types.ErrPoolTooShallow, "Pool Balance nil before adjusting asymmetry")
	}

	// Swapping between Native and External based on Asymmetry
	if msg.Asymmetry.IsPositive() {
		swapResult, _, _, swappedPool, err := SwapOne(GetSettlementAsset(), swapAmount, msg.ExternalAsset, pool)
		if err != nil {
			return nil, types.ErrSwapping
		}
		if !swapResult.IsZero() {
			swapCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(swapResult.Uint64()))
			swapAmountInCoin := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(swapAmount.Uint64()))
			externalAssetCoin = externalAssetCoin.Add(swapCoin)
			nativeAssetCoin = nativeAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	if msg.Asymmetry.IsNegative() {
		swapResult, _, _, swappedPool, err := SwapOne(msg.ExternalAsset, swapAmount, GetSettlementAsset(), pool)
		if err != nil {
			return nil, types.ErrSwapping
		}
		if !swapResult.IsZero() {
			swapCoin := sdk.NewCoin(GetSettlementAsset().Ticker, sdk.NewIntFromUint64(swapResult.Uint64()))
			swapAmountInCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(swapAmount.Uint64()))

			nativeAssetCoin = nativeAssetCoin.Add(swapCoin)
			externalAssetCoin = externalAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	//Calculate final withdraw amount after swap
	sendCoins := sdk.Coins{}
	if !externalAssetCoin.IsZero() && !externalAssetCoin.IsNegative() {
		sendCoins = sendCoins.Add(externalAssetCoin)
	}

	if !nativeAssetCoin.IsZero() && !nativeAssetCoin.IsNegative() {
		sendCoins = sendCoins.Add(nativeAssetCoin)
	}
	// Verify if Swap makes the pool too shallow in one of the assets
	if externalAssetCoin.Amount.GTE(sdk.Int(poolOriginalEB)) || nativeAssetCoin.Amount.GTE(sdk.Int(poolOriginalNB)) {
		return nil, errors.Wrap(types.ErrPoolTooShallow, "Pool Balance nil after adjusting asymmetry")
	}
	// Setting pool after all calculations of withdraw and then swap
	err = keeper.SetPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Send coins from pool to user
	if !sendCoins.Empty() {
		if !keeper.BankKeeper.HasCoins(ctx, pool.PoolAddress, sendCoins) {
			return nil, types.ErrNotEnoughLiquidity
		}
		err = keeper.BankKeeper.SendCoins(ctx, pool.PoolAddress, msg.Signer, sendCoins)
		if err != nil {
			return nil, err
		}
	}

	if lpUnitsLeft.IsZero() {
		keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress.String())
	} else {
		lp.LiquidityProviderUnits = lpUnitsLeft
		keeper.SetLiquidityProvider(ctx, lp)
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSwap(ctx sdk.Context, keeper Keeper, msg MsgSwap) (*sdk.Result, error) {
	var (
		liquidityFee sdk.Uint
		tradeSlip    sdk.Uint
	)
	liquidityFee = sdk.ZeroUint()
	tradeSlip = sdk.ZeroUint()
	sentAmount := msg.SentAmount
	sentAsset := msg.SentAsset
	receivedAsset := msg.ReceivedAsset
	// Get native asset
	nativeAsset := types.GetSettlementAsset()
	// If its one swap , this pool would be RWN:RWN ( Ex User sends RWN wants ETH)
	// If its two swap . this pool would be RWN:EXTERNAL1 ( Ex User sends ETH wants XCT , ETH is EXTERNAL1)
	//CASE 1 : RWN:ETH
	//CASE 2 : RWN:ETH
	inPool, err := keeper.GetPool(ctx, msg.SentAsset.Ticker)
	if err != nil {
		return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
	}
	// If its one swap , this pool would be RWN:EXTERNAL ( Ex User sends RWN wants ETH , ETH IS EXTERNAL )
	// If its two swap . this pool would be RWN:EXTERNAL2 ( Ex User sends ETH wants XCT , XCT is EXTERNAL2)
	//CASE 1 : RWN:ETH
	//CASE 2 : RWN:XCT
	outPool, err := keeper.GetPool(ctx, msg.ReceivedAsset.Ticker)
	if err != nil {
		return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.ReceivedAsset.String())
	}

	// Deducting Balance from the user , Sent Asset is the asset the user is sending to the Pool
	// Case 1 . Deducting his RWN and adding to RWN:ETH pool
	// Case 2 , Deduction his ETH and adding to RWN:ETH pool
	sentCoin := sdk.NewCoin(msg.SentAsset.Ticker, sdk.NewIntFromUint64(sentAmount.Uint64()))
	if !keeper.BankKeeper.HasCoins(ctx, msg.Signer, sdk.Coins{sentCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}
	err = keeper.BankKeeper.SendCoins(ctx, msg.Signer, inPool.PoolAddress, sdk.Coins{sentCoin})
	if err != nil {
		return nil, err
	}
	// Check if its a two way swap, swapping non native fro non native .
	// If its one way we can skip this if condition and add balance to users account from outpool
	if msg.SentAsset != nativeAsset && msg.ReceivedAsset != nativeAsset {

		emitAmount, lp, ts, finalPool, err := SwapOne(sentAsset, sentAmount, nativeAsset, inPool)
		if err != nil {
			return nil, err
		}
		err = keeper.SetPool(ctx, finalPool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
		sentAmount = emitAmount
		sentAsset = nativeAsset
		liquidityFee = liquidityFee.Add(lp)
		tradeSlip = tradeSlip.Add(ts)
		interpoolCoin := sdk.NewCoin(nativeAsset.Ticker, sdk.NewIntFromUint64(emitAmount.Uint64()))
		// Case 2 - Transfer from RWN:ETH -> RWN:DASH
		err = keeper.BankKeeper.SendCoins(ctx, outPool.PoolAddress, inPool.PoolAddress, sdk.Coins{interpoolCoin})
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToAddBalance, err.Error())
		}
	}
	// Calculating amount user receives
	emitAmount, lp, ts, finalPool, err := SwapOne(sentAsset, sentAmount, receivedAsset, outPool)
	if err != nil {
		return nil, err
	}
	err = keeper.SetPool(ctx, finalPool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	// Adding balance to users account ,Received Asset is the asset the user wants to receive
	// Case 1 . Adding his ETH and deducting from  RWN:ETH pool
	// Case 2 , Adding his XCT and deducting from  RWN:XCT pool
	sentCoin = sdk.NewCoin(msg.ReceivedAsset.Ticker, sdk.NewIntFromUint64(sentAmount.Uint64()))
	err = keeper.BankKeeper.SendCoins(ctx, outPool.PoolAddress, msg.Signer, sdk.Coins{sentCoin})
	if err != nil {
		return nil, err
	}
	liquidityFee = liquidityFee.Add(lp)
	tradeSlip = tradeSlip.Add(ts)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
			sdk.NewAttribute(types.AttributeKeyLiquidityFee, liquidityFee.String()),
			sdk.NewAttribute(types.AttributeKeyTradeSlip, tradeSlip.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

//------------------------------------------------------------------------------------------------------------------
// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
func SwapOne(from Asset, sentAmount sdk.Uint, to Asset, pool Pool) (sdk.Uint, sdk.Uint, sdk.Uint, Pool, error) {

	var X sdk.Uint
	var Y sdk.Uint

	if to == GetSettlementAsset() {
		Y = pool.NativeAssetBalance
		X = pool.ExternalAssetBalance
	} else {
		X = pool.NativeAssetBalance
		Y = pool.ExternalAssetBalance
	}
	x := sentAmount
	liquidityFee := calcLiquidityFee(X, x, Y)
	tradeSlip := calcTradeSlip(X, x)
	swapResult := calcSwapResult(X, x, Y)
	if swapResult.GTE(Y) {
		return sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), Pool{}, types.ErrNotEnoughAssetTokens
	}
	if from == GetSettlementAsset() {
		pool.NativeAssetBalance = X.Add(x)
		pool.ExternalAssetBalance = Y.Sub(swapResult)
	} else {
		pool.ExternalAssetBalance = X.Add(x)
		pool.NativeAssetBalance = Y.Sub(swapResult)
	}

	return swapResult, liquidityFee, tradeSlip, pool, nil
}

// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
func CalculateWithdrawal(poolUnits sdk.Uint, nativeAssetBalance string,
	externalAssetBalance string, lpUnits string, wBasisPoints string, asymmetry sdk.Int) (sdk.Uint, sdk.Uint, sdk.Uint, sdk.Uint) {
	poolUnitsF := sdk.NewDecFromBigInt(poolUnits.BigInt())

	nativeAssetBalanceF, err := sdk.NewDecFromStr(nativeAssetBalance)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", nativeAssetBalance, err))
	}
	externalAssetBalanceF, err := sdk.NewDecFromStr(externalAssetBalance)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", externalAssetBalance, err))
	}
	lpUnitsF, err := sdk.NewDecFromStr(lpUnits)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", lpUnits, err))
	}
	wBasisPointsF, err := sdk.NewDecFromStr(wBasisPoints)
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", wBasisPoints, err))
	}
	asymmetryF, err := sdk.NewDecFromStr(asymmetry.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", asymmetry.String(), err))
	}
	denominator := sdk.NewDec(10000).Quo(wBasisPointsF)
	unitsToClaim := lpUnitsF.Quo(denominator)
	withdrawExternalAssetAmount := externalAssetBalanceF.Quo(poolUnitsF.Quo(unitsToClaim))
	withdrawNativeAssetAmount := nativeAssetBalanceF.Quo(poolUnitsF.Quo(unitsToClaim))

	swapAmount := sdk.NewDec(0)
	//if asymmetry is positive we need to swap from native to external
	if asymmetry.IsPositive() {
		unitsToSwap := unitsToClaim.Quo(sdk.NewDec(10000).Quo(asymmetryF.Abs()))
		swapAmount = nativeAssetBalanceF.Quo(poolUnitsF.Quo(unitsToSwap))
	}
	//if asymmetry is negative we need to swap from external to native
	if asymmetry.IsNegative() {
		unitsToSwap := unitsToClaim.Quo(sdk.NewDec(10000).Quo(asymmetryF.Abs()))
		swapAmount = externalAssetBalanceF.Quo(poolUnitsF.Quo(unitsToSwap))
	}
	//if asymmetry is 0 we don't need to swap

	lpUnitsLeft := lpUnitsF.Sub(unitsToClaim)
	return sdk.NewUintFromBigInt(withdrawNativeAssetAmount.RoundInt().BigInt()),
		sdk.NewUintFromBigInt(withdrawExternalAssetAmount.RoundInt().BigInt()),
		sdk.NewUintFromBigInt(lpUnitsLeft.RoundInt().BigInt()),
		sdk.NewUintFromBigInt(swapAmount.RoundInt().BigInt())
}

// More details on the formula
// https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md

//native asset balance  : currently in pool before adding
//external asset balance : currently in pool before adding
//native asset to added  : the amount the user sends
//external asset amount to be added : the amount the user sends

// r = native asset added;
// a = external asset added
// R = native Balance (before)
// A = external Balance (before)
// P = existing Pool Units
// slipAdjustment = (1 - ABS((R a - r A)/((2 r + R) (a + A))))
// units = ((P (a R + A r))/(2 A R))*slidAdjustment

func calculatePoolUnits(oldPoolUnits, nativeAssetBalance, externalAssetBalance,
	nativeAssetAmount, externalAssetAmount sdk.Uint) (sdk.Uint, sdk.Uint, error) {
	if nativeAssetBalance.Add(nativeAssetAmount).IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), errors.New("total Native in the pool is zero")
	}
	if externalAssetBalance.Add(externalAssetAmount).IsZero() {
		return sdk.ZeroUint(), sdk.ZeroUint(), errors.New("total External in the pool is zero")
	}
	if nativeAssetBalance.IsZero() || externalAssetBalance.IsZero() {
		return nativeAssetAmount, nativeAssetAmount, nil
	}
	P, err := sdk.NewDecFromStr(oldPoolUnits.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", oldPoolUnits.String(), err))
	}
	R, err := sdk.NewDecFromStr(nativeAssetBalance.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", nativeAssetBalance.String(), err))
	}
	A, err := sdk.NewDecFromStr(externalAssetBalance.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", externalAssetBalance.String(), err))
	}
	r, err := sdk.NewDecFromStr(nativeAssetAmount.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", nativeAssetAmount.String(), err))
	}
	a, err := sdk.NewDecFromStr(externalAssetAmount.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", externalAssetAmount.String(), err))
	}

	// (2 r + R) (a + A)
	// (2 r + R) (a + A)
	slipAdjDenominator := (r.MulInt64(2).Add(R)).Mul(a.Add(A))
	// ABS((R a - r A)/((2 r + R) (a + A)))
	var slipAdjustment sdk.Dec
	if R.Mul(a).GT(r.Mul(A)) {
		slipAdjustment = R.Mul(a).Sub(r.Mul(A)).Quo(slipAdjDenominator)
	} else {
		slipAdjustment = r.Mul(A).Sub(R.Mul(a)).Quo(slipAdjDenominator)
	}
	// (1 - ABS((R a - r A)/((2 r + R) (a + A))))
	slipAdjustment = sdk.NewDec(1).Sub(slipAdjustment)

	// ((P (a R + A r))
	numerator := P.Mul(a.Mul(R).Add(A.Mul(r)))
	// 2AR
	denominator := sdk.NewDec(2).Mul(A).Mul(R)
	stakeUnits := numerator.Quo(denominator).Mul(slipAdjustment)
	newPoolUnit := P.Add(stakeUnits)

	return sdk.NewUintFromBigInt(newPoolUnit.RoundInt().BigInt()), sdk.NewUintFromBigInt(stakeUnits.RoundInt().BigInt()), nil

}

func calcLiquidityFee(X, x, Y sdk.Uint) sdk.Uint {
	d := x.Add(X)
	denom := d.Mul(d)
	return (x.Mul(x).Mul(Y)).Quo(denom)
}

func calcTradeSlip(X, x sdk.Uint) sdk.Uint {
	numerator := x.Mul(sdk.NewUint(2).Mul(X).Add(x))
	denom := X.Mul(X)
	return numerator.Quo(denom)
}

func calcSwapResult(X, x, Y sdk.Uint) sdk.Uint {
	d := x.Add(X)
	denom := d.Mul(d)
	return (x.Mul(X).Mul(Y)).Quo(denom)
}
