package keeper

import (
	"context"
	"github.com/Sifchain/sifnode/x/whitelist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Querier struct {
	types.Keeper
}

func (q Querier) Denoms(c context.Context, request *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	wl := q.GetDenomWhitelist(ctx)
	return &types.QueryDenomsResponse{List: &wl}, nil
}

var _ types.QueryServer = Querier{}
