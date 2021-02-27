package txs

// DONTCOVER

import (
	"log"
	"sync/atomic"

	"github.com/Sifchain/sifnode/x/ethbridge"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"go.uber.org/zap"
)

var (
	nextSequenceNumber uint64 = 0
)

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(cdc *codec.Codec, moniker, password string, claims []types.EthBridgeClaim, cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder, sugaredLogger *zap.SugaredLogger) error {
	var messages []sdk.Msg
	// log.Printf("ebrelayer RelayToCosmos with %d claims\n", len(claims))
	// log.Printf("ebrelayer RelayToCosmos nextSequenceNumber is %d\n", nextSequenceNumber)

	sugaredLogger.Infow("relay prophecies to cosmos.",
		"claim amount", len(claims),
		"next sequence number", nextSequenceNumber)

	for _, claim := range claims {
		// Packages the claim as a Tendermint message
		msg := ethbridge.NewMsgCreateEthBridgeClaim(claim)

		err := msg.ValidateBasic()
		if err != nil {
			// log.Println("failed to get message from claim with:", err.Error())
			sugaredLogger.Errorw("failed to get message from claim.",
				"error message", err.Error())
			continue
		} else {
			messages = append(messages, msg)
		}
	}

	// Prepare tx
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		// log.Println("error building tx: ", err)
		// log.Println("tx buidler response on error: ", txBldr)
		sugaredLogger.Errorw("failed to get tx builder.",
			"error message", err.Error(),
			"transaction builder", txBldr)
		return err
	}

	// log.Printf("ebrelayer RelayToCosmos sequenceNumber is %d from tx builder\n", txBldr.Sequence())
	sugaredLogger.Infow("relay sequenceNumber from builder.",
		"next sequence number", txBldr.Sequence())

	// If we start to control sequence
	if nextSequenceNumber > 0 {
		txBldr.WithSequence(nextSequenceNumber)
		sugaredLogger.Infow("txBldr.WithSequence(nextSequenceNumber) passed")
	}

	log.Println("building and signing")
	// Build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(moniker, password, messages)
	if err != nil {
		// log.Printf("ebrelayer failed to sign transaction ", err.Error())
		sugaredLogger.Errorw("failed to sign transaction.",
			"error message", err.Error())
		return err
	}

	log.Println("built tx, now broadcasting")
	// Broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTxAsync(txBytes)
	if err != nil {
		// log.Printf("ebrelayer RelayToCosmos error 2 is %s\n", err.Error())
		sugaredLogger.Errorw("failed to broadcast tx to sifchain.",
			"error message", err.Error())
		return err
	}
	log.Println("Broadcasted tx without error")

	if err = cliCtx.PrintOutput(res); err != nil {
		// log.Printf("ebrelayer RelayToCosmos error 3 is %s\n", err.Error())
		sugaredLogger.Errorw("failed to print out result.",
			"error message", err.Error())
		return err
	}

	// start to control sequence number after first successful tx
	if nextSequenceNumber == 0 {
		setNextSequenceNumber(txBldr.Sequence() + 1)
	} else {
		incrementNextSequenceNumber()
	}
	// log.Printf("ebrelayer RelayToCosmos nextSequenceNumber is %d after tx\n", nextSequenceNumber)
	sugaredLogger.Infow("relay next sequenceNumber from memory.",
		"next sequence number", nextSequenceNumber)

	return nil
}

func incrementNextSequenceNumber() {
	atomic.AddUint64(&nextSequenceNumber, 1)
}

func setNextSequenceNumber(sequenceNumber uint64) {
	atomic.StoreUint64(&nextSequenceNumber, sequenceNumber)
}
