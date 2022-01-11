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
	mtp, err := srv.keeper.GetMTP(sdk.UnwrapSDKContext(ctx), request.CollateralAsset, request.CustodyAsset, request.Address)
	if err != nil {
		return nil, err
	}

	return &types.MTPResponse{Mtp: &mtp}, nil
}
