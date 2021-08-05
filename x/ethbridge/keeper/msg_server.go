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

	logger.Info("sifnode emit lock event.", "message", msg)

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
			sdk.NewAttribute(types.AttributeKeycrossChainFee, msg.CrosschainFee.String()),
		),
	})

	return &types.MsgLockResponse{}, nil
}

func (srv msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	if !srv.Keeper.ExistsPeggyToken(ctx, msg.Symbol) {
		logger.Error("crosschain fee can't be burn.",
			"tokenSymbol", msg.Symbol)
		return nil, errors.Errorf("crosschain fee %s can't be burn.", msg.Symbol)
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

	logger.Info("sifnode emit burn event.", "message", msg)

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
			sdk.NewAttribute(types.AttributeKeycrossChainFee, msg.CrosschainFee.String()),
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

	claim := msg.EthBridgeClaim

	metadata := types.TokenMetadata{
		Decimals:          claim.Decimals,
		Name:              claim.TokenName,
		Symbol:            claim.Symbol,
		TokenAddress:      claim.TokenContractAddress,
		NetworkDescriptor: claim.NetworkDescriptor,
	}

	if !srv.Keeper.ExistsTokenMetadata(ctx, claim.DenomHash) {
		srv.Keeper.AddTokenMetadata(ctx, metadata)
	}

	logger.Info("sifnode emit create event.",
		"CosmosSender", claim.ValidatorAddress,
		"EthereumSender", claim.EthereumSender,
		"EthereumSenderNonce", strconv.FormatInt(claim.Nonce, 10),
		"CosmosReceiver", claim.CosmosReceiver,
		"Amount", claim.Amount.String(),
		"Symbol", claim.Symbol,
		"ClaimType", claim.ClaimType.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, claim.ValidatorAddress),
		),
		sdk.NewEvent(
			types.EventTypeCreateClaim,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, claim.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyEthereumSender, claim.EthereumSender),
			sdk.NewAttribute(types.AttributeKeyEthereumSenderNonce, strconv.FormatInt(claim.Nonce, 10)),
			sdk.NewAttribute(types.AttributeKeyCosmosReceiver, claim.CosmosReceiver),
			sdk.NewAttribute(types.AttributeKeyAmount, claim.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, claim.Symbol),
			sdk.NewAttribute(types.AttributeKeyTokenContract, claim.TokenContractAddress),
			sdk.NewAttribute(types.AttributeKeyClaimType, claim.ClaimType.String()),
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

func (srv msgServer) UpdateCrossChainFeeReceiverAccount(goCtx context.Context,
	msg *types.MsgUpdateCrossChainFeeReceiverAccount) (*types.MsgUpdateCrossChainFeeReceiverAccountResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		return nil, err
	}

	crossChainFeeReceiverAddress, err := sdk.AccAddressFromBech32(msg.CrosschainFeeReceiver)
	if err != nil {
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	err = srv.Keeper.ProcessUpdateCrossChainFeeReceiverAccount(ctx,
		cosmosSender, crossChainFeeReceiverAddress)
	if err != nil {
		logger.Error("keeper failed to process update crosschain fee receiver account.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit update crosschain fee receiver account event.",
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"CrossChainFeeReceiverAccount", msg.CrosschainFeeReceiver)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyCrossChainFeeReceiverAccount, msg.CrosschainFeeReceiver),
		),
	})

	return &types.MsgUpdateCrossChainFeeReceiverAccountResponse{}, nil
}

func (srv msgServer) RescueCrossChainFee(goCtx context.Context, msg *types.MsgRescueCrossChainFee) (*types.MsgRescueCrossChainFeeResponse, error) {
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

	if err := srv.Keeper.RescueCrossChainFees(ctx, msg); err != nil {
		logger.Error("keeper failed to process rescue crosschain_fee message.", errorMessageKey, err.Error())
		return nil, err
	}
	logger.Info("sifnode emit rescue crosschain_fee event.",
		"CosmosSender", msg.CosmosSender,
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"CosmosReceiver", msg.CosmosReceiver,
		"crossChainFee", msg.CrosschainFee)

	return &types.MsgRescueCrossChainFeeResponse{}, nil
}

func (srv msgServer) SetFeeInfo(goCtx context.Context, msg *types.MsgSetFeeInfo) (*types.MsgSetFeeInfoResponse, error) {
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

	if err := srv.Keeper.SetFeeInfo(ctx, msg); err != nil {
		logger.Error("keeper failed to process rescue crosschain fee message.", errorMessageKey, err.Error())
		return nil, err
	}
	logger.Info("sifnode emit set crosschain fee event.",
		"Message", msg)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeSetCrossChainFee,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyNetworkDescriptor, msg.NetworkDescriptor.String()),
			sdk.NewAttribute(types.AttributeKeyCrossChainFee, msg.FeeCurrency),
			sdk.NewAttribute(types.AttributeKeyCrossChainFeeGas, msg.FeeCurrencyGas.String()),
			sdk.NewAttribute(types.AttributeKeyMinimumLockCost, msg.MinimumLockCost.String()),
			sdk.NewAttribute(types.AttributeKeyMinimumBurnCost, msg.MinimumBurnCost.String()),
		),
	})

	return &types.MsgSetFeeInfoResponse{}, nil
}
