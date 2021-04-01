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
	errorMessageKey           = "errorMessage"
)

// SendGasPrice applies validator's signature to an updateGasPrice message containing
// information about Ethereum block number and gas price
func SendGasPrice(cdc *codec.Codec, moniker, password string, validator sdk.ValAddress, blockNumber sdk.Int, gasPrice sdk.Int, cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder, sugaredLogger *zap.SugaredLogger) error {

	msg := ethbridge.NewMsgUpdateGasPrice(validator, blockNumber, gasPrice)

	return sendMessagesToCosmos(cdc, moniker, password, []sdk.Msg{msg}, cliCtx, txBldr, sugaredLogger)
}

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(cdc *codec.Codec, moniker, password string, claims []types.EthBridgeClaim, cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder, sugaredLogger *zap.SugaredLogger) error {
	var messages []sdk.Msg

	sugaredLogger.Infow("relay prophecies to cosmos.",
		"claimAmount", len(claims),
		"nextSequenceNumber", nextSequenceNumber)

	for _, claim := range claims {
		// Packages the claim as a Tendermint message
		msg := ethbridge.NewMsgCreateEthBridgeClaim(claim)

		err := msg.ValidateBasic()
		if err != nil {
			sugaredLogger.Errorw("failed to get message from claim.",
				errorMessageKey, err.Error())
			continue
		} else {
			messages = append(messages, msg)
		}
	}

	return sendMessagesToCosmos(cdc, moniker, password, messages, cliCtx, txBldr, sugaredLogger)
}

// sendMessagesToCosmos send the messages to cosmos
func sendMessagesToCosmos(cdc *codec.Codec, moniker, password string, messages []sdk.Msg, cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder, sugaredLogger *zap.SugaredLogger) error {

	// Prepare tx
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		sugaredLogger.Errorw("failed to get tx builder.",
			errorMessageKey, err.Error(),
			"transactionBuilder", txBldr)
		return err
	}

	sugaredLogger.Infow("relay sequenceNumber from builder.",
		"nextSequenceNumber", txBldr.Sequence())

	// If we start to control sequence
	if nextSequenceNumber > 0 {
		txBldr.WithSequence(nextSequenceNumber)
		sugaredLogger.Infow("txBldr.WithSequence(nextSequenceNumber) passed")
	}

	log.Println("building and signing")
	// Build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(moniker, password, messages)
	if err != nil {
		sugaredLogger.Errorw("failed to sign transaction.",
			errorMessageKey, err.Error())
		return err
	}

	log.Println("built tx, now broadcasting")
	// Broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTxAsync(txBytes)
	if err != nil {
		sugaredLogger.Errorw("failed to broadcast tx to sifchain.",
			errorMessageKey, err.Error())
		return err
	}
	log.Println("Broadcasted tx without error")

	if err = cliCtx.PrintOutput(res); err != nil {
		sugaredLogger.Errorw("failed to print out result.",
			errorMessageKey, err.Error())
		return err
	}

	// start to control sequence number after first successful tx
	if nextSequenceNumber == 0 {
		setNextSequenceNumber(txBldr.Sequence() + 1)
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
