package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

const maxPageSize = 100

type queryServer struct {
	bankkeeper.Keeper
}

func NewSifQueryServer(keeper bankkeeper.Keeper) queryServer {
	return queryServer{keeper}
}

func (srv queryServer) AllBalances(ctx context.Context, req *types.QueryAllBalancesRequest) (*types.QueryAllBalancesResponse, error) {
	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{}
	}

	if req.Pagination.Limit > maxPageSize {
		req.Pagination.Limit = maxPageSize
	}

	return srv.Keeper.AllBalances(ctx, req)
}
