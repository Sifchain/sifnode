package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.TokenMetadataServiceServer = metadataServer{}

type metadataServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the ethbridge QueryServer interface,
// for the provided Keeper.
func NewTokenMetadataServer(keeper Keeper) types.TokenMetadataServiceServer {
	return &metadataServer{
		Keeper: keeper,
	}
}

func (srv metadataServer) Search(ctx context.Context, req *types.TokenMetadataRequest) (*types.TokenMetadataResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	denomHash := req.GetDenom()
	metadata := srv.Keeper.GetTokenMetadata(sdkCtx, denomHash)

	res := types.TokenMetadataResponse{
		Metadata: &metadata,
	}

	return &res, nil
}
