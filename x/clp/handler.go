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
	if pool.ExternalAssetBalance+pool.NativeAssetBalance > keeper.GetParams(ctx).MinCreatePoolThreshold {
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
		withdrawNativeAsset, withdrawExternalAsset, _, _ := calculateWithdrawal(poolUnits, nativeAssetBalance, externalAssetBalance, lp.LiquidityProviderUnits, 10000, 0)
		poolUnits = poolUnits - lp.LiquidityProviderUnits
		nativeAssetBalance = nativeAssetBalance - withdrawNativeAsset
		externalAssetBalance = externalAssetBalance - withdrawExternalAsset
		withdrawNativeCoins := sdk.NewCoin(GetNativeAsset().Ticker, sdk.NewIntFromUint64(uint64(withdrawNativeAsset)))
		withdrawExternalCoins := sdk.NewCoin(msg.Ticker, sdk.NewIntFromUint64(uint64(withdrawExternalAsset)))
		lpAddess, err := sdk.AccAddressFromBech32(lp.LiquidityProviderAddress)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToDestroyPool, err.Error())
		}
		err = keeper.BankKeeper.SendCoins(ctx, pool.PoolAddress, lpAddess, sdk.Coins{withdrawNativeCoins})
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToAddBalance, err.Error())
		}
		err = keeper.BankKeeper.SendCoins(ctx, pool.PoolAddress, lpAddess, sdk.Coins{withdrawExternalCoins})
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToAddBalance, err.Error())
		}
		keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress)
	}
	// TODO : Do we check if nativeBalance and external balance is still left in the pool before we delete it ?
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
	MinThreshold := keeper.GetParams(ctx).MinCreatePoolThreshold
	if (msg.ExternalAssetAmount + msg.NativeAssetAmount) < MinThreshold { // Need to verify
		return nil, types.ErrTotalAmountTooLow
	}
	// Check if pool already exists
	if keeper.ExistsPool(ctx, msg.ExternalAsset.Ticker) {
		return nil, types.ErrUnableToCreatePool
	}

	asset := msg.ExternalAsset
	// Verify user has balance to create a new pool
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(uint64(msg.ExternalAssetAmount)))
	nativeAssetCoin := sdk.NewCoin(GetNativeAsset().Ticker, sdk.NewIntFromUint64(uint64(msg.NativeAssetAmount)))
	if !keeper.BankKeeper.HasCoins(ctx, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}

	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	poolUnits, lpunits := calculatePoolUnits(0, 0, 0, nativeBalance, externalBalance)
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
	lp := NewLiquidityProvider(asset, lpunits, msg.Signer.String())
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
	newPoolUnits, lpUnits := calculatePoolUnits(pool.PoolUnits, pool.NativeAssetBalance, pool.ExternalAssetBalance, msg.NativeAssetAmount, msg.ExternalAssetAmount)
	// Get lp , if lp doesnt exist create lp
	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Ticker, msg.Signer.String())
	if err != nil {
		createNewLP = true
	}
	// Verify user has coins to add liquidity
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(uint64(msg.ExternalAssetAmount)))
	nativeAssetCoin := sdk.NewCoin(GetNativeAsset().Ticker, sdk.NewIntFromUint64(uint64(msg.NativeAssetAmount)))
	if !keeper.BankKeeper.HasCoins(ctx, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
		return nil, types.ErrBalanceNotAvailable
	}
	// Send from user to pool
	err = keeper.BankKeeper.SendCoins(ctx, msg.Signer, pool.PoolAddress, sdk.Coins{externalAssetCoin, nativeAssetCoin})
	if err != nil {
		return nil, err
	}

	pool.PoolUnits = newPoolUnits
	pool.NativeAssetBalance = pool.NativeAssetBalance + msg.NativeAssetAmount
	pool.ExternalAssetBalance = pool.ExternalAssetBalance + msg.ExternalAssetAmount
	// Create lp if needed
	if createNewLP {
		lp := NewLiquidityProvider(msg.ExternalAsset, lpUnits, msg.Signer.String())
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeCreateLiquidityProvider,
				sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
		})
	} else {
		lp.LiquidityProviderUnits = lp.LiquidityProviderUnits + lpUnits
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

	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := calculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance, pool.ExternalAssetBalance, lp.LiquidityProviderUnits,
		msg.WBasisPoints, msg.Asymmetry)

	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(uint64(withdrawExternalAssetAmount)))
	nativeAssetCoin := sdk.NewCoin(GetNativeAsset().Ticker, sdk.NewIntFromUint64(uint64(withdrawNativeAssetAmount)))
	// Send coins from pool to user
	pool.PoolUnits = pool.PoolUnits - lp.LiquidityProviderUnits + lpUnitsLeft
	pool.NativeAssetBalance = pool.NativeAssetBalance - withdrawNativeAssetAmount
	pool.ExternalAssetBalance = pool.ExternalAssetBalance - withdrawExternalAssetAmount
	if msg.Asymmetry > 0 {
		swapResult, _, _, swappedPool, err := swapOne(GetNativeAsset(), swapAmount, msg.ExternalAsset, pool)
		if err != nil {
			return nil, types.ErrSwapping
		}
		if swapResult != 0 {
			swapCoin := sdk.NewCoin(msg.ExternalAsset.Ticker, sdk.NewIntFromUint64(uint64(swapResult)))
			externalAssetCoin.Add(swapCoin)
		}
		err = keeper.SetPool(ctx, swappedPool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
	}
	//if asymmetry is negative we need to swap from external to native
	if msg.Asymmetry < 0 {
		swapResult, _, _, swappedPool, err := swapOne(msg.ExternalAsset, swapAmount, GetNativeAsset(), pool)
		if err != nil {
			return nil, types.ErrSwapping
		}
		if swapResult != 0 {
			swapCoin := sdk.NewCoin(GetNativeAsset().Ticker, sdk.NewIntFromUint64(uint64(swapResult)))
			nativeAssetCoin.Add(swapCoin)
		}
		err = keeper.SetPool(ctx, swappedPool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
	}
	// if asymmetry is 0 , just set pool
	if msg.Asymmetry == 0 {
		err = keeper.SetPool(ctx, pool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
	}
	sendCoins := sdk.Coins{}
	if !externalAssetCoin.IsZero() && !externalAssetCoin.IsNegative() {
		sendCoins = sendCoins.Add(externalAssetCoin)

	}
	if !nativeAssetCoin.IsZero() && !nativeAssetCoin.IsNegative() {
		sendCoins = sendCoins.Add(nativeAssetCoin)
	}
	if !sendCoins.Empty() {
		if !keeper.BankKeeper.HasCoins(ctx, pool.PoolAddress, sdk.Coins{externalAssetCoin, nativeAssetCoin}) {
			return nil, types.ErrNotEnoughLiquidity
		}
		err = keeper.BankKeeper.SendCoins(ctx, pool.PoolAddress, msg.Signer, sdk.Coins{externalAssetCoin, nativeAssetCoin})
		if err != nil {
			return nil, err
		}
	}

	if lpUnitsLeft == 0 {
		keeper.DestroyLiquidityProvider(ctx, lp.Asset.Ticker, lp.LiquidityProviderAddress)
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
		liquidityFee uint
		tradeSlip    uint
	)
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
	sentCoin := sdk.NewCoin(msg.SentAsset.Ticker, sdk.NewIntFromUint64(uint64(sentAmount)))
	err = keeper.BankKeeper.SendCoins(ctx, msg.Signer, inPool.PoolAddress, sdk.Coins{sentCoin})
	if err != nil {
		return nil, err
	}
	// Check if its a two way swap, swapping non native fro non native .
	// If its one way we can skip this if condition and add balance to users account from outpool
	if msg.SentAsset != nativeAsset && msg.ReceivedAsset != nativeAsset {

		emitAmount, lp, ts, finalPool, err := swapOne(sentAsset, sentAmount, nativeAsset, inPool)
		if err != nil {
			return nil, err
		}
		err = keeper.SetPool(ctx, finalPool)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
		sentAmount = emitAmount
		sentAsset = nativeAsset
		liquidityFee = liquidityFee + lp
		tradeSlip = tradeSlip + ts
		interpoolCoin := sdk.NewCoin(nativeAsset.Ticker, sdk.NewIntFromUint64(uint64(emitAmount)))
		// Case 2 - Transfer from RWN:ETH -> RWN:DASH
		err = keeper.BankKeeper.SendCoins(ctx, outPool.PoolAddress, inPool.PoolAddress, sdk.Coins{interpoolCoin})
	}
	// Calculating amount user receives
	emitAmount, lp, ts, finalPool, err := swapOne(sentAsset, sentAmount, receivedAsset, outPool)
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
	sentCoin = sdk.NewCoin(msg.ReceivedAsset.Ticker, sdk.NewIntFromUint64(uint64(sentAmount)))
	err = keeper.BankKeeper.SendCoins(ctx, outPool.PoolAddress, msg.Signer, sdk.Coins{sentCoin})
	if err != nil {
		return nil, err
	}
	liquidityFee = liquidityFee + lp
	tradeSlip = tradeSlip + ts
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeySwapAmount, strconv.FormatInt(int64(emitAmount), 10)),
			sdk.NewAttribute(types.AttributeKeyLiquidityFee, strconv.FormatInt(int64(liquidityFee), 10)),
			sdk.NewAttribute(types.AttributeKeyTradeSlip, strconv.FormatInt(int64(tradeSlip), 10)),
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

func swapOne(from Asset, sentAmount uint, to Asset, pool Pool) (uint, uint, uint, Pool, error) {

	var X uint
	var Y uint

	if to == GetNativeAsset() {
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
	if swapResult >= Y {
		return 0, 0, 0, Pool{}, types.ErrNotEnoughAssetTokens
	}
	if from == GetNativeAsset() {
		pool.NativeAssetBalance = X + x
		pool.ExternalAssetBalance = Y - swapResult
	} else {
		pool.ExternalAssetBalance = X + x
		pool.NativeAssetBalance = Y - swapResult
	}

	return swapResult, liquidityFee, tradeSlip, pool, nil
}

func calculateWithdrawal(poolUnits uint, nativeAssetBalance uint,
	externalAssetBalance uint, lpUnits uint, wBasisPoints int, asymmetry int) (uint, uint, uint, uint) {
	poolUnitsF := float64(poolUnits)
	nativeAssetBalanceF := float64(nativeAssetBalance)
	externalAssetBalanceF := float64(externalAssetBalance)
	lpUnitsF := float64(lpUnits)
	wBasisPointsF := float64(wBasisPoints)
	asymmetryF := float64(asymmetry)

	unitsToClaim := lpUnitsF / (10000 / (wBasisPointsF))
	withdrawExternalAssetAmount := externalAssetBalanceF / (poolUnitsF / unitsToClaim)
	withdrawNativeAssetAmount := nativeAssetBalanceF / (poolUnitsF / unitsToClaim)

	swapAmount := 0.0
	//if asymmetry is positive we need to swap from native to external
	if asymmetry > 0 {
		unitsToSwap := (unitsToClaim) / (10000 / (asymmetryF))
		swapAmount = (nativeAssetBalanceF) / (poolUnitsF / unitsToSwap)
	}
	//if asymmetry is negative we need to swap from external to native
	if asymmetry < 0 {
		unitsToSwap := (unitsToClaim) / (10000 / (-1 * asymmetryF))
		swapAmount = (externalAssetBalanceF) / (poolUnitsF / unitsToSwap)
	}
	//if asymmetry is 0 we don't need to swap

	lpUnitsLeft := lpUnitsF - unitsToClaim
	if withdrawNativeAssetAmount < 0 {
		withdrawNativeAssetAmount = 0
	}
	if withdrawExternalAssetAmount < 0 {
		withdrawExternalAssetAmount = 0
	}
	if lpUnitsLeft < 0 {
		lpUnitsLeft = 0
	}
	if swapAmount < 0 {
		swapAmount = 0
	}

	return uint(withdrawNativeAssetAmount), uint(withdrawExternalAssetAmount), uint(lpUnitsLeft), uint(swapAmount)
}

func calculatePoolUnits(oldPoolUnits uint, nativeAssetBalance uint, externalAssetBalance uint,
	nativeAssetAmount uint, externalAssetAmount uint) (uint, uint) {
	R := nativeAssetBalance + nativeAssetAmount
	A := externalAssetBalance + externalAssetAmount
	r := nativeAssetAmount
	a := externalAssetAmount
	lpUnits := ((R + A) * (r*A + R*a)) / (4 * R * A)
	poolUnits := oldPoolUnits + lpUnits
	return poolUnits, lpUnits
}

func calcLiquidityFee(X, x, Y uint) uint {
	return (x * x * Y) / ((x + X) * (x + X))
}

func calcTradeSlip(X, x uint) uint {
	return x * (2*X + x) / (X * X)
}

func calcSwapResult(X, x, Y uint) uint {
	return (x * X * Y) / ((x + X) * (x + X))
}
