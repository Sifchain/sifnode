package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
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
func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgBurnResponse{}, nil

}
func (k msgServer) CreateEthBridgeClaim(goCtx context.Context, msg *types.MsgCreateEthBridgeClaim) (*types.MsgCreateEthBridgeClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgCreateEthBridgeClaimResponse{}, nil

}
func (k msgServer) UpdateWhiteListValidator(goCtx context.Context, msg *types.MsgUpdateWhiteListValidator) (*types.MsgUpdateWhiteListValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgUpdateWhiteListValidatorResponse{}, nil

}

// // Handle a message to create a bridge claim
// func handleMsgCreateEthBridgeClaim(
// 	ctx sdk.Context, cdc *codec.Codec, bridgeKeeper Keeper, msg MsgCreateEthBridgeClaim,
// ) (*sdk.Result, error) {
// 	var mutex = &sync.RWMutex{}
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	status, err := bridgeKeeper.ProcessClaim(ctx, types.EthBridgeClaim(msg))
// 	if err != nil {
// 		fmt.Printf("Sifnode handleMsgCreateEthBridgeClaim 46 %s\n", err.Error())
// 		return nil, err
// 	}
// 	if status.Text == oracle.SuccessStatusText {
// 		if err = bridgeKeeper.ProcessSuccessfulClaim(ctx, status.FinalClaim); err != nil {
// 			fmt.Printf("Sifnode handleMsgCreateEthBridgeClaim 51 %s\n", err.Error())
// 			return nil, err
// 		}
// 	}
// 	// set mutex lock to false

// 	fmt.Printf("Sifnode handleMsgCreateEthBridgeClaim 56 all done, emit events statue is %s\n", status.Text.String())
// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
// 			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
// 		),
// 		sdk.NewEvent(
// 			types.EventTypeCreateClaim,
// 			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.ValidatorAddress.String()),
// 			sdk.NewAttribute(types.AttributeKeyEthereumSender, msg.EthereumSender.String()),
// 			sdk.NewAttribute(types.AttributeKeyEthereumSenderNonce, strconv.Itoa(msg.Nonce)),
// 			sdk.NewAttribute(types.AttributeKeyCosmosReceiver, msg.CosmosReceiver.String()),
// 			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
// 			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
// 			sdk.NewAttribute(types.AttributeKeyTokenContract, msg.TokenContractAddress.String()),
// 			sdk.NewAttribute(types.AttributeKeyClaimType, msg.ClaimType.String()),
// 		),
// 		sdk.NewEvent(
// 			types.EventTypeProphecyStatus,
// 			sdk.NewAttribute(types.AttributeKeyStatus, status.Text.String()),
// 		),
// 	})

// 	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
// }

// func handleMsgBurn(
// 	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
// 	bridgeKeeper Keeper, msg MsgBurn,
// ) (*sdk.Result, error) {
// 	if !bridgeKeeper.ExistsPeggyToken(ctx, msg.Symbol) {
// 		return nil, errors.Errorf("Native token %s can't be burn.", msg.Symbol)
// 	}

// 	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
// 	if account == nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
// 	}

// 	var coins sdk.Coins

// 	if msg.Symbol == CethSymbol {
// 		coins = sdk.NewCoins(sdk.NewCoin(CethSymbol, msg.CethAmount.Add(msg.Amount)))
// 	} else {
// 		coins = sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount), sdk.NewCoin(CethSymbol, msg.CethAmount))
// 	}
// 	if err := bridgeKeeper.ProcessBurn(ctx, msg.CosmosSender, coins); err != nil {
// 		return nil, err
// 	}

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
// 			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
// 		),
// 		sdk.NewEvent(
// 			types.EventTypeBurn,
// 			sdk.NewAttribute(types.AttributeKeyEthereumChainID, strconv.Itoa(msg.EthereumChainID)),
// 			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
// 			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
// 			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver.String()),
// 			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
// 			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
// 			sdk.NewAttribute(types.AttributeKeyCoins, coins.String()),
// 		),
// 	})

// 	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

// }

// func handleMsgLock(
// 	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
// 	bridgeKeeper Keeper, msg MsgLock,
// ) (*sdk.Result, error) {
// 	if bridgeKeeper.ExistsPeggyToken(ctx, msg.Symbol) {
// 		return nil, errors.Errorf("Pegged token %s can't be lock.", msg.Symbol)
// 	}

// 	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
// 	if account == nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
// 	}

// 	coins := sdk.NewCoins(sdk.NewCoin(msg.Symbol, msg.Amount), sdk.NewCoin(CethSymbol, msg.CethAmount))
// 	if err := bridgeKeeper.ProcessLock(ctx, msg.CosmosSender, coins); err != nil {
// 		return nil, err
// 	}

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
// 			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
// 		),
// 		sdk.NewEvent(
// 			types.EventTypeLock,
// 			sdk.NewAttribute(types.AttributeKeyEthereumChainID, strconv.Itoa(msg.EthereumChainID)),
// 			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
// 			sdk.NewAttribute(types.AttributeKeyCosmosSenderSequence, strconv.FormatUint(account.GetSequence(), 10)),
// 			sdk.NewAttribute(types.AttributeKeyEthereumReceiver, msg.EthereumReceiver.String()),
// 			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
// 			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
// 			sdk.NewAttribute(types.AttributeKeyCoins, coins.String()),
// 		),
// 	})

// 	return &sdk.Result{Events: ctx.EventManager().Events()}, nil

// }

// func handleMsgUpdateWhiteListValidator(
// 	ctx sdk.Context, cdc *codec.Codec, accountKeeper types.AccountKeeper,
// 	bridgeKeeper Keeper, msg MsgUpdateWhiteListValidator,
// ) (*sdk.Result, error) {
// 	account := accountKeeper.GetAccount(ctx, msg.CosmosSender)
// 	if account == nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
// 	}

// 	if err := bridgeKeeper.ProcessUpdateWhiteListValidator(ctx, msg.CosmosSender, msg.Validator, msg.OperationType); err != nil {
// 		return nil, err
// 	}

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
// 			sdk.NewAttribute(sdk.AttributeKeySender, msg.CosmosSender.String()),
// 		),
// 		sdk.NewEvent(
// 			types.EventTypeLock,
// 			sdk.NewAttribute(types.AttributeKeyCosmosSender, msg.CosmosSender.String()),
// 			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator.String()),
// 			sdk.NewAttribute(types.AttributeKeyEthereumChainID, msg.OperationType),
// 		),
// 	})

// 	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
// }
