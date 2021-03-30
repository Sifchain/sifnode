//nolint:dupl
package ethbridge

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/oracle"
)

const errorMessageKey = "errorMessage"

// NewHandler returns a handler for "ethbridge" type messages.
func NewHandler(
	accountKeeper types.AccountKeeper, bridgeKeeper Keeper,
	cdc *codec.Codec) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateEthBridgeClaim:
			return handleMsgCreateEthBridgeClaim(ctx, cdc, bridgeKeeper, msg)
		case MsgBurn:
			return handleMsgBurn(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		case MsgLock:
			return handleMsgLock(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		case MsgUpdateWhiteListValidator:
			return handleMsgUpdateWhiteListValidator(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized ethbridge message type: %v", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handle a message to create a bridge claim
func handleMsgCreateEthBridgeClaim(
	ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgCreateEthBridgeClaim,
) (*sdk.Result, error) {

	logger := bridgeKeeper.Logger(ctx)
	var mutex = &sync.RWMutex{}
	mutex.Lock()
	defer mutex.Unlock()

	status, err := bridgeKeeper.ProcessClaim(ctx, types.EthBridgeClaim(msg))
	if err != nil {
		logger.Error("bridge keeper failed to process claim.",
			errorMessageKey, err.Error())
		return nil, err
	}
	if status.Text == oracle.SuccessStatusText {
		if err = bridgeKeeper.ProcessSuccessfulClaim(ctx, status.FinalClaim); err != nil {
			logger.Error("bridge keeper failed to process successful claim.",
				errorMessageKey, err.Error())
			return nil, err
		}
	}
	// set mutex lock to false

	logger.Info("sifnode emit create event.",
		"CosmosSender", msg.ValidatorAddress.String(),
		"EthereumSender", msg.EthereumSender.String(),
		"EthereumSenderNonce", strconv.Itoa(msg.Nonce),
		"CosmosReceiver", msg.CosmosReceiver.String(),
		"Amount", msg.Amount.String(),
		"Symbol", msg.Symbol,
		"ClaimType", msg.ClaimType.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
		sdk.NewEvent(
			types.EventTypeCreateClaim,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.ValidatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyEthereumSender, msg.EthereumSender.String()),
			sdk.NewAttribute(types.AttributeKeyEthereumSenderNonce, strconv.Itoa(msg.Nonce)),
			sdk.NewAttribute(types.AttributeKeyCosmosReceiver, msg.CosmosReceiver.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyTokenContract, msg.TokenContractAddress.String()),
			sdk.NewAttribute(types.AttributeKeyClaimType, msg.ClaimType.String()),
		),
		sdk.NewEvent(
			types.EventTypeProphecyStatus,
			sdk.NewAttribute(types.AttributeKeyStatus, status.Text.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBurn(
	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
	bridgeKeeper Keeper, msg MsgBurn,
) (*sdk.Result, error) {
	logger := bridgeKeeper.Logger(ctx)

	if !bridgeKeeper.ExistsPeggyToken(ctx, msg.Symbol) {
		logger.Error("Native token can't be burn.",
			"tokenSymbol", msg.Symbol)
		return nil, errors.Errorf("Native token %s can't be burn.", msg.Symbol)
	}

	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender.String())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	var coins sdk.Coins

	if msg.Symbol == CethSymbol {
		coins = sdk.NewCoins(sdk.NewCoin(CethSymbol, msg.CethAmount.Add(msg.Amount)))
	} else {
		coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount), sdk.NewCoin(CethSymbol, msg.CethAmount))
	}
	if err := bridgeKeeper.ProcessBurn(ctx, msg.CosmosSender, coins); err != nil {
		logger.Error("bridge keeper failed to process burn.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit burn event.",
		"EthereumChainID", strconv.Itoa(msg.EthereumChainID),
		"CosmosSender", msg.CosmosSender.String(),
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"EthereumReceiver", msg.EthereumReceiver.String(),
		"Amount", msg.Amount.String(),
		"Symbol", msg.Symbol,
		"Coins", coins.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeKeyEthereumChainID, strconv.Itoa(msg.EthereumChainID)),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyCoins, coins.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}

func handleMsgLock(
	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
	bridgeKeeper Keeper, msg MsgLock,
) (*sdk.Result, error) {
	logger := bridgeKeeper.Logger(ctx)
	if bridgeKeeper.ExistsPeggyToken(ctx, msg.Symbol) {
		logger.Error("pegged token can't be lock.", "tokenSymbol", msg.Symbol)
		return nil, errors.Errorf("Pegged token %s can't be lock.", msg.Symbol)
	}

	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender.String())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount), sdk.NewCoin(CethSymbol, msg.CethAmount))
	if err := bridgeKeeper.ProcessLock(ctx, msg.CosmosSender, coins); err != nil {
		logger.Error("bridge keeper failed to process lock.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit lock event.",
		"EthereumChainID", strconv.Itoa(msg.EthereumChainID),
		"CosmosSender", msg.CosmosSender.String(),
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"EthereumReceiver", msg.EthereumReceiver.String(),
		"Amount", msg.Amount.String(),
		"Symbol", msg.Symbol,
		"Coins", coins.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyEthereumChainID, strconv.Itoa(msg.EthereumChainID)),
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyCoins, coins.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

}

func handleMsgUpdateWhiteListValidator(
	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
	bridgeKeeper Keeper, msg MsgUpdateWhiteListValidator,
) (*sdk.Result, error) {
	logger := bridgeKeeper.Logger(ctx)
	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
	if account == nil {
		logger.Error("account is nil.", "CosmosSender", msg.CosmosSender.String())

		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	if err := bridgeKeeper.ProcessUpdateWhiteListValidator(ctx, msg.CosmosSender, msg.Validator, msg.OperationType); err != nil {
		logger.Error("bridge keeper failed to process update validator.", errorMessageKey, err.Error())
		return nil, err
	}

	logger.Info("sifnode emit update whitelist validators event.",
		"CosmosSender", msg.CosmosSender.String(),
		"CosmosSenderSequence", strconv.FormatUint(account.GetSequence(), 10),
		"Validator", msg.Validator.String(),
		"OperationType", msg.OperationType)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
		),
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
			sdk.NewAttribute(types.AttributeKeyOperationType, msg.OperationType),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
