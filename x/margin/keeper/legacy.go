//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewLegacyQuerier(k types.Keeper, cdc *codec.LegacyAmino) sdk.Querier {
	querier := queryServer{k}
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryMTPsForAddress:
			return queryMtpsForAddress(ctx, req, cdc, querier)
		default:
			return nil, sdkerrors.Wrap(types.ErrUnknownRequest, "unknown request")
		}
	}
}

func NewLegacyHandler(k types.Keeper) sdk.Handler {
	msgServer := NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgOpen:
			res, err := msgServer.Open(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgClose:
			res, err := msgServer.Close(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgForceClose:
			res, err := msgServer.ForceClose(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdateParams:
			res, err := msgServer.UpdateParams(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUpdatePools:
			res, err := msgServer.UpdatePools(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized margin message type: %T", msg)
		}
	}
}

func queryMtpsForAddress(ctx sdk.Context, req abci.RequestQuery, cdc *codec.LegacyAmino, querier queryServer) ([]byte, error) {
	params := types.PositionsForAddressRequest{}
	err := cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	res, err := querier.GetPositionsForAddress(sdk.WrapSDKContext(ctx), &params)
	if err != nil {
		return nil, err
	}
	bz, err := codec.MarshalJSONIndent(cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}
