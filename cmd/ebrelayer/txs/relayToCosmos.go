package txs

// DONTCOVER

import (
	ethbridge "github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/instrumentation"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"go.uber.org/zap"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	errorMessageKey = "errorMessage"
)

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(factory tx.Factory, claims []*ethbridge.EthBridgeClaim, cliCtx client.Context, sugaredLogger *zap.SugaredLogger) error {
	var messages []sdk.Msg

	sugaredLogger.Infow(
		"relay prophecies to cosmos.",
		"claimAmount", len(claims),
	)
	instrumentation.PeggyCheckpointZap(sugaredLogger, "EthBridgeClaimArray", zap.Reflect("claims", claims))

	for _, claim := range claims {
		// Packages the claim as a Tendermint message
		msg := ethbridge.NewMsgCreateEthBridgeClaim(claim)

		err := msg.ValidateBasic()
		if err != nil {
			sugaredLogger.Errorw(
				"failed to get message from claim.",
				"message", msg,
				errorMessageKey, err.Error(),
			)
			continue
		} else {
			messages = append(messages, &msg)
		}
	}

	sugaredLogger.Infow("RelayToCosmos building, signing, and broadcasting", "messages", messages)
	instrumentation.PeggyCheckpointZap(sugaredLogger, "BroadcastTx", zap.Reflect("messages", messages))

	err := tx.BroadcastTx(
		cliCtx,
		factory.
			WithGas(1000000000000000000).
			WithFees("500000000000000000rowan"),
		messages...,
	)

	// Broadcast to a Tendermint node
	// open question as to how we handle this situation.
	//    do we retry,
	//        if so, how many times do we try?
	if err != nil {
		sugaredLogger.Errorw(
			"failed to broadcast tx to sifchain.",
			errorMessageKey, err.Error(),
		)
		return err
	}

	return nil
}

// SignProphecyToCosmos broadcasts the signed prophecy message to cosmos
func SignProphecyToCosmos(factory tx.Factory, signProphecy ethbridge.MsgSignProphecy, cliCtx client.Context, sugaredLogger *zap.SugaredLogger) {
	var messages []sdk.Msg

	messages = append(messages, &signProphecy)

	sugaredLogger.Infow("RelayToCosmos building, signing, and broadcasting", "messages", messages)
	// TODO this WithGas isn't correct
	err := tx.BroadcastTx(
		cliCtx,
		factory.
			WithGas(1000000000000000000).
			WithFees("500000000000000000rowan"),
		messages...,
	)

	// Broadcast to a Tendermint node
	// open question as to how we handle this situation.
	//    do we retry,
	//        if so, how many times do we try?
	if err == nil {
		sugaredLogger.Infow("Broadcast SignProphecyToCosmos tx without error")
	} else {
		sugaredLogger.Errorw(
			"failed to broadcast tx to sifchain.",
			errorMessageKey, err.Error(),
		)
	}
}
