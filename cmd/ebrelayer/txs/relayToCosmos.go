package txs

// DONTCOVER

import (
	"log"
	"sync/atomic"

	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"go.uber.org/zap"
)

var (
	nextSequenceNumber uint64 = 0
	errorMessageKey           = "errorMessage"
)

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(moniker, password string, claims []*types.EthBridgeClaim, cliCtx client.Context,
	txBldr client.TxBuilder, sugaredLogger *zap.SugaredLogger) error {
	var messages []sdk.Msg

	sugaredLogger.Infow("relay prophecies to cosmos.",
		"claimAmount", len(claims),
		"nextSequenceNumber", nextSequenceNumber)

	for _, claim := range claims {
		// Packages the claim as a Tendermint message
		msg := types.NewMsgCreateEthBridgeClaim(claim)

		e := msg.ValidateBasic()
		if e != nil {
			sugaredLogger.Errorw("failed to get message from claim.",
				errorMessageKey, e.Error())
			continue
		} else {
			messages = append(messages, &msg)
		}
	}

	sugaredLogger.Infow("relay sequenceNumber from builder.",
		"nextSequenceNumber", nextSequenceNumber)

	// If we start to control sequence
	if nextSequenceNumber > 0 {
		sugaredLogger.Infow("txBldr.WithSequence(nextSequenceNumber) passed")
	}

	log.Println("building and signing")

	log.Println("built tx, now broadcasting")
	// Broadcast to a Tendermint node
	err := tx.GenerateOrBroadcastTxCLI(cliCtx, nil, messages...)
	if err != nil {
		sugaredLogger.Errorw("failed to broadcast tx to sifchain.",
			errorMessageKey, err.Error())
		return err
	}
	log.Println("Broadcasted tx without error")

	// start to control sequence number after first successful tx
	if nextSequenceNumber == 0 {
		setNextSequenceNumber(nextSequenceNumber + 1)
	} else {
		incrementNextSequenceNumber()
	}
	sugaredLogger.Infow("relay next sequenceNumber from memory.",
		"nextSequenceNumber", nextSequenceNumber)

	return nil
}

func incrementNextSequenceNumber() {
	atomic.AddUint64(&nextSequenceNumber, 1)
}

func setNextSequenceNumber(sequenceNumber uint64) {
	atomic.StoreUint64(&nextSequenceNumber, sequenceNumber)
}
