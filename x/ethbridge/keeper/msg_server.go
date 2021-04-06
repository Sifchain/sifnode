package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"strconv"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the ethbridge MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Lock(goCtx context.Context, msg *types.MsgLock) (*types.MsgLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgLockResponse{}, nil

}
func (srv msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	if !srv.Keeper.ExistsPeggyToken(ctx, msg.Symbol) {
		logger.Error("Native token can't be burn.",
			"tokenSymbol", msg.Symbol)
		return nil, errors.Errorf("Native token %s can't be burn.", msg.Symbol)
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, sdk.AccAddress(msg.CosmosSender))
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if err := srv.Keeper.ProcessBurn(ctx, sdk.AccAddress(msg.CosmosSender), msg); err != nil {
		logger.Error("bridge keeper failed to process burn.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit burn event.",
		"EthereumChainID", strconv.FormatInt(msg.EthereumChainId, 10),
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"EthereumReceiver", msg.EthereumReceiver,
		"Amount", msg.Amount.String(),
		"Symbol", msg.Symbol,
		"CethAmount", msg.CethAmount.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeKeyEthereumChainID, strconv.FormatInt(msg.EthereumChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyCethAmount, msg.CethAmount.String()),
		),
	})

	return &types.MsgBurnResponse{}, nil

}
func (srv msgServer) CreateEthBridgeClaim(goCtx context.Context, msg *types.MsgCreateEthBridgeClaim) (*types.MsgCreateEthBridgeClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	logger := srv.Keeper.Logger(ctx)
	var mutex = &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	status, err := srv.Keeper.ProcessClaim(ctx, msg.EthBridgeClaim)
	if err != nil {
		logger.Error("bridge keeper failed to process claim.",
			errorMessageKey, err.Error())
		return nil, err
	}
	if status.Text == oracletypes.StatusText_STATUS_TEXT_SUCCESS {
		if err = srv.Keeper.ProcessSuccessfulClaim(ctx, status.FinalClaim); err != nil {
			logger.Error("bridge keeper failed to process successful claim.",
				errorMessageKey, err.Error())
			return nil, err
		}
	}
	// set mutex lock to false

	logger.Info("sifnode emit create event.",
		"CosmosSender", msg.EthBridgeClaim.ValidatorAddress,
		"EthereumSender", msg.EthBridgeClaim.EthereumSender,
		"EthereumSenderNonce", strconv.FormatInt(msg.EthBridgeClaim.Nonce, 10),
		"CosmosReceiver", msg.EthBridgeClaim.CosmosReceiver,
		"Amount", msg.EthBridgeClaim.Amount.String(),
		"Symbol", msg.EthBridgeClaim.Symbol,
		"ClaimType", msg.EthBridgeClaim.ClaimType.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.EthBridgeClaim.ValidatorAddress),
		),
		sdk.NewEvent(
			types.EventTypeCreateClaim,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.EthBridgeClaim.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyEthereumSender, msg.EthBridgeClaim.EthereumSender),
			sdk.NewAttribute(types.AttributeKeyEthereumSenderNonce, strconv.FormatInt(msg.EthBridgeClaim.Nonce, 10)),
			sdk.NewAttribute(types.AttributeKeyCosmosReceiver, msg.EthBridgeClaim.CosmosReceiver),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.EthBridgeClaim.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.EthBridgeClaim.Symbol),
			sdk.NewAttribute(types.AttributeKeyTokenContract, msg.EthBridgeClaim.TokenContractAddress),
			sdk.NewAttribute(types.AttributeKeyClaimType, msg.EthBridgeClaim.ClaimType.String()),
		),
		sdk.NewEvent(
			types.EventTypeProphecyStatus,
			sdk.NewAttribute(types.AttributeKeyStatus, status.Text.String()),
		),
	})

	return &types.MsgCreateEthBridgeClaimResponse{}, nil
}

func (k msgServer) UpdateWhiteListValidator(goCtx context.Context, msg *types.MsgUpdateWhiteListValidator) (*types.MsgUpdateWhiteListValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgUpdateWhiteListValidatorResponse{}, nil

}