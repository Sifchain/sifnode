package keeper

import (
	"context"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
	_ *types.QueryAllDistributionsRequest) (*types.QueryAllDistributionsResponse, error) {

	list := q.keeper.GetDistributions(sdk.UnwrapSDKContext(ctx))

	return &types.QueryAllDistributionsResponse{
		Distributions: list.Distributions,
	}, nil
}

func (q Querier) ClaimsByType(ctx context.Context,
	request *types.QueryClaimsByTypeRequest) (*types.QueryClaimsResponse, error) {

	claims := q.keeper.GetClaimsByType(sdk.UnwrapSDKContext(ctx), request.UserClaimType)
	Claims := make([]*types.UserClaim, len(claims))
	for i := range claims {
		Claims = append(Claims, &claims[i])
	}
	return &types.QueryClaimsResponse{
		Claims: Claims,
	}, nil
}

func (q Querier) RecordsByDistributionName(ctx context.Context, request *types.QueryRecordsByDistributionNameRequest) (*types.QueryRecordsByDistributionNameResponse, error) {
	records := new(types.DistributionRecords)
	*records = q.keeper.GetRecordsForName(sdk.UnwrapSDKContext(ctx), request.DistributionName, request.Status)
	if request.Status == types.DistributionStatus_DISTRIBUTION_STATUS_UNSPECIFIED {
		records.DistributionRecords = append(records.DistributionRecords,
			q.keeper.GetRecordsForName(sdk.UnwrapSDKContext(ctx), request.DistributionName, types.DistributionStatus_DISTRIBUTION_STATUS_PENDING).DistributionRecords...)
	}
	return &types.QueryRecordsByDistributionNameResponse{
		DistributionRecords: records,
	}, nil
}

func (q Querier) RecordsByRecipient(ctx context.Context, request *types.QueryRecordsByRecipientAddrRequest) (*types.QueryRecordsByRecipientAddrResponse, error) {
	records := q.keeper.GetRecordsForRecipient(sdk.UnwrapSDKContext(ctx), request.Address)

	return &types.QueryRecordsByRecipientAddrResponse{
		DistributionRecords: &records,
	}, nil
}
