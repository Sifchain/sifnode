package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxPageLimit = 200

type Querier struct {
	keeper Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{
		keeper: k,
	}
}

var _ types.QueryServer = Querier{}

func (q Querier) AllDistributions(ctx context.Context,
	req *types.QueryAllDistributionsRequest) (*types.QueryAllDistributionsResponse, error) {

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

	c := sdk.UnwrapSDKContext(ctx)
	list, pageRes, err := q.keeper.GetDistributionsPaginated(sdk.UnwrapSDKContext(ctx), req.Pagination)
	if err != nil {
		return nil, err
	}
	return &types.QueryAllDistributionsResponse{
		Distributions: list.Distributions,
		Height:        c.BlockHeight(),
		Pagination:    pageRes,
	}, nil
}

func (q Querier) ClaimsByType(ctx context.Context,
	request *types.QueryClaimsByTypeRequest) (*types.QueryClaimsResponse, error) {
	claims := q.keeper.GetClaimsByType(sdk.UnwrapSDKContext(ctx), request.UserClaimType)
	return &types.QueryClaimsResponse{
		Claims: claims.UserClaims,
	}, nil
}

func (q Querier) RecordsByDistributionName(ctx context.Context, request *types.QueryRecordsByDistributionNameRequest) (*types.QueryRecordsByDistributionNameResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if request.Pagination == nil {
		request.Pagination = &query.PageRequest{
			Limit: MaxPageLimit,
		}
	}

	if request.Pagination.Limit > MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", MaxPageLimit))
	}

	c := sdk.UnwrapSDKContext(ctx)
	records, pageRes, err := q.keeper.GetRecordsForNameAndStatusPaginated(sdk.UnwrapSDKContext(ctx), request.DistributionName, request.Status, request.Pagination)
	if err != nil {
		return nil, err
	}
	return &types.QueryRecordsByDistributionNameResponse{
		DistributionRecords: records,
		Height:              c.BlockHeight(),
		Pagination:          pageRes,
	}, nil
}

func (q Querier) RecordsByRecipient(ctx context.Context, request *types.QueryRecordsByRecipientAddrRequest) (*types.QueryRecordsByRecipientAddrResponse, error) {
	records := q.keeper.GetRecordsForRecipient(sdk.UnwrapSDKContext(ctx), request.Address)

	return &types.QueryRecordsByRecipientAddrResponse{
		DistributionRecords: records,
	}, nil
}
