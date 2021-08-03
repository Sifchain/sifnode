package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	querier := Querier{keeper}

	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryPool:
			return queryPool(ctx, path[1:], req, legacyQuerierCdc, querier)
		case types.QueryPools:
			return queryPools(ctx, path[1:], legacyQuerierCdc, querier)
		case types.QueryLiquidityProvider:
			return queryLiquidityProvider(ctx, path[1:], req, legacyQuerierCdc, querier)
		case types.QueryAssetList:
			return queryAssetList(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case types.QueryLPList:
			return queryLPList(ctx, path[1:], req, keeper, legacyQuerierCdc, querier)
		case types.QueryAllLP:
			return queryAllLP(ctx, path[1:], keeper, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown clp query endpoint")
		}
	}
}

func queryPool(ctx sdk.Context, path []string, req abci.RequestQuery,
	legacyQuerierCdc *codec.LegacyAmino, querier Querier) ([]byte, error) {

	var params types.PoolReq

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	res, err := querier.GetPool(sdk.WrapSDKContext(ctx), &params)
	if err != nil {
		return nil, err
	}

	bz, err := legacyQuerierCdc.MarshalJSON(&res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryPools(ctx sdk.Context, path []string, legacyQuerierCdc *codec.LegacyAmino, querier Querier) ([]byte, error) {
	res, err := querier.GetPools(sdk.WrapSDKContext(ctx), &types.PoolsReq{})
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryLiquidityProvider(ctx sdk.Context, path []string, req abci.RequestQuery,
	legacyQuerierCdc *codec.LegacyAmino, querier Querier) ([]byte, error) {
	var params types.LiquidityProviderReq

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	res, err := querier.GetLiquidityProvider(sdk.WrapSDKContext(ctx), &params)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryAssetList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.AssetListReq
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addr, err := sdk.AccAddressFromBech32(params.LpAddress)
	if err != nil {
		return nil, err
	}

	assetList := keeper.GetAssetsForLiquidityProvider(ctx, addr)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, assetList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryLPList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino, querier Querier) ([]byte, error) {
	var params types.LiquidityProviderListReq

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	res, err := querier.GetLiquidityProviderList(sdk.WrapSDKContext(ctx), &params)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryAllLP(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	lpList := keeper.GetAllLiquidityProviders(ctx)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, lpList)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
