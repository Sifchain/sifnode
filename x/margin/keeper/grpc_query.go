//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	keeper types.Keeper
}

func NewQueryServer(k types.Keeper) types.QueryServer {
	return &queryServer{k}
}

func (srv queryServer) GetMTP(ctx context.Context, request *types.MTPRequest) (*types.MTPResponse, error) {
	mtp, err := srv.keeper.GetMTP(sdk.UnwrapSDKContext(ctx), request.Address, request.Id)
	if err != nil {
		return nil, err
	}

	return &types.MTPResponse{Mtp: &mtp}, nil
}

func (srv queryServer) GetPositionsForAddress(goCtx context.Context, request *types.PositionsForAddressRequest) (*types.PositionsForAddressResponse, error) {
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	mtps, pageRes, err := srv.keeper.GetMTPsForAddress(sdk.UnwrapSDKContext(goCtx), addr, request.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.PositionsForAddressResponse{Mtps: mtps, Pagination: pageRes}, nil
}

func (srv queryServer) GetParams(ctx context.Context, request *types.ParamsRequest) (*types.ParamsResponse, error) {
	params := srv.keeper.GetParams(sdk.UnwrapSDKContext(ctx))

	return &types.ParamsResponse{Params: &params}, nil
}
