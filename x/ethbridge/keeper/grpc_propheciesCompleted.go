package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.ProphciesCompletedQueryServiceServer = prophciesCompletedQueryServiceServer{}

type prophciesCompletedQueryServiceServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the ethbridge QueryServer interface,
// for the provided Keeper.
func NewProphciesCompletedQueryServer(keeper Keeper) types.ProphciesCompletedQueryServiceServer {
	return &prophciesCompletedQueryServiceServer{
		Keeper: keeper,
	}
}

func (srv prophciesCompletedQueryServiceServer) Search(ctx context.Context, req *types.ProphciesCompletedQueryRequest) (*types.ProphciesCompletedQueryResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.NetworkDescriptor
	globalSequence := req.GlobalSequence

	prophecyInfo := srv.Keeper.oracleKeeper.GetProphecyInfoWithScopeGlobalSequence(sdkCtx, networkDescriptor, globalSequence)

	res := types.ProphciesCompletedQueryResponse{
		ProphecyInfo: prophecyInfo,
	}

	return &res, nil
}
