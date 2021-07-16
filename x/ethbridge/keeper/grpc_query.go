package keeper

import (
	"context"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the ethbridge QueryServer interface,
// for the provided Keeper.
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &queryServer{
		Keeper: keeper,
	}
}

func (srv queryServer) EthProphecy(ctx context.Context, req *types.QueryEthProphecyRequest) (*types.QueryEthProphecyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	id := req.GetProphecyId()
	prophecy, found := srv.Keeper.oracleKeeper.GetProphecy(sdkCtx, id)
	if !found {
		return nil, sdkerrors.Wrap(oracletypes.ErrProphecyNotFound, string(id))
	}

	res := types.NewQueryEthProphecyResponse(id, prophecy.Status, prophecy.ClaimValidators)

	return &res, nil
}

func (srv queryServer) CrosschainFeeConfig(ctx context.Context, req *types.QueryCrosschainFeeConfigRequest) (*types.QueryCrosschainFeeConfigResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.GetNetworkDescriptor()
	networkIdentity := oracletypes.NewNetworkIdentity(networkDescriptor)
	crosschainFeeConfig, err := srv.Keeper.oracleKeeper.GetCrossChainFeeConfig(sdkCtx, networkIdentity)
	if err != nil {
		return nil, sdkerrors.Wrap(oracletypes.ErrProphecyNotFound, networkDescriptor.String())
	}

	res := types.NewQueryCrosschainFeeConfigResponse(crosschainFeeConfig)

	return &res, nil
}
