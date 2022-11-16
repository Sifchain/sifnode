package keeper

import (
	"context"
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"

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

func (srv msgServer) SetPause(goCtx context.Context, msg *types.MsgPause) (*types.MsgPauseResponse, error) {
	response := &types.MsgPauseResponse{}
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}
	if !srv.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_ETHBRIDGE, signer) {
		return response, types.ErrNotEnoughPermissions
	}

	srv.Keeper.SetPause(ctx, &types.Pause{IsPaused: msg.IsPaused})
	return response, nil
}

func (srv msgServer) Lock(goCtx context.Context, msg *types.MsgLock) (*types.MsgLockResponse, error) {
	response := &types.MsgLockResponse{}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if srv.Keeper.IsPaused(ctx) {
		return response, types.ErrPaused
	}
	logger := srv.Keeper.Logger(ctx)
	if srv.Keeper.ExistsPeggyToken(ctx, msg.Symbol) {
		logger.Error("pegged token can't be lock.", "tokenSymbol", msg.Symbol)
		return response, errors.Errorf("pegged token %s can't be locked", msg.Symbol)
	}
	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return response, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return response, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}
	if err := srv.Keeper.ProcessLock(ctx, cosmosSender, msg); err != nil {
		logger.Error("bridge keeper failed to process lock.", errorMessageKey, err.Error())
		return response, err
	}
	logger.Info("sifnode emit lock event.",
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
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyEthereumChainID, strconv.FormatInt(msg.EthereumChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyCethAmount, msg.CethAmount.String()),
		),
	})

	return response, nil
}

func (srv msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	response := &types.MsgBurnResponse{}
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)
	if srv.Keeper.IsPaused(ctx) {
		return response, types.ErrPaused
	}
	if !srv.Keeper.ExistsPeggyToken(ctx, msg.Symbol) {
		logger.Error("Native token can't be burn.",
			"tokenSymbol", msg.Symbol)
		return response, errors.Errorf("native token %s can't be burned", msg.Symbol)
	}

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return response, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return response, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if err := srv.Keeper.ProcessBurn(ctx, cosmosSender, msg); err != nil {
		logger.Error("bridge keeper failed to process burn.", errorMessageKey, err.Error())
		return response, err
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

	return response, nil

}
func (srv msgServer) CreateEthBridgeClaim(goCtx context.Context, msg *types.MsgCreateEthBridgeClaim) (*types.MsgCreateEthBridgeClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	logger := srv.Keeper.Logger(ctx)

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

func (srv msgServer) UpdateWhiteListValidator(goCtx context.Context,
	msg *types.MsgUpdateWhiteListValidator) (*types.MsgUpdateWhiteListValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)

		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, msg.CosmosSender)
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Validator)
	}

	err = srv.Keeper.ProcessUpdateWhiteListValidator(ctx, cosmosSender,
		valAddr, msg.OperationType)
	if err != nil {
		logger.Error("bridge keeper failed to process update validator.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit update whitelist validators event.",
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"Validator", msg.Validator,
		"OperationType", msg.OperationType)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyOperationType, msg.OperationType),
		),
	})

	return &types.MsgUpdateWhiteListValidatorResponse{}, nil
}

func (srv msgServer) UpdateCethReceiverAccount(goCtx context.Context,
	msg *types.MsgUpdateCethReceiverAccount) (*types.MsgUpdateCethReceiverAccountResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return nil, err
	}

	cethReceiverAddress, err := sdk.AccAddressFromBech32(msg.CethReceiverAccount)
	if err != nil {
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	err = srv.Keeper.ProcessUpdateCethReceiverAccount(ctx,
		cosmosSender, cethReceiverAddress)
	if err != nil {
		logger.Error("keeper failed to process update ceth receiver account.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit update ceth receiver account event.",
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"CethReceiverAccount", msg.CethReceiverAccount)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyCethReceiverAccount, msg.CethReceiverAccount),
		),
	})

	return &types.MsgUpdateCethReceiverAccountResponse{}, nil
}

func (srv msgServer) RescueCeth(goCtx context.Context, msg *types.MsgRescueCeth) (*types.MsgRescueCethResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if err := srv.Keeper.ProcessRescueCeth(ctx, msg); err != nil {
		logger.Error("keeper failed to process rescue ceth message.", errorMessageKey, err.Error())
		return nil, err
	}
	logger.Info("sifnode emit rescue ceth event.",
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"CosmosReceiver", msg.CosmosReceiver,
		"CethAmount", msg.CethAmount)

	return &types.MsgRescueCethResponse{}, nil
}

func (srv msgServer) SetBlacklist(goCtx context.Context, msg *types.MsgSetBlacklist) (*types.MsgSetBlacklistResponse, error) {
	err := srv.Keeper.SetBlacklist(sdk.UnwrapSDKContext(goCtx), msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetBlacklistResponse{}, nil
}
