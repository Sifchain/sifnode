package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sifchain/sifnode/x/clp/types"
)

const MaxPageLimit = 200

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper Keeper
}

var _ types.QueryServer = Querier{}

func (k Querier) GetPool(c context.Context, req *types.PoolReq) (*types.PoolRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	pool, err := k.Keeper.GetPool(ctx, req.Symbol)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "validator %s not found", req.Symbol)
	}
	return &types.PoolRes{
		Pool:             &pool,
		Height:           ctx.BlockHeight(),
		ClpModuleAddress: types.GetCLPModuleAddress().String(),
	}, nil
}

func (k Querier) GetPools(c context.Context, req *types.PoolsReq) (*types.PoolsRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)

	pools, pageRes, err := k.Keeper.GetPoolsPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &types.PoolsRes{
		Pools:            pools,
		Height:           ctx.BlockHeight(),
		ClpModuleAddress: types.GetCLPModuleAddress().String(),
		Pagination:       pageRes,
	}, nil
}

func (k Querier) GetLiquidityProvider(c context.Context, req *types.LiquidityProviderReq) (*types.LiquidityProviderRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	lp, err := k.Keeper.GetLiquidityProvider(ctx, req.Symbol, req.LpAddress)
	if err != nil {
		return nil, err
	}
	pool, err := k.Keeper.GetPool(ctx, req.Symbol)
	if err != nil {
		return nil, err
	}
	native, external, _, _ := CalculateAllAssetsForLP(pool, lp)
	lpResponse := types.NewLiquidityProviderResponse(lp, ctx.BlockHeight(), native.String(), external.String())
	return &lpResponse, nil
}

func (k Querier) GetLiquidityProviderData(c context.Context, req *types.LiquidityProviderDataReq) (*types.LiquidityProviderDataRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.LpAddress)
	if err != nil {
		return nil, err
	}

	if req.Pagination.Limit > MaxPageLimit {
		req.Pagination.Limit = MaxPageLimit
	}
	assetList, _, err := k.Keeper.GetAssetsForLiquidityProviderPaginated(ctx, addr, &query.PageRequest{Limit: req.Pagination.Limit})
	if err != nil {
		return nil, err
	}

	lpDataList := make([]*types.LiquidityProviderData, 0, len(assetList))
	for i := range assetList {
		asset := assetList[i]
		pool, err := k.Keeper.GetPool(ctx, asset.Symbol)
		if err != nil {
			continue
		}
		lp, err := k.Keeper.GetLiquidityProvider(ctx, asset.Symbol, req.LpAddress)
		if err != nil {
			continue
		}
		native, external, _, _ := CalculateAllAssetsForLP(pool, lp)
		lpData := types.NewLiquidityProviderData(lp, native.String(), external.String())
		lpDataList = append(lpDataList, &lpData)
	}

	lpDataResponse := types.NewLiquidityProviderDataResponse(lpDataList, ctx.BlockHeight())
	return &lpDataResponse, nil
}

func (k Querier) GetAssetList(c context.Context, req *types.AssetListReq) (*types.AssetListRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.LpAddress)
	if err != nil {
		return nil, err
	}
	assetList, _, err := k.Keeper.GetAssetsForLiquidityProviderPaginated(ctx, addr, &query.PageRequest{Limit: MaxPageLimit})
	if err != nil {
		return nil, err
	}
	return &types.AssetListRes{
		Assets: assetList,
	}, nil
}

func (k Querier) GetLiquidityProviderList(c context.Context, req *types.LiquidityProviderListReq) (*types.LiquidityProviderListRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)
	searchingAsset := types.NewAsset(req.Symbol)

	lpList, pageRes, err := k.Keeper.GetLiquidityProvidersForAssetPaginated(ctx, searchingAsset, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.LiquidityProviderListRes{
		LiquidityProviders: lpList,
		Height:             ctx.BlockHeight(),
		Pagination:         pageRes,
	}, nil
}

func (k Querier) GetLiquidityProviders(c context.Context, req *types.LiquidityProvidersReq) (*types.LiquidityProvidersRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if req.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(c)

	lpList, pageRes, err := k.Keeper.GetAllLiquidityProvidersPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &types.LiquidityProvidersRes{
		LiquidityProviders: lpList,
		Height:             ctx.BlockHeight(),
		Pagination:         pageRes,
	}, nil
}

func (k Querier) GetPmtpParams(c context.Context, _ *types.PmtpParamsReq) (*types.PmtpParamsRes, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.Keeper.GetPmtpParams(ctx)

	rateParams := k.Keeper.GetPmtpRateParams(ctx)

	epoch := k.Keeper.GetPmtpEpoch(ctx)

	pmtpParamsResponse := types.NewPmtpParamsResponse(params, rateParams, epoch, ctx.BlockHeight())
	return &pmtpParamsResponse, nil
}

func (k Querier) GetParams(c context.Context, _ *types.ParamsReq) (*types.ParamsRes, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.Keeper.GetParams(ctx)
	threshold := k.Keeper.GetSymmetryThreshold(ctx)
	ratioThreshold := k.Keeper.GetSymmetryRatio(ctx)

	return &types.ParamsRes{
		Params:                 &params,
		SymmetryThreshold:      threshold,
		SymmetryRatioThreshold: ratioThreshold,
	}, nil
}

func (k Querier) GetRewardParams(c context.Context, _ *types.RewardParamsReq) (*types.RewardParamsRes, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.Keeper.GetRewardsParams(ctx)

	return &types.RewardParamsRes{Params: params}, nil
}

func (k Querier) GetLiquidityProtectionParams(c context.Context, _ *types.LiquidityProtectionParamsReq) (*types.LiquidityProtectionParamsRes, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.Keeper.GetLiquidityProtectionParams(ctx)
	rateParams := k.Keeper.GetLiquidityProtectionRateParams(ctx)
	response := types.NewLiquidityProtectionParamsResponse(params, rateParams, ctx.BlockHeight())
	return &response, nil
}

func (k Querier) GetProviderDistributionParams(c context.Context, _ *types.ProviderDistributionParamsReq) (*types.ProviderDistributionParamsRes, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.Keeper.GetProviderDistributionParams(ctx)

	return &types.ProviderDistributionParamsRes{Params: params}, nil
}

func (k Querier) GetSwapFeeParams(c context.Context, _ *types.SwapFeeParamsReq) (*types.SwapFeeParamsRes, error) {
	ctx := sdk.UnwrapSDKContext(c)
	swapFeeParams := k.Keeper.GetSwapFeeParams(ctx)

	return &types.SwapFeeParamsRes{DefaultSwapFeeRate: swapFeeParams.DefaultSwapFeeRate, TokenParams: swapFeeParams.TokenParams}, nil
}

func (k Querier) GetPoolShareEstimate(c context.Context, req *types.PoolShareEstimateReq) (*types.PoolShareEstimateRes, error) {
	ctx := sdk.UnwrapSDKContext(c)

	pool, err := k.Keeper.GetPool(ctx, req.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}

	pmtpCurrentRunningRate := k.Keeper.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	sellNativeSwapFeeRate := k.Keeper.GetSwapFeeRate(ctx, types.GetSettlementAsset(), false)
	buyNativeSwapFeeRate := k.Keeper.GetSwapFeeRate(ctx, *req.ExternalAsset, false)

	nativeAssetDepth, externalAssetDepth := pool.ExtractDebt(pool.NativeAssetBalance, pool.ExternalAssetBalance, false)

	newPoolUnits, lpUnits, swapStatus, swapAmount, err := CalculatePoolUnits(
		pool.PoolUnits,
		nativeAssetDepth,
		externalAssetDepth,
		req.NativeAssetAmount,
		req.ExternalAssetAmount,
		sellNativeSwapFeeRate,
		buyNativeSwapFeeRate,
		pmtpCurrentRunningRate)
	if err != nil {
		return nil, err
	}

	feeRate, swapResult, feeAmount, resSwapStatus := calculateSwapInfo(swapStatus, swapAmount, nativeAssetDepth, externalAssetDepth, sellNativeSwapFeeRate, buyNativeSwapFeeRate, pmtpCurrentRunningRate)

	newPoolUnitsD := sdk.NewDecFromBigInt(newPoolUnits.BigInt())
	lpUnitsD := sdk.NewDecFromBigInt(lpUnits.BigInt())

	newNativeAssetDepthD := sdk.NewDecFromBigInt(nativeAssetDepth.Add(req.NativeAssetAmount).BigInt())
	newExternalAssetDepthD := sdk.NewDecFromBigInt(externalAssetDepth.Add(req.ExternalAssetAmount).BigInt())

	percentage := lpUnitsD.Quo(newPoolUnitsD)

	nativeAssetAmountD := percentage.Mul(newNativeAssetDepthD)
	externalAssetAmountD := percentage.Mul(newExternalAssetDepthD)

	return &types.PoolShareEstimateRes{
		Percentage:          percentage,
		NativeAssetAmount:   sdk.NewUintFromBigInt(nativeAssetAmountD.TruncateInt().BigInt()),
		ExternalAssetAmount: sdk.NewUintFromBigInt(externalAssetAmountD.TruncateInt().BigInt()),
		SwapInfo: types.SwapInfo{
			Status:  resSwapStatus,
			Fee:     feeAmount,
			FeeRate: feeRate,
			Amount:  swapAmount,
			Result:  swapResult,
		},
	}, nil

}

func calculateSwapInfo(swapStatus int, swapAmount, nativeAssetDepth, externalAssetDepth sdk.Uint, sellNativeSwapFeeRate, buyNativeSwapFeeRate, pmtpCurrentRunningRate sdk.Dec) (sdk.Dec, sdk.Uint, sdk.Uint, types.SwapStatus) {
	switch swapStatus {
	case NoSwap:
		return sdk.ZeroDec(), sdk.ZeroUint(), sdk.ZeroUint(), types.SwapStatus_NO_SWAP
	case SellNative:
		swapResult, liquidityFee := CalcSwapResult(false, nativeAssetDepth, swapAmount, externalAssetDepth, pmtpCurrentRunningRate, sellNativeSwapFeeRate)
		return sellNativeSwapFeeRate, swapResult, liquidityFee, types.SwapStatus_SELL_NATIVE
	case BuyNative:
		swapResult, liquidityFee := CalcSwapResult(true, externalAssetDepth, swapAmount, nativeAssetDepth, pmtpCurrentRunningRate, buyNativeSwapFeeRate)
		return buyNativeSwapFeeRate, swapResult, liquidityFee, types.SwapStatus_BUY_NATIVE
	default:
		panic("expect not to reach here!")
	}
}

func (k Querier) GetRewardsBucketAll(goCtx context.Context, req *types.AllRewardsBucketReq) (*types.AllRewardsBucketRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var rewardsBuckets []types.RewardsBucket
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.Keeper.storeKey)
	rewardsBucketStore := prefix.NewStore(store, types.KeyPrefix(types.RewardsBucketKeyPrefix))

	pageRes, err := query.Paginate(rewardsBucketStore, req.Pagination, func(key []byte, value []byte) error {
		var rewardsBucket types.RewardsBucket
		if err := k.Keeper.cdc.Unmarshal(value, &rewardsBucket); err != nil {
			return err
		}

		rewardsBuckets = append(rewardsBuckets, rewardsBucket)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.AllRewardsBucketRes{RewardsBucket: rewardsBuckets, Pagination: pageRes}, nil
}

func (k Querier) GetRewardsBucket(goCtx context.Context, req *types.RewardsBucketReq) (*types.RewardsBucketRes, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.Keeper.GetRewardsBucket(
		ctx,
		req.Denom,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.RewardsBucketRes{RewardsBucket: val}, nil
}
