package keeper

import (
	"context"
	"github.com/Sifchain/sifnode/x/whitelist/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	types.Keeper
}

func NewQueryServer(k types.Keeper) types.QueryServer {
	return Querier{k}
}

func (q Querier) Denoms(c context.Context, _ *types.QueryDenomsRequest) (*types.QueryDenomsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	wl := q.GetDenomWhitelist(ctx)
	return &types.QueryDenomsResponse{List: &wl}, nil
}

var _ types.QueryServer = Querier{}
