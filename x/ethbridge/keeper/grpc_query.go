package keeper

import (
	"context"
	"strconv"

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

	id := strconv.FormatInt(int64(req.NetworkId), 10) + strconv.FormatInt(req.Nonce, 10) + req.EthereumSender

	prophecy, found := srv.Keeper.oracleKeeper.GetProphecy(sdkCtx, id)
	if !found {
		return nil, sdkerrors.Wrap(oracletypes.ErrProphecyNotFound, id)
	}

	bridgeClaims, err := types.MapOracleClaimsToEthBridgeClaims(
		req.NetworkId,
		types.NewEthereumAddress(req.BridgeContractAddress),
		req.Nonce,
		req.Symbol,
		types.NewEthereumAddress(req.TokenContractAddress),
		types.NewEthereumAddress(req.EthereumSender),
		prophecy.ValidatorClaims,
		types.CreateEthClaimFromOracleString,
	)
	if err != nil {
		return nil, err
	}

	res := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	return &res, nil
}
