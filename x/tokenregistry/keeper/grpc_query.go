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

func (q Querier) Entries(c context.Context, _ *types.QueryEntriesRequest) (*types.QueryEntriesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	wl := q.GetDenomWhitelist(ctx)
	return &types.QueryEntriesResponse{List: &wl}, nil
}

var _ types.QueryServer = Querier{}
