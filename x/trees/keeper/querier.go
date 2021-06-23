package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/trees/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for faucet clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetTreeById:
			return queryTreeById(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown faucet query endpoint")
		}
	}
}

func queryTreeById(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryReqGetTreeByID

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	tree := keeper.GetTree(ctx, params.Id)
	res, err := codec.MarshalJSONIndent(keeper.cdc, tree)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
