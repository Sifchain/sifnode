package txs

// DONTCOVER

import (
	"go.uber.org/zap"
	"log"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"

	// tx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	errorMessageKey           = "errorMessage"
)

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(factory tx.Factory, claims []*types.EthBridgeClaim, cliCtx client.Context, sugaredLogger *zap.SugaredLogger) error {
	var messages []sdk.Msg

	sugaredLogger.Infow(
		"relay prophecies to cosmos.",
		"claimAmount", len(claims),
	)

	for _, claim := range claims {
		// Packages the claim as a Tendermint message
		msg := types.NewMsgCreateEthBridgeClaim(claim)

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

	sugaredLogger.Infow(
		"relay sequenceNumber from builder.",
	)

	sugaredLogger.Infow("RelayToCosmos building, signing, and broadcasting", "messages", messages)
	err := tx.BroadcastTx(cliCtx, factory.WithGas(1000000000000000000).WithFees("500000000000000000rowan"), messages...)

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

	log.Println("Broadcasted tx without error")

	return nil
}