package txs

// DONTCOVER

import (
	"fmt"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/keyring"
	"github.com/Sifchain/sifnode/x/ethbridge"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// RelayToCosmos applies validator's signature to an EthBridgeClaim message containing
// information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayToCosmos(cdc *codec.Codec, moniker string, claim *types.EthBridgeClaim, cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder, singer *keyring.KeyRing) error {
	// Packages the claim as a Tendermint message
	msg := ethbridge.NewMsgCreateEthBridgeClaim(*claim)

	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	// Prepare tx
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	// txMsg, err := txBldr.BuildSignMsg([]sdk.Msg{msg})

	// Build and sign the transaction
	// txBytes, err := txBldr.BuildAndSign(moniker, keys.DefaultKeyPass, []sdk.Msg{msg})

	// How to encode msg to bytes not complete yet.
	sdkMsg := []sdk.Msg{msg}
	stdSignMsg, err := txBldr.BuildSignMsg(sdkMsg)
	if err != nil {
		fmt.Println("Message is wrong from building.")
		fmt.Println(err)
		return err
	}

	txBytes, _, err := singer.Sign(stdSignMsg.Bytes())
	if err != nil {
		fmt.Println(err)
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
	return nil
}
