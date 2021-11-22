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

func (srv queryServer) EthereumLockBurnSequence(ctx context.Context, req *types.QueryEthereumLockBurnSequenceRequest) (*types.QueryEthereumLockBurnSequenceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.GetNetworkDescriptor()
	relayerValAddress := req.RelayerValAddress

	address, err := sdk.ValAddressFromBech32(relayerValAddress)

	if err != nil {
		return nil, err
	}

	LockBurnSequence := srv.Keeper.GetEthereumLockBurnSequence(sdkCtx, networkDescriptor, address)

	res := types.NewEthereumLockBurnSequenceResponse(LockBurnSequence)

	return &res, nil
}

func (srv queryServer) WitnessLockBurnSequence(ctx context.Context, req *types.QueryWitnessLockBurnSequenceRequest) (*types.QueryWitnessLockBurnSequenceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.GetNetworkDescriptor()
	relayerValAddress := req.RelayerValAddress

	address, err := sdk.ValAddressFromBech32(relayerValAddress)

	if err != nil {
		return nil, err
	}

	LockBurnSequence := srv.Keeper.GetWitnessLockBurnSequence(sdkCtx, networkDescriptor, address)

	res := types.NewWitnessLockBurnSequenceResponse(LockBurnSequence)

	return &res, nil
}

func (srv queryServer) GlobalSequenceBlockNumber(ctx context.Context, req *types.QueryGlobalSequenceBlockNumberRequest) (*types.QueryGlobalSequenceBlockNumberResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.GetNetworkDescriptor()
	globalSequence := req.GlobalSequence

	blockNumber := srv.Keeper.GetGlobalSequenceToBlockNumber(sdkCtx, networkDescriptor, globalSequence)

	res := types.NewGlobalSequenceBlockNumberResponse(blockNumber)

	return &res, nil
}

func (srv queryServer) ProphciesCompleted(ctx context.Context, req *types.QueryProphciesCompletedRequest) (*types.QueryProphciesCompletedResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.NetworkDescriptor
	globalSequence := req.GlobalSequence

	prophecyInfo := srv.Keeper.oracleKeeper.GetProphecyInfoWithScopeGlobalSequence(sdkCtx, networkDescriptor, globalSequence)

	res := types.QueryProphciesCompletedResponse{
		ProphecyInfo: prophecyInfo,
	}

	return &res, nil
}
