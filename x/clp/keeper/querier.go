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
	pool, err := keeper.GetPool(ctx, params.Ticker)
	if err != nil {
		return nil, err
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, pool)
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
	res, err := codec.MarshalJSONIndent(keeper.cdc, poolList)
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
	lp, err := keeper.GetLiquidityProvider(ctx, params.Ticker, params.LpAddress.String())
	if err != nil {
		return nil, err
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, lp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
