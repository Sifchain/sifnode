package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryPool:
			return queryPool(ctx, req, keeper)
		case types.QueryPools:
			return queryPools(ctx, keeper)
		case types.QueryLiquidityProvider:
			return queryLiquidityProvider(ctx, req, keeper)
		case types.QueryAssetList:
			return queryAssetList(ctx, req, keeper)
		case types.QueryLPList:
			return queryLPList(ctx, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown clp query endpoint")
		}
	}
}

func queryPool(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryReqGetPool

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	pool, err := keeper.GetPool(ctx, params.Symbol)
	if err != nil {
		return nil, err
	}
	height := ctx.BlockHeight()
	poolResponse := types.NewPoolResponse(pool, height, types.GetCLPModuleAddress().String())
	res, err := codec.MarshalJSONIndent(keeper.cdc, poolResponse)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
func queryPools(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	poolList := keeper.GetPools(ctx)
	if len(poolList) == 0 {
		return nil, types.ErrPoolListIsEmpty
	}
	height := ctx.BlockHeight()
	poolsResponse := types.NewPoolsResponse(poolList, height, types.GetCLPModuleAddress().String())
	res, err := codec.MarshalJSONIndent(keeper.cdc, poolsResponse)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
func queryLiquidityProvider(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryReqLiquidityProvider

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	lp, err := keeper.GetLiquidityProvider(ctx, params.Symbol, params.LpAddress.String())
	if err != nil {
		return nil, err
	}
	pool, err := keeper.GetPool(ctx, params.Symbol)
	if err != nil {
		return nil, err
	}
	native, external, _, _ := CalculateAllAssetsForLP(pool, lp)
	lpResponse := types.NewLiquidityProviderResponse(lp, ctx.BlockHeight(), native.String(), external.String())
	res, err := codec.MarshalJSONIndent(keeper.cdc, lpResponse)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryAssetList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryReqGetAssetList
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	assetList := keeper.GetAssetsForLiquidityProvider(ctx, params.LpAddress)
	res, err := codec.MarshalJSONIndent(keeper.cdc, assetList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryLPList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryReqGetLiquidityProviderList
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	searchingAsset := types.NewAsset(params.Symbol)
	lpList := keeper.GetLiquidityProvidersForAsset(ctx, searchingAsset)
	res, err := codec.MarshalJSONIndent(keeper.cdc, lpList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
