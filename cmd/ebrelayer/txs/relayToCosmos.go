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
)

var (
	nextSequenceNumber uint64 = 0
)

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(cdc *codec.Codec, moniker, password string, claims []types.EthBridgeClaim, cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder) error {
	var messages []sdk.Msg
	log.Printf("ebrelayer RelayToCosmos with %d claims\n", len(claims))
	log.Printf("ebrelayer RelayToCosmos nextSequenceNumber is %d\n", nextSequenceNumber)

	for _, claim := range claims {
		// Packages the claim as a Tendermint message
		msg := ethbridge.NewMsgCreateEthBridgeClaim(claim)

		err := msg.ValidateBasic()
		if err != nil {
			log.Println("failed to get message from claim with:", err.Error())
			continue
		} else {
			messages = append(messages, msg)
		}
	}

	// Prepare tx
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		log.Println("error building tx: ", err)
		log.Println("tx buidler response on error: ", txBldr)
		return err
	}

	log.Printf("ebrelayer RelayToCosmos sequenceNumber is %d from tx builder\n", txBldr.Sequence())

	// If we start to control sequence
	if nextSequenceNumber > 0 {
		txBldr.WithSequence(nextSequenceNumber)
		log.Println("txBldr.WithSequence(nextSequenceNumber) passed")
	}

	log.Println("building and signing")
	// Build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(moniker, password, messages)
	if err != nil {
		log.Printf("ebrelayer RelayToCosmos error 1 is %s\n", err.Error())
		return err
	}

	log.Println("built tx, now broadcasting")
	// Broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTxAsync(txBytes)
	if err != nil {
		log.Printf("ebrelayer RelayToCosmos error 2 is %s\n", err.Error())
		return err
	}
	log.Println("Broadcasted tx without error")

	if err = cliCtx.PrintOutput(res); err != nil {
		log.Printf("ebrelayer RelayToCosmos error 3 is %s\n", err.Error())
		return err
	}
	log.Println("printed tx output")

	// start to control sequence number after first successful tx
	if nextSequenceNumber == 0 {
		setNextSequenceNumber(txBldr.Sequence() + 1)
	} else {
		incrementNextSequenceNumber()
	}
	log.Printf("ebrelayer RelayToCosmos nextSequenceNumber is %d after tx\n", nextSequenceNumber)
	return nil
}

func incrementNextSequenceNumber() {
	atomic.AddUint64(&nextSequenceNumber, 1)
}

func setNextSequenceNumber(sequenceNumber uint64) {
	atomic.StoreUint64(&nextSequenceNumber, sequenceNumber)
}
