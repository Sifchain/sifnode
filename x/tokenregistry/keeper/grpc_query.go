package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	types.Keeper
}

func NewQueryServer(k types.Keeper) types.QueryServer {
	return Querier{k}
}

func (q Querier) Entries(c context.Context, request *types.QueryEntriesRequest) (*types.QueryEntriesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	registry, err := q.GetRegistryPaginated(ctx, uint(request.Page), uint(request.Limit))
	if err != nil {
		return &types.QueryEntriesResponse{
			Registry: &registry,
		}, err
	}
	return &types.QueryEntriesResponse{
		Registry: &registry,
	}, nil
}

var _ types.QueryServer = Querier{}
