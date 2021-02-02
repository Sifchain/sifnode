package txs

// DONTCOVER

import (
	"fmt"
	"sync/atomic"

	"github.com/davecgh/go-spew/spew"

	"github.com/Sifchain/sifnode/x/ethbridge"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	nextSequenceNumber uint64 = 0
)

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(cdc *codec.Codec, moniker, password string, claims []types.EthBridgeClaim, cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder) error {
	var messages []sdk.Msg

	for _, claim := range claims {
		// Packages the claim as a Tendermint message
		msg := ethbridge.NewMsgCreateEthBridgeClaim(claim)

		err := msg.ValidateBasic()
		if err != nil {
			fmt.Println("failed to get message from claim with:", err.Error())
			continue
		} else {
			messages = append(messages, msg)
		}
	}

	// Prepare tx
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	// If we start to control sequence
	if nextSequenceNumber > 0 {
		txBldr.WithSequence(nextSequenceNumber)
	}

	spew.Dump("messages len in relayToCosmos: ", len(messages))
	// spew.Dump(moniker)
	// spew.Dump(messages)
	// spew.Dump("--------------------------------------")

	// Build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(moniker, password, messages)
	if err != nil {
		return err
	}

	// Broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		return err
	}

	if err = cliCtx.PrintOutput(res); err != nil {
		return err
	}
	// start to control sequence number after first successful tx
	if nextSequenceNumber == 0 {
		setNextSequenceNumber(txBldr.Sequence() + 1)
	} else {
		incrementNextSequenceNumber()
	}
	return nil
}

func incrementNextSequenceNumber() {
	atomic.AddUint64(&nextSequenceNumber, 1)
}

func setNextSequenceNumber(sequenceNumber uint64) {
	atomic.StoreUint64(&nextSequenceNumber, sequenceNumber)
}
