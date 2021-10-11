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

func (srv queryServer) EthereumLockBurnNonce(ctx context.Context, req *types.QueryEthereumLockBurnNonceRequest) (*types.QueryEthereumLockBurnNonceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.GetNetworkDescriptor()
	relayerValAddress := req.RelayerValAddress

	address, err := sdk.ValAddressFromBech32(relayerValAddress)

	if err != nil {
		return nil, err
	}

	lockBurnNonce := srv.Keeper.GetEthereumLockBurnNonce(sdkCtx, networkDescriptor, address)

	res := types.NewEthereumLockBurnNonceResponse(lockBurnNonce)

	return &res, nil
}

func (srv queryServer) WitnessLockBurnNonce(ctx context.Context, req *types.QueryWitnessLockBurnNonceRequest) (*types.QueryWitnessLockBurnNonceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.GetNetworkDescriptor()
	relayerValAddress := req.RelayerValAddress

	address, err := sdk.ValAddressFromBech32(relayerValAddress)

	if err != nil {
		return nil, err
	}

	lockBurnNonce := srv.Keeper.GetWitnessLockBurnNonce(sdkCtx, networkDescriptor, address)

	res := types.NewWitnessLockBurnNonceResponse(lockBurnNonce)

	return &res, nil
}

func (srv queryServer) GlocalNonceBlockNumber(ctx context.Context, req *types.QueryGlocalNonceBlockNumberRequest) (*types.QueryGlocalNonceBlockNumberResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networkDescriptor := req.GetNetworkDescriptor()
	globalNonce := req.GlobalNonce

	blockNumber, err := srv.Keeper.GetGlocalNonceToBlockNumber(sdkCtx, networkDescriptor, globalNonce)
	if err != nil {
		return nil, err
	}

	res := types.NewGlocalNonceBlockNumberResponse(blockNumber)

	return &res, nil
}
