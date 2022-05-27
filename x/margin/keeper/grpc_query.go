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

func (srv queryServer) GetPositionsForAddress(ctx context.Context, request *types.PositionsForAddressRequest) (*types.PositionsForAddressResponse, error) {
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	mtps := srv.keeper.GetMTPsForAddress(sdk.UnwrapSDKContext(ctx), addr)

	return &types.PositionsForAddressResponse{Mtps: mtps}, nil
}
