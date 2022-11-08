package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/admin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	Keeper
}

func (q Querier) ListAccounts(ctx context.Context, _ *types.ListAccountsRequest) (*types.ListAccountsResponse, error) {
	al := q.GetAdminAccounts(sdk.UnwrapSDKContext(ctx))
	return &types.ListAccountsResponse{
		Keys: al,
	}, nil
}

func (q Querier) GetParams(ctx context.Context, _ *types.GetParamsRequest) (*types.GetParamsResponse, error) {
	return &types.GetParamsResponse{Params: q.Keeper.GetParams(sdk.UnwrapSDKContext(ctx))}, nil
}

func NewQueryServer(k Keeper) types.QueryServer {
	return Querier{k}
}

var _ types.QueryServer = Querier{}
