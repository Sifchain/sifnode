package keeper

import (
	"context"
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

func (srv msgServer) Lock(goCtx context.Context, msg *types.MsgLock) (*types.MsgLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	logger := srv.Keeper.Logger(ctx)
	if srv.Keeper.ExistsPeggyToken(ctx, msg.Symbol) {
		logger.Error("pegged token can't be lock.", "tokenSymbol", msg.Symbol)
		return nil, errors.Errorf("Pegged token %s can't be lock.", msg.Symbol)
	}

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if err := srv.Keeper.ProcessLock(ctx, cosmosSender, msg); err != nil {
		logger.Error("bridge keeper failed to process lock.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit lock event.",
		"NetworkId", strconv.FormatInt(int64(msg.NetworkDescriptor), 10),
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"EthereumReceiver", msg.EthereumReceiver,
		"Amount", msg.Amount.String(),
		"Symbol", msg.Symbol,
		"NativeTokenAmount", msg.NativeTokenAmount.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyNetworkDescriptor, strconv.FormatInt(int64(msg.NetworkDescriptor), 10)),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyNativeTokenAmount, msg.NativeTokenAmount.String()),
		),
	})

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

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if err := srv.Keeper.ProcessBurn(ctx, cosmosSender, msg); err != nil {
		logger.Error("bridge keeper failed to process burn.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit burn event.",
		"NetworkId", strconv.FormatInt(int64(msg.NetworkDescriptor), 10),
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"EthereumReceiver", msg.EthereumReceiver,
		"Amount", msg.Amount.String(),
		"Symbol", msg.Symbol,
		"NativeTokenAmount", msg.NativeTokenAmount.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeKeyNetworkDescriptor, strconv.FormatInt(int64(msg.NetworkDescriptor), 10)),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyNativeTokenAmount, msg.NativeTokenAmount.String()),
		),
	})

	return &types.MsgBurnResponse{}, nil

}
func (srv msgServer) CreateEthBridgeClaim(goCtx context.Context, msg *types.MsgCreateEthBridgeClaim) (*types.MsgCreateEthBridgeClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	status, err := srv.Keeper.ProcessClaim(ctx, msg.EthBridgeClaim)

	if err != nil {
		if err != oracletypes.ErrProphecyFinalized {
			logger.Error("bridge keeper failed to process claim.",
				errorMessageKey, err.Error())
			return nil, err
		}

	} else if status == oracletypes.StatusText_STATUS_TEXT_SUCCESS {
		if err = srv.Keeper.ProcessSuccessfulClaim(ctx, msg.EthBridgeClaim); err != nil {
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
			sdk.NewAttribute(types.AttributeKeyStatus, status.String()),
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

	err = srv.Keeper.ProcessUpdateWhiteListValidator(ctx, msg.NetworkDescriptor, cosmosSender,
		valAddr, msg.Power)
	if err != nil {
		logger.Error("bridge keeper failed to process update validator.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit update whitelist validators event.",
		"NetworkDescriptor", msg.NetworkDescriptor,
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"Validator", msg.Validator,
		"Power", msg.Power)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyNetworkDescriptor, strconv.Itoa(int(msg.NetworkDescriptor))),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyPowerType, strconv.Itoa(int(msg.Power))),
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

	nativeTokenReceiverAddress, err := sdk.AccAddressFromBech32(msg.CethReceiverAccount)
	if err != nil {
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	err = srv.Keeper.ProcessUpdateCethReceiverAccount(ctx,
		cosmosSender, nativeTokenReceiverAddress)
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
		logger.Error("keeper failed to process rescue native_token message.", errorMessageKey, err.Error())
		return nil, err
	}
	logger.Info("sifnode emit rescue native_token event.",
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"CosmosReceiver", msg.CosmosReceiver,
		"NativeTokenAmount", msg.CethAmount)

	return &types.MsgRescueCethResponse{}, nil
}

func (srv msgServer) SetNativeToken(goCtx context.Context, msg *types.MsgSetNativeToken) (*types.MsgSetNativeTokenResponse, error) {
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

	if err := srv.Keeper.ProcessSetNativeToken(ctx, msg); err != nil {
		logger.Error("keeper failed to process rescue native_token message.", errorMessageKey, err.Error())
		return nil, err
	}
	logger.Info("sifnode emit rescue native_token event.",
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"NetworkDescriptor", msg.NetworkDescriptor,
		"NativeToken", msg.NativeToken)

	return &types.MsgSetNativeTokenResponse{}, nil
}
