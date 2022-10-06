package keeper

import (
	"context"
	"fmt"

	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (srv queryServer) GetPositionsByPool(ctx context.Context, request *types.PositionsByPoolRequest) (*types.PositionsByPoolResponse, error) {
	mtps, pageRes, err := srv.keeper.GetMTPsForPool(sdk.UnwrapSDKContext(ctx), request.Asset, request.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.PositionsByPoolResponse{
		Mtps:       mtps,
		Pagination: pageRes,
	}, nil
}

func (srv queryServer) GetPositions(ctx context.Context, request *types.PositionsRequest) (*types.PositionsResponse, error) {
	if request.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	mtps, page, err := srv.keeper.GetMTPs(sdk.UnwrapSDKContext(ctx), request.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.PositionsResponse{
		Mtps:       mtps,
		Pagination: page,
	}, nil
}

func (srv queryServer) GetStatus(ctx context.Context, request *types.StatusRequest) (*types.StatusResponse, error) {
	return &types.StatusResponse{
		OpenMtpCount:     srv.keeper.GetOpenMTPCount(sdk.UnwrapSDKContext(ctx)),
		LifetimeMtpCount: srv.keeper.GetMTPCount(sdk.UnwrapSDKContext(ctx)),
	}, nil
}

func (srv queryServer) GetWhitelist(ctx context.Context, request *types.WhitelistRequest) (*types.WhitelistResponse, error) {
	if request.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	whitelist, page, err := srv.keeper.GetWhitelist(sdk.UnwrapSDKContext(ctx), request.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.WhitelistResponse{
		Whitelist:  whitelist,
		Pagination: page,
	}, nil
}

func (srv queryServer) GetSQParams(ctx context.Context, request *types.GetSQParamsRequest) (*types.GetSQParamsResponse, error) {
	pool, err := srv.keeper.ClpKeeper().GetPool(sdk.UnwrapSDKContext(ctx), request.Pool)
	if err != nil {
		return nil, err
	}
	return &types.GetSQParamsResponse{
		BeginBlock: int64(srv.keeper.GetSQBeginBlock(sdk.UnwrapSDKContext(ctx), &pool)),
	}, nil
}

func (srv queryServer) IsWhitelisted(ctx context.Context, request *types.IsWhitelistedRequest) (*types.IsWhitelistedResponse, error) {
	return &types.IsWhitelistedResponse{
		Address:       request.Address,
		IsWhitelisted: srv.keeper.IsWhitelisted(sdk.UnwrapSDKContext(ctx), request.Address),
	}, nil
}
