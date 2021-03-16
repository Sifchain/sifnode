package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/clp/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the clp MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) DecommissionPool(goCtx context.Context, msg *types.MsgDecommissionPool) (*types.MsgDecommissionPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, err := k.Keeper.GetPool(ctx, msg.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}

	addAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	if !k.Keeper.ValidateAddress(ctx, addAddr) {
		return nil, errors.Wrap(types.ErrInvalid, "user does not have permission to decommission pool")
	}
	if pool.NativeAssetBalance.GTE(sdk.NewUintFromString(types.PoolThrehold)) {
		return nil, types.ErrBalanceTooHigh
	}
	// Get all LP's for the pool
	lpList := k.Keeper.GetLiquidityProvidersForAsset(ctx, pool.ExternalAsset)
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	// iterate over Lp list and refund them there tokens
	// Return both RWN and EXTERNAL ASSET
	for _, lp := range lpList {
		withdrawNativeAsset, withdrawExternalAsset, _, _ := clpkeeper.CalculateAllAssetsForLP(pool, lp)
		poolUnits = poolUnits.Sub(lp.LiquidityProviderUnits)
		nativeAssetBalance = nativeAssetBalance.Sub(withdrawNativeAsset)
		externalAssetBalance = externalAssetBalance.Sub(withdrawExternalAsset)

		withdrawNativeAssetInt, ok := k.Keeper.ParseToInt(withdrawNativeAsset.String())
		if !ok {
			return nil, types.ErrUnableToParseInt
		}
		withdrawExternalAssetInt, ok := k.Keeper.ParseToInt(withdrawExternalAsset.String())
		if !ok {
			return nil, types.ErrUnableToParseInt
		}
		withdrawNativeCoins := sdk.NewCoin(types.GetSettlementAsset().Symbol, withdrawNativeAssetInt)
		withdrawExternalCoins := sdk.NewCoin(msg.Symbol, withdrawExternalAssetInt)
		refundingCoins := sdk.Coins{withdrawExternalCoins, withdrawNativeCoins}
		err := k.Keeper.RemoveLiquidityProvider(ctx, refundingCoins, lp)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnableToRemoveLiquidityProvider, err.Error())
		}
	}
	// Pool should be empty at this point
	// Decommission the pool
	err = k.Keeper.DecommissionPool(ctx, pool)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToDecommissionPool, err.Error())
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})

	return &types.MsgDecommissionPoolResponse{}, nil
}

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 	var (
	// 		priceImpact sdk.Uint
	// 	)

	// 	liquidityFeeNative := sdk.ZeroUint()
	// 	liquidityFeeExternal := sdk.ZeroUint()
	// 	totalLiquidityFee := sdk.ZeroUint()
	// 	priceImpact = sdk.ZeroUint()
	// 	sentAmount := msg.SentAmount

	// 	sentAsset := msg.SentAsset
	// 	receivedAsset := msg.ReceivedAsset
	// 	// Get native asset
	// 	nativeAsset := types.GetSettlementAsset()

	// 	inPool, outPool := types.Pool{}, types.Pool{}
	// 	err := errors.New("Swap Error")
	// 	// If sending rowan ,deduct directly from the Native balance  instead of fetching from rowan pool
	// 	if msg.SentAsset != types.GetSettlementAsset() {
	// 		inPool, err = keeper.GetPool(ctx, msg.SentAsset.Symbol)
	// 		if err != nil {
	// 			return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
	// 		}
	// 	}
	// 	fmt.Println(err)
	// 	sentAmountInt, ok := keeper.ParseToInt(sentAmount.String())
	// 	if !ok {
	// 		return nil, types.ErrUnableToParseInt
	// 	}
	// 	sentCoin := sdk.NewCoin(msg.SentAsset.Symbol, sentAmountInt)
	// 	err = keeper.InitiateSwap(ctx, sentCoin, msg.Signer)
	// 	if err != nil {
	// 		return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
	// 	}
	// 	// Check if its a two way swap, swapping non native fro non native .
	// 	// If its one way we can skip this if condition and add balance to users account from outpool

	// 	if msg.SentAsset != nativeAsset && msg.ReceivedAsset != nativeAsset {
	// 		emitAmount, lp, ts, finalPool, err := clpkeeper.SwapOne(sentAsset, sentAmount, nativeAsset, inPool)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		err = keeper.SetPool(ctx, finalPool)
	// 		if err != nil {
	// 			return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	// 		}
	// 		sentAmount = emitAmount
	// 		sentAsset = nativeAsset
	// 		priceImpact = priceImpact.Add(ts)
	// 		liquidityFeeNative = liquidityFeeNative.Add(lp)
	// 	}
	// 	// If receiving  rowan , add directly to  Native balance  instead of fetching from rowan pool
	// 	if msg.ReceivedAsset == types.GetSettlementAsset() {
	// 		outPool, err = keeper.GetPool(ctx, msg.SentAsset.Symbol)
	// 		if err != nil {
	// 			return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
	// 		}
	// 	} else {
	// 		outPool, err = keeper.GetPool(ctx, msg.ReceivedAsset.Symbol)
	// 		if err != nil {
	// 			return nil, errors.Wrap(types.ErrPoolDoesNotExist, msg.ReceivedAsset.String())
	// 		}
	// 	}

	// 	// Calculating amount user receives
	// 	emitAmount, lp, ts, finalPool, err := clpkeeper.SwapOne(sentAsset, sentAmount, receivedAsset, outPool)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if emitAmount.LT(msg.MinReceivingAmount) {
	// 		ctx.EventManager().EmitEvents(sdk.Events{
	// 			sdk.NewEvent(
	// 				types.EventTypeSwapFailed,
	// 				sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
	// 				sdk.NewAttribute(types.AttributeKeyThreshold, msg.MinReceivingAmount.String()),
	// 				sdk.NewAttribute(types.AttributeKeyInPool, inPool.String()),
	// 				sdk.NewAttribute(types.AttributeKeyOutPool, outPool.String()),
	// 				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	// 			),
	// 			sdk.NewEvent(
	// 				sdk.EventTypeMessage,
	// 				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
	// 			),
	// 		})
	// 		return &sdk.Result{Events: ctx.EventManager().Events()}, types.ErrReceivedAmountBelowExpected

	return &types.MsgSwapResponse{}, nil
}
func (k msgServer) RemoveLiquidity(goCtx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 	// Get pool
	// 	pool, err := keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	// 	if err != nil {
	// 		return nil, types.ErrPoolDoesNotExist
	// 	}
	// 	//Get LP
	// 	lp, err := keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer.String())
	// 	if err != nil {
	// 		return nil, types.ErrLiquidityProviderDoesNotExist
	// 	}
	// 	poolOriginalEB := pool.ExternalAssetBalance
	// 	poolOriginalNB := pool.NativeAssetBalance
	// 	//Calculate amount to withdraw
	// 	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := clpkeeper.CalculateWithdrawal(pool.PoolUnits,
	// 		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
	// 		msg.WBasisPoints.String(), msg.Asymmetry)

	// 	withdrawExternalAssetAmountInt, ok := keeper.ParseToInt(withdrawExternalAssetAmount.String())
	// 	if !ok {
	// 		return nil, types.ErrUnableToParseInt
	// 	}
	// 	withdrawNativeAssetAmountInt, ok := keeper.ParseToInt(withdrawNativeAssetAmount.String())
	// 	if !ok {
	// 		return nil, types.ErrUnableToParseInt
	// 	}
	// 	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, withdrawExternalAssetAmountInt)
	// 	nativeAssetCoin := sdk.NewCoin(GetSettlementAsset().Symbol, withdrawNativeAssetAmountInt)

	// 	// Subtract Value from pool
	// 	pool.PoolUnits = pool.PoolUnits.Sub(lp.LiquidityProviderUnits).Add(lpUnitsLeft)
	// 	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	// 	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	// 	// Check if withdrawal makes pool too shallow , checking only for asymetric withdraw.

	// 	if !msg.Asymmetry.IsZero() && (pool.ExternalAssetBalance.IsZero() || pool.NativeAssetBalance.IsZero()) {
	// 		return nil, errors.Wrap(types.ErrPoolTooShallow, "pool balance nil before adjusting asymmetry")
	// 	}

	// 	// Swapping between Native and External based on Asymmetry
	// 	if msg.Asymmetry.IsPositive() {
	// 		swapResult, _, _, swappedPool, err := clpkeeper.SwapOne(GetSettlementAsset(), swapAmount, msg.ExternalAsset, pool)
	// 		if err != nil {
	// 			return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
	// 		}
	// 		if !swapResult.IsZero() {
	// 			swapResultInt, ok := keeper.ParseToInt(swapResult.String())
	// 			if !ok {
	// 				return nil, types.ErrUnableToParseInt
	// 			}
	// 			swapAmountInt, ok := keeper.ParseToInt(swapAmount.String())
	// 			if !ok {
	// 				return nil, types.ErrUnableToParseInt
	// 			}
	// 			swapCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapResultInt)
	// 			swapAmountInCoin := sdk.NewCoin(GetSettlementAsset().Symbol, swapAmountInt)
	// 			externalAssetCoin = externalAssetCoin.Add(swapCoin)
	// 			nativeAssetCoin = nativeAssetCoin.Sub(swapAmountInCoin)
	// 		}
	// 		pool = swappedPool
	// 	}
	// 	if msg.Asymmetry.IsNegative() {
	// 		swapResult, _, _, swappedPool, err := clpkeeper.SwapOne(msg.ExternalAsset, swapAmount, GetSettlementAsset(), pool)
	// 		if err != nil {
	// 			return nil, errors.Wrap(types.ErrUnableToSwap, err.Error())
	// 		}
	// 		if !swapResult.IsZero() {
	// 			swapInt, ok := keeper.ParseToInt(swapResult.String())
	// 			if !ok {
	// 				return nil, types.ErrUnableToParseInt
	// 			}
	// 			swapAmountInt, ok := keeper.ParseToInt(swapAmount.String())
	// 			if !ok {
	// 				return nil, types.ErrUnableToParseInt
	// 			}
	// 			swapCoin := sdk.NewCoin(GetSettlementAsset().Symbol, swapInt)
	// 			swapAmountInCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapAmountInt)

	// 			nativeAssetCoin = nativeAssetCoin.Add(swapCoin)
	// 			externalAssetCoin = externalAssetCoin.Sub(swapAmountInCoin)
	// 		}
	// 		pool = swappedPool
	// 	}
	// 	// Check and  remove Liquidity
	// 	err = keeper.RemoveLiquidity(ctx, pool, externalAssetCoin, nativeAssetCoin, lp, lpUnitsLeft, poolOriginalEB, poolOriginalNB)
	// 	if err != nil {
	// 		return nil, errors.Wrap(types.ErrUnableToRemoveLiquidity, err.Error())
	// 	}
	// 	ctx.EventManager().EmitEvents(sdk.Events{
	// 		sdk.NewEvent(
	// 			types.EventTypeRemoveLiquidity,
	// 			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
	// 			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
	// 		),
	// 		sdk.NewEvent(
	// 			sdk.EventTypeMessage,
	// 			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
	// 		),
	// 	})

	return &types.MsgRemoveLiquidityResponse{}, nil
}

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify min threshold

	MinThreshold := sdk.NewUintFromString(types.PoolThrehold)

	if msg.NativeAssetAmount.LT(MinThreshold) { // Need to verify
		return nil, types.ErrTotalAmountTooLow
	}
	// Check if pool already exists
	if k.Keeper.ExistsPool(ctx, msg.ExternalAsset.Symbol) {
		return nil, types.ErrUnableToCreatePool
	}

	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	poolUnits, lpunits, err := CalculatePoolUnits(msg.ExternalAsset.Symbol, sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(), nativeBalance, externalBalance)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}
	// Create Pool
	pool, err := k.Keeper.CreatePool(ctx, poolUnits, msg)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToSetPool, err.Error())
	}

	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	lp := k.Keeper.CreateLiquidityProvider(ctx, msg.ExternalAsset, lpunits, accAddr)

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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})

	return &types.MsgCreatePoolResponse{}, nil
}

func (k msgServer) AddLiquidity(goCtx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get pool
	pool, err := k.Keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}

	newPoolUnits, lpUnits, err := CalculatePoolUnits(
		msg.ExternalAsset.Symbol,
		pool.PoolUnits,
		pool.NativeAssetBalance,
		pool.ExternalAssetBalance,
		msg.NativeAssetAmount,
		msg.ExternalAssetAmount)
	if err != nil {
		return nil, err
	}

	// Get lp , if lp doesnt exist create lp

	lp, err := k.Keeper.AddLiquidity(ctx, msg, pool, newPoolUnits, lpUnits)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnableToAddLiquidity, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})

	return &types.MsgAddLiquidityResponse{}, nil
}
