package keeper

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

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
		return nil, sdkerrors.Wrap(types.ErrInvalid, "user does not have permission to decommission pool")
	}
	if pool.NativeAssetBalance.GTE(sdk.NewUintFromString(types.PoolThrehold)) {
		return nil, types.ErrBalanceTooHigh
	}
	// Get all LP's for the pool
	if pool.ExternalAsset == nil {
		return nil, errors.New("nill external asset")
	}
	req := query.PageRequest{
		Limit: uint64(math.MaxUint64),
	}
	lpList, _, err := k.Keeper.GetLiquidityProvidersForAssetPaginated(ctx, *pool.ExternalAsset, &req)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrLiquidityProviderDoesNotExist, err.Error())
	}
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	// iterate over Lp list and refund them there tokens
	// Return both RWN and EXTERNAL ASSET
	for _, lp := range lpList {
		withdrawNativeAsset, withdrawExternalAsset, _, _ := CalculateAllAssetsForLP(pool, *lp)
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
		refundingCoins := sdk.NewCoins(withdrawExternalCoins, withdrawNativeCoins)
		err := k.Keeper.RemoveLiquidityProvider(ctx, refundingCoins, *lp)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToRemoveLiquidityProvider, err.Error())
		}
	}
	// Pool should be empty at this point
	// Decommission the pool
	err = k.Keeper.DecommissionPool(ctx, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToDecommissionPool, err.Error())
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
	var (
		priceImpact sdk.Uint
	)

	liquidityFeeNative := sdk.ZeroUint()
	liquidityFeeExternal := sdk.ZeroUint()
	totalLiquidityFee := sdk.ZeroUint()
	priceImpact = sdk.ZeroUint()
	sentAmount := msg.SentAmount

	sentAsset := msg.SentAsset
	receivedAsset := msg.ReceivedAsset
	// Get native asset
	nativeAsset := types.GetSettlementAsset()

	inPool, outPool := types.Pool{}, types.Pool{}
	err := errors.New("Swap Error")
	// If sending rowan ,deduct directly from the Native balance  instead of fetching from rowan pool
	if !msg.SentAsset.Equals(types.GetSettlementAsset()) {
		inPool, err = k.Keeper.GetPool(ctx, msg.SentAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
	}
	fmt.Println(err)
	sentAmountInt, ok := k.Keeper.ParseToInt(sentAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}

	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	sentCoin := sdk.NewCoin(msg.SentAsset.Symbol, sentAmountInt)
	err = k.Keeper.InitiateSwap(ctx, sentCoin, accAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	// Check if its a two way swap, swapping non native fro non native .
	// If its one way we can skip this if condition and add balance to users account from outpool

	if !msg.SentAsset.Equals(nativeAsset) && !msg.ReceivedAsset.Equals(nativeAsset) {
		emitAmount, lp, ts, finalPool, err := SwapOne(*sentAsset, sentAmount, nativeAsset, inPool)
		if err != nil {
			return nil, err
		}
		err = k.Keeper.SetPool(ctx, &finalPool)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
		sentAmount = emitAmount
		sentAsset = &nativeAsset
		priceImpact = priceImpact.Add(ts)
		liquidityFeeNative = liquidityFeeNative.Add(lp)
	}
	// If receiving  rowan , add directly to  Native balance  instead of fetching from rowan pool
	if msg.ReceivedAsset.Equals(types.GetSettlementAsset()) {
		outPool, err = k.Keeper.GetPool(ctx, msg.SentAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
	} else {
		outPool, err = k.Keeper.GetPool(ctx, msg.ReceivedAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.ReceivedAsset.String())
		}
	}

	// Calculating amount user receives
	emitAmount, lp, ts, finalPool, err := SwapOne(*sentAsset, sentAmount, *receivedAsset, outPool)
	if err != nil {
		return nil, err
	}

	if emitAmount.LT(msg.MinReceivingAmount) {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeSwapFailed,
				sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
				sdk.NewAttribute(types.AttributeKeyThreshold, msg.MinReceivingAmount.String()),
				sdk.NewAttribute(types.AttributeKeyInPool, inPool.String()),
				sdk.NewAttribute(types.AttributeKeyOutPool, outPool.String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
		})
		return &types.MsgSwapResponse{}, types.ErrReceivedAmountBelowExpected
	}

	// todo nil pointer deref test
	err = k.Keeper.FinalizeSwap(ctx, emitAmount.String(), finalPool, *msg)

	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	if liquidityFeeNative.GT(sdk.ZeroUint()) {
		liquidityFeeExternal = liquidityFeeExternal.Add(lp)
		firstSwapFeeInOutputAsset := GetSwapFee(liquidityFeeNative, *msg.ReceivedAsset, outPool)
		totalLiquidityFee = liquidityFeeExternal.Add(firstSwapFeeInOutputAsset)
	} else {
		totalLiquidityFee = liquidityFeeNative.Add(lp)
	}

	priceImpact = priceImpact.Add(ts)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
			sdk.NewAttribute(types.AttributeKeyLiquidityFee, totalLiquidityFee.String()),
			sdk.NewAttribute(types.AttributeKeyPriceImpact, priceImpact.String()),
			sdk.NewAttribute(types.AttributeKeyInPool, inPool.String()),
			sdk.NewAttribute(types.AttributeKeyOutPool, outPool.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	return &types.MsgSwapResponse{}, nil
}

func (k msgServer) RemoveLiquidity(goCtx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, err := k.Keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	//Get LP
	lp, err := k.Keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer)
	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}
	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WBasisPoints.String(), msg.Asymmetry)

	withdrawExternalAssetAmountInt, ok := k.Keeper.ParseToInt(withdrawExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	withdrawNativeAssetAmountInt, ok := k.Keeper.ParseToInt(withdrawNativeAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, withdrawExternalAssetAmountInt)
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, withdrawNativeAssetAmountInt)

	// Subtract Value from pool
	pool.PoolUnits = pool.PoolUnits.Sub(lp.LiquidityProviderUnits).Add(lpUnitsLeft)
	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	// Check if withdrawal makes pool too shallow , checking only for asymetric withdraw.

	if !msg.Asymmetry.IsZero() && (pool.ExternalAssetBalance.IsZero() || pool.NativeAssetBalance.IsZero()) {
		return nil, sdkerrors.Wrap(types.ErrPoolTooShallow, "pool balance nil before adjusting asymmetry")
	}

	// Swapping between Native and External based on Asymmetry
	if msg.Asymmetry.IsPositive() {
		swapResult, _, _, swappedPool, err := SwapOne(types.GetSettlementAsset(), swapAmount, *msg.ExternalAsset, pool)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapResultInt, ok := k.Keeper.ParseToInt(swapResult.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapAmountInt, ok := k.Keeper.ParseToInt(swapAmount.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapResultInt)
			swapAmountInCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, swapAmountInt)
			externalAssetCoin = externalAssetCoin.Add(swapCoin)
			nativeAssetCoin = nativeAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	if msg.Asymmetry.IsNegative() {
		swapResult, _, _, swappedPool, err := SwapOne(*msg.ExternalAsset, swapAmount, types.GetSettlementAsset(), pool)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapInt, ok := k.Keeper.ParseToInt(swapResult.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapAmountInt, ok := k.Keeper.ParseToInt(swapAmount.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, swapInt)
			swapAmountInCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapAmountInt)

			nativeAssetCoin = nativeAssetCoin.Add(swapCoin)
			externalAssetCoin = externalAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	// Check and  remove Liquidity
	err = k.Keeper.RemoveLiquidity(ctx, pool, externalAssetCoin, nativeAssetCoin, lp, lpUnitsLeft, poolOriginalEB, poolOriginalNB)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, err.Error())
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
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})

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
		return nil, sdkerrors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}
	// Create Pool
	pool, err := k.Keeper.CreatePool(ctx, poolUnits, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
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
		return nil, sdkerrors.Wrap(types.ErrUnableToAddLiquidity, err.Error())
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
