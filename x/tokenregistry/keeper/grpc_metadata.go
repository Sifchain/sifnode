package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.TokenMetadataServiceServer = metadataServer{}

type metadataServer struct {
	keeper
}

// NewQueryServer returns an implementation of the ethbridge QueryServer interface,
// for the provided Keeper.
func NewTokenMetadataServer(keeper keeper) types.TokenMetadataServiceServer {
	return &metadataServer{
		keeper: keeper,
	}
}

func (srv metadataServer) Search(ctx context.Context, req *types.TokenMetadataSearchRequest) (*types.TokenMetadataSearchResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	denomHash := req.GetDenom()
	metadata, _ := srv.keeper.GetTokenMetadata(sdkCtx, denomHash)

	res := types.TokenMetadataSearchResponse{
		Metadata: &metadata,
	}

	return &res, nil
}
