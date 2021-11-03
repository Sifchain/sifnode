package keeper

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
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

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		logger.Error("validator address is wrong", errorMessageKey, err.Error())
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	tokenMetadata, ok := srv.Keeper.GetTokenMetadata(ctx, msg.DenomHash)
	if !ok {
		logger.Error("token metadata not available", "DenomHash", msg.DenomHash)
		return &types.MsgLockResponse{}, fmt.Errorf("token metadata not available for %s", msg.DenomHash)
	}

	doublePeg := tokenMetadata.NetworkDescriptor != msg.NetworkDescriptor
	firstDoublePeg := doublePeg && srv.tokenRegistryKeeper.GetFirstLockDoublePeg(ctx, msg.DenomHash, msg.NetworkDescriptor)

	prophecyID, err := srv.Keeper.ProcessLock(ctx, cosmosSender, account.GetSequence(), msg, tokenMetadata, firstDoublePeg)

	if err != nil {
		logger.Error("bridge keeper failed to process lock.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit lock event.", "message", msg)
	globalNonce := srv.Keeper.GetGlobalNonce(ctx, msg.NetworkDescriptor)
	srv.Keeper.UpdateGlobalNonce(ctx, msg.NetworkDescriptor, uint64(ctx.BlockHeight()))
	srv.tokenRegistryKeeper.SetFirstLockDoublePeg(ctx, msg.DenomHash, msg.NetworkDescriptor)

	err = srv.oracleKeeper.SetProphecyInfo(ctx,
		prophecyID,
		msg.NetworkDescriptor,
		cosmosSender.String(),
		account.GetSequence(),
		msg.EthereumReceiver,
		msg.DenomHash,
		tokenMetadata.TokenAddress,
		msg.Amount,
		msg.CrosschainFee,
		doublePeg,
		globalNonce,
	)

	if err != nil {
		logger.Error("bridge keeper failed to set prophecy info.", errorMessageKey, err.Error())
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyNetworkDescriptor, strconv.FormatInt(int64(msg.NetworkDescriptor), 10)),
			sdk.NewAttribute(types.AttributeKeyProphecyID, string(prophecyID[:])),
			sdk.NewAttribute(types.AttributeKeyGlobalNonce, strconv.FormatInt(int64(globalNonce), 10)),
		),
	})

	return &types.MsgLockResponse{}, nil
}

func (srv msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		logger.Error("validator address is wrong", errorMessageKey, err.Error())
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	tokenMetadata, ok := srv.Keeper.GetTokenMetadata(ctx, msg.DenomHash)
	if !ok {
		logger.Error("token metadata not available", "DenomHash", msg.DenomHash)
		return nil, fmt.Errorf("token metadata not available for %s", msg.DenomHash)
	}

	prophecyID, err := srv.Keeper.ProcessBurn(ctx, cosmosSender, account.GetSequence(), msg, tokenMetadata)

	if err != nil {
		logger.Error("bridge keeper failed to process burn.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit burn event.", "message", msg)
	globalNonce := srv.Keeper.GetGlobalNonce(ctx, msg.NetworkDescriptor)
	srv.Keeper.UpdateGlobalNonce(ctx, msg.NetworkDescriptor, uint64(ctx.BlockHeight()))

	doublePeg := tokenMetadata.NetworkDescriptor != msg.NetworkDescriptor

	err = srv.oracleKeeper.SetProphecyInfo(ctx,
		prophecyID,
		msg.NetworkDescriptor,
		cosmosSender.String(),
		account.GetSequence(),
		msg.EthereumReceiver,
		msg.DenomHash,
		tokenMetadata.TokenAddress,
		msg.Amount,
		msg.CrosschainFee,
		doublePeg,
		globalNonce,
	)

	if err != nil {
		logger.Error("bridge keeper failed to set prophecy info.", errorMessageKey, err.Error())
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeKeyNetworkDescriptor, strconv.FormatInt(int64(msg.NetworkDescriptor), 10)),
			sdk.NewAttribute(types.AttributeKeyProphecyID, string(prophecyID[:])),
			sdk.NewAttribute(types.AttributeKeyGlobalNonce, strconv.FormatInt(int64(globalNonce), 10)),
		),
	})

	return &types.MsgBurnResponse{}, nil

}
func (srv msgServer) CreateEthBridgeClaim(goCtx context.Context, msg *types.MsgCreateEthBridgeClaim) (*types.MsgCreateEthBridgeClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	logger.Debug(types.PeggyTestMarker, "kind", "CreateEthBridgeClaim", "message", zap.Reflect("message", msg))

	// check the account
	cosmosSender := msg.EthBridgeClaim.ValidatorAddress
	valAddress, err := sdk.ValAddressFromBech32(cosmosSender)
	if err != nil {
		logger.Error("validator address is wrong", errorMessageKey, err.Error())
		return nil, err
	}

	// check the lock burn nonce
	lockBurnNonce := srv.Keeper.GetEthereumLockBurnNonce(ctx, msg.EthBridgeClaim.NetworkDescriptor, valAddress)

	newLockBurnNonce := msg.EthBridgeClaim.EthereumLockBurnNonce

	if newLockBurnNonce != 0 && newLockBurnNonce != lockBurnNonce+1 {
		logger.Error("lock burn nonce out of order", errorMessageKey, err.Error())
		return nil, errors.New("lock burn nonce out of order")
	}

	status, err := srv.Keeper.ProcessClaim(ctx, msg.EthBridgeClaim)

	if err != nil {
		if err != oracletypes.ErrProphecyFinalized {
			logger.Error("bridge keeper failed to process claim.", errorMessageKey, err.Error())
			return nil, err
		}

	} else if status == oracletypes.StatusText_STATUS_TEXT_SUCCESS {
		if err = srv.Keeper.ProcessSuccessfulClaim(ctx, msg.EthBridgeClaim); err != nil {
			logger.Error("bridge keeper failed to process successful claim.", errorMessageKey, err.Error())
			return nil, err
		}
	}

	claim := msg.EthBridgeClaim
	metadata := tokenregistrytypes.TokenMetadata{
		Decimals:          claim.Decimals,
		Name:              claim.TokenName,
		Symbol:            claim.Symbol,
		TokenAddress:      claim.TokenContractAddress,
		NetworkDescriptor: claim.NetworkDescriptor,
	}

	srv.Keeper.AddTokenMetadata(ctx, metadata)

	// update lock burn nonce in keeper
	srv.Keeper.SetEthereumLockBurnNonce(ctx, msg.EthBridgeClaim.NetworkDescriptor, valAddress, newLockBurnNonce)

	logger.Info("sifnode emit create event.",
		"CosmosSender", claim.ValidatorAddress,
		"EthereumSender", claim.EthereumSender,
		"EthereumSenderNonce", strconv.FormatInt(claim.Nonce, 10),
		"CosmosReceiver", claim.CosmosReceiver,
		"Amount", claim.Amount.String(),
		"Symbol", claim.Symbol,
		"ClaimType", claim.ClaimType.String(),
		"DenomHash", claim.DenomHash,
	)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, claim.ValidatorAddress),
		),
		sdk.NewEvent(
			types.EventTypeProphecyStatus,
			sdk.NewAttribute(types.AttributeKeyStatus, status.String()),
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
	})

	return &types.MsgCreateEthBridgeClaimResponse{}, nil
}

func (srv msgServer) UpdateWhiteListValidator(goCtx context.Context,
	msg *types.MsgUpdateWhiteListValidator) (*types.MsgUpdateWhiteListValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		logger.Error("cosmos address is wrong", errorMessageKey, err.Error())
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, msg.CosmosSender)
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		logger.Error("validator address is wrong", errorMessageKey, err.Error())
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
		logger.Error("cosmos address is wrong", errorMessageKey, err.Error())
		return nil, err
	}

	crossChainFeeReceiverAddress, err := sdk.AccAddressFromBech32(msg.CrosschainFeeReceiver)
	if err != nil {
		logger.Error("cosmos receiver address is wrong", errorMessageKey, err.Error())
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
		logger.Error("validator address is wrong", errorMessageKey, err.Error())
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
		logger.Error("cosmos address is wrong", errorMessageKey, err.Error())
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	if err := srv.Keeper.SetFeeInfo(ctx, msg); err != nil {
		logger.Error("keeper failed to process setting crosschain fee message.", errorMessageKey, err.Error())
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

// SignProphecy relayer sign the prophecy ID and send to Sifchain after receive the burn/lock events
func (srv msgServer) SignProphecy(goCtx context.Context, msg *types.MsgSignProphecy) (*types.MsgSignProphecyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := srv.Keeper.Logger(ctx)

	cosmosSender, err := sdk.AccAddressFromBech32(msg.CosmosSender)
	if err != nil {
		logger.Error("cosmos address is wrong", errorMessageKey, err.Error())
		return nil, err
	}

	account := srv.Keeper.accountKeeper.GetAccount(ctx, cosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender)
	}

	err = srv.Keeper.ProcessSignProphecy(ctx, msg)

	// if error is ErrProphecyFinalized, will continue and emit event, not return error.
	if err != nil && err != oracletypes.ErrProphecyFinalized {
		logger.Error("keeper failed to process rescue native_token message.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode received the sign prophecy message.",
		"Message", msg)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender),
		),
		sdk.NewEvent(
			types.EventTypeSignProphecy,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender),
			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
			sdk.NewAttribute(types.AttributeKeyNetworkDescriptor, msg.NetworkDescriptor.String()),
			sdk.NewAttribute(types.AttributeKeyProphecyID, string(msg.ProphecyId)),
		),
	})

	return &types.MsgSignProphecyResponse{}, nil
}
