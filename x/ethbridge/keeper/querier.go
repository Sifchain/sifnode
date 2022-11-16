package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

// TODO: move to x/oracle

// NewLegacyQuerier is the module level router for state queries
func NewLegacyQuerier(keeper Keeper, cdc *codec.LegacyAmino) sdk.Querier { //nolint

	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryEthProphecy:
			return legacyQueryEthProphecy(ctx, cdc, req, keeper)
		case types.QueryBlacklist:
			return legacyQueryBlacklist(ctx, cdc, req, keeper)
		case types.QueryPause:
			return legacyQueryPause(ctx, cdc, req, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown ethbridge query endpoint")
		}
	}
}

func legacyQueryEthProphecy(ctx sdk.Context, cdc *codec.LegacyAmino, query abci.RequestQuery, keeper Keeper) ([]byte, error) { //nolint
	var req types.QueryEthProphecyRequest

	if err := cdc.UnmarshalJSON(query.Data, &req); err != nil {
		return nil, sdkerrors.Wrap(types.ErrJSONMarshalling, fmt.Sprintf("failed to parse req: %s", err.Error()))
	}

	queryServer := NewQueryServer(keeper)
	response, err := queryServer.EthProphecy(sdk.WrapSDKContext(ctx), &req)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSONIndent(response, "", "  ")
}

func legacyQueryBlacklist(ctx sdk.Context, cdc *codec.LegacyAmino, query abci.RequestQuery, keeper Keeper) ([]byte, error) { //nolint
	var req types.QueryBlacklistRequest

	if err := cdc.UnmarshalJSON(query.Data, &req); err != nil {
		return nil, sdkerrors.Wrap(types.ErrJSONMarshalling, fmt.Sprintf("failed to parse req: %s", err.Error()))
	}

	queryServer := NewQueryServer(keeper)
	response, err := queryServer.GetBlacklist(sdk.WrapSDKContext(ctx), &req)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSONIndent(response, "", "  ")
}

func legacyQueryPause(ctx sdk.Context, cdc *codec.LegacyAmino, query abci.RequestQuery, keeper Keeper) ([]byte, error) { //nolint
	var req types.QueryPauseRequest
	if err := cdc.UnmarshalJSON(query.Data, &req); err != nil {
		return nil, sdkerrors.Wrap(types.ErrJSONMarshalling, fmt.Sprintf("failed to parse req: %s", err.Error()))
	}
	queryServer := NewQueryServer(keeper)
	response, err := queryServer.GetPauseStatus(sdk.WrapSDKContext(ctx), &req)
	if err != nil {
		return nil, err
	}
	return cdc.MarshalJSONIndent(response, "", "  ")
}
