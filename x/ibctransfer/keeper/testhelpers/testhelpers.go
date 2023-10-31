package testhelpers

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktransferkeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	sdktransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
)

// Can be passed to sctransferkeeper.MsgServer as the SDK stub,
// so that sctransferkeeper.MsgServer can be used in other tests as well.
type MsgServerStub struct {
	transferKeeper sdktransferkeeper.Keeper
	bankKeeper     sdktransfertypes.BankKeeper
}

func (srv *MsgServerStub) Transfer(ctx context.Context, msg *sdktransfertypes.MsgTransfer) (*sdktransfertypes.MsgTransferResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	_, err = SendStub(sdk.UnwrapSDKContext(ctx), srv.transferKeeper, srv.bankKeeper, msg.Token, sender, msg.SourcePort, msg.SourceChannel)
	if err != nil {
		return nil, err
	}
	return &sdktransfertypes.MsgTransferResponse{}, nil
}

func SendStub(ctx sdk.Context, transferKeeper sdktransferkeeper.Keeper, bankKeeper sdktransfertypes.BankKeeper, token sdk.Coin, sender sdk.AccAddress, sourcePort, sourceChannel string) (string, error) {
	// deconstruct the token denomination into the denomination trace info
	// to determine if the sender is the source chain
	fullDenomPath := token.Denom
	var err error
	if strings.HasPrefix(token.Denom, "ibc/") {
		fullDenomPath, err = transferKeeper.DenomPathFromHash(ctx, token.Denom)
		if err != nil {
			return "", err
		}
	}
	if sdktransfertypes.SenderChainIsSource(sourcePort, sourceChannel, fullDenomPath) {
		// create the escrow address for the tokens
		escrowAddress := sdktransfertypes.GetEscrowAddress(sourcePort, sourceChannel)
		// escrow source tokens. It fails if balance insufficient.
		if err := bankKeeper.SendCoins(
			ctx, sender, escrowAddress, sdk.NewCoins(token),
		); err != nil {
			return "", err
		}
	} else {
		// transfer the coins to the module account and burn them
		if err := bankKeeper.SendCoinsFromAccountToModule(
			ctx, sender, sdktransfertypes.ModuleName, sdk.NewCoins(token),
		); err != nil {
			return "", err
		}
		if err := bankKeeper.BurnCoins(ctx, sdktransfertypes.ModuleName, sdk.NewCoins(token)); err != nil {
			// NOTE: should not happen as the module account was
			// retrieved on the step above and it has enough balace
			// to burn.
			panic(fmt.Sprintf("cannot burn coins after a successful send to a module account: %v", err))
		}
	}
	return fullDenomPath, nil
}
