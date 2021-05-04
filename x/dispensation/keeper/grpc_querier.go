package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/dispensation/types"
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

func (q Querier) RecordsByDistributionName(ctx context.Context, request *types.QueryRecordsByDistributionNameRequest) (*types.QueryRecordsByDistributionNameResponse, error) {
	records := new(types.DistributionRecords)
	switch request.Status {
	case types.ClaimStatus_CLAIM_STATUS_PENDING:
		*records = q.keeper.GetRecordsForNamePending(sdk.UnwrapSDKContext(ctx), request.DistributionName)
	case types.ClaimStatus_CLAIM_STATUS_COMPLETED:
		*records = q.keeper.GetRecordsForNameCompleted(sdk.UnwrapSDKContext(ctx), request.DistributionName)
	default:
		*records = q.keeper.GetRecordsForNameAll(sdk.UnwrapSDKContext(ctx), request.DistributionName)
	}

	return &types.QueryRecordsByDistributionNameResponse{
		DistributionRecords: records,
	}, nil
}

func (q Querier) RecordsByRecipient(ctx context.Context, request *types.QueryRecordsByRecipientAddrRequest) (*types.QueryRecordsByRecipientAddrResponse, error) {
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	records := q.keeper.GetRecordsForRecipient(sdk.UnwrapSDKContext(ctx), addr)

	return &types.QueryRecordsByRecipientAddrResponse{
		DistributionRecords: records,
	}, nil
}

