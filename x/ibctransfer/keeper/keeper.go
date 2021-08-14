package keeper

import (
	"context"

	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	bankKeeper          bankkeeper.Keeper
	tokenRegistryKeeper tokenregistrytypes.Keeper
	sdkMsgServer        types.MsgServer
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl( /* bankKeeper, tokenRegistryKeeper */ ) types.MsgServer {
	return &msgServer{}
}

// Transfer defines a rpc handler method for MsgTransfer.
func (srv msgServer) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// get token registry entry for sent token
	registryEntry := srv.tokenRegistryKeeper.GetDenom(ctx, msg.Token.Denom)
	// check if registry entry has an IBC decimal field
	if registryEntry.IbcDenom != "" && registryEntry.Decimals > registryEntry.IbcDecimals {

		po := registryEntry.Decimals - registryEntry.IbcDecimals
		decAmount := sdk.NewDecFromInt(msg.Token.Amount)
		convAmountDec := ReducePrecision(decAmount, po)

		convAmount := sdk.NewIntFromBigInt(convAmountDec.RoundInt().BigInt())

		// convAmount := msg.Token.Amount / (10 **(uint64(registryEntry.Decimals) - uint64(registryEntry.IbcDecimals)))
		convToken := sdk.NewCoin(registryEntry.IbcDenom, convAmount)
		// send coins from account to module
		token := sdk.NewCoin(msg.Token.Denom, msg.Token.Amount)
		err = srv.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(token))
		if err != nil {
			return nil, err
		}
		// mint ibcdenom coins
		err = srv.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(convToken))
		if err != nil {
			return nil, err
		}
		// send coins from module account to address
		err = srv.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(convToken))
		if err != nil {
			return nil, err
		}
		msg.Token = convToken
	}

	return srv.sdkMsgServer.Transfer(goCtx, msg)
}

func ReducePrecision(dec sdk.Dec, po int64) sdk.Dec {
	p := sdk.NewDec(10).Power(uint64(po))
	return dec.Quo(p)
}
